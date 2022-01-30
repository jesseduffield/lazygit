package filetree

import (
	"fmt"
	"sync"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/sirupsen/logrus"
)

type FileTreeDisplayFilter int

const (
	DisplayAll FileTreeDisplayFilter = iota
	DisplayStaged
	DisplayUnstaged
	// this shows files with merge conflicts
	DisplayConflicted
)

type FileTreeViewModel struct {
	getFiles       func() []*models.File
	tree           *FileNode
	showTree       bool
	log            *logrus.Entry
	filter         FileTreeDisplayFilter
	collapsedPaths CollapsedPaths
	sync.RWMutex
}

func NewFileTreeViewModel(getFiles func() []*models.File, log *logrus.Entry, showTree bool) *FileTreeViewModel {
	viewModel := &FileTreeViewModel{
		getFiles:       getFiles,
		log:            log,
		showTree:       showTree,
		filter:         DisplayAll,
		collapsedPaths: CollapsedPaths{},
		RWMutex:        sync.RWMutex{},
	}

	return viewModel
}

func (self *FileTreeViewModel) InTreeMode() bool {
	return self.showTree
}

func (self *FileTreeViewModel) ExpandToPath(path string) {
	self.collapsedPaths.ExpandToPath(path)
}

func (self *FileTreeViewModel) GetFilesForDisplay() []*models.File {
	switch self.filter {
	case DisplayAll:
		return self.getFiles()
	case DisplayStaged:
		return self.FilterFiles(func(file *models.File) bool { return file.HasStagedChanges })
	case DisplayUnstaged:
		return self.FilterFiles(func(file *models.File) bool { return file.HasUnstagedChanges })
	case DisplayConflicted:
		return self.FilterFiles(func(file *models.File) bool { return file.HasMergeConflicts })
	default:
		panic(fmt.Sprintf("Unexpected files display filter: %d", self.filter))
	}
}

func (self *FileTreeViewModel) FilterFiles(test func(*models.File) bool) []*models.File {
	result := make([]*models.File, 0)
	for _, file := range self.getFiles() {
		if test(file) {
			result = append(result, file)
		}
	}
	return result
}

func (self *FileTreeViewModel) SetFilter(filter FileTreeDisplayFilter) {
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

func (self *FileTreeViewModel) GetFile(path string) *models.File {
	for _, file := range self.getFiles() {
		if file.Name == path {
			return file
		}
	}

	return nil
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
	return self.getFiles()
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

func (self *FileTreeViewModel) GetFilter() FileTreeDisplayFilter {
	return self.filter
}
