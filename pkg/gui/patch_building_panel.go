package gui

import (
	"github.com/jesseduffield/gocui"
)

func (gui *Gui) refreshPatchBuildingPanel() error {
	gui.State.SplitMainPanel = true

	// get diff from commit file that's currently selected
	commitFile := gui.getSelectedCommitFile(gui.g)
	if commitFile == nil {
		return gui.renderString(gui.g, "commitFiles", gui.Tr.SLocalize("NoCommiteFiles"))
	}

	diff, err := gui.GitCommand.ShowCommitFile(commitFile.Sha, commitFile.Name, true)
	if err != nil {
		return err
	}

	secondaryDiff := gui.GitCommand.PatchManager.RenderPatchForFile(commitFile.Name, true, false, true)
	if err != nil {
		return err
	}

	gui.Log.Warn(secondaryDiff)

	empty, err := gui.refreshLineByLinePanel(diff, secondaryDiff, false)
	if err != nil {
		return err
	}

	if empty {
		return gui.handleStagingEscape(gui.g, nil)
	}

	return nil
}

func (gui *Gui) handleAddSelectionToPatch(g *gocui.Gui, v *gocui.View) error {
	state := gui.State.Panels.LineByLine

	// add range of lines to those set for the file
	commitFile := gui.getSelectedCommitFile(gui.g)
	if commitFile == nil {
		return gui.renderString(gui.g, "commitFiles", gui.Tr.SLocalize("NoCommiteFiles"))
	}

	gui.GitCommand.PatchManager.AddFileLineRange(commitFile.Name, state.FirstLineIdx, state.LastLineIdx)

	if err := gui.refreshCommitFilesView(); err != nil {
		return err
	}

	if err := gui.refreshPatchBuildingPanel(); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) handleRemoveSelectionFromPatch(g *gocui.Gui, v *gocui.View) error {
	state := gui.State.Panels.LineByLine

	// add range of lines to those set for the file
	commitFile := gui.getSelectedCommitFile(gui.g)
	if commitFile == nil {
		return gui.renderString(gui.g, "commitFiles", gui.Tr.SLocalize("NoCommiteFiles"))
	}

	gui.GitCommand.PatchManager.RemoveFileLineRange(commitFile.Name, state.FirstLineIdx, state.LastLineIdx)

	if err := gui.refreshCommitFilesView(); err != nil {
		return err
	}

	if err := gui.refreshPatchBuildingPanel(); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) handleEscapePatchBuildingPanel(g *gocui.Gui, v *gocui.View) error {
	gui.State.Panels.LineByLine = nil

	return gui.switchFocus(gui.g, nil, gui.getCommitFilesView())
}
