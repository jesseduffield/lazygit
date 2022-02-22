package gui

import (
	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/jesseduffield/lazygit/pkg/gui/lbl"
)

// Currently there are two 'pseudo-panels' that make use of this 'pseudo-panel'.
// One is the staging panel where we stage files line-by-line, the other is the
// patch building panel where we add lines of an old commit's file to a patch.
// This file contains the logic around selecting lines and displaying the diffs
// staging_panel.go and patch_building_panel.go have functions specific to their
// use cases

// returns whether the patch is empty so caller can escape if necessary
// both diffs should be non-coloured because we'll parse them and colour them here
func (gui *Gui) refreshLineByLinePanel(diff string, secondaryDiff string, secondaryFocused bool, selectedLineIdx int) (bool, error) {
	gui.splitMainPanel(true)

	var oldState *lbl.State
	if gui.State.Panels.LineByLine != nil {
		oldState = gui.State.Panels.LineByLine.State
	}

	state := lbl.NewState(diff, selectedLineIdx, oldState, gui.Log)
	if state == nil {
		return true, nil
	}

	gui.State.Panels.LineByLine = &LblPanelState{
		State:            state,
		SecondaryFocused: secondaryFocused,
	}

	if err := gui.refreshMainViewForLineByLine(gui.State.Panels.LineByLine); err != nil {
		return false, err
	}

	if err := gui.focusSelection(gui.State.Panels.LineByLine); err != nil {
		return false, err
	}

	gui.Views.Secondary.Highlight = true
	gui.Views.Secondary.Wrap = false

	secondaryPatchParser := patch.NewPatchParser(gui.Log, secondaryDiff)

	gui.setViewContent(gui.Views.Secondary, secondaryPatchParser.Render(-1, -1, nil))

	return false, nil
}

func (gui *Gui) handleSelectPrevLine() error {
	return gui.withLBLActiveCheck(func(state *LblPanelState) error {
		state.CycleSelection(false)

		return gui.refreshAndFocusLblPanel(state)
	})
}

func (gui *Gui) handleSelectNextLine() error {
	return gui.withLBLActiveCheck(func(state *LblPanelState) error {
		state.CycleSelection(true)

		return gui.refreshAndFocusLblPanel(state)
	})
}

func (gui *Gui) handleSelectPrevHunk() error {
	return gui.withLBLActiveCheck(func(state *LblPanelState) error {
		state.CycleHunk(false)

		return gui.refreshAndFocusLblPanel(state)
	})
}

func (gui *Gui) handleSelectNextHunk() error {
	return gui.withLBLActiveCheck(func(state *LblPanelState) error {
		state.CycleHunk(true)

		return gui.refreshAndFocusLblPanel(state)
	})
}

func (gui *Gui) copySelectedToClipboard() error {
	return gui.withLBLActiveCheck(func(state *LblPanelState) error {
		selected := state.PlainRenderSelected()

		gui.c.LogAction(gui.c.Tr.Actions.CopySelectedTextToClipboard)
		if err := gui.os.CopyToClipboard(selected); err != nil {
			return gui.c.Error(err)
		}

		return nil
	})
}

func (gui *Gui) refreshAndFocusLblPanel(state *LblPanelState) error {
	if err := gui.refreshMainViewForLineByLine(state); err != nil {
		return err
	}

	return gui.focusSelection(state)
}

func (gui *Gui) handleLBLMouseDown() error {
	return gui.withLBLActiveCheck(func(state *LblPanelState) error {
		state.SelectNewLineForRange(gui.Views.Main.SelectedLineIdx())

		return gui.refreshAndFocusLblPanel(state)
	})
}

func (gui *Gui) handleMouseDrag() error {
	return gui.withLBLActiveCheck(func(state *LblPanelState) error {
		state.SelectLine(gui.Views.Main.SelectedLineIdx())

		return gui.refreshAndFocusLblPanel(state)
	})
}

