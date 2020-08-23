package gui

import (
	"fmt"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
)

// Currently there are two 'pseudo-panels' that make use of this 'pseudo-panel'.
// One is the staging panel where we stage files line-by-line, the other is the
// patch building panel where we add lines of an old commit's file to a patch.
// This file contains the logic around selecting lines and displaying the diffs
// staging_panel.go and patch_building_panel.go have functions specific to their
// use cases

// these represent what select mode we're in
const (
	LINE = iota
	RANGE
	HUNK
)

// returns whether the patch is empty so caller can escape if necessary
// both diffs should be non-coloured because we'll parse them and colour them here
func (gui *Gui) refreshLineByLinePanel(diff string, secondaryDiff string, secondaryFocused bool, selectedLineIdx int) (bool, error) {
	state := gui.State.Panels.LineByLine

	patchParser, err := patch.NewPatchParser(gui.Log, diff)
	if err != nil {
		return false, nil
	}

	if len(patchParser.StageableLines) == 0 {
		return true, nil
	}

	var firstLineIdx int
	var lastLineIdx int
	selectMode := LINE
	// if we have clicked from the outside to focus the main view we'll pass in a non-negative line index so that we can instantly select that line
	if selectedLineIdx >= 0 {
		selectMode = RANGE
		firstLineIdx, lastLineIdx = selectedLineIdx, selectedLineIdx
	} else if state != nil {
		if state.SelectMode == HUNK {
			// this is tricky: we need to find out which hunk we just staged based on our old `state.PatchParser` (as opposed to the new `patchParser`)
			// we do this by getting the first line index of the original hunk, then
			// finding the next stageable line, then getting its containing hunk
			// in the new diff
			selectMode = HUNK
			prevNewHunk := state.PatchParser.GetHunkContainingLine(state.SelectedLineIdx, 0)
			selectedLineIdx = patchParser.GetNextStageableLineIndex(prevNewHunk.FirstLineIdx)
			newHunk := patchParser.GetHunkContainingLine(selectedLineIdx, 0)
			firstLineIdx, lastLineIdx = newHunk.FirstLineIdx, newHunk.LastLineIdx()
		} else {
			selectedLineIdx = patchParser.GetNextStageableLineIndex(state.SelectedLineIdx)
			firstLineIdx, lastLineIdx = selectedLineIdx, selectedLineIdx
		}
	} else {
		selectedLineIdx = patchParser.StageableLines[0]
		firstLineIdx, lastLineIdx = selectedLineIdx, selectedLineIdx
	}

	gui.State.Panels.LineByLine = &lineByLinePanelState{
		PatchParser:      patchParser,
		SelectedLineIdx:  selectedLineIdx,
		SelectMode:       selectMode,
		FirstLineIdx:     firstLineIdx,
		LastLineIdx:      lastLineIdx,
		Diff:             diff,
		SecondaryFocused: secondaryFocused,
	}

	if err := gui.refreshMainViewForLineByLine(); err != nil {
		return false, err
	}

	if err := gui.focusSelection(selectMode == HUNK); err != nil {
		return false, err
	}

	secondaryView := gui.getSecondaryView()
	secondaryView.Highlight = true
	secondaryView.Wrap = false

	secondaryPatchParser, err := patch.NewPatchParser(gui.Log, secondaryDiff)
	if err != nil {
		return false, nil
	}

	gui.g.Update(func(*gocui.Gui) error {
		gui.setViewContent(gui.getSecondaryView(), secondaryPatchParser.Render(-1, -1, nil))
		return nil
	})

	return false, nil
}

func (gui *Gui) handleSelectPrevLine(g *gocui.Gui, v *gocui.View) error {
	return gui.handleCycleLine(-1)
}

func (gui *Gui) handleSelectNextLine(g *gocui.Gui, v *gocui.View) error {
	return gui.handleCycleLine(+1)
}

func (gui *Gui) handleSelectPrevHunk(g *gocui.Gui, v *gocui.View) error {
	state := gui.State.Panels.LineByLine
	newHunk := state.PatchParser.GetHunkContainingLine(state.SelectedLineIdx, -1)

	return gui.selectNewHunk(newHunk)
}

func (gui *Gui) handleSelectNextHunk(g *gocui.Gui, v *gocui.View) error {
	state := gui.State.Panels.LineByLine
	newHunk := state.PatchParser.GetHunkContainingLine(state.SelectedLineIdx, 1)

	return gui.selectNewHunk(newHunk)
}

func (gui *Gui) selectNewHunk(newHunk *patch.PatchHunk) error {
	state := gui.State.Panels.LineByLine
	state.SelectedLineIdx = state.PatchParser.GetNextStageableLineIndex(newHunk.FirstLineIdx)
	if state.SelectMode == HUNK {
		state.FirstLineIdx, state.LastLineIdx = newHunk.FirstLineIdx, newHunk.LastLineIdx()
	} else {
		state.FirstLineIdx, state.LastLineIdx = state.SelectedLineIdx, state.SelectedLineIdx
	}

	if err := gui.refreshMainViewForLineByLine(); err != nil {
		return err
	}

	return gui.focusSelection(true)
}

func (gui *Gui) handleCycleLine(change int) error {
	state := gui.State.Panels.LineByLine

	if state.SelectMode == HUNK {
		newHunk := state.PatchParser.GetHunkContainingLine(state.SelectedLineIdx, change)
		return gui.selectNewHunk(newHunk)
	}

	return gui.handleSelectNewLine(state.SelectedLineIdx + change)
}

func (gui *Gui) handleSelectNewLine(newSelectedLineIdx int) error {
	state := gui.State.Panels.LineByLine

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

	if err := gui.refreshMainViewForLineByLine(); err != nil {
		return err
	}

	return gui.focusSelection(false)
}

