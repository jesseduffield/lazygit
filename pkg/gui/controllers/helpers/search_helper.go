package helpers

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// NOTE: this helper supports both filtering and searching. Filtering is when
// the contents of the list are filtered, whereas searching does not actually
// change the contents of the list but instead just highlights the search.
// The general term we use to capture both searching and filtering is...
// 'searching', which is unfortunate but I can't think of a better name.

type SearchHelper struct {
	c *HelperCommon
}

func NewSearchHelper(
	c *HelperCommon,
) *SearchHelper {
	return &SearchHelper{
		c: c,
	}
}

func (self *SearchHelper) OpenFilterPrompt(context types.IFilterableContext) error {
	state := self.searchState()

	state.Context = context

	self.searchPrefixView().SetContent(self.c.Tr.FilterPrefix)
	promptView := self.promptView()
	promptView.ClearTextArea()
	self.OnPromptContentChanged("")
	promptView.RenderTextArea()

	self.c.Context().Push(self.c.Contexts().Search)

	return self.c.ResetKeybindings()
}

func (self *SearchHelper) OpenSearchPrompt(context types.ISearchableContext) error {
	state := self.searchState()

	state.PrevSearchIndex = -1

	state.Context = context

	self.searchPrefixView().SetContent(self.c.Tr.SearchPrefix)
	promptView := self.promptView()
	promptView.ClearTextArea()
	promptView.RenderTextArea()

	self.c.Context().Push(self.c.Contexts().Search)

	return self.c.ResetKeybindings()
}

func (self *SearchHelper) DisplayFilterStatus(context types.IFilterableContext) {
	state := self.searchState()

	state.Context = context
	searchString := context.GetFilter()

	self.searchPrefixView().SetContent(self.c.Tr.FilterPrefix)

	promptView := self.promptView()
	keybindingConfig := self.c.UserConfig().Keybinding
	promptView.SetContent(fmt.Sprintf("matches for '%s' ", searchString) + theme.OptionsFgColor.Sprintf(self.c.Tr.ExitTextFilterMode, keybindings.Label(keybindingConfig.Universal.Return)))
}

func (self *SearchHelper) DisplaySearchStatus(context types.ISearchableContext) {
	state := self.searchState()

	state.Context = context

	self.searchPrefixView().SetContent(self.c.Tr.SearchPrefix)
	index, totalCount := context.GetView().GetSearchStatus()
	context.RenderSearchStatus(index, totalCount)
}

func (self *SearchHelper) searchState() *types.SearchState {
	return self.c.State().GetRepoState().GetSearchState()
}

func (self *SearchHelper) searchPrefixView() *gocui.View {
	return self.c.Views().SearchPrefix
}

func (self *SearchHelper) promptView() *gocui.View {
	return self.c.Contexts().Search.GetView()
}

func (self *SearchHelper) promptContent() string {
	return self.c.Contexts().Search.GetView().TextArea.GetContent()
}

func (self *SearchHelper) Confirm() error {
	state := self.searchState()
	if self.promptContent() == "" {
		return self.CancelPrompt()
	}

	var err error
	switch state.SearchType() {
	case types.SearchTypeFilter:
		self.ConfirmFilter()
	case types.SearchTypeSearch:
		err = self.ConfirmSearch()
	case types.SearchTypeNone:
		self.c.Context().Pop()
	}

	if err != nil {
		return err
	}

	return self.c.ResetKeybindings()
}

func (self *SearchHelper) ConfirmFilter() {
	// We also do this on each keypress but we do it here again just in case
	state := self.searchState()

	context, ok := state.Context.(types.IFilterableContext)
	if !ok {
		self.c.Log.Warnf("Context %s is not filterable", state.Context.GetKey())
		return
	}

	self.OnPromptContentChanged(self.promptContent())
	filterString := self.promptContent()
	if filterString != "" {
		context.GetSearchHistory().Push(filterString)
	}

	self.c.Context().Pop()
}

func (self *SearchHelper) ConfirmSearch() error {
	state := self.searchState()

	context, ok := state.Context.(types.ISearchableContext)
	if !ok {
		self.c.Log.Warnf("Context %s is searchable", state.Context.GetKey())
		return nil
	}

	searchString := self.promptContent()
	context.SetSearchString(searchString)
	if searchString != "" {
		context.GetSearchHistory().Push(searchString)
	}

	self.c.Context().Pop()

	return context.GetView().Search(searchString, modelSearchResults(context))
}

func modelSearchResults(context types.ISearchableContext) []gocui.SearchPosition {
	searchString := context.GetSearchString()

	var normalizedSearchStr string
	// if we have any uppercase characters we'll do a case-sensitive search
	caseSensitive := utils.ContainsUppercase(searchString)
	if caseSensitive {
		normalizedSearchStr = searchString
	} else {
		normalizedSearchStr = strings.ToLower(searchString)
	}

	return context.ModelSearchResults(normalizedSearchStr, caseSensitive)
}

