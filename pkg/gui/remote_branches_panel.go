package gui

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// list panel functions

func (gui *Gui) getSelectedRemoteBranch() *commands.RemoteBranch {
	selectedLine := gui.State.Panels.RemoteBranches.SelectedLine
	if selectedLine == -1 || len(gui.State.RemoteBranches) == 0 {
		return nil
	}

	return gui.State.RemoteBranches[selectedLine]
}

func (gui *Gui) handleRemoteBranchesClick(g *gocui.Gui, v *gocui.View) error {
	itemCount := len(gui.State.RemoteBranches)
	handleSelect := gui.handleRemoteBranchSelect
	selectedLine := &gui.State.Panels.RemoteBranches.SelectedLine

	return gui.handleClick(v, itemCount, selectedLine, handleSelect)
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
		return gui.renderString(g, "main", "No branches for this remote")
	}

	gui.focusPoint(0, gui.State.Panels.Menu.SelectedLine, gui.State.MenuItemCount, v)
	if err := gui.focusPoint(0, gui.State.Panels.RemoteBranches.SelectedLine, len(gui.State.RemoteBranches), v); err != nil {
		return err
	}

	go func() {
		graph, err := gui.GitCommand.GetBranchGraph(fmt.Sprintf("%s/%s", remote.Name, remoteBranch.Name))
		if err != nil && strings.HasPrefix(graph, "fatal: ambiguous argument") {
			graph = gui.Tr.SLocalize("NoTrackingThisBranch")
		}
		_ = gui.renderString(g, "main", fmt.Sprintf("%s/%s\n\n%s", utils.ColoredString(remote.Name, color.FgRed), utils.ColoredString(remoteBranch.Name, color.FgGreen), graph))
	}()

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
	if err := gui.handleRemoteBranchSelect(gui.g, branchesView); err != nil {
		return err
	}

	return nil
}
