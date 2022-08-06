package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/mergeconflicts"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type MergeConflictsContext struct {
	types.Context
	viewModel *ConflictsViewModel
	c         *types.HelperCommon
}

type ConflictsViewModel struct {
	state *mergeconflicts.State

	// userVerticalScrolling tells us if the user has started scrolling through the file themselves
	// in which case we won't auto-scroll to a conflict.
	userVerticalScrolling bool
}

func NewMergeConflictsContext(
	view *gocui.View,

	opts ContextCallbackOpts,

	c *types.HelperCommon,
	getOptionsMap func() map[string]string,
) *MergeConflictsContext {
	viewModel := &ConflictsViewModel{
		state:                 mergeconflicts.NewState(),
		userVerticalScrolling: false,
	}

	return &MergeConflictsContext{
		viewModel: viewModel,
		Context: NewSimpleContext(
			NewBaseContext(NewBaseContextOpts{
				Kind:            types.MAIN_CONTEXT,
				View:            view,
				WindowName:      "main",
				Key:             MERGE_CONFLICTS_CONTEXT_KEY,
				OnGetOptionsMap: getOptionsMap,
				Focusable:       true,
			}),
			opts,
		),
		c: c,
	}
}

func (self *MergeConflictsContext) SetUserScrolling(isScrolling bool) {
	self.viewModel.userVerticalScrolling = isScrolling
}

func (self *MergeConflictsContext) IsUserScrolling() bool {
	return self.viewModel.userVerticalScrolling
}

func (self *MergeConflictsContext) State() *mergeconflicts.State {
	return self.viewModel.state
}
