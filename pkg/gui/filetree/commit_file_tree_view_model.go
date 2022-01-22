package filetree

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/sirupsen/logrus"
)

type CommitFileTreeViewModel struct {
	files          []*models.CommitFile
	tree           *CommitFileNode
	showTree       bool
	log            *logrus.Entry
	collapsedPaths CollapsedPaths
	// parent is the identifier of the parent object e.g. a commit SHA if this commit file is for a commit, or a stash entry ref like 'stash@{1}'
	parent string
}

func (self *CommitFileTreeViewModel) GetParent() string {
	return self.parent
}

func (self *CommitFileTreeViewModel) SetParent(parent string) {
	self.parent = parent
}

func NewCommitFileTreeViewModel(files []*models.CommitFile, log *logrus.Entry, showTree bool) *CommitFileTreeViewModel {
	viewModel := &CommitFileTreeViewModel{
		log:            log,
		showTree:       showTree,
		collapsedPaths: CollapsedPaths{},
	}

	viewModel.SetFiles(files)

	return viewModel
}

func (self *CommitFileTreeViewModel) ExpandToPath(path string) {
	self.collapsedPaths.ExpandToPath(path)
}

func (self *CommitFileTreeViewModel) ToggleShowTree() {
	self.showTree = !self.showTree
	self.SetTree()
}

func (self *CommitFileTreeViewModel) GetItemAtIndex(index int) *CommitFileNode {
	// need to traverse the three depth first until we get to the index.
	return self.tree.GetNodeAtIndex(index+1, self.collapsedPaths) // ignoring root
}

func (self *CommitFileTreeViewModel) GetIndexForPath(path string) (int, bool) {
	index, found := self.tree.GetIndexForPath(path, self.collapsedPaths)
	return index - 1, found
}

func (self *CommitFileTreeViewModel) GetAllItems() []*CommitFileNode {
	if self.tree == nil {
		return nil
	}

	return self.tree.Flatten(self.collapsedPaths)[1:] // ignoring root
}

func (self *CommitFileTreeViewModel) GetItemsLength() int {
	return self.tree.Size(self.collapsedPaths) - 1 // ignoring root
}

func (self *CommitFileTreeViewModel) GetAllFiles() []*models.CommitFile {
	return self.files
}

func (self *CommitFileTreeViewModel) SetFiles(files []*models.CommitFile) {
	self.files = files

	self.SetTree()
}

func (self *CommitFileTreeViewModel) SetTree() {
	if self.showTree {
		self.tree = BuildTreeFromCommitFiles(self.files)
	} else {
		self.tree = BuildFlatTreeFromCommitFiles(self.files)
	}
}

func (self *CommitFileTreeViewModel) IsCollapsed(path string) bool {
	return self.collapsedPaths.IsCollapsed(path)
}

func (self *CommitFileTreeViewModel) ToggleCollapsed(path string) {
	self.collapsedPaths.ToggleCollapsed(path)
}

func (self *CommitFileTreeViewModel) Tree() INode {
	return self.tree
}

func (self *CommitFileTreeViewModel) CollapsedPaths() CollapsedPaths {
	return self.collapsedPaths
}
