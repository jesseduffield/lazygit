package context

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type ListContextTrait struct {
	types.Context

	c                 *types.HelperCommon
	list              types.IList
	viewTrait         *ViewTrait
	getDisplayStrings func(startIdx int, length int) [][]string
}

func (self *ListContextTrait) GetList() types.IList {
	return self.list
}

func (self *ListContextTrait) GetViewTrait() types.IViewTrait {
	return self.viewTrait
}

func (self *ListContextTrait) FocusLine() {
	// we need a way of knowing whether we've rendered to the view yet.
	self.viewTrait.FocusPoint(self.list.GetSelectedLineIdx())
	self.setFooter()
}

func (self *ListContextTrait) setFooter() {
	self.viewTrait.SetFooter(formatListFooter(self.list.GetSelectedLineIdx(), self.list.Len()))
}

func formatListFooter(selectedLineIdx int, length int) string {
	return fmt.Sprintf("%d of %d", selectedLineIdx+1, length)
}

func (self *ListContextTrait) HandleFocus(opts ...types.OnFocusOpts) error {
	self.FocusLine()

	self.viewTrait.SetHighlight(self.list.Len() > 0)

	return self.Context.HandleFocus(opts...)
}

func (self *ListContextTrait) HandleFocusLost() error {
	self.viewTrait.SetOriginX(0)

	return self.Context.HandleFocusLost()
}

// OnFocus assumes that the content of the context has already been rendered to the view. OnRender is the function which actually renders the content to the view
func (self *ListContextTrait) HandleRender() error {
	self.list.RefreshSelectedIdx()
	content := utils.RenderDisplayStrings(self.getDisplayStrings(0, self.list.Len()))
	self.viewTrait.SetContent(content)
	self.c.Render()
	self.setFooter()

	return nil
}

func (self *ListContextTrait) OnSearchSelect(selectedLineIdx int) error {
	self.GetList().SetSelectedLineIdx(selectedLineIdx)
	return self.HandleFocus()
}
