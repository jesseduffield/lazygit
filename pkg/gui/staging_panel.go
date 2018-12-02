package gui

import (
	"errors"
	"io/ioutil"

	"github.com/davecgh/go-spew/spew"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/git"
)

func (gui *Gui) refreshStagingPanel() error {
	// get the currently selected file. Get the diff of that file directly, not
	// using any custom diff tools.
	// parse the file to find out where the chunks and unstaged changes are

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

	gui.Log.WithField("staging", "staging").Info("DIFF IS:")
	gui.Log.WithField("staging", "staging").Info(spew.Sdump(diff))
	gui.Log.WithField("staging", "staging").Info("hello")

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

	var currentLineIndex int
	if gui.State.StagingState != nil {
		end := len(stageableLines) - 1
		if end < gui.State.StagingState.CurrentLineIndex {
			currentLineIndex = end
		} else {
			currentLineIndex = gui.State.StagingState.CurrentLineIndex
		}
	} else {
		currentLineIndex = 0
	}

	gui.State.StagingState = &stagingState{
		StageableLines:   stageableLines,
		HunkStarts:       hunkStarts,
		CurrentLineIndex: currentLineIndex,
		Diff:             diff,
	}

	if len(stageableLines) == 0 {
		return errors.New("No lines to stage")
	}

	stagingView := gui.getStagingView(gui.g)
	stagingView.SetCursor(0, stageableLines[currentLineIndex])
	stagingView.SetOrigin(0, 0)
	return gui.renderString(gui.g, "staging", colorDiff)
}

func (gui *Gui) handleStagingEscape(g *gocui.Gui, v *gocui.View) error {
	if _, err := gui.g.SetViewOnBottom("staging"); err != nil {
		return err
	}

	return gui.switchFocus(gui.g, nil, gui.getFilesView(gui.g))
}

// nextNumber returns the next index, cycling if we reach the end
func nextIndex(numbers []int, currentNumber int) int {
	for index, number := range numbers {
		if number > currentNumber {
			return index
		}
	}
	return 0
}

// prevNumber returns the next number, cycling if we reach the end
func prevIndex(numbers []int, currentNumber int) int {
	end := len(numbers) - 1
	for i := end; i >= 0; i -= 1 {
		if numbers[i] < currentNumber {
			return i
		}
	}
	return end
}

func (gui *Gui) handleStagingKeyUp(g *gocui.Gui, v *gocui.View) error {
	return gui.handleCycleLine(true)
}

func (gui *Gui) handleStagingKeyDown(g *gocui.Gui, v *gocui.View) error {
	return gui.handleCycleLine(false)
}

func (gui *Gui) handleCycleLine(up bool) error {
	state := gui.State.StagingState
	lineNumbers := state.StageableLines
	currentLine := lineNumbers[state.CurrentLineIndex]
	var newIndex int
	if up {
		newIndex = prevIndex(lineNumbers, currentLine)
	} else {
		newIndex = nextIndex(lineNumbers, currentLine)
	}

	state.CurrentLineIndex = newIndex
	stagingView := gui.getStagingView(gui.g)
	stagingView.SetCursor(0, lineNumbers[newIndex])
	stagingView.SetOrigin(0, 0)
	return nil
}

func (gui *Gui) handleStageLine(g *gocui.Gui, v *gocui.View) error {
	state := gui.State.StagingState
	p, err := git.NewPatchModifier(gui.Log)
	if err != nil {
		return err
	}

	currentLine := state.StageableLines[state.CurrentLineIndex]
	patch, err := p.ModifyPatch(state.Diff, currentLine)
	if err != nil {
		return err
	}

	// for logging purposes
	ioutil.WriteFile("patch.diff", []byte(patch), 0600)

	// apply the patch then refresh this panel
	// create a new temp file with the patch, then call git apply with that patch
	_, err = gui.GitCommand.ApplyPatch(patch)
	if err != nil {
		panic(err)
	}

	gui.refreshStagingPanel()
	gui.refreshFiles(gui.g)
	return nil
}
