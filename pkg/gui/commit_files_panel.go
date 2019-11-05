package gui

import (
	"github.com/go-errors/errors"
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
	commitText, err := gui.GitCommand.ShowCommitFile(commitFile.Sha, commitFile.Name, false)
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
	if ok, err := gui.validateNormalWorkingTreeState(); !ok {
		return err
	}

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
	if err := gui.refreshPatchPanel(); err != nil {
		return err
	}

	commit := gui.getSelectedCommit(gui.g)
	if commit == nil {
		return nil
	}

	files, err := gui.GitCommand.GetCommitFiles(commit.Sha, gui.GitCommand.PatchManager)
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

func (gui *Gui) handleToggleFileForPatch(g *gocui.Gui, v *gocui.View) error {
	if ok, err := gui.validateNormalWorkingTreeState(); !ok {
		return err
	}

	commitFile := gui.getSelectedCommitFile(g)
	if commitFile == nil {
		return gui.renderString(g, "commitFiles", gui.Tr.SLocalize("NoCommiteFiles"))
	}

	toggleTheFile := func() error {
		if gui.GitCommand.PatchManager == nil {
			if err := gui.createPatchManager(); err != nil {
				return err
			}
		}

		gui.GitCommand.PatchManager.ToggleFileWhole(commitFile.Name)

		return gui.refreshCommitFilesView()
	}

	if gui.GitCommand.PatchManager != nil && gui.GitCommand.PatchManager.CommitSha != commitFile.Sha {
		return gui.createConfirmationPanel(g, v, gui.Tr.SLocalize("DiscardPatch"), gui.Tr.SLocalize("DiscardPatchConfirm"), func(g *gocui.Gui, v *gocui.View) error {
			gui.GitCommand.PatchManager = nil
			return toggleTheFile()
		}, nil)
	}

	return toggleTheFile()
}

func (gui *Gui) createPatchManager() error {
	diffMap := map[string]string{}
	for _, commitFile := range gui.State.CommitFiles {
		commitText, err := gui.GitCommand.ShowCommitFile(commitFile.Sha, commitFile.Name, true)
		if err != nil {
			return err
		}
		diffMap[commitFile.Name] = commitText
	}

	commit := gui.getSelectedCommit(gui.g)
	if commit == nil {
		return errors.New("No commit selected")
	}

	gui.GitCommand.PatchManager = commands.NewPatchManager(gui.Log, gui.GitCommand.ApplyPatch, commit.Sha, diffMap)
	return nil
}

func (gui *Gui) handleEnterCommitFile(g *gocui.Gui, v *gocui.View) error {
	if ok, err := gui.validateNormalWorkingTreeState(); !ok {
		return err
	}

	commitFile := gui.getSelectedCommitFile(g)
	if commitFile == nil {
		return gui.renderString(g, "commitFiles", gui.Tr.SLocalize("NoCommiteFiles"))
	}

	enterTheFile := func() error {
		if gui.GitCommand.PatchManager == nil {
			if err := gui.createPatchManager(); err != nil {
				return err
			}
		}

		if err := gui.changeContext("main", "staging"); err != nil {
			return err
		}
		if err := gui.switchFocus(g, v, gui.getMainView()); err != nil {
			return err
		}
		return gui.refreshStagingPanel()
	}

	if gui.GitCommand.PatchManager != nil && gui.GitCommand.PatchManager.CommitSha != commitFile.Sha {
		return gui.createConfirmationPanel(g, v, gui.Tr.SLocalize("DiscardPatch"), gui.Tr.SLocalize("DiscardPatchConfirm"), func(g *gocui.Gui, v *gocui.View) error {
			gui.GitCommand.PatchManager = nil
			return enterTheFile()
		}, nil)
	}

	return enterTheFile()
}
