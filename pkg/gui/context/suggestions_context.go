package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context/traits"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type SuggestionsContext struct {
	*SuggestionsViewModel
	*ListContextTrait
}

var _ types.IListContext = (*SuggestionsContext)(nil)

func NewSuggestionsContext(
	getModel func() []*types.Suggestion,
	view *gocui.View,
	getDisplayStrings func(startIdx int, length int) [][]string,

	onFocus func(...types.OnFocusOpts) error,
	onRenderToMain func(...types.OnFocusOpts) error,
	onFocusLost func() error,

	c *types.ControllerCommon,
) *SuggestionsContext {
	viewModel := NewSuggestionsViewModel(getModel)

	return &SuggestionsContext{
		SuggestionsViewModel: viewModel,
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

type SuggestionsViewModel struct {
	*traits.ListCursor
	getModel func() []*types.Suggestion
}

func NewSuggestionsViewModel(getModel func() []*types.Suggestion) *SuggestionsViewModel {
	self := &SuggestionsViewModel{
		getModel: getModel,
	}

	self.ListCursor = traits.NewListCursor(self)

	return self
}

func (self *SuggestionsViewModel) GetItemsLength() int {
	return len(self.getModel())
}

func (self *SuggestionsViewModel) GetSelected() *types.Suggestion {
	if self.GetItemsLength() == 0 {
		return nil
	}

	return self.getModel()[self.GetSelectedLineIdx()]
}
