package context

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

type WorktreesContext struct {
	*FilteredListViewModel[*models.Worktree]
	*ListContextTrait
}

var _ types.IListContext = (*WorktreesContext)(nil)

func NewWorktreesContext(c *ContextCommon) *WorktreesContext {
	viewModel := NewFilteredListViewModel(
		func() []*models.Worktree { return c.Model().Worktrees },
		func(Worktree *models.Worktree) []string {
			return []string{Worktree.Name()}
		},
	)

	getDisplayStrings := func(startIdx int, length int) [][]string {
		return lo.Map(c.Model().Worktrees, func(worktree *models.Worktree, _ int) []string {
			return presentation.GetWorktreeDisplayString(
				c.Git().Worktree.IsCurrentWorktree(worktree),
				c.Git().Worktree.IsWorktreePathMissing(worktree),
				worktree)
		})
	}

	return &WorktreesContext{
		FilteredListViewModel: viewModel,
		ListContextTrait: &ListContextTrait{
			Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
				View:       c.Views().Worktrees,
				WindowName: "branches",
				Key:        WORKTREES_CONTEXT_KEY,
				Kind:       types.SIDE_CONTEXT,
				Focusable:  true,
			})),
			list:              viewModel,
			getDisplayStrings: getDisplayStrings,
			c:                 c,
		},
	}
}

func (self *WorktreesContext) GetSelectedItemId() string {
	item := self.GetSelected()
	if item == nil {
		return ""
	}

	return item.ID()
}
