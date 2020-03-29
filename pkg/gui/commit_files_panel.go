package gui

import (
	"github.com/go-errors/errors"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
)

func (gui *Gui) getSelectedCommitFile() *commands.CommitFile {
	selectedLine := gui.State.Panels.CommitFiles.SelectedLine
	if selectedLine == -1 {
		return nil
	}

	return gui.State.CommitFiles[selectedLine]
}

func (gui *Gui) handleCommitFilesClick(g *gocui.Gui, v *gocui.View) error {
	itemCount := len(gui.State.CommitFiles)
	handleSelect := gui.handleCommitFileSelect
	selectedLine := &gui.State.Panels.CommitFiles.SelectedLine

	return gui.handleClick(v, itemCount, selectedLine, handleSelect)
}

func (gui *Gui) handleCommitFileSelect(g *gocui.Gui, v *gocui.View) error {
	if gui.popupPanelFocused() {
		return nil
	}

	gui.getMainView().Title = "Patch"
	if gui.currentViewName() == "commitFiles" {
		gui.handleEscapeLineByLinePanel()
	}

	commitFile := gui.getSelectedCommitFile()
	if commitFile == nil {
		gui.renderString(g, "commitFiles", gui.Tr.SLocalize("NoCommiteFiles"))
		return nil
	}

	if err := gui.refreshSecondaryPatchPanel(); err != nil {
		return err
	}

	v.FocusPoint(0, gui.State.Panels.CommitFiles.SelectedLine)

	cmd := gui.OSCommand.ExecutableFromString(
		gui.GitCommand.ShowCommitFileCmdStr(commitFile.Sha, commitFile.Name, false),
	)
	if err := gui.newPtyTask("main", cmd); err != nil {
		gui.Log.Error(err)
	}

	return nil
}

func (gui *Gui) handleSwitchToCommitsPanel(g *gocui.Gui, v *gocui.View) error {
	return gui.switchFocus(g, v, gui.getCommitsView())
}

func (gui *Gui) handleCheckoutCommitFile(g *gocui.Gui, v *gocui.View) error {
	file := gui.State.CommitFiles[gui.State.Panels.CommitFiles.SelectedLine]

	if err := gui.GitCommand.CheckoutFile(file.Sha, file.Name); err != nil {
		return gui.surfaceError(err)
	}

	return gui.refreshSidePanels(refreshOptions{mode: ASYNC})
}

func (gui *Gui) handleDiscardOldFileChange(g *gocui.Gui, v *gocui.View) error {
	if ok, err := gui.validateNormalWorkingTreeState(); !ok {
		return err
	}

	fileName := gui.State.CommitFiles[gui.State.Panels.CommitFiles.SelectedLine].Name

	return gui.createConfirmationPanel(gui.g, v, true, gui.Tr.SLocalize("DiscardFileChangesTitle"), gui.Tr.SLocalize("DiscardFileChangesPrompt"), func(g *gocui.Gui, v *gocui.View) error {
		return gui.WithWaitingStatus(gui.Tr.SLocalize("RebasingStatus"), func() error {
			if err := gui.GitCommand.DiscardOldFileChanges(gui.State.Commits, gui.State.Panels.Commits.SelectedLine, fileName); err != nil {
				if err := gui.handleGenericMergeCommandResult(err); err != nil {
					return err
				}
			}

			return gui.refreshSidePanels(refreshOptions{mode: BLOCK_UI})
		})
	}, nil)
}

func (gui *Gui) refreshCommitFilesView() error {
	if err := gui.refreshSecondaryPatchPanel(); err != nil {
		return err
	}

	if err := gui.refreshPatchBuildingPanel(-1); err != nil {
		return err
	}

	commit := gui.getSelectedCommit()
	if commit == nil {
		return nil
	}

	files, err := gui.GitCommand.GetCommitFiles(commit.Sha, gui.GitCommand.PatchManager)
	if err != nil {
		return gui.surfaceError(err)
	}
	gui.State.CommitFiles = files

	gui.refreshSelectedLine(&gui.State.Panels.CommitFiles.SelectedLine, len(gui.State.CommitFiles))

	commitsFileView := gui.getCommitFilesView()
	displayStrings := presentation.GetCommitFileListDisplayStrings(gui.State.CommitFiles, gui.State.Diff.Ref)
	gui.renderDisplayStrings(commitsFileView, displayStrings)

	return gui.handleCommitFileSelect(gui.g, commitsFileView)
}

