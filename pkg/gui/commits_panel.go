package gui

import (
	"fmt"
	"sync"

	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/loaders"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// after selecting the 200th commit, we'll load in all the rest
const COMMIT_THRESHOLD = 200

// list panel functions

func (gui *Gui) getSelectedLocalCommit() *models.Commit {
	selectedLine := gui.State.Panels.Commits.SelectedLineIdx
	if selectedLine == -1 || selectedLine > len(gui.State.Commits)-1 {
		return nil
	}

	return gui.State.Commits[selectedLine]
}

func (gui *Gui) onCommitFocus() error {
	state := gui.State.Panels.Commits
	if state.SelectedLineIdx > COMMIT_THRESHOLD && state.LimitCommits {
		state.LimitCommits = false
		go utils.Safe(func() {
			if err := gui.refreshCommitsWithLimit(); err != nil {
				_ = gui.surfaceError(err)
			}
		})
	}

	gui.escapeLineByLinePanel()

	return nil
}

func (gui *Gui) branchCommitsRenderToMain() error {
	var task updateTask
	commit := gui.getSelectedLocalCommit()
	if commit == nil {
		task = NewRenderStringTask(gui.Tr.NoCommitsThisBranch)
	} else {
		cmdObj := gui.GitCommand.ShowCmdObj(commit.Sha, gui.State.Modes.Filtering.GetPath())
		task = NewRunPtyTask(cmdObj.GetCmd())
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
		context, ok := gui.State.Contexts.CommitFiles.GetParentContext()
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

	loader := commands.NewCommitLoader(gui.Common, gui.GitCommand, gui.OSCommand)

	commits, err := loader.GetCommits(
		loaders.GetCommitsOptions{
			Limit:                gui.State.Panels.Commits.LimitCommits,
			FilterPath:           gui.State.Modes.Filtering.GetPath(),
			IncludeRebaseCommits: true,
			RefName:              "HEAD",
			All:                  gui.State.ShowWholeGitGraph,
		},
	)
	if err != nil {
		return err
	}
	gui.State.Commits = commits

	return gui.postRefreshUpdate(gui.State.Contexts.BranchCommits)
}

func (gui *Gui) refreshRebaseCommits() error {
	gui.Mutexes.BranchCommitsMutex.Lock()
	defer gui.Mutexes.BranchCommitsMutex.Unlock()

	loader := commands.NewCommitLoader(gui.Common, gui.GitCommand, gui.OSCommand)

	updatedCommits, err := loader.MergeRebasingCommits(gui.State.Commits)
	if err != nil {
		return err
	}
	gui.State.Commits = updatedCommits

	return gui.postRefreshUpdate(gui.State.Contexts.BranchCommits)
}

// specific functions

func (gui *Gui) handleCommitSquashDown() error {
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
				err := gui.GitCommand.WithSpan(gui.Tr.Spans.SquashCommitDown).InteractiveRebase(gui.State.Commits, gui.State.Panels.Commits.SelectedLineIdx, "squash")
				return gui.handleGenericMergeCommandResult(err)
			})
		},
	})
}

func (gui *Gui) handleCommitFixup() error {
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
				err := gui.GitCommand.WithSpan(gui.Tr.Spans.FixupCommit).InteractiveRebase(gui.State.Commits, gui.State.Panels.Commits.SelectedLineIdx, "fixup")
				return gui.handleGenericMergeCommandResult(err)
			})
		},
	})
}

func (gui *Gui) handleRenameCommit() error {
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

	return gui.prompt(promptOpts{
		title:          gui.Tr.LcRenameCommit,
		initialContent: message,
		handleConfirm: func(response string) error {
			if err := gui.GitCommand.WithSpan(gui.Tr.Spans.RewordCommit).RenameCommit(response); err != nil {
				return gui.surfaceError(err)
			}

			return gui.refreshSidePanels(refreshOptions{mode: ASYNC})
		},
	})
}

