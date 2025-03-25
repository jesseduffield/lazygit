package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// This controller is for all contexts that can focus their main view.

var _ types.IController = &SwitchToFocusedMainViewController{}

type SwitchToFocusedMainViewController struct {
	baseController
	c       *ControllerCommon
	context types.Context
}

func NewSwitchToFocusedMainViewController(
	c *ControllerCommon,
	context types.Context,
) *SwitchToFocusedMainViewController {
	return &SwitchToFocusedMainViewController{
		baseController: baseController{},
		c:              c,
		context:        context,
	}
}

func (self *SwitchToFocusedMainViewController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Universal.FocusMainView),
			Handler:     self.handleFocusMainView,
			Description: self.c.Tr.FocusMainView,
		},
	}

	return bindings
}

func (self *SwitchToFocusedMainViewController) Context() types.Context {
	return self.context
}

func (self *SwitchToFocusedMainViewController) handleFocusMainView() error {
	mainViewContext := self.c.Helpers().Window.GetContextForWindow("main")
	mainViewContext.SetParentContext(self.context)
	self.c.Context().Push(mainViewContext, types.OnFocusOpts{})
	return nil
}
