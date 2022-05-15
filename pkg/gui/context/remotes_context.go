package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type RemotesContext struct {
	*FilteredListViewModel[*models.Remote]
	*ListContextTrait
}

var _ types.IListContext = (*RemotesContext)(nil)

func NewRemotesContext(
	getItems func() []*models.Remote,
	guiContextState GuiContextState,
	view *gocui.View,

	onFocus func(...types.OnFocusOpts) error,
	onRenderToMain func(...types.OnFocusOpts) error,
	onFocusLost func() error,

	c *types.HelperCommon,
) *RemotesContext {
	viewModel := NewFilteredListViewModel(getItems, guiContextState.Needle, func(item *models.Remote) string {
		return item.Name
	})

	getDisplayStrings := func(startIdx int, length int) [][]string {
		return presentation.GetRemoteListDisplayStrings(viewModel.getModel(), guiContextState.Modes().Diffing.Ref)
	}

	return &RemotesContext{
		FilteredListViewModel: viewModel,
		ListContextTrait: &ListContextTrait{
			Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
				ViewName:   "branches",
				WindowName: "branches",
				Key:        REMOTES_CONTEXT_KEY,
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

func (self *RemotesContext) GetSelectedItemId() string {
	item := self.GetSelected()
	if item == nil {
		return ""
	}

	return item.ID()
}
