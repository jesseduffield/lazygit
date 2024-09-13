package filetree

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/mergeconflicts"
)

// FileNode wraps a node and provides some file-specific methods for it.
type FileNode struct {
	*Node[models.File]
}

var _ models.IFile = &FileNode{}

func NewFileNode(node *Node[models.File]) *FileNode {
	if node == nil {
		return nil
	}

	return &FileNode{Node: node}
}

// returns the underlying node, without any file-specific methods attached
func (self *FileNode) Raw() *Node[models.File] {
	if self == nil {
		return nil
	}

	return self.Node
}

func (self *FileNode) GetHasUnstagedChanges() bool {
	return self.SomeFile(func(file *models.File) bool { return file.HasUnstagedChanges })
}

func (self *FileNode) GetHasStagedOrTrackedChanges() bool {
	if !self.GetHasStagedChanges() {
		return self.SomeFile(func(t *models.File) bool {
			return t.Tracked
		})
	}
	return true
}

func (self *FileNode) GetHasStagedChanges() bool {
	return self.SomeFile(func(file *models.File) bool { return file.HasStagedChanges })
}

func (self *FileNode) GetHasInlineMergeConflicts() bool {
	return self.SomeFile(func(file *models.File) bool {
		if !file.HasInlineMergeConflicts {
			return false
		}
		hasConflicts, _ := mergeconflicts.FileHasConflictMarkers(file.Name)
		return hasConflicts
	})
}

func (self *FileNode) GetIsTracked() bool {
	return self.SomeFile(func(file *models.File) bool { return file.Tracked })
}

func (self *FileNode) GetIsFile() bool {
	return self.IsFile()
}

func (self *FileNode) GetPreviousPath() string {
	if self.File == nil {
		return ""
	}

	return self.File.PreviousName
}
