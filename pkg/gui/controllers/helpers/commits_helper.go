package helpers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"sync/atomic"
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

	// set to 1 while AI commit message generation is in progress
	generating atomic.Int32
}

func (self *CommitsHelper) IsGenerating() bool {
	return self.generating.Load() == 1
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

	aiConfig := self.c.UserConfig().Git.Commit.AI
	aiConfigured := aiConfig.CLI.Command != "" || aiConfig.API.Endpoint != ""
	var disabledReasonForAI *types.DisabledReason
	if !aiConfigured {
		disabledReasonForAI = &types.DisabledReason{
			Text: self.c.Tr.NoAIConfigured,
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
		{
			Label: self.c.Tr.GenerateCommitMessageWithAI,
			OnPress: func() error {
				return self.generateCommitMessageWithAI()
			},
			Key:            'a',
			DisabledReason: disabledReasonForAI,
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

func (self *CommitsHelper) generateCommitMessageWithAI() error {
	self.generating.Store(1)
	self.c.Views().CommitMessage.Editable = false
	self.c.Views().CommitDescription.Editable = false
	originalTitle := self.c.Views().CommitMessage.Title
	self.c.Views().CommitMessage.Title = self.c.Tr.GeneratingCommitMessageStatus

	restore := func() {
		self.c.OnUIThread(func() error {
			self.generating.Store(0)
			self.c.Views().CommitMessage.Editable = true
			self.c.Views().CommitDescription.Editable = true
			self.c.Views().CommitMessage.Title = originalTitle
			return nil
		})
	}

	return self.c.WithWaitingStatus(self.c.Tr.GeneratingCommitMessageStatus, func(_ gocui.Task) error {
		defer restore()

		diff, err := self.c.Git().Diff.GetDiff(true)
		if err != nil {
			return err
		}
		if strings.TrimSpace(diff) == "" {
			self.c.OnUIThread(func() error {
				self.c.ErrorToast(self.c.Tr.NoStagedChangesForAI)
				return nil
			})
			return nil
		}

		cfg := self.c.UserConfig().Git.Commit.AI
		var message string
		if cfg.CLI.Command != "" {
			message, err = self.generateViaAICLI(cfg.CLI.Command, diff)
		} else {
			message, err = generateViaAIAPI(cfg.API.Endpoint, cfg.API.Model, cfg.API.APIKey, cfg.API.SystemPrompt, diff)
		}
		if err != nil {
			return err
		}

		message = parseAIOutput(message)
		if message == "" {
			return errors.New("AI returned an empty commit message")
		}

		self.c.OnUIThread(func() error {
			self.SetMessageAndDescriptionInView(message)
			return nil
		})
		return nil
	})
}

func (self *CommitsHelper) generateViaAICLI(command string, diff string) (string, error) {
	return self.c.OS().Cmd.NewShell(command, "").SetStdin(diff).DontLog().RunWithOutput()
}

type aiAPIRequest struct {
	Model    string       `json:"model"`
	Messages []aiMessage  `json:"messages"`
}

type aiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type aiAPIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func generateViaAIAPI(endpoint, model, apiKey, systemPrompt, diff string) (string, error) {
	if systemPrompt == "" {
		systemPrompt = "Generate a conventional commit message for the following diff. Output only the commit message (subject and optional body), with no additional explanation or markdown fences."
	}

	reqBody := aiAPIRequest{
		Model: model,
		Messages: []aiMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: diff},
		},
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", strings.TrimRight(endpoint, "/")+"/chat/completions", bytes.NewReader(bodyBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("AI API returned status %d", resp.StatusCode)
	}

	var result aiAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if len(result.Choices) == 0 {
		return "", errors.New("AI API returned no choices")
	}
	return result.Choices[0].Message.Content, nil
}

// parseAIOutput strips markdown code fences and any preamble before them.
// Handles tools like opencode that emit lines like "> build • glm-5" before the fence.
func parseAIOutput(s string) string {
	s = strings.TrimSpace(s)
	if start := strings.Index(s, "```"); start >= 0 {
		// skip past the opening fence line (e.g. "```" or "```text")
		rest := s[start:]
		if idx := strings.Index(rest, "\n"); idx >= 0 {
			rest = rest[idx+1:]
		}
		// drop the closing fence
		if idx := strings.LastIndex(rest, "```"); idx >= 0 {
			rest = rest[:idx]
		}
		s = strings.TrimSpace(rest)
	}
	return s
}
