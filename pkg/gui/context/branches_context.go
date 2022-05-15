package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type BranchesContext struct {
	*FilteredListViewModel[*models.Branch]
	*ListContextTrait
}

var _ types.IListContext = (*BranchesContext)(nil)

func NewBranchesContext(
	getItems func() []*models.Branch,
	guiContextState GuiContextState,
	view *gocui.View,

	onFocus func(...types.OnFocusOpts) error,
	onRenderToMain func(...types.OnFocusOpts) error,
	onFocusLost func() error,

	c *types.HelperCommon,
) *BranchesContext {
	viewModel := NewFilteredListViewModel(getItems, guiContextState.Needle, func(item *models.Branch) string {
		return item.Name
	})

	getDisplayStrings := func(startIdx int, length int) [][]string {
		return presentation.GetBranchListDisplayStrings(viewModel.getModel(), guiContextState.ScreenMode() != types.SCREEN_NORMAL, guiContextState.Modes().Diffing.Ref, c.Tr)
	}

	return &BranchesContext{
		FilteredListViewModel: viewModel,
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

func (self *BranchesContext) GetSelectedRef() types.Ref {
	branch := self.GetSelected()
	if branch == nil {
		return nil
	}
	return branch
}
