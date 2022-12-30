package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type WorkingTreeContext struct {
	*filetree.FileTreeViewModel
	*ListContextTrait
}

var _ types.IListContext = (*WorkingTreeContext)(nil)

func NewWorkingTreeContext(
	getModel func() []*models.File,
	view *gocui.View,
	getDisplayStrings func(startIdx int, length int) [][]string,

	c *types.HelperCommon,
) *WorkingTreeContext {
	viewModel := filetree.NewFileTreeViewModel(getModel, c.Log, c.UserConfig.Gui.ShowFileTree)

	return &WorkingTreeContext{
		FileTreeViewModel: viewModel,
		ListContextTrait: &ListContextTrait{
			Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
				View:       view,
				WindowName: "files",
				Key:        FILES_CONTEXT_KEY,
				Kind:       types.SIDE_CONTEXT,
				Focusable:  true,
			}), ContextCallbackOpts{}),
			list:              viewModel,
			getDisplayStrings: getDisplayStrings,
			c:                 c,
		},
	}
}

func (self *WorkingTreeContext) GetSelectedItemId() string {
	item := self.GetSelected()
	if item == nil {
		return ""
	}

	return item.ID()
}
