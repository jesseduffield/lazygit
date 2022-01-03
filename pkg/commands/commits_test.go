package commands

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/stretchr/testify/assert"
)

func TestGitCommandRenameCommit(t *testing.T) {
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs([]string{"commit", "--allow-empty", "--amend", "--only", "-m", "test"}, "", nil)
	gitCmd := NewDummyGitCommandWithRunner(runner)

	assert.NoError(t, gitCmd.RenameCommit("test"))
	runner.CheckForMissingCalls()
}

func TestGitCommandResetToCommit(t *testing.T) {
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs([]string{"reset", "--hard", "78976bc"}, "", nil)
	gitCmd := NewDummyGitCommandWithRunner(runner)

	assert.NoError(t, gitCmd.ResetToCommit("78976bc", "hard", []string{}))
	runner.CheckForMissingCalls()
}

func TestGitCommandCommitObj(t *testing.T) {
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
			expected: `git commit -m "test"`,
		},
		{
			testName: "Commit with --no-verify flag",
			message:  "test",
			flags:    "--no-verify",
			expected: `git commit --no-verify -m "test"`,
		},
		{
			testName: "Commit with multiline message",
			message:  "line1\nline2",
			flags:    "",
			expected: `git commit -m "line1" -m "line2"`,
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCmd := NewDummyGitCommand()
			cmdStr := gitCmd.CommitCmdObj(s.message, s.flags).ToString()
			assert.Equal(t, s.expected, cmdStr)
		})
	}
}

func TestGitCommandCreateFixupCommit(t *testing.T) {
	type scenario struct {
		testName string
		sha      string
		runner   *oscommands.FakeCmdObjRunner
		test     func(error)
	}

	scenarios := []scenario{
		{
			testName: "valid case",
			sha:      "12345",
			runner: oscommands.NewFakeRunner(t).
				Expect(`git commit --fixup=12345`, "", nil),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCmd := NewDummyGitCommandWithRunner(s.runner)
			s.test(gitCmd.CreateFixupCommit(s.sha))
			s.runner.CheckForMissingCalls()
		})
	}
}

func TestGitCommandShowCmdObj(t *testing.T) {
	type scenario struct {
		testName    string
		filterPath  string
		contextSize int
		expected    string
	}

	scenarios := []scenario{
		{
			testName:    "Default case without filter path",
			filterPath:  "",
			contextSize: 3,
			expected:    "git show --submodule --color=always --unified=3 --no-renames --stat -p 1234567890 ",
		},
		{
			testName:    "Default case with filter path",
			filterPath:  "file.txt",
			contextSize: 3,
			expected:    `git show --submodule --color=always --unified=3 --no-renames --stat -p 1234567890  -- "file.txt"`,
		},
		{
			testName:    "Show diff with custom context size",
			filterPath:  "",
			contextSize: 77,
			expected:    "git show --submodule --color=always --unified=77 --no-renames --stat -p 1234567890 ",
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCmd := NewDummyGitCommand()
			gitCmd.UserConfig.Git.DiffContextSize = s.contextSize
			cmdStr := gitCmd.ShowCmdObj("1234567890", s.filterPath).ToString()
			assert.Equal(t, s.expected, cmdStr)
		})
	}
}
