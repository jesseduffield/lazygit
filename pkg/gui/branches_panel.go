package gui

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
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
		task = NewRenderStringTask(gui.c.Tr.NoBranchesThisRepo)
	} else {
		cmdObj := gui.git.Branch.GetGraphCmdObj(branch.Name)

		task = NewRunPtyTask(cmdObj.GetCmd())
	}

	return gui.refreshMainViews(refreshMainOpts{
		main: &viewUpdateOpts{
			title: "Log",
			task:  task,
		},
	})
}

// specific functions

func (gui *Gui) handleBranchPress() error {
	if gui.State.Panels.Branches.SelectedLineIdx == -1 {
		return nil
	}
	if gui.State.Panels.Branches.SelectedLineIdx == 0 {
		return gui.c.ErrorMsg(gui.c.Tr.AlreadyCheckedOutBranch)
	}
	branch := gui.getSelectedBranch()
	gui.c.LogAction(gui.c.Tr.Actions.CheckoutBranch)
	return gui.helpers.refs.CheckoutRef(branch.Name, types.CheckoutRefOptions{})
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

	branchExistsOnRemote := gui.git.Remote.CheckRemoteBranchExists(branch.Name)

	if !branchExistsOnRemote {
		return gui.c.Error(errors.New(gui.c.Tr.NoBranchOnRemote))
	}

	url, err := hostingServiceMgr.GetPullRequestURL(branch.Name, "")
	if err != nil {
		return gui.c.Error(err)
	}
	gui.c.LogAction(gui.c.Tr.Actions.CopyPullRequestURL)
	if err := gui.OSCommand.CopyToClipboard(url); err != nil {
		return gui.c.Error(err)
	}

	gui.c.Toast(gui.c.Tr.PullRequestURLCopiedToClipboard)

	return nil
}

func (gui *Gui) handleGitFetch() error {
	return gui.c.WithLoaderPanel(gui.c.Tr.FetchWait, func() error {
		if err := gui.fetch(); err != nil {
			_ = gui.c.Error(err)
		}
		return gui.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
	})
}

func (gui *Gui) handleForceCheckout() error {
	branch := gui.getSelectedBranch()
	message := gui.c.Tr.SureForceCheckout
	title := gui.c.Tr.ForceCheckoutBranch

	return gui.c.Ask(types.AskOpts{
		Title:  title,
		Prompt: message,
		HandleConfirm: func() error {
			gui.c.LogAction(gui.c.Tr.Actions.ForceCheckoutBranch)
			if err := gui.git.Branch.Checkout(branch.Name, git_commands.CheckoutOptions{Force: true}); err != nil {
				_ = gui.c.Error(err)
			}
			return gui.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
		},
	})
}

