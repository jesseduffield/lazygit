package controllers

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type ListControllerFactory struct {
	c *ControllerCommon
}

func NewListControllerFactory(c *ControllerCommon) *ListControllerFactory {
	return &ListControllerFactory{
		c: c,
	}
}

func (self *ListControllerFactory) Create(context types.IListContext) *ListController {
	return &ListController{
		baseController: baseController{},
		c:              self.c,
		context:        context,
	}
}

type ListController struct {
	baseController
	c *ControllerCommon

	context types.IListContext
}

func (self *ListController) Context() types.Context {
	return self.context
}

func (self *ListController) HandlePrevLine() error {
	return self.handleLineChange(-1)
}

func (self *ListController) HandleNextLine() error {
	return self.handleLineChange(1)
}

func (self *ListController) HandleScrollLeft() error {
	return self.scrollHorizontal(self.context.GetViewTrait().ScrollLeft)
}

func (self *ListController) HandleScrollRight() error {
	return self.scrollHorizontal(self.context.GetViewTrait().ScrollRight)
}

func (self *ListController) HandleScrollUp() error {
	scrollHeight := self.c.UserConfig.Gui.ScrollHeight
	self.context.GetViewTrait().ScrollUp(scrollHeight)

	return nil
}

func (self *ListController) HandleScrollDown() error {
	scrollHeight := self.c.UserConfig.Gui.ScrollHeight
	self.context.GetViewTrait().ScrollDown(scrollHeight)

	return nil
}

func (self *ListController) scrollHorizontal(scrollFunc func()) error {
	scrollFunc()

	return self.context.HandleFocus(types.OnFocusOpts{})
}

func (self *ListController) handleLineChange(change int) error {
	return self.handleLineChangeAux(
		self.context.GetList().MoveSelectedLine, change,
	)
}

func (self *ListController) HandleRangeSelectChange(change int) error {
	return self.handleLineChangeAux(
		self.context.GetList().ExpandNonStickyRange, change,
	)
}

func (self *ListController) handleLineChangeAux(f func(int), change int) error {
	list := self.context.GetList()

	rangeBefore := list.IsSelectingRange()
	before := list.GetSelectedLineIdx()
	f(change)
	rangeAfter := list.IsSelectingRange()
	after := list.GetSelectedLineIdx()

	if err := self.pushContextIfNotFocused(); err != nil {
		return err
	}

	// doing this check so that if we're holding the up key at the start of the list
	// we're not constantly re-rendering the main view.
	cursorMoved := before != after
	if cursorMoved {
		if change == -1 {
			checkScrollUp(self.context.GetViewTrait(), self.c.UserConfig,
				self.context.ModelIndexToViewIndex(before), self.context.ModelIndexToViewIndex(after))
		} else if change == 1 {
			checkScrollDown(self.context.GetViewTrait(), self.c.UserConfig,
				self.context.ModelIndexToViewIndex(before), self.context.ModelIndexToViewIndex(after))
		}
	}

	if cursorMoved || rangeBefore != rangeAfter {
		return self.context.HandleFocus(types.OnFocusOpts{})
	}

	return nil
}

func (self *ListController) HandlePrevPage() error {
	return self.handleLineChange(-self.context.GetViewTrait().PageDelta())
}

func (self *ListController) HandleNextPage() error {
	return self.handleLineChange(self.context.GetViewTrait().PageDelta())
}

func (self *ListController) HandleGotoTop() error {
	return self.handleLineChange(-self.context.GetList().Len())
}

func (self *ListController) HandleGotoBottom() error {
	return self.handleLineChange(self.context.GetList().Len())
}

func (self *ListController) HandleToggleRangeSelect() error {
	list := self.context.GetList()

	list.ToggleStickyRange()

	return self.context.HandleFocus(types.OnFocusOpts{})
}

func (self *ListController) HandleRangeSelectDown() error {
	return self.HandleRangeSelectChange(1)
}

func (self *ListController) HandleRangeSelectUp() error {
	return self.HandleRangeSelectChange(-1)
}

