package controllers

import (
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

func (self *CommandLogController) GetOnFocusLost() func(types.OnFocusLostOpts) {
	return func(types.OnFocusLostOpts) {
		self.c.Views().Extras.Autoscroll = true
	}
}

func (self *CommandLogController) Context() types.Context {
	return self.context()
}

func (self *CommandLogController) context() types.Context {
	return self.c.Contexts().CommandLog
}
