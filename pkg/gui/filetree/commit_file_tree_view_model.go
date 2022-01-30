package filetree

import (
	"sync"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context/traits"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/sirupsen/logrus"
)

type ICommitFileTreeViewModel interface {
	ICommitFileTree
	types.IListCursor

	GetRefName() string
	SetRefName(string)
	GetCanRebase() bool
	SetCanRebase(bool)
}

type CommitFileTreeViewModel struct {
	sync.RWMutex
	ICommitFileTree
	types.IListCursor

	// this is e.g. the commit SHA of the commit for which we're viewing the files
	refName string

	// we set this to true when you're viewing the files within the checked-out branch's commits.
	// If you're viewing the files of some random other branch we can't do any rebase stuff.
	canRebase bool
}

var _ ICommitFileTreeViewModel = &CommitFileTreeViewModel{}

func NewCommitFileTreeViewModel(getFiles func() []*models.CommitFile, log *logrus.Entry, showTree bool) *CommitFileTreeViewModel {
	fileTree := NewCommitFileTree(getFiles, log, showTree)
	listCursor := traits.NewListCursor(fileTree)
	return &CommitFileTreeViewModel{
		ICommitFileTree: fileTree,
		IListCursor:     listCursor,
		refName:         "",
		canRebase:       false,
	}
}

func (self *CommitFileTreeViewModel) GetRefName() string {
	return self.refName
}

func (self *CommitFileTreeViewModel) SetRefName(refName string) {
	self.refName = refName
}

func (self *CommitFileTreeViewModel) GetCanRebase() bool {
	return self.canRebase
}

func (self *CommitFileTreeViewModel) SetCanRebase(canRebase bool) {
	self.canRebase = canRebase
}

func (self *CommitFileTreeViewModel) GetSelectedFileNode() *CommitFileNode {
	if self.GetItemsLength() == 0 {
		return nil
	}

	return self.GetItemAtIndex(self.GetSelectedLineIdx())
}

// duplicated from file_tree_view_model.go. Generics will help here
func (self *CommitFileTreeViewModel) ToggleShowTree() {
	selectedNode := self.GetSelectedFileNode()

	self.ICommitFileTree.ToggleShowTree()

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
