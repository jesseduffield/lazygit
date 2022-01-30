package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type CommitFilesContext struct {
	*filetree.CommitFileTreeViewModel
	*BaseContext
	*ListContextTrait
}

var _ types.IListContext = (*CommitFilesContext)(nil)

func NewCommitFilesContext(
	getModel func() []*models.CommitFile,
	getView func() *gocui.View,
	getDisplayStrings func(startIdx int, length int) [][]string,

	onFocus func(...types.OnFocusOpts) error,
	onRenderToMain func(...types.OnFocusOpts) error,
	onFocusLost func() error,

	c *types.ControllerCommon,
) *CommitFilesContext {
	baseContext := NewBaseContext(NewBaseContextOpts{
		ViewName:   "commitFiles",
		WindowName: "commits",
		Key:        COMMIT_FILES_CONTEXT_KEY,
		Kind:       types.SIDE_CONTEXT,
	})

	self := &CommitFilesContext{}
	takeFocus := func() error { return c.PushContext(self) }

	viewModel := filetree.NewCommitFileTreeViewModel(getModel, c.Log, c.UserConfig.Gui.ShowFileTree)
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
	self.CommitFileTreeViewModel = viewModel

	return self
}

func (self *CommitFilesContext) GetSelectedItem() (types.ListItem, bool) {
	item := self.CommitFileTreeViewModel.GetSelectedFileNode()
	return item, item != nil
}
