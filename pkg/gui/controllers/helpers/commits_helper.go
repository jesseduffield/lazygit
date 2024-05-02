package helpers

import (
	"errors"
	"path/filepath"
	"strings"
	"time"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

type ICommitsHelper interface {
	UpdateCommitPanelView(message string)
}

type CommitsHelper struct {
	c *HelperCommon

	getCommitSummary              func() string
	setCommitSummary              func(string)
	getCommitDescription          func() string
	getUnwrappedCommitDescription func() string
	setCommitDescription          func(string)
}

var _ ICommitsHelper = &CommitsHelper{}

func NewCommitsHelper(
	c *HelperCommon,
	getCommitSummary func() string,
	setCommitSummary func(string),
	getCommitDescription func() string,
	getUnwrappedCommitDescription func() string,
	setCommitDescription func(string),
) *CommitsHelper {
	return &CommitsHelper{
		c:                             c,
		getCommitSummary:              getCommitSummary,
		setCommitSummary:              setCommitSummary,
		getCommitDescription:          getCommitDescription,
		getUnwrappedCommitDescription: getUnwrappedCommitDescription,
		setCommitDescription:          setCommitDescription,
	}
}

func (self *CommitsHelper) SplitCommitMessageAndDescription(message string) (string, string) {
	msg, description, _ := strings.Cut(message, "\n")
	return msg, strings.TrimSpace(description)
}

func (self *CommitsHelper) SetMessageAndDescriptionInView(message string) {
	summary, description := self.SplitCommitMessageAndDescription(message)

	self.setCommitSummary(summary)
	self.setCommitDescription(description)
	self.c.Contexts().CommitMessage.RenderCommitLength()
}

func (self *CommitsHelper) JoinCommitMessageAndUnwrappedDescription() string {
	if len(self.getUnwrappedCommitDescription()) == 0 {
		return self.getCommitSummary()
	}
	return self.getCommitSummary() + "\n" + self.getUnwrappedCommitDescription()
}

func TryRemoveHardLineBreaks(message string, autoWrapWidth int) string {
	messageRunes := []rune(message)
	lastHardLineStart := 0
	for i, r := range messageRunes {
		if r == '\n' {
			// Try to make this a soft linebreak by turning it into a space, and
			// checking whether it still wraps to the same result then.
			messageRunes[i] = ' '

			_, cursorMapping := gocui.AutoWrapContent(messageRunes[lastHardLineStart:], autoWrapWidth)

			// Look at the cursorMapping to check whether auto-wrapping inserted
			// a line break. If it did, there will be a cursorMapping entry with
			// Orig pointing to the position after the inserted line break.
			if len(cursorMapping) == 0 || cursorMapping[0].Orig != i-lastHardLineStart+1 {
				// It didn't, so change it back to a newline
				messageRunes[i] = '\n'
			}
			lastHardLineStart = i + 1
		}
	}

	return string(messageRunes)
}

func (self *CommitsHelper) SwitchToEditor() error {
	message := lo.Ternary(len(self.getCommitDescription()) == 0,
		self.getCommitSummary(),
		self.getCommitSummary()+"\n\n"+self.getCommitDescription())
	filepath := filepath.Join(self.c.OS().GetTempDir(), self.c.Git().RepoPaths.RepoName(), time.Now().Format("Jan _2 15.04.05.000000000")+".msg")
	err := self.c.OS().CreateFileWithContent(filepath, message)
	if err != nil {
		return err
	}

	err = self.CloseCommitMessagePanel()
	if err != nil {
		return err
	}

	return self.c.Contexts().CommitMessage.SwitchToEditor(filepath)
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
	OnSwitchToEditor func(string) error
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
		opts.OnSwitchToEditor,
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
		return errors.New(self.c.Tr.CommitWithoutMessageErr)
	}

	err := self.c.Contexts().CommitMessage.OnConfirm(summary, description)
	if err != nil {
		return err
	}

	return nil
}

func (self *CommitsHelper) CloseCommitMessagePanel() error {
	if self.c.Contexts().CommitMessage.GetPreserveMessage() {
		message := self.JoinCommitMessageAndUnwrappedDescription()

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

func (self *CommitsHelper) OpenCommitMenu(suggestionFunc func(string) []*types.Suggestion) error {
	var disabledReasonForOpenInEditor *types.DisabledReason
	if !self.c.Contexts().CommitMessage.CanSwitchToEditor() {
		disabledReasonForOpenInEditor = &types.DisabledReason{
			Text: self.c.Tr.CommandDoesNotSupportOpeningInEditor,
		}
	}

	menuItems := []*types.MenuItem{
		{
			Label: self.c.Tr.OpenInEditor,
			OnPress: func() error {
				return self.SwitchToEditor()
			},
			Key:            'e',
			DisabledReason: disabledReasonForOpenInEditor,
		},
		{
			Label: self.c.Tr.AddCoAuthor,
			OnPress: func() error {
				return self.addCoAuthor(suggestionFunc)
			},
			Key: 'c',
		},
	}
	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.CommitMenuTitle,
		Items: menuItems,
	})
}

func (self *CommitsHelper) addCoAuthor(suggestionFunc func(string) []*types.Suggestion) error {
	return self.c.Prompt(types.PromptOpts{
		Title:               self.c.Tr.AddCoAuthorPromptTitle,
		FindSuggestionsFunc: suggestionFunc,
		HandleConfirm: func(value string) error {
			commitDescription := self.getCommitDescription()
			commitDescription = git_commands.AddCoAuthorToDescription(commitDescription, value)
			self.setCommitDescription(commitDescription)
			return nil
		},
	})
}
