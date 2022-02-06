package controllers

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type ListControllerFactory struct {
	c *types.HelperCommon
}

func NewListControllerFactory(c *types.HelperCommon) *ListControllerFactory {
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
	c *types.HelperCommon

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
	return self.scroll(self.context.GetViewTrait().ScrollLeft)
}

func (self *ListController) HandleScrollRight() error {
	return self.scroll(self.context.GetViewTrait().ScrollRight)
}

func (self *ListController) scroll(scrollFunc func()) error {
	scrollFunc()

	return self.context.HandleFocus()
}

func (self *ListController) handleLineChange(change int) error {
	before := self.context.GetList().GetSelectedLineIdx()
	self.context.GetList().MoveSelectedLine(change)
	after := self.context.GetList().GetSelectedLineIdx()

	// doing this check so that if we're holding the up key at the start of the list
	// we're not constantly re-rendering the main view.
	if before != after {
		return self.context.HandleFocus()
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
	return self.handleLineChange(-self.context.GetList().GetItemsLength())
}

func (self *ListController) HandleGotoBottom() error {
	return self.handleLineChange(self.context.GetList().GetItemsLength())
}

func (self *ListController) HandleClick(onClick func() error) error {
	prevSelectedLineIdx := self.context.GetList().GetSelectedLineIdx()
	// because we're handling a click, we need to determine the new line idx based
	// on the view itself.
	newSelectedLineIdx := self.context.GetViewTrait().SelectedLineIdx()

	currentContextKey := self.c.CurrentContext().GetKey()
	alreadyFocused := currentContextKey == self.context.GetKey()

	// we need to focus the view
	if !alreadyFocused {
		if err := self.c.PushContext(self.context); err != nil {
			return err
		}
	}

	if newSelectedLineIdx > self.context.GetList().GetItemsLength()-1 {
		return nil
	}

	self.context.GetList().SetSelectedLineIdx(newSelectedLineIdx)

	if prevSelectedLineIdx == newSelectedLineIdx && alreadyFocused && onClick != nil {
		return onClick()
	}
	return self.context.HandleFocus()
}

func (self *ListController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	return []*types.Binding{
		{Tag: "navigation", Key: opts.GetKey(opts.Config.Universal.PrevItemAlt), Modifier: gocui.ModNone, Handler: self.HandlePrevLine},
		{Tag: "navigation", Key: opts.GetKey(opts.Config.Universal.PrevItem), Modifier: gocui.ModNone, Handler: self.HandlePrevLine},
		{Tag: "navigation", Key: gocui.MouseWheelUp, Modifier: gocui.ModNone, Handler: self.HandlePrevLine},
		{Tag: "navigation", Key: opts.GetKey(opts.Config.Universal.NextItemAlt), Modifier: gocui.ModNone, Handler: self.HandleNextLine},
		{Tag: "navigation", Key: opts.GetKey(opts.Config.Universal.NextItem), Modifier: gocui.ModNone, Handler: self.HandleNextLine},
		{Tag: "navigation", Key: opts.GetKey(opts.Config.Universal.PrevPage), Modifier: gocui.ModNone, Handler: self.HandlePrevPage, Description: self.c.Tr.LcPrevPage},
		{Tag: "navigation", Key: opts.GetKey(opts.Config.Universal.NextPage), Modifier: gocui.ModNone, Handler: self.HandleNextPage, Description: self.c.Tr.LcNextPage},
		{Tag: "navigation", Key: opts.GetKey(opts.Config.Universal.GotoTop), Modifier: gocui.ModNone, Handler: self.HandleGotoTop, Description: self.c.Tr.LcGotoTop},
		{Key: gocui.MouseLeft, Modifier: gocui.ModNone, Handler: func() error { return self.HandleClick(nil) }},
		{Tag: "navigation", Key: gocui.MouseWheelDown, Modifier: gocui.ModNone, Handler: self.HandleNextLine},
		{Tag: "navigation", Key: opts.GetKey(opts.Config.Universal.ScrollLeft), Modifier: gocui.ModNone, Handler: self.HandleScrollLeft},
		{Tag: "navigation", Key: opts.GetKey(opts.Config.Universal.ScrollRight), Modifier: gocui.ModNone, Handler: self.HandleScrollRight},
		{
			Key:         opts.GetKey(opts.Config.Universal.StartSearch),
			Handler:     func() error { self.c.OpenSearch(); return nil },
			Description: self.c.Tr.LcStartSearch,
			Tag:         "navigation",
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.GotoBottom),
			Description: self.c.Tr.LcGotoBottom,
			Handler:     self.HandleGotoBottom,
			Tag:         "navigation",
		},
	}
}
