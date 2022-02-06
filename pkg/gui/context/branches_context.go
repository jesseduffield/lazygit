package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context/traits"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type BranchesContext struct {
	*BranchesViewModel
	*ListContextTrait
}

var _ types.IListContext = (*BranchesContext)(nil)

func NewBranchesContext(
	getModel func() []*models.Branch,
	view *gocui.View,
	getDisplayStrings func(startIdx int, length int) [][]string,

	onFocus func(...types.OnFocusOpts) error,
	onRenderToMain func(...types.OnFocusOpts) error,
	onFocusLost func() error,

	c *types.ControllerCommon,
) *BranchesContext {
	viewModel := NewBranchesViewModel(getModel)

	return &BranchesContext{
		BranchesViewModel: viewModel,
		ListContextTrait: &ListContextTrait{
			Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
				ViewName:   "branches",
				WindowName: "branches",
				Key:        LOCAL_BRANCHES_CONTEXT_KEY,
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

func (self *BranchesContext) GetSelectedItemId() string {
	item := self.GetSelected()
	if item == nil {
		return ""
	}

	return item.ID()
}

type BranchesViewModel struct {
	*traits.ListCursor
	getModel func() []*models.Branch
}

func NewBranchesViewModel(getModel func() []*models.Branch) *BranchesViewModel {
	self := &BranchesViewModel{
		getModel: getModel,
	}

	self.ListCursor = traits.NewListCursor(self)

	return self
}

func (self *BranchesViewModel) GetItemsLength() int {
	return len(self.getModel())
}

func (self *BranchesViewModel) GetSelected() *models.Branch {
	if self.GetItemsLength() == 0 {
		return nil
	}

	return self.getModel()[self.GetSelectedLineIdx()]
}

func (self *BranchesViewModel) GetSelectedRefName() string {
	item := self.GetSelected()
	if item == nil {
		return ""
	}

	return item.RefName()
}
