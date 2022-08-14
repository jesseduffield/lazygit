package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type RemotesContext struct {
	*BasicViewModel[*models.Remote]
	*ListContextTrait
}

var _ types.IListContext = (*RemotesContext)(nil)

func NewRemotesContext(
	getModel func() []*models.Remote,
	view *gocui.View,
	getDisplayStrings func(startIdx int, length int) [][]string,

	onFocus func(types.OnFocusOpts) error,
	onRenderToMain func() error,
	onFocusLost func(opts types.OnFocusLostOpts) error,

	c *types.HelperCommon,
) *RemotesContext {
	viewModel := NewBasicViewModel(getModel)

	return &RemotesContext{
		BasicViewModel: viewModel,
		ListContextTrait: &ListContextTrait{
			Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
				View:       view,
				WindowName: "branches",
				Key:        REMOTES_CONTEXT_KEY,
				Kind:       types.SIDE_CONTEXT,
				Focusable:  true,
			}), ContextCallbackOpts{
				OnFocus:        onFocus,
				OnFocusLost:    onFocusLost,
				OnRenderToMain: onRenderToMain,
			}),
			list:              viewModel,
			getDisplayStrings: getDisplayStrings,
			c:                 c,
		},
	}
}

func (self *RemotesContext) GetSelectedItemId() string {
	item := self.GetSelected()
	if item == nil {
		return ""
	}

	return item.ID()
}
