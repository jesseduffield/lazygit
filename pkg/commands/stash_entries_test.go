package commands

import (
	"os/exec"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/secureexec"
	"github.com/stretchr/testify/assert"
)

// TestGitCommandStashDo is a function.
func TestGitCommandStashDo(t *testing.T) {
	gitCmd := NewDummyGit()
	gitCmd.GetOSCommand().Command = func(cmd string, args ...string) *exec.Cmd {
		assert.EqualValues(t, "git", cmd)
		assert.EqualValues(t, []string{"stash", "drop", "stash@{1}"}, args)

		return secureexec.Command("echo")
	}

	assert.NoError(t, gitCmd.StashDo(1, "drop"))
}

// TestGitCommandStashSave is a function.
func TestGitCommandStashSave(t *testing.T) {
	gitCmd := NewDummyGit()
	gitCmd.GetOSCommand().Command = func(cmd string, args ...string) *exec.Cmd {
		assert.EqualValues(t, "git", cmd)
		assert.EqualValues(t, []string{"stash", "save", "A stash message"}, args)

		return secureexec.Command("echo")
	}

	assert.NoError(t, gitCmd.StashSave("A stash message"))
}
