package patch_exploring

import (
	"github.com/jesseduffield/generics/set"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/sirupsen/logrus"
)

// State represents the current state of the patch explorer context i.e. when
// you're staging a file or you're building a patch from an existing commit
// this struct holds the info about the diff you're interacting with and what's currently selected.
type State struct {
	selectedLineIdx   int
	rangeStartLineIdx int
	// If a range is sticky, it means we expand the range when we move up or down.
	// Otherwise, we cancel the range when we move up or down.
	rangeIsSticky bool
	diff          string
	patch         *patch.Patch
	selectMode    selectMode
}

// these represent what select mode we're in
type selectMode int

const (
	LINE selectMode = iota
	RANGE
	HUNK
)

func NewState(diff string, selectedLineIdx int, oldState *State, log *logrus.Entry) *State {
	if oldState != nil && diff == oldState.diff && selectedLineIdx == -1 {
		// if we're here then we can return the old state. If selectedLineIdx was not -1
		// then that would mean we were trying to click and potentiall drag a range, which
		// is why in that case we continue below
		return oldState
	}

	patch := patch.Parse(diff)

	if !patch.ContainsChanges() {
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
		selectedLineIdx = patch.GetNextChangeIdx(oldState.selectedLineIdx)
	} else {
		selectedLineIdx = patch.GetNextChangeIdx(0)
	}

	return &State{
		patch:             patch,
		selectedLineIdx:   selectedLineIdx,
		selectMode:        selectMode,
		rangeStartLineIdx: rangeStartLineIdx,
		rangeIsSticky:     false,
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

func (s *State) ToggleStickySelectRange() {
	s.ToggleSelectRange(true)
}

func (s *State) ToggleSelectRange(sticky bool) {
	if s.selectMode == RANGE {
		s.selectMode = LINE
	} else {
		s.selectMode = RANGE
		s.rangeStartLineIdx = s.selectedLineIdx
		s.rangeIsSticky = sticky
	}
}

func (s *State) SetRangeIsSticky(value bool) {
	s.rangeIsSticky = value
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

// For when you move the cursor without holding shift (meaning if we're in
// a non-sticky range select, we'll cancel it)
func (s *State) SelectLine(newSelectedLineIdx int) {
	if s.selectMode == RANGE && !s.rangeIsSticky {
		s.selectMode = LINE
	}

	s.selectLineWithoutRangeCheck(newSelectedLineIdx)
}

// This just moves the cursor without caring about range select
func (s *State) selectLineWithoutRangeCheck(newSelectedLineIdx int) {
	if newSelectedLineIdx < 0 {
		newSelectedLineIdx = 0
	} else if newSelectedLineIdx > s.patch.LineCount()-1 {
		newSelectedLineIdx = s.patch.LineCount() - 1
	}

	s.selectedLineIdx = newSelectedLineIdx
}

func (s *State) SelectNewLineForRange(newSelectedLineIdx int) {
	s.rangeStartLineIdx = newSelectedLineIdx

	s.selectMode = RANGE

	s.selectLineWithoutRangeCheck(newSelectedLineIdx)
}

func (s *State) DragSelectLine(newSelectedLineIdx int) {
	s.selectMode = RANGE

	s.selectLineWithoutRangeCheck(newSelectedLineIdx)
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

	hunkIdx := s.patch.HunkContainingLine(s.selectedLineIdx)
	if hunkIdx != -1 {
		newHunkIdx := hunkIdx + change
		if newHunkIdx >= 0 && newHunkIdx < s.patch.HunkCount() {
			start := s.patch.HunkStartIdx(newHunkIdx)
			s.selectedLineIdx = s.patch.GetNextChangeIdx(start)
		}
	}
}

func (s *State) CycleLine(forward bool) {
	change := 1
	if !forward {
		change = -1
	}

	s.SelectLine(s.selectedLineIdx + change)
}

// This is called when we use shift+arrow to expand the range (i.e. a non-sticky
// range)
func (s *State) CycleRange(forward bool) {
	if !s.SelectingRange() {
		s.ToggleSelectRange(false)
	}

	s.SetRangeIsSticky(false)

	change := 1
	if !forward {
		change = -1
	}

	s.selectLineWithoutRangeCheck(s.selectedLineIdx + change)
}

// returns first and last patch line index of current hunk
func (s *State) CurrentHunkBounds() (int, int) {
	hunkIdx := s.patch.HunkContainingLine(s.selectedLineIdx)
	start := s.patch.HunkStartIdx(hunkIdx)
	end := s.patch.HunkEndIdx(hunkIdx)
	return start, end
}

func (s *State) SelectedRange() (int, int) {
	switch s.selectMode {
	case HUNK:
		return s.CurrentHunkBounds()
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
	return s.patch.LineNumberOfLine(s.selectedLineIdx)
}

func (s *State) AdjustSelectedLineIdx(change int) {
	s.SelectLine(s.selectedLineIdx + change)
}

func (s *State) RenderForLineIndices(isFocused bool, includedLineIndices []int) string {
	includedLineIndicesSet := set.NewFromSlice(includedLineIndices)
	return s.patch.FormatView(patch.FormatViewOpts{
		IncLineIndices: includedLineIndicesSet,
	})
}

func (s *State) PlainRenderSelected() string {
	firstLineIdx, lastLineIdx := s.SelectedRange()
	return s.patch.FormatRangePlain(firstLineIdx, lastLineIdx)
}

func (s *State) SelectBottom() {
	s.SetLineSelectMode()
	s.SelectLine(s.patch.LineCount() - 1)
}

func (s *State) SelectTop() {
	s.SetLineSelectMode()
	s.SelectLine(0)
}

func (s *State) CalculateOrigin(currentOrigin int, bufferHeight int, numLines int) int {
	firstLineIdx, lastLineIdx := s.SelectedRange()

	return calculateOrigin(currentOrigin, bufferHeight, numLines, firstLineIdx, lastLineIdx, s.GetSelectedLineIdx(), s.selectMode)
}

func (s *State) RangeStartLineIdx() (int, bool) {
	if s.selectMode == RANGE {
		return s.rangeStartLineIdx, true
	}

	return 0, false
}
