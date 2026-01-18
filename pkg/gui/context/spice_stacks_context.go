package context

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type SpiceStacksContext struct {
	*FilteredListViewModel[*models.SpiceStackItem]
	*ListContextTrait
}

var (
	_ types.IListContext    = (*SpiceStacksContext)(nil)
	_ types.DiffableContext = (*SpiceStacksContext)(nil)
)

func NewSpiceStacksContext(c *ContextCommon) *SpiceStacksContext {
	viewModel := NewFilteredListViewModel(
		func() []*models.SpiceStackItem {
			// Return ALL items (including commits) so indices match the display
			return c.Model().SpiceStackItems
		},
		func(item *models.SpiceStackItem) []string {
			return []string{item.Name}
		},
	)

	getDisplayStrings := func(_ int, _ int) [][]string {
		return presentation.GetSpiceStackDisplayStrings(
			c.Model().SpiceStackItems,
			c.State().GetItemOperation,
			c.Modes().Diffing.Ref,
			c.Tr,
			c.UserConfig(),
		)
	}

	return &SpiceStacksContext{
		FilteredListViewModel: viewModel,
		ListContextTrait: &ListContextTrait{
			Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
				View:       c.Views().SpiceStacks,
				WindowName: "branches",
				Key:        SPICE_STACKS_CONTEXT_KEY,
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

func (self *SpiceStacksContext) GetSelectedRef() models.Ref {
	item := self.GetSelected()
	if item == nil {
		return nil
	}
	return item
}

func (self *SpiceStacksContext) GetDiffTerminals() []string {
	itemId := self.GetSelectedItemId()
	return []string{itemId}
}

func (self *SpiceStacksContext) RefForAdjustingLineNumberInDiff() string {
	return self.GetSelectedItemId()
}

func (self *SpiceStacksContext) ShowBranchHeadsInSubCommits() bool {
	return true
}