func (gui *Gui) handleRenameCommitEditor() error {
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

	subProcess, err := gui.GitCommand.WithSpan(gui.Tr.Spans.RewordCommit).RewordCommit(gui.State.Commits, gui.State.Panels.Commits.SelectedLineIdx)
	if err != nil {
		return gui.surfaceError(err)
	}
	if subProcess != nil {
		return gui.runSubprocessWithSuspenseAndRefresh(subProcess)
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

	gui.OnRunCommand(oscommands.NewCmdLogEntry(
		fmt.Sprintf("Updating rebase action of commit %s to '%s'", selectedCommit.ShortSha(), action),
		"Update rebase TODO",
		false,
	))

	if err := gui.GitCommand.EditRebaseTodo(gui.State.Panels.Commits.SelectedLineIdx, action); err != nil {
		return false, gui.surfaceError(err)
	}

	return true, gui.refreshRebaseCommits()
}

func (gui *Gui) handleCommitDelete() error {
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
				err := gui.GitCommand.WithSpan(gui.Tr.Spans.DropCommit).InteractiveRebase(gui.State.Commits, gui.State.Panels.Commits.SelectedLineIdx, "drop")
				return gui.handleGenericMergeCommandResult(err)
			})
		},
	})
}

func (gui *Gui) handleCommitMoveDown() error {
	if ok, err := gui.validateNotInFilterMode(); err != nil || !ok {
		return err
	}

	span := gui.Tr.Spans.MoveCommitDown

	index := gui.State.Panels.Commits.SelectedLineIdx
	selectedCommit := gui.State.Commits[index]
	if selectedCommit.Status == "rebasing" {
		if gui.State.Commits[index+1].Status != "rebasing" {
			return nil
		}

		// logging directly here because MoveTodoDown doesn't have enough information
		// to provide a useful log
		gui.OnRunCommand(oscommands.NewCmdLogEntry(
			fmt.Sprintf("Moving commit %s down", selectedCommit.ShortSha()),
			span,
			false,
		))

		if err := gui.GitCommand.MoveTodoDown(index); err != nil {
			return gui.surfaceError(err)
		}
		gui.State.Panels.Commits.SelectedLineIdx++
		return gui.refreshRebaseCommits()
	}

	return gui.WithWaitingStatus(gui.Tr.MovingStatus, func() error {
		err := gui.GitCommand.WithSpan(span).MoveCommitDown(gui.State.Commits, index)
		if err == nil {
			gui.State.Panels.Commits.SelectedLineIdx++
		}
		return gui.handleGenericMergeCommandResult(err)
	})
}

func (gui *Gui) handleCommitMoveUp() error {
	if ok, err := gui.validateNotInFilterMode(); err != nil || !ok {
		return err
	}

	index := gui.State.Panels.Commits.SelectedLineIdx
	if index == 0 {
		return nil
	}

	span := gui.Tr.Spans.MoveCommitUp

	selectedCommit := gui.State.Commits[index]
	if selectedCommit.Status == "rebasing" {
		// logging directly here because MoveTodoDown doesn't have enough information
		// to provide a useful log
		gui.OnRunCommand(oscommands.NewCmdLogEntry(
			fmt.Sprintf("Moving commit %s up", selectedCommit.ShortSha()),
			span,
			false,
		))

		if err := gui.GitCommand.MoveTodoDown(index - 1); err != nil {
			return gui.surfaceError(err)
		}
		gui.State.Panels.Commits.SelectedLineIdx--
		return gui.refreshRebaseCommits()
	}

	return gui.WithWaitingStatus(gui.Tr.MovingStatus, func() error {
		err := gui.GitCommand.WithSpan(span).MoveCommitDown(gui.State.Commits, index-1)
		if err == nil {
			gui.State.Panels.Commits.SelectedLineIdx--
		}
		return gui.handleGenericMergeCommandResult(err)
	})
}

