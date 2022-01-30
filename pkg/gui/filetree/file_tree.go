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

type IFileTree interface {
	InTreeMode() bool
	ExpandToPath(path string)
	FilterFiles(test func(*models.File) bool) []*models.File
	SetFilter(filter FileTreeDisplayFilter)
	ToggleShowTree()

	GetItemAtIndex(index int) *FileNode
	GetFile(path string) *models.File
	GetIndexForPath(path string) (int, bool)
	GetAllItems() []*FileNode
	GetItemsLength() int
	GetAllFiles() []*models.File

	SetTree()
	IsCollapsed(path string) bool
	ToggleCollapsed(path string)
	Tree() INode
	CollapsedPaths() CollapsedPaths
	GetFilter() FileTreeDisplayFilter
}

type FileTree struct {
	getFiles       func() []*models.File
	tree           *FileNode
	showTree       bool
	log            *logrus.Entry
	filter         FileTreeDisplayFilter
	collapsedPaths CollapsedPaths

	sync.RWMutex
}

func NewFileTree(getFiles func() []*models.File, log *logrus.Entry, showTree bool) *FileTree {
	return &FileTree{
		getFiles:       getFiles,
		log:            log,
		showTree:       showTree,
		filter:         DisplayAll,
		collapsedPaths: CollapsedPaths{},
		RWMutex:        sync.RWMutex{},
	}
}

func (self *FileTree) InTreeMode() bool {
	return self.showTree
}

func (self *FileTree) ExpandToPath(path string) {
	self.collapsedPaths.ExpandToPath(path)
}

func (self *FileTree) getFilesForDisplay() []*models.File {
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

func (self *FileTree) FilterFiles(test func(*models.File) bool) []*models.File {
	result := make([]*models.File, 0)
	for _, file := range self.getFiles() {
		if test(file) {
			result = append(result, file)
		}
	}
	return result
}

func (self *FileTree) SetFilter(filter FileTreeDisplayFilter) {
	self.filter = filter
	self.SetTree()
}

func (self *FileTree) ToggleShowTree() {
	self.showTree = !self.showTree
	self.SetTree()
}

func (self *FileTree) GetItemAtIndex(index int) *FileNode {
	// need to traverse the three depth first until we get to the index.
	return self.tree.GetNodeAtIndex(index+1, self.collapsedPaths) // ignoring root
}

func (self *FileTree) GetFile(path string) *models.File {
	for _, file := range self.getFiles() {
		if file.Name == path {
			return file
		}
	}

	return nil
}

func (self *FileTree) GetIndexForPath(path string) (int, bool) {
	index, found := self.tree.GetIndexForPath(path, self.collapsedPaths)
	return index - 1, found
}

// note: this gets all items when the filter is taken into consideration. There may
// be hidden files that aren't included here. Files off the screen however will
// be included
func (self *FileTree) GetAllItems() []*FileNode {
	if self.tree == nil {
		return nil
	}

	return self.tree.Flatten(self.collapsedPaths)[1:] // ignoring root
}

func (self *FileTree) GetItemsLength() int {
	return self.tree.Size(self.collapsedPaths) - 1 // ignoring root
}

func (self *FileTree) GetAllFiles() []*models.File {
	return self.getFiles()
}

func (self *FileTree) SetTree() {
	filesForDisplay := self.getFilesForDisplay()
	if self.showTree {
		self.tree = BuildTreeFromFiles(filesForDisplay)
	} else {
		self.tree = BuildFlatTreeFromFiles(filesForDisplay)
	}
}

func (self *FileTree) IsCollapsed(path string) bool {
	return self.collapsedPaths.IsCollapsed(path)
}

func (self *FileTree) ToggleCollapsed(path string) {
	self.collapsedPaths.ToggleCollapsed(path)
}

func (self *FileTree) Tree() INode {
	return self.tree
}

func (self *FileTree) CollapsedPaths() CollapsedPaths {
	return self.collapsedPaths
}

func (self *FileTree) GetFilter() FileTreeDisplayFilter {
	return self.filter
}
