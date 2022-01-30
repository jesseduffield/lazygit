package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/sirupsen/logrus"
)

type WorkingTreeContext struct {
	*WorkingTreeViewModal
	*BaseContext
	*ListContextTrait
}

var _ types.IListContext = (*WorkingTreeContext)(nil)

func NewWorkingTreeContext(
	getModel func() []*models.File,
	getView func() *gocui.View,
	getDisplayStrings func(startIdx int, length int) [][]string,

	onFocus func(...types.OnFocusOpts) error,
	onRenderToMain func(...types.OnFocusOpts) error,
	onFocusLost func() error,

	c *types.ControllerCommon,
) *WorkingTreeContext {
	baseContext := NewBaseContext(NewBaseContextOpts{
		ViewName:   "files",
		WindowName: "files",
		Key:        FILES_CONTEXT_KEY,
		Kind:       types.SIDE_CONTEXT,
	})

	self := &WorkingTreeContext{}
	takeFocus := func() error { return c.PushContext(self) }

	list := NewWorkingTreeViewModal(getModel, c.Log, c.UserConfig.Gui.ShowFileTree)
	viewTrait := NewViewTrait(getView)
	listContextTrait := &ListContextTrait{
		base:      baseContext,
		listTrait: list.ListTrait,
		viewTrait: viewTrait,

		GetDisplayStrings: getDisplayStrings,
		OnFocus:           onFocus,
		OnRenderToMain:    onRenderToMain,
		OnFocusLost:       onFocusLost,
		takeFocus:         takeFocus,

		// TODO: handle this in a trait
		RenderSelection: false,

		c: c,
	}

	self.BaseContext = baseContext
	self.ListContextTrait = listContextTrait
	self.WorkingTreeViewModal = list

	return self
}

type WorkingTreeViewModal struct {
	*ListTrait
	*filetree.FileTreeViewModel
	getModel func() []*models.File
}

func (self *WorkingTreeViewModal) GetItemsLength() int {
	return self.FileTreeViewModel.GetItemsLength()
}

func (self *WorkingTreeViewModal) GetSelectedFileNode() *filetree.FileNode {
	if self.GetItemsLength() == 0 {
		return nil
	}

	return self.FileTreeViewModel.GetItemAtIndex(self.selectedIdx)
}

func (self *WorkingTreeViewModal) GetSelectedItem() (types.ListItem, bool) {
	item := self.GetSelectedFileNode()
	return item, item != nil
}

func NewWorkingTreeViewModal(getModel func() []*models.File, log *logrus.Entry, showTree bool) *WorkingTreeViewModal {
	self := &WorkingTreeViewModal{
		getModel:          getModel,
		FileTreeViewModel: filetree.NewFileTreeViewModel(getModel, log, showTree),
	}

	self.ListTrait = &ListTrait{
		selectedIdx: 0,
		HasLength:   self,
	}

	return self
}
