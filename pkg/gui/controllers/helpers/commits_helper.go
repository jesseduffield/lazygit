package helpers

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type ICommitsHelper interface {
	UpdateCommitPanelView(message string)
}

type CommitsHelper struct {
	c *types.HelperCommon

	model                *types.Model
	contexts             *context.ContextTree
	getCommitSummary     func() string
	setCommitSummary     func(string)
	getCommitDescription func() string
	setCommitDescription func(string)
	renderCommitLength   func()
}

var _ ICommitsHelper = &CommitsHelper{}

func NewCommitsHelper(
	c *types.HelperCommon,
	model *types.Model,
	contexts *context.ContextTree,
	getCommitSummary func() string,
	setCommitSummary func(string),
	getCommitDescription func() string,
	setCommitDescription func(string),
	renderCommitLength func(),
) *CommitsHelper {
	return &CommitsHelper{
		c:                    c,
		model:                model,
		contexts:             contexts,
		getCommitSummary:     getCommitSummary,
		setCommitSummary:     setCommitSummary,
		getCommitDescription: getCommitDescription,
		setCommitDescription: setCommitDescription,
		renderCommitLength:   renderCommitLength,
	}
}

func (self *CommitsHelper) SplitCommitMessageAndDescription(message string) (string, string) {
	for _, separator := range []string{"\n\n", "\n\r\n\r", "\n", "\n\r"} {
		msg, description, found := strings.Cut(message, separator)
		if found {
			return msg, description
		}
	}
	return message, ""
}

func (self *CommitsHelper) SetMessageAndDescriptionInView(message string) {
	summary, description := self.SplitCommitMessageAndDescription(message)

	self.setCommitSummary(summary)
	self.setCommitDescription(description)
	self.renderCommitLength()
}

func (self *CommitsHelper) JoinCommitMessageAndDescription() string {
	if len(self.getCommitDescription()) == 0 {
		return self.getCommitSummary()
	}
	return self.getCommitSummary() + "\n" + self.getCommitDescription()
}

func (self *CommitsHelper) UpdateCommitPanelView(message string) {
	// first try the passed in message, if not fallback to context -> view in that order
	if message != "" {
		self.SetMessageAndDescriptionInView(message)
		return
	}
	message = self.contexts.CommitMessage.GetPreservedMessage()
	if message != "" {
		self.SetMessageAndDescriptionInView(message)
	} else {
		self.SetMessageAndDescriptionInView(self.getCommitSummary())
	}
}

type OpenCommitMessagePanelOpts struct {
	CommitIndex     int
	Title           string
	PreserveMessage bool
	OnConfirm       func(string) error
	InitialMessage  string
}

func (self *CommitsHelper) OpenCommitMessagePanel(opts *OpenCommitMessagePanelOpts) error {
	self.contexts.CommitMessage.SetPanelState(
		opts.CommitIndex,
		opts.Title,
		opts.PreserveMessage,
		opts.OnConfirm,
	)

	self.UpdateCommitPanelView(opts.InitialMessage)

	return self.pushCommitMessageContexts()
}

func (self *CommitsHelper) OnCommitSuccess() {
	// if we have a preserved message we want to clear it on success
	if self.contexts.CommitMessage.GetPreserveMessage() {
		self.contexts.CommitMessage.SetPreservedMessage("")
	}
	self.SetMessageAndDescriptionInView("")
}

func (self *CommitsHelper) HandleCommitConfirm() error {
	fullMessage := self.JoinCommitMessageAndDescription()

	if fullMessage == "" {
		return self.c.ErrorMsg(self.c.Tr.CommitWithoutMessageErr)
	}

	err := self.contexts.CommitMessage.OnConfirm(fullMessage)
	if err != nil {
		return err
	}

	return nil
}

func (self *CommitsHelper) CloseCommitMessagePanel() error {
	if self.contexts.CommitMessage.GetPreserveMessage() {
		message := self.JoinCommitMessageAndDescription()

		self.contexts.CommitMessage.SetPreservedMessage(message)
	} else {
		self.SetMessageAndDescriptionInView("")
	}

	self.contexts.CommitMessage.SetHistoryMessage("")

	return self.PopCommitMessageContexts()
}

func (self *CommitsHelper) PopCommitMessageContexts() error {
	return self.c.RemoveContexts(self.commitMessageContexts())
}

func (self *CommitsHelper) pushCommitMessageContexts() error {
	for _, context := range self.commitMessageContexts() {
		if err := self.c.PushContext(context); err != nil {
			return err
		}
	}

	return nil
}

func (self *CommitsHelper) commitMessageContexts() []types.Context {
	return []types.Context{
		self.contexts.CommitDescription,
		self.contexts.CommitMessage,
	}
}
