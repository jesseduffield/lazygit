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

	if err := gui.focusPoint(0, gui.State.Panels.CommitFiles.SelectedLine, len(gui.State.CommitFiles), v); err != nil {
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
	file := gui.State.CommitFiles[gui.State.Panels.CommitFiles.SelectedLine]

	if err := gui.GitCommand.CheckoutFile(file.Sha, file.Name); err != nil {
		return gui.createErrorPanel(gui.g, err.Error())
	}

	return gui.refreshFiles()
}

func (gui *Gui) handleDiscardOldFileChange(g *gocui.Gui, v *gocui.View) error {
	fileName := gui.State.CommitFiles[gui.State.Panels.CommitFiles.SelectedLine].Name

	return gui.createConfirmationPanel(gui.g, v, gui.Tr.SLocalize("DiscardFileChangesTitle"), gui.Tr.SLocalize("DiscardFileChangesPrompt"), func(g *gocui.Gui, v *gocui.View) error {
		return gui.WithWaitingStatus(gui.Tr.SLocalize("RebasingStatus"), func() error {
			if err := gui.GitCommand.DiscardOldFileChanges(gui.State.Commits, gui.State.Panels.Commits.SelectedLine, fileName); err != nil {
				if err := gui.handleGenericMergeCommandResult(err); err != nil {
					return err
				}
			}

			return gui.refreshSidePanels(gui.g)
		})
	}, nil)
}

func (gui *Gui) refreshCommitFilesView() error {
	commit := gui.getSelectedCommit(gui.g)
	if commit == nil {
		return nil
	}

	files, err := gui.GitCommand.GetCommitFiles(commit.Sha)
	if err != nil {
		return gui.createErrorPanel(gui.g, err.Error())
	}

	gui.State.CommitFiles = files

	gui.refreshSelectedLine(&gui.State.Panels.CommitFiles.SelectedLine, len(gui.State.CommitFiles))

	if err := gui.renderListPanel(gui.getCommitFilesView(), gui.State.CommitFiles); err != nil {
		return err
	}

	return gui.handleCommitFileSelect(gui.g, gui.getCommitFilesView())
}

func (gui *Gui) handleOpenOldCommitFile(g *gocui.Gui, v *gocui.View) error {
	file := gui.getSelectedCommitFile(g)
	return gui.openFile(file.Name)
}
