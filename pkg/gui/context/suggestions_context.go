package context

import (
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/tasks"
)

type SuggestionsContext struct {
	*ListViewModel[*types.Suggestion]
	*ListContextTrait

	State *SuggestionsContextState
}

type SuggestionsContextState struct {
	Suggestions        []*types.Suggestion
	OnConfirm          func() error
	OnClose            func() error
	OnDeleteSuggestion func() error
	AsyncHandler       *tasks.AsyncHandler

	AllowEditSuggestion bool

	// FindSuggestions will take a string that the user has typed into a prompt
	// and return a slice of suggestions which match that string.
	FindSuggestions func(string) []*types.Suggestion
}

var _ types.IListContext = (*SuggestionsContext)(nil)

func NewSuggestionsContext(
	c *ContextCommon,
) *SuggestionsContext {
	state := &SuggestionsContextState{
		AsyncHandler: tasks.NewAsyncHandler(c.OnWorker),
	}
	getModel := func() []*types.Suggestion {
		return state.Suggestions
	}

	getDisplayStrings := func(_ int, _ int) [][]string {
		return presentation.GetSuggestionListDisplayStrings(state.Suggestions)
	}

	viewModel := NewListViewModel(getModel)

	return &SuggestionsContext{
		State:         state,
		ListViewModel: viewModel,
		ListContextTrait: &ListContextTrait{
			Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
				View:                  c.Views().Suggestions,
				WindowName:            "suggestions",
				Key:                   SUGGESTIONS_CONTEXT_KEY,
				Kind:                  types.PERSISTENT_POPUP,
				Focusable:             true,
				HasUncontrolledBounds: true,
			})),
			ListRenderer: ListRenderer{
				list:              viewModel,
				getDisplayStrings: getDisplayStrings,
			},
			c: c,
		},
	}
}

func (self *SuggestionsContext) SetSuggestions(suggestions []*types.Suggestion) {
	self.State.Suggestions = suggestions
	self.SetSelection(0)
	self.c.ResetViewOrigin(self.GetView())
	self.HandleRender()
}

func (self *SuggestionsContext) RefreshSuggestions() {
	self.State.AsyncHandler.Do(func() func() {
		findSuggestionsFn := self.State.FindSuggestions
		if findSuggestionsFn != nil {
			suggestions := findSuggestionsFn(self.c.GetPromptInput())
			return func() { self.SetSuggestions(suggestions) }
		} else {
			return func() {}
		}
	})
}

// There is currently no need to use range-select in the suggestions view so we're disabling it.
func (self *SuggestionsContext) RangeSelectEnabled() bool {
	return false
}
