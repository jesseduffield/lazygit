package gui

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
)

type WorkingTreeHelper struct {
	getFiles func() []*models.File
}

func NewWorkingTreeHelper(getFiles func() []*models.File) *WorkingTreeHelper {
	return &WorkingTreeHelper{
		getFiles: getFiles,
	}
}

func (self *WorkingTreeHelper) AnyStagedFiles() bool {
	files := self.getFiles()
	for _, file := range files {
		if file.HasStagedChanges {
			return true
		}
	}
	return false
}

func (self *WorkingTreeHelper) AnyTrackedFiles() bool {
	files := self.getFiles()
	for _, file := range files {
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
	for _, file := range self.getFiles() {
		if file.IsSubmodule([]*models.SubmoduleConfig{submodule}) {
			return file
		}
	}

	return nil
}
