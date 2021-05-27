package commands

import (
	"os/exec"
	"regexp"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/test"
	"github.com/stretchr/testify/assert"
)

// TestGitCommandRebaseBranch is a function.
func TestGitCommandRebaseBranch(t *testing.T) {
	type scenario struct {
		testName string
		arg      string
		command  func(string, ...string) *exec.Cmd
		test     func(error)
	}

	scenarios := []scenario{
		{
			"successful rebase",
			"master",
			test.CreateMockCommand(t, []*test.CommandSwapper{
				{
					Expect:  "git rebase --interactive --autostash --keep-empty master",
					Replace: "echo",
				},
			}),
			func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			"unsuccessful rebase",
			"master",
			test.CreateMockCommand(t, []*test.CommandSwapper{
				{
					Expect:  "git rebase --interactive --autostash --keep-empty master",
					Replace: "test",
				},
			}),
			func(err error) {
				assert.Error(t, err)
			},
		},
	}

	gitCmd := NewDummyGitCommand()

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCmd.OSCommand.Command = s.command
			s.test(gitCmd.RebaseBranch(s.arg))
		})
	}
}

// TestGitCommandSkipEditorCommand confirms that SkipEditorCommand injects
// environment variables that suppress an interactive editor
func TestGitCommandSkipEditorCommand(t *testing.T) {
	cmd := NewDummyGitCommand()

	cmd.OSCommand.SetBeforeExecuteCmd(func(cmd *exec.Cmd) {
		test.AssertContainsMatch(
			t,
			cmd.Env,
			regexp.MustCompile("^VISUAL="),
			"expected VISUAL to be set for a non-interactive external command",
		)

		test.AssertContainsMatch(
			t,
			cmd.Env,
			regexp.MustCompile("^EDITOR="),
			"expected EDITOR to be set for a non-interactive external command",
		)

		test.AssertContainsMatch(
			t,
			cmd.Env,
			regexp.MustCompile("^GIT_EDITOR="),
			"expected GIT_EDITOR to be set for a non-interactive external command",
		)

		test.AssertContainsMatch(
			t,
			cmd.Env,
			regexp.MustCompile("^LAZYGIT_CLIENT_COMMAND=EXIT_IMMEDIATELY$"),
			"expected LAZYGIT_CLIENT_COMMAND to be set for a non-interactive external command",
		)
	})

	_ = cmd.runSkipEditorCommand("true")
}
