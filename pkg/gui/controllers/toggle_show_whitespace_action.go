package controllers

import "github.com/jesseduffield/lazygit/pkg/gui/types"

type ToggleShowWhitespaceAction struct {
	c *ControllerCommon
}

func (self *ToggleShowWhitespaceAction) Call() error {
	self.c.UserConfig().Gui.ShowWhitespace = !self.c.UserConfig().Gui.ShowWhitespace

	if self.c.UserConfig().Gui.ShowWhitespace {
		self.c.Toast(self.c.Tr.ShowWhitespaceIndicatorOn)
	}

	self.c.Context().CurrentSide().HandleFocus(types.OnFocusOpts{})
	return nil
}
