package context

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type RemotesContext struct {
	*FilteredListViewModel[*models.Remote]
	*ListContextTrait
}

var (
	_ types.IListContext    = (*RemotesContext)(nil)
	_ types.DiffableContext = (*RemotesContext)(nil)
)

func NewRemotesContext(c *ContextCommon) *RemotesContext {
	viewModel := NewFilteredListViewModel(
		func() []*models.Remote { return c.Model().Remotes },
		func(remote *models.Remote) []string {
			return []string{remote.Name}
		},
	)

	getDisplayStrings := func(_ int, _ int) [][]string {
		return presentation.GetRemoteListDisplayStrings(
			viewModel.GetItems(), c.Modes().Diffing.Ref, c.State().GetItemOperation, c.Tr, c.UserConfig())
	}

	return &RemotesContext{
		FilteredListViewModel: viewModel,
		ListContextTrait: &ListContextTrait{
			Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
				View:       c.Views().Remotes,
				WindowName: "branches",
				Key:        REMOTES_CONTEXT_KEY,
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

func (self *RemotesContext) GetDiffTerminals() []string {
	itemId := self.GetSelectedItemId()

	return []string{itemId}
}
