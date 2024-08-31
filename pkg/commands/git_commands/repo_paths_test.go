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
					// --show-superproject-working-tree
				}, []string{
					// --show-toplevel
					"/path/to/repo",
					// --git-dir
					"/path/to/repo/.git",
					// --git-common-dir
					"/path/to/repo/.git",
					// --is-bare-repository
					"false",
					// --show-superproject-working-tree
				})
				runner.ExpectGitArgs(
					append(getRevParseArgs(), "--show-toplevel", "--absolute-git-dir", "--git-common-dir", "--is-bare-repository", "--show-superproject-working-tree"),
					strings.Join(mockOutput, "\n"),
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
				// setup for main worktree
				mockOutput := lo.Ternary(runtime.GOOS == "windows", []string{
					// --show-toplevel
					`C:\path\to\repo`,
					// --git-dir
					`C:\path\to\bare_repo\bare.git`,
					// --git-common-dir
					`C:\path\to\bare_repo\bare.git`,
					// --is-bare-repository
					`true`,
					// --show-superproject-working-tree
				}, []string{
					// --show-toplevel
					"/path/to/repo",
					// --git-dir
					"/path/to/bare_repo/bare.git",
					// --git-common-dir
					"/path/to/bare_repo/bare.git",
					// --is-bare-repository
					"true",
					// --show-superproject-working-tree
				})
				runner.ExpectGitArgs(
					append(getRevParseArgs(), "--show-toplevel", "--absolute-git-dir", "--git-common-dir", "--is-bare-repository", "--show-superproject-working-tree"),
					strings.Join(mockOutput, "\n"),
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
				mockOutput := lo.Ternary(runtime.GOOS == "windows", []string{
					// --show-toplevel
					`C:\path\to\repo\submodule1`,
					// --git-dir
					`C:\path\to\repo\.git\modules\submodule1`,
					// --git-common-dir
					`C:\path\to\repo\.git\modules\submodule1`,
					// --is-bare-repository
					`false`,
					// --show-superproject-working-tree
					`C:\path\to\repo`,
				}, []string{
					// --show-toplevel
					"/path/to/repo/submodule1",
					// --git-dir
					"/path/to/repo/.git/modules/submodule1",
					// --git-common-dir
					"/path/to/repo/.git/modules/submodule1",
					// --is-bare-repository
					"false",
					// --show-superproject-working-tree
					"/path/to/repo",
				})
				runner.ExpectGitArgs(
					append(getRevParseArgs(), "--show-toplevel", "--absolute-git-dir", "--git-common-dir", "--is-bare-repository", "--show-superproject-working-tree"),
					strings.Join(mockOutput, "\n"),
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
			Name: "git rev-parse returns an error",
			BeforeFunc: func(runner *oscommands.FakeCmdObjRunner, getRevParseArgs argFn) {
				runner.ExpectGitArgs(
					append(getRevParseArgs(), "--show-toplevel", "--absolute-git-dir", "--git-common-dir", "--is-bare-repository", "--show-superproject-working-tree"),
					"",
					errors.New("fatal: invalid gitfile format: /path/to/repo/worktree2/.git"))
			},
			Path:     "/path/to/repo/worktree2",
			Expected: nil,
			Err: func(getRevParseArgs argFn) error {
				args := strings.Join(getRevParseArgs(), " ")
				return errors.New(
					fmt.Sprintf("'git %v --show-toplevel --absolute-git-dir --git-common-dir --is-bare-repository --show-superproject-working-tree' failed: fatal: invalid gitfile format: /path/to/repo/worktree2/.git", args),
				)
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.Name, func(t *testing.T) {
			runner := oscommands.NewFakeRunner(t)
			cmd := oscommands.NewDummyCmdObjBuilder(runner)

			version, err := GetGitVersion(oscommands.NewDummyOSCommand())
			if err != nil {
				t.Fatal(err)
			}

			getRevParseArgs := func() []string {
				args := []string{"rev-parse"}
				if version.IsAtLeast(2, 31, 0) {
					args = append(args, "--path-format=absolute")
				}
				return args
			}
			// prepare the filesystem for the scenario
			s.BeforeFunc(runner, getRevParseArgs)

			repoPaths, err := GetRepoPathsForDir("", cmd, version)

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
