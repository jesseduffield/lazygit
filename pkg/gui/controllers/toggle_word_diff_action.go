package controllers

import "github.com/jesseduffield/lazygit/pkg/gui/types"

type ToggleWordDiffAction struct {
	c *ControllerCommon
}

func (self *ToggleWordDiffAction) Call() error {
	self.c.UserConfig().Git.UseWordDiffInDiffView = !self.c.UserConfig().Git.UseWordDiffInDiffView

	self.c.Context().CurrentSide().HandleFocus(types.OnFocusOpts{})
	return nil
}
