package filetree

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/sirupsen/logrus"
)

type ICommitFileTree interface {
	ITree

	Get(index int) *CommitFileNode
	GetFile(path string) *models.CommitFile
	GetAllItems() []*CommitFileNode
	GetAllFiles() []*models.CommitFile
}

type CommitFileTree struct {
	getFiles       func() []*models.CommitFile
	tree           *CommitFileNode
	showTree       bool
	log            *logrus.Entry
	collapsedPaths *CollapsedPaths
}

var _ ICommitFileTree = &CommitFileTree{}

func NewCommitFileTree(getFiles func() []*models.CommitFile, log *logrus.Entry, showTree bool) *CommitFileTree {
	return &CommitFileTree{
		getFiles:       getFiles,
		log:            log,
		showTree:       showTree,
		collapsedPaths: NewCollapsedPaths(),
	}
}

func (self *CommitFileTree) ExpandToPath(path string) {
	self.collapsedPaths.ExpandToPath(path)
}

func (self *CommitFileTree) ToggleShowTree() {
	self.showTree = !self.showTree
	self.SetTree()
}

func (self *CommitFileTree) Get(index int) *CommitFileNode {
	// need to traverse the three depth first until we get to the index.
	return self.tree.GetNodeAtIndex(index+1, self.collapsedPaths) // ignoring root
}

func (self *CommitFileTree) GetIndexForPath(path string) (int, bool) {
	index, found := self.tree.GetIndexForPath(path, self.collapsedPaths)
	return index - 1, found
}

func (self *CommitFileTree) GetAllItems() []*CommitFileNode {
	if self.tree == nil {
		return nil
	}

	return self.tree.Flatten(self.collapsedPaths)[1:] // ignoring root
}

func (self *CommitFileTree) Len() int {
	return self.tree.Size(self.collapsedPaths) - 1 // ignoring root
}

func (self *CommitFileTree) GetAllFiles() []*models.CommitFile {
	return self.getFiles()
}

func (self *CommitFileTree) SetTree() {
	if self.showTree {
		self.tree = BuildTreeFromCommitFiles(self.getFiles())
	} else {
		self.tree = BuildFlatTreeFromCommitFiles(self.getFiles())
	}
}

func (self *CommitFileTree) IsCollapsed(path string) bool {
	return self.collapsedPaths.IsCollapsed(path)
}

func (self *CommitFileTree) ToggleCollapsed(path string) {
	self.collapsedPaths.ToggleCollapsed(path)
}

func (self *CommitFileTree) Tree() INode {
	return self.tree
}

func (self *CommitFileTree) CollapsedPaths() *CollapsedPaths {
	return self.collapsedPaths
}

func (self *CommitFileTree) GetFile(path string) *models.CommitFile {
	for _, file := range self.getFiles() {
		if file.Name == path {
			return file
		}
	}

	return nil
}

func (self *CommitFileTree) InTreeMode() bool {
	return self.showTree
}
