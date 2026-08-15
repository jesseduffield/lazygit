package controllers

import (
	"errors"

	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
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

	if lo.Contains(contextsThatDontSupportIgnoringWhitespace, self.c.Context().Current().GetKey()) {
		// Ignoring whitespace is not supported in these views. Let the user
		// know that it's not going to work in case they try to turn it on.
		return errors.New(self.c.Tr.IgnoreWhitespaceNotSupportedHere)
	}

	self.c.UserConfig().Git.IgnoreWhitespaceInDiffView = !self.c.UserConfig().Git.IgnoreWhitespaceInDiffView

	// You toggle this to see whether what you are looking at is more than
	// reindentation, so that is the thing to keep in front of you — even though
	// ignoring whitespace, unlike the other ways of re-rendering a diff, can take
	// the line away entirely along with the hunk or file it was in.
	self.c.Helpers().DiffLine.PreserveDiffPositionOnRerender(self.c.Contexts().Normal.GetView())
	self.c.Helpers().DiffLine.PreserveDiffPositionOnRerender(self.c.Contexts().NormalSecondary.GetView())
	self.c.Context().CurrentSide().HandleRenderToMain()
	return nil
}
