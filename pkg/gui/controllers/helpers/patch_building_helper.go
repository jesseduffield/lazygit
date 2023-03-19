package helpers

import (
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/types/enums"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type IPatchBuildingHelper interface {
	ValidateNormalWorkingTreeState() (bool, error)
}

type PatchBuildingHelper struct {
	c        *types.HelperCommon
	git      *commands.GitCommand
	contexts *context.ContextTree
}

func NewPatchBuildingHelper(
	c *types.HelperCommon,
	git *commands.GitCommand,
	contexts *context.ContextTree,
) *PatchBuildingHelper {
	return &PatchBuildingHelper{
		c:        c,
		git:      git,
		contexts: contexts,
	}
}

func (self *PatchBuildingHelper) ValidateNormalWorkingTreeState() (bool, error) {
	if self.git.Status.WorkingTreeState() != enums.REBASE_MODE_NONE {
		return false, self.c.ErrorMsg(self.c.Tr.CantPatchWhileRebasingError)
	}
	return true, nil
}

// takes us from the patch building panel back to the commit files panel
func (self *PatchBuildingHelper) Escape() error {
	return self.c.PopContext()
}

// kills the custom patch and returns us back to the commit files panel if needed
func (self *PatchBuildingHelper) Reset() error {
	self.git.Patch.PatchBuilder.Reset()

	if self.c.CurrentStaticContext().GetKind() != types.SIDE_CONTEXT {
		if err := self.Escape(); err != nil {
			return err
		}
	}

	if err := self.c.Refresh(types.RefreshOptions{
		Scope: []types.RefreshableView{types.COMMIT_FILES},
	}); err != nil {
		return err
	}

	// refreshing the current context so that the secondary panel is hidden if necessary.
	return self.c.PostRefreshUpdate(self.c.CurrentContext())
}
