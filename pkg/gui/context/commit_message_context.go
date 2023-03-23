package context

import (
	"strconv"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type CommitMessageContext struct {
	*SimpleContext
	c *ContextCommon
}

var _ types.Context = (*CommitMessageContext)(nil)

func NewCommitMessageContext(
	c *ContextCommon,
) *CommitMessageContext {
	return &CommitMessageContext{
		c: c,
		SimpleContext: NewSimpleContext(
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

func (self *CommitMessageContext) RenderCommitLength() {
	if !self.c.UserConfig.Gui.CommitLength.Show {
		return
	}

	self.c.Views().CommitMessage.Subtitle = getBufferLength(self.c.Views().CommitMessage)
}

func getBufferLength(view *gocui.View) string {
	return " " + strconv.Itoa(strings.Count(view.TextArea.GetContent(), "")-1) + " "
}
