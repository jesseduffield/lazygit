package filetree

import (
	"sync"

	"github.com/lobes/lazytask/pkg/commands/models"
	"github.com/lobes/lazytask/pkg/gui/context/traits"
	"github.com/lobes/lazytask/pkg/gui/types"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

type ICommitFileTreeViewModel interface {
	ICommitFileTree
	types.IListCursor

	GetRef() types.Ref
	SetRef(types.Ref)
	GetCanRebase() bool
	SetCanRebase(bool)
}

type CommitFileTreeViewModel struct {
	sync.RWMutex
	types.IListCursor
	ICommitFileTree

	// this is e.g. the commit for which we're viewing the files
	ref types.Ref

	// we set this to true when you're viewing the files within the checked-out branch's commits.
	// If you're viewing the files of some random other branch we can't do any rebase stuff.
	canRebase bool
}

var _ ICommitFileTreeViewModel = &CommitFileTreeViewModel{}

func NewCommitFileTreeViewModel(getFiles func() []*models.CommitFile, log *logrus.Entry, showTree bool) *CommitFileTreeViewModel {
	fileTree := NewCommitFileTree(getFiles, log, showTree)
	listCursor := traits.NewListCursor(fileTree.Len)
	return &CommitFileTreeViewModel{
		ICommitFileTree: fileTree,
		IListCursor:     listCursor,
		ref:             nil,
		canRebase:       false,
	}
}

func (self *CommitFileTreeViewModel) GetRef() types.Ref {
	return self.ref
}

func (self *CommitFileTreeViewModel) SetRef(ref types.Ref) {
	self.ref = ref
}

func (self *CommitFileTreeViewModel) GetCanRebase() bool {
	return self.canRebase
}

func (self *CommitFileTreeViewModel) SetCanRebase(canRebase bool) {
	self.canRebase = canRebase
}

func (self *CommitFileTreeViewModel) GetSelected() *CommitFileNode {
	if self.Len() == 0 {
		return nil
	}

	return self.Get(self.GetSelectedLineIdx())
}

func (self *CommitFileTreeViewModel) GetSelectedItemId() string {
	item := self.GetSelected()
	if item == nil {
		return ""
	}

	return item.ID()
}

func (self *CommitFileTreeViewModel) GetSelectedItems() ([]*CommitFileNode, int, int) {
	if self.Len() == 0 {
		return nil, 0, 0
	}

	startIdx, endIdx := self.GetSelectionRange()

	nodes := []*CommitFileNode{}
	for i := startIdx; i <= endIdx; i++ {
		nodes = append(nodes, self.Get(i))
	}

	return nodes, startIdx, endIdx
}

func (self *CommitFileTreeViewModel) GetSelectedItemIds() ([]string, int, int) {
	selectedItems, startIdx, endIdx := self.GetSelectedItems()

	ids := lo.Map(selectedItems, func(item *CommitFileNode, _ int) string {
		return item.ID()
	})

	return ids, startIdx, endIdx
}

func (self *CommitFileTreeViewModel) GetSelectedFile() *models.CommitFile {
	node := self.GetSelected()
	if node == nil {
		return nil
	}

	return node.File
}

func (self *CommitFileTreeViewModel) GetSelectedPath() string {
	node := self.GetSelected()
	if node == nil {
		return ""
	}

	return node.GetPath()
}

// duplicated from file_tree_view_model.go. Generics will help here
func (self *CommitFileTreeViewModel) ToggleShowTree() {
	selectedNode := self.GetSelected()

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
		self.SetSelection(index)
	}
}
