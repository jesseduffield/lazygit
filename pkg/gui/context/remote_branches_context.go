package context

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

type RemoteBranchesContext struct {
	*FilteredListViewModel[*models.RemoteBranch]
	*ListContextTrait
	*DynamicTitleBuilder
}

var (
	_ types.IListContext    = (*RemoteBranchesContext)(nil)
	_ types.DiffableContext = (*RemoteBranchesContext)(nil)
)

func NewRemoteBranchesContext(
	c *ContextCommon,
) *RemoteBranchesContext {
	viewModel := NewFilteredListViewModel(
		func() []*models.RemoteBranch { return c.Model().RemoteBranches },
		func(remoteBranch *models.RemoteBranch) []string {
			return []string{remoteBranch.Name}
		},
	)

	getDisplayStrings := func(_ int, _ int) [][]string {
		return presentation.GetRemoteBranchListDisplayStrings(viewModel.GetItems(), c.Modes().Diffing.Ref)
	}

	return &RemoteBranchesContext{
		FilteredListViewModel: viewModel,
		DynamicTitleBuilder:   NewDynamicTitleBuilder(c.Tr.RemoteBranchesDynamicTitle),
		ListContextTrait: &ListContextTrait{
			Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
				View:                       c.Views().RemoteBranches,
				WindowName:                 "branches",
				Key:                        REMOTE_BRANCHES_CONTEXT_KEY,
				Kind:                       types.SIDE_CONTEXT,
				Focusable:                  true,
				Transient:                  true,
				NeedsRerenderOnWidthChange: true,
			})),
			ListRenderer: ListRenderer{
				list:              viewModel,
				getDisplayStrings: getDisplayStrings,
			},
			c: c,
		},
	}
}

func (self *RemoteBranchesContext) GetSelectedRef() types.Ref {
	remoteBranch := self.GetSelected()
	if remoteBranch == nil {
		return nil
	}
	return remoteBranch
}

func (self *RemoteBranchesContext) GetSelectedRefs() ([]types.Ref, int, int) {
	items, startIdx, endIdx := self.GetSelectedItems()

	refs := lo.Map(items, func(item *models.RemoteBranch, _ int) types.Ref {
		return item
	})

	return refs, startIdx, endIdx
}

func (self *RemoteBranchesContext) GetDiffTerminals() []string {
	itemId := self.GetSelectedItemId()

	return []string{itemId}
}

func (self *RemoteBranchesContext) ShowBranchHeadsInSubCommits() bool {
	return true
}
