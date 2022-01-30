package context

import (
	"fmt"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type ListContextTrait struct {
	base      types.IBaseContext
	list      types.IList
	viewTrait *ViewTrait

	takeFocus func() error

	GetDisplayStrings func(startIdx int, length int) [][]string
	OnFocus           func(...types.OnFocusOpts) error
	OnRenderToMain    func(...types.OnFocusOpts) error
	OnFocusLost       func() error

	// if this is true, we'll call GetDisplayStrings for just the visible part of the
	// view and re-render that. This is useful when you need to render different
	// content based on the selection (e.g. for showing the selected commit)
	RenderSelection bool

	c *types.ControllerCommon
}

func (self *ListContextTrait) GetPanelState() types.IListPanelState {
	return self.list
}

func (self *ListContextTrait) FocusLine() {
	// we need a way of knowing whether we've rendered to the view yet.
	self.viewTrait.FocusPoint(self.list.GetSelectedLineIdx())
	if self.RenderSelection {
		min, max := self.viewTrait.ViewPortYBounds()
		displayStrings := self.GetDisplayStrings(min, max)
		content := utils.RenderDisplayStrings(displayStrings)
		self.viewTrait.SetViewPortContent(content)
	}
	self.viewTrait.SetFooter(formatListFooter(self.list.GetSelectedLineIdx(), self.list.GetItemsLength()))
}

func formatListFooter(selectedLineIdx int, length int) string {
	return fmt.Sprintf("%d of %d", selectedLineIdx+1, length)
}

// OnFocus assumes that the content of the context has already been rendered to the view. OnRender is the function which actually renders the content to the view
func (self *ListContextTrait) HandleRender() error {
	if self.GetDisplayStrings != nil {
		self.list.RefreshSelectedIdx()
		content := utils.RenderDisplayStrings(self.GetDisplayStrings(0, self.list.GetItemsLength()))
		self.viewTrait.SetContent(content)
		self.c.Render()
	}

	return nil
}

func (self *ListContextTrait) HandleFocusLost() error {
	if self.OnFocusLost != nil {
		return self.OnFocusLost()
	}

	self.viewTrait.SetOriginX(0)

	return nil
}

func (self *ListContextTrait) HandleFocus(opts ...types.OnFocusOpts) error {
	self.FocusLine()

	if self.OnFocus != nil {
		if err := self.OnFocus(opts...); err != nil {
			return err
		}
	}

	if self.OnRenderToMain != nil {
		if err := self.OnRenderToMain(opts...); err != nil {
			return err
		}
	}

	return nil
}

func (self *ListContextTrait) HandlePrevLine() error {
	return self.handleLineChange(-1)
}

func (self *ListContextTrait) HandleNextLine() error {
	return self.handleLineChange(1)
}

func (self *ListContextTrait) HandleScrollLeft() error {
	return self.scroll(self.viewTrait.ScrollLeft)
}

func (self *ListContextTrait) HandleScrollRight() error {
	return self.scroll(self.viewTrait.ScrollRight)
}

func (self *ListContextTrait) scroll(scrollFunc func()) error {
	scrollFunc()

	return self.HandleFocus()
}

func (self *ListContextTrait) handleLineChange(change int) error {
	before := self.list.GetSelectedLineIdx()
	self.list.MoveSelectedLine(change)
	after := self.list.GetSelectedLineIdx()

	// doing this check so that if we're holding the up key at the start of the list
	// we're not constantly re-rendering the main view.
	if before != after {
		return self.HandleFocus()
	}

	return nil
}

func (self *ListContextTrait) HandlePrevPage() error {
	return self.handleLineChange(-self.viewTrait.PageDelta())
}

func (self *ListContextTrait) HandleNextPage() error {
	return self.handleLineChange(self.viewTrait.PageDelta())
}

func (self *ListContextTrait) HandleGotoTop() error {
	return self.handleLineChange(-self.list.GetItemsLength())
}

func (self *ListContextTrait) HandleGotoBottom() error {
	return self.handleLineChange(self.list.GetItemsLength())
}

func (self *ListContextTrait) HandleClick(onClick func() error) error {
	prevSelectedLineIdx := self.list.GetSelectedLineIdx()
	// because we're handling a click, we need to determine the new line idx based
	// on the view itself.
	newSelectedLineIdx := self.viewTrait.SelectedLineIdx()

	currentContextKey := self.c.CurrentContext().GetKey()
	alreadyFocused := currentContextKey == self.base.GetKey()

	// we need to focus the view
	if !alreadyFocused {
		if err := self.takeFocus(); err != nil {
			return err
		}
	}

	if newSelectedLineIdx > self.list.GetItemsLength()-1 {
		return nil
	}

	self.list.SetSelectedLineIdx(newSelectedLineIdx)

	if prevSelectedLineIdx == newSelectedLineIdx && alreadyFocused && onClick != nil {
		return onClick()
	}
	return self.HandleFocus()
}

func (self *ListContextTrait) OnSearchSelect(selectedLineIdx int) error {
	self.list.SetSelectedLineIdx(selectedLineIdx)
	return self.HandleFocus()
}

func (self *ListContextTrait) HandleRenderToMain() error {
	if self.OnRenderToMain != nil {
		return self.OnRenderToMain()
	}

	return nil
}

func (self *ListContextTrait) Keybindings(
	getKey func(key string) interface{},
	config config.KeybindingConfig,
	guards types.KeybindingGuards,
) []*types.Binding {
	return []*types.Binding{
		{Tag: "navigation", Key: getKey(config.Universal.PrevItemAlt), Modifier: gocui.ModNone, Handler: self.HandlePrevLine},
		{Tag: "navigation", Key: getKey(config.Universal.PrevItem), Modifier: gocui.ModNone, Handler: self.HandlePrevLine},
		{Tag: "navigation", Key: gocui.MouseWheelUp, Modifier: gocui.ModNone, Handler: self.HandlePrevLine},
		{Tag: "navigation", Key: getKey(config.Universal.NextItemAlt), Modifier: gocui.ModNone, Handler: self.HandleNextLine},
		{Tag: "navigation", Key: getKey(config.Universal.NextItem), Modifier: gocui.ModNone, Handler: self.HandleNextLine},
		{Tag: "navigation", Key: getKey(config.Universal.PrevPage), Modifier: gocui.ModNone, Handler: self.HandlePrevPage, Description: self.c.Tr.LcPrevPage},
		{Tag: "navigation", Key: getKey(config.Universal.NextPage), Modifier: gocui.ModNone, Handler: self.HandleNextPage, Description: self.c.Tr.LcNextPage},
		{Tag: "navigation", Key: getKey(config.Universal.GotoTop), Modifier: gocui.ModNone, Handler: self.HandleGotoTop, Description: self.c.Tr.LcGotoTop},
		{Key: gocui.MouseLeft, Modifier: gocui.ModNone, Handler: func() error { return self.HandleClick(nil) }},
		{Tag: "navigation", Key: gocui.MouseWheelDown, Modifier: gocui.ModNone, Handler: self.HandleNextLine},
		{Tag: "navigation", Key: getKey(config.Universal.ScrollLeft), Modifier: gocui.ModNone, Handler: self.HandleScrollLeft},
		{Tag: "navigation", Key: getKey(config.Universal.ScrollRight), Modifier: gocui.ModNone, Handler: self.HandleScrollRight},
		{
			Key:         getKey(config.Universal.StartSearch),
			Handler:     func() error { self.c.OpenSearch(); return nil },
			Description: self.c.Tr.LcStartSearch,
			Tag:         "navigation",
		},
		{
			Key:         getKey(config.Universal.GotoBottom),
			Description: self.c.Tr.LcGotoBottom,
			Handler:     self.HandleGotoBottom,
			Tag:         "navigation",
		},
	}
}
