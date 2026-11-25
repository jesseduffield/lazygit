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
	"github.com/jesseduffield/lazygit/pkg/utils"
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
	runner := oscommands.NewFakeRunner(t).ExpectFunc("matches editor env var", func(cmdObj *oscommands.CmdObj) bool {
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
	err := skipEditor(instance.cmd.New(cmdArgs)).Run()
	assert.NoError(t, err)
	runner.CheckForMissingCalls()
}

func TestRebaseDiscardOldFileChanges(t *testing.T) {
	type scenario struct {
		testName               string
		gitConfigMockResponses map[string]string
		commitOpts             []models.NewCommitOpts
		commitIndex            int
		fileName               []string
		runner                 *oscommands.FakeCmdObjRunner
		test                   func(error)
	}

	scenarios := []scenario{
		{
			testName:               "returns error when index outside of range of commits",
			gitConfigMockResponses: nil,
			commitOpts:             []models.NewCommitOpts{},
			commitIndex:            0,
			fileName:               []string{"test999.txt"},
			runner:                 oscommands.NewFakeRunner(t),
			test: func(err error) {
				assert.Error(t, err)
			},
		},
		{
			testName:               "returns error when using gpg",
			gitConfigMockResponses: map[string]string{"commit.gpgSign": "true"},
			commitOpts:             []models.NewCommitOpts{{Name: "commit", Hash: "123456"}},
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
			commitOpts: []models.NewCommitOpts{
				{Name: "commit", Hash: "123456"},
				{Name: "commit2", Hash: "abcdef"},
			},
			commitIndex: 0,
			fileName:    []string{"test999.txt"},
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"rebase", "--interactive", "--autostash", "--keep-empty", "--no-autosquash", "--rebase-merges", "abcdef"}, "", nil).
				ExpectGitArgs([]string{"ls-tree", "--name-only", "HEAD^", "--", "test999.txt"}, "test999.txt\n", nil).
				ExpectGitArgs([]string{"checkout", "HEAD^", "--", "test999.txt"}, "", nil).
				ExpectGitArgs([]string{"commit", "--amend", "--no-edit", "--allow-empty", "--allow-empty-message"}, "", nil).
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

			hashPool := &utils.StringPool{}
			commits := lo.Map(s.commitOpts,
				func(opts models.NewCommitOpts, _ int) *models.Commit { return models.NewCommit(hashPool, opts) })

			s.test(instance.DiscardOldFileChanges(commits, s.commitIndex, s.fileName))
			s.runner.CheckForMissingCalls()
		})
	}
}

func TestGetBaseHashOrRoot(t *testing.T) {
	hashPool := &utils.StringPool{}
	newCommit := func(hash string, offBranch bool, parents ...string) *models.Commit {
		return models.NewCommit(hashPool, models.NewCommitOpts{Hash: hash, OffBranch: offBranch, Parents: parents})
	}

	t.Run("the commit below is loaded", func(t *testing.T) {
		commits := []*models.Commit{newCommit("b", false, "a"), newCommit("a", false)}
		assert.Equal(t, "a", getBaseHashOrRoot(commits, 0, 1))
	})

	t.Run("the commit below is the parent of the root commit", func(t *testing.T) {
		commits := []*models.Commit{newCommit("b", false, "a"), newCommit("a", false)}
		assert.Equal(t, "--root", getBaseHashOrRoot(commits, 1, 1))
	})

	t.Run("the commit below hasn't been loaded", func(t *testing.T) {
		commits := []*models.Commit{newCommit("c", false, "b"), newCommit("b", false, "a")}
		/* EXPECTED:
		assert.Equal(t, "b~1", getBaseHashOrRoot(commits, 0, 2))
		ACTUAL: */
		assert.Equal(t, "--root", getBaseHashOrRoot(commits, 0, 2))
	})

	t.Run("only commits of other branches are loaded below", func(t *testing.T) {
		commits := []*models.Commit{newCommit("b", false, "a"), newCommit("x", true, "a")}
		/* EXPECTED:
		assert.Equal(t, "b~1", getBaseHashOrRoot(commits, 0, 1))
		ACTUAL: */
		assert.Equal(t, "--root", getBaseHashOrRoot(commits, 0, 1))
	})
}