func (gui *Gui) handleOpenOldCommitFile(g *gocui.Gui, v *gocui.View) error {
	file := gui.getSelectedCommitFile()
	return gui.openFile(file.Name)
}

func (gui *Gui) handleToggleFileForPatch(g *gocui.Gui, v *gocui.View) error {
	if ok, err := gui.validateNormalWorkingTreeState(); !ok {
		return err
	}

	commitFile := gui.getSelectedCommitFile()
	if commitFile == nil {
		gui.renderString(g, "commitFiles", gui.Tr.SLocalize("NoCommiteFiles"))
		return nil
	}

	toggleTheFile := func() error {
		if !gui.GitCommand.PatchManager.CommitSelected() {
			if err := gui.startPatchManager(); err != nil {
				return err
			}
		}

		gui.GitCommand.PatchManager.ToggleFileWhole(commitFile.Name)

		return gui.refreshCommitFilesView()
	}

	if gui.GitCommand.PatchManager.CommitSelected() && gui.GitCommand.PatchManager.CommitSha != commitFile.Sha {
		return gui.createConfirmationPanel(g, v, true, gui.Tr.SLocalize("DiscardPatch"), gui.Tr.SLocalize("DiscardPatchConfirm"), func(g *gocui.Gui, v *gocui.View) error {
			gui.GitCommand.PatchManager.Reset()
			return toggleTheFile()
		}, nil)
	}

	return toggleTheFile()
}

func (gui *Gui) startPatchManager() error {
	diffMap := map[string]string{}
	for _, commitFile := range gui.State.CommitFiles {
		commitText, err := gui.GitCommand.ShowCommitFile(commitFile.Sha, commitFile.Name, true)
		if err != nil {
			return err
		}
		diffMap[commitFile.Name] = commitText
	}

	commit := gui.getSelectedCommit()
	if commit == nil {
		return errors.New("No commit selected")
	}

	gui.GitCommand.PatchManager.Start(commit.Sha, diffMap)
	return nil
}

func (gui *Gui) handleEnterCommitFile(g *gocui.Gui, v *gocui.View) error {
	return gui.enterCommitFile(-1)
}

func (gui *Gui) enterCommitFile(selectedLineIdx int) error {
	if ok, err := gui.validateNormalWorkingTreeState(); !ok {
		return err
	}

	commitFile := gui.getSelectedCommitFile()
	if commitFile == nil {
		gui.renderString(gui.g, "commitFiles", gui.Tr.SLocalize("NoCommiteFiles"))
		return nil
	}

	enterTheFile := func(selectedLineIdx int) error {
		if !gui.GitCommand.PatchManager.CommitSelected() {
			if err := gui.startPatchManager(); err != nil {
				return err
			}
		}

		gui.changeMainViewsContext("patch-building")
		if err := gui.switchFocus(gui.g, gui.getCommitFilesView(), gui.getMainView()); err != nil {
			return err
		}
		return gui.refreshPatchBuildingPanel(selectedLineIdx)
	}

	if gui.GitCommand.PatchManager.CommitSelected() && gui.GitCommand.PatchManager.CommitSha != commitFile.Sha {
		return gui.createConfirmationPanel(gui.g, gui.getCommitFilesView(), false, gui.Tr.SLocalize("DiscardPatch"), gui.Tr.SLocalize("DiscardPatchConfirm"), func(g *gocui.Gui, v *gocui.View) error {
			gui.GitCommand.PatchManager.Reset()
			return enterTheFile(selectedLineIdx)
		}, func(g *gocui.Gui, v *gocui.View) error {
			return gui.switchFocus(gui.g, nil, gui.getCommitFilesView())
		})
	}

	return enterTheFile(selectedLineIdx)
}

func (gui *Gui) onCommitFilesPanelSearchSelect(selectedLine int) error {
	gui.State.Panels.CommitFiles.SelectedLine = selectedLine
	return gui.handleCommitFileSelect(gui.g, gui.getCommitFilesView())
}
