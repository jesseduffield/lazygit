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
	self.c.PostRefreshUpdate(self.c.Context().Current())
	return nil
}
