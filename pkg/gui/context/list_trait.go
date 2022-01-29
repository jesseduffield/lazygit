package context

import "github.com/jesseduffield/lazygit/pkg/gui/types"

type HasLength interface {
	GetItemsLength() int
}

type ListTrait struct {
	selectedIdx int
	HasLength
}

var _ types.IListPanelState = (*ListTrait)(nil)

func (self *ListTrait) GetSelectedLineIdx() int {
	return self.selectedIdx
}

func (self *ListTrait) SetSelectedLineIdx(value int) {
	self.selectedIdx = clamp(value, 0, self.GetItemsLength()-1)
}

// moves the cursor up or down by the given amount
func (self *ListTrait) MoveSelectedLine(value int) {
	self.SetSelectedLineIdx(self.selectedIdx + value)
}

// to be called when the model might have shrunk so that our selection is not not out of bounds
func (self *ListTrait) RefreshSelectedIdx() {
	self.SetSelectedLineIdx(self.selectedIdx)
}
