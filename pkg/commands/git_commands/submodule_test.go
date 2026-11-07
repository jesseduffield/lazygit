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

func TestParseGitmodules(t *testing.T) {
	scenarios := []struct {
		testName string
		lines    []string
		expected []*models.SubmoduleConfig
	}{
		{
			testName: "two submodules",
			lines: []string{
				`[submodule "one"]`,
				"\tpath = path/to/one",
				"\turl = https://example.com/one.git",
				`[submodule "two"]`,
				"\tpath = two",
				"\turl = ../two",
			},
			expected: []*models.SubmoduleConfig{
				{Name: "one", Path: "path/to/one", Url: "https://example.com/one.git"},
				{Name: "two", Path: "two", Url: "../two"},
			},
		},
		{
			testName: "url that contains path=",
			lines: []string{
				`[submodule "sub"]`,
				"\tpath = sub",
				"\turl = https://example.com/repo?path=other",
			},
			expected: []*models.SubmoduleConfig{
				{Name: "sub", Path: "sub", Url: "https://example.com/repo?path=other"},
			},
		},
		{
			testName: "commented-out section",
			lines: []string{
				`[submodule "sub"]`,
				"\tpath = sub",
				"\turl = ../sub",
				`#[submodule "old"]`,
				"#\tpath = old",
				"#\turl = ../old",
			},
			expected: []*models.SubmoduleConfig{
				{Name: "sub", Path: "sub", Url: "../sub"},
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			configs, err := parseGitmodules(strings.NewReader(strings.Join(s.lines, "\n")), nil)
			assert.NoError(t, err)
			assert.Equal(t, s.expected, configs)
		})
	}
}

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

func TestIsSuspiciousSubmodulePath(t *testing.T) {
	dir := t.TempDir()
	parentDir := filepath.Join(dir, "repo")
	assert.NoError(t, os.MkdirAll(filepath.Join(parentDir, "sub"), 0o755))
	assert.NoError(t, os.MkdirAll(filepath.Join(parentDir, "path", "to", "sub"), 0o755))
	assert.NoError(t, os.MkdirAll(filepath.Join(dir, "other"), 0o755))
	assert.NoError(t, os.WriteFile(filepath.Join(parentDir, "file"), []byte("content"), 0o644))

	scenarios := []struct {
		path       string
		suspicious bool
	}{
		{path: "sub", suspicious: false},
		{path: "sub/", suspicious: false},
		{path: "path/to/sub", suspicious: false},
		{path: "not_checked_out", suspicious: false},
		{path: "path/not_checked_out/sub", suspicious: false},
		{path: "", suspicious: true},
		{path: ".", suspicious: true},
		{path: "./", suspicious: true},
		{path: "sub/..", suspicious: true},
		{path: "..", suspicious: true},
		{path: "../other", suspicious: true},
		{path: "sub/../../other", suspicious: true},
		{path: filepath.Join(dir, "other"), suspicious: true},
		{path: "file", suspicious: true},
		{path: "file/sub", suspicious: true},
	}

	for _, s := range scenarios {
		t.Run(s.path, func(t *testing.T) {
			assert.Equal(t, s.suspicious, isSuspiciousSubmodulePath(parentDir, s.path))
		})
	}
}

func TestIsSuspiciousSubmodulePathWithSymlinks(t *testing.T) {
	dir := t.TempDir()
	parentDir := filepath.Join(dir, "repo")
	assert.NoError(t, os.MkdirAll(filepath.Join(parentDir, "sub"), 0o755))
	assert.NoError(t, os.MkdirAll(filepath.Join(dir, "other", "sub"), 0o755))
	if err := os.Symlink(filepath.Join("..", "other"), filepath.Join(parentDir, "link_outside")); err != nil {
		t.Skipf("Can't create symlinks: %v", err)
	}
	assert.NoError(t, os.Symlink("sub", filepath.Join(parentDir, "link_inside")))

	scenarios := []struct {
		path       string
		suspicious bool
	}{
		{path: "sub", suspicious: false},
		{path: "link_outside", suspicious: true},
		{path: "link_outside/sub", suspicious: true},
		{path: "link_inside", suspicious: true},
	}

	for _, s := range scenarios {
		t.Run(s.path, func(t *testing.T) {
			assert.Equal(t, s.suspicious, isSuspiciousSubmodulePath(parentDir, s.path))
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
