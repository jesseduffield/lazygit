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
		nextOutput   string
		expected     map[string]string
	}
	scenarios := []scenario{
		{
			testName:     "empty inputs",
			legacyOutput: "",
			nextOutput:   "",
			expected:     map[string]string{},
		},
		{
			testName:     "legacy only",
			legacyOutput: "gitflow.prefix.feature feature/\ngitflow.prefix.hotfix hotfix/",
			nextOutput:   "",
			expected:     map[string]string{"feature/": "feature", "hotfix/": "hotfix"},
		},
		{
			testName:     "next only",
			legacyOutput: "",
			nextOutput:   "gitflow.branch.feature.prefix feature/\ngitflow.branch.release.prefix release/",
			expected:     map[string]string{"feature/": "feature", "release/": "release"},
		},
		{
			testName:     "legacy wins on duplicate prefix",
			legacyOutput: "gitflow.prefix.feature feature/",
			nextOutput:   "gitflow.branch.feature.prefix feature/",
			expected:     map[string]string{"feature/": "feature"},
		},
		{
			testName:     "prefix normalized with trailing slash from legacy",
			legacyOutput: "gitflow.prefix.feature feature",
			nextOutput:   "",
			expected:     map[string]string{"feature/": "feature"},
		},
		{
			testName:     "malformed legacy lines skipped",
			legacyOutput: "gitflow.prefix.feature feature/\nnot-a-valid-line\ngitflow.prefix.hotfix hotfix/",
			nextOutput:   "",
			expected:     map[string]string{"feature/": "feature", "hotfix/": "hotfix"},
		},
		{
			testName:     "blank lines and whitespace ignored",
			legacyOutput: "  \n  gitflow.prefix.feature feature/  \n  \n  ",
			nextOutput:   "",
			expected:     map[string]string{"feature/": "feature"},
		},
	}
	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			got := parseGitFlowPrefixMap(s.legacyOutput, s.nextOutput)
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
			testName:               "empty when both queries empty",
			gitConfigMockResponses: nil,
			expected:               map[string]string{},
		},
		{
			testName: "correct map from legacy-only output",
			gitConfigMockResponses: map[string]string{
				"--local --get-regexp gitflow.prefix": "gitflow.prefix.feature feature/\ngitflow.prefix.hotfix hotfix/",
			},
			expected: map[string]string{"feature/": "feature", "hotfix/": "hotfix"},
		},
		{
			testName: "correct map from git-flow-next-only output",
			gitConfigMockResponses: map[string]string{
				"--local --get-regexp gitflow\\.branch\\..*\\.prefix": "gitflow.branch.feature.prefix feature/\ngitflow.branch.release.prefix release/",
			},
			expected: map[string]string{"feature/": "feature", "release/": "release"},
		},
		{
			testName: "merged map with legacy winning when both have same prefix",
			gitConfigMockResponses: map[string]string{
				"--local --get-regexp gitflow.prefix":                 "gitflow.prefix.feature feature/",
				"--local --get-regexp gitflow\\.branch\\..*\\.prefix": "gitflow.branch.feature.prefix feature/",
			},
			expected: map[string]string{"feature/": "feature"},
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
			config := NewConfigCommands(common.NewDummyCommon(), git_config.NewFakeGitConfig(s.gitConfigMockResponses), nil)
			got := config.GetGitFlowPrefixMap()
			assert.Equal(t, s.expected, got)
		})
	}
}

func TestBranches_NilRepo(t *testing.T) {
	config := NewConfigCommands(common.NewDummyCommon(), git_config.NewFakeGitConfig(nil), nil)
	branches, err := config.Branches()
	assert.Error(t, err)
	assert.Nil(t, branches)
	assert.Equal(t, "repository is nil", err.Error())
}
