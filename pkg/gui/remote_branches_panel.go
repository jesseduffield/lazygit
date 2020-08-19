package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
)

// list panel functions

func (gui *Gui) getSelectedRemoteBranch() *commands.RemoteBranch {
	selectedLine := gui.State.Panels.RemoteBranches.SelectedLine
	if selectedLine == -1 || len(gui.State.RemoteBranches) == 0 {
		return nil
	}

	return gui.State.RemoteBranches[selectedLine]
}

func (gui *Gui) handleRemoteBranchSelect() error {
	if gui.popupPanelFocused() {
		return nil
	}

	if gui.inDiffMode() {
		return gui.renderDiff()
	}

	var task updateTask
	remoteBranch := gui.getSelectedRemoteBranch()
	if remoteBranch == nil {
		task = gui.createRenderStringTask("No branches for this remote")
	} else {
		cmd := gui.OSCommand.ExecutableFromString(
			gui.GitCommand.GetBranchGraphCmdStr(remoteBranch.FullName()),
		)
		task = gui.createRunCommandTask(cmd)
	}

	return gui.refreshMain(refreshMainOpts{
		main: &viewUpdateOpts{
			title: "Remote Branch",
			task:  task,
		},
	})
}

func (gui *Gui) handleRemoteBranchesEscape(g *gocui.Gui, v *gocui.View) error {
	return gui.switchContext(gui.Contexts.Remotes.Context)
}

func (gui *Gui) renderRemoteBranchesWithSelection() error {
	branchesView := gui.getBranchesView()

	gui.refreshSelectedLine(&gui.State.Panels.RemoteBranches.SelectedLine, len(gui.State.RemoteBranches))
	displayStrings := presentation.GetRemoteBranchListDisplayStrings(gui.State.RemoteBranches, gui.State.Diff.Ref)
	gui.renderDisplayStrings(branchesView, displayStrings)

	return nil
}

func (gui *Gui) handleCheckoutRemoteBranch(g *gocui.Gui, v *gocui.View) error {
	remoteBranch := gui.getSelectedRemoteBranch()
	if remoteBranch == nil {
		return nil
	}
	if err := gui.handleCheckoutRef(remoteBranch.FullName(), handleCheckoutRefOptions{}); err != nil {
		return err
	}
	return gui.switchContext(gui.Contexts.Branches.Context)
}

func (gui *Gui) handleMergeRemoteBranch(g *gocui.Gui, v *gocui.View) error {
	selectedBranchName := gui.getSelectedRemoteBranch().FullName()
	return gui.mergeBranchIntoCheckedOutBranch(selectedBranchName)
}

func (gui *Gui) handleDeleteRemoteBranch(g *gocui.Gui, v *gocui.View) error {
	remoteBranch := gui.getSelectedRemoteBranch()
	if remoteBranch == nil {
		return nil
	}
	message := fmt.Sprintf("%s '%s'?", gui.Tr.SLocalize("DeleteRemoteBranchMessage"), remoteBranch.FullName())

	return gui.ask(askOpts{
		returnToView:       v,
		returnFocusOnClose: true,
		title:              gui.Tr.SLocalize("DeleteRemoteBranch"),
		prompt:             message,
		handleConfirm: func() error {
			return gui.WithWaitingStatus(gui.Tr.SLocalize("DeletingStatus"), func() error {
				if err := gui.GitCommand.DeleteRemoteBranch(remoteBranch.RemoteName, remoteBranch.Name); err != nil {
					return err
				}

				return gui.refreshSidePanels(refreshOptions{scope: []int{BRANCHES, REMOTES}})
			})
		},
	})
}

func (gui *Gui) handleRebaseOntoRemoteBranch(g *gocui.Gui, v *gocui.View) error {
	selectedBranchName := gui.getSelectedRemoteBranch().FullName()
	return gui.handleRebaseOntoBranch(selectedBranchName)
}

func (gui *Gui) handleSetBranchUpstream(g *gocui.Gui, v *gocui.View) error {
	selectedBranch := gui.getSelectedRemoteBranch()
	checkedOutBranch := gui.getCheckedOutBranch()

	message := gui.Tr.TemplateLocalize(
		"SetUpstreamMessage",
		Teml{
			"checkedOut": checkedOutBranch.Name,
			"selected":   selectedBranch.FullName(),
		},
	)

	return gui.ask(askOpts{
		returnToView:       v,
		returnFocusOnClose: true,
		title:              gui.Tr.SLocalize("SetUpstreamTitle"),
		prompt:             message,
		handleConfirm: func() error {
			if err := gui.GitCommand.SetBranchUpstream(selectedBranch.RemoteName, selectedBranch.Name, checkedOutBranch.Name); err != nil {
				return err
			}

			return gui.refreshSidePanels(refreshOptions{scope: []int{BRANCHES, REMOTES}})
		},
	})
}

func (gui *Gui) handleCreateResetToRemoteBranchMenu(g *gocui.Gui, v *gocui.View) error {
	selectedBranch := gui.getSelectedRemoteBranch()
	if selectedBranch == nil {
		return nil
	}

	return gui.createResetMenu(selectedBranch.FullName())
}

func (gui *Gui) handleNewBranchOffRemote(g *gocui.Gui, v *gocui.View) error {
	branch := gui.getSelectedRemoteBranch()
	if branch == nil {
		return nil
	}
	message := gui.Tr.TemplateLocalize(
		"NewBranchNameBranchOff",
		Teml{
			"branchName": branch.FullName(),
		},
	)
	return gui.prompt(v, message, branch.FullName(), func(response string) error {
		if err := gui.GitCommand.NewBranch(response, branch.FullName()); err != nil {
			return gui.surfaceError(err)
		}
		gui.State.Panels.Branches.SelectedLine = 0

		if err := gui.switchContext(gui.Contexts.Branches.Context); err != nil {
			return err
		}
		return gui.refreshSidePanels(refreshOptions{mode: ASYNC})
	})
}
