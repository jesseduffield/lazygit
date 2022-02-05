package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context/traits"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type RemoteBranchesContext struct {
	*RemoteBranchesViewModel
	*ListContextTrait
}

var _ types.IListContext = (*RemoteBranchesContext)(nil)

func NewRemoteBranchesContext(
	getModel func() []*models.RemoteBranch,
	view *gocui.View,
	getDisplayStrings func(startIdx int, length int) [][]string,

	onFocus func(...types.OnFocusOpts) error,
	onRenderToMain func(...types.OnFocusOpts) error,
	onFocusLost func() error,

	c *types.ControllerCommon,
) *RemoteBranchesContext {
	viewModel := NewRemoteBranchesViewModel(getModel)

	return &RemoteBranchesContext{
		RemoteBranchesViewModel: viewModel,
		ListContextTrait: &ListContextTrait{
			Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
				ViewName:   "branches",
				WindowName: "branches",
				Key:        REMOTE_BRANCHES_CONTEXT_KEY,
				Kind:       types.SIDE_CONTEXT,
				Focusable:  true,
			}), ContextCallbackOpts{
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

func (self *RemoteBranchesContext) GetSelectedItemId() string {
	item := self.GetSelected()
	if item == nil {
		return ""
	}

	return item.ID()
}

type RemoteBranchesViewModel struct {
	*traits.ListCursor
	getModel func() []*models.RemoteBranch
}

func NewRemoteBranchesViewModel(getModel func() []*models.RemoteBranch) *RemoteBranchesViewModel {
	self := &RemoteBranchesViewModel{
		getModel: getModel,
	}

	self.ListCursor = traits.NewListCursor(self)

	return self
}

func (self *RemoteBranchesViewModel) GetItemsLength() int {
	return len(self.getModel())
}

func (self *RemoteBranchesViewModel) GetSelected() *models.RemoteBranch {
	if self.GetItemsLength() == 0 {
		return nil
	}

	return self.getModel()[self.GetSelectedLineIdx()]
}
