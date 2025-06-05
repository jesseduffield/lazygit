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

type CommitsHelper struct {
	c *HelperCommon

	getCommitSummary              func() string
	setCommitSummary              func(string)
	getCommitDescription          func() string
	getUnwrappedCommitDescription func() string
	setCommitDescription          func(string)
}

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
	self.c.Contexts().CommitMessage.RenderSubtitle()
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

	self.CloseCommitMessagePanel()

	return self.c.Contexts().CommitMessage.SwitchToEditor(filepath)
}

func (self *CommitsHelper) UpdateCommitPanelView(message string) {
	if message != "" {
		self.SetMessageAndDescriptionInView(message)
		return
	}

	if self.c.Contexts().CommitMessage.GetPreserveMessage() {
		preservedMessage := self.c.Contexts().CommitMessage.GetPreservedMessageAndLogError()
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

	// The following two fields are only for the display of the "(hooks
	// disabled)" display in the commit message panel. They have no effect on
	// the actual behavior; make sure what you are passing in matches that.
	// Leave unassigned if the concept of skipping hooks doesn't make sense for
	// what you are doing, e.g. when creating a tag.
	ForceSkipHooks  bool
	SkipHooksPrefix string
}

func (self *CommitsHelper) OpenCommitMessagePanel(opts *OpenCommitMessagePanelOpts) {
	onConfirm := func(summary string, description string) error {
		self.CloseCommitMessagePanel()

		return opts.OnConfirm(summary, description)
	}

	self.c.Contexts().CommitMessage.SetPanelState(
		opts.CommitIndex,
		opts.SummaryTitle,
		opts.DescriptionTitle,
		opts.PreserveMessage,
		opts.InitialMessage,
		onConfirm,
		opts.OnSwitchToEditor,
		opts.ForceSkipHooks,
		opts.SkipHooksPrefix,
	)

	self.UpdateCommitPanelView(opts.InitialMessage)

	self.c.Context().Push(self.c.Contexts().CommitMessage, types.OnFocusOpts{})
}

func (self *CommitsHelper) ClearPreservedCommitMessage() {
	self.c.Contexts().CommitMessage.SetPreservedMessageAndLogError("")
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

func (self *CommitsHelper) CloseCommitMessagePanel() {
	if self.c.Contexts().CommitMessage.GetPreserveMessage() {
		message := self.JoinCommitMessageAndUnwrappedDescription()
		if message != self.c.Contexts().CommitMessage.GetInitialMessage() {
			self.c.Contexts().CommitMessage.SetPreservedMessageAndLogError(message)
		}
	} else {
		self.SetMessageAndDescriptionInView("")
	}

	self.c.Contexts().CommitMessage.SetHistoryMessage("")

	self.c.Views().CommitMessage.Visible = false
	self.c.Views().CommitDescription.Visible = false

	self.c.Context().Pop()
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
		{
			Label: self.c.Tr.PasteCommitMessageFromClipboard,
			OnPress: func() error {
				return self.pasteCommitMessageFromClipboard()
			},
			Key: 'p',
		},
	}
	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.CommitMenuTitle,
		Items: menuItems,
	})
}

func (self *CommitsHelper) addCoAuthor(suggestionFunc func(string) []*types.Suggestion) error {
	self.c.Prompt(types.PromptOpts{
		Title:               self.c.Tr.AddCoAuthorPromptTitle,
		FindSuggestionsFunc: suggestionFunc,
		HandleConfirm: func(value string) error {
			commitDescription := self.getCommitDescription()
			commitDescription = git_commands.AddCoAuthorToDescription(commitDescription, value)
			self.setCommitDescription(commitDescription)
			return nil
		},
	})

	return nil
}

func (self *CommitsHelper) pasteCommitMessageFromClipboard() error {
	message, err := self.c.OS().PasteFromClipboard()
	if err != nil {
		return err
	}
	if message == "" {
		return nil
	}

	if currentMessage := self.JoinCommitMessageAndUnwrappedDescription(); currentMessage == "" {
		self.SetMessageAndDescriptionInView(message)
		return nil
	}

	// Confirm before overwriting the commit message
	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.PasteCommitMessageFromClipboard,
		Prompt: self.c.Tr.SurePasteCommitMessage,
		HandleConfirm: func() error {
			self.SetMessageAndDescriptionInView(message)
			return nil
		},
	})

	return nil
}
