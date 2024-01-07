package traits

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type HasLength interface {
	Len() int
}

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
	list          HasLength
}

func NewListCursor(list HasLength) *ListCursor {
	return &ListCursor{
		selectedIdx:     0,
		rangeStartIdx:   0,
		rangeSelectMode: RangeSelectModeNone,
		list:            list,
	}
}

var _ types.IListCursor = (*ListCursor)(nil)

func (self *ListCursor) GetSelectedLineIdx() int {
	return self.selectedIdx
}

func (self *ListCursor) SetSelectedLineIdx(value int) {
	self.selectedIdx = self.clampValue(value)
}

func (self *ListCursor) clampValue(value int) int {
	clampedValue := -1
	if self.list.Len() > 0 {
		clampedValue = utils.Clamp(value, 0, self.list.Len()-1)
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
	return self.list.Len()
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

func (self *ListCursor) IsSelectingRange() bool {
	return self.rangeSelectMode != RangeSelectModeNone
}

func (self *ListCursor) GetSelectionRange() (int, int) {
	if self.IsSelectingRange() {
		return utils.MinMax(self.selectedIdx, self.rangeStartIdx)
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
