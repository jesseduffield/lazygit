package controllers

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

type ViewSelectionControllerFactory struct {
	c *ControllerCommon
}

func NewViewSelectionControllerFactory(c *ControllerCommon) *ViewSelectionControllerFactory {
	return &ViewSelectionControllerFactory{
		c: c,
	}
}

func (self *ViewSelectionControllerFactory) Create(context types.Context) types.IController {
	return &ViewSelectionController{
		baseController: baseController{},
		c:              self.c,
		context:        context,
	}
}

type ViewSelectionController struct {
	baseController
	c *ControllerCommon

	context types.Context
}

func (self *ViewSelectionController) Context() types.Context {
	return self.context
}

func (self *ViewSelectionController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	return []*types.Binding{
		{Tag: "navigation", Key: opts.GetKey(opts.Config.Universal.PrevItem), Handler: self.handlePrevLine},
		{Tag: "navigation", Key: opts.GetKey(opts.Config.Universal.PrevItemAlt), Handler: self.handlePrevLine},
		{Tag: "navigation", Key: opts.GetKey(opts.Config.Universal.NextItem), Handler: self.handleNextLine},
		{Tag: "navigation", Key: opts.GetKey(opts.Config.Universal.NextItemAlt), Handler: self.handleNextLine},
		{Tag: "navigation", Key: opts.GetKey(opts.Config.Universal.PrevPage), Handler: self.handlePrevPage, Description: self.c.Tr.PrevPage},
		{Tag: "navigation", Key: opts.GetKey(opts.Config.Universal.NextPage), Handler: self.handleNextPage, Description: self.c.Tr.NextPage},
		{Tag: "navigation", Key: opts.GetKey(opts.Config.Universal.GotoTop), Handler: self.handleGotoTop, Description: self.c.Tr.GotoTop, Alternative: "<home>"},
		{Tag: "navigation", Key: opts.GetKey(opts.Config.Universal.GotoBottom), Handler: self.handleGotoBottom, Description: self.c.Tr.GotoBottom, Alternative: "<end>"},
		{Tag: "navigation", Key: opts.GetKey(opts.Config.Universal.GotoTopAlt), Handler: self.handleGotoTop},
		{Tag: "navigation", Key: opts.GetKey(opts.Config.Universal.GotoBottomAlt), Handler: self.handleGotoBottom},
	}
}

func (self *ViewSelectionController) GetMouseKeybindings(opts types.KeybindingsOpts) []*gocui.ViewMouseBinding {
	return []*gocui.ViewMouseBinding{}
}

func (self *ViewSelectionController) handleLineChange(delta int) {
	if delta > 0 {
		if manager := self.c.GetViewBufferManagerForView(self.context.GetView()); manager != nil {
			manager.ReadLines(delta)
		}
	}

	v := self.Context().GetView()
	if self.context.GetView().Highlight {
		lineIdxBefore := v.CursorY() + v.OriginY()
		lineIdxAfter := lo.Clamp(lineIdxBefore+delta, 0, v.ViewLinesHeight()-1)
		if delta == -1 {
			checkScrollUp(self.Context().GetViewTrait(), self.c.UserConfig(), lineIdxBefore, lineIdxAfter)
		} else if delta == 1 {
			checkScrollDown(self.Context().GetViewTrait(), self.c.UserConfig(), lineIdxBefore, lineIdxAfter)
		}
		v.FocusPoint(0, lineIdxAfter)
	} else {
		if delta < 0 {
			v.ScrollUp(-delta)
		} else {
			v.ScrollDown(delta)
		}
	}
}

func (self *ViewSelectionController) handlePrevLine() error {
	self.handleLineChange(-1)
	return nil
}

func (self *ViewSelectionController) handleNextLine() error {
	self.handleLineChange(1)
	return nil
}

func (self *ViewSelectionController) handlePrevPage() error {
	self.handleLineChange(-self.context.GetViewTrait().PageDelta())
	return nil
}

func (self *ViewSelectionController) handleNextPage() error {
	self.handleLineChange(self.context.GetViewTrait().PageDelta())
	return nil
}

func (self *ViewSelectionController) handleGotoTop() error {
	v := self.Context().GetView()
	if self.context.GetView().Highlight {
		v.FocusPoint(0, 0)
	} else {
		self.handleLineChange(-v.ViewLinesHeight())
	}
	return nil
}

func (self *ViewSelectionController) handleGotoBottom() error {
	if manager := self.c.GetViewBufferManagerForView(self.context.GetView()); manager != nil {
		manager.ReadToEnd(func() {
			self.c.OnUIThread(func() error {
				v := self.Context().GetView()
				if self.context.GetView().Highlight {
					v.FocusPoint(0, v.ViewLinesHeight()-1)
				} else {
					self.handleLineChange(v.ViewLinesHeight())
				}
				return nil
			})
		})
	}

	return nil
}
