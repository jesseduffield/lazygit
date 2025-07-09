package context

import (
	"fmt"
	"time"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

type SubCommitsContext struct {
	c *ContextCommon

	*SubCommitsViewModel
	*ListContextTrait
	*DynamicTitleBuilder
	*SearchTrait
}

var (
	_ types.IListContext       = (*SubCommitsContext)(nil)
	_ types.DiffableContext    = (*SubCommitsContext)(nil)
	_ types.ISearchableContext = (*SubCommitsContext)(nil)
)

func NewSubCommitsContext(
	c *ContextCommon,
) *SubCommitsContext {
	viewModel := &SubCommitsViewModel{
		ListViewModel: NewListViewModel(
			func() []*models.Commit { return c.Model().SubCommits },
		),
		ref:          nil,
		limitCommits: true,
	}

	getDisplayStrings := func(startIdx int, endIdx int) [][]string {
		// This can happen if a sub-commits view is asked to be rerendered while
		// it is invisible; for example when switching screen modes, which
		// rerenders all views.
		if viewModel.GetRef() == nil {
			return [][]string{}
		}

		var selectedCommitHashPtr *string
		if c.Context().Current().GetKey() == SUB_COMMITS_CONTEXT_KEY {
			selectedCommit := viewModel.GetSelected()
			if selectedCommit != nil {
				selectedCommitHashPtr = selectedCommit.HashPtr()
			}
		}
		branches := []*models.Branch{}
		if viewModel.GetShowBranchHeads() {
			branches = c.Model().Branches
		}
		hasRebaseUpdateRefsConfig := c.Git().Config.GetRebaseUpdateRefs()
		return presentation.GetCommitListDisplayStrings(
			c.Common,
			c.Model().SubCommits,
			branches,
			viewModel.GetRef().RefName(),
			hasRebaseUpdateRefsConfig,
			c.State().GetRepoState().GetScreenMode() != types.SCREEN_NORMAL,
			c.Modes().CherryPicking.SelectedHashSet(),
			c.Modes().Diffing.Ref,
			"",
			c.UserConfig().Gui.TimeFormat,
			c.UserConfig().Gui.ShortTimeFormat,
			time.Now(),
			c.UserConfig().Git.ParseEmoji,
			selectedCommitHashPtr,
			startIdx,
			endIdx,
			shouldShowGraph(c),
			shouldDisplayForScreenMode(c, viewModel.GetShowTags()),
			git_commands.NewNullBisectInfo(),
		)
	}

	getNonModelItems := func() []*NonModelItem {
		result := []*NonModelItem{}
		if viewModel.GetRefToShowDivergenceFrom() != "" {
			_, upstreamIdx, found := lo.FindIndexOf(
				c.Model().SubCommits, func(c *models.Commit) bool { return c.Divergence == models.DivergenceRight })
			if !found {
				upstreamIdx = 0
			}
			result = append(result, &NonModelItem{
				Index:   upstreamIdx,
				Content: fmt.Sprintf("--- %s ---", c.Tr.DivergenceSectionHeaderRemote),
			})

			_, localIdx, found := lo.FindIndexOf(
				c.Model().SubCommits, func(c *models.Commit) bool { return c.Divergence == models.DivergenceLeft })
			if !found {
				localIdx = len(c.Model().SubCommits)
			}
			result = append(result, &NonModelItem{
				Index:   localIdx,
				Content: fmt.Sprintf("--- %s ---", c.Tr.DivergenceSectionHeaderLocal),
			})
		}

		return result
	}

	ctx := &SubCommitsContext{
		c:                   c,
		SubCommitsViewModel: viewModel,
		SearchTrait:         NewSearchTrait(c),
		DynamicTitleBuilder: NewDynamicTitleBuilder(c.Tr.SubCommitsDynamicTitle),
		ListContextTrait: &ListContextTrait{
			Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
				View:                        c.Views().SubCommits,
				WindowName:                  "branches",
				Key:                         SUB_COMMITS_CONTEXT_KEY,
				Kind:                        types.SIDE_CONTEXT,
				Focusable:                   true,
				Transient:                   true,
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

type SubCommitsViewModel struct {
	// name of the ref that the sub-commits are shown for
	ref                     types.Ref
	refToShowDivergenceFrom string
	*ListViewModel[*models.Commit]

	limitCommits    bool
	showBranchHeads bool
	showTags        string
}

func (self *SubCommitsViewModel) SetRef(ref types.Ref) {
	self.ref = ref
}

func (self *SubCommitsViewModel) GetRef() types.Ref {
	return self.ref
}

func (self *SubCommitsViewModel) SetRefToShowDivergenceFrom(ref string) {
	self.refToShowDivergenceFrom = ref
}

func (self *SubCommitsViewModel) GetRefToShowDivergenceFrom() string {
	return self.refToShowDivergenceFrom
}

func (self *SubCommitsViewModel) SetShowBranchHeads(value bool) {
	self.showBranchHeads = value
}

func (self *SubCommitsViewModel) GetShowBranchHeads() bool {
	return self.showBranchHeads
}

func (self *SubCommitsViewModel) SetShowTags(value string) {
	self.showTags = value
}

func (self *SubCommitsViewModel) GetShowTags() string {
	return self.showTags
}

func (self *SubCommitsContext) CanRebase() bool {
	return false
}

func (self *SubCommitsContext) GetSelectedRef() types.Ref {
	commit := self.GetSelected()
	if commit == nil {
		return nil
	}
	return commit
}

func (self *SubCommitsContext) GetSelectedRefRangeForDiffFiles() *types.RefRange {
	commits, startIdx, endIdx := self.GetSelectedItems()
	if commits == nil || startIdx == endIdx {
		return nil
	}
	from := commits[len(commits)-1]
	to := commits[0]
	if from.Divergence != to.Divergence {
		return nil
	}
	return &types.RefRange{From: from, To: to}
}

func (self *SubCommitsContext) GetCommits() []*models.Commit {
	return self.getModel()
}

func (self *SubCommitsContext) SetLimitCommits(value bool) {
	self.limitCommits = value
}

func (self *SubCommitsContext) GetLimitCommits() bool {
	return self.limitCommits
}

func (self *SubCommitsContext) GetDiffTerminals() []string {
	itemId := self.GetSelectedItemId()

	return []string{itemId}
}

func (self *SubCommitsContext) RefForAdjustingLineNumberInDiff() string {
	commits, _, _ := self.GetSelectedItems()
	if commits == nil {
		return ""
	}
	return commits[0].Hash()
}

func (self *SubCommitsContext) ModelSearchResults(searchStr string, caseSensitive bool) []gocui.SearchPosition {
	return searchModelCommits(caseSensitive, self.GetCommits(), self.ColumnPositions(), searchStr)
}

func (self *SubCommitsContext) IndexForGotoBottom() int {
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
