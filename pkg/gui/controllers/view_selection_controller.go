package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gocui"
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
		{Tag: "navigation", Keys: opts.GetKeys(opts.Config.Universal.PrevItem), Handler: self.handlePrevLine},
		{Tag: "navigation", Keys: opts.GetKeys(opts.Config.Universal.NextItem), Handler: self.handleNextLine},
		{Tag: "navigation", Keys: opts.GetKeys(opts.Config.Universal.PrevPage), Handler: self.handlePrevPage, Description: self.c.Tr.PrevPage},
		{Tag: "navigation", Keys: opts.GetKeys(opts.Config.Universal.NextPage), Handler: self.handleNextPage, Description: self.c.Tr.NextPage},
		{Tag: "navigation", Keys: opts.GetKeys(opts.Config.Universal.GotoTop), Handler: self.handleGotoTop, Description: self.c.Tr.GotoTop},
		{Tag: "navigation", Keys: opts.GetKeys(opts.Config.Universal.GotoBottom), Handler: self.handleGotoBottom, Description: self.c.Tr.GotoBottom},
	}
}

func (self *ViewSelectionController) GetMouseKeybindings(opts types.KeybindingsOpts) []*gocui.ViewMouseBinding {
	return []*gocui.ViewMouseBinding{}
}

func (self *ViewSelectionController) handleLineChange(delta int) {
	v := self.Context().GetView()
	if self.context.GetView().Highlight {
		lineIdxBefore := v.CursorY() + v.OriginY()
		lineIdxAfter := lo.Clamp(lineIdxBefore+delta, 0, v.ViewLinesHeight()-1)
		if delta == -1 {
			checkScrollUp(self.Context().GetViewTrait(), self.c.UserConfig(), lineIdxBefore, lineIdxAfter)
		} else if delta == 1 {
			checkScrollDown(self.Context().GetViewTrait(), self.c.UserConfig(), lineIdxBefore, lineIdxAfter)
		}
		v.FocusPoint(0, lineIdxAfter, true)
	} else {
		if delta < 0 {
			v.ScrollUp(-delta)
		} else {
			v.ScrollDown(delta)
			self.c.ReadLinesToFillView(v)
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
		v.FocusPoint(0, 0, true)
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
					v.FocusPoint(0, v.ViewLinesHeight()-1, true)
				} else {
					self.handleLineChange(v.ViewLinesHeight())
				}
				return nil
			})
		})
	}

	return nil
}
