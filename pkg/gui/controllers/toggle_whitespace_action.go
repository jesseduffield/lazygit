package controllers

import (
	"github.com/lobes/lazytask/pkg/gui/context"
	"github.com/lobes/lazytask/pkg/gui/types"
	"github.com/samber/lo"
)

type ToggleWhitespaceAction struct {
	c *ControllerCommon
}

func (self *ToggleWhitespaceAction) Call() error {
	contextsThatDontSupportIgnoringWhitespace := []types.ContextKey{
		context.STAGING_MAIN_CONTEXT_KEY,
		context.STAGING_SECONDARY_CONTEXT_KEY,
		context.PATCH_BUILDING_MAIN_CONTEXT_KEY,
	}

	if lo.Contains(contextsThatDontSupportIgnoringWhitespace, self.c.CurrentContext().GetKey()) {
		// Ignoring whitespace is not supported in these views. Let the user
		// know that it's not going to work in case they try to turn it on.
		return self.c.ErrorMsg(self.c.Tr.IgnoreWhitespaceNotSupportedHere)
	}

	self.c.GetAppState().IgnoreWhitespaceInDiffView = !self.c.GetAppState().IgnoreWhitespaceInDiffView
	self.c.SaveAppStateAndLogError()

	return self.c.CurrentSideContext().HandleFocus(types.OnFocusOpts{})
}
