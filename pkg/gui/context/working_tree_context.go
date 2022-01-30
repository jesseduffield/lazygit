package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type WorkingTreeContext struct {
	*filetree.FileTreeViewModel
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

	viewModel := filetree.NewFileTreeViewModel(getModel, c.Log, c.UserConfig.Gui.ShowFileTree)
	viewTrait := NewViewTrait(getView)
	listContextTrait := &ListContextTrait{
		base:      baseContext,
		list:      viewModel,
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
	self.FileTreeViewModel = viewModel

	return self
}

func (self *WorkingTreeContext) GetSelectedItem() (types.ListItem, bool) {
	item := self.FileTreeViewModel.GetSelectedFileNode()
	return item, item != nil
}
