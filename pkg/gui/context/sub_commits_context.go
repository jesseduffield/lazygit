package context

import (
	"fmt"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type SubCommitsContext struct {
	*SubCommitsViewModel
	*ViewportListContextTrait
	*DynamicTitleBuilder
}

var _ types.IListContext = (*SubCommitsContext)(nil)

func NewSubCommitsContext(
	getModel func() []*models.Commit,
	view *gocui.View,
	getDisplayStrings func(startIdx int, length int) [][]string,

	onFocus func(types.OnFocusOpts) error,
	onRenderToMain func() error,
	onFocusLost func(opts types.OnFocusLostOpts) error,

	c *types.HelperCommon,
) *SubCommitsContext {
	viewModel := &SubCommitsViewModel{
		BasicViewModel: NewBasicViewModel(getModel),
		ref:            nil,
		limitCommits:   true,
	}

	return &SubCommitsContext{
		SubCommitsViewModel: viewModel,
		DynamicTitleBuilder: NewDynamicTitleBuilder(c.Tr.SubCommitsDynamicTitle),
		ViewportListContextTrait: &ViewportListContextTrait{
			ListContextTrait: &ListContextTrait{
				Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
					View:       view,
					WindowName: "branches",
					Key:        SUB_COMMITS_CONTEXT_KEY,
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
		},
	}
}

type SubCommitsViewModel struct {
	// name of the ref that the sub-commits are shown for
	ref types.Ref
	*BasicViewModel[*models.Commit]

	limitCommits bool
}

func (self *SubCommitsViewModel) SetRef(ref types.Ref) {
	self.ref = ref
}

func (self *SubCommitsViewModel) GetRef() types.Ref {
	return self.ref
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

func (self *SubCommitsContext) GetSelectedRef() types.Ref {
	commit := self.GetSelected()
	if commit == nil {
		return nil
	}
	return commit
}

func (self *SubCommitsContext) GetCommits() []*models.Commit {
	return self.getModel()
}

func (self *SubCommitsContext) Title() string {
	return fmt.Sprintf(self.c.Tr.SubCommitsDynamicTitle, utils.TruncateWithEllipsis(self.ref.RefName(), 50))
}

func (self *SubCommitsContext) SetLimitCommits(value bool) {
	self.limitCommits = value
}

func (self *SubCommitsContext) GetLimitCommits() bool {
	return self.limitCommits
}
