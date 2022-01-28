package gui

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/popup"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// list panel functions

func (gui *Gui) getSelectedBranch() *models.Branch {
	if len(gui.State.Branches) == 0 {
		return nil
	}

	selectedLine := gui.State.Panels.Branches.SelectedLineIdx
	if selectedLine == -1 {
		return nil
	}

	return gui.State.Branches[selectedLine]
}

func (gui *Gui) branchesRenderToMain() error {
	var task updateTask
	branch := gui.getSelectedBranch()
	if branch == nil {
		task = NewRenderStringTask(gui.Tr.NoBranchesThisRepo)
	} else {
		cmdObj := gui.Git.Branch.GetGraphCmdObj(branch.Name)

		task = NewRunPtyTask(cmdObj.GetCmd())
	}

	return gui.refreshMainViews(refreshMainOpts{
		main: &viewUpdateOpts{
			title: "Log",
			task:  task,
		},
	})
}

// gui.refreshStatus is called at the end of this because that's when we can
// be sure there is a state.Branches array to pick the current branch from
func (gui *Gui) refreshBranches() {
	reflogCommits := gui.State.FilteredReflogCommits
	if gui.State.Modes.Filtering.Active() {
		// in filter mode we filter our reflog commits to just those containing the path
		// however we need all the reflog entries to populate the recencies of our branches
		// which allows us to order them correctly. So if we're filtering we'll just
		// manually load all the reflog commits here
		var err error
		reflogCommits, _, err = gui.Git.Loaders.ReflogCommits.GetReflogCommits(nil, "")
		if err != nil {
			gui.Log.Error(err)
		}
	}

	branches, err := gui.Git.Loaders.Branches.Load(reflogCommits)
	if err != nil {
		_ = gui.PopupHandler.Error(err)
	}

	gui.State.Branches = branches

	if err := gui.postRefreshUpdate(gui.State.Contexts.Branches); err != nil {
		gui.Log.Error(err)
	}

	gui.refreshStatus()
}

// specific functions

func (gui *Gui) handleBranchPress() error {
	if gui.State.Panels.Branches.SelectedLineIdx == -1 {
		return nil
	}
	if gui.State.Panels.Branches.SelectedLineIdx == 0 {
		return gui.PopupHandler.ErrorMsg(gui.Tr.AlreadyCheckedOutBranch)
	}
	branch := gui.getSelectedBranch()
	gui.logAction(gui.Tr.Actions.CheckoutBranch)
	return gui.handleCheckoutRef(branch.Name, handleCheckoutRefOptions{})
}

func (gui *Gui) handleCreatePullRequestPress() error {
	branch := gui.getSelectedBranch()
	return gui.createPullRequest(branch.Name, "")
}

func (gui *Gui) handleCreatePullRequestMenu() error {
	selectedBranch := gui.getSelectedBranch()
	if selectedBranch == nil {
		return nil
	}
	checkedOutBranch := gui.getCheckedOutBranch()

	return gui.createPullRequestMenu(selectedBranch, checkedOutBranch)
}

func (gui *Gui) handleCopyPullRequestURLPress() error {
	hostingServiceMgr := gui.getHostingServiceMgr()

	branch := gui.getSelectedBranch()

	branchExistsOnRemote := gui.Git.Remote.CheckRemoteBranchExists(branch.Name)

	if !branchExistsOnRemote {
		return gui.PopupHandler.Error(errors.New(gui.Tr.NoBranchOnRemote))
	}

	url, err := hostingServiceMgr.GetPullRequestURL(branch.Name, "")
	if err != nil {
		return gui.PopupHandler.Error(err)
	}
	gui.logAction(gui.Tr.Actions.CopyPullRequestURL)
	if err := gui.OSCommand.CopyToClipboard(url); err != nil {
		return gui.PopupHandler.Error(err)
	}

	gui.raiseToast(gui.Tr.PullRequestURLCopiedToClipboard)

	return nil
}

func (gui *Gui) handleGitFetch() error {
	return gui.PopupHandler.WithLoaderPanel(gui.Tr.FetchWait, func() error {
		if err := gui.fetch(); err != nil {
			_ = gui.PopupHandler.Error(err)
		}
		return gui.refreshSidePanels(types.RefreshOptions{Mode: types.ASYNC})
	})
}

func (gui *Gui) handleForceCheckout() error {
	branch := gui.getSelectedBranch()
	message := gui.Tr.SureForceCheckout
	title := gui.Tr.ForceCheckoutBranch

	return gui.PopupHandler.Ask(popup.AskOpts{
		Title:  title,
		Prompt: message,
		HandleConfirm: func() error {
			gui.logAction(gui.Tr.Actions.ForceCheckoutBranch)
			if err := gui.Git.Branch.Checkout(branch.Name, git_commands.CheckoutOptions{Force: true}); err != nil {
				_ = gui.PopupHandler.Error(err)
			}
			return gui.refreshSidePanels(types.RefreshOptions{Mode: types.ASYNC})
		},
	})
}

