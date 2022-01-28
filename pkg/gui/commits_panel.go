package gui

import (
	"fmt"
	"sync"

	"github.com/jesseduffield/lazygit/pkg/commands/loaders"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/popup"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
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
				_ = gui.PopupHandler.Error(err)
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
		cmdObj := gui.Git.Commit.ShowCmdObj(commit.Sha, gui.State.Modes.Filtering.GetPath())
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
func (gui *Gui) refreshCommits() {
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
}

func (gui *Gui) refreshCommitsWithLimit() error {
	gui.Mutexes.BranchCommitsMutex.Lock()
	defer gui.Mutexes.BranchCommitsMutex.Unlock()

	commits, err := gui.Git.Loaders.Commits.GetCommits(
		loaders.GetCommitsOptions{
			Limit:                gui.State.Panels.Commits.LimitCommits,
			FilterPath:           gui.State.Modes.Filtering.GetPath(),
			IncludeRebaseCommits: true,
			RefName:              gui.refForLog(),
			All:                  gui.ShowWholeGitGraph,
		},
	)
	if err != nil {
		return err
	}
	gui.State.Commits = commits

	return gui.postRefreshUpdate(gui.State.Contexts.BranchCommits)
}

func (gui *Gui) refForLog() string {
	bisectInfo := gui.Git.Bisect.GetInfo()
	gui.State.BisectInfo = bisectInfo

	if !bisectInfo.Started() {
		return "HEAD"
	}

	// need to see if our bisect's current commit is reachable from our 'new' ref.
	if bisectInfo.Bisecting() && !gui.Git.Bisect.ReachableFromStart(bisectInfo) {
		return bisectInfo.GetNewSha()
	}

	return bisectInfo.GetStartSha()
}

func (gui *Gui) refreshRebaseCommits() error {
	gui.Mutexes.BranchCommitsMutex.Lock()
	defer gui.Mutexes.BranchCommitsMutex.Unlock()

	updatedCommits, err := gui.Git.Loaders.Commits.MergeRebasingCommits(gui.State.Commits)
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
		return gui.PopupHandler.ErrorMsg(gui.Tr.YouNoCommitsToSquash)
	}

	applied, err := gui.handleMidRebaseCommand("squash")
	if err != nil {
		return err
	}
	if applied {
		return nil
	}

	return gui.PopupHandler.Ask(popup.AskOpts{
		Title:  gui.Tr.Squash,
		Prompt: gui.Tr.SureSquashThisCommit,
		HandleConfirm: func() error {
			return gui.PopupHandler.WithWaitingStatus(gui.Tr.SquashingStatus, func() error {
				gui.logAction(gui.Tr.Actions.SquashCommitDown)
				err := gui.Git.Rebase.InteractiveRebase(gui.State.Commits, gui.State.Panels.Commits.SelectedLineIdx, "squash")
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
		return gui.PopupHandler.ErrorMsg(gui.Tr.YouNoCommitsToSquash)
	}

	applied, err := gui.handleMidRebaseCommand("fixup")
	if err != nil {
		return err
	}
	if applied {
		return nil
	}

	return gui.PopupHandler.Ask(popup.AskOpts{
		Title:  gui.Tr.Fixup,
		Prompt: gui.Tr.SureFixupThisCommit,
		HandleConfirm: func() error {
			return gui.PopupHandler.WithWaitingStatus(gui.Tr.FixingStatus, func() error {
				gui.logAction(gui.Tr.Actions.FixupCommit)
				err := gui.Git.Rebase.InteractiveRebase(gui.State.Commits, gui.State.Panels.Commits.SelectedLineIdx, "fixup")
				return gui.handleGenericMergeCommandResult(err)
			})
		},
	})
}

func (gui *Gui) handleRewordCommit() error {
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

	commit := gui.getSelectedLocalCommit()
	if commit == nil {
		return nil
	}

	message, err := gui.Git.Commit.GetCommitMessage(commit.Sha)
	if err != nil {
		return gui.PopupHandler.Error(err)
	}

	// TODO: use the commit message panel here
	return gui.PopupHandler.Prompt(popup.PromptOpts{
		Title:          gui.Tr.LcRewordCommit,
		InitialContent: message,
		HandleConfirm: func(response string) error {
			gui.logAction(gui.Tr.Actions.RewordCommit)
			if err := gui.Git.Rebase.RewordCommit(gui.State.Commits, gui.State.Panels.Commits.SelectedLineIdx, response); err != nil {
				return gui.PopupHandler.Error(err)
			}

			return gui.refreshSidePanels(types.RefreshOptions{Mode: types.ASYNC})
		},
	})
}

func (gui *Gui) handleRewordCommitEditor() error {
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

	gui.logAction(gui.Tr.Actions.RewordCommit)
	subProcess, err := gui.Git.Rebase.RewordCommitInEditor(gui.State.Commits, gui.State.Panels.Commits.SelectedLineIdx)
	if err != nil {
		return gui.PopupHandler.Error(err)
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
		return true, gui.PopupHandler.ErrorMsg(gui.Tr.LcRewordNotSupported)
	}

	gui.logAction("Update rebase TODO")
	gui.logCommand(
		fmt.Sprintf("Updating rebase action of commit %s to '%s'", selectedCommit.ShortSha(), action),
		false,
	)

	if err := gui.Git.Rebase.EditRebaseTodo(gui.State.Panels.Commits.SelectedLineIdx, action); err != nil {
		return false, gui.PopupHandler.Error(err)
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

	return gui.PopupHandler.Ask(popup.AskOpts{
		Title:  gui.Tr.DeleteCommitTitle,
		Prompt: gui.Tr.DeleteCommitPrompt,
		HandleConfirm: func() error {
			return gui.PopupHandler.WithWaitingStatus(gui.Tr.DeletingStatus, func() error {
				gui.logAction(gui.Tr.Actions.DropCommit)
				err := gui.Git.Rebase.InteractiveRebase(gui.State.Commits, gui.State.Panels.Commits.SelectedLineIdx, "drop")
				return gui.handleGenericMergeCommandResult(err)
			})
		},
	})
}

func (gui *Gui) handleCommitMoveDown() error {
	if ok, err := gui.validateNotInFilterMode(); err != nil || !ok {
		return err
	}

	index := gui.State.Panels.Commits.SelectedLineIdx
	selectedCommit := gui.State.Commits[index]
	if selectedCommit.Status == "rebasing" {
		if gui.State.Commits[index+1].Status != "rebasing" {
			return nil
		}

		// logging directly here because MoveTodoDown doesn't have enough information
		// to provide a useful log
		gui.logAction(gui.Tr.Actions.MoveCommitDown)
		gui.logCommand(fmt.Sprintf("Moving commit %s down", selectedCommit.ShortSha()), false)

		if err := gui.Git.Rebase.MoveTodoDown(index); err != nil {
			return gui.PopupHandler.Error(err)
		}
		gui.State.Panels.Commits.SelectedLineIdx++
		return gui.refreshRebaseCommits()
	}

	return gui.PopupHandler.WithWaitingStatus(gui.Tr.MovingStatus, func() error {
		gui.logAction(gui.Tr.Actions.MoveCommitDown)
		err := gui.Git.Rebase.MoveCommitDown(gui.State.Commits, index)
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

	selectedCommit := gui.State.Commits[index]
	if selectedCommit.Status == "rebasing" {
		// logging directly here because MoveTodoDown doesn't have enough information
		// to provide a useful log
		gui.logAction(gui.Tr.Actions.MoveCommitUp)
		gui.logCommand(
			fmt.Sprintf("Moving commit %s up", selectedCommit.ShortSha()),
			false,
		)

		if err := gui.Git.Rebase.MoveTodoDown(index - 1); err != nil {
			return gui.PopupHandler.Error(err)
		}
		gui.State.Panels.Commits.SelectedLineIdx--
		return gui.refreshRebaseCommits()
	}

	return gui.PopupHandler.WithWaitingStatus(gui.Tr.MovingStatus, func() error {
		gui.logAction(gui.Tr.Actions.MoveCommitUp)
		err := gui.Git.Rebase.MoveCommitDown(gui.State.Commits, index-1)
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

	return gui.PopupHandler.WithWaitingStatus(gui.Tr.RebasingStatus, func() error {
		gui.logAction(gui.Tr.Actions.EditCommit)
		err = gui.Git.Rebase.InteractiveRebase(gui.State.Commits, gui.State.Panels.Commits.SelectedLineIdx, "edit")
		return gui.handleGenericMergeCommandResult(err)
	})
}

func (gui *Gui) handleCommitAmendTo() error {
	if ok, err := gui.validateNotInFilterMode(); err != nil || !ok {
		return err
	}

	return gui.PopupHandler.Ask(popup.AskOpts{
		Title:  gui.Tr.AmendCommitTitle,
		Prompt: gui.Tr.AmendCommitPrompt,
		HandleConfirm: func() error {
			return gui.PopupHandler.WithWaitingStatus(gui.Tr.AmendingStatus, func() error {
				gui.logAction(gui.Tr.Actions.AmendCommit)
				err := gui.Git.Rebase.AmendTo(gui.State.Commits[gui.State.Panels.Commits.SelectedLineIdx].Sha)
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
		return gui.PopupHandler.Ask(popup.AskOpts{
			Title: gui.Tr.Actions.RevertCommit,
			Prompt: utils.ResolvePlaceholderString(
				gui.Tr.ConfirmRevertCommit,
				map[string]string{
					"selectedCommit": commit.ShortSha(),
				}),
			HandleConfirm: func() error {
				gui.logAction(gui.Tr.Actions.RevertCommit)
				if err := gui.Git.Commit.Revert(commit.Sha); err != nil {
					return gui.PopupHandler.Error(err)
				}
				return gui.afterRevertCommit()
			},
		})
	}
}

func (gui *Gui) createRevertMergeCommitMenu(commit *models.Commit) error {
	menuItems := make([]*popup.MenuItem, len(commit.Parents))
	for i, parentSha := range commit.Parents {
		i := i
		message, err := gui.Git.Commit.GetCommitMessageFirstLine(parentSha)
		if err != nil {
			return gui.PopupHandler.Error(err)
		}

		menuItems[i] = &popup.MenuItem{
			DisplayString: fmt.Sprintf("%s: %s", utils.SafeTruncate(parentSha, 8), message),
			OnPress: func() error {
				parentNumber := i + 1
				gui.logAction(gui.Tr.Actions.RevertCommit)
				if err := gui.Git.Commit.RevertMerge(commit.Sha, parentNumber); err != nil {
					return gui.PopupHandler.Error(err)
				}
				return gui.afterRevertCommit()
			},
		}
	}

	return gui.PopupHandler.Menu(popup.CreateMenuOptions{Title: gui.Tr.SelectParentCommitForMerge, Items: menuItems})
}

func (gui *Gui) afterRevertCommit() error {
	gui.State.Panels.Commits.SelectedLineIdx++
	return gui.refreshSidePanels(types.RefreshOptions{Mode: types.BLOCK_UI, Scope: []types.RefreshableView{types.COMMITS, types.BRANCHES}})
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

	return gui.PopupHandler.Ask(popup.AskOpts{
		Title:  gui.Tr.CreateFixupCommit,
		Prompt: prompt,
		HandleConfirm: func() error {
			gui.logAction(gui.Tr.Actions.CreateFixupCommit)
			if err := gui.Git.Commit.CreateFixupCommit(commit.Sha); err != nil {
				return gui.PopupHandler.Error(err)
			}

			return gui.refreshSidePanels(types.RefreshOptions{Mode: types.ASYNC})
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

	return gui.PopupHandler.Ask(popup.AskOpts{
		Title:  gui.Tr.SquashAboveCommits,
		Prompt: prompt,
		HandleConfirm: func() error {
			return gui.PopupHandler.WithWaitingStatus(gui.Tr.SquashingStatus, func() error {
				gui.logAction(gui.Tr.Actions.SquashAllAboveFixupCommits)
				err := gui.Git.Rebase.SquashAllAboveFixupCommits(commit.Sha)
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
	return gui.PopupHandler.Menu(popup.CreateMenuOptions{
		Title: gui.Tr.TagMenuTitle,
		Items: []*popup.MenuItem{
			{
				DisplayString: gui.Tr.LcLightweightTag,
				OnPress: func() error {
					return gui.handleCreateLightweightTag(commitSha)
				},
			},
			{
				DisplayString: gui.Tr.LcAnnotatedTag,
				OnPress: func() error {
					return gui.handleCreateAnnotatedTag(commitSha)
				},
			},
		},
	})
}

func (gui *Gui) afterTagCreate() error {
	gui.State.Panels.Tags.SelectedLineIdx = 0 // Set to the top
	return gui.refreshSidePanels(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.COMMITS, types.TAGS}})
}

func (gui *Gui) handleCreateAnnotatedTag(commitSha string) error {
	return gui.PopupHandler.Prompt(popup.PromptOpts{
		Title: gui.Tr.TagNameTitle,
		HandleConfirm: func(tagName string) error {
			return gui.PopupHandler.Prompt(popup.PromptOpts{
				Title: gui.Tr.TagMessageTitle,
				HandleConfirm: func(msg string) error {
					gui.logAction(gui.Tr.Actions.CreateAnnotatedTag)
					if err := gui.Git.Tag.CreateAnnotated(tagName, commitSha, msg); err != nil {
						return gui.PopupHandler.Error(err)
					}
					return gui.afterTagCreate()
				},
			})
		},
	})
}

func (gui *Gui) handleCreateLightweightTag(commitSha string) error {
	return gui.PopupHandler.Prompt(popup.PromptOpts{
		Title: gui.Tr.TagNameTitle,
		HandleConfirm: func(tagName string) error {
			gui.logAction(gui.Tr.Actions.CreateLightweightTag)
			if err := gui.Git.Tag.CreateLightweight(tagName, commitSha); err != nil {
				return gui.PopupHandler.Error(err)
			}
			return gui.afterTagCreate()
		},
	})
}

func (gui *Gui) handleCheckoutCommit() error {
	commit := gui.getSelectedLocalCommit()
	if commit == nil {
		return nil
	}

	return gui.PopupHandler.Ask(popup.AskOpts{
		Title:  gui.Tr.LcCheckoutCommit,
		Prompt: gui.Tr.SureCheckoutThisCommit,
		HandleConfirm: func() error {
			gui.logAction(gui.Tr.Actions.CheckoutCommit)
			return gui.handleCheckoutRef(commit.Sha, handleCheckoutRefOptions{})
		},
	})
}

func (gui *Gui) handleCreateCommitResetMenu() error {
	commit := gui.getSelectedLocalCommit()
	if commit == nil {
		return gui.PopupHandler.ErrorMsg(gui.Tr.NoCommitsThisBranch)
	}

	return gui.createResetMenu(commit.Sha)
}

func (gui *Gui) handleOpenSearchForCommitsPanel(string) error {
	// we usually lazyload these commits but now that we're searching we need to load them now
	if gui.State.Panels.Commits.LimitCommits {
		gui.State.Panels.Commits.LimitCommits = false
		if err := gui.refreshSidePanels(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.COMMITS}}); err != nil {
			return err
		}
	}

	return gui.handleOpenSearch("commits")
}

func (gui *Gui) handleGotoBottomForCommitsPanel() error {
	// we usually lazyload these commits but now that we're searching we need to load them now
	if gui.State.Panels.Commits.LimitCommits {
		gui.State.Panels.Commits.LimitCommits = false
		if err := gui.refreshSidePanels(types.RefreshOptions{Mode: types.SYNC, Scope: []types.RefreshableView{types.COMMITS}}); err != nil {
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

	message, err := gui.Git.Commit.GetCommitMessage(commit.Sha)
	if err != nil {
		return gui.PopupHandler.Error(err)
	}

	gui.logAction(gui.Tr.Actions.CopyCommitMessageToClipboard)
	if err := gui.OSCommand.CopyToClipboard(message); err != nil {
		return gui.PopupHandler.Error(err)
	}

	gui.raiseToast(gui.Tr.CommitMessageCopiedToClipboard)

	return nil
}

func (gui *Gui) handleOpenLogMenu() error {
	return gui.PopupHandler.Menu(popup.CreateMenuOptions{
		Title: gui.Tr.LogMenuTitle,
		Items: []*popup.MenuItem{
			{
				DisplayString: gui.Tr.ToggleShowGitGraphAll,
				OnPress: func() error {
					gui.ShowWholeGitGraph = !gui.ShowWholeGitGraph

					if gui.ShowWholeGitGraph {
						gui.State.Panels.Commits.LimitCommits = false
					}

					return gui.PopupHandler.WithWaitingStatus(gui.Tr.LcLoadingCommits, func() error {
						return gui.refreshSidePanels(types.RefreshOptions{Mode: types.SYNC, Scope: []types.RefreshableView{types.COMMITS}})
					})
				},
			},
			{
				DisplayString: gui.Tr.ShowGitGraph,
				OpensMenu:     true,
				OnPress: func() error {
					onPress := func(value string) func() error {
						return func() error {
							gui.UserConfig.Git.Log.ShowGraph = value
							gui.render()
							return nil
						}
					}
					return gui.PopupHandler.Menu(popup.CreateMenuOptions{
						Title: gui.Tr.LogMenuTitle,
						Items: []*popup.MenuItem{
							{
								DisplayString: "always",
								OnPress:       onPress("always"),
							},
							{
								DisplayString: "never",
								OnPress:       onPress("never"),
							},
							{
								DisplayString: "when maximised",
								OnPress:       onPress("when-maximised"),
							},
						},
					})
				},
			},
			{
				DisplayString: gui.Tr.SortCommits,
				OpensMenu:     true,
				OnPress: func() error {
					onPress := func(value string) func() error {
						return func() error {
							gui.UserConfig.Git.Log.Order = value
							return gui.PopupHandler.WithWaitingStatus(gui.Tr.LcLoadingCommits, func() error {
								return gui.refreshSidePanels(types.RefreshOptions{Mode: types.SYNC, Scope: []types.RefreshableView{types.COMMITS}})
							})
						}
					}

					return gui.PopupHandler.Menu(popup.CreateMenuOptions{
						Title: gui.Tr.LogMenuTitle,
						Items: []*popup.MenuItem{
							{
								DisplayString: "topological (topo-order)",
								OnPress:       onPress("topo-order"),
							},
							{
								DisplayString: "date-order",
								OnPress:       onPress("date-order"),
							},
							{
								DisplayString: "author-date-order",
								OnPress:       onPress("author-date-order"),
							},
						},
					})
				},
			},
		},
	})
}

func (gui *Gui) handleOpenCommitInBrowser() error {
	commit := gui.getSelectedLocalCommit()
	if commit == nil {
		return nil
	}

	hostingServiceMgr := gui.getHostingServiceMgr()

	url, err := hostingServiceMgr.GetCommitURL(commit.Sha)
	if err != nil {
		return gui.PopupHandler.Error(err)
	}

	gui.logAction(gui.Tr.Actions.OpenCommitInBrowser)
	if err := gui.OSCommand.OpenLink(url); err != nil {
		return gui.PopupHandler.Error(err)
	}

	return nil
}
