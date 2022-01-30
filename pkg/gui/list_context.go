package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type ListContext struct {
	GetItemsLength    func() int
	GetDisplayStrings func(startIdx int, length int) [][]string
	OnFocus           func(...types.OnFocusOpts) error
	OnRenderToMain    func(...types.OnFocusOpts) error
	OnFocusLost       func() error

	OnGetSelectedItemId func() string
	OnGetPanelState     func() types.IListPanelState
	// if this is true, we'll call GetDisplayStrings for just the visible part of the
	// view and re-render that. This is useful when you need to render different
	// content based on the selection (e.g. for showing the selected commit)
	RenderSelection bool

	Gui *Gui

	*context.BaseContext
}

var _ types.IListContext = &ListContext{}

func (self *ListContext) GetPanelState() types.IListPanelState {
	return self.OnGetPanelState()
}

func (self *ListContext) FocusLine() {
	view, err := self.Gui.g.View(self.ViewName)
	if err != nil {
		// ignoring error for now
		return
	}

	// we need a way of knowing whether we've rendered to the view yet.
	view.FocusPoint(view.OriginX(), self.GetPanelState().GetSelectedLineIdx())
	if self.RenderSelection {
		_, originY := view.Origin()
		displayStrings := self.GetDisplayStrings(originY, view.InnerHeight()+1)
		self.Gui.renderDisplayStringsInViewPort(view, displayStrings)
	}
	view.Footer = formatListFooter(self.GetPanelState().GetSelectedLineIdx(), self.GetItemsLength())
}

func formatListFooter(selectedLineIdx int, length int) string {
	return fmt.Sprintf("%d of %d", selectedLineIdx+1, length)
}

func (self *ListContext) GetSelectedItemId() string {
	return self.OnGetSelectedItemId()
}

// OnFocus assumes that the content of the context has already been rendered to the view. OnRender is the function which actually renders the content to the view
func (self *ListContext) HandleRender() error {
	view, err := self.Gui.g.View(self.ViewName)
	if err != nil {
		return nil
	}

	if self.GetDisplayStrings != nil {
		self.Gui.refreshSelectedLine(self.GetPanelState(), self.GetItemsLength())
		self.Gui.renderDisplayStrings(view, self.GetDisplayStrings(0, self.GetItemsLength()))
		self.Gui.render()
	}

	return nil
}

func (self *ListContext) HandleFocusLost() error {
	if self.OnFocusLost != nil {
		return self.OnFocusLost()
	}

	view, err := self.Gui.g.View(self.ViewName)
	if err != nil {
		return nil
	}

	_ = view.SetOriginX(0)

	return nil
}

