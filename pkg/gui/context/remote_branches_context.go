package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type RemoteBranchesContext struct {
	*BasicViewModel[*models.RemoteBranch]
	*ListContextTrait
	*DynamicTitleBuilder
}

var _ types.IListContext = (*RemoteBranchesContext)(nil)

func NewRemoteBranchesContext(
	getModel func() []*models.RemoteBranch,
	view *gocui.View,
	getDisplayStrings func(startIdx int, length int) [][]string,

	onFocus func(types.OnFocusOpts) error,
	onRenderToMain func() error,
	onFocusLost func(opts types.OnFocusLostOpts) error,

	c *types.HelperCommon,
) *RemoteBranchesContext {
	viewModel := NewBasicViewModel(getModel)

	return &RemoteBranchesContext{
		BasicViewModel:      viewModel,
		DynamicTitleBuilder: NewDynamicTitleBuilder(c.Tr.RemoteBranchesDynamicTitle),
		ListContextTrait: &ListContextTrait{
			Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
				View:       view,
				WindowName: "branches",
				Key:        REMOTE_BRANCHES_CONTEXT_KEY,
				Kind:       types.SIDE_CONTEXT,
				Focusable:  true,
				Transient:  true,
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

func (self *RemoteBranchesContext) GetSelectedItemId() string {
	item := self.GetSelected()
	if item == nil {
		return ""
	}

	return item.ID()
}

func (self *RemoteBranchesContext) GetSelectedRef() types.Ref {
	remoteBranch := self.GetSelected()
	if remoteBranch == nil {
		return nil
	}
	return remoteBranch
}
