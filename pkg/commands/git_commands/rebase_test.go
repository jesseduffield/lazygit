package git_commands

import (
	"regexp"
	"strconv"
	"testing"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/app/daemon"
	"github.com/jesseduffield/lazygit/pkg/commands/git_config"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestRebaseRebaseBranch(t *testing.T) {
	type scenario struct {
		testName   string
		arg        string
		gitVersion *GitVersion
		runner     *oscommands.FakeCmdObjRunner
		test       func(error)
	}

	scenarios := []scenario{
		{
			testName:   "successful rebase",
			arg:        "master",
			gitVersion: &GitVersion{2, 26, 0, ""},
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"rebase", "--interactive", "--autostash", "--keep-empty", "--no-autosquash", "--rebase-merges", "master"}, "", nil),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			testName:   "unsuccessful rebase",
			arg:        "master",
			gitVersion: &GitVersion{2, 26, 0, ""},
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"rebase", "--interactive", "--autostash", "--keep-empty", "--no-autosquash", "--rebase-merges", "master"}, "", errors.New("error")),
			test: func(err error) {
				assert.Error(t, err)
			},
		},
		{
			testName:   "successful rebase (< 2.26.0)",
			arg:        "master",
			gitVersion: &GitVersion{2, 25, 5, ""},
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"rebase", "--interactive", "--autostash", "--keep-empty", "--no-autosquash", "--rebase-merges", "master"}, "", nil),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			testName:   "successful rebase (< 2.22.0)",
			arg:        "master",
			gitVersion: &GitVersion{2, 21, 9, ""},
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"rebase", "--interactive", "--autostash", "--keep-empty", "--no-autosquash", "master"}, "", nil),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			instance := buildRebaseCommands(commonDeps{runner: s.runner, gitVersion: s.gitVersion})
			s.test(instance.RebaseBranch(s.arg))
		})
	}
}

// TestRebaseSkipEditorCommand confirms that SkipEditorCommand injects
// environment variables that suppress an interactive editor
func TestRebaseSkipEditorCommand(t *testing.T) {
	cmdArgs := []string{"git", "blah"}
	runner := oscommands.NewFakeRunner(t).ExpectFunc("matches editor env var", func(cmdObj oscommands.ICmdObj) bool {
		assert.EqualValues(t, cmdArgs, cmdObj.Args())
		envVars := cmdObj.GetEnvVars()
		for _, regexStr := range []string{
			`^VISUAL=.*$`,
			`^EDITOR=.*$`,
			`^GIT_EDITOR=.*$`,
			`^GIT_SEQUENCE_EDITOR=.*$`,
			"^" + daemon.DaemonKindEnvKey + "=" + strconv.Itoa(int(daemon.DaemonKindExitImmediately)) + "$",
		} {
			foundMatch := lo.ContainsBy(envVars, func(envVar string) bool {
				return regexp.MustCompile(regexStr).MatchString(envVar)
			})
			if !foundMatch {
				return false
			}
		}
		return true
	}, "", nil)
	instance := buildRebaseCommands(commonDeps{runner: runner})
	err := instance.runSkipEditorCommand(instance.cmd.New(cmdArgs))
	assert.NoError(t, err)
	runner.CheckForMissingCalls()
}

func TestRebaseDiscardOldFileChanges(t *testing.T) {
	type scenario struct {
		testName               string
		gitConfigMockResponses map[string]string
		commits                []*models.Commit
		commitIndex            int
		fileName               []string
		runner                 *oscommands.FakeCmdObjRunner
		test                   func(error)
	}

	scenarios := []scenario{
		{
			testName:               "returns error when index outside of range of commits",
			gitConfigMockResponses: nil,
			commits:                []*models.Commit{},
			commitIndex:            0,
			fileName:               []string{"test999.txt"},
			runner:                 oscommands.NewFakeRunner(t),
			test: func(err error) {
				assert.Error(t, err)
			},
		},
		{
			testName:               "returns error when using gpg",
			gitConfigMockResponses: map[string]string{"commit.gpgsign": "true"},
			commits:                []*models.Commit{{Name: "commit", Hash: "123456"}},
			commitIndex:            0,
			fileName:               []string{"test999.txt"},
			runner:                 oscommands.NewFakeRunner(t),
			test: func(err error) {
				assert.Error(t, err)
			},
		},
		{
			testName:               "checks out file if it already existed",
			gitConfigMockResponses: nil,
			commits: []*models.Commit{
				{Name: "commit", Hash: "123456"},
				{Name: "commit2", Hash: "abcdef"},
			},
			commitIndex: 0,
			fileName:    []string{"test999.txt"},
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"rebase", "--interactive", "--autostash", "--keep-empty", "--no-autosquash", "--rebase-merges", "abcdef"}, "", nil).
				ExpectGitArgs([]string{"cat-file", "-e", "HEAD^:test999.txt"}, "", nil).
				ExpectGitArgs([]string{"checkout", "HEAD^", "--", "test999.txt"}, "", nil).
				ExpectGitArgs([]string{"commit", "--amend", "--no-edit", "--allow-empty"}, "", nil).
				ExpectGitArgs([]string{"rebase", "--continue"}, "", nil),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
		// test for when the file was created within the commit requires a refactor to support proper mocks
		// currently we'd need to mock out the os.Remove function and that's gonna introduce tech debt
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			instance := buildRebaseCommands(commonDeps{
				runner:     s.runner,
				gitVersion: &GitVersion{2, 26, 0, ""},
				gitConfig:  git_config.NewFakeGitConfig(s.gitConfigMockResponses),
			})

			s.test(instance.DiscardOldFileChanges(s.commits, s.commitIndex, s.fileName))
			s.runner.CheckForMissingCalls()
		})
	}
}
