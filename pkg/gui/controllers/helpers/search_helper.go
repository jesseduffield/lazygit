package helpers

import (
	"fmt"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/theme"
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
	promptView.TextArea.TypeString(context.GetFilter())
	promptView.RenderTextArea()

	if err := self.c.PushContext(self.c.Contexts().Search); err != nil {
		return err
	}

	return nil
}

func (self *SearchHelper) OpenSearchPrompt(context types.ISearchableContext) error {
	state := self.searchState()

	state.Context = context
	searchString := context.GetSearchString()

	self.searchPrefixView().SetContent(self.c.Tr.SearchPrefix)
	promptView := self.promptView()
	promptView.ClearTextArea()
	promptView.TextArea.TypeString(searchString)
	promptView.RenderTextArea()

	if err := self.c.PushContext(self.c.Contexts().Search); err != nil {
		return err
	}

	return nil
}

func (self *SearchHelper) DisplayFilterStatus(context types.IFilterableContext) {
	state := self.searchState()

	state.Context = context
	searchString := context.GetFilter()

	self.searchPrefixView().SetContent(self.c.Tr.FilterPrefix)

	promptView := self.promptView()
	keybindingConfig := self.c.UserConfig.Keybinding
	promptView.SetContent(fmt.Sprintf("matches for '%s' ", searchString) + theme.OptionsFgColor.Sprintf(self.c.Tr.ExitTextFilterMode, keybindings.Label(keybindingConfig.Universal.Return)))
}

func (self *SearchHelper) DisplaySearchStatus(context types.ISearchableContext) {
	state := self.searchState()

	state.Context = context

	self.searchPrefixView().SetContent(self.c.Tr.SearchPrefix)
	_ = context.GetView().SelectCurrentSearchResult()
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

	switch state.SearchType() {
	case types.SearchTypeFilter:
		return self.ConfirmFilter()
	case types.SearchTypeSearch:
		return self.ConfirmSearch()
	case types.SearchTypeNone:
		return self.c.PopContext()
	}

	return nil
}

func (self *SearchHelper) ConfirmFilter() error {
	// We also do this on each keypress but we do it here again just in case
	state := self.searchState()

	_, ok := state.Context.(types.IFilterableContext)
	if !ok {
		self.c.Log.Warnf("Context %s is not filterable", state.Context.GetKey())
		return nil
	}

	self.OnPromptContentChanged(self.promptContent())

	return self.c.PopContext()
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

	view := context.GetView()

	if err := self.c.PopContext(); err != nil {
		return err
	}

	if err := view.Search(searchString); err != nil {
		return err
	}

	return nil
}

func (self *SearchHelper) CancelPrompt() error {
	self.Cancel()

	return self.c.PopContext()
}

func (self *SearchHelper) Cancel() {
	state := self.searchState()

	switch context := state.Context.(type) {
	case types.IFilterableContext:
		context.ClearFilter()
		_ = self.c.PostRefreshUpdate(context)
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
		context.SetSelectedLineIdx(0)
		_ = context.GetView().SetOriginY(0)
		context.SetFilter(searchString)
		_ = self.c.PostRefreshUpdate(context)
	case types.ISearchableContext:
		// do nothing
	default:
		// do nothing (shouldn't land here)
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