func (gui *Gui) handleCheckoutByName() error {
	return gui.c.Prompt(types.PromptOpts{
		Title:               gui.c.Tr.BranchName + ":",
		FindSuggestionsFunc: gui.helpers.suggestions.GetRefsSuggestionsFunc(),
		HandleConfirm: func(response string) error {
			gui.c.LogAction("Checkout branch")
			return gui.helpers.refs.CheckoutRef(response, types.CheckoutRefOptions{
				OnRefNotFound: func(ref string) error {
					return gui.c.Ask(types.AskOpts{
						Title:  gui.c.Tr.BranchNotFoundTitle,
						Prompt: fmt.Sprintf("%s %s%s", gui.c.Tr.BranchNotFoundPrompt, ref, "?"),
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

	if err := gui.git.Branch.New(newBranchName, branch.Name); err != nil {
		return gui.c.Error(err)
	}

	gui.State.Panels.Branches.SelectedLineIdx = 0
	return gui.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
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
		return gui.c.ErrorMsg(gui.c.Tr.CantDeleteCheckOutBranch)
	}
	return gui.deleteNamedBranch(selectedBranch, force)
}

func (gui *Gui) deleteNamedBranch(selectedBranch *models.Branch, force bool) error {
	title := gui.c.Tr.DeleteBranch
	var templateStr string
	if force {
		templateStr = gui.c.Tr.ForceDeleteBranchMessage
	} else {
		templateStr = gui.c.Tr.DeleteBranchMessage
	}
	message := utils.ResolvePlaceholderString(
		templateStr,
		map[string]string{
			"selectedBranchName": selectedBranch.Name,
		},
	)

	return gui.c.Ask(types.AskOpts{
		Title:  title,
		Prompt: message,
		HandleConfirm: func() error {
			gui.c.LogAction(gui.c.Tr.Actions.DeleteBranch)
			if err := gui.git.Branch.Delete(selectedBranch.Name, force); err != nil {
				errMessage := err.Error()
				if !force && strings.Contains(errMessage, "git branch -D ") {
					return gui.deleteNamedBranch(selectedBranch, true)
				}
				return gui.c.ErrorMsg(errMessage)
			}
			return gui.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.BRANCHES}})
		},
	})
}

func (gui *Gui) mergeBranchIntoCheckedOutBranch(branchName string) error {
	if gui.git.Branch.IsHeadDetached() {
		return gui.c.ErrorMsg("Cannot merge branch in detached head state. You might have checked out a commit directly or a remote branch, in which case you should checkout the local branch you want to be on")
	}
	checkedOutBranchName := gui.getCheckedOutBranch().Name
	if checkedOutBranchName == branchName {
		return gui.c.ErrorMsg(gui.c.Tr.CantMergeBranchIntoItself)
	}
	prompt := utils.ResolvePlaceholderString(
		gui.c.Tr.ConfirmMerge,
		map[string]string{
			"checkedOutBranch": checkedOutBranchName,
			"selectedBranch":   branchName,
		},
	)

	return gui.c.Ask(types.AskOpts{
		Title:  gui.c.Tr.MergingTitle,
		Prompt: prompt,
		HandleConfirm: func() error {
			gui.c.LogAction(gui.c.Tr.Actions.Merge)
			err := gui.git.Branch.Merge(branchName, git_commands.MergeOpts{})
			return gui.helpers.rebase.CheckMergeOrRebase(err)
		},
	})
}

func (gui *Gui) handleMerge() error {
	selectedBranchName := gui.getSelectedBranch().Name
	return gui.mergeBranchIntoCheckedOutBranch(selectedBranchName)
}

func (gui *Gui) handleRebaseOntoLocalBranch() error {
	selectedBranchName := gui.getSelectedBranch().Name
	return gui.handleRebaseOntoBranch(selectedBranchName)
}

func (gui *Gui) handleRebaseOntoBranch(selectedBranchName string) error {
	checkedOutBranch := gui.getCheckedOutBranch().Name
	if selectedBranchName == checkedOutBranch {
		return gui.c.ErrorMsg(gui.c.Tr.CantRebaseOntoSelf)
	}
	prompt := utils.ResolvePlaceholderString(
		gui.c.Tr.ConfirmRebase,
		map[string]string{
			"checkedOutBranch": checkedOutBranch,
			"selectedBranch":   selectedBranchName,
		},
	)

	return gui.c.Ask(types.AskOpts{
		Title:  gui.c.Tr.RebasingTitle,
		Prompt: prompt,
		HandleConfirm: func() error {
			gui.c.LogAction(gui.c.Tr.Actions.RebaseBranch)
			err := gui.git.Rebase.RebaseBranch(selectedBranchName)
			return gui.helpers.rebase.CheckMergeOrRebase(err)
		},
	})
}

func (gui *Gui) handleFastForward() error {
	branch := gui.getSelectedBranch()
	if branch == nil || !branch.IsRealBranch() {
		return nil
	}

	if !branch.IsTrackingRemote() {
		return gui.c.ErrorMsg(gui.c.Tr.FwdNoUpstream)
	}
	if !branch.RemoteBranchStoredLocally() {
		return gui.c.ErrorMsg(gui.c.Tr.FwdNoLocalUpstream)
	}
	if branch.HasCommitsToPush() {
		return gui.c.ErrorMsg(gui.c.Tr.FwdCommitsToPush)
	}

	action := gui.c.Tr.Actions.FastForwardBranch

	message := utils.ResolvePlaceholderString(
		gui.c.Tr.Fetching,
		map[string]string{
			"from": fmt.Sprintf("%s/%s", branch.UpstreamRemote, branch.UpstreamBranch),
			"to":   branch.Name,
		},
	)

	return gui.c.WithLoaderPanel(message, func() error {
		if gui.State.Panels.Branches.SelectedLineIdx == 0 {
			gui.c.LogAction(action)

			err := gui.git.Sync.Pull(
				git_commands.PullOptions{
					RemoteName:      branch.UpstreamRemote,
					BranchName:      branch.Name,
					FastForwardOnly: true,
				},
			)
			if err != nil {
				_ = gui.c.Error(err)
			}

			return gui.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
		} else {
			gui.c.LogAction(action)
			err := gui.git.Sync.FastForward(branch.Name, branch.UpstreamRemote, branch.UpstreamBranch)
			if err != nil {
				_ = gui.c.Error(err)
			}
			_ = gui.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.BRANCHES}})
		}

		return nil
	})
}

func (gui *Gui) handleCreateResetToBranchMenu() error {
	branch := gui.getSelectedBranch()
	if branch == nil {
		return nil
	}

	return gui.helpers.refs.CreateGitResetMenu(branch.Name)
}

func (gui *Gui) handleRenameBranch() error {
	branch := gui.getSelectedBranch()
	if branch == nil || !branch.IsRealBranch() {
		return nil
	}

	promptForNewName := func() error {
		return gui.c.Prompt(types.PromptOpts{
			Title:          gui.c.Tr.NewBranchNamePrompt + " " + branch.Name + ":",
			InitialContent: branch.Name,
			HandleConfirm: func(newBranchName string) error {
				gui.c.LogAction(gui.c.Tr.Actions.RenameBranch)
				if err := gui.git.Branch.Rename(branch.Name, newBranchName); err != nil {
					return gui.c.Error(err)
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

	return gui.c.Ask(types.AskOpts{
		Title:         gui.c.Tr.LcRenameBranch,
		Prompt:        gui.c.Tr.RenameBranchWarning,
		HandleConfirm: promptForNewName,
	})
}

// sanitizedBranchName will remove all spaces in favor of a dash "-" to meet
// git's branch naming requirement.
func sanitizedBranchName(input string) string {
	return strings.Replace(input, " ", "-", -1)
}

func (gui *Gui) handleEnterBranch() error {
	branch := gui.getSelectedBranch()
	if branch == nil {
		return nil
	}

	return gui.switchToSubCommitsContext(branch.RefName())
}

func (gui *Gui) handleNewBranchOffBranch() error {
	selectedBranch := gui.getSelectedBranch()
	if selectedBranch == nil {
		return nil
	}

	return gui.helpers.refs.NewBranch(selectedBranch.RefName(), selectedBranch.RefName(), "")
}
