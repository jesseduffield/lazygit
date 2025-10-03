package ai

import "errors"

// This file contains the data models and interfaces for the AI service.
// These types define the contract between the AI service components and
// external AI providers for generating commit messages.

// Provider interface for different AI providers (OpenAI, GitHub Copilot, etc.)
type Provider interface {
	GenerateMessage(prompt string) (string, error)
	Name() string
	ValidateConfig() error
}

// GenerateRequest contains the context for generating a commit message
type GenerateRequest struct {
	StagedDiff      string       `json:"staged_diff"`
	FileChanges     []FileChange `json:"file_changes"`
	BranchName      string       `json:"branch_name"`
	RecentCommits   []string     `json:"recent_commits"`
	ProjectName     string       `json:"project_name"`
	CommitType      string       `json:"commit_type"`                // "new" or "reword"
	ExistingMessage string       `json:"existing_message,omitempty"` // For reword operations
}

// GenerateResponse contains the AI-generated commit message
type GenerateResponse struct {
	Message     string  `json:"message"`
	Description string  `json:"description,omitempty"`
	Confidence  float64 `json:"confidence"`
	Provider    string  `json:"provider"`
}

// FileChange represents a change to a file in the commit
type FileChange struct {
	Path         string `json:"path"`
	Status       string `json:"status"`   // "added", "modified", "deleted", "renamed"
	Language     string `json:"language"` // Programming language detected
	LinesAdded   int    `json:"lines_added"`
	LinesDeleted int    `json:"lines_deleted"`
	IsBinary     bool   `json:"is_binary"`
}

// GitContext contains git repository context for message generation
type GitContext struct {
	BranchName          string   `json:"branch_name"`
	ProjectName         string   `json:"project_name"`
	RecentCommits       []string `json:"recent_commits"`
	RepositoryType      string   `json:"repository_type"`      // detected from files
	ConventionalCommits bool     `json:"conventional_commits"` // whether to use conventional format
}

// Common errors
var (
	ErrUnsupportedProvider  = errors.New("unsupported AI provider")
	ErrNotConfigured        = errors.New("AI service is not configured")
	ErrInvalidAPIKey        = errors.New("invalid API key")
	ErrNetworkError         = errors.New("network error while calling AI service")
	ErrInvalidResponse      = errors.New("invalid response from AI service")
	ErrEmptyDiff            = errors.New("no staged changes to generate commit message for")
	ErrMessageTooLong       = errors.New("generated commit message is too long")
	ErrInappropriateContent = errors.New("generated message contains inappropriate content")
)
