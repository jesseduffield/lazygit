package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context/traits"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type LocalCommitsContext struct {
	*LocalCommitsViewModel
	*ViewportListContextTrait
}

var _ types.IListContext = (*LocalCommitsContext)(nil)

func NewLocalCommitsContext(
	getModel func() []*models.Commit,
	view *gocui.View,
	getDisplayStrings func(startIdx int, length int) [][]string,

	onFocus func(...types.OnFocusOpts) error,
	onRenderToMain func(...types.OnFocusOpts) error,
	onFocusLost func() error,

	c *types.ControllerCommon,
) *LocalCommitsContext {
	viewModel := NewLocalCommitsViewModel(getModel)

	return &LocalCommitsContext{
		LocalCommitsViewModel: viewModel,
		ViewportListContextTrait: &ViewportListContextTrait{
			ListContextTrait: &ListContextTrait{
				Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
					ViewName:   "commits",
					WindowName: "commits",
					Key:        BRANCH_COMMITS_CONTEXT_KEY,
					Kind:       types.SIDE_CONTEXT,
					Focusable:  true,
				}), ContextCallbackOpts{
					OnFocus:        onFocus,
					OnFocusLost:    onFocusLost,
					OnRenderToMain: onRenderToMain,
				}),
				list:              viewModel,
				viewTrait:         NewViewTrait(view),
				getDisplayStrings: getDisplayStrings,
				c:                 c,
			}},
	}
}

func (self *LocalCommitsContext) GetSelectedItemId() string {
	item := self.GetSelected()
	if item == nil {
		return ""
	}

	return item.ID()
}

type LocalCommitsViewModel struct {
	*traits.ListCursor
	limitCommits bool
	getModel     func() []*models.Commit
}

func NewLocalCommitsViewModel(getModel func() []*models.Commit) *LocalCommitsViewModel {
	self := &LocalCommitsViewModel{
		getModel:     getModel,
		limitCommits: true,
	}

	self.ListCursor = traits.NewListCursor(self)

	return self
}

func (self *LocalCommitsViewModel) GetItemsLength() int {
	return len(self.getModel())
}

func (self *LocalCommitsViewModel) GetSelected() *models.Commit {
	if self.GetItemsLength() == 0 {
		return nil
	}

	return self.getModel()[self.GetSelectedLineIdx()]
}

func (self *LocalCommitsViewModel) SetLimitCommits(value bool) {
	self.limitCommits = value
}

func (self *LocalCommitsViewModel) GetLimitCommits() bool {
	return self.limitCommits
}
