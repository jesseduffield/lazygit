package context

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type ListContextTrait struct {
	types.Context

	c                 *types.ControllerCommon
	list              types.IList
	viewTrait         *ViewTrait
	getDisplayStrings func(startIdx int, length int) [][]string
}

func (self *ListContextTrait) GetList() types.IList {
	return self.list
}

// TODO: remove
func (self *ListContextTrait) GetPanelState() types.IListPanelState {
	return self.list
}

func (self *ListContextTrait) GetViewTrait() types.IViewTrait {
	return self.viewTrait
}

func (self *ListContextTrait) FocusLine() {
	// we need a way of knowing whether we've rendered to the view yet.
	self.viewTrait.FocusPoint(self.list.GetSelectedLineIdx())
	self.viewTrait.SetFooter(formatListFooter(self.list.GetSelectedLineIdx(), self.list.GetItemsLength()))
}

func formatListFooter(selectedLineIdx int, length int) string {
	return fmt.Sprintf("%d of %d", selectedLineIdx+1, length)
}

func (self *ListContextTrait) HandleFocus(opts ...types.OnFocusOpts) error {
	self.FocusLine()

	return self.Context.HandleFocus(opts...)
}

func (self *ListContextTrait) HandleFocusLost() error {
	self.viewTrait.SetOriginX(0)

	return self.Context.HandleFocus()
}

// OnFocus assumes that the content of the context has already been rendered to the view. OnRender is the function which actually renders the content to the view
func (self *ListContextTrait) HandleRender() error {
	self.list.RefreshSelectedIdx()
	content := utils.RenderDisplayStrings(self.getDisplayStrings(0, self.list.GetItemsLength()))
	self.viewTrait.SetContent(content)
	self.c.Render()

	return nil
}

func (self *ListContextTrait) OnSearchSelect(selectedLineIdx int) error {
	self.GetList().SetSelectedLineIdx(selectedLineIdx)
	return self.HandleFocus()
}
