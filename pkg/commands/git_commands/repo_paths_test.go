package git_commands

import (
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

type (
	argFn func() []string
	errFn func(getRevParseArgs argFn) error
)

type Scenario struct {
	Name       string
	BeforeFunc func(runner *oscommands.FakeCmdObjRunner, getRevParseArgs argFn)
	Path       string
	Expected   *RepoPaths
	Err        errFn
}

// primaryRevParseArgs matches the first GetRepoPathsForDir rev-parse (path flags only).
func primaryRevParseArgs(getRevParseArgs argFn) []string {
	return append(getRevParseArgs(),
		"--show-toplevel", "--absolute-git-dir", "--git-common-dir", "--is-bare-repository")
}

// superprojectRevParseArgs matches the second rev-parse (--show-superproject-working-tree alone).
func superprojectRevParseArgs(getRevParseArgs argFn) []string {
	return append(getRevParseArgs(), "--show-superproject-working-tree")
}

func TestGetRepoPaths(t *testing.T) {
	scenarios := []Scenario{
		{
			Name: "typical case",
			BeforeFunc: func(runner *oscommands.FakeCmdObjRunner, getRevParseArgs argFn) {
				// setup for main worktree
				mockOutput := lo.Ternary(runtime.GOOS == "windows", []string{
					// --show-toplevel
					`C:\path\to\repo`,
					// --git-dir
					`C:\path\to\repo\.git`,
					// --git-common-dir
					`C:\path\to\repo\.git`,
					// --is-bare-repository
					"false",
				}, []string{
					// --show-toplevel
					"/path/to/repo",
					// --git-dir
					"/path/to/repo/.git",
					// --git-common-dir
					"/path/to/repo/.git",
					// --is-bare-repository
					"false",
				})
				runner.ExpectGitArgs(
					primaryRevParseArgs(getRevParseArgs),
					strings.Join(mockOutput, "\n"),
					nil)
				// --show-superproject-working-tree (empty: not inside a submodule checkout)
				runner.ExpectGitArgs(
					superprojectRevParseArgs(getRevParseArgs),
					"",
					nil)
			},
			Path: "/path/to/repo",
			Expected: lo.Ternary(runtime.GOOS == "windows", &RepoPaths{
				worktreePath:       `C:\path\to\repo`,
				worktreeGitDirPath: `C:\path\to\repo\.git`,
				repoPath:           `C:\path\to\repo`,
				repoGitDirPath:     `C:\path\to\repo\.git`,
				repoName:           `repo`,
				isBareRepo:         false,
			}, &RepoPaths{
				worktreePath:       "/path/to/repo",
				worktreeGitDirPath: "/path/to/repo/.git",
				repoPath:           "/path/to/repo",
				repoGitDirPath:     "/path/to/repo/.git",
				repoName:           "repo",
				isBareRepo:         false,
			}),
			Err: nil,
		},
		{
			Name: "bare repo",
			BeforeFunc: func(runner *oscommands.FakeCmdObjRunner, getRevParseArgs argFn) {
				// setup for main worktree with a separate bare git dir
				mockOutput := lo.Ternary(runtime.GOOS == "windows", []string{
					// --show-toplevel
					`C:\path\to\repo`,
					// --git-dir
					`C:\path\to\bare_repo\bare.git`,
					// --git-common-dir
					`C:\path\to\bare_repo\bare.git`,
					// --is-bare-repository
					`true`,
				}, []string{
					// --show-toplevel
					"/path/to/repo",
					// --git-dir
					"/path/to/bare_repo/bare.git",
					// --git-common-dir
					"/path/to/bare_repo/bare.git",
					// --is-bare-repository
					"true",
				})
				runner.ExpectGitArgs(
					primaryRevParseArgs(getRevParseArgs),
					strings.Join(mockOutput, "\n"),
					nil)
				// --show-superproject-working-tree (empty)
				runner.ExpectGitArgs(
					superprojectRevParseArgs(getRevParseArgs),
					"",
					nil)
			},
			Path: "/path/to/repo",
			Expected: lo.Ternary(runtime.GOOS == "windows", &RepoPaths{
				worktreePath:       `C:\path\to\repo`,
				worktreeGitDirPath: `C:\path\to\bare_repo\bare.git`,
				repoPath:           `C:\path\to\bare_repo`,
				repoGitDirPath:     `C:\path\to\bare_repo\bare.git`,
				repoName:           `bare_repo`,
				isBareRepo:         true,
			}, &RepoPaths{
				worktreePath:       "/path/to/repo",
				worktreeGitDirPath: "/path/to/bare_repo/bare.git",
				repoPath:           "/path/to/bare_repo",
				repoGitDirPath:     "/path/to/bare_repo/bare.git",
				repoName:           "bare_repo",
				isBareRepo:         true,
			}),
			Err: nil,
		},
		{
			Name: "submodule",
			BeforeFunc: func(runner *oscommands.FakeCmdObjRunner, getRevParseArgs argFn) {
				mockPrimary := lo.Ternary(runtime.GOOS == "windows", []string{
					// --show-toplevel
					`C:\path\to\repo\submodule1`,
					// --git-dir
					`C:\path\to\repo\.git\modules\submodule1`,
					// --git-common-dir
					`C:\path\to\repo\.git\modules\submodule1`,
					// --is-bare-repository
					`false`,
				}, []string{
					// --show-toplevel
					"/path/to/repo/submodule1",
					// --git-dir
					"/path/to/repo/.git/modules/submodule1",
					// --git-common-dir
					"/path/to/repo/.git/modules/submodule1",
					// --is-bare-repository
					"false",
				})
				// --show-superproject-working-tree (superproject worktree path)
				superOut := lo.Ternary(runtime.GOOS == "windows", `C:\path\to\repo`, "/path/to/repo")
				runner.ExpectGitArgs(
					primaryRevParseArgs(getRevParseArgs),
					strings.Join(mockPrimary, "\n"),
					nil)
				runner.ExpectGitArgs(
					superprojectRevParseArgs(getRevParseArgs),
					superOut,
					nil)
			},
			Path: "/path/to/repo/submodule1",
			Expected: lo.Ternary(runtime.GOOS == "windows", &RepoPaths{
				worktreePath:       `C:\path\to\repo\submodule1`,
				worktreeGitDirPath: `C:\path\to\repo\.git\modules\submodule1`,
				repoPath:           `C:\path\to\repo\submodule1`,
				repoGitDirPath:     `C:\path\to\repo\.git\modules\submodule1`,
				repoName:           `submodule1`,
				isBareRepo:         false,
			}, &RepoPaths{
				worktreePath:       "/path/to/repo/submodule1",
				worktreeGitDirPath: "/path/to/repo/.git/modules/submodule1",
				repoPath:           "/path/to/repo/submodule1",
				repoGitDirPath:     "/path/to/repo/.git/modules/submodule1",
				repoName:           "submodule1",
				isBareRepo:         false,
			}),
			Err: nil,
		},
		{
			Name: "superproject rev-parse fails (fallback to non-submodule repoPath)",
			BeforeFunc: func(runner *oscommands.FakeCmdObjRunner, getRevParseArgs argFn) {
				// Primary rev-parse succeeds (e.g. repo-tool symlinked .git); secondary can still error
				// (e.g. BUG: submodule.c) — we ignore superproject failure and use repoPath from common-dir.
				mockOutput := lo.Ternary(runtime.GOOS == "windows", []string{
					// --show-toplevel
					`C:\path\to\repo`,
					// --absolute-git-dir
					`C:\path\to\repo\.git`,
					// --git-common-dir
					`C:\path\to\repo\.git`,
					// --is-bare-repository
					"false",
				}, []string{
					// --show-toplevel
					"/path/to/repo",
					// --absolute-git-dir
					"/path/to/repo/.git",
					// --git-common-dir
					"/path/to/repo/.git",
					// --is-bare-repository
					"false",
				})
				runner.ExpectGitArgs(
					primaryRevParseArgs(getRevParseArgs),
					strings.Join(mockOutput, "\n"),
					nil)
				// --show-superproject-working-tree (Git errors, e.g. submodule.c internal BUG)
				runner.ExpectGitArgs(
					superprojectRevParseArgs(getRevParseArgs),
					"",
					errors.New("BUG: submodule.c:2455: returned path string doesn't match cwd?"))
			},
			Path: "/path/to/repo",
			Expected: lo.Ternary(runtime.GOOS == "windows", &RepoPaths{
				worktreePath:       `C:\path\to\repo`,
				worktreeGitDirPath: `C:\path\to\repo\.git`,
				repoPath:           `C:\path\to\repo`,
				repoGitDirPath:     `C:\path\to\repo\.git`,
				repoName:           `repo`,
				isBareRepo:         false,
			}, &RepoPaths{
				worktreePath:       "/path/to/repo",
				worktreeGitDirPath: "/path/to/repo/.git",
				repoPath:           "/path/to/repo",
				repoGitDirPath:     "/path/to/repo/.git",
				repoName:           "repo",
				isBareRepo:         false,
			}),
			Err: nil,
		},
		{
			Name: "git rev-parse returns an error",
			BeforeFunc: func(runner *oscommands.FakeCmdObjRunner, getRevParseArgs argFn) {
				// Primary rev-parse fails; superproject call is never run
				runner.ExpectGitArgs(
					primaryRevParseArgs(getRevParseArgs),
					"",
					errors.New("fatal: invalid gitfile format: /path/to/repo/worktree2/.git"))
			},
			Path:     "/path/to/repo/worktree2",
			Expected: nil,
			Err: func(getRevParseArgs argFn) error {
				args := strings.Join(getRevParseArgs(), " ")
				return fmt.Errorf("'git %v --show-toplevel --absolute-git-dir --git-common-dir --is-bare-repository' failed: fatal: invalid gitfile format: /path/to/repo/worktree2/.git", args)
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.Name, func(t *testing.T) {
			runner := oscommands.NewFakeRunner(t)
			cmd := oscommands.NewDummyCmdObjBuilder(runner)

			getRevParseArgs := func() []string {
				return []string{"rev-parse", "--path-format=absolute"}
			}
			// prepare the filesystem for the scenario
			s.BeforeFunc(runner, getRevParseArgs)

			repoPaths, err := GetRepoPathsForDir("", cmd)

			// check the error and the paths
			if s.Err != nil {
				scenarioErr := s.Err(getRevParseArgs)
				assert.Error(t, err)
				assert.EqualError(t, err, scenarioErr.Error())
			} else {
				assert.Nil(t, err)
				assert.Equal(t, s.Expected, repoPaths)
			}
		})
	}
}
