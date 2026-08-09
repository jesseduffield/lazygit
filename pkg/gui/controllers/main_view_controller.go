package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type MainViewController struct {
	baseController
	c *ControllerCommon

	context      *context.MainContext
	otherContext *context.MainContext
}

var _ types.IController = &MainViewController{}

func NewMainViewController(
	c *ControllerCommon,
	context *context.MainContext,
	otherContext *context.MainContext,
) *MainViewController {
	return &MainViewController{
		baseController: baseController{},
		c:              c,
		context:        context,
		otherContext:   otherContext,
	}
}

func (self *MainViewController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	return []*types.Binding{
		{
			Keys:            opts.GetKeys(opts.Config.Universal.TogglePanel),
			Handler:         self.togglePanel,
			Description:     self.c.Tr.ToggleStagingView,
			Tooltip:         self.c.Tr.ToggleStagingViewTooltip,
			DisplayOnScreen: true,
		},
		{
			Keys:            opts.GetKeys(opts.Config.Universal.Return),
			Handler:         self.escape,
			Description:     self.c.Tr.ExitFocusedMainView,
			DisplayOnScreen: true,
		},
		{
			// overriding this because we want to read all of the task's output before we start searching
			Keys:        opts.GetKeys(opts.Config.Universal.StartSearch),
			Handler:     self.openSearch,
			Description: self.c.Tr.StartSearch,
			Tag:         "navigation",
		},
		{Tag: "navigation", Keys: opts.GetKeys(opts.Config.Universal.PrevItem), Handler: self.handlePrevLine},
		{Tag: "navigation", Keys: opts.GetKeys(opts.Config.Universal.NextItem), Handler: self.handleNextLine},
		{Tag: "navigation", Keys: opts.GetKeys(opts.Config.Universal.PrevPage), Handler: self.handlePrevPage, Description: self.c.Tr.PrevPage},
		{Tag: "navigation", Keys: opts.GetKeys(opts.Config.Universal.NextPage), Handler: self.handleNextPage, Description: self.c.Tr.NextPage},
		{Tag: "navigation", Keys: opts.GetKeys(opts.Config.Universal.GotoTop), Handler: self.handleGotoTop, Description: self.c.Tr.GotoTop},
		{Tag: "navigation", Keys: opts.GetKeys(opts.Config.Universal.GotoBottom), Handler: self.handleGotoBottom, Description: self.c.Tr.GotoBottom},
	}
}

func (self *MainViewController) GetMouseKeybindings(opts types.KeybindingsOpts) []*gocui.ViewMouseBinding {
	return []*gocui.ViewMouseBinding{
		{
			ViewName:    self.context.GetViewName(),
			Key:         gocui.MouseLeft,
			Handler:     self.onClickInAlreadyFocusedView,
			FocusedView: self.context.GetViewName(),
		},
		{
			ViewName:    self.context.GetViewName(),
			Key:         gocui.MouseLeft,
			Handler:     self.onClickInOtherViewOfMainViewPair,
			FocusedView: self.otherContext.GetViewName(),
		},
	}
}

func (self *MainViewController) Context() types.Context {
	return self.context
}

func (self *MainViewController) togglePanel() error {
	if self.otherContext.GetView().Visible {
		self.c.Context().Push(self.otherContext, types.OnFocusOpts{})
	}

	return nil
}

func (self *MainViewController) escape() error {
	self.c.Context().Pop()
	return nil
}

func (self *MainViewController) onClickInAlreadyFocusedView(opts gocui.ViewMouseBindingOpts) error {
	sidePanelContext := self.c.Context().NextInStack(self.context)
	if sidePanelContext != nil && sidePanelContext.GetOnClickFocusedMainView() != nil {
		return sidePanelContext.GetOnClickFocusedMainView()(self.context.GetViewName(), opts.Y)
	}
	return nil
}

func (self *MainViewController) onClickInOtherViewOfMainViewPair(opts gocui.ViewMouseBindingOpts) error {
	self.c.Context().Push(self.context, types.OnFocusOpts{
		ClickedWindowName:  self.context.GetWindowName(),
		ClickedViewLineIdx: opts.Y,
	})

	return nil
}

func (self *MainViewController) handleLineChange(delta int) {
	v := self.context.GetView()
	if delta < 0 {
		v.ScrollUp(-delta)
	} else {
		v.ScrollDown(delta)
		self.c.ReadLinesToFillView(v)
	}
}

func (self *MainViewController) handlePrevLine() error {
	self.handleLineChange(-1)
	return nil
}

func (self *MainViewController) handleNextLine() error {
	self.handleLineChange(1)
	return nil
}

func (self *MainViewController) handlePrevPage() error {
	self.handleLineChange(-self.context.GetViewTrait().PageDelta())
	return nil
}

func (self *MainViewController) handleNextPage() error {
	self.handleLineChange(self.context.GetViewTrait().PageDelta())
	return nil
}

func (self *MainViewController) handleGotoTop() error {
	v := self.context.GetView()
	self.handleLineChange(-v.ViewLinesHeight())
	return nil
}

func (self *MainViewController) handleGotoBottom() error {
	if manager := self.c.GetViewBufferManagerForView(self.context.GetView()); manager != nil {
		manager.ReadToEnd(func() {
			self.c.OnUIThread(func() error {
				v := self.context.GetView()
				self.handleLineChange(v.ViewLinesHeight())
				return nil
			})
		})
	}

	return nil
}

func (self *MainViewController) openSearch() error {
	if manager := self.c.GetViewBufferManagerForView(self.context.GetView()); manager != nil {
		manager.ReadToEnd(func() {
			self.c.OnUIThread(func() error {
				self.c.Helpers().Search.OpenSearchPrompt(self.context)
				return nil
			})
		})
	}

	return nil
}
