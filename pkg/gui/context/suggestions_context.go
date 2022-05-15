package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type SuggestionsContext struct {
	*FilteredListViewModel[*types.Suggestion]
	*ListContextTrait
}

var _ types.IListContext = (*SuggestionsContext)(nil)

func NewSuggestionsContext(
	getItems func() []*types.Suggestion,
	guiContextState GuiContextState,
	view *gocui.View,

	onFocus func(...types.OnFocusOpts) error,
	onRenderToMain func(...types.OnFocusOpts) error,
	onFocusLost func() error,

	c *types.HelperCommon,
) *SuggestionsContext {
	viewModel := NewFilteredListViewModel(getItems, guiContextState.Needle, func(item *types.Suggestion) string {
		return item.Label
	})

	getDisplayStrings := func(startIdx int, length int) [][]string {
		return presentation.GetSuggestionListDisplayStrings(viewModel.getModel())
	}

	return &SuggestionsContext{
		FilteredListViewModel: viewModel,
		ListContextTrait: &ListContextTrait{
			Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
				ViewName:   "suggestions",
				WindowName: "suggestions",
				Key:        SUGGESTIONS_CONTEXT_KEY,
				Kind:       types.PERSISTENT_POPUP,
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

func (self *SuggestionsContext) GetSelectedItemId() string {
	item := self.GetSelected()
	if item == nil {
		return ""
	}

	return item.Value
}
