package gui

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/git"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) refreshStagingPanel() error {
	file, err := gui.getSelectedFile(gui.g)
	if err != nil {
		if err != gui.Errors.ErrNoFiles {
			return err
		}
		return gui.handleStagingEscape(gui.g, nil)
	}

	if !file.HasUnstagedChanges {
		return gui.handleStagingEscape(gui.g, nil)
	}

	// note for custom diffs, we'll need to send a flag here saying not to use the custom diff
	diff := gui.GitCommand.Diff(file, true)
	colorDiff := gui.GitCommand.Diff(file, false)

	if len(diff) < 2 {
		return gui.handleStagingEscape(gui.g, nil)
	}

	// parse the diff and store the line numbers of hunks and stageable lines
	// TODO: maybe instantiate this at application start
	p, err := git.NewPatchParser(gui.Log)
	if err != nil {
		return nil
	}
	hunkStarts, stageableLines, err := p.ParsePatch(diff)
	if err != nil {
		return nil
	}

	var selectedLine int
	if gui.State.Panels.Staging != nil {
		end := len(stageableLines) - 1
		if end < gui.State.Panels.Staging.SelectedLine {
			selectedLine = end
		} else {
			selectedLine = gui.State.Panels.Staging.SelectedLine
		}
	} else {
		selectedLine = 0
	}

	gui.State.Panels.Staging = &stagingPanelState{
		StageableLines: stageableLines,
		HunkStarts:     hunkStarts,
		SelectedLine:   selectedLine,
		Diff:           diff,
	}

	if len(stageableLines) == 0 {
		return gui.createErrorPanel(gui.g, "No lines to stage")
	}

	if err := gui.focusLineAndHunk(); err != nil {
		return err
	}

	mainView := gui.getMainView()
	mainView.Highlight = true
	mainView.Wrap = false

	gui.g.Update(func(*gocui.Gui) error {
		return gui.setViewContent(gui.g, gui.getMainView(), colorDiff)
	})

	return nil
}

func (gui *Gui) handleStagingEscape(g *gocui.Gui, v *gocui.View) error {
	gui.State.Panels.Staging = nil

	return gui.switchFocus(gui.g, nil, gui.getFilesView())
}

func (gui *Gui) handleStagingPrevLine(g *gocui.Gui, v *gocui.View) error {
	return gui.handleCycleLine(true)
}

func (gui *Gui) handleStagingNextLine(g *gocui.Gui, v *gocui.View) error {
	return gui.handleCycleLine(false)
}

func (gui *Gui) handleStagingPrevHunk(g *gocui.Gui, v *gocui.View) error {
	return gui.handleCycleHunk(true)
}

func (gui *Gui) handleStagingNextHunk(g *gocui.Gui, v *gocui.View) error {
	return gui.handleCycleHunk(false)
}

func (gui *Gui) handleCycleHunk(prev bool) error {
	state := gui.State.Panels.Staging
	lineNumbers := state.StageableLines
	currentLine := lineNumbers[state.SelectedLine]
	currentHunkIndex := utils.PrevIndex(state.HunkStarts, currentLine)
	var newHunkIndex int
	if prev {
		if currentHunkIndex == 0 {
			newHunkIndex = len(state.HunkStarts) - 1
		} else {
			newHunkIndex = currentHunkIndex - 1
		}
	} else {
		if currentHunkIndex == len(state.HunkStarts)-1 {
			newHunkIndex = 0
		} else {
			newHunkIndex = currentHunkIndex + 1
		}
	}

	state.SelectedLine = utils.NextIndex(lineNumbers, state.HunkStarts[newHunkIndex])

	return gui.focusLineAndHunk()
}

func (gui *Gui) handleCycleLine(prev bool) error {
	state := gui.State.Panels.Staging
	lineNumbers := state.StageableLines
	currentLine := lineNumbers[state.SelectedLine]
	var newIndex int
	if prev {
		newIndex = utils.PrevIndex(lineNumbers, currentLine)
	} else {
		newIndex = utils.NextIndex(lineNumbers, currentLine)
	}
	state.SelectedLine = newIndex

	return gui.focusLineAndHunk()
}

// focusLineAndHunk works out the best focus for the staging panel given the
// selected line and size of the hunk
func (gui *Gui) focusLineAndHunk() error {
	stagingView := gui.getMainView()
	state := gui.State.Panels.Staging

	lineNumber := state.StageableLines[state.SelectedLine]

	// we want the bottom line of the view buffer to ideally be the bottom line
	// of the hunk, but if the hunk is too big we'll just go three lines beyond
	// the currently selected line so that the user can see the context
	var bottomLine int
	nextHunkStartIndex := utils.NextIndex(state.HunkStarts, lineNumber)
	if nextHunkStartIndex == 0 {
		// for now linesHeight is an efficient means of getting the number of lines
		// in the patch. However if we introduce word wrap we'll need to update this
		bottomLine = stagingView.LinesHeight() - 1
	} else {
		bottomLine = state.HunkStarts[nextHunkStartIndex] - 1
	}

	hunkStartIndex := utils.PrevIndex(state.HunkStarts, lineNumber)
	hunkStart := state.HunkStarts[hunkStartIndex]
	// if it's the first hunk we'll also show the diff header
	if hunkStartIndex == 0 {
		hunkStart = 0
	}

	_, height := stagingView.Size()
	// if this hunk is too big, we will just ensure that the user can at least
	// see three lines of context below the cursor
	if bottomLine-hunkStart > height {
		bottomLine = lineNumber + 3
	}

	return gui.generalFocusLine(lineNumber, bottomLine, stagingView)
}

func (gui *Gui) handleStageHunk(g *gocui.Gui, v *gocui.View) error {
	return gui.handleStageLineOrHunk(true)
}

func (gui *Gui) handleStageLine(g *gocui.Gui, v *gocui.View) error {
	return gui.handleStageLineOrHunk(false)
}

func (gui *Gui) handleStageLineOrHunk(hunk bool) error {
	state := gui.State.Panels.Staging
	p, err := git.NewPatchModifier(gui.Log)
	if err != nil {
		return err
	}

	currentLine := state.StageableLines[state.SelectedLine]
	var patch string
	if hunk {
		patch, err = p.ModifyPatchForHunk(state.Diff, state.HunkStarts, currentLine)
	} else {
		patch, err = p.ModifyPatchForLine(state.Diff, currentLine)
	}
	if err != nil {
		return err
	}

	// for logging purposes
	// ioutil.WriteFile("patch.diff", []byte(patch), 0600)

	// apply the patch then refresh this panel
	// create a new temp file with the patch, then call git apply with that patch
	_, err = gui.GitCommand.ApplyPatch(patch)
	if err != nil {
		return err
	}

	if err := gui.refreshFiles(); err != nil {
		return err
	}
	if err := gui.refreshStagingPanel(); err != nil {
		return err
	}
	return nil
}
