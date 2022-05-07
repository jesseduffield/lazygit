package lbl

import (
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/sirupsen/logrus"
)

// State represents the current state of the line-by-line context i.e. when
// you're staging a file line-by-line or you're building a patch from an existing commit
// this struct holds the info about the diff you're interacting with and what's currently selected.
type State struct {
	selectedLineIdx   int
	rangeStartLineIdx int
	diff              string
	patchParser       *patch.PatchParser
	selectMode        selectMode
}

// these represent what select mode we're in
type selectMode int

const (
	LINE selectMode = iota
	RANGE
	HUNK
)

func NewState(diff string, selectedLineIdx int, oldState *State, log *logrus.Entry) *State {
	patchParser := patch.NewPatchParser(log, diff)

	if len(patchParser.StageableLines) == 0 {
		return nil
	}

	rangeStartLineIdx := 0
	if oldState != nil {
		rangeStartLineIdx = oldState.rangeStartLineIdx
	}

	selectMode := LINE
	// if we have clicked from the outside to focus the main view we'll pass in a non-negative line index so that we can instantly select that line
	if selectedLineIdx >= 0 {
		selectMode = RANGE
		rangeStartLineIdx = selectedLineIdx
	} else if oldState != nil {
		// if we previously had a selectMode of RANGE, we want that to now be line again
		if oldState.selectMode == HUNK {
			selectMode = HUNK
		}
		selectedLineIdx = patchParser.GetNextStageableLineIndex(oldState.selectedLineIdx)
	} else {
		selectedLineIdx = patchParser.StageableLines[0]
	}

	return &State{
		patchParser:       patchParser,
		selectedLineIdx:   selectedLineIdx,
		selectMode:        selectMode,
		rangeStartLineIdx: rangeStartLineIdx,
		diff:              diff,
	}
}

func (s *State) GetSelectedLineIdx() int {
	return s.selectedLineIdx
}

func (s *State) GetDiff() string {
	return s.diff
}

func (s *State) ToggleSelectHunk() {
	if s.selectMode == HUNK {
		s.selectMode = LINE
	} else {
		s.selectMode = HUNK
	}
}

func (s *State) ToggleSelectRange() {
	if s.selectMode == RANGE {
		s.selectMode = LINE
	} else {
		s.selectMode = RANGE
		s.rangeStartLineIdx = s.selectedLineIdx
	}
}

func (s *State) SelectingHunk() bool {
	return s.selectMode == HUNK
}

func (s *State) SelectingRange() bool {
	return s.selectMode == RANGE
}

func (s *State) SelectingLine() bool {
	return s.selectMode == LINE
}

func (s *State) SetLineSelectMode() {
	s.selectMode = LINE
}

func (s *State) SelectLine(newSelectedLineIdx int) {
	if newSelectedLineIdx < 0 {
		newSelectedLineIdx = 0
	} else if newSelectedLineIdx > len(s.patchParser.PatchLines)-1 {
		newSelectedLineIdx = len(s.patchParser.PatchLines) - 1
	}

	s.selectedLineIdx = newSelectedLineIdx
}

func (s *State) SelectNewLineForRange(newSelectedLineIdx int) {
	s.rangeStartLineIdx = newSelectedLineIdx

	s.selectMode = RANGE

	s.SelectLine(newSelectedLineIdx)
}

func (s *State) CycleSelection(forward bool) {
	if s.SelectingHunk() {
		s.CycleHunk(forward)
	} else {
		s.CycleLine(forward)
	}
}

func (s *State) CycleHunk(forward bool) {
	change := 1
	if !forward {
		change = -1
	}

	newHunk := s.patchParser.GetHunkContainingLine(s.selectedLineIdx, change)
	s.selectedLineIdx = s.patchParser.GetNextStageableLineIndex(newHunk.FirstLineIdx)
}

func (s *State) CycleLine(forward bool) {
	change := 1
	if !forward {
		change = -1
	}

	s.SelectLine(s.selectedLineIdx + change)
}

func (s *State) CurrentHunk() *patch.PatchHunk {
	return s.patchParser.GetHunkContainingLine(s.selectedLineIdx, 0)
}

func (s *State) SelectedRange() (int, int) {
	switch s.selectMode {
	case HUNK:
		hunk := s.CurrentHunk()
		return hunk.FirstLineIdx, hunk.LastLineIdx()
	case RANGE:
		if s.rangeStartLineIdx > s.selectedLineIdx {
			return s.selectedLineIdx, s.rangeStartLineIdx
		} else {
			return s.rangeStartLineIdx, s.selectedLineIdx
		}
	case LINE:
		return s.selectedLineIdx, s.selectedLineIdx
	default:
		// should never happen
		return 0, 0
	}
}

func (s *State) CurrentLineNumber() int {
	return s.CurrentHunk().LineNumberOfLine(s.selectedLineIdx)
}

func (s *State) AdjustSelectedLineIdx(change int) {
	s.SelectLine(s.selectedLineIdx + change)
}

func (s *State) RenderForLineIndices(includedLineIndices []int) string {
	firstLineIdx, lastLineIdx := s.SelectedRange()
	return s.patchParser.Render(firstLineIdx, lastLineIdx, includedLineIndices)
}

func (s *State) PlainRenderSelected() string {
	firstLineIdx, lastLineIdx := s.SelectedRange()
	return s.patchParser.PlainRenderLines(firstLineIdx, lastLineIdx)
}

func (s *State) SelectBottom() {
	s.SetLineSelectMode()
	s.SelectLine(len(s.patchParser.PatchLines) - 1)
}

func (s *State) SelectTop() {
	s.SetLineSelectMode()
	s.SelectLine(0)
}

func (s *State) CalculateOrigin(currentOrigin int, bufferHeight int) int {
	firstLineIdx, lastLineIdx := s.SelectedRange()

	return calculateOrigin(currentOrigin, bufferHeight, firstLineIdx, lastLineIdx, s.GetSelectedLineIdx(), s.selectMode)
}
