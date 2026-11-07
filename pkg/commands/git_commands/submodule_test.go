package git_commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/env"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestIsSuspiciousSubmoduleName(t *testing.T) {
	scenarios := []struct {
		name       string
		suspicious bool
	}{
		{name: "sub", suspicious: false},
		{name: "path/to/sub", suspicious: false},
		{name: "path\\to\\sub", suspicious: false},
		{name: ".", suspicious: false},
		{name: "...", suspicious: false},
		{name: "..sub", suspicious: false},
		{name: "sub..", suspicious: false},
		{name: "a..b/c", suspicious: false},
		{name: "", suspicious: true},
		{name: "..", suspicious: true},
		{name: "../victim", suspicious: true},
		{name: "../../../victim", suspicious: true},
		{name: "sub/..", suspicious: true},
		{name: "sub/../../victim", suspicious: true},
		{name: "..\\victim", suspicious: true},
		{name: "sub\\..\\..\\victim", suspicious: true},
		{name: "sub/..\\..\\victim", suspicious: true},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			assert.Equal(t, s.suspicious, isSuspiciousSubmoduleName(s.name))
		})
	}
}

func TestSubmoduleDeleteOutsideModulesDir(t *testing.T) {
	scenarios := []struct {
		testName string
		name     string
		// relative to the directory that contains the repo
		victimPath string
	}{
		{
			testName:   "modules dir itself",
			name:       ".",
			victimPath: "repo/.git/modules/other_submodule",
		},
		{
			testName:   "outside the repo",
			name:       "../../../victim",
			victimPath: "victim",
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			dir := t.TempDir()
			victimFile := filepath.Join(dir, s.victimPath, "file")
			assert.NoError(t, os.MkdirAll(filepath.Dir(victimFile), 0o755))
			assert.NoError(t, os.WriteFile(victimFile, []byte("content"), 0o644))

			runner := oscommands.NewFakeRunner(t)
			instance := buildSubmoduleCommands(commonDeps{
				runner:    runner,
				repoPaths: MockRepoPaths(filepath.Join(dir, "repo")),
			})

			err := instance.Delete(&models.SubmoduleConfig{Name: s.name, Path: "sub"})
			assert.Error(t, err)
			assert.FileExists(t, victimFile)
			runner.CheckForMissingCalls()
		})
	}
}

func TestSubmoduleGetConflictCommits(t *testing.T) {
	type scenario struct {
		testName       string
		output         string
		expectedBase   string
		expectedOurs   string
		expectedTheirs string
	}

	scenarios := []scenario{
		{
			testName:       "all three stages present (both modified)",
			output:         "160000 aaaaaaa 1\tmysub\x00160000 bbbbbbb 2\tmysub\x00160000 ccccccc 3\tmysub\x00",
			expectedBase:   "aaaaaaa",
			expectedOurs:   "bbbbbbb",
			expectedTheirs: "ccccccc",
		},
		{
			testName:       "only our and their stages (added on both sides)",
			output:         "160000 bbbbbbb 2\tmysub\x00160000 ccccccc 3\tmysub\x00",
			expectedBase:   "",
			expectedOurs:   "bbbbbbb",
			expectedTheirs: "ccccccc",
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			runner := oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"ls-files", "-u", "-z", "--", "mysub"}, s.output, nil)
			instance := buildSubmoduleCommands(commonDeps{runner: runner})

			base, ours, theirs, err := instance.GetConflictCommits("mysub")
			assert.NoError(t, err)
			assert.Equal(t, s.expectedBase, base)
			assert.Equal(t, s.expectedOurs, ours)
			assert.Equal(t, s.expectedTheirs, theirs)
			runner.CheckForMissingCalls()
		})
	}
}

func TestSubmoduleGetConflictCommitsError(t *testing.T) {
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs([]string{"ls-files", "-u", "-z", "--", "mysub"}, "", errors.New("error"))
	instance := buildSubmoduleCommands(commonDeps{runner: runner})

	_, _, _, err := instance.GetConflictCommits("mysub")
	assert.Error(t, err)
	runner.CheckForMissingCalls()
}

func TestSubmoduleGetCommitSummary(t *testing.T) {
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs([]string{"-c", "log.showsignature=false", "-C", "mysub", "log", "--format=%h %s", "--max-count=1", "bbbbbbb"}, "bbbbbbb the subject\n", nil)
	instance := buildSubmoduleCommands(commonDeps{runner: runner})

	summary, err := instance.GetCommitSummary("mysub", "bbbbbbb")
	assert.NoError(t, err)
	assert.Equal(t, "bbbbbbb the subject", summary)
	runner.CheckForMissingCalls()
}

func TestSubmoduleCheckoutConflictCommit(t *testing.T) {
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs([]string{"-C", "mysub", "checkout", "bbbbbbb"}, "", nil)
	instance := buildSubmoduleCommands(commonDeps{runner: runner})

	assert.NoError(t, instance.CheckoutConflictCommit("mysub", "bbbbbbb"))
	runner.CheckForMissingCalls()
}

// A command that runs inside a submodule mustn't inherit the GIT_DIR and
// GIT_WORK_TREE that say where the superproject is; git would answer it from
// there instead, and the answer would look perfectly plausible.
func TestSubmoduleCommandDoesntUseOurGitLocation(t *testing.T) {
	t.Setenv(env.GitDirEnvVar, "/path/to/repo/.git")
	t.Setenv(env.GitWorkTreeEnvVar, "/path/to/repo")

	runner := oscommands.NewFakeRunner(t).
		ExpectFunc("has neither GIT_DIR nor GIT_WORK_TREE", func(cmdObj *oscommands.CmdObj) bool {
			return lo.NoneBy(cmdObj.GetEnvVars(), func(envVar string) bool {
				return strings.HasPrefix(envVar, env.GitDirEnvVar+"=") ||
					strings.HasPrefix(envVar, env.GitWorkTreeEnvVar+"=")
			})
		}, "bbbbbbb the subject\n", nil)
	instance := buildSubmoduleCommands(commonDeps{runner: runner})

	_, err := instance.GetCommitSummary("mysub", "bbbbbbb")
	assert.NoError(t, err)
	runner.CheckForMissingCalls()
}

func TestSubmoduleConflictSideLog(t *testing.T) {
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs([]string{"-C", "mysub", "log", "--oneline", "--color=always", "ccccccc..bbbbbbb"}, "bbbbbbb left\n", nil)
	instance := buildSubmoduleCommands(commonDeps{runner: runner})

	output, err := instance.ConflictSideLog("mysub", "bbbbbbb", "ccccccc")
	assert.NoError(t, err)
	assert.Equal(t, "bbbbbbb left\n", output)
	runner.CheckForMissingCalls()
}
