package filetree

import (
	"sync"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context/traits"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/sirupsen/logrus"
)

type IFileTreeViewModel interface {
	IFileTree
	types.IListCursor
}

// This combines our FileTree struct with a cursor that retains information about
// which item is selected. It also contains logic for repositioning that cursor
// after the files are refreshed
type FileTreeViewModel struct {
	sync.RWMutex
	IFileTree
	types.IListCursor
}

var _ IFileTreeViewModel = &FileTreeViewModel{}

func NewFileTreeViewModel(getFiles func() []*models.File, log *logrus.Entry, showTree bool) *FileTreeViewModel {
	fileTree := NewFileTree(getFiles, log, showTree)
	listCursor := traits.NewListCursor(fileTree)
	return &FileTreeViewModel{
		IFileTree:   fileTree,
		IListCursor: listCursor,
	}
}

func (self *FileTreeViewModel) GetSelected() *FileNode {
	if self.Len() == 0 {
		return nil
	}

	return self.Get(self.GetSelectedLineIdx())
}

func (self *FileTreeViewModel) GetSelectedFile() *models.File {
	node := self.GetSelected()
	if node == nil {
		return nil
	}

	return node.File
}

func (self *FileTreeViewModel) GetSelectedPath() string {
	node := self.GetSelected()
	if node == nil {
		return ""
	}

	return node.GetPath()
}

func (self *FileTreeViewModel) SetTree() {
	newFiles := self.GetAllFiles()
	selectedNode := self.GetSelected()

	// for when you stage the old file of a rename and the new file is in a collapsed dir
	for _, file := range newFiles {
		if selectedNode != nil && selectedNode.Path != "" && file.PreviousName == selectedNode.Path {
			self.ExpandToPath(file.Name)
		}
	}

	prevNodes := self.GetAllItems()
	prevSelectedLineIdx := self.GetSelectedLineIdx()

	self.IFileTree.SetTree()

	if selectedNode != nil {
		newNodes := self.GetAllItems()
		newIdx := self.findNewSelectedIdx(prevNodes[prevSelectedLineIdx:], newNodes)
		if newIdx != -1 && newIdx != prevSelectedLineIdx {
			self.SetSelectedLineIdx(newIdx)
		}
	}

	self.RefreshSelectedIdx()
}

// Let's try to find our file again and move the cursor to that.
// If we can't find our file, it was probably just removed by the user. In that
// case, we go looking for where the next file has been moved to. Given that the
// user could have removed a whole directory, we continue iterating through the old
// nodes until we find one that exists in the new set of nodes, then move the cursor
// to that.
// prevNodes starts from our previously selected node because we don't need to consider anything above that
func (self *FileTreeViewModel) findNewSelectedIdx(prevNodes []*FileNode, currNodes []*FileNode) int {
	getPaths := func(node *FileNode) []string {
		if node == nil {
			return nil
		}
		if node.File != nil && node.File.IsRename() {
			return node.File.Names()
		} else {
			return []string{node.Path}
		}
	}

	for _, prevNode := range prevNodes {
		selectedPaths := getPaths(prevNode)

		for idx, node := range currNodes {
			paths := getPaths(node)

			// If you started off with a rename selected, and now it's broken in two, we want you to jump to the new file, not the old file.
			// This is because the new should be in the same position as the rename was meaning less cursor jumping
			foundOldFileInRename := prevNode.File != nil && prevNode.File.IsRename() && node.Path == prevNode.File.PreviousName
			foundNode := utils.StringArraysOverlap(paths, selectedPaths) && !foundOldFileInRename
			if foundNode {
				return idx
			}
		}
	}

	return -1
}

func (self *FileTreeViewModel) SetFilter(filter FileTreeDisplayFilter) {
	self.IFileTree.SetFilter(filter)
	self.IListCursor.SetSelectedLineIdx(0)
}

// If we're going from flat to tree we want to select the same file.
// If we're going from tree to flat and we have a file selected we want to select that.
// If instead we've selected a directory we need to select the first file in that directory.
func (self *FileTreeViewModel) ToggleShowTree() {
	selectedNode := self.GetSelected()

	self.IFileTree.ToggleShowTree()

	if selectedNode == nil {
		return
	}
	path := selectedNode.Path

	if self.InTreeMode() {
		self.ExpandToPath(path)
	} else if len(selectedNode.Children) > 0 {
		path = selectedNode.GetLeaves()[0].Path
	}

	index, found := self.GetIndexForPath(path)
	if found {
		self.SetSelectedLineIdx(index)
	}
}
