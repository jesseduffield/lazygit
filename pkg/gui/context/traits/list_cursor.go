package traits

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type RangeSelectMode int

const (
	// None means we are not selecting a range
	RangeSelectModeNone RangeSelectMode = iota
	// Sticky range select is started by pressing 'v', then the range is expanded
	// when you move up or down. It is cancelled by pressing 'v' again.
	RangeSelectModeSticky
	// Nonsticky range select is started by pressing shift+arrow and cancelled
	// when pressing up/down without shift, or by pressing 'v'
	RangeSelectModeNonSticky
)

type ListCursor struct {
	selectedIdx     int
	rangeSelectMode RangeSelectMode
	// value is ignored when rangeSelectMode is RangeSelectModeNone
	rangeStartIdx int
	// Get the length of the list. We use this to clamp the selection so that
	// the selected index is always valid
	getLength func() int
}

func NewListCursor(getLength func() int) *ListCursor {
	return &ListCursor{
		selectedIdx:     0,
		rangeStartIdx:   0,
		rangeSelectMode: RangeSelectModeNone,
		getLength:       getLength,
	}
}

var _ types.IListCursor = (*ListCursor)(nil)

func (self *ListCursor) GetSelectedLineIdx() int {
	return self.selectedIdx
}

// Sets the selected line index. Note, you probably don't want to use this directly,
// because it doesn't affect the range select mode or range start index. You should only
// use this for navigation situations where e.g. the user wants to jump to the top of
// a list while in range select mode so that the selection ends up being between
// the top of the list and the previous selection
func (self *ListCursor) SetSelectedLineIdx(value int) {
	self.selectedIdx = self.clampValue(value)
}

// Sets the selected index and cancels the range. You almost always want to use
// this instead of SetSelectedLineIdx. For example, if you want to jump the cursor
// to the top of a list after checking out a branch, you should use this method,
// or you may end up with a large range selection from the previous cursor position
// to the top of the list.
func (self *ListCursor) SetSelection(value int) {
	self.selectedIdx = self.clampValue(value)
	self.CancelRangeSelect()
}

func (self *ListCursor) SetSelectionRangeAndMode(selectedIdx, rangeStartIdx int, mode RangeSelectMode) {
	self.selectedIdx = self.clampValue(selectedIdx)
	self.rangeStartIdx = self.clampValue(rangeStartIdx)
	if mode == RangeSelectModeNonSticky && selectedIdx == rangeStartIdx {
		self.rangeSelectMode = RangeSelectModeNone
	} else {
		self.rangeSelectMode = mode
	}
}

// Returns the selectedIdx, the rangeStartIdx, and the mode of the current selection.
func (self *ListCursor) GetSelectionRangeAndMode() (int, int, RangeSelectMode) {
	if self.IsSelectingRange() {
		return self.selectedIdx, self.rangeStartIdx, self.rangeSelectMode
	} else {
		return self.selectedIdx, self.selectedIdx, self.rangeSelectMode
	}
}

func (self *ListCursor) clampValue(value int) int {
	clampedValue := -1
	length := self.getLength()
	if length > 0 {
		clampedValue = utils.Clamp(value, 0, length-1)
	}

	return clampedValue
}

// Moves the cursor up or down by the given amount.
// If we are in non-sticky range select mode, this will cancel the range select
func (self *ListCursor) MoveSelectedLine(change int) {
	if self.rangeSelectMode == RangeSelectModeNonSticky {
		self.CancelRangeSelect()
	}

	self.SetSelectedLineIdx(self.selectedIdx + change)
}

// Moves the cursor up or down by the given amount, and also moves the range start
// index by the same amount
func (self *ListCursor) MoveSelection(delta int) {
	self.selectedIdx = self.clampValue(self.selectedIdx + delta)
	if self.IsSelectingRange() {
		self.rangeStartIdx = self.clampValue(self.rangeStartIdx + delta)
	}
}

// To be called when the model might have shrunk so that our selection is not out of bounds
func (self *ListCursor) ClampSelection() {
	self.selectedIdx = self.clampValue(self.selectedIdx)
	self.rangeStartIdx = self.clampValue(self.rangeStartIdx)
}

func (self *ListCursor) Len() int {
	// The length of the model slice can change at any time, so the selection may
	// become out of bounds. To reduce the likelihood of this, we clamp the selection
	// whenever we obtain the length of the model.
	self.ClampSelection()

	return self.getLength()
}

func (self *ListCursor) GetRangeStartIdx() (int, bool) {
	if self.IsSelectingRange() {
		return self.rangeStartIdx, true
	}

	return 0, false
}

func (self *ListCursor) CancelRangeSelect() {
	self.rangeSelectMode = RangeSelectModeNone
}

// Returns true if we are in range select mode. Note that we may be in range select
// mode and still only selecting a single item. See AreMultipleItemsSelected below.
func (self *ListCursor) IsSelectingRange() bool {
	return self.rangeSelectMode != RangeSelectModeNone
}

// Returns true if we are in range select mode and selecting multiple items
func (self *ListCursor) AreMultipleItemsSelected() bool {
	startIdx, endIdx := self.GetSelectionRange()
	return startIdx != endIdx
}

func (self *ListCursor) GetSelectionRange() (int, int) {
	if self.IsSelectingRange() {
		return utils.SortRange(self.selectedIdx, self.rangeStartIdx)
	}

	return self.selectedIdx, self.selectedIdx
}

func (self *ListCursor) ToggleStickyRange() {
	if self.IsSelectingRange() {
		self.CancelRangeSelect()
	} else {
		self.rangeStartIdx = self.selectedIdx
		self.rangeSelectMode = RangeSelectModeSticky
	}
}

func (self *ListCursor) ExpandNonStickyRange(change int) {
	if !self.IsSelectingRange() {
		self.rangeStartIdx = self.selectedIdx
	}

	self.rangeSelectMode = RangeSelectModeNonSticky

	self.SetSelectedLineIdx(self.selectedIdx + change)
}
