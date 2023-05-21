package git_commands

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/git_config"
	"github.com/stretchr/testify/assert"
)

func TestStartCmdObj(t *testing.T) {
	scenarios := []struct {
		testName   string
		branchType string
		name       string
		expected   []string
	}{
		{
			testName:   "basic",
			branchType: "feature",
			name:       "test",
			expected:   []string{"git", "flow", "feature", "start", "test"},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			instance := buildFlowCommands(commonDeps{})

			assert.Equal(t,
				instance.StartCmdObj(s.branchType, s.name).Args(),
				s.expected,
			)
		})
	}
}

func TestFinishCmdObj(t *testing.T) {
	scenarios := []struct {
		testName               string
		branchName             string
		expected               []string
		expectedError          string
		gitConfigMockResponses map[string]string
	}{
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
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			instance := buildFlowCommands(commonDeps{
				gitConfig: git_config.NewFakeGitConfig(s.gitConfigMockResponses),
			})

			cmd, err := instance.FinishCmdObj(s.branchName)

			if s.expectedError != "" {
				if err == nil {
					t.Errorf("Expected error, got nil")
				} else {
					assert.Equal(t, err.Error(), s.expectedError)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, cmd.Args(), s.expected)
			}
		})
	}
}
