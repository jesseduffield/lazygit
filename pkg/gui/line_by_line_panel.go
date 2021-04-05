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
type SelectMode int

const (
	LINE SelectMode = iota
	RANGE
	HUNK
)

// returns whether the patch is empty so caller can escape if necessary
// both diffs should be non-coloured because we'll parse them and colour them here
func (gui *Gui) refreshLineByLinePanel(diff string, secondaryDiff string, secondaryFocused bool, selectedLineIdx int, state *lBlPanelState) (bool, error) {
	gui.splitMainPanel(true)

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

	state = &lBlPanelState{
		PatchParser:      patchParser,
		SelectedLineIdx:  selectedLineIdx,
		SelectMode:       selectMode,
		FirstLineIdx:     firstLineIdx,
		LastLineIdx:      lastLineIdx,
		Diff:             diff,
		SecondaryFocused: secondaryFocused,
	}
	gui.State.Panels.LineByLine = state

	if err := gui.refreshMainViewForLineByLine(state); err != nil {
		return false, err
	}

	if err := gui.focusSelection(selectMode == HUNK, state); err != nil {
		return false, err
	}

	gui.Views.Secondary.Highlight = true
	gui.Views.Secondary.Wrap = false

	secondaryPatchParser, err := patch.NewPatchParser(gui.Log, secondaryDiff)
	if err != nil {
		return false, nil
	}

	gui.g.Update(func(*gocui.Gui) error {
		gui.setViewContent(gui.Views.Secondary, secondaryPatchParser.Render(-1, -1, nil))
		return nil
	})

	return false, nil
}

func (gui *Gui) handleSelectPrevLine() error {
	return gui.withLBLActiveCheck(func(state *lBlPanelState) error {
		return gui.LBLCycleLine(-1, state)
	})
}

func (gui *Gui) handleSelectNextLine() error {
	return gui.withLBLActiveCheck(func(state *lBlPanelState) error {
		return gui.LBLCycleLine(+1, state)
	})
}

func (gui *Gui) handleSelectPrevHunk() error {
	return gui.withLBLActiveCheck(func(state *lBlPanelState) error {
		newHunk := state.PatchParser.GetHunkContainingLine(state.SelectedLineIdx, -1)

		return gui.selectNewHunk(newHunk, state)
	})
}

func (gui *Gui) handleSelectNextHunk() error {
	return gui.withLBLActiveCheck(func(state *lBlPanelState) error {
		newHunk := state.PatchParser.GetHunkContainingLine(state.SelectedLineIdx, 1)

		return gui.selectNewHunk(newHunk, state)
	})
}

func (gui *Gui) selectNewHunk(newHunk *patch.PatchHunk, state *lBlPanelState) error {
	state.SelectedLineIdx = state.PatchParser.GetNextStageableLineIndex(newHunk.FirstLineIdx)
	if state.SelectMode == HUNK {
		state.FirstLineIdx, state.LastLineIdx = newHunk.FirstLineIdx, newHunk.LastLineIdx()
	} else {
		state.FirstLineIdx, state.LastLineIdx = state.SelectedLineIdx, state.SelectedLineIdx
	}

	if err := gui.refreshMainViewForLineByLine(state); err != nil {
		return err
	}

	return gui.focusSelection(true, state)
}

func (gui *Gui) LBLCycleLine(change int, state *lBlPanelState) error {
	if state.SelectMode == HUNK {
		newHunk := state.PatchParser.GetHunkContainingLine(state.SelectedLineIdx, change)
		return gui.selectNewHunk(newHunk, state)
	}

	return gui.LBLSelectLine(state.SelectedLineIdx+change, state)
}

func (gui *Gui) LBLSelectLine(newSelectedLineIdx int, state *lBlPanelState) error {
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

	if err := gui.refreshMainViewForLineByLine(state); err != nil {
		return err
	}

	return gui.focusSelection(false, state)
}

func (gui *Gui) handleLBLMouseDown() error {
	return gui.withLBLActiveCheck(func(state *lBlPanelState) error {
		if gui.popupPanelFocused() {
			return nil
		}

		newSelectedLineIdx := gui.Views.Main.SelectedLineIdx()
		state.FirstLineIdx = newSelectedLineIdx
		state.LastLineIdx = newSelectedLineIdx

		state.SelectMode = RANGE

		return gui.LBLSelectLine(newSelectedLineIdx, state)
	})
}

func (gui *Gui) handleMouseDrag() error {
	return gui.withLBLActiveCheck(func(state *lBlPanelState) error {
		if gui.popupPanelFocused() {
			return nil
		}

		return gui.LBLSelectLine(gui.Views.Main.SelectedLineIdx(), state)
	})
}

func (gui *Gui) getSelectedCommitFileName() string {
	idx := gui.State.Panels.CommitFiles.SelectedLineIdx

	return gui.State.CommitFileManager.GetItemAtIndex(idx).GetPath()
}

