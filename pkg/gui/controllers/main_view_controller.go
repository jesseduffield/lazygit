package controllers

import (
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
	}
}

func (self *MainViewController) Context() types.Context {
	return self.context
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