func (gui *Gui) handleCommitEdit() error {
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
		err = gui.GitCommand.WithSpan(gui.Tr.Spans.EditCommit).InteractiveRebase(gui.State.Commits, gui.State.Panels.Commits.SelectedLineIdx, "edit")
		return gui.handleGenericMergeCommandResult(err)
	})
}

func (gui *Gui) handleCommitAmendTo() error {
	if ok, err := gui.validateNotInFilterMode(); err != nil || !ok {
		return err
	}

	return gui.ask(askOpts{
		title:  gui.Tr.AmendCommitTitle,
		prompt: gui.Tr.AmendCommitPrompt,
		handleConfirm: func() error {
			return gui.WithWaitingStatus(gui.Tr.AmendingStatus, func() error {
				err := gui.GitCommand.WithSpan(gui.Tr.Spans.AmendCommit).AmendTo(gui.State.Commits[gui.State.Panels.Commits.SelectedLineIdx].Sha)
				return gui.handleGenericMergeCommandResult(err)
			})
		},
	})
}

func (gui *Gui) handleCommitPick() error {
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
	return gui.handlePullFiles()
}

func (gui *Gui) handleCommitRevert() error {
	if ok, err := gui.validateNotInFilterMode(); err != nil || !ok {
		return err
	}

	commit := gui.getSelectedLocalCommit()

	if commit.IsMerge() {
		return gui.createRevertMergeCommitMenu(commit)
	} else {
		if err := gui.GitCommand.WithSpan(gui.Tr.Spans.RevertCommit).Revert(commit.Sha); err != nil {
			return gui.surfaceError(err)
		}
		return gui.afterRevertCommit()
	}
}

func (gui *Gui) createRevertMergeCommitMenu(commit *models.Commit) error {
	menuItems := make([]*menuItem, len(commit.Parents))
	for i, parentSha := range commit.Parents {
		i := i
		message, err := gui.GitCommand.GetCommitMessageFirstLine(parentSha)
		if err != nil {
			return gui.surfaceError(err)
		}

		menuItems[i] = &menuItem{
			displayString: fmt.Sprintf("%s: %s", utils.SafeTruncate(parentSha, 8), message),
			onPress: func() error {
				parentNumber := i + 1
				if err := gui.GitCommand.WithSpan(gui.Tr.Spans.RevertCommit).RevertMerge(commit.Sha, parentNumber); err != nil {
					return gui.surfaceError(err)
				}
				return gui.afterRevertCommit()
			},
		}
	}

	return gui.createMenu(gui.Tr.SelectParentCommitForMerge, menuItems, createMenuOptions{showCancel: true})
}

func (gui *Gui) afterRevertCommit() error {
	gui.State.Panels.Commits.SelectedLineIdx++
	return gui.refreshSidePanels(refreshOptions{mode: BLOCK_UI, scope: []RefreshableView{COMMITS, BRANCHES}})
}

func (gui *Gui) handleViewCommitFiles() error {
	commit := gui.getSelectedLocalCommit()
	if commit == nil {
		return nil
	}

	return gui.switchToCommitFilesContext(commit.Sha, true, gui.State.Contexts.BranchCommits, "commits")
}

func (gui *Gui) handleCreateFixupCommit() error {
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
			if err := gui.GitCommand.WithSpan(gui.Tr.Spans.CreateFixupCommit).CreateFixupCommit(commit.Sha); err != nil {
				return gui.surfaceError(err)
			}

			return gui.refreshSidePanels(refreshOptions{mode: ASYNC})
		},
	})
}

func (gui *Gui) handleSquashAllAboveFixupCommits() error {
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
				err := gui.GitCommand.WithSpan(gui.Tr.Spans.SquashAllAboveFixupCommits).SquashAllAboveFixupCommits(commit.Sha)
				return gui.handleGenericMergeCommandResult(err)
			})
		},
	})
}

