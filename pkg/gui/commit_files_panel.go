package gui

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
)

func (gui *Gui) getSelectedCommitFile(g *gocui.Gui) *commands.CommitFile {
	selectedLine := gui.State.Panels.CommitFiles.SelectedLine
	if selectedLine == -1 {
		return nil
	}

	return gui.State.CommitFiles[selectedLine]
}

func (gui *Gui) handleCommitFileSelect(g *gocui.Gui, v *gocui.View) error {
	commitFile := gui.getSelectedCommitFile(g)
	if commitFile == nil {
		return gui.renderString(g, "commitFiles", gui.Tr.SLocalize("NoCommiteFiles"))
	}

	if err := gui.focusPoint(0, gui.State.Panels.CommitFiles.SelectedLine, v); err != nil {
		return err
	}
	commitText, err := gui.GitCommand.ShowCommitFile(commitFile.Sha, commitFile.Name)
	if err != nil {
		return err
	}
	return gui.renderString(g, "main", commitText)
}

func (gui *Gui) handleCommitFilesNextLine(g *gocui.Gui, v *gocui.View) error {
	panelState := gui.State.Panels.CommitFiles
	gui.changeSelectedLine(&panelState.SelectedLine, len(gui.State.CommitFiles), false)

	return gui.handleCommitFileSelect(gui.g, v)
}

func (gui *Gui) handleCommitFilesPrevLine(g *gocui.Gui, v *gocui.View) error {
	panelState := gui.State.Panels.CommitFiles
	gui.changeSelectedLine(&panelState.SelectedLine, len(gui.State.CommitFiles), true)

	return gui.handleCommitFileSelect(gui.g, v)
}

func (gui *Gui) handleSwitchToCommitsPanel(g *gocui.Gui, v *gocui.View) error {
	commitsView, err := g.View("commits")
	if err != nil {
		return err
	}
	return gui.switchFocus(g, v, commitsView)
}

func (gui *Gui) handleCheckoutCommitFile(g *gocui.Gui, v *gocui.View) error {
	commitSha := gui.State.Commits[gui.State.Panels.Commits.SelectedLine].Sha
	fileName := gui.State.CommitFiles[gui.State.Panels.CommitFiles.SelectedLine].Name

	if err := gui.GitCommand.CheckoutFile(commitSha, fileName); err != nil {
		return gui.createErrorPanel(gui.g, err.Error())
	}

	return gui.refreshFiles()
}
