package gui

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
)

type WorkingTreeHelper struct {
	fileTreeViewModel *filetree.FileTreeViewModel
}

func NewWorkingTreeHelper(fileTreeViewModel *filetree.FileTreeViewModel) *WorkingTreeHelper {
	return &WorkingTreeHelper{
		fileTreeViewModel: fileTreeViewModel,
	}
}

func (self *WorkingTreeHelper) AnyStagedFiles() bool {
	files := self.fileTreeViewModel.GetAllFiles()
	for _, file := range files {
		if file.HasStagedChanges {
			return true
		}
	}
	return false
}

func (self *WorkingTreeHelper) AnyTrackedFiles() bool {
	files := self.fileTreeViewModel.GetAllFiles()
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
	for _, file := range self.fileTreeViewModel.GetAllFiles() {
		if file.IsSubmodule([]*models.SubmoduleConfig{submodule}) {
			return file
		}
	}

	return nil
}
