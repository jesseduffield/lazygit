package gui

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
)

type WorkingTreeHelper struct {
	getFileTreeViewModel func() *filetree.FileTreeViewModel
}

func NewWorkingTreeHelper(getFileTreeViewModel func() *filetree.FileTreeViewModel) *WorkingTreeHelper {
	return &WorkingTreeHelper{
		getFileTreeViewModel: getFileTreeViewModel,
	}
}

func (self *WorkingTreeHelper) AnyStagedFiles() bool {
	files := self.getFileTreeViewModel().GetAllFiles()
	for _, file := range files {
		if file.HasStagedChanges {
			return true
		}
	}
	return false
}

func (self *WorkingTreeHelper) AnyTrackedFiles() bool {
	files := self.getFileTreeViewModel().GetAllFiles()
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
	for _, file := range self.getFileTreeViewModel().GetAllFiles() {
		if file.IsSubmodule([]*models.SubmoduleConfig{submodule}) {
			return file
		}
	}

	return nil
}
