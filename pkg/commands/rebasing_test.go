package commands

import (
	"regexp"
	"testing"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestGitCommandRebaseBranch(t *testing.T) {
	type scenario struct {
		testName string
		arg      string
		runner   *oscommands.FakeCmdObjRunner
		test     func(error)
	}

	scenarios := []scenario{
		{
			testName: "successful rebase",
			arg:      "master",
			runner: oscommands.NewFakeRunner(t).
				Expect(`git rebase --interactive --autostash --keep-empty master`, "", nil),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			testName: "unsuccessful rebase",
			arg:      "master",
			runner: oscommands.NewFakeRunner(t).
				Expect(`git rebase --interactive --autostash --keep-empty master`, "", errors.New("error")),
			test: func(err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCmd := NewDummyGitCommandWithRunner(s.runner)
			s.test(gitCmd.RebaseBranch(s.arg))
		})
	}
}

// TestGitCommandSkipEditorCommand confirms that SkipEditorCommand injects
// environment variables that suppress an interactive editor
func TestGitCommandSkipEditorCommand(t *testing.T) {
	commandStr := "git blah"
	runner := oscommands.NewFakeRunner(t).ExpectFunc(func(cmdObj oscommands.ICmdObj) (string, error) {
		assert.Equal(t, commandStr, cmdObj.ToString())
		envVars := cmdObj.GetEnvVars()
		for _, regexStr := range []string{
			`^VISUAL=.*$`,
			`^EDITOR=.*$`,
			`^GIT_EDITOR=.*$`,
			"^LAZYGIT_CLIENT_COMMAND=EXIT_IMMEDIATELY$",
		} {
			foundMatch := utils.IncludesStringFunc(envVars, func(envVar string) bool {
				return regexp.MustCompile(regexStr).MatchString(envVar)
			})
			if !foundMatch {
				t.Errorf("expected environment variable %s to be set", regexStr)
			}
		}
		return "", nil
	})
	gitCmd := NewDummyGitCommandWithRunner(runner)
	err := gitCmd.runSkipEditorCommand(commandStr)
	assert.NoError(t, err)
	runner.CheckForMissingCalls()
}
