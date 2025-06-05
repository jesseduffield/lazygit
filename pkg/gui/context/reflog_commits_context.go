package context

import (
	"time"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type ReflogCommitsContext struct {
	*FilteredListViewModel[*models.Commit]
	*ListContextTrait
}

var (
	_ types.IListContext    = (*ReflogCommitsContext)(nil)
	_ types.DiffableContext = (*ReflogCommitsContext)(nil)
)

func NewReflogCommitsContext(c *ContextCommon) *ReflogCommitsContext {
	viewModel := NewFilteredListViewModel(
		func() []*models.Commit { return c.Model().FilteredReflogCommits },
		func(commit *models.Commit) []string {
			return []string{commit.ShortHash(), commit.Name}
		},
	)

	getDisplayStrings := func(_ int, _ int) [][]string {
		return presentation.GetReflogCommitListDisplayStrings(
			viewModel.GetItems(),
			c.State().GetRepoState().GetScreenMode() != types.SCREEN_NORMAL,
			c.Modes().CherryPicking.SelectedHashSet(),
			c.Modes().Diffing.Ref,
			time.Now(),
			c.UserConfig().Gui.TimeFormat,
			c.UserConfig().Gui.ShortTimeFormat,
			c.UserConfig().Git.ParseEmoji,
		)
	}

	return &ReflogCommitsContext{
		FilteredListViewModel: viewModel,
		ListContextTrait: &ListContextTrait{
			Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
				View:                       c.Views().ReflogCommits,
				WindowName:                 "commits",
				Key:                        REFLOG_COMMITS_CONTEXT_KEY,
				Kind:                       types.SIDE_CONTEXT,
				Focusable:                  true,
				NeedsRerenderOnWidthChange: types.NEEDS_RERENDER_ON_WIDTH_CHANGE_WHEN_SCREEN_MODE_CHANGES,
			})),
			ListRenderer: ListRenderer{
				list:              viewModel,
				getDisplayStrings: getDisplayStrings,
			},
			c: c,
		},
	}
}

func (self *ReflogCommitsContext) CanRebase() bool {
	return false
}

func (self *ReflogCommitsContext) GetSelectedRef() types.Ref {
	commit := self.GetSelected()
	if commit == nil {
		return nil
	}
	return commit
}

func (self *ReflogCommitsContext) GetSelectedRefRangeForDiffFiles() *types.RefRange {
	// It doesn't make much sense to show a range diff between two reflog entries.
	return nil
}

func (self *ReflogCommitsContext) GetCommits() []*models.Commit {
	return self.getModel()
}

func (self *ReflogCommitsContext) GetDiffTerminals() []string {
	itemId := self.GetSelectedItemId()

	return []string{itemId}
}

func (self *ReflogCommitsContext) RefForAdjustingLineNumberInDiff() string {
	return self.GetSelectedItemId()
}

func (self *ReflogCommitsContext) ShowBranchHeadsInSubCommits() bool {
	return false
}
