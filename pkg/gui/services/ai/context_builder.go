package ai

import (
	"path/filepath"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
)

// ContextBuilder builds AI context from git repository state
type ContextBuilder struct {
	c *helpers.HelperCommon
}

// NewContextBuilder creates a new context builder
func NewContextBuilder(c *helpers.HelperCommon) *ContextBuilder {
	return &ContextBuilder{c: c}
}

// BuildContext builds a GenerateRequest from current git state
func (cb *ContextBuilder) BuildContext(commitType string, existingMessage string) (*GenerateRequest, error) {
	// TODO: Implement context building
	// 1. Get staged diff
	// 2. Analyze file changes
	// 3. Get branch name
	// 4. Get recent commits for style reference
	// 5. Detect project characteristics
	
	request := &GenerateRequest{
		CommitType:      commitType,
		ExistingMessage: existingMessage,
	}
	
	// Get staged diff
	if err := cb.setStagedDiff(request); err != nil {
		return nil, err
	}
	
	// Analyze file changes
	if err := cb.setFileChanges(request); err != nil {
		return nil, err
	}
	
	// Set git context
	cb.setGitContext(request)
	
	return request, nil
}

// setStagedDiff gets the staged diff and sets it in the request
func (cb *ContextBuilder) setStagedDiff(request *GenerateRequest) error {
	// TODO: Execute `git diff --staged --no-color` and set StagedDiff
	// Handle case where there are no staged changes
	request.StagedDiff = ""
	return nil
}

// setFileChanges analyzes the staged changes and extracts file information
func (cb *ContextBuilder) setFileChanges(request *GenerateRequest) error {
	// TODO: Parse git diff to extract file changes
	// 1. Get list of changed files with status
	// 2. Detect programming languages
	// 3. Count lines added/deleted
	// 4. Detect binary files
	
	request.FileChanges = []FileChange{}
	return nil
}

// setGitContext sets git repository context information
func (cb *ContextBuilder) setGitContext(request *GenerateRequest) {
	// TODO: Get current branch name
	request.BranchName = cb.getBranchName()
	
	// TODO: Get project name from repository
	request.ProjectName = cb.getProjectName()
	
	// TODO: Get recent commits for style analysis
	request.RecentCommits = cb.getRecentCommits()
}

// getBranchName returns the current git branch name
func (cb *ContextBuilder) getBranchName() string {
	// TODO: Execute `git branch --show-current` or equivalent
	return ""
}

// getProjectName extracts the project name from repository path or remote
func (cb *ContextBuilder) getProjectName() string {
	// TODO: Get project name from:
	// 1. Repository directory name
	// 2. Remote origin URL
	// 3. package.json, go.mod, etc.
	return ""
}

// getRecentCommits gets recent commit messages for style analysis
func (cb *ContextBuilder) getRecentCommits() []string {
	// TODO: Execute `git log --oneline -10` to get recent commits
	// Filter out merge commits and format consistently
	return []string{}
}

// detectLanguage detects the programming language from file extension
func (cb *ContextBuilder) detectLanguage(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	
	languageMap := map[string]string{
		".go":   "Go",
		".js":   "JavaScript", 
		".ts":   "TypeScript",
		".py":   "Python",
		".java": "Java",
		".cpp":  "C++",
		".c":    "C",
		".rs":   "Rust",
		".rb":   "Ruby",
		".php":  "PHP",
		".cs":   "C#",
		".swift": "Swift",
		".kt":   "Kotlin",
		".scala": "Scala",
		".sh":   "Shell",
		".yml":  "YAML",
		".yaml": "YAML",
		".json": "JSON",
		".xml":  "XML",
		".md":   "Markdown",
		".sql":  "SQL",
		".dockerfile": "Docker",
	}
	
	if lang, exists := languageMap[ext]; exists {
		return lang
	}
	
	// Check for special files
	filename := strings.ToLower(filepath.Base(filePath))
	switch filename {
	case "dockerfile":
		return "Docker"
	case "makefile":
		return "Make"
	case "rakefile":
		return "Ruby"
	default:
		return "Text"
	}
}

// BuildPrompt builds the AI prompt from the request context
func (cb *ContextBuilder) BuildPrompt(request *GenerateRequest) string {
	// TODO: Build a well-structured prompt for the AI
	// Include:
	// 1. Task description
	// 2. Context about the changes
	// 3. Style guidelines
	// 4. Examples from recent commits
	// 5. Specific requirements (conventional commits, etc.)
	
	return ""
}