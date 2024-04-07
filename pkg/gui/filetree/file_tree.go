package filetree

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
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

type ITree[T any] interface {
	InTreeMode() bool
	ExpandToPath(path string)
	ToggleShowTree()
	GetIndexForPath(path string) (int, bool)
	Len() int
	GetItem(index int) types.HasUrn
	SetTree()
	IsCollapsed(path string) bool
	ToggleCollapsed(path string)
	CollapsedPaths() *CollapsedPaths
}

type IFileTree interface {
	ITree[models.File]

	FilterFiles(test func(*models.File) bool) []*models.File
	SetStatusFilter(filter FileTreeDisplayFilter)
	Get(index int) *FileNode
	GetFile(path string) *models.File
	GetAllItems() []*FileNode
	GetAllFiles() []*models.File
	GetFilter() FileTreeDisplayFilter
	GetRoot() *FileNode
}

type FileTree struct {
	getFiles       func() []*models.File
	tree           *Node[models.File]
	showTree       bool
	log            *logrus.Entry
	filter         FileTreeDisplayFilter
	collapsedPaths *CollapsedPaths
}

var _ IFileTree = &FileTree{}

func NewFileTree(getFiles func() []*models.File, log *logrus.Entry, showTree bool) *FileTree {
	return &FileTree{
		getFiles:       getFiles,
		log:            log,
		showTree:       showTree,
		filter:         DisplayAll,
		collapsedPaths: NewCollapsedPaths(),
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
	return lo.Filter(self.getFiles(), func(file *models.File, _ int) bool { return test(file) })
}

func (self *FileTree) SetStatusFilter(filter FileTreeDisplayFilter) {
	self.filter = filter
	self.SetTree()
}

func (self *FileTree) ToggleShowTree() {
	self.showTree = !self.showTree
	self.SetTree()
}

func (self *FileTree) Get(index int) *FileNode {
	// need to traverse the tree depth first until we get to the index.
	return NewFileNode(self.tree.GetNodeAtIndex(index+1, self.collapsedPaths)) // ignoring root
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

	// ignoring root
	return lo.Map(self.tree.Flatten(self.collapsedPaths)[1:], func(node *Node[models.File], _ int) *FileNode {
		return NewFileNode(node)
	})
}

func (self *FileTree) Len() int {
	// -1 because we're ignoring the root
	return max(self.tree.Size(self.collapsedPaths)-1, 0)
}

func (self *FileTree) GetItem(index int) types.HasUrn {
	// Unimplemented because we don't yet need to show inlines statuses in commit file views
	return nil
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

func (self *FileTree) Tree() *FileNode {
	return NewFileNode(self.tree)
}

func (self *FileTree) GetRoot() *FileNode {
	return NewFileNode(self.tree)
}

func (self *FileTree) CollapsedPaths() *CollapsedPaths {
	return self.collapsedPaths
}

func (self *FileTree) GetFilter() FileTreeDisplayFilter {
	return self.filter
}
