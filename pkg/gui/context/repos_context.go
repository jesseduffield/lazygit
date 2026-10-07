package context

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type ReposContext struct {
	*FilteredListViewModel[*models.Repo]
	*ListContextTrait
}

var _ types.IListContext = (*ReposContext)(nil)

func NewReposContext(c *ContextCommon) *ReposContext {
	viewModel := NewFilteredListViewModel(
		func() []*models.Repo { return c.State().GetRepoList() },
		func(repo *models.Repo) []string {
			return []string{repo.Name, repo.Branch}
		},
	)

	getDisplayStrings := func(_ int, _ int) [][]string {
		return presentation.GetRepoListDisplayStrings(
			viewModel.GetFilteredList(),
			c.Git().RepoPaths.WorktreePath(),
		)
	}

	return &ReposContext{
		FilteredListViewModel: viewModel,
		ListContextTrait: &ListContextTrait{
			Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
				View:       c.Views().Repos,
				WindowName: "repos",
				Key:        REPOS_CONTEXT_KEY,
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

// SelectRepo moves the selection to the repo with the given path, if the
// filtered list shows it.
func (self *ReposContext) SelectRepo(path string) {
	for i, repo := range self.GetFilteredList() {
		if repo.Path == path {
			self.SetSelection(i)
			return
		}
	}
}
