package context

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type StashContext struct {
	*FilteredListViewModel[*models.StashEntry]
	*ListContextTrait
}

var (
	_ types.IListContext    = (*StashContext)(nil)
	_ types.DiffableContext = (*StashContext)(nil)
)

func NewStashContext(
	c *ContextCommon,
) *StashContext {
	viewModel := NewFilteredListViewModel(
		func() []*models.StashEntry { return c.Model().StashEntries },
		func(stashEntry *models.StashEntry) []string {
			return []string{stashEntry.Name}
		},
	)

	getDisplayStrings := func(_ int, _ int) [][]string {
		return presentation.GetStashEntryListDisplayStrings(viewModel.GetItems(), c.Modes().Diffing.Ref)
	}

	return &StashContext{
		FilteredListViewModel: viewModel,
		ListContextTrait: &ListContextTrait{
			Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
				View:       c.Views().Stash,
				WindowName: "stash",
				Key:        STASH_CONTEXT_KEY,
				Kind:       types.SIDE_CONTEXT,
				Focusable:  true,
			})),
			ListRenderer: ListRenderer{
				list:              viewModel,
				getDisplayStrings: getDisplayStrings,
			},
			c: c,
		},
	}
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

func (self *StashContext) GetSelectedRefRangeForDiffFiles() *types.RefRange {
	// It doesn't make much sense to show a range diff between two stash entries.
	return nil
}

func (self *StashContext) GetDiffTerminals() []string {
	itemId := self.GetSelectedItemId()

	return []string{itemId}
}

func (self *StashContext) RefForAdjustingLineNumberInDiff() string {
	return self.GetSelectedItemId()
}
