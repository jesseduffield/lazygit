package helpers

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/gocui"
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

// SplitCommitMessageAndDescription splits a message in git's canonical format
// (summary and body separated by a blank line) into summary and description.
func (self *CommitsHelper) SplitCommitMessageAndDescription(message string) (string, string) {
	summary, description, _ := strings.Cut(message, "\n")
	description = strings.TrimPrefix(description, "\n")
	return summary, description
}

// SplitPreservedCommitMessage splits a message in our preservation format
// (summary and description joined by a single "\n") into summary and description.
// It is lossless: round-tripping through JoinCommitMessageAndUnwrappedDescription
// preserves the exact content.
func (self *CommitsHelper) SplitPreservedCommitMessage(message string) (string, string) {
	summary, description, _ := strings.Cut(message, "\n")
	return summary, description
}

func (self *CommitsHelper) SetMessageAndDescriptionInView(message string) {
	summary, description := self.SplitCommitMessageAndDescription(message)
	self.setSummaryAndDescriptionInView(summary, description)
}

func (self *CommitsHelper) SetPreservedMessageInView(message string) {
	summary, description := self.SplitPreservedCommitMessage(message)
	self.setSummaryAndDescriptionInView(summary, description)
}

func (self *CommitsHelper) setSummaryAndDescriptionInView(summary, description string) {
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
	lastHardLineStart := 0
	result := message
	for i, b := range message {
		if b == '\n' {
			// Try to make this a soft linebreak by turning it into a space, and
			// checking whether it still wraps to the same result then.
			str := message[lastHardLineStart:i] + " " + message[i+1:]
			softLineBreakIndices := gocui.AutoWrapContent(str, autoWrapWidth)

			// See if auto-wrapping inserted a soft line break:
			if len(softLineBreakIndices) > 0 && softLineBreakIndices[0] == i-lastHardLineStart+1 {
				// It did, so change it to a space in the result.
				result = result[:i] + " " + result[i+1:]
			}
			lastHardLineStart = i + 1
		}
	}

	return result
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

	// When there's no explicit initial message but we're in a preservation
	// context, fall back to any previously preserved message. This is stored as
	// the "initial" value so the unchanged-message check on close still works
	// correctly (in particular, clearing the panel then escaping will notice
	// the difference and delete the preserved file).
	initialMessage := opts.InitialMessage
	initialMessageIsPreserved := false
	if opts.PreserveMessage && initialMessage == "" {
		initialMessage = self.c.Contexts().CommitMessage.GetPreservedMessageAndLogError()
		initialMessageIsPreserved = true
	}

	self.c.Contexts().CommitMessage.SetPanelState(
		opts.CommitIndex,
		opts.SummaryTitle,
		opts.DescriptionTitle,
		opts.PreserveMessage,
		initialMessage,
		onConfirm,
		opts.OnSwitchToEditor,
		opts.ForceSkipHooks,
		opts.SkipHooksPrefix,
	)

	if initialMessageIsPreserved {
		self.SetPreservedMessageInView(initialMessage)
	} else {
		self.SetMessageAndDescriptionInView(initialMessage)
	}

	self.c.Context().Push(self.c.Contexts().CommitMessage, types.OnFocusOpts{})
}

func (self *CommitsHelper) ClearPreservedCommitMessage() {
	self.c.Contexts().CommitMessage.SetPreservedMessageAndLogError("")
}

func (self *CommitsHelper) HandleCommitConfirm() error {
	summary, description := self.getCommitSummary(), self.getCommitDescription()

	if strings.TrimSpace(summary) == "" {
		return errors.New(self.c.Tr.CommitWithoutMessageErr)
	}

	err := self.c.Contexts().CommitMessage.OnConfirm(summary, description)
	if err != nil {
		return err
	}

	return nil
}

func (self *CommitsHelper) PreserveCommitMessage() {
	if self.c.Contexts().CommitMessage.GetPreserveMessage() {
		message := self.JoinCommitMessageAndUnwrappedDescription()
		if message != self.c.Contexts().CommitMessage.GetInitialMessage() {
			self.c.Contexts().CommitMessage.SetPreservedMessageAndLogError(message)
		}
	}
}