func (gui *Gui) handleMouseDown(g *gocui.Gui, v *gocui.View) error {
	state := gui.State.Panels.LineByLine

	if gui.popupPanelFocused() {
		return nil
	}

	newSelectedLineIdx := v.SelectedLineIdx()
	state.FirstLineIdx = newSelectedLineIdx
	state.LastLineIdx = newSelectedLineIdx

	state.SelectMode = RANGE

	return gui.handleSelectNewLine(newSelectedLineIdx)
}

func (gui *Gui) handleMouseDrag(g *gocui.Gui, v *gocui.View) error {
	if gui.popupPanelFocused() {
		return nil
	}

	return gui.handleSelectNewLine(v.SelectedLineIdx())
}

func (gui *Gui) handleMouseScrollUp(g *gocui.Gui, v *gocui.View) error {
	state := gui.State.Panels.LineByLine

	if gui.popupPanelFocused() {
		return nil
	}

	state.SelectMode = LINE

	return gui.handleCycleLine(-1)
}

func (gui *Gui) handleMouseScrollDown(g *gocui.Gui, v *gocui.View) error {
	state := gui.State.Panels.LineByLine

	if gui.popupPanelFocused() {
		return nil
	}

	state.SelectMode = LINE

	return gui.handleCycleLine(1)
}

func (gui *Gui) getSelectedCommitFileName() string {
	return gui.State.CommitFiles[gui.State.Panels.CommitFiles.SelectedLineIdx].Name
}

func (gui *Gui) refreshMainViewForLineByLine() error {
	state := gui.State.Panels.LineByLine

	var includedLineIndices []int
	// I'd prefer not to have knowledge of contexts using this file but I'm not sure
	// how to get around this
	if gui.currentContext().GetKey() == gui.Contexts.PatchBuilding.Context.GetKey() {
		filename := gui.getSelectedCommitFileName()
		var err error
		includedLineIndices, err = gui.GitCommand.PatchManager.GetFileIncLineIndices(filename)
		if err != nil {
			return err
		}
	}
	colorDiff := state.PatchParser.Render(state.FirstLineIdx, state.LastLineIdx, includedLineIndices)

	mainView := gui.getMainView()
	mainView.Highlight = true
	mainView.Wrap = false

	gui.g.Update(func(*gocui.Gui) error {
		gui.setViewContent(gui.getMainView(), colorDiff)
		return nil
	})

	return nil
}

// focusSelection works out the best focus for the staging panel given the
// selected line and size of the hunk
func (gui *Gui) focusSelection(includeCurrentHunk bool) error {
	stagingView := gui.getMainView()
	state := gui.State.Panels.LineByLine

	_, viewHeight := stagingView.Size()
	bufferHeight := viewHeight - 1
	_, origin := stagingView.Origin()

	firstLineIdx := state.SelectedLineIdx
	lastLineIdx := state.SelectedLineIdx

	if includeCurrentHunk {
		hunk := state.PatchParser.GetHunkContainingLine(state.SelectedLineIdx, 0)
		firstLineIdx = hunk.FirstLineIdx
		lastLineIdx = hunk.LastLineIdx()
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

func (gui *Gui) handleToggleSelectRange(g *gocui.Gui, v *gocui.View) error {
	state := gui.State.Panels.LineByLine
	if state.SelectMode == RANGE {
		state.SelectMode = LINE
	} else {
		state.SelectMode = RANGE
	}
	state.FirstLineIdx, state.LastLineIdx = state.SelectedLineIdx, state.SelectedLineIdx

	return gui.refreshMainViewForLineByLine()
}

func (gui *Gui) handleToggleSelectHunk(g *gocui.Gui, v *gocui.View) error {
	state := gui.State.Panels.LineByLine

	if state.SelectMode == HUNK {
		state.SelectMode = LINE
		state.FirstLineIdx, state.LastLineIdx = state.SelectedLineIdx, state.SelectedLineIdx
	} else {
		state.SelectMode = HUNK
		selectedHunk := state.PatchParser.GetHunkContainingLine(state.SelectedLineIdx, 0)
		state.FirstLineIdx, state.LastLineIdx = selectedHunk.FirstLineIdx, selectedHunk.LastLineIdx()
	}

	if err := gui.refreshMainViewForLineByLine(); err != nil {
		return err
	}

	return gui.focusSelection(state.SelectMode == HUNK)
}

func (gui *Gui) handleEscapeLineByLinePanel() {
	gui.State.Panels.LineByLine = nil
}

func (gui *Gui) handleOpenFileAtLine() error {
	// again, would be good to use inheritance here (or maybe even composition)
	var filename string
	switch gui.State.MainContext {
	case gui.Contexts.PatchBuilding.Context.GetKey():
		filename = gui.getSelectedCommitFileName()
	case gui.Contexts.Staging.Context.GetKey():
		file := gui.getSelectedFile()
		if file == nil {
			return nil
		}
		filename = file.Name
	default:
		return errors.Errorf("unknown main context: %s", gui.State.MainContext)
	}

	state := gui.State.Panels.LineByLine
	// need to look at current index, then work out what my hunk's header information is, and see how far my line is away from the hunk header
	selectedHunk := state.PatchParser.GetHunkContainingLine(state.SelectedLineIdx, 0)
	lineNumber := selectedHunk.LineNumberOfLine(state.SelectedLineIdx)
	filenameWithLineNum := fmt.Sprintf("%s:%d", filename, lineNumber)
	if err := gui.OSCommand.OpenFile(filenameWithLineNum); err != nil {
		return err
	}

	return nil
}
