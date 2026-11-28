package context

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type BranchesContext struct {
	*FilteredListViewModel[*models.Branch]
	*ListContextTrait
}

var (
	_ types.IListContext    = (*BranchesContext)(nil)
	_ types.DiffableContext = (*BranchesContext)(nil)
)

func NewBranchesContext(c *ContextCommon) *BranchesContext {
	viewModel := NewFilteredListViewModel(
		func() []*models.Branch { return c.Model().Branches },
		func(branch *models.Branch) []string {
			return []string{branch.Name}
		},
	)

	getDisplayStrings := func(_ int, _ int) [][]string {
		branches := viewModel.GetItems()
		prs := c.Model().PullRequestsMap
		fullDescription := c.State().GetRepoState().GetScreenMode() != types.SCREEN_NORMAL

		// The rows are as wide as the view plus its horizontal scroll position, so that
		// scrolling right reveals more of a long branch name. Once every name fits,
		// rendering wider would only move the right-aligned divergence from the base
		// branch further right, and scrolling right could go on forever.
		view := c.Views().Branches
		width := view.InnerWidth()
		if view.OriginX() > 0 {
			untruncatedWidth := presentation.GetBranchListUntruncatedWidth(
				branches, c.State().GetItemOperation, prs, fullDescription, c.Tr, c.UserConfig(), c.Model().Worktrees)
			width += min(view.OriginX(), max(0, untruncatedWidth-width))
		}

		return presentation.GetBranchListDisplayStrings(
			branches,
			c.State().GetItemOperation,
			prs,
			fullDescription,
			c.Modes().Diffing.Ref,
			width,
			c.Tr,
			c.UserConfig(),
			c.Model().Worktrees,
		)
	}

	self := &BranchesContext{
		FilteredListViewModel: viewModel,
		ListContextTrait: &ListContextTrait{
			Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
				View:                       c.Views().Branches,
				WindowName:                 "branches",
				Key:                        LOCAL_BRANCHES_CONTEXT_KEY,
				Kind:                       types.SIDE_CONTEXT,
				Focusable:                  true,
				NeedsRerenderOnWidthChange: types.NEEDS_RERENDER_ON_WIDTH_CHANGE_WHEN_WIDTH_CHANGES,
			})),
			ListRenderer: ListRenderer{
				list:              viewModel,
				getDisplayStrings: getDisplayStrings,
			},
			c: c,
		},
	}

	return self
}

func (self *BranchesContext) GetSelectedRef() models.Ref {
	branch := self.GetSelected()
	if branch == nil {
		return nil
	}
	return branch
}

func (self *BranchesContext) GetDiffTerminals() []string {
	// for our local branches we want to include both the branch and its upstream
	branch := self.GetSelected()
	if branch != nil {
		names := []string{branch.ID()}
		if branch.IsTrackingRemote() {
			names = append(names, branch.ID()+"@{u}")
		}
		return names
	}
	return nil
}

func (self *BranchesContext) RefForAdjustingLineNumberInDiff() string {
	branch := self.GetSelected()
	if branch != nil {
		return branch.ID()
	}
	return ""
}

func (self *BranchesContext) ShowBranchHeadsInSubCommits() bool {
	return true
}
