package gui

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
)

// these represent what select mode we're in
const (
	LINE = iota
	RANGE
	HUNK
)

func (gui *Gui) refreshStagingPanel() error {
	state := gui.State.Panels.Staging

	// file, err := gui.getSelectedFile(gui.g)
	// if err != nil {
	// 	if err != gui.Errors.ErrNoFiles {
	// 		return err
	// 	}
	// 	return gui.handleStagingEscape(gui.g, nil)
	// }

	gui.State.SplitMainPanel = true

	secondaryFocused := false
	if state != nil {
		secondaryFocused = state.SecondaryFocused
	}

	// if !file.HasUnstagedChanges && !file.HasStagedChanges {
	// 	return gui.handleStagingEscape(gui.g, nil)
	// }

	// if (secondaryFocused && !file.HasStagedChanges) || (!secondaryFocused && !file.HasUnstagedChanges) {
	// 	secondaryFocused = !secondaryFocused
	// }

	// getDiffs := func() (string, string) {
	// 	// note for custom diffs, we'll need to send a flag here saying not to use the custom diff
	// 	diff := gui.GitCommand.Diff(file, true, secondaryFocused)
	// 	secondaryColorDiff := gui.GitCommand.Diff(file, false, !secondaryFocused)
	// 	return diff, secondaryColorDiff
	// }

	// diff, secondaryColorDiff := getDiffs()

	// // if we have e.g. a deleted file with nothing else to the diff will have only
	// // 4-5 lines in which case we'll swap panels
	// if len(strings.Split(diff, "\n")) < 5 {
	// 	if len(strings.Split(secondaryColorDiff, "\n")) < 5 {
	// 		return gui.handleStagingEscape(gui.g, nil)
	// 	}
	// 	secondaryFocused = !secondaryFocused
	// 	diff, secondaryColorDiff = getDiffs()
	// }

	// get diff from commit file that's currently selected
	commitFile := gui.getSelectedCommitFile(gui.g)
	if commitFile == nil {
		return gui.renderString(gui.g, "commitFiles", gui.Tr.SLocalize("NoCommiteFiles"))
	}

	diff, err := gui.GitCommand.ShowCommitFile(commitFile.Sha, commitFile.Name, true)
	if err != nil {
		return err
	}

	secondaryColorDiff := gui.GitCommand.PatchManager.RenderPatchForFile(commitFile.Name, false, false)
	if err != nil {
		return err
	}

	patchParser, err := commands.NewPatchParser(gui.Log, diff)
	if err != nil {
		return nil
	}

	if len(patchParser.StageableLines) == 0 {
		return gui.handleStagingEscape(gui.g, nil)
	}

	var selectedLineIdx int
	var firstLineIdx int
	var lastLineIdx int
	selectMode := LINE
	if state != nil {
		if state.SelectMode == HUNK {
			// this is tricky: we need to find out which hunk we just staged based on our old `state.PatchParser` (as opposed to the new `patchParser`)
			// we do this by getting the first line index of the original hunk, then
			// finding the next stageable line, then getting its containing hunk
			// in the new diff
			selectMode = HUNK
			prevNewHunk := state.PatchParser.GetHunkContainingLine(state.SelectedLineIdx, 0)
			selectedLineIdx = patchParser.GetNextStageableLineIndex(prevNewHunk.FirstLineIdx)
			newHunk := patchParser.GetHunkContainingLine(selectedLineIdx, 0)
			firstLineIdx, lastLineIdx = newHunk.FirstLineIdx, newHunk.LastLineIdx
		} else {
			selectedLineIdx = patchParser.GetNextStageableLineIndex(state.SelectedLineIdx)
			firstLineIdx, lastLineIdx = selectedLineIdx, selectedLineIdx
		}
	} else {
		selectedLineIdx = patchParser.StageableLines[0]
		firstLineIdx, lastLineIdx = selectedLineIdx, selectedLineIdx
	}

	gui.State.Panels.Staging = &stagingPanelState{
		PatchParser:      patchParser,
		SelectedLineIdx:  selectedLineIdx,
		SelectMode:       selectMode,
		FirstLineIdx:     firstLineIdx,
		LastLineIdx:      lastLineIdx,
		Diff:             diff,
		SecondaryFocused: secondaryFocused,
	}

	if err := gui.refreshView(); err != nil {
		return err
	}

	if err := gui.focusSelection(selectMode == HUNK); err != nil {
		return err
	}

	secondaryView := gui.getSecondaryView()
	secondaryView.Highlight = true
	secondaryView.Wrap = false

	gui.g.Update(func(*gocui.Gui) error {
		return gui.setViewContent(gui.g, gui.getSecondaryView(), secondaryColorDiff)
	})

	return nil
}

