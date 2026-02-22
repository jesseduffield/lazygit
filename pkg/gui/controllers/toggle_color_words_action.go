package controllers

import (
	"errors"

	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

type ToggleColorWordsAction struct {
	c *ControllerCommon
}

func (self *ToggleColorWordsAction) Call() error {
	contextsThatDontSupportColorWords := []types.ContextKey{
		context.STAGING_MAIN_CONTEXT_KEY,
		context.STAGING_SECONDARY_CONTEXT_KEY,
		context.PATCH_BUILDING_MAIN_CONTEXT_KEY,
	}

	if lo.Contains(contextsThatDontSupportColorWords, self.c.Context().Current().GetKey()) {
		return errors.New(self.c.Tr.ColorWordsNotSupportedHere)
	}

	self.c.UserConfig().Git.ColorWordsInDiffView = !self.c.UserConfig().Git.ColorWordsInDiffView

	self.c.Context().CurrentSide().HandleFocus(types.OnFocusOpts{})
	return nil
}