type handleCheckoutRefOptions struct {
	WaitingStatus string
	EnvVars       []string
	onRefNotFound func(ref string) error
}

func (gui *Gui) handleCheckoutRef(ref string, options handleCheckoutRefOptions) error {
	waitingStatus := options.WaitingStatus
	if waitingStatus == "" {
		waitingStatus = gui.Tr.CheckingOutStatus
	}

	cmdOptions := git_commands.CheckoutOptions{Force: false, EnvVars: options.EnvVars}

	onSuccess := func() {
		gui.State.Panels.Branches.SelectedLineIdx = 0
		gui.State.Panels.Commits.SelectedLineIdx = 0
		// loading a heap of commits is slow so we limit them whenever doing a reset
		gui.State.Panels.Commits.LimitCommits = true
	}

	return gui.PopupHandler.WithWaitingStatus(waitingStatus, func() error {
		if err := gui.Git.Branch.Checkout(ref, cmdOptions); err != nil {
			// note, this will only work for english-language git commands. If we force git to use english, and the error isn't this one, then the user will receive an english command they may not understand. I'm not sure what the best solution to this is. Running the command once in english and a second time in the native language is one option

			if options.onRefNotFound != nil && strings.Contains(err.Error(), "did not match any file(s) known to git") {
				return options.onRefNotFound(ref)
			}

			if strings.Contains(err.Error(), "Please commit your changes or stash them before you switch branch") {
				// offer to autostash changes
				return gui.PopupHandler.Ask(popup.AskOpts{

					Title:  gui.Tr.AutoStashTitle,
					Prompt: gui.Tr.AutoStashPrompt,
					HandleConfirm: func() error {
						if err := gui.Git.Stash.Save(gui.Tr.StashPrefix + ref); err != nil {
							return gui.PopupHandler.Error(err)
						}
						if err := gui.Git.Branch.Checkout(ref, cmdOptions); err != nil {
							return gui.PopupHandler.Error(err)
						}

						onSuccess()
						if err := gui.Git.Stash.Pop(0); err != nil {
							if err := gui.refreshSidePanels(types.RefreshOptions{Mode: types.BLOCK_UI}); err != nil {
								return err
							}
							return gui.PopupHandler.Error(err)
						}
						return gui.refreshSidePanels(types.RefreshOptions{Mode: types.BLOCK_UI})
					},
				})
			}

			if err := gui.PopupHandler.Error(err); err != nil {
				return err
			}
		}
		onSuccess()

		return gui.refreshSidePanels(types.RefreshOptions{Mode: types.BLOCK_UI})
	})
}

func (gui *Gui) handleCheckoutByName() error {
	return gui.PopupHandler.Prompt(popup.PromptOpts{
		Title:               gui.Tr.BranchName + ":",
		FindSuggestionsFunc: gui.getRefsSuggestionsFunc(),
		HandleConfirm: func(response string) error {
			gui.logAction("Checkout branch")
			return gui.handleCheckoutRef(response, handleCheckoutRefOptions{
				onRefNotFound: func(ref string) error {
					return gui.PopupHandler.Ask(popup.AskOpts{
						Title:  gui.Tr.BranchNotFoundTitle,
						Prompt: fmt.Sprintf("%s %s%s", gui.Tr.BranchNotFoundPrompt, ref, "?"),
						HandleConfirm: func() error {
							return gui.createNewBranchWithName(ref)
						},
					})
				},
			})
		}},
	)
}

func (gui *Gui) getCheckedOutBranch() *models.Branch {
	if len(gui.State.Branches) == 0 {
		return nil
	}

	return gui.State.Branches[0]
}

func (gui *Gui) createNewBranchWithName(newBranchName string) error {
	branch := gui.getSelectedBranch()
	if branch == nil {
		return nil
	}

	if err := gui.Git.Branch.New(newBranchName, branch.Name); err != nil {
		return gui.PopupHandler.Error(err)
	}

	gui.State.Panels.Branches.SelectedLineIdx = 0
	return gui.refreshSidePanels(types.RefreshOptions{Mode: types.ASYNC})
}

func (gui *Gui) handleDeleteBranch() error {
	return gui.deleteBranch(false)
}

func (gui *Gui) deleteBranch(force bool) error {
	selectedBranch := gui.getSelectedBranch()
	if selectedBranch == nil {
		return nil
	}
	checkedOutBranch := gui.getCheckedOutBranch()
	if checkedOutBranch.Name == selectedBranch.Name {
		return gui.PopupHandler.ErrorMsg(gui.Tr.CantDeleteCheckOutBranch)
	}
	return gui.deleteNamedBranch(selectedBranch, force)
}

