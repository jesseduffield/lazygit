package context

import (
	"strconv"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type CommitMessageContext struct {
	c *ContextCommon
	types.Context
	viewModel *CommitMessageViewModel
}

var _ types.Context = (*CommitMessageContext)(nil)

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
	onConfirm func(string, string) error

	// The message typed in before cycling through history
	// We store this separately to 'preservedMessage' because 'preservedMessage'
	// is specifically for committing staged files and we don't want this affected
	// by cycling through history in the context of rewording an old commit.
	historyMessage string
}

func NewCommitMessageContext(
	c *ContextCommon,
) *CommitMessageContext {
	viewModel := &CommitMessageViewModel{}
	return &CommitMessageContext{
		c:         c,
		viewModel: viewModel,
		Context: NewSimpleContext(
			NewBaseContext(NewBaseContextOpts{
				Kind:                  types.PERSISTENT_POPUP,
				View:                  c.Views().CommitMessage,
				WindowName:            "commitMessage",
				Key:                   COMMIT_MESSAGE_CONTEXT_KEY,
				Focusable:             true,
				HasUncontrolledBounds: true,
			}),
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

func (self *CommitMessageContext) GetPreservedMessage() string {
	return self.viewModel.preservedMessage
}

func (self *CommitMessageContext) SetPreservedMessage(message string) {
	self.viewModel.preservedMessage = message
}

func (self *CommitMessageContext) GetHistoryMessage() string {
	return self.viewModel.historyMessage
}

func (self *CommitMessageContext) SetHistoryMessage(message string) {
	self.viewModel.historyMessage = message
}

func (self *CommitMessageContext) OnConfirm(summary string, description string) error {
	return self.viewModel.onConfirm(summary, description)
}

func (self *CommitMessageContext) SetPanelState(
	index int,
	summaryTitle string,
	descriptionTitle string,
	preserveMessage bool,
	onConfirm func(string, string) error,
) {
	self.viewModel.selectedindex = index
	self.viewModel.preserveMessage = preserveMessage
	self.viewModel.onConfirm = onConfirm
	self.GetView().Title = summaryTitle
	self.c.Views().CommitDescription.Title = descriptionTitle
}

func (self *CommitMessageContext) RenderCommitLength() {
	if !self.c.UserConfig.Gui.CommitLength.Show {
		return
	}

	self.c.Views().CommitMessage.Subtitle = getBufferLength(self.c.Views().CommitMessage)
}

func getBufferLength(view *gocui.View) string {
	return " " + strconv.Itoa(strings.Count(view.TextArea.GetContent(), "")-1) + " "
}
