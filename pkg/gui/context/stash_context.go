package context

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type StashContext struct {
	*BasicViewModel[*models.StashEntry]
	*ListContextTrait
}

var (
	_ types.IListContext    = (*StashContext)(nil)
	_ types.DiffableContext = (*StashContext)(nil)
)

func NewStashContext(
	c *types.HelperCommon,
) *StashContext {
	viewModel := NewBasicViewModel(func() []*models.StashEntry { return c.Model().StashEntries })

	getDisplayStrings := func(startIdx int, length int) [][]string {
		return presentation.GetStashEntryListDisplayStrings(c.Model().StashEntries, c.Modes().Diffing.Ref)
	}

	return &StashContext{
		BasicViewModel: viewModel,
		ListContextTrait: &ListContextTrait{
			Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
				View:       c.Views().Stash,
				WindowName: "stash",
				Key:        STASH_CONTEXT_KEY,
				Kind:       types.SIDE_CONTEXT,
				Focusable:  true,
			})),
			list:              viewModel,
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

func (self *StashContext) GetDiffTerminals() []string {
	itemId := self.GetSelectedItemId()

	return []string{itemId}
}