func (self *ListContext) HandleFocus(opts ...types.OnFocusOpts) error {
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

func (self *ListContext) HandlePrevLine() error {
	return self.handleLineChange(-1)
}

func (self *ListContext) HandleNextLine() error {
	return self.handleLineChange(1)
}

func (self *ListContext) HandleScrollLeft() error {
	return self.scroll(self.Gui.scrollLeft)
}

func (self *ListContext) HandleScrollRight() error {
	return self.scroll(self.Gui.scrollRight)
}

func (self *ListContext) scroll(scrollFunc func(*gocui.View)) error {
	if self.ignoreKeybinding() {
		return nil
	}

	// get the view, move the origin
	view, err := self.Gui.g.View(self.ViewName)
	if err != nil {
		return nil
	}

	scrollFunc(view)

	return self.HandleFocus()
}

func (self *ListContext) ignoreKeybinding() bool {
	return !self.Gui.isPopupPanel(self.ViewName) && self.Gui.popupPanelFocused()
}

func (self *ListContext) handleLineChange(change int) error {
	if self.ignoreKeybinding() {
		return nil
	}

	selectedLineIdx := self.GetPanelState().GetSelectedLineIdx()
	if (change < 0 && selectedLineIdx == 0) || (change > 0 && selectedLineIdx == self.GetItemsLength()-1) {
		return nil
	}

	self.Gui.changeSelectedLine(self.GetPanelState(), self.GetItemsLength(), change)

	return self.HandleFocus()
}

func (self *ListContext) HandleNextPage() error {
	view, err := self.Gui.g.View(self.ViewName)
	if err != nil {
		return nil
	}
	delta := self.Gui.pageDelta(view)

	return self.handleLineChange(delta)
}

func (self *ListContext) HandleGotoTop() error {
	return self.handleLineChange(-self.GetItemsLength())
}

func (self *ListContext) HandleGotoBottom() error {
	return self.handleLineChange(self.GetItemsLength())
}

func (self *ListContext) HandlePrevPage() error {
	view, err := self.Gui.g.View(self.ViewName)
	if err != nil {
		return nil
	}

	delta := self.Gui.pageDelta(view)

	return self.handleLineChange(-delta)
}

func (self *ListContext) HandleClick(onClick func() error) error {
	if self.ignoreKeybinding() {
		return nil
	}

	view, err := self.Gui.g.View(self.ViewName)
	if err != nil {
		return nil
	}

	prevSelectedLineIdx := self.GetPanelState().GetSelectedLineIdx()
	newSelectedLineIdx := view.SelectedLineIdx()

	// we need to focus the view
	if err := self.Gui.c.PushContext(self); err != nil {
		return err
	}

	if newSelectedLineIdx > self.GetItemsLength()-1 {
		return nil
	}

	self.GetPanelState().SetSelectedLineIdx(newSelectedLineIdx)

	prevViewName := self.Gui.currentViewName()
	if prevSelectedLineIdx == newSelectedLineIdx && prevViewName == self.ViewName && onClick != nil {
		return onClick()
	}
	return self.HandleFocus()
}

func (self *ListContext) OnSearchSelect(selectedLineIdx int) error {
	self.GetPanelState().SetSelectedLineIdx(selectedLineIdx)
	return self.HandleFocus()
}

func (self *ListContext) HandleRenderToMain() error {
	if self.OnRenderToMain != nil {
		return self.OnRenderToMain()
	}

	return nil
}

func (self *ListContext) Keybindings(
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
		{Tag: "navigation", Key: getKey(config.Universal.PrevPage), Modifier: gocui.ModNone, Handler: self.HandlePrevPage, Description: self.Gui.c.Tr.LcPrevPage},
		{Tag: "navigation", Key: getKey(config.Universal.NextPage), Modifier: gocui.ModNone, Handler: self.HandleNextPage, Description: self.Gui.c.Tr.LcNextPage},
		{Tag: "navigation", Key: getKey(config.Universal.GotoTop), Modifier: gocui.ModNone, Handler: self.HandleGotoTop, Description: self.Gui.c.Tr.LcGotoTop},
		{Key: gocui.MouseLeft, Modifier: gocui.ModNone, Handler: func() error { return self.HandleClick(nil) }},
		{Tag: "navigation", Key: gocui.MouseWheelDown, Modifier: gocui.ModNone, Handler: self.HandleNextLine},
		{Tag: "navigation", Key: getKey(config.Universal.ScrollLeft), Modifier: gocui.ModNone, Handler: self.HandleScrollLeft},
		{Tag: "navigation", Key: getKey(config.Universal.ScrollRight), Modifier: gocui.ModNone, Handler: self.HandleScrollRight},
		{
			Key:         getKey(config.Universal.StartSearch),
			Handler:     func() error { return self.Gui.handleOpenSearch(self.GetViewName()) },
			Description: self.Gui.c.Tr.LcStartSearch,
			Tag:         "navigation",
		},
		{
			Key:         getKey(config.Universal.GotoBottom),
			Description: self.Gui.c.Tr.LcGotoBottom,
			Handler:     self.HandleGotoBottom,
			Tag:         "navigation",
		},
	}
}
