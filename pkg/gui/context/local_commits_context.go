package context

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

type LocalCommitsContext struct {
	*LocalCommitsViewModel
	*ListContextTrait
	*SearchTrait
}

var (
	_ types.IListContext       = (*LocalCommitsContext)(nil)
	_ types.DiffableContext    = (*LocalCommitsContext)(nil)
	_ types.ISearchableContext = (*LocalCommitsContext)(nil)
)

func NewLocalCommitsContext(c *ContextCommon) *LocalCommitsContext {
	viewModel := NewLocalCommitsViewModel(
		func() []*models.Commit { return c.Model().Commits },
		c,
	)

	getDisplayStrings := func(startIdx int, endIdx int) [][]string {
		var selectedCommitHashPtr *string

		if c.Context().Current().GetKey() == LOCAL_COMMITS_CONTEXT_KEY {
			selectedCommit := viewModel.GetSelected()
			if selectedCommit != nil {
				selectedCommitHashPtr = selectedCommit.HashPtr()
			}
		}

		hasRebaseUpdateRefsConfig := c.Git().Config.GetRebaseUpdateRefs()

		return presentation.GetCommitListDisplayStrings(
			c.Common,
			c.Model().Commits,
			c.Model().Branches,
			c.Model().CheckedOutBranch,
			hasRebaseUpdateRefsConfig,
			c.State().GetRepoState().GetScreenMode() != types.SCREEN_NORMAL,
			c.Modes().CherryPicking.SelectedHashSet(),
			c.Modes().Diffing.Ref,
			c.Modes().MarkedBaseCommit.GetHash(),
			c.UserConfig().Gui.TimeFormat,
			c.UserConfig().Gui.ShortTimeFormat,
			time.Now(),
			c.UserConfig().Git.ParseEmoji,
			selectedCommitHashPtr,
			startIdx,
			endIdx,
			shouldShowGraph(c),
			c.Model().BisectInfo,
		)
	}

	getNonModelItems := func() []*NonModelItem {
		result := []*NonModelItem{}
		if c.Model().WorkingTreeStateAtLastCommitRefresh.CanShowTodos() {
			if c.Model().WorkingTreeStateAtLastCommitRefresh.Rebasing {
				result = append(result, &NonModelItem{
					Index:   0,
					Content: fmt.Sprintf("--- %s ---", c.Tr.PendingRebaseTodosSectionHeader),
				})
			}

			if c.Model().WorkingTreeStateAtLastCommitRefresh.CherryPicking ||
				c.Model().WorkingTreeStateAtLastCommitRefresh.Reverting {
				_, firstCherryPickOrRevertTodo, found := lo.FindIndexOf(
					c.Model().Commits, func(c *models.Commit) bool {
						return c.Status == models.StatusCherryPickingOrReverting ||
							c.Status == models.StatusConflicted
					})
				if !found {
					firstCherryPickOrRevertTodo = 0
				}
				label := lo.Ternary(c.Model().WorkingTreeStateAtLastCommitRefresh.CherryPicking,
					c.Tr.PendingCherryPicksSectionHeader,
					c.Tr.PendingRevertsSectionHeader)
				result = append(result, &NonModelItem{
					Index:   firstCherryPickOrRevertTodo,
					Content: fmt.Sprintf("--- %s ---", label),
				})
			}

			_, firstRealCommit, found := lo.FindIndexOf(
				c.Model().Commits, func(c *models.Commit) bool {
					return !c.IsTODO()
				})
			if !found {
				firstRealCommit = 0
			}
			result = append(result, &NonModelItem{
				Index:   firstRealCommit,
				Content: fmt.Sprintf("--- %s ---", c.Tr.CommitsSectionHeader),
			})
		}

		return result
	}

	ctx := &LocalCommitsContext{
		LocalCommitsViewModel: viewModel,
		SearchTrait:           NewSearchTrait(c),
		ListContextTrait: &ListContextTrait{
			Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
				View:                        c.Views().Commits,
				WindowName:                  "commits",
				Key:                         LOCAL_COMMITS_CONTEXT_KEY,
				Kind:                        types.SIDE_CONTEXT,
				Focusable:                   true,
				NeedsRerenderOnWidthChange:  types.NEEDS_RERENDER_ON_WIDTH_CHANGE_WHEN_SCREEN_MODE_CHANGES,
				NeedsRerenderOnHeightChange: true,
			})),
			ListRenderer: ListRenderer{
				list:              viewModel,
				getDisplayStrings: getDisplayStrings,
				getNonModelItems:  getNonModelItems,
			},
			c:                       c,
			refreshViewportOnChange: true,
			renderOnlyVisibleLines:  true,
		},
	}

	ctx.GetView().SetOnSelectItem(ctx.SearchTrait.onSelectItemWrapper(ctx.OnSearchSelect))

	return ctx
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
		showWholeGitGraph: c.UserConfig().Git.Log.ShowWholeGraph,
	}

	return self
}

