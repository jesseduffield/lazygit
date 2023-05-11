package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type ToggleWhitespaceAction struct {
	c *ControllerCommon
}

func (self *ToggleWhitespaceAction) Call() error {
	self.c.State().SetIgnoreWhitespaceInDiffView(!self.c.State().GetIgnoreWhitespaceInDiffView())

	toastMessage := self.c.Tr.ShowingWhitespaceInDiffView
	if self.c.State().GetIgnoreWhitespaceInDiffView() {
		toastMessage = self.c.Tr.IgnoringWhitespaceInDiffView
	}
	self.c.Toast(toastMessage)

	return self.c.CurrentSideContext().HandleFocus(types.OnFocusOpts{})
}