func (gui *Gui) deleteNamedBranch(selectedBranch *models.Branch, force bool) error {
	title := gui.Tr.DeleteBranch
	var templateStr string
	if force {
		templateStr = gui.Tr.ForceDeleteBranchMessage
	} else {
		templateStr = gui.Tr.DeleteBranchMessage
	}
	message := utils.ResolvePlaceholderString(
		templateStr,
		map[string]string{
			"selectedBranchName": selectedBranch.Name,
		},
	)

	return gui.PopupHandler.Ask(popup.AskOpts{
		Title:  title,
		Prompt: message,
		HandleConfirm: func() error {
			gui.logAction(gui.Tr.Actions.DeleteBranch)
			if err := gui.Git.Branch.Delete(selectedBranch.Name, force); err != nil {
				errMessage := err.Error()
				if !force && strings.Contains(errMessage, "git branch -D ") {
					return gui.deleteNamedBranch(selectedBranch, true)
				}
				return gui.PopupHandler.ErrorMsg(errMessage)
			}
			return gui.refreshSidePanels(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.BRANCHES}})
		},
	})
}

func (gui *Gui) mergeBranchIntoCheckedOutBranch(branchName string) error {
	if ok, err := gui.validateNotInFilterMode(); err != nil || !ok {
		return err
	}

	if gui.Git.Branch.IsHeadDetached() {
		return gui.PopupHandler.ErrorMsg("Cannot merge branch in detached head state. You might have checked out a commit directly or a remote branch, in which case you should checkout the local branch you want to be on")
	}
	checkedOutBranchName := gui.getCheckedOutBranch().Name
	if checkedOutBranchName == branchName {
		return gui.PopupHandler.ErrorMsg(gui.Tr.CantMergeBranchIntoItself)
	}
	prompt := utils.ResolvePlaceholderString(
		gui.Tr.ConfirmMerge,
		map[string]string{
			"checkedOutBranch": checkedOutBranchName,
			"selectedBranch":   branchName,
		},
	)

	return gui.PopupHandler.Ask(popup.AskOpts{
		Title:  gui.Tr.MergingTitle,
		Prompt: prompt,
		HandleConfirm: func() error {
			gui.logAction(gui.Tr.Actions.Merge)
			err := gui.Git.Branch.Merge(branchName, git_commands.MergeOpts{})
			return gui.handleGenericMergeCommandResult(err)
		},
	})
}

func (gui *Gui) handleMerge() error {
	if ok, err := gui.validateNotInFilterMode(); err != nil || !ok {
		return err
	}

	selectedBranchName := gui.getSelectedBranch().Name
	return gui.mergeBranchIntoCheckedOutBranch(selectedBranchName)
}

func (gui *Gui) handleRebaseOntoLocalBranch() error {
	selectedBranchName := gui.getSelectedBranch().Name
	return gui.handleRebaseOntoBranch(selectedBranchName)
}

func (gui *Gui) handleRebaseOntoBranch(selectedBranchName string) error {
	if ok, err := gui.validateNotInFilterMode(); err != nil || !ok {
		return err
	}

	checkedOutBranch := gui.getCheckedOutBranch().Name
	if selectedBranchName == checkedOutBranch {
		return gui.PopupHandler.ErrorMsg(gui.Tr.CantRebaseOntoSelf)
	}
	prompt := utils.ResolvePlaceholderString(
		gui.Tr.ConfirmRebase,
		map[string]string{
			"checkedOutBranch": checkedOutBranch,
			"selectedBranch":   selectedBranchName,
		},
	)

	return gui.PopupHandler.Ask(popup.AskOpts{
		Title:  gui.Tr.RebasingTitle,
		Prompt: prompt,
		HandleConfirm: func() error {
			gui.logAction(gui.Tr.Actions.RebaseBranch)
			err := gui.Git.Rebase.RebaseBranch(selectedBranchName)
			return gui.handleGenericMergeCommandResult(err)
		},
	})
}

