package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// list panel functions

func (gui *Gui) remoteBranchesRenderToMain() error {
	var task updateTask
	remoteBranch := gui.State.Contexts.RemoteBranches.GetSelected()
	if remoteBranch == nil {
		task = NewRenderStringTask("No branches for this remote")
	} else {
		cmdObj := gui.git.Branch.GetGraphCmdObj(remoteBranch.FullName())
		task = NewRunCommandTask(cmdObj.GetCmd())
	}

	return gui.refreshMainViews(refreshMainOpts{
		main: &viewUpdateOpts{
			title: "Remote Branch",
			task:  task,
		},
	})
}

func (gui *Gui) handleRemoteBranchesEscape() error {
	return gui.c.PushContext(gui.State.Contexts.Remotes)
}

func (gui *Gui) handleMergeRemoteBranch() error {
	selectedBranchName := gui.State.Contexts.RemoteBranches.GetSelected().FullName()
	return gui.mergeBranchIntoCheckedOutBranch(selectedBranchName)
}

func (gui *Gui) handleDeleteRemoteBranch() error {
	remoteBranch := gui.State.Contexts.RemoteBranches.GetSelected()
	if remoteBranch == nil {
		return nil
	}
	message := fmt.Sprintf("%s '%s'?", gui.c.Tr.DeleteRemoteBranchMessage, remoteBranch.FullName())

	return gui.c.Ask(types.AskOpts{
		Title:  gui.c.Tr.DeleteRemoteBranch,
		Prompt: message,
		HandleConfirm: func() error {
			return gui.c.WithWaitingStatus(gui.c.Tr.DeletingStatus, func() error {
				gui.c.LogAction(gui.c.Tr.Actions.DeleteRemoteBranch)
				err := gui.git.Remote.DeleteRemoteBranch(remoteBranch.RemoteName, remoteBranch.Name)
				if err != nil {
					_ = gui.c.Error(err)
				}

				return gui.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.BRANCHES, types.REMOTES}})
			})
		},
	})
}

func (gui *Gui) handleRebaseOntoRemoteBranch() error {
	selectedBranchName := gui.State.Contexts.RemoteBranches.GetSelected().FullName()
	return gui.handleRebaseOntoBranch(selectedBranchName)
}

func (gui *Gui) handleSetBranchUpstream() error {
	selectedBranch := gui.State.Contexts.RemoteBranches.GetSelected()
	checkedOutBranch := gui.getCheckedOutBranch()

	message := utils.ResolvePlaceholderString(
		gui.c.Tr.SetUpstreamMessage,
		map[string]string{
			"checkedOut": checkedOutBranch.Name,
			"selected":   selectedBranch.FullName(),
		},
	)

	return gui.c.Ask(types.AskOpts{
		Title:  gui.c.Tr.SetUpstreamTitle,
		Prompt: message,
		HandleConfirm: func() error {
			gui.c.LogAction(gui.c.Tr.Actions.SetBranchUpstream)
			if err := gui.git.Branch.SetUpstream(selectedBranch.RemoteName, selectedBranch.Name, checkedOutBranch.Name); err != nil {
				return gui.c.Error(err)
			}

			return gui.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.BRANCHES, types.REMOTES}})
		},
	})
}

func (gui *Gui) handleCreateResetToRemoteBranchMenu() error {
	selectedBranch := gui.State.Contexts.RemoteBranches.GetSelected()
	if selectedBranch == nil {
		return nil
	}

	return gui.helpers.Refs.CreateGitResetMenu(selectedBranch.FullName())
}

func (gui *Gui) handleEnterRemoteBranch() error {
	selectedBranch := gui.State.Contexts.RemoteBranches.GetSelected()
	if selectedBranch == nil {
		return nil
	}

	return gui.switchToSubCommitsContext(selectedBranch.RefName())
}

func (gui *Gui) handleNewBranchOffRemoteBranch() error {
	selectedBranch := gui.State.Contexts.RemoteBranches.GetSelected()
	if selectedBranch == nil {
		return nil
	}

	// will set to the remote's branch name without the remote name
	nameSuggestion := strings.SplitAfterN(selectedBranch.RefName(), "/", 2)[1]

	return gui.helpers.Refs.NewBranch(selectedBranch.RefName(), selectedBranch.RefName(), nameSuggestion)
}
