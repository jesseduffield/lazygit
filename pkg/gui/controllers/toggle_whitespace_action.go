package controllers

type ToggleWhitespaceAction struct {
	c *ControllerCommon
}

func (self *ToggleWhitespaceAction) Call() error {
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
