package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context/traits"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type ReflogCommitsContext struct {
	*ReflogCommitsViewModel
	*ListContextTrait
}

var _ types.IListContext = (*ReflogCommitsContext)(nil)

func NewReflogCommitsContext(
	getModel func() []*models.Commit,
	view *gocui.View,
	getDisplayStrings func(startIdx int, length int) [][]string,

	onFocus func(...types.OnFocusOpts) error,
	onRenderToMain func(...types.OnFocusOpts) error,
	onFocusLost func() error,

	c *types.HelperCommon,
) *ReflogCommitsContext {
	viewModel := NewReflogCommitsViewModel(getModel)

	return &ReflogCommitsContext{
		ReflogCommitsViewModel: viewModel,
		ListContextTrait: &ListContextTrait{
			Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
				ViewName:   "commits",
				WindowName: "commits",
				Key:        REFLOG_COMMITS_CONTEXT_KEY,
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

func (self *ReflogCommitsContext) GetSelectedRefName() string {
	item := self.GetSelected()

	if item == nil {
		return ""
	}

	return item.RefName()
}

type ReflogCommitsViewModel struct {
	*traits.ListCursor
	getModel func() []*models.Commit
}

func NewReflogCommitsViewModel(getModel func() []*models.Commit) *ReflogCommitsViewModel {
	self := &ReflogCommitsViewModel{
		getModel: getModel,
	}

	self.ListCursor = traits.NewListCursor(self)

	return self
}

func (self *ReflogCommitsViewModel) GetItemsLength() int {
	return len(self.getModel())
}

func (self *ReflogCommitsViewModel) GetSelected() *models.Commit {
	if self.GetItemsLength() == 0 {
		return nil
	}

	return self.getModel()[self.GetSelectedLineIdx()]
}
