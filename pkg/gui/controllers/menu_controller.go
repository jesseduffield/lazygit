package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type MenuController struct {
	baseController

	c       *types.ControllerCommon
	context *context.MenuContext
}

var _ types.IController = &MenuController{}

func NewMenuController(
	c *types.ControllerCommon,
	context *context.MenuContext,
) *MenuController {
	return &MenuController{
		baseController: baseController{},

		c:       c,
		context: context,
	}
}

func (self *MenuController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:     opts.GetKey(opts.Config.Universal.Select),
			Handler: self.press,
		},
		{
			Key:     opts.GetKey(opts.Config.Universal.Confirm),
			Handler: self.press,
		},
		{
			Key:     opts.GetKey(opts.Config.Universal.ConfirmAlt1),
			Handler: self.press,
		},
		// {
		// 	Key:     gocui.MouseLeft,
		// 	Handler: func() error { return self.context.HandleClick(self.press) },
		// },
	}

	return bindings
}

func (self *MenuController) press() error {
	selectedItem := self.context.GetSelected()

	if err := self.c.PopContext(); err != nil {
		return err
	}

	if err := selectedItem.OnPress(); err != nil {
		return err
	}

	return nil
}

func (self *MenuController) Context() types.Context {
	return self.context
}
