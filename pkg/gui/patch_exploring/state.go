package patch_exploring

import (
	"strings"

	"github.com/jesseduffield/generics/set"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
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

	// whether the user has switched to hunk mode manually; if hunk mode is on
	// but this is false, then hunk mode was enabled because the config makes it
	// on by default.
	// this makes a difference for whether we want to escape out of hunk mode
	userEnabledHunkMode bool
}

// these represent what select mode we're in
type selectMode int

const (
	LINE selectMode = iota
	RANGE
	HUNK
)

func NewState(diff string, selectedLineIdx int, view *gocui.View, oldState *State, useHunkModeByDefault bool) *State {
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
	if useHunkModeByDefault && !patch.IsSingleHunkForWholeFile() {
		selectMode = HUNK
	}

	userEnabledHunkMode := false
	if oldState != nil {
		userEnabledHunkMode = oldState.userEnabledHunkMode
	}

	// if we have clicked from the outside to focus the main view we'll pass in a non-negative line index so that we can instantly select that line
	if selectedLineIdx >= 0 {
		// Clamp to the number of wrapped view lines; index might be out of
		// bounds if a custom pager is being used which produces more lines
		selectedLineIdx = min(selectedLineIdx, len(viewLineIndices)-1)

		selectMode = RANGE
		rangeStartLineIdx = selectedLineIdx
	} else if oldState != nil {
		// if we previously had a selectMode of RANGE, we want that to now be line again (or hunk, if that's the default)
		if oldState.selectMode != RANGE {
			selectMode = oldState.selectMode
		}
		selectedLineIdx = viewLineIndices[patch.GetNextChangeIdx(oldState.patchLineIndices[oldState.selectedLineIdx])]
	} else {
		selectedLineIdx = viewLineIndices[patch.GetNextChangeIdx(0)]
	}

	return &State{
		patch:               patch,
		selectedLineIdx:     selectedLineIdx,
		selectMode:          selectMode,
		rangeStartLineIdx:   rangeStartLineIdx,
		rangeIsSticky:       false,
		diff:                diff,
		viewLineIndices:     viewLineIndices,
		patchLineIndices:    patchLineIndices,
		userEnabledHunkMode: userEnabledHunkMode,
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
		s.userEnabledHunkMode = true

		// If we are not currently on a change line, select the next one (or the
		// previous one if there is no next one):
		s.selectedLineIdx = s.viewLineIndices[s.patch.GetNextChangeIdx(
			s.patchLineIndices[s.selectedLineIdx])]
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

func (s *State) SelectingHunkEnabledByUser() bool {
	return s.selectMode == HUNK && s.userEnabledHunkMode
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

func (s *State) clampLineIdx(lineIdx int) int {
	return lo.Clamp(lineIdx, 0, len(s.patchLineIndices)-1)
}

// This just moves the cursor without caring about range select
func (s *State) selectLineWithoutRangeCheck(newSelectedLineIdx int) {
	s.selectedLineIdx = s.clampLineIdx(newSelectedLineIdx)
}

func (s *State) SelectNewLineForRange(newSelectedLineIdx int) {
	s.rangeStartLineIdx = s.clampLineIdx(newSelectedLineIdx)

	s.selectMode = RANGE

	s.selectLineWithoutRangeCheck(newSelectedLineIdx)
}

func (s *State) DragSelectLine(newSelectedLineIdx int) {
	s.selectMode = RANGE

	s.selectLineWithoutRangeCheck(newSelectedLineIdx)
}

func (s *State) CycleSelection(forward bool) {
	if s.SelectingHunk() {
		if forward {
			s.SelectNextHunk()
		} else {
			s.SelectPreviousHunk()
		}
	} else {
		s.CycleLine(forward)
	}
}

func (s *State) SelectPreviousHunk() {
	patchLines := s.patch.Lines()
	patchLineIdx := s.patchLineIndices[s.selectedLineIdx]
	nextNonChangeLine := patchLineIdx
	for nextNonChangeLine >= 0 && patchLines[nextNonChangeLine].IsChange() {
		nextNonChangeLine--
	}
	nextChangeLine := nextNonChangeLine
	for nextChangeLine >= 0 && !patchLines[nextChangeLine].IsChange() {
		nextChangeLine--
	}
	if nextChangeLine >= 0 {
		// Now we found a previous hunk, but we're on its last line. Skip to the beginning.
		for nextChangeLine > 0 && patchLines[nextChangeLine-1].IsChange() {
			nextChangeLine--
		}
		s.selectedLineIdx = s.viewLineIndices[nextChangeLine]
	}
}

func (s *State) SelectNextHunk() {
	patchLines := s.patch.Lines()
	patchLineIdx := s.patchLineIndices[s.selectedLineIdx]
	nextNonChangeLine := patchLineIdx
	for nextNonChangeLine < len(patchLines) && patchLines[nextNonChangeLine].IsChange() {
		nextNonChangeLine++
	}
	nextChangeLine := nextNonChangeLine
	for nextChangeLine < len(patchLines) && !patchLines[nextChangeLine].IsChange() {
		nextChangeLine++
	}
	if nextChangeLine < len(patchLines) {
		s.selectedLineIdx = s.viewLineIndices[nextChangeLine]
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

func (s *State) selectionRangeForCurrentBlockOfChanges() (int, int) {
	patchLines := s.patch.Lines()
	patchLineIdx := s.patchLineIndices[s.selectedLineIdx]

	patchStart := patchLineIdx
	for patchStart > 0 && patchLines[patchStart-1].IsChange() {
		patchStart--
	}

	patchEnd := patchLineIdx
	for patchEnd < len(patchLines)-1 && patchLines[patchEnd+1].IsChange() {
		patchEnd++
	}

	viewStart, viewEnd := s.viewLineIndices[patchStart], s.viewLineIndices[patchEnd]

	// Increase viewEnd in case the last patch line is wrapped to more than one view line.
	for viewEnd < len(s.patchLineIndices)-1 && s.patchLineIndices[viewEnd] == s.patchLineIndices[viewEnd+1] {
		viewEnd++
	}

	return viewStart, viewEnd
}

func (s *State) SelectedViewRange() (int, int) {
	switch s.selectMode {
	case HUNK:
		return s.selectionRangeForCurrentBlockOfChanges()
	case RANGE:
		if s.rangeStartLineIdx > s.selectedLineIdx {
			return s.selectedLineIdx, s.rangeStartLineIdx
		}
		return s.rangeStartLineIdx, s.selectedLineIdx
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

// Returns the line indices of the selected patch range that are changes (i.e. additions or deletions)
func (s *State) LineIndicesOfAddedOrDeletedLinesInSelectedPatchRange() []int {
	viewStart, viewEnd := s.SelectedViewRange()
	patchStart, patchEnd := s.patchLineIndices[viewStart], s.patchLineIndices[viewEnd]
	lines := s.patch.Lines()
	indices := []int{}
	for i := patchStart; i <= patchEnd; i++ {
		if lines[i].IsChange() {
			indices = append(indices, i)
		}
	}
	return indices
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
		view.Wrap, view.Editable, strings.TrimSuffix(diff, "\n"), view.InnerWidth(), view.TabWidth)
	return viewLineIndices, patchLineIndices
}

func (s *State) SelectNextStageableLineOfSameIncludedState(includedLines []int, included bool) {
	_, lastLineIdx := s.SelectedPatchRange()
	patchLineIdx, found := s.patch.GetNextChangeIdxOfSameIncludedState(lastLineIdx+1, includedLines, included)
	if found {
		s.SelectLine(s.viewLineIndices[patchLineIdx])
	}
}
