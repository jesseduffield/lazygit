package helpers

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type CommitFilesHelper struct {
	c *HelperCommon
}

func NewCommitFilesHelper(c *HelperCommon) *CommitFilesHelper {
	return &CommitFilesHelper{
		c: c,
	}
}

func (self *CommitFilesHelper) EnterCommitFile(node *filetree.CommitFileNode, opts types.OnFocusOpts) error {
	if node.File == nil {
		self.handleToggleCommitFileDirCollapsed(node)
		return nil
	}

	if self.c.AppState.DiffContextSize == 0 {
		return fmt.Errorf(self.c.Tr.Actions.NotEnoughContextToStage,
			keybindings.Label(self.c.UserConfig().Keybinding.Universal.IncreaseContextInDiffView))
	}

	enterTheFile := func() error {
		if !self.c.Git().Patch.PatchBuilder.Active() {
			if err := self.StartPatchBuilder(); err != nil {
				return err
			}
		}

		self.c.Context().Push(self.c.Contexts().CustomPatchBuilder, opts)
		return nil
	}

	from, to, reverse := self.CurrentFromToReverseForPatchBuilding()
	if self.c.Git().Patch.PatchBuilder.Active() && self.c.Git().Patch.PatchBuilder.NewPatchRequired(from, to, reverse) {
		self.c.Confirm(types.ConfirmOpts{
			Title:  self.c.Tr.DiscardPatch,
			Prompt: self.c.Tr.DiscardPatchConfirm,
			HandleConfirm: func() error {
				self.c.Git().Patch.PatchBuilder.Reset()
				return enterTheFile()
			},
		})

		return nil
	}

	return enterTheFile()
}

func (self *CommitFilesHelper) context() *context.CommitFilesContext {
	return self.c.Contexts().CommitFiles
}

func (self *CommitFilesHelper) handleToggleCommitFileDirCollapsed(node *filetree.CommitFileNode) {
	self.context().CommitFileTreeViewModel.ToggleCollapsed(node.GetInternalPath())

	self.c.PostRefreshUpdate(self.context())
}

func (self *CommitFilesHelper) StartPatchBuilder() error {
	commitFilesContext := self.context()

	canRebase := commitFilesContext.GetCanRebase()
	from, to, reverse := self.CurrentFromToReverseForPatchBuilding()

	self.c.Git().Patch.PatchBuilder.Start(from, to, reverse, canRebase)
	return nil
}

func (self *CommitFilesHelper) CurrentFromToReverseForPatchBuilding() (string, string, bool) {
	commitFilesContext := self.context()

	from, to := commitFilesContext.GetFromAndToForDiff()
	from, reverse := self.c.Modes().Diffing.GetFromAndReverseArgsForDiff(from)
	return from, to, reverse
}
