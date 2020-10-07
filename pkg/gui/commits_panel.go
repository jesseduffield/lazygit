package gui

import (
	"sync"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// list panel functions

func (gui *Gui) getSelectedLocalCommit() *models.Commit {
	selectedLine := gui.State.Panels.Commits.SelectedLineIdx
	if selectedLine == -1 {
		return nil
	}

	return gui.State.Commits[selectedLine]
}

func (gui *Gui) handleCommitSelect() error {
	state := gui.State.Panels.Commits
	if state.SelectedLineIdx > 290 && state.LimitCommits {
		state.LimitCommits = false
		go utils.Safe(func() {
			if err := gui.refreshCommitsWithLimit(); err != nil {
				_ = gui.surfaceError(err)
			}
		})
	}

	gui.escapeLineByLinePanel()

	var task updateTask
	commit := gui.getSelectedLocalCommit()
	if commit == nil {
		task = gui.createRenderStringTask(gui.Tr.NoCommitsThisBranch)
	} else {
		cmd := gui.OSCommand.ExecutableFromString(
			gui.GitCommand.ShowCmdStr(commit.Sha, gui.State.Modes.Filtering.Path),
		)
		task = gui.createRunPtyTask(cmd)
	}

	return gui.refreshMainViews(refreshMainOpts{
		main: &viewUpdateOpts{
			title: "Patch",
			task:  task,
		},
		secondary: gui.secondaryPatchPanelUpdateOpts(),
	})
}

// during startup, the bottleneck is fetching the reflog entries. We need these
// on startup to sort the branches by recency. So we have two phases: INITIAL, and COMPLETE.
// In the initial phase we don't get any reflog commits, but we asynchronously get them
// and refresh the branches after that
func (gui *Gui) refreshReflogCommitsConsideringStartup() {
	switch gui.State.StartupStage {
	case INITIAL:
		go utils.Safe(func() {
			_ = gui.refreshReflogCommits()
			gui.refreshBranches()
			gui.State.StartupStage = COMPLETE
		})

	case COMPLETE:
		_ = gui.refreshReflogCommits()
	}
}

// whenever we change commits, we should update branches because the upstream/downstream
// counts can change. Whenever we change branches we should probably also change commits
// e.g. in the case of switching branches.
func (gui *Gui) refreshCommits() error {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go utils.Safe(func() {
		gui.refreshReflogCommitsConsideringStartup()

		gui.refreshBranches()
		wg.Done()
	})

	go utils.Safe(func() {
		_ = gui.refreshCommitsWithLimit()
		context, ok := gui.Contexts.CommitFiles.Context.GetParentContext()
		if ok && context.GetKey() == BRANCH_COMMITS_CONTEXT_KEY {
			// This makes sense when we've e.g. just amended a commit, meaning we get a new commit SHA at the same position.
			// However if we've just added a brand new commit, it pushes the list down by one and so we would end up
			// showing the contents of a different commit than the one we initially entered.
			// Ideally we would know when to refresh the commit files context and when not to,
			// or perhaps we could just pop that context off the stack whenever cycling windows.
			// For now the awkwardness remains.
			commit := gui.getSelectedLocalCommit()
			if commit != nil {
				gui.State.Panels.CommitFiles.refName = commit.RefName()
				_ = gui.refreshCommitFilesView()
			}
		}
		wg.Done()
	})

	wg.Wait()

	return nil
}

func (gui *Gui) refreshCommitsWithLimit() error {
	gui.Mutexes.BranchCommitsMutex.Lock()
	defer gui.Mutexes.BranchCommitsMutex.Unlock()

	builder := commands.NewCommitListBuilder(gui.Log, gui.GitCommand, gui.OSCommand, gui.Tr)

	commits, err := builder.GetCommits(
		commands.GetCommitsOptions{
			Limit:                gui.State.Panels.Commits.LimitCommits,
			FilterPath:           gui.State.Modes.Filtering.Path,
			IncludeRebaseCommits: true,
			RefName:              "HEAD",
		},
	)
	if err != nil {
		return err
	}
	gui.State.Commits = commits

	return gui.postRefreshUpdate(gui.Contexts.BranchCommits.Context)
}

func (gui *Gui) refreshRebaseCommits() error {
	gui.Mutexes.BranchCommitsMutex.Lock()
	defer gui.Mutexes.BranchCommitsMutex.Unlock()

	builder := commands.NewCommitListBuilder(gui.Log, gui.GitCommand, gui.OSCommand, gui.Tr)

	updatedCommits, err := builder.MergeRebasingCommits(gui.State.Commits)
	if err != nil {
		return err
	}
	gui.State.Commits = updatedCommits

	return gui.postRefreshUpdate(gui.Contexts.BranchCommits.Context)
}

// specific functions

func (gui *Gui) handleCommitSquashDown(g *gocui.Gui, v *gocui.View) error {
	if ok, err := gui.validateNotInFilterMode(); err != nil || !ok {
		return err
	}

	if len(gui.State.Commits) <= 1 {
		return gui.createErrorPanel(gui.Tr.YouNoCommitsToSquash)
	}

	applied, err := gui.handleMidRebaseCommand("squash")
	if err != nil {
		return err
	}
	if applied {
		return nil
	}

	return gui.ask(askOpts{
		title:  gui.Tr.Squash,
		prompt: gui.Tr.SureSquashThisCommit,
		handleConfirm: func() error {
			return gui.WithWaitingStatus(gui.Tr.SquashingStatus, func() error {
				err := gui.GitCommand.InteractiveRebase(gui.State.Commits, gui.State.Panels.Commits.SelectedLineIdx, "squash")
				return gui.handleGenericMergeCommandResult(err)
			})
		},
	})
}

func (gui *Gui) handleCommitFixup(g *gocui.Gui, v *gocui.View) error {
	if ok, err := gui.validateNotInFilterMode(); err != nil || !ok {
		return err
	}

	if len(gui.State.Commits) <= 1 {
		return gui.createErrorPanel(gui.Tr.YouNoCommitsToSquash)
	}

	applied, err := gui.handleMidRebaseCommand("fixup")
	if err != nil {
		return err
	}
	if applied {
		return nil
	}

	return gui.ask(askOpts{
		title:  gui.Tr.Fixup,
		prompt: gui.Tr.SureFixupThisCommit,
		handleConfirm: func() error {
			return gui.WithWaitingStatus(gui.Tr.FixingStatus, func() error {
				err := gui.GitCommand.InteractiveRebase(gui.State.Commits, gui.State.Panels.Commits.SelectedLineIdx, "fixup")
				return gui.handleGenericMergeCommandResult(err)
			})
		},
	})
}

func (gui *Gui) handleRenameCommit(g *gocui.Gui, v *gocui.View) error {
	if ok, err := gui.validateNotInFilterMode(); err != nil || !ok {
		return err
	}

	applied, err := gui.handleMidRebaseCommand("reword")
	if err != nil {
		return err
	}
	if applied {
		return nil
	}

	if gui.State.Panels.Commits.SelectedLineIdx != 0 {
		return gui.createErrorPanel(gui.Tr.OnlyRenameTopCommit)
	}

	commit := gui.getSelectedLocalCommit()
	if commit == nil {
		return nil
	}

	message, err := gui.GitCommand.GetCommitMessage(commit.Sha)
	if err != nil {
		return gui.surfaceError(err)
	}

	return gui.prompt(gui.Tr.LcRenameCommit, message, func(response string) error {
		if err := gui.GitCommand.RenameCommit(response); err != nil {
			return gui.surfaceError(err)
		}

		return gui.refreshSidePanels(refreshOptions{mode: ASYNC})
	})
}

func (gui *Gui) handleRenameCommitEditor(g *gocui.Gui, v *gocui.View) error {
	if ok, err := gui.validateNotInFilterMode(); err != nil || !ok {
		return err
	}

	applied, err := gui.handleMidRebaseCommand("reword")
	if err != nil {
		return err
	}
	if applied {
		return nil
	}

	subProcess, err := gui.GitCommand.RewordCommit(gui.State.Commits, gui.State.Panels.Commits.SelectedLineIdx)
	if err != nil {
		return gui.surfaceError(err)
	}
	if subProcess != nil {
		gui.SubProcess = subProcess
		return gui.Errors.ErrSubProcess
	}

	return nil
}

// handleMidRebaseCommand sees if the selected commit is in fact a rebasing
// commit meaning you are trying to edit the todo file rather than actually
// begin a rebase. It then updates the todo file with that action
func (gui *Gui) handleMidRebaseCommand(action string) (bool, error) {
	selectedCommit := gui.State.Commits[gui.State.Panels.Commits.SelectedLineIdx]
	if selectedCommit.Status != "rebasing" {
		return false, nil
	}

	// for now we do not support setting 'reword' because it requires an editor
	// and that means we either unconditionally wait around for the subprocess to ask for
	// our input or we set a lazygit client as the EDITOR env variable and have it
	// request us to edit the commit message when prompted.
	if action == "reword" {
		return true, gui.createErrorPanel(gui.Tr.LcRewordNotSupported)
	}

	if err := gui.GitCommand.EditRebaseTodo(gui.State.Panels.Commits.SelectedLineIdx, action); err != nil {
		return false, gui.surfaceError(err)
	}

	return true, gui.refreshRebaseCommits()
}

func (gui *Gui) handleCommitDelete(g *gocui.Gui, v *gocui.View) error {
	if ok, err := gui.validateNotInFilterMode(); err != nil || !ok {
		return err
	}

	applied, err := gui.handleMidRebaseCommand("drop")
	if err != nil {
		return err
	}
	if applied {
		return nil
	}

	return gui.ask(askOpts{
		title:  gui.Tr.DeleteCommitTitle,
		prompt: gui.Tr.DeleteCommitPrompt,
		handleConfirm: func() error {
			return gui.WithWaitingStatus(gui.Tr.DeletingStatus, func() error {
				err := gui.GitCommand.InteractiveRebase(gui.State.Commits, gui.State.Panels.Commits.SelectedLineIdx, "drop")
				return gui.handleGenericMergeCommandResult(err)
			})
		},
	})
}

func (gui *Gui) handleCommitMoveDown(g *gocui.Gui, v *gocui.View) error {
	if ok, err := gui.validateNotInFilterMode(); err != nil || !ok {
		return err
	}

	index := gui.State.Panels.Commits.SelectedLineIdx
	selectedCommit := gui.State.Commits[index]
	if selectedCommit.Status == "rebasing" {
		if gui.State.Commits[index+1].Status != "rebasing" {
			return nil
		}
		if err := gui.GitCommand.MoveTodoDown(index); err != nil {
			return gui.surfaceError(err)
		}
		gui.State.Panels.Commits.SelectedLineIdx++
		return gui.refreshRebaseCommits()
	}

	return gui.WithWaitingStatus(gui.Tr.MovingStatus, func() error {
		err := gui.GitCommand.MoveCommitDown(gui.State.Commits, index)
		if err == nil {
			gui.State.Panels.Commits.SelectedLineIdx++
		}
		return gui.handleGenericMergeCommandResult(err)
	})
}

func (gui *Gui) handleCommitMoveUp(g *gocui.Gui, v *gocui.View) error {
	if ok, err := gui.validateNotInFilterMode(); err != nil || !ok {
		return err
	}

	index := gui.State.Panels.Commits.SelectedLineIdx
	if index == 0 {
		return nil
	}
	selectedCommit := gui.State.Commits[index]
	if selectedCommit.Status == "rebasing" {
		if err := gui.GitCommand.MoveTodoDown(index - 1); err != nil {
			return gui.surfaceError(err)
		}
		gui.State.Panels.Commits.SelectedLineIdx--
		return gui.refreshRebaseCommits()
	}

	return gui.WithWaitingStatus(gui.Tr.MovingStatus, func() error {
		err := gui.GitCommand.MoveCommitDown(gui.State.Commits, index-1)
		if err == nil {
			gui.State.Panels.Commits.SelectedLineIdx--
		}
		return gui.handleGenericMergeCommandResult(err)
	})
}

func (gui *Gui) handleCommitEdit(g *gocui.Gui, v *gocui.View) error {
	if ok, err := gui.validateNotInFilterMode(); err != nil || !ok {
		return err
	}

	applied, err := gui.handleMidRebaseCommand("edit")
	if err != nil {
		return err
	}
	if applied {
		return nil
	}

	return gui.WithWaitingStatus(gui.Tr.RebasingStatus, func() error {
		err = gui.GitCommand.InteractiveRebase(gui.State.Commits, gui.State.Panels.Commits.SelectedLineIdx, "edit")
		return gui.handleGenericMergeCommandResult(err)
	})
}

func (gui *Gui) handleCommitAmendTo(g *gocui.Gui, v *gocui.View) error {
	if ok, err := gui.validateNotInFilterMode(); err != nil || !ok {
		return err
	}

	return gui.ask(askOpts{
		title:  gui.Tr.AmendCommitTitle,
		prompt: gui.Tr.AmendCommitPrompt,
		handleConfirm: func() error {
			return gui.WithWaitingStatus(gui.Tr.AmendingStatus, func() error {
				err := gui.GitCommand.AmendTo(gui.State.Commits[gui.State.Panels.Commits.SelectedLineIdx].Sha)
				return gui.handleGenericMergeCommandResult(err)
			})
		},
	})
}

func (gui *Gui) handleCommitPick(g *gocui.Gui, v *gocui.View) error {
	if ok, err := gui.validateNotInFilterMode(); err != nil || !ok {
		return err
	}

	applied, err := gui.handleMidRebaseCommand("pick")
	if err != nil {
		return err
	}
	if applied {
		return nil
	}

	// at this point we aren't actually rebasing so we will interpret this as an
	// attempt to pull. We might revoke this later after enabling configurable keybindings
	return gui.handlePullFiles(g, v)
}

func (gui *Gui) handleCommitRevert(g *gocui.Gui, v *gocui.View) error {
	if ok, err := gui.validateNotInFilterMode(); err != nil || !ok {
		return err
	}

	if err := gui.GitCommand.Revert(gui.State.Commits[gui.State.Panels.Commits.SelectedLineIdx].Sha); err != nil {
		return gui.surfaceError(err)
	}
	gui.State.Panels.Commits.SelectedLineIdx++
	return gui.refreshSidePanels(refreshOptions{mode: BLOCK_UI, scope: []int{COMMITS, BRANCHES}})
}

func (gui *Gui) handleViewCommitFiles() error {
	commit := gui.getSelectedLocalCommit()
	if commit == nil {
		return nil
	}

	return gui.switchToCommitFilesContext(commit.Sha, true, gui.Contexts.BranchCommits.Context, "commits")
}

func (gui *Gui) hasCommit(commits []*models.Commit, target string) (int, bool) {
	for idx, commit := range commits {
		if commit.Sha == target {
			return idx, true
		}
	}
	return -1, false
}

func (gui *Gui) unchooseCommit(commits []*models.Commit, i int) []*models.Commit {
	return append(commits[:i], commits[i+1:]...)
}

func (gui *Gui) handleCreateFixupCommit(g *gocui.Gui, v *gocui.View) error {
	if ok, err := gui.validateNotInFilterMode(); err != nil || !ok {
		return err
	}

	commit := gui.getSelectedLocalCommit()
	if commit == nil {
		return nil
	}

	prompt := utils.ResolvePlaceholderString(
		gui.Tr.SureCreateFixupCommit,
		map[string]string{
			"commit": commit.Sha,
		},
	)

	return gui.ask(askOpts{
		title:  gui.Tr.CreateFixupCommit,
		prompt: prompt,
		handleConfirm: func() error {
			if err := gui.GitCommand.CreateFixupCommit(commit.Sha); err != nil {
				return gui.surfaceError(err)
			}

			return gui.refreshSidePanels(refreshOptions{mode: ASYNC})
		},
	})
}

func (gui *Gui) handleSquashAllAboveFixupCommits(g *gocui.Gui, v *gocui.View) error {
	if ok, err := gui.validateNotInFilterMode(); err != nil || !ok {
		return err
	}

	commit := gui.getSelectedLocalCommit()
	if commit == nil {
		return nil
	}

	prompt := utils.ResolvePlaceholderString(
		gui.Tr.SureSquashAboveCommits,
		map[string]string{
			"commit": commit.Sha,
		},
	)

	return gui.ask(askOpts{
		title:  gui.Tr.SquashAboveCommits,
		prompt: prompt,
		handleConfirm: func() error {
			return gui.WithWaitingStatus(gui.Tr.SquashingStatus, func() error {
				err := gui.GitCommand.SquashAllAboveFixupCommits(commit.Sha)
				return gui.handleGenericMergeCommandResult(err)
			})
		},
	})
}

func (gui *Gui) handleTagCommit(g *gocui.Gui, v *gocui.View) error {
	// TODO: bring up menu asking if you want to make a lightweight or annotated tag
	// if annotated, switch to a subprocess to create the message

	commit := gui.getSelectedLocalCommit()
	if commit == nil {
		return nil
	}

	return gui.handleCreateLightweightTag(commit.Sha)
}

func (gui *Gui) handleCreateLightweightTag(commitSha string) error {
	return gui.prompt(gui.Tr.TagNameTitle, "", func(response string) error {
		if err := gui.GitCommand.CreateLightweightTag(response, commitSha); err != nil {
			return gui.surfaceError(err)
		}
		return gui.refreshSidePanels(refreshOptions{mode: ASYNC, scope: []int{COMMITS, TAGS}})
	})
}

func (gui *Gui) handleCheckoutCommit(g *gocui.Gui, v *gocui.View) error {
	commit := gui.getSelectedLocalCommit()
	if commit == nil {
		return nil
	}

	return gui.ask(askOpts{
		title:  gui.Tr.LcCheckoutCommit,
		prompt: gui.Tr.SureCheckoutThisCommit,
		handleConfirm: func() error {
			return gui.handleCheckoutRef(commit.Sha, handleCheckoutRefOptions{})
		},
	})
}

func (gui *Gui) handleCreateCommitResetMenu(g *gocui.Gui, v *gocui.View) error {
	commit := gui.getSelectedLocalCommit()
	if commit == nil {
		return gui.createErrorPanel(gui.Tr.NoCommitsThisBranch)
	}

	return gui.createResetMenu(commit.Sha)
}

func (gui *Gui) handleOpenSearchForCommitsPanel(g *gocui.Gui, v *gocui.View) error {
	// we usually lazyload these commits but now that we're searching we need to load them now
	if gui.State.Panels.Commits.LimitCommits {
		gui.State.Panels.Commits.LimitCommits = false
		if err := gui.refreshSidePanels(refreshOptions{mode: ASYNC, scope: []int{COMMITS}}); err != nil {
			return err
		}
	}

	return gui.handleOpenSearch(gui.g, v)
}

func (gui *Gui) handleGotoBottomForCommitsPanel(g *gocui.Gui, v *gocui.View) error {
	// we usually lazyload these commits but now that we're searching we need to load them now
	if gui.State.Panels.Commits.LimitCommits {
		gui.State.Panels.Commits.LimitCommits = false
		if err := gui.refreshSidePanels(refreshOptions{mode: SYNC, scope: []int{COMMITS}}); err != nil {
			return err
		}
	}

	for _, context := range gui.getListContexts() {
		if context.ViewName == "commits" {
			return context.handleGotoBottom(g, v)
		}
	}

	return nil
}
