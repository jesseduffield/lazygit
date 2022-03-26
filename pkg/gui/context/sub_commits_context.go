package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type SubCommitsContext struct {
	*BasicViewModel[*models.Commit]
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
	viewModel := NewBasicViewModel(getModel)

	return &SubCommitsContext{
		BasicViewModel: viewModel,
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
			},
		},
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

func (self *SubCommitsContext) GetSelectedRefName() string {
	item := self.GetSelected()

	if item == nil {
		return ""
	}

	return item.RefName()
}

func (self *SubCommitsContext) GetCommits() []*models.Commit {
	return self.getModel()
}