func (self *ListController) HandleClick(opts gocui.ViewMouseBindingOpts) error {
	prevSelectedLineIdx := self.context.GetList().GetSelectedLineIdx()
	newSelectedLineIdx := self.context.ViewIndexToModelIndex(opts.Y)
	alreadyFocused := self.isFocused()

	if err := self.pushContextIfNotFocused(); err != nil {
		return err
	}

	if newSelectedLineIdx > self.context.GetList().Len()-1 {
		return nil
	}

	self.context.GetList().SetSelection(newSelectedLineIdx)

	if prevSelectedLineIdx == newSelectedLineIdx && alreadyFocused && self.context.GetOnClick() != nil {
		return self.context.GetOnClick()()
	}
	return self.context.HandleFocus(types.OnFocusOpts{})
}

func (self *ListController) pushContextIfNotFocused() error {
	if !self.isFocused() {
		if err := self.c.PushContext(self.context); err != nil {
			return err
		}
	}

	return nil
}

func (self *ListController) isFocused() bool {
	return self.c.CurrentContext().GetKey() == self.context.GetKey()
}

func (self *ListController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{Tag: "navigation", Key: opts.GetKey(opts.Config.Universal.PrevItemAlt), Handler: self.HandlePrevLine},
		{Tag: "navigation", Key: opts.GetKey(opts.Config.Universal.PrevItem), Handler: self.HandlePrevLine},
		{Tag: "navigation", Key: opts.GetKey(opts.Config.Universal.NextItemAlt), Handler: self.HandleNextLine},
		{Tag: "navigation", Key: opts.GetKey(opts.Config.Universal.NextItem), Handler: self.HandleNextLine},
		{Tag: "navigation", Key: opts.GetKey(opts.Config.Universal.PrevPage), Handler: self.HandlePrevPage, Description: self.c.Tr.PrevPage},
		{Tag: "navigation", Key: opts.GetKey(opts.Config.Universal.NextPage), Handler: self.HandleNextPage, Description: self.c.Tr.NextPage},
		{Tag: "navigation", Key: opts.GetKey(opts.Config.Universal.GotoTop), Handler: self.HandleGotoTop, Description: self.c.Tr.GotoTop},
		{Tag: "navigation", Key: opts.GetKey(opts.Config.Universal.ScrollLeft), Handler: self.HandleScrollLeft},
		{Tag: "navigation", Key: opts.GetKey(opts.Config.Universal.ScrollRight), Handler: self.HandleScrollRight},
		{Tag: "navigation", Key: opts.GetKey(opts.Config.Universal.GotoBottom), Handler: self.HandleGotoBottom, Description: self.c.Tr.GotoBottom},
	}

	if self.context.RangeSelectEnabled() {
		bindings = append(bindings,
			[]*types.Binding{
				{Tag: "navigation", Key: opts.GetKey(opts.Config.Universal.ToggleRangeSelect), Handler: self.HandleToggleRangeSelect, Description: self.c.Tr.ToggleRangeSelect},
				{Tag: "navigation", Key: opts.GetKey(opts.Config.Universal.RangeSelectDown), Handler: self.HandleRangeSelectDown, Description: self.c.Tr.RangeSelectDown},
				{Tag: "navigation", Key: opts.GetKey(opts.Config.Universal.RangeSelectUp), Handler: self.HandleRangeSelectUp, Description: self.c.Tr.RangeSelectUp},
			}...,
		)
	}

	return bindings
}

func (self *ListController) GetMouseKeybindings(opts types.KeybindingsOpts) []*gocui.ViewMouseBinding {
	return []*gocui.ViewMouseBinding{
		{
			ViewName: self.context.GetViewName(),
			Key:      gocui.MouseWheelUp,
			Handler:  func(gocui.ViewMouseBindingOpts) error { return self.HandleScrollUp() },
		},
		{
			ViewName: self.context.GetViewName(),
			Key:      gocui.MouseLeft,
			Handler:  func(opts gocui.ViewMouseBindingOpts) error { return self.HandleClick(opts) },
		},
		{
			ViewName: self.context.GetViewName(),
			Key:      gocui.MouseWheelDown,
			Handler:  func(gocui.ViewMouseBindingOpts) error { return self.HandleScrollDown() },
		},
	}
}
