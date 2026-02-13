package git_commands

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/git_config"
	"github.com/stretchr/testify/assert"
)

func TestGitFlowEnabled(t *testing.T) {
	type scenario struct {
		testName               string
		expected               bool
		gitConfigMockResponses map[string]string
	}
	scenarios := []scenario{
		{
			testName:               "disabled when no config",
			expected:               false,
			gitConfigMockResponses: nil,
		},
		{
			testName: "enabled with legacy config",
			expected: true,
			gitConfigMockResponses: map[string]string{
				"--local --get-regexp gitflow.prefix": "gitflow.prefix.feature feature/",
			},
		},
		{
			testName: "enabled with git-flow-next only config",
			expected: true,
			gitConfigMockResponses: map[string]string{
				"--local --get-regexp gitflow\\.branch\\..*\\.prefix": "gitflow.branch.feature.prefix feature/",
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			instance := buildFlowCommands(commonDeps{
				gitConfig: git_config.NewFakeGitConfig(s.gitConfigMockResponses),
			})
			assert.Equal(t, s.expected, instance.GitFlowEnabled())
		})
	}
}

func TestStartCmdObj(t *testing.T) {
	type scenario struct {
		testName   string
		branchType string
		branchName string
		expected   []string
	}
	scenarios := []scenario{
		{
			testName:   "basic",
			branchType: "feature",
			branchName: "test",
			expected:   []string{"git", "flow", "feature", "start", "test"},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			instance := buildFlowCommands(commonDeps{})

			assert.Equal(t,
				instance.StartCmdObj(s.branchType, s.branchName).Args(),
				s.expected,
			)
		})
	}
}

func TestFinishCmdObj(t *testing.T) {
	type scenario struct {
		testName               string
		branchName             string
		expected               []string
		expectedError          string
		gitConfigMockResponses map[string]string
	}
	scenarios := []scenario{
		{
			testName:               "not a git flow branch",
			branchName:             "mybranch",
			expected:               nil,
			expectedError:          "This does not seem to be a git flow branch",
			gitConfigMockResponses: nil,
		},
		{
			testName:               "feature branch without config",
			branchName:             "feature/mybranch",
			expected:               nil,
			expectedError:          "This does not seem to be a git flow branch",
			gitConfigMockResponses: nil,
		},
		{
			testName:      "feature branch with config",
			branchName:    "feature/mybranch",
			expected:      []string{"git", "flow", "feature", "finish", "mybranch"},
			expectedError: "",
			gitConfigMockResponses: map[string]string{
				"--local --get-regexp gitflow.prefix": "gitflow.prefix.feature feature/",
			},
		},
		{
			testName:      "feature branch with git-flow-next only config",
			branchName:    "feature/mybranch",
			expected:      []string{"git", "flow", "feature", "finish", "mybranch"},
			expectedError: "",
			gitConfigMockResponses: map[string]string{
				"--local --get-regexp gitflow\\.branch\\..*\\.prefix": "gitflow.branch.feature.prefix feature/",
			},
		},
		{
			testName:      "legacy wins when both configs have same prefix",
			branchName:    "feature/mybranch",
			expected:      []string{"git", "flow", "feature", "finish", "mybranch"},
			expectedError: "",
			gitConfigMockResponses: map[string]string{
				"--local --get-regexp gitflow.prefix":                 "gitflow.prefix.feature feature/",
				"--local --get-regexp gitflow\\.branch\\..*\\.prefix": "gitflow.branch.feature.prefix feature/",
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			instance := buildFlowCommands(commonDeps{
				gitConfig: git_config.NewFakeGitConfig(s.gitConfigMockResponses),
			})

			cmd, err := instance.FinishCmdObj(s.branchName)

			if s.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, s.expectedError, err.Error())
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, s.expected, cmd.Args())
		})
	}
}
