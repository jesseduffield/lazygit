package patch_exploring

import (
	"strings"

	"github.com/jesseduffield/generics/set"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// State represents the current state of the patch explorer context i.e. when
// you're staging a file or you're building a patch from an existing commit
// this struct holds the info about the diff you're interacting with and what's currently selected.
type State struct {
	// These are in terms of view lines (wrapped), not patch lines
	selectedLineIdx   int
	rangeStartLineIdx int
	// If a range is sticky, it means we expand the range when we move up or down.
	// Otherwise, we cancel the range when we move up or down.
	rangeIsSticky bool
	diff          string
	patch         *patch.Patch
	selectMode    selectMode

	// Array of indices of the wrapped lines indexed by a patch line index
	viewLineIndices []int
	// Array of indices of the original patch lines indexed by a wrapped view line index
	patchLineIndices []int
}

// these represent what select mode we're in
type selectMode int

const (
	LINE selectMode = iota
	RANGE
	HUNK
)

func NewState(diff string, selectedLineIdx int, view *gocui.View, oldState *State) *State {
	if oldState != nil && diff == oldState.diff && selectedLineIdx == -1 {
		// if we're here then we can return the old state. If selectedLineIdx was not -1
		// then that would mean we were trying to click and potentially drag a range, which
		// is why in that case we continue below
		return oldState
	}

	patch := patch.Parse(diff)

	if !patch.ContainsChanges() {
		return nil
	}

	viewLineIndices, patchLineIndices := wrapPatchLines(diff, view)

	rangeStartLineIdx := 0
	if oldState != nil {
		rangeStartLineIdx = oldState.rangeStartLineIdx
	}

	selectMode := LINE
	// if we have clicked from the outside to focus the main view we'll pass in a non-negative line index so that we can instantly select that line
	if selectedLineIdx >= 0 {
		// Clamp to the number of wrapped view lines; index might be out of
		// bounds if a custom pager is being used which produces more lines
		selectedLineIdx = min(selectedLineIdx, len(viewLineIndices)-1)

		selectMode = RANGE
		rangeStartLineIdx = selectedLineIdx
	} else if oldState != nil {
		// if we previously had a selectMode of RANGE, we want that to now be line again
		if oldState.selectMode == HUNK {
			selectMode = HUNK
		}
		selectedLineIdx = viewLineIndices[patch.GetNextChangeIdx(oldState.patchLineIndices[oldState.selectedLineIdx])]
	} else {
		selectedLineIdx = viewLineIndices[patch.GetNextChangeIdx(0)]
	}

	return &State{
		patch:             patch,
		selectedLineIdx:   selectedLineIdx,
		selectMode:        selectMode,
		rangeStartLineIdx: rangeStartLineIdx,
		rangeIsSticky:     false,
		diff:              diff,
		viewLineIndices:   viewLineIndices,
		patchLineIndices:  patchLineIndices,
	}
}

func (s *State) OnViewWidthChanged(view *gocui.View) {
	if !view.Wrap {
		return
	}

	selectedPatchLineIdx := s.patchLineIndices[s.selectedLineIdx]
	var rangeStartPatchLineIdx int
	if s.selectMode == RANGE {
		rangeStartPatchLineIdx = s.patchLineIndices[s.rangeStartLineIdx]
	}
	s.viewLineIndices, s.patchLineIndices = wrapPatchLines(s.diff, view)
	s.selectedLineIdx = s.viewLineIndices[selectedPatchLineIdx]
	if s.selectMode == RANGE {
		s.rangeStartLineIdx = s.viewLineIndices[rangeStartPatchLineIdx]
	}
}

func (s *State) GetSelectedPatchLineIdx() int {
	return s.patchLineIndices[s.selectedLineIdx]
}

func (s *State) GetSelectedViewLineIdx() int {
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
	if s.SelectingRange() {
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
	return s.selectMode == RANGE && (s.rangeIsSticky || s.rangeStartLineIdx != s.selectedLineIdx)
}

func (s *State) SelectingLine() bool {
	return s.selectMode == LINE
}

func (s *State) SetLineSelectMode() {
	s.selectMode = LINE
}

func (s *State) DismissHunkSelectMode() {
	if s.SelectingHunk() {
		s.selectMode = LINE
	}
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
	} else if newSelectedLineIdx > len(s.patchLineIndices)-1 {
		newSelectedLineIdx = len(s.patchLineIndices) - 1
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

	hunkIdx := s.patch.HunkContainingLine(s.patchLineIndices[s.selectedLineIdx])
	if hunkIdx != -1 {
		newHunkIdx := hunkIdx + change
		if newHunkIdx >= 0 && newHunkIdx < s.patch.HunkCount() {
			start := s.patch.HunkStartIdx(newHunkIdx)
			s.selectedLineIdx = s.viewLineIndices[s.patch.GetNextChangeIdx(start)]
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
	hunkIdx := s.patch.HunkContainingLine(s.patchLineIndices[s.selectedLineIdx])
	start := s.patch.HunkStartIdx(hunkIdx)
	end := s.patch.HunkEndIdx(hunkIdx)
	return start, end
}

func (s *State) SelectedViewRange() (int, int) {
	switch s.selectMode {
	case HUNK:
		start, end := s.CurrentHunkBounds()
		return s.viewLineIndices[start], s.viewLineIndices[end]
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

func (s *State) SelectedPatchRange() (int, int) {
	start, end := s.SelectedViewRange()
	return s.patchLineIndices[start], s.patchLineIndices[end]
}

func (s *State) CurrentLineNumber() int {
	return s.patch.LineNumberOfLine(s.patchLineIndices[s.selectedLineIdx])
}

func (s *State) AdjustSelectedLineIdx(change int) {
	s.DismissHunkSelectMode()
	s.SelectLine(s.selectedLineIdx + change)
}

func (s *State) RenderForLineIndices(includedLineIndices []int) string {
	includedLineIndicesSet := set.NewFromSlice(includedLineIndices)
	return s.patch.FormatView(patch.FormatViewOpts{
		IncLineIndices: includedLineIndicesSet,
	})
}

func (s *State) PlainRenderSelected() string {
	firstLineIdx, lastLineIdx := s.SelectedPatchRange()
	return s.patch.FormatRangePlain(firstLineIdx, lastLineIdx)
}

func (s *State) SelectBottom() {
	s.DismissHunkSelectMode()
	s.SelectLine(len(s.patchLineIndices) - 1)
}

func (s *State) SelectTop() {
	s.DismissHunkSelectMode()
	s.SelectLine(0)
}

func (s *State) CalculateOrigin(currentOrigin int, bufferHeight int, numLines int) int {
	firstLineIdx, lastLineIdx := s.SelectedViewRange()

	return calculateOrigin(currentOrigin, bufferHeight, numLines, firstLineIdx, lastLineIdx, s.GetSelectedViewLineIdx(), s.selectMode)
}

func wrapPatchLines(diff string, view *gocui.View) ([]int, []int) {
	_, viewLineIndices, patchLineIndices := utils.WrapViewLinesToWidth(
		view.Wrap, strings.TrimSuffix(diff, "\n"), view.InnerWidth())
	return viewLineIndices, patchLineIndices
}
