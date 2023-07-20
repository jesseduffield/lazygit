package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type MenuController struct {
	baseController
	c *ControllerCommon
}

var _ types.IController = &MenuController{}

func NewMenuController(
	common *ControllerCommon,
) *MenuController {
	return &MenuController{
		baseController: baseController{},
		c:              common,
	}
}

// NOTE: if you add a new keybinding here, you'll also need to add it to
// `reservedKeys` in `pkg/gui/context/menu_context.go`
func (self *MenuController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:     opts.GetKey(opts.Config.Universal.Select),
			Handler: self.press,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Confirm),
			Handler:     self.press,
			Description: self.c.Tr.Execute,
			Display:     true,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Return),
			Handler:     self.close,
			Description: self.c.Tr.Close,
			Display:     true,
		},
	}

	return bindings
}

func (self *MenuController) GetOnClick() func() error {
	return self.press
}

func (self *MenuController) GetOnFocus() func(types.OnFocusOpts) error {
	return func(types.OnFocusOpts) error {
		selectedMenuItem := self.context().GetSelected()
		if selectedMenuItem != nil {
			self.c.Views().Tooltip.SetContent(selectedMenuItem.Tooltip)
		}
		return nil
	}
}

func (self *MenuController) press() error {
	return self.context().OnMenuPress(self.context().GetSelected())
}

func (self *MenuController) close() error {
	return self.c.PopContext()
}

func (self *MenuController) Context() types.Context {
	return self.context()
}

func (self *MenuController) context() *context.MenuContext {
	return self.c.Contexts().Menu
}
