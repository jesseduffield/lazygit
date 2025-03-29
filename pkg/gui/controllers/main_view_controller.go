package controllers

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type MainViewController struct {
	baseController
	c *ControllerCommon

	context      types.Context
	otherContext types.Context
}

var _ types.IController = &MainViewController{}

func NewMainViewController(
	c *ControllerCommon,
	context types.Context,
	otherContext types.Context,
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
			Key:             opts.GetKey(opts.Config.Universal.TogglePanel),
			Handler:         self.togglePanel,
			Description:     self.c.Tr.ToggleStagingView,
			Tooltip:         self.c.Tr.ToggleStagingViewTooltip,
			DisplayOnScreen: true,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Return),
			Handler:     self.escape,
			Description: self.c.Tr.ExitFocusedMainView,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.GoInto),
			Handler:     self.enter,
			Description: self.c.Tr.EnterStaging,
		},
	}
}

func (self *MainViewController) GetMouseKeybindings(opts types.KeybindingsOpts) []*gocui.ViewMouseBinding {
	return []*gocui.ViewMouseBinding{
		{
			ViewName: self.context.GetViewName(),
			Key:      gocui.MouseLeft,
			Handler: func(opts gocui.ViewMouseBindingOpts) error {
				if self.isFocused() {
					return self.onClick(opts)
				}

				self.context.SetParentContext(self.otherContext.GetParentContext())
				self.c.Context().Push(self.context, types.OnFocusOpts{
					ClickedWindowName:  self.context.GetWindowName(),
					ClickedViewLineIdx: opts.Y,
				})

				return nil
			},
		},
	}
}

func (self *MainViewController) Context() types.Context {
	return self.context
}

func (self *MainViewController) GetOnFocus() func(types.OnFocusOpts) {
	return func(opts types.OnFocusOpts) {
		if opts.ClickedWindowName != "" {
			self.context.GetView().FocusPoint(0, opts.ClickedViewLineIdx)
		}
	}
}

func (self *MainViewController) togglePanel() error {
	if self.otherContext.GetView().Visible {
		self.otherContext.SetParentContext(self.context.GetParentContext())
		self.c.Context().Push(self.otherContext, types.OnFocusOpts{})
	}

	return nil
}

func (self *MainViewController) escape() error {
	self.c.Context().Pop()
	return nil
}

func (self *MainViewController) enter() error {
	parentCtx := self.context.GetParentContext()
	if parentCtx.GetOnClickFocusedMainView() != nil {
		return parentCtx.GetOnClickFocusedMainView()(self.context.GetViewName(), self.context.GetView().SelectedLineIdx())
	}
	return nil
}

func (self *MainViewController) onClick(opts gocui.ViewMouseBindingOpts) error {
	if opts.Y != opts.PreviousY {
		return nil
	}

	parentCtx := self.context.GetParentContext()
	if parentCtx.GetOnClickFocusedMainView() != nil {
		return parentCtx.GetOnClickFocusedMainView()(self.context.GetViewName(), opts.Y)
	}
	return nil
}

func (self *MainViewController) isFocused() bool {
	return self.c.Context().Current().GetKey() == self.context.GetKey()
}
