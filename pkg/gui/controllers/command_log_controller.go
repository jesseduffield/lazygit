package controllers

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type CommandLogController struct {
	baseController
	c *ControllerCommon
}

var _ types.IController = &CommandLogController{}

func NewCommandLogController(
	c *ControllerCommon,
) *CommandLogController {
	return &CommandLogController{
		baseController: baseController{},
		c:              c,
	}
}

func (self *CommandLogController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{}

	return bindings
}

func (self *CommandLogController) GetMouseKeybindings(opts types.KeybindingsOpts) []*gocui.ViewMouseBinding {
	return []*gocui.ViewMouseBinding{
		{
			ViewName: self.Context().GetViewName(),
			Key:      gocui.MouseLeft,
			Handler:  self.onClick,
		},
	}
}

func (self *CommandLogController) onClick(opts gocui.ViewMouseBindingOpts) error {
	if err := self.c.PushContext(self.context(), types.OnFocusOpts{
		ClickedWindowName:  self.context().GetWindowName(),
		ClickedViewLineIdx: opts.Y,
	}); err != nil {
		return err
	}
	return self.c.HandleGenericClick(self.c.Views().Extras)
}

func (self *CommandLogController) GetOnFocusLost() func(types.OnFocusLostOpts) error {
	return func(types.OnFocusLostOpts) error {
		self.c.Views().Extras.Autoscroll = true
		return nil
	}
}

func (self *CommandLogController) Context() types.Context {
	return self.context()
}

func (self *CommandLogController) context() types.Context {
	return self.c.Contexts().CommandLog
}
