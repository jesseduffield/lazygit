package controllers

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type MenuController struct {
	c       *types.ControllerCommon
	context types.IListContext

	getSelectedMenuItem func() *types.MenuItem
}

var _ types.IController = &MenuController{}

func NewMenuController(
	c *types.ControllerCommon,
	context types.IListContext,
	getSelectedMenuItem func() *types.MenuItem,
) *MenuController {
	return &MenuController{
		c:                   c,
		context:             context,
		getSelectedMenuItem: getSelectedMenuItem,
	}
}

func (self *MenuController) Keybindings(getKey func(key string) interface{}, config config.KeybindingConfig, guards types.KeybindingGuards) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:     getKey(config.Universal.Select),
			Handler: self.press,
		},
		{
			Key:     getKey(config.Universal.Confirm),
			Handler: self.press,
		},
		{
			Key:     getKey(config.Universal.ConfirmAlt1),
			Handler: self.press,
		},
		{
			Key:     gocui.MouseLeft,
			Handler: func() error { return self.context.HandleClick(self.press) },
		},
	}

	return append(bindings, self.context.Keybindings(getKey, config, guards)...)
}

func (self *MenuController) press() error {
	selectedItem := self.getSelectedMenuItem()

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
