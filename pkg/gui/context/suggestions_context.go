package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type SuggestionsContext struct {
	*BasicViewModel[*types.Suggestion]
	*ListContextTrait
}

var _ types.IListContext = (*SuggestionsContext)(nil)

func NewSuggestionsContext(
	getModel func() []*types.Suggestion,
	view *gocui.View,
	getDisplayStrings func(startIdx int, length int) [][]string,

	c *types.HelperCommon,
) *SuggestionsContext {
	viewModel := NewBasicViewModel(getModel)

	return &SuggestionsContext{
		BasicViewModel: viewModel,
		ListContextTrait: &ListContextTrait{
			Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
				View:                  view,
				WindowName:            "suggestions",
				Key:                   SUGGESTIONS_CONTEXT_KEY,
				Kind:                  types.PERSISTENT_POPUP,
				Focusable:             true,
				HasUncontrolledBounds: true,
			}), ContextCallbackOpts{}),
			list:              viewModel,
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