func (gui *Gui) handleTogglePanel(g *gocui.Gui, v *gocui.View) error {
	state := gui.State.Panels.Staging

	state.SecondaryFocused = !state.SecondaryFocused
	return gui.refreshStagingPanel()
}

func (gui *Gui) handleStagingEscape(g *gocui.Gui, v *gocui.View) error {
	gui.State.Panels.Staging = nil

	return gui.switchFocus(gui.g, nil, gui.getCommitFilesView())
}

func (gui *Gui) handleStagingPrevLine(g *gocui.Gui, v *gocui.View) error {
	return gui.handleCycleLine(-1)
}

func (gui *Gui) handleStagingNextLine(g *gocui.Gui, v *gocui.View) error {
	return gui.handleCycleLine(1)
}

func (gui *Gui) handleStagingPrevHunk(g *gocui.Gui, v *gocui.View) error {
	return gui.handleCycleHunk(-1)
}

func (gui *Gui) handleStagingNextHunk(g *gocui.Gui, v *gocui.View) error {
	return gui.handleCycleHunk(1)
}

func (gui *Gui) handleCycleHunk(change int) error {
	state := gui.State.Panels.Staging
	newHunk := state.PatchParser.GetHunkContainingLine(state.SelectedLineIdx, change)
	state.SelectedLineIdx = state.PatchParser.GetNextStageableLineIndex(newHunk.FirstLineIdx)
	if state.SelectMode == HUNK {
		state.FirstLineIdx, state.LastLineIdx = newHunk.FirstLineIdx, newHunk.LastLineIdx
	} else {
		state.FirstLineIdx, state.LastLineIdx = state.SelectedLineIdx, state.SelectedLineIdx
	}

	if err := gui.refreshView(); err != nil {
		return err
	}

	return gui.focusSelection(true)
}

func (gui *Gui) handleCycleLine(change int) error {
	state := gui.State.Panels.Staging

	if state.SelectMode == HUNK {
		return gui.handleCycleHunk(change)
	}

	newSelectedLineIdx := state.SelectedLineIdx + change
	if newSelectedLineIdx < 0 {
		newSelectedLineIdx = 0
	} else if newSelectedLineIdx > len(state.PatchParser.PatchLines)-1 {
		newSelectedLineIdx = len(state.PatchParser.PatchLines) - 1
	}

	state.SelectedLineIdx = newSelectedLineIdx

	if state.SelectMode == RANGE {
		if state.SelectedLineIdx < state.FirstLineIdx {
			state.FirstLineIdx = state.SelectedLineIdx
		} else {
			state.LastLineIdx = state.SelectedLineIdx
		}
	} else {
		state.LastLineIdx = state.SelectedLineIdx
		state.FirstLineIdx = state.SelectedLineIdx
	}

	if err := gui.refreshView(); err != nil {
		return err
	}

	return gui.focusSelection(false)
}

func (gui *Gui) refreshView() error {
	state := gui.State.Panels.Staging

	filename := gui.State.CommitFiles[gui.State.Panels.CommitFiles.SelectedLine].Name

	colorDiff := state.PatchParser.Render(state.FirstLineIdx, state.LastLineIdx, gui.GitCommand.PatchManager.GetFileIncLineIndices(filename))

	mainView := gui.getMainView()
	mainView.Highlight = true
	mainView.Wrap = false

	gui.g.Update(func(*gocui.Gui) error {
		return gui.setViewContent(gui.g, gui.getMainView(), colorDiff)
	})

	return nil
}