func (gui *Gui) refreshMainViewForLineByLine(state *lBlPanelState) error {
	var includedLineIndices []int
	// I'd prefer not to have knowledge of contexts using this file but I'm not sure
	// how to get around this
	if gui.currentContext().GetKey() == gui.State.Contexts.PatchBuilding.GetKey() {
		filename := gui.getSelectedCommitFileName()
		var err error
		includedLineIndices, err = gui.GitCommand.PatchManager.GetFileIncLineIndices(filename)
		if err != nil {
			return err
		}
	}
	colorDiff := state.PatchParser.Render(state.FirstLineIdx, state.LastLineIdx, includedLineIndices)

	gui.Views.Main.Highlight = true
	gui.Views.Main.Wrap = false

	gui.g.Update(func(*gocui.Gui) error {
		gui.setViewContent(gui.Views.Main, colorDiff)
		return nil
	})

	return nil
}

// focusSelection works out the best focus for the staging panel given the
// selected line and size of the hunk
func (gui *Gui) focusSelection(includeCurrentHunk bool, state *lBlPanelState) error {
	stagingView := gui.Views.Main

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

func (gui *Gui) handleToggleSelectRange() error {
	return gui.withLBLActiveCheck(func(state *lBlPanelState) error {
		if state.SelectMode == RANGE {
			state.SelectMode = LINE
		} else {
			state.SelectMode = RANGE
		}
		state.FirstLineIdx, state.LastLineIdx = state.SelectedLineIdx, state.SelectedLineIdx

		return gui.refreshMainViewForLineByLine(state)
	})
}

func (gui *Gui) handleToggleSelectHunk() error {
	return gui.withLBLActiveCheck(func(state *lBlPanelState) error {
		if state.SelectMode == HUNK {
			state.SelectMode = LINE
			state.FirstLineIdx, state.LastLineIdx = state.SelectedLineIdx, state.SelectedLineIdx
		} else {
			state.SelectMode = HUNK
			selectedHunk := state.PatchParser.GetHunkContainingLine(state.SelectedLineIdx, 0)
			state.FirstLineIdx, state.LastLineIdx = selectedHunk.FirstLineIdx, selectedHunk.LastLineIdx()
		}

		if err := gui.refreshMainViewForLineByLine(state); err != nil {
			return err
		}

		return gui.focusSelection(state.SelectMode == HUNK, state)
	})
}

func (gui *Gui) escapeLineByLinePanel() {
	gui.State.Panels.LineByLine = nil
}

func (gui *Gui) handleOpenFileAtLine() error {
	return gui.withLBLActiveCheck(func(state *lBlPanelState) error {
		// again, would be good to use inheritance here (or maybe even composition)
		var filename string
		switch gui.State.MainContext {
		case gui.State.Contexts.PatchBuilding.GetKey():
			filename = gui.getSelectedCommitFileName()
		case gui.State.Contexts.Staging.GetKey():
			file := gui.getSelectedFile()
			if file == nil {
				return nil
			}
			filename = file.Name
		default:
			return errors.Errorf("unknown main context: %s", gui.State.MainContext)
		}

		// need to look at current index, then work out what my hunk's header information is, and see how far my line is away from the hunk header
		selectedHunk := state.PatchParser.GetHunkContainingLine(state.SelectedLineIdx, 0)
		lineNumber := selectedHunk.LineNumberOfLine(state.SelectedLineIdx)
		filenameWithLineNum := fmt.Sprintf("%s:%d", filename, lineNumber)
		if err := gui.OSCommand.OpenFile(filenameWithLineNum); err != nil {
			return err
		}

		return nil
	})
}

func (gui *Gui) handleLineByLineNextPage() error {
	return gui.withLBLActiveCheck(func(state *lBlPanelState) error {
		newSelectedLineIdx := state.SelectedLineIdx + gui.pageDelta(gui.Views.Main)

		return gui.lineByLineNavigateTo(newSelectedLineIdx, state)
	})
}

func (gui *Gui) handleLineByLinePrevPage() error {
	return gui.withLBLActiveCheck(func(state *lBlPanelState) error {
		newSelectedLineIdx := state.SelectedLineIdx - gui.pageDelta(gui.Views.Main)

		return gui.lineByLineNavigateTo(newSelectedLineIdx, state)
	})
}

func (gui *Gui) handleLineByLineGotoBottom() error {
	return gui.withLBLActiveCheck(func(state *lBlPanelState) error {
		newSelectedLineIdx := len(state.PatchParser.PatchLines) - 1

		return gui.lineByLineNavigateTo(newSelectedLineIdx, state)
	})
}

func (gui *Gui) handleLineByLineGotoTop() error {
	return gui.withLBLActiveCheck(func(state *lBlPanelState) error {
		return gui.lineByLineNavigateTo(0, state)
	})
}

func (gui *Gui) handlelineByLineNavigateTo(selectedLineIdx int) error {
	return gui.withLBLActiveCheck(func(state *lBlPanelState) error {
		return gui.lineByLineNavigateTo(selectedLineIdx, state)
	})
}

func (gui *Gui) lineByLineNavigateTo(selectedLineIdx int, state *lBlPanelState) error {
	state.SelectMode = LINE

	return gui.LBLSelectLine(selectedLineIdx, state)
}

func (gui *Gui) withLBLActiveCheck(f func(*lBlPanelState) error) error {
	gui.Mutexes.LineByLinePanelMutex.Lock()
	defer gui.Mutexes.LineByLinePanelMutex.Unlock()

	state := gui.State.Panels.LineByLine
	if state == nil {
		return nil
	}

	return f(state)
}
