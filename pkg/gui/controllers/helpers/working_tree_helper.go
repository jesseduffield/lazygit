package helpers

import (
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type IWorkingTreeHelper interface {
	AnyStagedFiles() bool
	AnyTrackedFiles() bool
	IsWorkingTreeDirty() bool
	FileForSubmodule(submodule *models.SubmoduleConfig) *models.File
}

type WorkingTreeHelper struct {
	c   *types.HelperCommon
	git *commands.GitCommand

	model *types.Model
}

func NewWorkingTreeHelper(c *types.HelperCommon, git *commands.GitCommand, model *types.Model) *WorkingTreeHelper {
	return &WorkingTreeHelper{
		c:     c,
		git:   git,
		model: model,
	}
}

func (self *WorkingTreeHelper) AnyStagedFiles() bool {
	for _, file := range self.model.Files {
		if file.HasStagedChanges {
			return true
		}
	}
	return false
}

func (self *WorkingTreeHelper) AnyTrackedFiles() bool {
	for _, file := range self.model.Files {
		if file.Tracked {
			return true
		}
	}
	return false
}

func (self *WorkingTreeHelper) IsWorkingTreeDirty() bool {
	return self.AnyStagedFiles() || self.AnyTrackedFiles()
}

func (self *WorkingTreeHelper) FileForSubmodule(submodule *models.SubmoduleConfig) *models.File {
	for _, file := range self.model.Files {
		if file.IsSubmodule([]*models.SubmoduleConfig{submodule}) {
			return file
		}
	}

	return nil
}

func (self *WorkingTreeHelper) OpenMergeTool() error {
	return self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.MergeToolTitle,
		Prompt: self.c.Tr.MergeToolPrompt,
		HandleConfirm: func() error {
			self.c.LogAction(self.c.Tr.Actions.OpenMergeTool)
			return self.c.RunSubprocessAndRefresh(
				self.git.WorkingTree.OpenMergeToolCmdObj(),
			)
		},
	})
}