func (self *LocalCommitsContext) CanRebase() bool {
	return true
}

func (self *LocalCommitsContext) GetSelectedRef() models.Ref {
	commit := self.GetSelected()
	if commit == nil {
		return nil
	}
	return commit
}

func (self *LocalCommitsContext) GetSelectedRefRangeForDiffFiles() *types.RefRange {
	commits, startIdx, endIdx := self.GetSelectedItems()
	if commits == nil || startIdx == endIdx {
		return nil
	}
	from := commits[len(commits)-1]
	to := commits[0]
	if from.IsTODO() || to.IsTODO() {
		return nil
	}
	return &types.RefRange{From: from, To: to}
}

// Returns the commit hash of the selected commit, or an empty string if no
// commit is selected
func (self *LocalCommitsContext) GetSelectedCommitHash() string {
	commit := self.GetSelected()
	if commit == nil {
		return ""
	}
	return commit.Hash()
}

func (self *LocalCommitsContext) SelectCommitByHash(hash string) bool {
	if hash == "" {
		return false
	}

	if _, idx, found := lo.FindIndexOf(self.GetItems(), func(c *models.Commit) bool { return c.Hash() == hash }); found {
		self.SetSelection(idx)
		return true
	}

	return false
}

func (self *LocalCommitsContext) GetDiffTerminals() []string {
	itemId := self.GetSelectedItemId()

	return []string{itemId}
}

func (self *LocalCommitsContext) RefForAdjustingLineNumberInDiff() string {
	commits, _, _ := self.GetSelectedItems()
	if commits == nil {
		return ""
	}
	return commits[0].Hash()
}

func (self *LocalCommitsContext) ModelSearchResults(searchStr string, caseSensitive bool) []gocui.SearchPosition {
	return searchModelCommits(caseSensitive, self.GetCommits(), self.ColumnPositions(), self.ModelIndexToViewIndex, searchStr)
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

	value := c.UserConfig().Git.Log.ShowGraph

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

func searchModelCommits(caseSensitive bool, commits []*models.Commit, columnPositions []int,
	modelToViewIndex func(int) int, searchStr string,
) []gocui.SearchPosition {
	if columnPositions == nil {
		// This should never happen. We are being called at a time where our
		// entire view content is scrolled out of view, so that we didn't draw
		// anything the last time we rendered. If we run into a scenario where
		// this happens, we should fix it, but until we found them all, at least
		// make sure we don't crash.
		return []gocui.SearchPosition{}
	}

	normalize := lo.Ternary(caseSensitive, func(s string) string { return s }, strings.ToLower)
	return lo.FilterMap(commits, func(commit *models.Commit, idx int) (gocui.SearchPosition, bool) {
		// The XStart and XEnd values are only used if the search string can't
		// be found in the view. This can really only happen if the user is
		// searching for a commit hash that is longer than the truncated hash
		// that we render. So we just set the XStart and XEnd values to the
		// start and end of the commit hash column, which is the second one.
		result := gocui.SearchPosition{XStart: columnPositions[1], XEnd: columnPositions[2] - 1, Y: modelToViewIndex(idx)}
		return result, strings.Contains(normalize(commit.Hash()), searchStr) ||
			strings.Contains(normalize(commit.Name), searchStr) ||
			strings.Contains(normalize(commit.ExtraInfo), searchStr) // allow searching for tags
	})
}

func (self *LocalCommitsContext) IndexForGotoBottom() int {
	commits := self.GetCommits()
	selectedIdx := self.GetSelectedLineIdx()
	if selectedIdx >= 0 && selectedIdx < len(commits)-1 {
		if commits[selectedIdx+1].Status != models.StatusMerged {
			_, idx, found := lo.FindIndexOf(commits, func(c *models.Commit) bool {
				return c.Status == models.StatusMerged
			})
			if found {
				return idx - 1
			}
		}
	}

	return self.list.Len() - 1
}
