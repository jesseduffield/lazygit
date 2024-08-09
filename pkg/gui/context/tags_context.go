package context

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type TagsContext struct {
	*FilteredListViewModel[*models.Tag]
	*ListContextTrait
}

var (
	_ types.IListContext    = (*TagsContext)(nil)
	_ types.DiffableContext = (*TagsContext)(nil)
)

func NewTagsContext(
	c *ContextCommon,
) *TagsContext {
	viewModel := NewFilteredListViewModel(
		func() []*models.Tag { return c.Model().Tags },
		func(tag *models.Tag) []string {
			return []string{tag.Name, tag.Message}
		},
	)

	getDisplayStrings := func(_ int, _ int) [][]string {
		return presentation.GetTagListDisplayStrings(
			viewModel.GetItems(),
			c.State().GetItemOperation,
			c.Modes().Diffing.Ref, c.Tr, c.UserConfig())
	}

	return &TagsContext{
		FilteredListViewModel: viewModel,
		ListContextTrait: &ListContextTrait{
			Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
				View:       c.Views().Tags,
				WindowName: "branches",
				Key:        TAGS_CONTEXT_KEY,
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

func (self *TagsContext) GetSelectedRef() types.Ref {
	tag := self.GetSelected()
	if tag == nil {
		return nil
	}
	return tag
}

func (self *TagsContext) GetDiffTerminals() []string {
	itemId := self.GetSelectedItemId()

	return []string{itemId}
}

func (self *TagsContext) ShowBranchHeadsInSubCommits() bool {
	return true
}