func (gui *Gui) handleTagCommit() error {
	commit := gui.getSelectedLocalCommit()
	if commit == nil {
		return nil
	}

	return gui.createTagMenu(commit.Sha)
}

func (gui *Gui) createTagMenu(commitSha string) error {
	items := []*menuItem{
		{
			displayString: gui.Tr.LcLightweightTag,
			onPress: func() error {
				return gui.handleCreateLightweightTag(commitSha)
			},
		},
		{
			displayString: gui.Tr.LcAnnotatedTag,
			onPress: func() error {
				return gui.handleCreateAnnotatedTag(commitSha)
			},
		},
	}

	return gui.createMenu(gui.Tr.TagMenuTitle, items, createMenuOptions{showCancel: true})
}

func (gui *Gui) afterTagCreate(tagName string) error {
	gui.State.Panels.Tags.SelectedLineIdx = 0 // Set to the top
	return gui.refreshSidePanels(refreshOptions{mode: ASYNC, scope: []RefreshableView{COMMITS, TAGS}})
}

func (gui *Gui) handleCreateAnnotatedTag(commitSha string) error {
	return gui.prompt(promptOpts{
		title: gui.Tr.TagNameTitle,
		handleConfirm: func(tagName string) error {
			return gui.prompt(promptOpts{
				title: gui.Tr.TagMessageTitle,
				handleConfirm: func(msg string) error {
					if err := gui.GitCommand.WithSpan(gui.Tr.Spans.CreateAnnotatedTag).CreateAnnotatedTag(tagName, commitSha, msg); err != nil {
						return gui.surfaceError(err)
					}
					return gui.afterTagCreate(tagName)
				},
			})
		},
	})
}

func (gui *Gui) handleCreateLightweightTag(commitSha string) error {
	return gui.prompt(promptOpts{
		title: gui.Tr.TagNameTitle,
		handleConfirm: func(tagName string) error {
			if err := gui.GitCommand.WithSpan(gui.Tr.Spans.CreateLightweightTag).CreateLightweightTag(tagName, commitSha); err != nil {
				return gui.surfaceError(err)
			}
			return gui.afterTagCreate(tagName)
		},
	})
}

func (gui *Gui) handleCheckoutCommit() error {
	commit := gui.getSelectedLocalCommit()
	if commit == nil {
		return nil
	}

	return gui.ask(askOpts{
		title:  gui.Tr.LcCheckoutCommit,
		prompt: gui.Tr.SureCheckoutThisCommit,
		handleConfirm: func() error {
			return gui.handleCheckoutRef(commit.Sha, handleCheckoutRefOptions{span: gui.Tr.Spans.CheckoutCommit})
		},
	})
}

func (gui *Gui) handleCreateCommitResetMenu() error {
	commit := gui.getSelectedLocalCommit()
	if commit == nil {
		return gui.createErrorPanel(gui.Tr.NoCommitsThisBranch)
	}

	return gui.createResetMenu(commit.Sha)
}

func (gui *Gui) handleOpenSearchForCommitsPanel(_viewName string) error {
	// we usually lazyload these commits but now that we're searching we need to load them now
	if gui.State.Panels.Commits.LimitCommits {
		gui.State.Panels.Commits.LimitCommits = false
		if err := gui.refreshSidePanels(refreshOptions{mode: ASYNC, scope: []RefreshableView{COMMITS}}); err != nil {
			return err
		}
	}

	return gui.handleOpenSearch("commits")
}

func (gui *Gui) handleGotoBottomForCommitsPanel() error {
	// we usually lazyload these commits but now that we're searching we need to load them now
	if gui.State.Panels.Commits.LimitCommits {
		gui.State.Panels.Commits.LimitCommits = false
		if err := gui.refreshSidePanels(refreshOptions{mode: SYNC, scope: []RefreshableView{COMMITS}}); err != nil {
			return err
		}
	}

	for _, context := range gui.getListContexts() {
		if context.GetViewName() == "commits" {
			return context.handleGotoBottom()
		}
	}

	return nil
}

