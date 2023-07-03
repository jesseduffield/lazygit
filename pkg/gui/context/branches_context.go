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

	getDisplayStrings := func(startIdx int, length int) [][]string {
		return presentation.GetBranchListDisplayStrings(
			viewModel.GetItems(),
			c.State().GetRepoState().GetScreenMode() != types.SCREEN_NORMAL,
			c.Modes().Diffing.Ref,
			c.Tr,
			c.UserConfig,
		)
	}

	self := &BranchesContext{
		FilteredListViewModel: viewModel,
		ListContextTrait: &ListContextTrait{
			Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
				View:       c.Views().Branches,
				WindowName: "branches",
				Key:        LOCAL_BRANCHES_CONTEXT_KEY,
				Kind:       types.SIDE_CONTEXT,
				Focusable:  true,
			})),
			list:              viewModel,
			getDisplayStrings: getDisplayStrings,
			c:                 c,
		},
	}

	return self
}

func (self *BranchesContext) GetSelectedItemId() string {
	item := self.GetSelected()
	if item == nil {
		return ""
	}

	return item.ID()
}

func (self *BranchesContext) GetSelectedRef() types.Ref {
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
