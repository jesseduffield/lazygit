package ai

import (
	"strings"
	"unicode/utf8"
)

// MessageValidator validates AI-generated commit messages
type MessageValidator struct {
	maxSubjectLength int
	maxBodyLength    int
	requireSubject   bool
}

// NewMessageValidator creates a new message validator with default settings
func NewMessageValidator() *MessageValidator {
	return &MessageValidator{
		maxSubjectLength: 72, // Standard git convention
		maxBodyLength:    72, // Per line in body
		requireSubject:   true,
	}
}

// ValidateMessage validates a generated commit message
func (mv *MessageValidator) ValidateMessage(response *GenerateResponse) error {
	if response == nil {
		return ErrInvalidResponse
	}

	// Check if message is empty
	if strings.TrimSpace(response.Message) == "" {
		return ErrInvalidResponse
	}

	// Validate subject line
	if err := mv.validateSubject(response.Message); err != nil {
		return err
	}

	// Validate message length
	if err := mv.validateLength(response.Message); err != nil {
		return err
	}

	// Check for inappropriate content
	if err := mv.validateContent(response.Message); err != nil {
		return err
	}

	return nil
}

// validateSubject validates the commit message subject line
func (mv *MessageValidator) validateSubject(message string) error {
	lines := strings.Split(message, "\n")
	if len(lines) == 0 {
		return ErrInvalidResponse
	}

	subject := strings.TrimSpace(lines[0])

	// Check if subject is required and present
	if mv.requireSubject && subject == "" {
		return ErrInvalidResponse
	}

	// Check subject length
	if utf8.RuneCountInString(subject) > mv.maxSubjectLength {
		return ErrMessageTooLong
	}

	// TODO: Add more subject validation rules:
	// - No trailing period
	// - Capitalized first letter
	// - Imperative mood check
	// - Conventional commit format validation

	return nil
}

// validateLength validates the overall message length
func (mv *MessageValidator) validateLength(message string) error {
	lines := strings.Split(message, "\n")

	// Check body line lengths (skip empty lines and subject)
	for i, line := range lines {
		if i == 0 || strings.TrimSpace(line) == "" {
			continue // Skip subject line and empty lines
		}

		if utf8.RuneCountInString(line) > mv.maxBodyLength {
			return ErrMessageTooLong
		}
	}

	return nil
}

// validateContent checks for inappropriate content in the message
func (mv *MessageValidator) validateContent(message string) error {
	// TODO: Implement content filtering
	// Check for:
	// - Profanity or inappropriate language
	// - Personal information (emails, API keys, etc.)
	// - Nonsensical or irrelevant content
	// - Obvious AI-generated artifacts

	// Basic checks for now
	message = strings.ToLower(message)

	// Check for placeholder text that AI might generate
	prohibitedPhrases := []string{
		"lorem ipsum",
		"placeholder",
		"todo:",
		"fixme:",
		"[your text here]",
		"replace this",
	}

	for _, phrase := range prohibitedPhrases {
		if strings.Contains(message, phrase) {
			return ErrInappropriateContent
		}
	}

	return nil
}

// SanitizeMessage cleans up the generated message
func (mv *MessageValidator) SanitizeMessage(message string) string {
	// TODO: Implement message sanitization
	// - Remove extra whitespace
	// - Fix capitalization
	// - Remove trailing periods from subject
	// - Ensure proper line breaks

	// Basic cleanup for now
	lines := strings.Split(message, "\n")
	var cleanLines []string

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if i == 0 {
			// Subject line: remove trailing period, capitalize first letter
			line = strings.TrimSuffix(line, ".")
			if len(line) > 0 {
				line = strings.ToUpper(string(line[0])) + line[1:]
			}
		}
		if line != "" || (i > 0 && i < len(lines)-1) {
			cleanLines = append(cleanLines, line)
		}
	}

	return strings.Join(cleanLines, "\n")
}

// IsConventionalCommit checks if the message follows conventional commit format
func (mv *MessageValidator) IsConventionalCommit(message string) bool {
	// TODO: Implement conventional commit validation
	// Check for format: type(scope): description
	// Common types: feat, fix, docs, style, refactor, test, chore

	lines := strings.Split(message, "\n")
	if len(lines) == 0 {
		return false
	}

	subject := strings.TrimSpace(lines[0])

	// Basic pattern matching for conventional commits
	conventionalTypes := []string{
		"feat:", "fix:", "docs:", "style:", "refactor:",
		"test:", "chore:", "perf:", "ci:", "build:",
		"revert:", "merge:", "release:",
	}

	for _, ctype := range conventionalTypes {
		if strings.HasPrefix(strings.ToLower(subject), ctype) {
			return true
		}
	}

	// Check for scoped format: type(scope):
	if strings.Contains(subject, "(") && strings.Contains(subject, "):") {
		return true
	}

	return false
}