func (self *SearchHelper) CancelPrompt() error {
	self.Cancel()

	self.c.Context().Pop()

	return self.c.ResetKeybindings()
}

func (self *SearchHelper) ScrollHistory(scrollIncrement int) {
	state := self.searchState()

	context, ok := state.Context.(types.ISearchHistoryContext)
	if !ok {
		return
	}

	states := context.GetSearchHistory()

	if val, err := states.PeekAt(state.PrevSearchIndex + scrollIncrement); err == nil {
		state.PrevSearchIndex += scrollIncrement
		promptView := self.promptView()
		promptView.ClearTextArea()
		promptView.TextArea.TypeString(val)
		promptView.RenderTextArea()
		self.OnPromptContentChanged(val)
	}
}

func (self *SearchHelper) Cancel() {
	state := self.searchState()

	switch context := state.Context.(type) {
	case types.IFilterableContext:
		context.ClearFilter()
		self.c.PostRefreshUpdate(context)
	case types.ISearchableContext:
		context.ClearSearchString()
		context.GetView().ClearSearch()
	default:
		// do nothing
	}

	self.HidePrompt()
}

func (self *SearchHelper) OnPromptContentChanged(searchString string) {
	state := self.searchState()
	switch context := state.Context.(type) {
	case types.IFilterableContext:
		context.SetSelection(0)
		context.GetView().SetOriginY(0)
		context.SetFilter(searchString, self.c.UserConfig().Gui.UseFuzzySearch())
		self.c.PostRefreshUpdate(context)
	case types.ISearchableContext:
		// do nothing
	default:
		// do nothing (shouldn't land here)
	}
}

func (self *SearchHelper) ReApplyFilter(context types.Context) {
	filterableContext, ok := context.(types.IFilterableContext)
	if ok {
		state := self.searchState()
		if context == state.Context {
			filterableContext.SetSelection(0)
			filterableContext.GetView().SetOriginY(0)
		}
		filterableContext.ReApplyFilter(self.c.UserConfig().Gui.UseFuzzySearch())
	}
}

func (self *SearchHelper) ReApplySearch(ctx types.Context) {
	// Reapply the search if the model has changed. This is needed for contexts
	// that use the model for searching, to pass the new model search positions
	// to the view.
	searchableContext, ok := ctx.(types.ISearchableContext)
	if ok {
		ctx.GetView().UpdateSearchResults(searchableContext.GetSearchString(), modelSearchResults(searchableContext))

		state := self.searchState()
		if ctx == state.Context {
			// Re-render the "x of y" search status, unless the search prompt is
			// open for typing.
			if self.c.Context().Current().GetKey() != context.SEARCH_CONTEXT_KEY {
				self.RenderSearchStatus(searchableContext)
			}
		}
	}
}

func (self *SearchHelper) RenderSearchStatus(c types.Context) {
	if c.GetKey() == context.SEARCH_CONTEXT_KEY {
		return
	}

	if searchableContext, ok := c.(types.ISearchableContext); ok {
		if searchableContext.IsSearching() {
			self.setSearchingFrameColor()
			self.DisplaySearchStatus(searchableContext)
			return
		}
	}
	if filterableContext, ok := c.(types.IFilterableContext); ok {
		if filterableContext.IsFiltering() {
			self.setSearchingFrameColor()
			self.DisplayFilterStatus(filterableContext)
			return
		}
	}

	self.HidePrompt()
}

func (self *SearchHelper) CancelSearchIfSearching(c types.Context) {
	if searchableContext, ok := c.(types.ISearchableContext); ok {
		view := searchableContext.GetView()
		if view != nil && view.IsSearching() {
			view.ClearSearch()
			searchableContext.ClearSearchString()
			self.Cancel()
		}
		return
	}

	if filterableContext, ok := c.(types.IFilterableContext); ok {
		if filterableContext.IsFiltering() {
			filterableContext.ClearFilter()
			self.Cancel()
		}
		return
	}
}

func (self *SearchHelper) HidePrompt() {
	self.setNonSearchingFrameColor()

	state := self.searchState()
	state.Context = nil
}

func (self *SearchHelper) setSearchingFrameColor() {
	self.c.GocuiGui().SelFgColor = theme.SearchingActiveBorderColor
	self.c.GocuiGui().SelFrameColor = theme.SearchingActiveBorderColor
}

func (self *SearchHelper) setNonSearchingFrameColor() {
	self.c.GocuiGui().SelFgColor = theme.ActiveBorderColor
	self.c.GocuiGui().SelFrameColor = theme.ActiveBorderColor
}
