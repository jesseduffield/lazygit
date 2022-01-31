package controllers

import (
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
	model *types.Model
}

func NewWorkingTreeHelper(model *types.Model) *WorkingTreeHelper {
	return &WorkingTreeHelper{
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
