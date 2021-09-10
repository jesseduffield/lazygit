package commands

import (
	"os/exec"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/secureexec"
	"github.com/stretchr/testify/assert"
)

// TestGitCommandStashDo is a function.
func TestGitCommandStashDo(t *testing.T) {
	gitCmd := NewDummyGitCommand()
	gitCmd.OSCommand.Command = func(cmd string, args ...string) *exec.Cmd {
		assert.EqualValues(t, "git", cmd)
		assert.EqualValues(t, []string{"stash", "drop", "stash@{1}"}, args)

		return secureexec.Command("echo")
	}

	assert.NoError(t, gitCmd.StashDo(1, "drop"))
}

// TestGitCommandStashSave is a function.
func TestGitCommandStashSave(t *testing.T) {
	gitCmd := NewDummyGitCommand()
	gitCmd.OSCommand.Command = func(cmd string, args ...string) *exec.Cmd {
		assert.EqualValues(t, "git", cmd)
		assert.EqualValues(t, []string{"stash", "save", "A stash message"}, args)

		return secureexec.Command("echo")
	}

	assert.NoError(t, gitCmd.StashSave("A stash message"))
}

// TestGitCommandShowStashEntryCmdStr is a function.
func TestGitCommandShowStashEntryCmdStr(t *testing.T) {
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
			gitCmd.Config.GetUserConfig().Git.DiffContextSize = s.contextSize
			cmdStr := gitCmd.ShowStashEntryCmdStr(s.index)
			assert.Equal(t, s.expected, cmdStr)
		})
	}
}
