package gui

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// list panel functions

func (gui *Gui) getSelectedRemoteBranch() *models.RemoteBranch {
	selectedLine := gui.State.Panels.RemoteBranches.SelectedLineIdx
	if selectedLine == -1 || len(gui.State.RemoteBranches) == 0 {
		return nil
	}

	return gui.State.RemoteBranches[selectedLine]
}

func (gui *Gui) remoteBranchesRenderToMain() error {
	var task updateTask
	remoteBranch := gui.getSelectedRemoteBranch()
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
	selectedBranchName := gui.getSelectedRemoteBranch().FullName()
	return gui.mergeBranchIntoCheckedOutBranch(selectedBranchName)
}

func (gui *Gui) handleDeleteRemoteBranch() error {
	remoteBranch := gui.getSelectedRemoteBranch()
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
	selectedBranchName := gui.getSelectedRemoteBranch().FullName()
	return gui.handleRebaseOntoBranch(selectedBranchName)
}

func (gui *Gui) handleSetBranchUpstream() error {
	selectedBranch := gui.getSelectedRemoteBranch()
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
	selectedBranch := gui.getSelectedRemoteBranch()
	if selectedBranch == nil {
		return nil
	}

	return gui.helpers.refs.CreateGitResetMenu(selectedBranch.FullName())
}
