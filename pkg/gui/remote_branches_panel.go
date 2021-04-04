package gui

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
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

func (gui *Gui) handleRemoteBranchSelect() error {
	var task updateTask
	remoteBranch := gui.getSelectedRemoteBranch()
	if remoteBranch == nil {
		task = NewRenderStringTask("No branches for this remote")
	} else {
		cmd := gui.OSCommand.ExecutableFromString(
			gui.GitCommand.GetBranchGraphCmdStr(remoteBranch.FullName()),
		)
		task = NewRunCommandTask(cmd)
	}

	return gui.refreshMainViews(refreshMainOpts{
		main: &viewUpdateOpts{
			title: "Remote Branch",
			task:  task,
		},
	})
}

func (gui *Gui) handleRemoteBranchesEscape() error {
	return gui.pushContext(gui.State.Contexts.Remotes)
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
	message := fmt.Sprintf("%s '%s'?", gui.Tr.DeleteRemoteBranchMessage, remoteBranch.FullName())

	return gui.ask(askOpts{
		title:  gui.Tr.DeleteRemoteBranch,
		prompt: message,
		handleConfirm: func() error {
			return gui.WithWaitingStatus(gui.Tr.DeletingStatus, func() error {
				err := gui.GitCommand.DeleteRemoteBranch(remoteBranch.RemoteName, remoteBranch.Name, gui.promptUserForCredential)
				gui.handleCredentialsPopup(err)

				return gui.refreshSidePanels(refreshOptions{scope: []RefreshableView{BRANCHES, REMOTES}})
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
		gui.Tr.SetUpstreamMessage,
		map[string]string{
			"checkedOut": checkedOutBranch.Name,
			"selected":   selectedBranch.FullName(),
		},
	)

	return gui.ask(askOpts{
		title:  gui.Tr.SetUpstreamTitle,
		prompt: message,
		handleConfirm: func() error {
			if err := gui.GitCommand.SetBranchUpstream(selectedBranch.RemoteName, selectedBranch.Name, checkedOutBranch.Name); err != nil {
				return err
			}

			return gui.refreshSidePanels(refreshOptions{scope: []RefreshableView{BRANCHES, REMOTES}})
		},
	})
}

func (gui *Gui) handleCreateResetToRemoteBranchMenu() error {
	selectedBranch := gui.getSelectedRemoteBranch()
	if selectedBranch == nil {
		return nil
	}

	return gui.createResetMenu(selectedBranch.FullName())
}
