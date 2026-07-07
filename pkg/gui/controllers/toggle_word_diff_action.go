package controllers

import (
	"errors"

	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

type ToggleWordDiffAction struct {
	c *ControllerCommon
}

func (self *ToggleWordDiffAction) Call() error {
	if lo.Contains(contextsThatDontSupportDiffViewOptions, self.c.Context().Current().GetKey()) {
		// Word-diff collapses the +/- lines that line-by-line staging needs, so
		// it can't be applied in these views. Let the user know rather than
		// silently doing nothing.
		return errors.New(self.c.Tr.WordDiffNotSupportedHere)
	}

	self.c.UserConfig().Git.WordDiffInDiffView = !self.c.UserConfig().Git.WordDiffInDiffView

	self.c.Context().CurrentSide().HandleFocus(types.OnFocusOpts{})
	return nil
}
