package context

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type SubmodulesContext struct {
	*FilteredListViewModel[*models.SubmoduleConfig]
	*ListContextTrait
}

var _ types.IListContext = (*SubmodulesContext)(nil)

func NewSubmodulesContext(c *ContextCommon) *SubmodulesContext {
	viewModel := NewFilteredListViewModel(
		func() []*models.SubmoduleConfig { return c.Model().Submodules },
		func(submodule *models.SubmoduleConfig) []string {
			return []string{submodule.FullName()}
		},
	)

	getDisplayStrings := func(_ int, _ int) [][]string {
		return presentation.GetSubmoduleListDisplayStrings(viewModel.GetItems())
	}

	return &SubmodulesContext{
		FilteredListViewModel: viewModel,
		ListContextTrait: &ListContextTrait{
			Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
				View:       c.Views().Submodules,
				WindowName: "files",
				Key:        SUBMODULES_CONTEXT_KEY,
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