func (gui *Gui) handleCopySelectedCommitMessageToClipboard() error {
	commit := gui.getSelectedLocalCommit()
	if commit == nil {
		return nil
	}

	message, err := gui.GitCommand.GetCommitMessage(commit.Sha)
	if err != nil {
		return gui.surfaceError(err)
	}

	if err := gui.OSCommand.WithSpan(gui.Tr.Spans.CopyCommitMessageToClipboard).CopyToClipboard(message); err != nil {
		return gui.surfaceError(err)
	}

	gui.raiseToast(gui.Tr.CommitMessageCopiedToClipboard)

	return nil
}

func (gui *Gui) handleOpenLogMenu() error {
	return gui.createMenu(gui.Tr.LogMenuTitle, []*menuItem{
		{
			displayString: gui.Tr.ToggleShowGitGraphAll,
			onPress: func() error {
				gui.State.ShowWholeGitGraph = !gui.State.ShowWholeGitGraph

				if gui.State.ShowWholeGitGraph {
					gui.State.Panels.Commits.LimitCommits = false
				}

				return gui.WithWaitingStatus(gui.Tr.LcLoadingCommits, func() error {
					return gui.refreshSidePanels(refreshOptions{mode: SYNC, scope: []RefreshableView{COMMITS}})
				})
			},
		},
		{
			displayString: gui.Tr.ShowGitGraph,
			opensMenu:     true,
			onPress: func() error {
				onSelect := func(value string) {
					gui.UserConfig.Git.Log.ShowGraph = value
					gui.render()
				}
				return gui.createMenu(gui.Tr.LogMenuTitle, []*menuItem{
					{
						displayString: "always",
						onPress: func() error {
							onSelect("always")
							return nil
						},
					},
					{
						displayString: "never",
						onPress: func() error {
							onSelect("never")
							return nil
						},
					},
					{
						displayString: "when maximised",
						onPress: func() error {
							onSelect("when-maximised")
							return nil
						},
					},
				}, createMenuOptions{showCancel: true})
			},
		},
		{
			displayString: gui.Tr.SortCommits,
			opensMenu:     true,
			onPress: func() error {
				onSelect := func(value string) error {
					gui.UserConfig.Git.Log.Order = value
					return gui.WithWaitingStatus(gui.Tr.LcLoadingCommits, func() error {
						return gui.refreshSidePanels(refreshOptions{mode: SYNC, scope: []RefreshableView{COMMITS}})
					})
				}
				return gui.createMenu(gui.Tr.LogMenuTitle, []*menuItem{
					{
						displayString: "topological (topo-order)",
						onPress: func() error {
							return onSelect("topo-order")
						},
					},
					{
						displayString: "date-order",
						onPress: func() error {
							return onSelect("date-order")
						},
					},
					{
						displayString: "author-date-order",
						onPress: func() error {
							return onSelect("author-date-order")
						},
					},
				}, createMenuOptions{showCancel: true})
			},
		},
	}, createMenuOptions{showCancel: true})
}

func (gui *Gui) handleOpenCommitInBrowser() error {
	commit := gui.getSelectedLocalCommit()
	if commit == nil {
		return nil
	}

	hostingServiceMgr := gui.getHostingServiceMgr()

	url, err := hostingServiceMgr.GetCommitURL(commit.Sha)
	if err != nil {
		return gui.surfaceError(err)
	}

	if err := gui.GitCommand.OSCommand.OpenLink(url); err != nil {
		return gui.surfaceError(err)
	}

	gui.OnRunCommand(oscommands.NewCmdLogEntry(fmt.Sprintf(gui.Tr.OpeningCommitInBrowser, url), gui.Tr.CreatePullRequest, false))

	return nil
}
