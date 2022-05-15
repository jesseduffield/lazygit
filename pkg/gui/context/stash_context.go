package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type StashContext struct {
	*FilteredListViewModel[*models.StashEntry]
	*ListContextTrait
}

var _ types.IListContext = (*StashContext)(nil)

func NewStashContext(
	getItems func() []*models.StashEntry,
	guiContextState GuiContextState,
	view *gocui.View,

	onFocus func(...types.OnFocusOpts) error,
	onRenderToMain func(...types.OnFocusOpts) error,
	onFocusLost func() error,

	c *types.HelperCommon,
) *StashContext {
	viewModel := NewFilteredListViewModel(getItems, guiContextState.Needle, func(item *models.StashEntry) string {
		return item.Name
	})

	getDisplayStrings := func(startIdx int, length int) [][]string {
		return presentation.GetStashEntryListDisplayStrings(viewModel.getModel(), guiContextState.Modes().Diffing.Ref)
	}

	return &StashContext{
		FilteredListViewModel: viewModel,
		ListContextTrait: &ListContextTrait{
			Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
				ViewName:   "stash",
				WindowName: "stash",
				Key:        STASH_CONTEXT_KEY,
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

func (self *StashContext) GetSelectedItemId() string {
	item := self.GetSelected()
	if item == nil {
		return ""
	}

	return item.ID()
}

func (self *StashContext) CanRebase() bool {
	return false
}

func (self *StashContext) GetSelectedRef() types.Ref {
	stash := self.GetSelected()
	if stash == nil {
		return nil
	}
	return stash
}
