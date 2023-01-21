package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type CommitMessageContext struct {
	types.Context
	viewModel *CommitMessageViewModel
}

// when selectedIndex (see below) is set to this value, it means that we're not
// currently viewing a commit message of an existing commit: instead we're making our own
// new commit message
const NoCommitIndex = -1

type CommitMessageViewModel struct {
	// index of the commit message, where -1 is 'no commit', 0 is the HEAD commit, 1
	// is the prior commit, and so on
	selectedindex int
	// if true, then upon escaping from the commit message panel, we will preserve
	// the message so that it's still shown next time we open the panel
	preserveMessage bool
	// the full preserved message (combined summary and description)
	preservedMessage string
	// invoked when pressing enter in the commit message panel
	onConfirm func(string) error
}

func NewCommitMessageContext(
	view *gocui.View,
	opts ContextCallbackOpts,
) *CommitMessageContext {
	viewModel := &CommitMessageViewModel{}
	return &CommitMessageContext{
		viewModel: viewModel,
		Context: NewSimpleContext(
			NewBaseContext(NewBaseContextOpts{
				Kind:                  types.PERSISTENT_POPUP,
				View:                  view,
				WindowName:            "commitMessage",
				Key:                   COMMIT_MESSAGE_CONTEXT_KEY,
				Focusable:             true,
				HasUncontrolledBounds: true,
			}),
			opts,
		),
	}
}

func (self *CommitMessageContext) SetSelectedIndex(value int) {
	self.viewModel.selectedindex = value
}

func (self *CommitMessageContext) GetSelectedIndex() int {
	return self.viewModel.selectedindex
}

func (self *CommitMessageContext) GetPreserveMessage() bool {
	return self.viewModel.preserveMessage
}

func (self *CommitMessageContext) OnConfirm(message string) error {
	return self.viewModel.onConfirm(message)
}

func (self *CommitMessageContext) SetPanelState(index int, title string, preserveMessage bool, onConfirm func(string) error) {
	self.viewModel.selectedindex = index
	self.viewModel.preserveMessage = preserveMessage
	self.viewModel.onConfirm = onConfirm
	self.GetView().Title = title
}

func (self *CommitMessageContext) SetPreservedMessage(message string) {
	self.viewModel.preservedMessage = message
}

func (self *CommitMessageContext) GetPreservedMessage() string {
	return self.viewModel.preservedMessage
}
