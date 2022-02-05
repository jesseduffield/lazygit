package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type CommitFilesContext struct {
	*filetree.CommitFileTreeViewModel
	*ListContextTrait
}

var _ types.IListContext = (*CommitFilesContext)(nil)

func NewCommitFilesContext(
	getModel func() []*models.CommitFile,
	view *gocui.View,
	getDisplayStrings func(startIdx int, length int) [][]string,

	onFocus func(...types.OnFocusOpts) error,
	onRenderToMain func(...types.OnFocusOpts) error,
	onFocusLost func() error,

	c *types.ControllerCommon,
) *CommitFilesContext {
	viewModel := filetree.NewCommitFileTreeViewModel(getModel, c.Log, c.UserConfig.Gui.ShowFileTree)

	return &CommitFilesContext{
		CommitFileTreeViewModel: viewModel,
		ListContextTrait: &ListContextTrait{
			Context: NewSimpleContext(
				NewBaseContext(NewBaseContextOpts{
					ViewName:   "commitFiles",
					WindowName: "commits",
					Key:        COMMIT_FILES_CONTEXT_KEY,
					Kind:       types.SIDE_CONTEXT,
					Focusable:  true,
				}),
				ContextCallbackOpts{
					OnFocus:        onFocus,
					OnFocusLost:    onFocusLost,
					OnRenderToMain: onRenderToMain,
				}),
			list:              viewModel,
			viewTrait:         NewViewTrait(view),
			getDisplayStrings: getDisplayStrings,
			c:                 c,
		},
	}
}

func (self *CommitFilesContext) GetSelectedItemId() string {
	item := self.GetSelectedFileNode()
	if item == nil {
		return ""
	}

	return item.ID()
}