// focusSelection works out the best focus for the staging panel given the
// selected line and size of the hunk
func (gui *Gui) focusSelection(includeCurrentHunk bool) error {
	stagingView := gui.getMainView()
	state := gui.State.Panels.Staging

	_, viewHeight := stagingView.Size()
	bufferHeight := viewHeight - 1
	_, origin := stagingView.Origin()

	firstLineIdx := state.SelectedLineIdx
	lastLineIdx := state.SelectedLineIdx

	if includeCurrentHunk {
		hunk := state.PatchParser.GetHunkContainingLine(state.SelectedLineIdx, 0)
		firstLineIdx = hunk.FirstLineIdx
		lastLineIdx = hunk.LastLineIdx
	}

	margin := 0 // we may want to have a margin in place to show context  but right now I'm thinking we keep this at zero

	var newOrigin int
	if firstLineIdx-origin < margin {
		newOrigin = firstLineIdx - margin
	} else if lastLineIdx-origin > bufferHeight-margin {
		newOrigin = lastLineIdx - bufferHeight + margin
	} else {
		newOrigin = origin
	}

	gui.g.Update(func(*gocui.Gui) error {
		if err := stagingView.SetOrigin(0, newOrigin); err != nil {
			return err
		}

		return stagingView.SetCursor(0, state.SelectedLineIdx-newOrigin)
	})

	return nil
}

func (gui *Gui) handleStageSelection(g *gocui.Gui, v *gocui.View) error {
	state := gui.State.Panels.Staging

	// add range of lines to those set for the file
	commitFile := gui.getSelectedCommitFile(gui.g)
	if commitFile == nil {
		return gui.renderString(gui.g, "commitFiles", gui.Tr.SLocalize("NoCommiteFiles"))
	}

	gui.GitCommand.PatchManager.AddFileLineRange(commitFile.Name, state.FirstLineIdx, state.LastLineIdx)

	if err := gui.refreshCommitFilesView(); err != nil {
		return err
	}

	if err := gui.refreshStagingPanel(); err != nil {
		return err
	}

	return nil

	// return gui.applySelection(false)
}

func (gui *Gui) handleResetSelection(g *gocui.Gui, v *gocui.View) error {
	state := gui.State.Panels.Staging

	// add range of lines to those set for the file
	commitFile := gui.getSelectedCommitFile(gui.g)
	if commitFile == nil {
		return gui.renderString(gui.g, "commitFiles", gui.Tr.SLocalize("NoCommiteFiles"))
	}

	gui.GitCommand.PatchManager.RemoveFileLineRange(commitFile.Name, state.FirstLineIdx, state.LastLineIdx)

	if err := gui.refreshCommitFilesView(); err != nil {
		return err
	}

	if err := gui.refreshStagingPanel(); err != nil {
		return err
	}

	return nil

	// return gui.applySelection(true)
}

func (gui *Gui) applySelection(reverse bool) error {
	state := gui.State.Panels.Staging

	if !reverse && state.SecondaryFocused {
		return gui.createErrorPanel(gui.g, gui.Tr.SLocalize("CantStageStaged"))
	}

	file, err := gui.getSelectedFile(gui.g)
	if err != nil {
		return err
	}

	patch := commands.ModifiedPatchForRange(gui.Log, file.Name, state.Diff, state.FirstLineIdx, state.LastLineIdx, reverse, false)

	if patch == "" {
		return nil
	}

	// apply the patch then refresh this panel
	// create a new temp file with the patch, then call git apply with that patch
	err = gui.GitCommand.ApplyPatch(patch, false, !reverse || state.SecondaryFocused, "")
	if err != nil {
		return err
	}

	if state.SelectMode == RANGE {
		state.SelectMode = LINE
	}

	if err := gui.refreshFiles(); err != nil {
		return err
	}
	if err := gui.refreshStagingPanel(); err != nil {
		return err
	}
	return nil
}

func (gui *Gui) handleToggleSelectRange(g *gocui.Gui, v *gocui.View) error {
	state := gui.State.Panels.Staging
	if state.SelectMode == RANGE {
		state.SelectMode = LINE
	} else {
		state.SelectMode = RANGE
	}
	state.FirstLineIdx, state.LastLineIdx = state.SelectedLineIdx, state.SelectedLineIdx

	return gui.refreshView()
}

func (gui *Gui) handleToggleSelectHunk(g *gocui.Gui, v *gocui.View) error {
	state := gui.State.Panels.Staging

	if state.SelectMode == HUNK {
		state.SelectMode = LINE
		state.FirstLineIdx, state.LastLineIdx = state.SelectedLineIdx, state.SelectedLineIdx
	} else {
		state.SelectMode = HUNK
		selectedHunk := state.PatchParser.GetHunkContainingLine(state.SelectedLineIdx, 0)
		state.FirstLineIdx, state.LastLineIdx = selectedHunk.FirstLineIdx, selectedHunk.LastLineIdx
	}

	if err := gui.refreshView(); err != nil {
		return err
	}

	return gui.focusSelection(state.SelectMode == HUNK)
}
