package commands

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/stretchr/testify/assert"
)

func TestGitCommandStashDo(t *testing.T) {
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs([]string{"stash", "drop", "stash@{1}"}, "", nil)
	gitCmd := NewDummyGitCommandWithRunner(runner)

	assert.NoError(t, gitCmd.StashDo(1, "drop"))
	runner.CheckForMissingCalls()
}

func TestGitCommandStashSave(t *testing.T) {
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs([]string{"stash", "save", "A stash message"}, "", nil)
	gitCmd := NewDummyGitCommandWithRunner(runner)

	assert.NoError(t, gitCmd.StashSave("A stash message"))
	runner.CheckForMissingCalls()
}

func TestGitCommandShowStashEntryCmdObj(t *testing.T) {
	type scenario struct {
		testName    string
		index       int
		contextSize int
		expected    string
	}

	scenarios := []scenario{
		{
			testName:    "Default case",
			index:       5,
			contextSize: 3,
			expected:    "git stash show -p --stat --color=always --unified=3 stash@{5}",
		},
		{
			testName:    "Show diff with custom context size",
			index:       5,
			contextSize: 77,
			expected:    "git stash show -p --stat --color=always --unified=77 stash@{5}",
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCmd := NewDummyGitCommand()
			gitCmd.UserConfig.Git.DiffContextSize = s.contextSize
			cmdStr := gitCmd.ShowStashEntryCmdObj(s.index).ToString()
			assert.Equal(t, s.expected, cmdStr)
		})
	}
}
