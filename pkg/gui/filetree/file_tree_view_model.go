package filetree

import (
	"strings"
	"sync"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/gui/context/traits"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
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
	types.IListCursor
	IFileTree
}

var _ IFileTreeViewModel = &FileTreeViewModel{}

func NewFileTreeViewModel(getFiles func() []*models.File, common *common.Common, showTree bool) *FileTreeViewModel {
	fileTree := NewFileTree(getFiles, common, showTree)
	listCursor := traits.NewListCursor(fileTree.Len)
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

func (self *FileTreeViewModel) GetSelectedItemId() string {
	item := self.GetSelected()
	if item == nil {
		return ""
	}

	return item.ID()
}

func (self *FileTreeViewModel) GetSelectedItems() ([]*FileNode, int, int) {
	if self.Len() == 0 {
		return nil, 0, 0
	}

	startIdx, endIdx := self.GetSelectionRange()

	nodes := []*FileNode{}
	for i := startIdx; i <= endIdx; i++ {
		nodes = append(nodes, self.Get(i))
	}

	return nodes, startIdx, endIdx
}

func (self *FileTreeViewModel) GetSelectedItemIds() ([]string, int, int) {
	selectedItems, startIdx, endIdx := self.GetSelectedItems()

	ids := lo.Map(selectedItems, func(item *FileNode, _ int) string {
		return item.ID()
	})

	return ids, startIdx, endIdx
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
		if selectedNode != nil && selectedNode.path != "" && file.PreviousPath == selectedNode.path {
			self.ExpandToPath(file.Path)
		}
	}

	prevNodes := self.GetAllItems()
	prevSelectedLineIdx := self.GetSelectedLineIdx()

	self.IFileTree.SetTree()

	if selectedNode != nil {
		newNodes := self.GetAllItems()
		newIdx := self.findNewSelectedIdx(prevNodes[prevSelectedLineIdx:], newNodes)
		if newIdx != -1 && newIdx != prevSelectedLineIdx {
			self.SetSelection(newIdx)
		}
	}

	self.ClampSelection()
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
		}
		return []string{node.path}
	}

	for _, prevNode := range prevNodes {
		selectedPaths := getPaths(prevNode)

		for idx, node := range currNodes {
			paths := getPaths(node)

			// If you started off with a rename selected, and now it's broken in two, we want you to jump to the new file, not the old file.
			// This is because the new should be in the same position as the rename was meaning less cursor jumping
			foundOldFileInRename := prevNode.File != nil && prevNode.File.IsRename() && node.path == prevNode.File.PreviousPath
			foundNode := utils.StringArraysOverlap(paths, selectedPaths) && !foundOldFileInRename
			if foundNode {
				return idx
			}
		}
	}

	return -1
}

func (self *FileTreeViewModel) SetStatusFilter(filter FileTreeDisplayFilter) {
	self.IFileTree.SetStatusFilter(filter)
	self.IListCursor.SetSelection(0)
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
	path := selectedNode.path

	if self.InTreeMode() {
		self.ExpandToPath(path)
	} else if len(selectedNode.Children) > 0 {
		path = selectedNode.GetLeaves()[0].path
	}

	index, found := self.GetIndexForPath(path)
	if found {
		self.SetSelectedLineIdx(index)
	}
}

func (self *FileTreeViewModel) CollapseAll() {
	selectedNode := self.GetSelected()

	self.IFileTree.CollapseAll()
	if selectedNode == nil {
		return
	}

	topLevelPath := strings.Split(selectedNode.path, "/")[0]
	index, found := self.GetIndexForPath(topLevelPath)
	if found {
		self.SetSelectedLineIdx(index)
	}
}

func (self *FileTreeViewModel) ExpandAll() {
	selectedNode := self.GetSelected()

	self.IFileTree.ExpandAll()

	if selectedNode == nil {
		return
	}

	index, found := self.GetIndexForPath(selectedNode.path)
	if found {
		self.SetSelectedLineIdx(index)
	}
}
