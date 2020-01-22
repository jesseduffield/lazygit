package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
)

// list panel functions

func (gui *Gui) getSelectedRemoteBranch() *commands.RemoteBranch {
	selectedLine := gui.State.Panels.RemoteBranches.SelectedLine
	if selectedLine == -1 || len(gui.State.RemoteBranches) == 0 {
		return nil
	}

	return gui.State.RemoteBranches[selectedLine]
}

func (gui *Gui) handleRemoteBranchSelect(g *gocui.Gui, v *gocui.View) error {
	if gui.popupPanelFocused() {
		return nil
	}

	gui.State.SplitMainPanel = false

	if _, err := gui.g.SetCurrentView(v.Name()); err != nil {
		return err
	}

	gui.getMainView().Title = "Remote Branch"

	remote := gui.getSelectedRemote()
	remoteBranch := gui.getSelectedRemoteBranch()
	if remoteBranch == nil {
		return gui.newStringTask("main", "No branches for this remote")
	}

	gui.focusPoint(0, gui.State.Panels.Menu.SelectedLine, gui.State.MenuItemCount, v)
	if err := gui.focusPoint(0, gui.State.Panels.RemoteBranches.SelectedLine, len(gui.State.RemoteBranches), v); err != nil {
		return err
	}

	branchName := fmt.Sprintf("%s/%s", remote.Name, remoteBranch.Name)

	cmd := gui.OSCommand.ExecutableFromString(
		gui.GitCommand.GetBranchGraphCmdStr(branchName),
	)
	if err := gui.newCmdTask("main", cmd); err != nil {
		gui.Log.Error(err)
	}

	return nil
}

func (gui *Gui) handleRemoteBranchesEscape(g *gocui.Gui, v *gocui.View) error {
	return gui.switchBranchesPanelContext("remotes")
}

func (gui *Gui) renderRemoteBranchesWithSelection() error {
	branchesView := gui.getBranchesView()

	gui.refreshSelectedLine(&gui.State.Panels.RemoteBranches.SelectedLine, len(gui.State.RemoteBranches))
	if err := gui.renderListPanel(branchesView, gui.State.RemoteBranches); err != nil {
		return err
	}
	if gui.g.CurrentView() == branchesView && branchesView.Context == "remote-branches" {
		if err := gui.handleRemoteBranchSelect(gui.g, branchesView); err != nil {
			return err
		}
	}

	return nil
}

func (gui *Gui) handleCheckoutRemoteBranch(g *gocui.Gui, v *gocui.View) error {
	remoteBranch := gui.getSelectedRemoteBranch()
	if remoteBranch == nil {
		return nil
	}
	if err := gui.handleCheckoutRef(remoteBranch.RemoteName + "/" + remoteBranch.Name); err != nil {
		return err
	}
	return gui.switchBranchesPanelContext("local-branches")
}

func (gui *Gui) handleMergeRemoteBranch(g *gocui.Gui, v *gocui.View) error {
	selectedBranchName := gui.getSelectedRemoteBranch().Name
	return gui.mergeBranchIntoCheckedOutBranch(selectedBranchName)
}

func (gui *Gui) handleDeleteRemoteBranch(g *gocui.Gui, v *gocui.View) error {
	remoteBranch := gui.getSelectedRemoteBranch()
	if remoteBranch == nil {
		return nil
	}
	message := fmt.Sprintf("%s '%s/%s'?", gui.Tr.SLocalize("DeleteRemoteBranchMessage"), remoteBranch.RemoteName, remoteBranch.Name)
	return gui.createConfirmationPanel(g, v, true, gui.Tr.SLocalize("DeleteRemoteBranch"), message, func(*gocui.Gui, *gocui.View) error {
		return gui.WithWaitingStatus(gui.Tr.SLocalize("DeletingStatus"), func() error {
			if err := gui.GitCommand.DeleteRemoteBranch(remoteBranch.RemoteName, remoteBranch.Name); err != nil {
				return err
			}

			return gui.refreshRemotes()
		})
	}, nil)
}

func (gui *Gui) handleRebaseOntoRemoteBranch(g *gocui.Gui, v *gocui.View) error {
	selectedBranchName := gui.getSelectedRemoteBranch().Name
	return gui.handleRebaseOntoBranch(selectedBranchName)
}

func (gui *Gui) handleSetBranchUpstream(g *gocui.Gui, v *gocui.View) error {
	selectedBranch := gui.getSelectedRemoteBranch()
	checkedOutBranch := gui.getCheckedOutBranch()

	message := gui.Tr.TemplateLocalize(
		"SetUpstreamMessage",
		Teml{
			"checkedOut": checkedOutBranch.Name,
			"selected":   selectedBranch.RemoteName + "/" + selectedBranch.Name,
		},
	)

	return gui.createConfirmationPanel(g, v, true, gui.Tr.SLocalize("SetUpstreamTitle"), message, func(*gocui.Gui, *gocui.View) error {
		if err := gui.GitCommand.SetBranchUpstream(selectedBranch.RemoteName, selectedBranch.Name, checkedOutBranch.Name); err != nil {
			return err
		}

		return gui.refreshSidePanels(gui.g)
	}, nil)
}
