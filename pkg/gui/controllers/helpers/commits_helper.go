package helpers

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type ICommitsHelper interface {
	UpdateCommitPanelView(message string)
}

type CommitsHelper struct {
	c *HelperCommon

	getCommitSummary     func() string
	setCommitSummary     func(string)
	getCommitDescription func() string
	setCommitDescription func(string)
}

var _ ICommitsHelper = &CommitsHelper{}

func NewCommitsHelper(
	c *HelperCommon,
	getCommitSummary func() string,
	setCommitSummary func(string),
	getCommitDescription func() string,
	setCommitDescription func(string),
) *CommitsHelper {
	return &CommitsHelper{
		c:                    c,
		getCommitSummary:     getCommitSummary,
		setCommitSummary:     setCommitSummary,
		getCommitDescription: getCommitDescription,
		setCommitDescription: setCommitDescription,
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
	self.c.Contexts().CommitMessage.RenderCommitLength()
}

func (self *CommitsHelper) JoinCommitMessageAndDescription() string {
	if len(self.getCommitDescription()) == 0 {
		return self.getCommitSummary()
	}
	return self.getCommitSummary() + "\n" + self.getCommitDescription()
}

func (self *CommitsHelper) UpdateCommitPanelView(message string) {
	if message != "" {
		self.SetMessageAndDescriptionInView(message)
		return
	}

	if self.c.Contexts().CommitMessage.GetPreserveMessage() {
		preservedMessage := self.c.Contexts().CommitMessage.GetPreservedMessage()
		self.SetMessageAndDescriptionInView(preservedMessage)
		return
	}

	self.SetMessageAndDescriptionInView("")
}

type OpenCommitMessagePanelOpts struct {
	CommitIndex      int
	SummaryTitle     string
	DescriptionTitle string
	PreserveMessage  bool
	OnConfirm        func(summary string, description string) error
	InitialMessage   string
}

func (self *CommitsHelper) OpenCommitMessagePanel(opts *OpenCommitMessagePanelOpts) error {
	onConfirm := func(summary string, description string) error {
		if err := self.CloseCommitMessagePanel(); err != nil {
			return err
		}

		return opts.OnConfirm(summary, description)
	}

	self.c.Contexts().CommitMessage.SetPanelState(
		opts.CommitIndex,
		opts.SummaryTitle,
		opts.DescriptionTitle,
		opts.PreserveMessage,
		onConfirm,
	)

	self.UpdateCommitPanelView(opts.InitialMessage)

	return self.pushCommitMessageContexts()
}

func (self *CommitsHelper) OnCommitSuccess() {
	// if we have a preserved message we want to clear it on success
	if self.c.Contexts().CommitMessage.GetPreserveMessage() {
		self.c.Contexts().CommitMessage.SetPreservedMessage("")
	}
}

func (self *CommitsHelper) HandleCommitConfirm() error {
	summary, description := self.getCommitSummary(), self.getCommitDescription()

	if summary == "" {
		return self.c.ErrorMsg(self.c.Tr.CommitWithoutMessageErr)
	}

	err := self.c.Contexts().CommitMessage.OnConfirm(summary, description)
	if err != nil {
		return err
	}

	return nil
}

func (self *CommitsHelper) CloseCommitMessagePanel() error {
	if self.c.Contexts().CommitMessage.GetPreserveMessage() {
		message := self.JoinCommitMessageAndDescription()

		self.c.Contexts().CommitMessage.SetPreservedMessage(message)
	} else {
		self.SetMessageAndDescriptionInView("")
	}

	self.c.Contexts().CommitMessage.SetHistoryMessage("")

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
		self.c.Contexts().CommitDescription,
		self.c.Contexts().CommitMessage,
	}
}
