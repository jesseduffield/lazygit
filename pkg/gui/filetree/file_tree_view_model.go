package filetree

import (
	"sync"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/sirupsen/logrus"
)

type FileTreeDisplayFilter int

const (
	DisplayAll FileTreeDisplayFilter = iota
	DisplayStaged
	DisplayUnstaged
)

type FileTreeViewModel struct {
	files          []*models.File
	tree           *FileNode
	showTree       bool
	log            *logrus.Entry
	filter         FileTreeDisplayFilter
	collapsedPaths CollapsedPaths
	sync.RWMutex
}

func NewFileTreeViewModel(files []*models.File, log *logrus.Entry, showTree bool) *FileTreeViewModel {
	viewModel := &FileTreeViewModel{
		log:            log,
		showTree:       showTree,
		filter:         DisplayAll,
		collapsedPaths: CollapsedPaths{},
		RWMutex:        sync.RWMutex{},
	}

	viewModel.SetFiles(files)

	return viewModel
}

func (self *FileTreeViewModel) InTreeMode() bool {
	return self.showTree
}

func (self *FileTreeViewModel) ExpandToPath(path string) {
	self.collapsedPaths.ExpandToPath(path)
}

func (self *FileTreeViewModel) GetFilesForDisplay() []*models.File {
	files := self.files
	if self.filter == DisplayAll {
		return files
	}

	result := make([]*models.File, 0)
	if self.filter == DisplayStaged {
		for _, file := range files {
			if file.HasStagedChanges {
				result = append(result, file)
			}
		}
	} else {
		for _, file := range files {
			if !file.HasStagedChanges {
				result = append(result, file)
			}
		}
	}

	return result
}

func (self *FileTreeViewModel) SetDisplayFilter(filter FileTreeDisplayFilter) {
	self.filter = filter
	self.SetTree()
}

func (self *FileTreeViewModel) ToggleShowTree() {
	self.showTree = !self.showTree
	self.SetTree()
}

func (self *FileTreeViewModel) GetItemAtIndex(index int) *FileNode {
	// need to traverse the three depth first until we get to the index.
	return self.tree.GetNodeAtIndex(index+1, self.collapsedPaths) // ignoring root
}

func (self *FileTreeViewModel) GetIndexForPath(path string) (int, bool) {
	index, found := self.tree.GetIndexForPath(path, self.collapsedPaths)
	return index - 1, found
}

func (self *FileTreeViewModel) GetAllItems() []*FileNode {
	if self.tree == nil {
		return nil
	}

	return self.tree.Flatten(self.collapsedPaths)[1:] // ignoring root
}

func (self *FileTreeViewModel) GetItemsLength() int {
	return self.tree.Size(self.collapsedPaths) - 1 // ignoring root
}

func (self *FileTreeViewModel) GetAllFiles() []*models.File {
	return self.files
}

func (self *FileTreeViewModel) SetFiles(files []*models.File) {
	self.files = files

	self.SetTree()
}

func (self *FileTreeViewModel) SetTree() {
	filesForDisplay := self.GetFilesForDisplay()
	if self.showTree {
		self.tree = BuildTreeFromFiles(filesForDisplay)
	} else {
		self.tree = BuildFlatTreeFromFiles(filesForDisplay)
	}
}

func (self *FileTreeViewModel) IsCollapsed(path string) bool {
	return self.collapsedPaths.IsCollapsed(path)
}

func (self *FileTreeViewModel) ToggleCollapsed(path string) {
	self.collapsedPaths.ToggleCollapsed(path)
}

func (self *FileTreeViewModel) Tree() INode {
	return self.tree
}

func (self *FileTreeViewModel) CollapsedPaths() CollapsedPaths {
	return self.collapsedPaths
}
