package controllers

import (
	"errors"

	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

// These contexts render their diff for line-by-line staging, so diff-view
// rendering options like ignoring whitespace or word-diff can't be applied to
// them: those options rewrite the diff in ways that break line selection.
var contextsThatDontSupportDiffViewOptions = []types.ContextKey{
	context.STAGING_MAIN_CONTEXT_KEY,
	context.STAGING_SECONDARY_CONTEXT_KEY,
	context.PATCH_BUILDING_MAIN_CONTEXT_KEY,
}

type ToggleWhitespaceAction struct {
	c *ControllerCommon
}

func (self *ToggleWhitespaceAction) Call() error {
	if lo.Contains(contextsThatDontSupportDiffViewOptions, self.c.Context().Current().GetKey()) {
		// Ignoring whitespace is not supported in these views. Let the user
		// know that it's not going to work in case they try to turn it on.
		return errors.New(self.c.Tr.IgnoreWhitespaceNotSupportedHere)
	}

	self.c.UserConfig().Git.IgnoreWhitespaceInDiffView = !self.c.UserConfig().Git.IgnoreWhitespaceInDiffView

	self.c.Context().CurrentSide().HandleFocus(types.OnFocusOpts{})
	return nil
}