func (self *CommitsHelper) CloseCommitMessagePanel() {
	self.PreserveCommitMessage()

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
			Keys:           menuKey('e'),
			DisabledReason: disabledReasonForOpenInEditor,
		},
	}

	if strings.TrimSpace(self.c.UserConfig().Git.Commit.MessageGeneratorCommand) != "" {
		menuItems = append(menuItems, &types.MenuItem{
			Label: self.c.Tr.GenerateCommitMessage,
			OnPress: func() error {
				return self.generateCommitMessage()
			},
			Keys: menuKey('g'),
		})
	}

	menuItems = append(menuItems,
		&types.MenuItem{
			Label: self.c.Tr.AddCoAuthor,
			OnPress: func() error {
				return self.addCoAuthor(suggestionFunc)
			},
			Keys: menuKey('c'),
		},
		&types.MenuItem{
			Label: self.c.Tr.PasteCommitMessageFromClipboard,
			OnPress: func() error {
				return self.pasteCommitMessageFromClipboard()
			},
			Keys: menuKey('p'),
		},
	)
	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.CommitMenuTitle,
		Items: menuItems,
	})
}

func (self *CommitsHelper) generateCommitMessage() error {
	if self.c.Contexts().CommitMessage.IsGeneratingCommitMessage() {
		return nil
	}

	command := strings.TrimSpace(self.c.UserConfig().Git.Commit.MessageGeneratorCommand)
	repoRoot := self.c.Git().RepoPaths.WorktreePath()
	commandWithRepoRoot := command + " " + self.c.OS().Quote(repoRoot)

	cmdObj := self.c.OS().Cmd.NewShell(commandWithRepoRoot, self.c.UserConfig().OS.ShellFunctionsFile).SetWd(repoRoot)
	cmd := cmdObj.GetCmd()

	var stdoutBuffer, stderrBuffer bytes.Buffer
	cmd.Stdout = &stdoutBuffer
	cmd.Stderr = &stderrBuffer

	var mutex sync.Mutex
	cancelled := false
	var process *os.Process
	done := make(chan struct{})

	cancel := func() {
		mutex.Lock()
		if cancelled {
			mutex.Unlock()
			return
		}
		cancelled = true
		processToTerminate := process
		mutex.Unlock()

		if processToTerminate == nil {
			return
		}

		if err := oscommands.TerminateProcessGracefully(cmd); err != nil {
			self.c.Log.Errorf("error when trying to terminate commit message generator: %v; Command: %v %v", err, cmd.Path, cmd.Args)
		}

		go func() {
			select {
			case <-done:
			case <-time.After(2 * time.Second):
				_ = processToTerminate.Kill()
			}
		}()
	}

	self.c.Contexts().CommitMessage.StartGeneratingCommitMessage(cancel)
	stopRendering := self.renderGeneratingCommitMessageStatus()

	return self.c.WithWaitingStatus(self.c.Tr.GeneratingCommitMessageStatus, func(gocui.Task) error {
		defer close(done)
		defer close(stopRendering)

		err := cmd.Start()
		mutex.Lock()
		if err == nil {
			process = cmd.Process
		}
		shouldTerminate := cancelled && process != nil
		mutex.Unlock()

		if shouldTerminate {
			_ = oscommands.TerminateProcessGracefully(cmd)
		}

		if err == nil {
			err = cmd.Wait()
		}

		mutex.Lock()
		wasCancelled := cancelled
		mutex.Unlock()

		self.c.OnUIThread(func() error {
			self.c.Contexts().CommitMessage.StopGeneratingCommitMessage()
			if wasCancelled {
				return nil
			}

			if err != nil {
				message := strings.TrimSpace(stderrBuffer.String())
				if message == "" {
					message = err.Error()
				}
				self.c.Alert(self.c.Tr.GenerateCommitMessageFailed, message)
				return nil
			}

			self.SetMessageAndDescriptionInView(stdoutBuffer.String())
			return nil
		})

		return nil
	})
}

func (self *CommitsHelper) renderGeneratingCommitMessageStatus() chan struct{} {
	stop := make(chan struct{})

	self.c.OnWorker(func(gocui.Task) error {
		ticker := time.NewTicker(time.Millisecond * time.Duration(self.c.UserConfig().Gui.Spinner.Rate))
		defer ticker.Stop()

		for {
			select {
			case <-stop:
				return nil
			case <-ticker.C:
				self.c.OnUIThreadContentOnly(func() error {
					self.c.Contexts().CommitMessage.RenderCommitDescriptionSubtitle()
					return nil
				})
			}
		}
	})

	return stop
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

	currentMessage := self.JoinCommitMessageAndUnwrappedDescription()
	return self.c.ConfirmIf(currentMessage != "", types.ConfirmOpts{
		Title:  self.c.Tr.PasteCommitMessageFromClipboard,
		Prompt: self.c.Tr.SurePasteCommitMessage,
		HandleConfirm: func() error {
			self.SetMessageAndDescriptionInView(message)
			return nil
		},
	})
}
