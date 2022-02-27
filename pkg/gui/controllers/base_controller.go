package controllers

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type baseController struct{}

func (self *baseController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	return nil
}

func (self *baseController) GetMouseKeybindings(opts types.KeybindingsOpts) []*gocui.ViewMouseBinding {
	return nil
}

func (self *baseController) GetOnClick() func() error {
	return nil
}
