package context

import (
	"log"
	"time"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/types/enums"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type LocalCommitsContext struct {
	*LocalCommitsViewModel
	*ListContextTrait
	*SearchTrait
}

var (
	_ types.IListContext    = (*LocalCommitsContext)(nil)
	_ types.DiffableContext = (*LocalCommitsContext)(nil)
)

func NewLocalCommitsContext(c *ContextCommon) *LocalCommitsContext {
	viewModel := NewLocalCommitsViewModel(
		func() []*models.Commit { return c.Model().Commits },
		c,
	)

	getDisplayStrings := func(startIdx int, length int) [][]string {
		selectedCommitSha := ""

		if c.CurrentContext().GetKey() == LOCAL_COMMITS_CONTEXT_KEY {
			selectedCommit := viewModel.GetSelected()
			if selectedCommit != nil {
				selectedCommitSha = selectedCommit.Sha
			}
		}

		showYouAreHereLabel := c.Model().WorkingTreeStateAtLastCommitRefresh == enums.REBASE_MODE_REBASING

		return presentation.GetCommitListDisplayStrings(
			c.Common,
			c.Model().Commits,
			c.State().GetRepoState().GetScreenMode() != types.SCREEN_NORMAL,
			c.Modes().CherryPicking.SelectedShaSet(),
			c.Modes().Diffing.Ref,
			c.UserConfig.Gui.TimeFormat,
			c.UserConfig.Gui.ShortTimeFormat,
			time.Now(),
			c.UserConfig.Git.ParseEmoji,
			selectedCommitSha,
			startIdx,
			length,
			shouldShowGraph(c),
			c.Model().BisectInfo,
			showYouAreHereLabel,
		)
	}

	ctx := &LocalCommitsContext{
		LocalCommitsViewModel: viewModel,
		SearchTrait:           NewSearchTrait(c),
		ListContextTrait: &ListContextTrait{
			Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
				View:       c.Views().Commits,
				WindowName: "commits",
				Key:        LOCAL_COMMITS_CONTEXT_KEY,
				Kind:       types.SIDE_CONTEXT,
				Focusable:  true,
			})),
			list:                    viewModel,
			getDisplayStrings:       getDisplayStrings,
			c:                       c,
			refreshViewportOnChange: true,
		},
	}

	ctx.GetView().SetOnSelectItem(ctx.SearchTrait.onSelectItemWrapper(func(selectedLineIdx int) error {
		ctx.GetList().SetSelectedLineIdx(selectedLineIdx)
		return ctx.HandleFocus(types.OnFocusOpts{})
	}))

	return ctx
}

func (self *LocalCommitsContext) GetSelectedItemId() string {
	item := self.GetSelected()
	if item == nil {
		return ""
	}

	return item.ID()
}

type LocalCommitsViewModel struct {
	*ListViewModel[*models.Commit]

	// If this is true we limit the amount of commits we load, for the sake of keeping things fast.
	// If the user attempts to scroll past the end of the list, we will load more commits.
	limitCommits bool

	// If this is true we'll use git log --all when fetching the commits.
	showWholeGitGraph bool
}

func NewLocalCommitsViewModel(getModel func() []*models.Commit, c *ContextCommon) *LocalCommitsViewModel {
	self := &LocalCommitsViewModel{
		ListViewModel:     NewListViewModel(getModel),
		limitCommits:      true,
		showWholeGitGraph: c.UserConfig.Git.Log.ShowWholeGraph,
	}

	return self
}

func (self *LocalCommitsContext) CanRebase() bool {
	return true
}

func (self *LocalCommitsContext) GetSelectedRef() types.Ref {
	commit := self.GetSelected()
	if commit == nil {
		return nil
	}
	return commit
}

func (self *LocalCommitsContext) GetDiffTerminals() []string {
	itemId := self.GetSelectedItemId()

	return []string{itemId}
}

func (self *LocalCommitsViewModel) SetLimitCommits(value bool) {
	self.limitCommits = value
}

func (self *LocalCommitsViewModel) GetLimitCommits() bool {
	return self.limitCommits
}

func (self *LocalCommitsViewModel) SetShowWholeGitGraph(value bool) {
	self.showWholeGitGraph = value
}

func (self *LocalCommitsViewModel) GetShowWholeGitGraph() bool {
	return self.showWholeGitGraph
}

func (self *LocalCommitsViewModel) GetCommits() []*models.Commit {
	return self.getModel()
}

func shouldShowGraph(c *ContextCommon) bool {
	if c.Modes().Filtering.Active() {
		return false
	}

	value := c.UserConfig.Git.Log.ShowGraph
	switch value {
	case "always":
		return true
	case "never":
		return false
	case "when-maximised":
		return c.State().GetRepoState().GetScreenMode() != types.SCREEN_NORMAL
	}

	log.Fatalf("Unknown value for git.log.showGraph: %s. Expected one of: 'always', 'never', 'when-maximised'", value)
	return false
}