func (gui *Gui) refreshMainViewForLineByLine(state *LblPanelState) error {
	var includedLineIndices []int
	// I'd prefer not to have knowledge of contexts using this file but I'm not sure
	// how to get around this
	if gui.currentContext().GetKey() == gui.State.Contexts.PatchBuilding.GetKey() {
		filename := gui.getSelectedCommitFileName()
		var err error
		includedLineIndices, err = gui.git.Patch.PatchManager.GetFileIncLineIndices(filename)
		if err != nil {
			return err
		}
	}
	colorDiff := state.RenderForLineIndices(includedLineIndices)

	gui.Views.Main.Highlight = true
	gui.Views.Main.Wrap = false

	gui.setViewContent(gui.Views.Main, colorDiff)

	return nil
}

// focusSelection works out the best focus for the staging panel given the
// selected line and size of the hunk
func (gui *Gui) focusSelection(state *LblPanelState) error {
	stagingView := gui.Views.Main

	_, viewHeight := stagingView.Size()
	bufferHeight := viewHeight - 1
	_, origin := stagingView.Origin()

	selectedLineIdx := state.GetSelectedLineIdx()

	newOrigin := state.CalculateOrigin(origin, bufferHeight)

	if err := stagingView.SetOriginY(newOrigin); err != nil {
		return err
	}

	return stagingView.SetCursor(0, selectedLineIdx-newOrigin)
}

func (gui *Gui) handleToggleSelectRange() error {
	return gui.withLBLActiveCheck(func(state *LblPanelState) error {
		state.ToggleSelectRange()

		return gui.refreshMainViewForLineByLine(state)
	})
}

func (gui *Gui) handleToggleSelectHunk() error {
	return gui.withLBLActiveCheck(func(state *LblPanelState) error {
		state.ToggleSelectHunk()

		return gui.refreshAndFocusLblPanel(state)
	})
}

func (gui *Gui) escapeLineByLinePanel() {
	gui.State.Panels.LineByLine = nil
}

func (gui *Gui) handleOpenFileAtLine() error {
	return gui.withLBLActiveCheck(func(state *LblPanelState) error {
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
		lineNumber := state.CurrentLineNumber()
		if err := gui.os.OpenFileAtLine(filename, lineNumber); err != nil {
			return err
		}

		return nil
	})
}

func (gui *Gui) handleLineByLineNextPage() error {
	return gui.withLBLActiveCheck(func(state *LblPanelState) error {
		state.SetLineSelectMode()
		state.AdjustSelectedLineIdx(gui.pageDelta(gui.Views.Main))

		return gui.refreshAndFocusLblPanel(state)
	})
}

func (gui *Gui) handleLineByLinePrevPage() error {
	return gui.withLBLActiveCheck(func(state *LblPanelState) error {
		state.SetLineSelectMode()
		state.AdjustSelectedLineIdx(-gui.pageDelta(gui.Views.Main))

		return gui.refreshAndFocusLblPanel(state)
	})
}

func (gui *Gui) handleLineByLineGotoBottom() error {
	return gui.withLBLActiveCheck(func(state *LblPanelState) error {
		state.SelectBottom()

		return gui.refreshAndFocusLblPanel(state)
	})
}

func (gui *Gui) handleLineByLineGotoTop() error {
	return gui.withLBLActiveCheck(func(state *LblPanelState) error {
		state.SelectTop()

		return gui.refreshAndFocusLblPanel(state)
	})
}

func (gui *Gui) handlelineByLineNavigateTo(selectedLineIdx int) error {
	return gui.withLBLActiveCheck(func(state *LblPanelState) error {
		state.SetLineSelectMode()
		state.SelectLine(selectedLineIdx)

		return gui.refreshAndFocusLblPanel(state)
	})
}

func (gui *Gui) withLBLActiveCheck(f func(*LblPanelState) error) error {
	gui.Mutexes.LineByLinePanelMutex.Lock()
	defer gui.Mutexes.LineByLinePanelMutex.Unlock()

	state := gui.State.Panels.LineByLine
	if state == nil {
		return nil
	}

	return f(state)
}

func (gui *Gui) handleLineByLineEdit() error {
	file := gui.getSelectedFile()
	if file == nil {
		return nil
	}

	lineNumber := gui.State.Panels.LineByLine.CurrentLineNumber()
	return gui.helpers.Files.EditFileAtLine(file.Name, lineNumber)
}
