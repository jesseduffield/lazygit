package gui

import (
	"sync"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
)

// list panel functions

func (gui *Gui) getSelectedLocalCommit() *commands.Commit {
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
		go func() {
			if err := gui.refreshCommitsWithLimit(); err != nil {
				_ = gui.surfaceError(err)
			}
		}()
	}

	gui.handleEscapeLineByLinePanel()

	var task updateTask
	commit := gui.getSelectedLocalCommit()
	if commit == nil {
		task = gui.createRenderStringTask(gui.Tr.SLocalize("NoCommitsThisBranch"))
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
		go func() {
			_ = gui.refreshReflogCommits()
			gui.refreshBranches()
			gui.State.StartupStage = COMPLETE
		}()

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

	go func() {
		gui.refreshReflogCommitsConsideringStartup()

		gui.refreshBranches()
		wg.Done()
	}()

	go func() {
		_ = gui.refreshCommitsWithLimit()
		if gui.g.CurrentView() == gui.getCommitFilesView() || (gui.currentContext().GetKey() == gui.Contexts.PatchBuilding.Context.GetKey()) {
			_ = gui.refreshCommitFilesView()
		}
		wg.Done()
	}()

	wg.Wait()

	return nil
}

func (gui *Gui) refreshCommitsWithLimit() error {
	builder := commands.NewCommitListBuilder(gui.Log, gui.GitCommand, gui.OSCommand, gui.Tr, gui.State.Modes.CherryPicking.CherryPickedCommits)

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

// specific functions

func (gui *Gui) handleCommitSquashDown(g *gocui.Gui, v *gocui.View) error {
	if ok, err := gui.validateNotInFilterMode(); err != nil || !ok {
		return err
	}

	if len(gui.State.Commits) <= 1 {
		return gui.createErrorPanel(gui.Tr.SLocalize("YouNoCommitsToSquash"))
	}

	applied, err := gui.handleMidRebaseCommand("squash")
	if err != nil {
		return err
	}
	if applied {
		return nil
	}

	return gui.ask(askOpts{
		title:  gui.Tr.SLocalize("Squash"),
		prompt: gui.Tr.SLocalize("SureSquashThisCommit"),
		handleConfirm: func() error {
			return gui.WithWaitingStatus(gui.Tr.SLocalize("SquashingStatus"), func() error {
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
		return gui.createErrorPanel(gui.Tr.SLocalize("YouNoCommitsToSquash"))
	}

	applied, err := gui.handleMidRebaseCommand("fixup")
	if err != nil {
		return err
	}
	if applied {
		return nil
	}

	return gui.ask(askOpts{
		title:  gui.Tr.SLocalize("Fixup"),
		prompt: gui.Tr.SLocalize("SureFixupThisCommit"),
		handleConfirm: func() error {
			return gui.WithWaitingStatus(gui.Tr.SLocalize("FixingStatus"), func() error {
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
		return gui.createErrorPanel(gui.Tr.SLocalize("OnlyRenameTopCommit"))
	}

	commit := gui.getSelectedLocalCommit()
	if commit == nil {
		return nil
	}

	message, err := gui.GitCommand.GetCommitMessage(commit.Sha)
	if err != nil {
		return gui.surfaceError(err)
	}

	return gui.prompt(gui.Tr.SLocalize("renameCommit"), message, func(response string) error {
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
		return true, gui.createErrorPanel(gui.Tr.SLocalize("rewordNotSupported"))
	}

	if err := gui.GitCommand.EditRebaseTodo(gui.State.Panels.Commits.SelectedLineIdx, action); err != nil {
		return false, gui.surfaceError(err)
	}
	// TODO: consider doing this in a way that is less expensive. We don't actually
	// need to reload all the commits, just the TODO commits.
	return true, gui.refreshSidePanels(refreshOptions{scope: []int{COMMITS}})
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
		title:  gui.Tr.SLocalize("DeleteCommitTitle"),
		prompt: gui.Tr.SLocalize("DeleteCommitPrompt"),
		handleConfirm: func() error {
			return gui.WithWaitingStatus(gui.Tr.SLocalize("DeletingStatus"), func() error {
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
		return gui.refreshSidePanels(refreshOptions{mode: BLOCK_UI, scope: []int{COMMITS, BRANCHES}})
	}

	return gui.WithWaitingStatus(gui.Tr.SLocalize("MovingStatus"), func() error {
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
		return gui.refreshSidePanels(refreshOptions{mode: BLOCK_UI, scope: []int{COMMITS, BRANCHES}})
	}

	return gui.WithWaitingStatus(gui.Tr.SLocalize("MovingStatus"), func() error {
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

	return gui.WithWaitingStatus(gui.Tr.SLocalize("RebasingStatus"), func() error {
		err = gui.GitCommand.InteractiveRebase(gui.State.Commits, gui.State.Panels.Commits.SelectedLineIdx, "edit")
		return gui.handleGenericMergeCommandResult(err)
	})
}

func (gui *Gui) handleCommitAmendTo(g *gocui.Gui, v *gocui.View) error {
	if ok, err := gui.validateNotInFilterMode(); err != nil || !ok {
		return err
	}

	return gui.ask(askOpts{
		title:  gui.Tr.SLocalize("AmendCommitTitle"),
		prompt: gui.Tr.SLocalize("AmendCommitPrompt"),
		handleConfirm: func() error {
			return gui.WithWaitingStatus(gui.Tr.SLocalize("AmendingStatus"), func() error {
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

func (gui *Gui) hasCommit(commits []*commands.Commit, target string) (int, bool) {
	for idx, commit := range commits {
		if commit.Sha == target {
			return idx, true
		}
	}
	return -1, false
}

func (gui *Gui) unchooseCommit(commits []*commands.Commit, i int) []*commands.Commit {
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

	return gui.ask(askOpts{
		title: gui.Tr.SLocalize("CreateFixupCommit"),
		prompt: gui.Tr.TemplateLocalize(
			"SureCreateFixupCommit",
			Teml{
				"commit": commit.Sha,
			},
		),
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

	return gui.ask(askOpts{
		title: gui.Tr.SLocalize("SquashAboveCommits"),
		prompt: gui.Tr.TemplateLocalize(
			"SureSquashAboveCommits",
			Teml{
				"commit": commit.Sha,
			},
		),
		handleConfirm: func() error {
			return gui.WithWaitingStatus(gui.Tr.SLocalize("SquashingStatus"), func() error {
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
	return gui.prompt(gui.Tr.SLocalize("TagNameTitle"), "", func(response string) error {
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
		title:  gui.Tr.SLocalize("checkoutCommit"),
		prompt: gui.Tr.SLocalize("SureCheckoutThisCommit"),
		handleConfirm: func() error {
			return gui.handleCheckoutRef(commit.Sha, handleCheckoutRefOptions{})
		},
	})
}

func (gui *Gui) handleCreateCommitResetMenu(g *gocui.Gui, v *gocui.View) error {
	commit := gui.getSelectedLocalCommit()
	if commit == nil {
		return gui.createErrorPanel(gui.Tr.SLocalize("NoCommitsThisBranch"))
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
