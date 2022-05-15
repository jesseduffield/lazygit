package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type RemoteBranchesContext struct {
	*FilteredListViewModel[*models.RemoteBranch]
	*ListContextTrait
	*DynamicTitleBuilder
}

var _ types.IListContext = (*RemoteBranchesContext)(nil)

func NewRemoteBranchesContext(
	getItems func() []*models.RemoteBranch,
	guiContextState GuiContextState,
	view *gocui.View,

	onFocus func(...types.OnFocusOpts) error,
	onRenderToMain func(...types.OnFocusOpts) error,
	onFocusLost func() error,

	c *types.HelperCommon,
) *RemoteBranchesContext {
	viewModel := NewFilteredListViewModel(getItems, guiContextState.Needle, func(item *models.RemoteBranch) string {
		return item.FullName()
	})

	getDisplayStrings := func(startIdx int, length int) [][]string {
		return presentation.GetRemoteBranchListDisplayStrings(viewModel.getModel(), guiContextState.Modes().Diffing.Ref)
	}

	return &RemoteBranchesContext{
		FilteredListViewModel: viewModel,
		DynamicTitleBuilder:   NewDynamicTitleBuilder(c.Tr.RemoteBranchesDynamicTitle),
		ListContextTrait: &ListContextTrait{
			Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
				ViewName:   "remoteBranches",
				WindowName: "branches",
				Key:        REMOTE_BRANCHES_CONTEXT_KEY,
				Kind:       types.SIDE_CONTEXT,
				Focusable:  true,
				Transient:  true,
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

func (self *RemoteBranchesContext) GetSelectedRef() types.Ref {
	remoteBranch := self.GetSelected()
	if remoteBranch == nil {
		return nil
	}
	return remoteBranch
}
