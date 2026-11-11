// Package ai generates commit messages using any OpenAI-compatible chat
// completions API (Groq, OpenRouter, OpenAI, Ollama, LM Studio, ...).
package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const DefaultSystemPrompt = `You are an expert software engineer writing git commit messages.
Given a staged diff, write a single commit message that describes the change.

Rules:
- First line: a concise summary in the imperative mood ("Add", "Fix", "Refactor"), at most 72 characters, no trailing period.
- If the change is non-trivial, add a blank line and then a short body explaining what changed and why. Wrap body lines at 72 characters.
- Match the style of the recent commit messages provided (e.g. use Conventional Commits prefixes like "feat:" only if the repository uses them).
- Do not invent details that are not supported by the diff.
- Output only the commit message itself: no preamble, no quotes, no markdown code fences.`

type Options struct {
	// Base URL of an OpenAI-compatible API, e.g. https://api.groq.com/openai/v1
	BaseURL string
	APIKey  string
	Model   string
	// System prompt; DefaultSystemPrompt is used if empty.
	SystemPrompt string
	// Diffs longer than this many bytes are truncated. 0 means no limit.
	MaxDiffLength int
	HTTPClient    *http.Client
}

type Input struct {
	// Output of `git diff --staged --stat`
	DiffStat string
	// Output of `git diff --staged`
	Diff string
	// Subjects of recent commits, newest first, to let the model match the
	// repository's style.
	RecentCommitSubjects []string
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type chatResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

var ErrNoChanges = errors.New("no staged changes")

// GenerateCommitMessage asks the model for a commit message describing input.
func GenerateCommitMessage(ctx context.Context, opts Options, input Input) (string, error) {
	if strings.TrimSpace(input.Diff) == "" {
		return "", ErrNoChanges
	}

	systemPrompt := opts.SystemPrompt
	if strings.TrimSpace(systemPrompt) == "" {
		systemPrompt = DefaultSystemPrompt
	}

	body, err := json.Marshal(chatRequest{
		Model: opts.Model,
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: BuildUserPrompt(input, opts.MaxDiffLength)},
		},
	})
	if err != nil {
		return "", err
	}

	url := strings.TrimRight(opts.BaseURL, "/") + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	if opts.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+opts.APIKey)
	}

	client := opts.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 90 * time.Second}
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("AI request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return "", fmt.Errorf("AI request failed: %w", err)
	}

	var parsed chatResponse
	jsonErr := json.Unmarshal(respBody, &parsed)

	if resp.StatusCode != http.StatusOK {
		detail := strings.TrimSpace(string(respBody))
		if jsonErr == nil && parsed.Error != nil && parsed.Error.Message != "" {
			detail = parsed.Error.Message
		}
		if len(detail) > 500 {
			detail = detail[:500] + "..."
		}
		return "", fmt.Errorf("AI request failed (HTTP %d): %s", resp.StatusCode, detail)
	}
	if jsonErr != nil {
		return "", fmt.Errorf("AI request returned invalid JSON: %w", jsonErr)
	}
	if len(parsed.Choices) == 0 {
		return "", errors.New("AI request returned no choices")
	}

	message := CleanMessage(parsed.Choices[0].Message.Content)
	if message == "" {
		return "", errors.New("AI returned an empty commit message")
	}
	return message, nil
}

// BuildUserPrompt assembles the user message sent to the model.
func BuildUserPrompt(input Input, maxDiffLength int) string {
	var sb strings.Builder

	if len(input.RecentCommitSubjects) > 0 {
		sb.WriteString("Recent commit messages in this repository (for style reference):\n")
		for _, subject := range input.RecentCommitSubjects {
			sb.WriteString("- " + subject + "\n")
		}
		sb.WriteString("\n")
	}

	if stat := strings.TrimSpace(input.DiffStat); stat != "" {
		sb.WriteString("Summary of staged changes:\n" + stat + "\n\n")
	}

	diff := input.Diff
	if maxDiffLength > 0 && len(diff) > maxDiffLength {
		diff = diff[:maxDiffLength] + "\n[... diff truncated ...]"
	}
	sb.WriteString("Staged diff:\n" + diff)

	return sb.String()
}

var (
	thinkBlockRegex = regexp.MustCompile(`(?s)<think>.*?</think>`)
	codeFenceRegex  = regexp.MustCompile("(?s)^```[a-zA-Z]*\n(.*?)\n?```$")
)

// CleanMessage strips artifacts that models commonly wrap around the message:
// reasoning blocks, markdown code fences, surrounding quotes and whitespace.
func CleanMessage(message string) string {
	message = thinkBlockRegex.ReplaceAllString(message, "")
	message = strings.TrimSpace(message)
	if m := codeFenceRegex.FindStringSubmatch(message); m != nil {
		message = strings.TrimSpace(m[1])
	}
	if len(message) >= 2 && message[0] == '"' && message[len(message)-1] == '"' && !strings.Contains(message, "\n") {
		message = message[1 : len(message)-1]
	}
	return strings.TrimSpace(message)
}
