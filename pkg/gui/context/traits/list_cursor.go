package traits

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type HasLength interface {
	Len() int
}

type ListCursor struct {
	selectedIdx int
	list        HasLength
}

func NewListCursor(list HasLength) *ListCursor {
	return &ListCursor{selectedIdx: 0, list: list}
}

var _ types.IListCursor = (*ListCursor)(nil)

func (self *ListCursor) GetSelectedLineIdx() int {
	return self.selectedIdx
}

func (self *ListCursor) SetSelectedLineIdx(value int) {
	clampedValue := -1
	if self.list.Len() > 0 {
		clampedValue = utils.Clamp(value, 0, self.list.Len()-1)
	}

	self.selectedIdx = clampedValue
}

// moves the cursor up or down by the given amount
func (self *ListCursor) MoveSelectedLine(delta int) {
	self.SetSelectedLineIdx(self.selectedIdx + delta)
}

// to be called when the model might have shrunk so that our selection is not not out of bounds
func (self *ListCursor) RefreshSelectedIdx() {
	self.SetSelectedLineIdx(self.selectedIdx)
}

func (self *ListCursor) Len() int {
	return self.list.Len()
}
