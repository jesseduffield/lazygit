package git_commands

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/git_config"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/stretchr/testify/assert"
)

func TestParseGitFlowPrefixMap(t *testing.T) {
	type scenario struct {
		testName     string
		legacyOutput string
		expected     map[string]string
	}
	scenarios := []scenario{
		{
			testName:     "empty input",
			legacyOutput: "",
			expected:     map[string]string{},
		},
		{
			testName:     "feature and hotfix",
			legacyOutput: "gitflow.prefix.feature feature/\ngitflow.prefix.hotfix hotfix/",
			expected:     map[string]string{"feature/": "feature", "hotfix/": "hotfix"},
		},
		{
			testName:     "prefix normalized with trailing slash",
			legacyOutput: "gitflow.prefix.feature feature",
			expected:     map[string]string{"feature/": "feature"},
		},
		{
			testName:     "malformed lines skipped",
			legacyOutput: "gitflow.prefix.feature feature/\nnot-a-valid-line\ngitflow.prefix.hotfix hotfix/",
			expected:     map[string]string{"feature/": "feature", "hotfix/": "hotfix"},
		},
		{
			testName:     "blank lines and whitespace ignored",
			legacyOutput: "  \n  gitflow.prefix.feature feature/  \n  \n  ",
			expected:     map[string]string{"feature/": "feature"},
		},
	}
	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			got := parseGitFlowPrefixMap(s.legacyOutput)
			assert.Equal(t, s.expected, got)
		})
	}
}

func TestGetGitFlowPrefixMap(t *testing.T) {
	type scenario struct {
		testName               string
		gitConfigMockResponses map[string]string
		expected               map[string]string
	}
	scenarios := []scenario{
		{
			testName:               "empty when no config",
			gitConfigMockResponses: nil,
			expected:               map[string]string{},
		},
		{
			testName: "correct map from legacy output",
			gitConfigMockResponses: map[string]string{
				"--local --get-regexp gitflow.prefix": "gitflow.prefix.feature feature/\ngitflow.prefix.hotfix hotfix/",
			},
			expected: map[string]string{"feature/": "feature", "hotfix/": "hotfix"},
		},
		{
			testName: "prefix normalized with trailing slash",
			gitConfigMockResponses: map[string]string{
				"--local --get-regexp gitflow.prefix": "gitflow.prefix.feature feature",
			},
			expected: map[string]string{"feature/": "feature"},
		},
		{
			testName: "malformed lines skipped",
			gitConfigMockResponses: map[string]string{
				"--local --get-regexp gitflow.prefix": "gitflow.prefix.feature feature/\nnot-a-valid-line\n",
			},
			expected: map[string]string{"feature/": "feature"},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			config := NewConfigCommands(common.NewDummyCommon(), git_config.NewFakeGitConfig(s.gitConfigMockResponses))
			got := config.GetGitFlowPrefixMap()
			assert.Equal(t, s.expected, got)
		})
	}
}
