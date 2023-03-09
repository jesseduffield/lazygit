package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type MenuController struct {
	baseController
	*controllerCommon
}

var _ types.IController = &MenuController{}

func NewMenuController(
	common *controllerCommon,
) *MenuController {
	return &MenuController{
		baseController:   baseController{},
		controllerCommon: common,
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
			Key:     opts.GetKey(opts.Config.Universal.Return),
			Handler: self.close,
		},
	}

	return bindings
}

func (self *MenuController) GetOnClick() func() error {
	return self.press
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
	return self.contexts.Menu
}
