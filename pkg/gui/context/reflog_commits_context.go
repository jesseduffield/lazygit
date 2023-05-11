package context

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type ReflogCommitsContext struct {
	*BasicViewModel[*models.Commit]
	*ListContextTrait
}

var (
	_ types.IListContext    = (*ReflogCommitsContext)(nil)
	_ types.DiffableContext = (*ReflogCommitsContext)(nil)
)

func NewReflogCommitsContext(c *ContextCommon) *ReflogCommitsContext {
	viewModel := NewBasicViewModel(func() []*models.Commit { return c.Model().FilteredReflogCommits })

	getDisplayStrings := func(startIdx int, length int) [][]string {
		return presentation.GetReflogCommitListDisplayStrings(
			c.Model().FilteredReflogCommits,
			c.State().GetRepoState().GetScreenMode() != types.SCREEN_NORMAL,
			c.Modes().CherryPicking.SelectedShaSet(),
			c.Modes().Diffing.Ref,
			c.UserConfig.Gui.TimeFormat,
			c.UserConfig.Git.ParseEmoji,
		)
	}

	return &ReflogCommitsContext{
		BasicViewModel: viewModel,
		ListContextTrait: &ListContextTrait{
			Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
				View:       c.Views().ReflogCommits,
				WindowName: "commits",
				Key:        REFLOG_COMMITS_CONTEXT_KEY,
				Kind:       types.SIDE_CONTEXT,
				Focusable:  true,
			})),
			list:              viewModel,
			getDisplayStrings: getDisplayStrings,
			c:                 c,
		},
	}
}

func (self *ReflogCommitsContext) GetSelectedItemId() string {
	item := self.GetSelected()
	if item == nil {
		return ""
	}

	return item.ID()
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

func (self *ReflogCommitsContext) GetCommits() []*models.Commit {
	return self.getModel()
}

func (self *ReflogCommitsContext) GetDiffTerminals() []string {
	itemId := self.GetSelectedItemId()

	return []string{itemId}
}
