package ai

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSystemPrompt(t *testing.T) {
	tests := []struct {
		name     string
		style    string
		wrap     int
		expected []string // Expected strings to be present in the prompt
	}{
		{
			name:  "conventional style",
			style: "conventional",
			wrap:  72,
			expected: []string{
				"Write a high-quality git commit message from the provided diff.",
				"Use Conventional Commits format: <type>(<scope>): <subject>",
				"Types: feat, fix, docs, style, refactor, perf, test, chore, build, ci",
				"Subject should be concise, ideally <= 72 chars.",
				"Wrap body at ~72 chars when useful.",
				"Focus on what and why; avoid restating diff line-by-line.",
				"Do not include code fences or markdown headings.",
			},
		},
		{
			name:  "conventional style with conv alias",
			style: "conv",
			wrap:  80,
			expected: []string{
				"Write a high-quality git commit message from the provided diff.",
				"Use Conventional Commits format: <type>(<scope>): <subject>",
				"Types: feat, fix, docs, style, refactor, perf, test, chore, build, ci",
				"Subject should be concise, ideally <= 72 chars.",
				"Wrap body at ~80 chars when useful.",
				"Focus on what and why; avoid restating diff line-by-line.",
				"Do not include code fences or markdown headings.",
			},
		},
		{
			name:  "plain style",
			style: "plain",
			wrap:  50,
			expected: []string{
				"Write a high-quality git commit message from the provided diff.",
				"Subject should be concise, ideally <= 72 chars.",
				"Wrap body at ~50 chars when useful.",
				"Focus on what and why; avoid restating diff line-by-line.",
				"Do not include code fences or markdown headings.",
			},
		},
		{
			name:  "empty style defaults to plain",
			style: "",
			wrap:  72,
			expected: []string{
				"Write a high-quality git commit message from the provided diff.",
				"Subject should be concise, ideally <= 72 chars.",
				"Wrap body at ~72 chars when useful.",
				"Focus on what and why; avoid restating diff line-by-line.",
				"Do not include code fences or markdown headings.",
			},
		},
		{
			name:  "unknown style defaults to plain",
			style: "unknown",
			wrap:  72,
			expected: []string{
				"Write a high-quality git commit message from the provided diff.",
				"Subject should be concise, ideally <= 72 chars.",
				"Wrap body at ~72 chars when useful.",
				"Focus on what and why; avoid restating diff line-by-line.",
				"Do not include code fences or markdown headings.",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompt := systemPrompt(tt.style, tt.wrap)
			
			// Check that all expected strings are present
			for _, expected := range tt.expected {
				assert.Contains(t, prompt, expected, "Expected string not found in prompt")
			}
			
			// Check that conventional commit rules are only present for conventional style
			hasConventionalRules := strings.Contains(prompt, "Use Conventional Commits format")
			if tt.style == "conventional" || tt.style == "conv" {
				assert.True(t, hasConventionalRules, "Conventional commit rules should be present")
			} else {
				assert.False(t, hasConventionalRules, "Conventional commit rules should not be present")
			}
		})
	}
}

func TestSystemPromptCaseInsensitive(t *testing.T) {
	tests := []struct {
		style    string
		hasConv  bool
	}{
		{"CONVENTIONAL", true},
		{"Conventional", true},
		{"conventional", true},
		{"CONV", true},
		{"Conv", true},
		{"conv", true},
		{"PLAIN", false},
		{"Plain", false},
		{"plain", false},
	}

	for _, tt := range tests {
		t.Run(tt.style, func(t *testing.T) {
			prompt := systemPrompt(tt.style, 72)
			hasConventionalRules := strings.Contains(prompt, "Use Conventional Commits format")
			assert.Equal(t, tt.hasConv, hasConventionalRules)
		})
	}
}

func TestBuildUserPrompt(t *testing.T) {
	tests := []struct {
		name     string
		diff     string
		expected string
	}{
		{
			name: "simple diff",
			diff: "diff --git a/file.go b/file.go\n+func newFeature() {}",
			expected: "Diff of changes (unified):\ndiff --git a/file.go b/file.go\n+func newFeature() {}\n\nReturn subject on first line, optional body after.",
		},
		{
			name: "empty diff",
			diff: "",
			expected: "Diff of changes (unified):\n\n\nReturn subject on first line, optional body after.",
		},
		{
			name: "multiline diff",
			diff: "diff --git a/file.go b/file.go\n-old line\n+new line\n@@ -1,3 +1,3 @@",
			expected: "Diff of changes (unified):\ndiff --git a/file.go b/file.go\n-old line\n+new line\n@@ -1,3 +1,3 @@\n\nReturn subject on first line, optional body after.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildUserPrompt(tt.diff)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPromptStructure(t *testing.T) {
	// Test that the system prompt has a consistent structure
	prompt := systemPrompt("conventional", 72)
	
	// Should be multiple lines
	lines := strings.Split(prompt, "\n")
	assert.Greater(t, len(lines), 3, "Prompt should have multiple lines")
	
	// Should start with the main instruction
	assert.True(t, strings.HasPrefix(prompt, "Write a high-quality git commit message"))
	
	// Should contain all the basic rules
	basicRules := []string{
		"Subject should be concise",
		"Focus on what and why",
		"Do not include code fences",
	}
	
	for _, rule := range basicRules {
		assert.Contains(t, prompt, rule)
	}
}

func TestPromptWrapParameter(t *testing.T) {
	tests := []int{50, 72, 80, 100, 120}

	for _, wrap := range tests {
		t.Run(fmt.Sprintf("wrap_%d", wrap), func(t *testing.T) {
			prompt := systemPrompt("plain", wrap)
			expectedWrapText := fmt.Sprintf("Wrap body at ~%d chars when useful.", wrap)
			assert.Contains(t, prompt, expectedWrapText)
		})
	}
}

func TestUserPromptFormat(t *testing.T) {
	diff := "test diff content"
	prompt := buildUserPrompt(diff)
	
	// Should have the expected structure
	assert.Contains(t, prompt, "Diff of changes (unified):")
	assert.Contains(t, prompt, diff)
	assert.Contains(t, prompt, "Return subject on first line, optional body after.")
	
	// Should have proper spacing
	parts := strings.Split(prompt, "\n\n")
	assert.Len(t, parts, 2, "Should have two main parts separated by double newline")
}
