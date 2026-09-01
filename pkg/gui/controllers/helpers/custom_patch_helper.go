package helpers

import "github.com/jesseduffield/lazygit/pkg/gui/types"

type CustomPatchHelper struct {
	c *HelperCommon
}

func NewCustomPatchHelper(c *HelperCommon) *CustomPatchHelper {
	return &CustomPatchHelper{c: c}
}

func (self *CustomPatchHelper) Reset() error {
	self.c.Git().Patch.PatchBuilder.Reset()
	self.c.Refresh(types.RefreshOptions{
		Scope: []types.RefreshableView{types.COMMIT_FILES},
	})

	// Render again so that the pane that was previewing the patch goes with it. The
	// panel asked to do that is the side panel rather than whichever context has the
	// focus. Both main panes are rendered by the panel beneath them, so a reset from
	// within the focused main view has to go through that panel too.
	self.c.PostRefreshUpdate(self.c.Context().CurrentSide())
	return nil
}
