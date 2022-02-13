package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context/traits"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type SubCommitsContext struct {
	*SubCommitsViewModel
	*ViewportListContextTrait
}

var _ types.IListContext = (*SubCommitsContext)(nil)

func NewSubCommitsContext(
	getModel func() []*models.Commit,
	view *gocui.View,
	getDisplayStrings func(startIdx int, length int) [][]string,

	onFocus func(...types.OnFocusOpts) error,
	onRenderToMain func(...types.OnFocusOpts) error,
	onFocusLost func() error,

	c *types.HelperCommon,
) *SubCommitsContext {
	viewModel := NewSubCommitsViewModel(getModel)

	return &SubCommitsContext{
		SubCommitsViewModel: viewModel,
		ViewportListContextTrait: &ViewportListContextTrait{
			ListContextTrait: &ListContextTrait{
				Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
					ViewName:   "branches",
					WindowName: "branches",
					Key:        SUB_COMMITS_CONTEXT_KEY,
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

func (self *SubCommitsContext) GetSelectedItemId() string {
	item := self.GetSelected()
	if item == nil {
		return ""
	}

	return item.ID()
}

func (self *SubCommitsContext) CanRebase() bool {
	return false
}

type SubCommitsViewModel struct {
	*traits.ListCursor
	getModel func() []*models.Commit
}

func NewSubCommitsViewModel(getModel func() []*models.Commit) *SubCommitsViewModel {
	self := &SubCommitsViewModel{
		getModel: getModel,
	}

	self.ListCursor = traits.NewListCursor(self)

	return self
}

func (self *SubCommitsViewModel) GetItemsLength() int {
	return len(self.getModel())
}

func (self *SubCommitsViewModel) GetSelected() *models.Commit {
	if self.GetItemsLength() == 0 {
		return nil
	}

	return self.getModel()[self.GetSelectedLineIdx()]
}
