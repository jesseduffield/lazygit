package commands

import (
	"os/exec"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/secureexec"
	"github.com/jesseduffield/lazygit/pkg/test"
	"github.com/stretchr/testify/assert"
)

// TestGitCommandRenameCommit is a function.
func TestGitCommandRenameCommit(t *testing.T) {
	gitCmd := NewDummyGitCommand()
	gitCmd.OSCommand.Command = func(cmd string, args ...string) *exec.Cmd {
		assert.EqualValues(t, "git", cmd)
		assert.EqualValues(t, []string{"commit", "--allow-empty", "--amend", "--only", "-m", "test"}, args)

		return secureexec.Command("echo")
	}

	assert.NoError(t, gitCmd.RenameCommit("test"))
}

// TestGitCommandResetToCommit is a function.
func TestGitCommandResetToCommit(t *testing.T) {
	gitCmd := NewDummyGitCommand()
	gitCmd.OSCommand.Command = func(cmd string, args ...string) *exec.Cmd {
		assert.EqualValues(t, "git", cmd)
		assert.EqualValues(t, []string{"reset", "--hard", "78976bc"}, args)

		return secureexec.Command("echo")
	}

	assert.NoError(t, gitCmd.ResetToCommit("78976bc", "hard", oscommands.RunCommandOptions{}))
}

// TestGitCommandCommitStr is a function.
func TestGitCommandCommitStr(t *testing.T) {
	gitCmd := NewDummyGitCommand()

	type scenario struct {
		testName string
		message  string
		flags    string
		expected string
	}

	scenarios := []scenario{
		{
			testName: "Commit",
			message:  "test",
			flags:    "",
			expected: "git commit -m " + gitCmd.OSCommand.Quote("test"),
		},
		{
			testName: "Commit with --no-verify flag",
			message:  "test",
			flags:    "--no-verify",
			expected: "git commit --no-verify -m " + gitCmd.OSCommand.Quote("test"),
		},
		{
			testName: "Commit with multiline message",
			message:  "line1\nline2",
			flags:    "",
			expected: "git commit -m " + gitCmd.OSCommand.Quote("line1") + " -m " + gitCmd.OSCommand.Quote("line2"),
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			cmdStr := gitCmd.CommitCmdStr(s.message, s.flags)
			assert.Equal(t, s.expected, cmdStr)
		})
	}
}

// TestGitCommandCreateFixupCommit is a function.
func TestGitCommandCreateFixupCommit(t *testing.T) {
	type scenario struct {
		testName string
		sha      string
		command  func(string, ...string) *exec.Cmd
		test     func(error)
	}

	scenarios := []scenario{
		{
			"valid case",
			"12345",
			test.CreateMockCommand(t, []*test.CommandSwapper{
				{
					Expect:  `git commit --fixup=12345`,
					Replace: "echo",
				},
			}),
			func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	gitCmd := NewDummyGitCommand()

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCmd.OSCommand.Command = s.command
			s.test(gitCmd.CreateFixupCommit(s.sha))
		})
	}
}