func (gui *Gui) handleFastForward() error {
	branch := gui.getSelectedBranch()
	if branch == nil || !branch.IsRealBranch() {
		return nil
	}

	if !branch.IsTrackingRemote() {
		return gui.PopupHandler.ErrorMsg(gui.Tr.FwdNoUpstream)
	}
	if !branch.RemoteBranchStoredLocally() {
		return gui.PopupHandler.ErrorMsg(gui.Tr.FwdNoLocalUpstream)
	}
	if branch.HasCommitsToPush() {
		return gui.PopupHandler.ErrorMsg(gui.Tr.FwdCommitsToPush)
	}

	action := gui.Tr.Actions.FastForwardBranch

	message := utils.ResolvePlaceholderString(
		gui.Tr.Fetching,
		map[string]string{
			"from": fmt.Sprintf("%s/%s", branch.UpstreamRemote, branch.UpstreamBranch),
			"to":   branch.Name,
		},
	)

	return gui.PopupHandler.WithLoaderPanel(message, func() error {
		if gui.State.Panels.Branches.SelectedLineIdx == 0 {
			_ = gui.pullWithLock(PullFilesOptions{action: action, FastForwardOnly: true})
		} else {
			gui.logAction(action)
			err := gui.Git.Sync.FastForward(branch.Name, branch.UpstreamRemote, branch.UpstreamBranch)
			if err != nil {
				_ = gui.PopupHandler.Error(err)
			}
			_ = gui.refreshSidePanels(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.BRANCHES}})
		}

		return nil
	})
}

func (gui *Gui) handleCreateResetToBranchMenu() error {
	branch := gui.getSelectedBranch()
	if branch == nil {
		return nil
	}

	return gui.createResetMenu(branch.Name)
}

func (gui *Gui) handleRenameBranch() error {
	branch := gui.getSelectedBranch()
	if branch == nil || !branch.IsRealBranch() {
		return nil
	}

	promptForNewName := func() error {
		return gui.PopupHandler.Prompt(popup.PromptOpts{
			Title:          gui.Tr.NewBranchNamePrompt + " " + branch.Name + ":",
			InitialContent: branch.Name,
			HandleConfirm: func(newBranchName string) error {
				gui.logAction(gui.Tr.Actions.RenameBranch)
				if err := gui.Git.Branch.Rename(branch.Name, newBranchName); err != nil {
					return gui.PopupHandler.Error(err)
				}

				// need to find where the branch is now so that we can re-select it. That means we need to refetch the branches synchronously and then find our branch
				gui.refreshBranches()

				// now that we've got our stuff again we need to find that branch and reselect it.
				for i, newBranch := range gui.State.Branches {
					if newBranch.Name == newBranchName {
						gui.State.Panels.Branches.SetSelectedLineIdx(i)
						if err := gui.State.Contexts.Branches.HandleRender(); err != nil {
							return err
						}
					}
				}

				return nil
			},
		})
	}

	// I could do an explicit check here for whether the branch is tracking a remote branch
	// but if we've selected it we'll already know that via Pullables and Pullables.
	// Bit of a hack but I'm lazy.
	if !branch.IsTrackingRemote() {
		return promptForNewName()
	}

	return gui.PopupHandler.Ask(popup.AskOpts{
		Title:         gui.Tr.LcRenameBranch,
		Prompt:        gui.Tr.RenameBranchWarning,
		HandleConfirm: promptForNewName,
	})
}

func (gui *Gui) currentBranch() *models.Branch {
	if len(gui.State.Branches) == 0 {
		return nil
	}
	return gui.State.Branches[0]
}

func (gui *Gui) handleNewBranchOffCurrentItem() error {
	context := gui.currentSideListContext()

	item, ok := context.GetSelectedItem()
	if !ok {
		return nil
	}

	message := utils.ResolvePlaceholderString(
		gui.Tr.NewBranchNameBranchOff,
		map[string]string{
			"branchName": item.Description(),
		},
	)

	prefilledName := ""
	if context.GetKey() == REMOTE_BRANCHES_CONTEXT_KEY {
		// will set to the remote's branch name without the remote name
		prefilledName = strings.SplitAfterN(item.ID(), "/", 2)[1]
	}

	return gui.PopupHandler.Prompt(popup.PromptOpts{
		Title:          message,
		InitialContent: prefilledName,
		HandleConfirm: func(response string) error {
			gui.logAction(gui.Tr.Actions.CreateBranch)
			if err := gui.Git.Branch.New(sanitizedBranchName(response), item.ID()); err != nil {
				return err
			}

			// if we're currently in the branch commits context then the selected commit
			// is about to go to the top of the list
			if context.GetKey() == BRANCH_COMMITS_CONTEXT_KEY {
				context.GetPanelState().SetSelectedLineIdx(0)
			}

			if context.GetKey() != gui.State.Contexts.Branches.GetKey() {
				if err := gui.pushContext(gui.State.Contexts.Branches); err != nil {
					return err
				}
			}

			gui.State.Panels.Branches.SelectedLineIdx = 0

			return gui.refreshSidePanels(types.RefreshOptions{Mode: types.ASYNC})
		},
	})
}

// sanitizedBranchName will remove all spaces in favor of a dash "-" to meet
// git's branch naming requirement.
func sanitizedBranchName(input string) string {
	return strings.Replace(input, " ", "-", -1)
}
