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
					// --show-superproject-working-tree
				}, []string{
					// --show-toplevel
					"/path/to/repo",
					// --git-dir
					"/path/to/repo/.git",
					// --git-common-dir
					"/path/to/repo/.git",
					// --show-superproject-working-tree
				})
				runner.ExpectGitArgs(
					append(getRevParseArgs(), "--show-toplevel", "--absolute-git-dir", "--git-common-dir", "--show-superproject-working-tree"),
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
			// git refuses to answer --show-toplevel when there's no work tree, so
			// we have to ask a second time without it.
			Name: "bare repo",
			BeforeFunc: func(runner *oscommands.FakeCmdObjRunner, getRevParseArgs argFn) {
				runner.ExpectGitArgs(
					append(getRevParseArgs(), "--show-toplevel", "--absolute-git-dir", "--git-common-dir", "--show-superproject-working-tree"),
					"",
					errors.New("fatal: this operation must be run in a work tree"))

				mockOutput := lo.Ternary(runtime.GOOS == "windows", []string{
					// --git-dir
					`C:\path\to\project\bare.git`,
					// --git-common-dir
					`C:\path\to\project\bare.git`,
				}, []string{
					// --git-dir
					"/path/to/project/bare.git",
					// --git-common-dir
					"/path/to/project/bare.git",
				})
				runner.ExpectGitArgs(
					append(getRevParseArgs(), "--absolute-git-dir", "--git-common-dir"),
					strings.Join(mockOutput, "\n"),
					nil)
			},
			Path: "/path/to/project",
			Expected: lo.Ternary(runtime.GOOS == "windows", &RepoPaths{
				worktreePath:       "",
				worktreeGitDirPath: `C:\path\to\project\bare.git`,
				repoPath:           `C:\path\to\project`,
				repoGitDirPath:     `C:\path\to\project\bare.git`,
				repoName:           `project`,
				isBareRepo:         true,
			}, &RepoPaths{
				worktreePath:       "",
				worktreeGitDirPath: "/path/to/project/bare.git",
				repoPath:           "/path/to/project",
				repoGitDirPath:     "/path/to/project/bare.git",
				repoName:           "project",
				isBareRepo:         true,
			}),
			Err: nil,
		},
		{
			// Standing in the .git dir of an ordinary repo: git refuses to name a
			// work tree, but the directory holding the .git is one, so we open the
			// repo from there.
			Name: "in a repo's .git dir",
			BeforeFunc: func(runner *oscommands.FakeCmdObjRunner, getRevParseArgs argFn) {
				gitDir := lo.Ternary(runtime.GOOS == "windows", `C:\path\to\repo\.git`, "/path/to/repo/.git")
				worktree := lo.Ternary(runtime.GOOS == "windows", `C:\path\to\repo`, "/path/to/repo")

				runner.ExpectGitArgs(
					append(getRevParseArgs(), "--show-toplevel", "--absolute-git-dir", "--git-common-dir", "--show-superproject-working-tree"),
					"",
					errors.New("fatal: this operation must be run in a work tree"))
				runner.ExpectGitArgs(
					append(getRevParseArgs(), "--absolute-git-dir", "--git-common-dir"),
					strings.Join([]string{gitDir, gitDir}, "\n"),
					nil)

				// asking again from the directory holding the .git
				runner.ExpectGitArgs(
					append(append([]string{"-C", worktree}, getRevParseArgs()...), "--show-toplevel", "--absolute-git-dir", "--git-common-dir", "--show-superproject-working-tree"),
					strings.Join([]string{worktree, gitDir, gitDir}, "\n"),
					nil)
			},
			Path: "/path/to/repo/.git",
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
			// A repo whose work tree lives somewhere else entirely, as set up by
			// core.worktree or by --work-tree. We're in the main worktree, but the
			// git dir is not inside it.
			Name: "repo with a separate work tree",
			BeforeFunc: func(runner *oscommands.FakeCmdObjRunner, getRevParseArgs argFn) {
				mockOutput := lo.Ternary(runtime.GOOS == "windows", []string{
					// --show-toplevel
					`C:\path\to\worktree`,
					// --git-dir
					`C:\path\to\repo\.git`,
					// --git-common-dir
					`C:\path\to\repo\.git`,
					// --show-superproject-working-tree
				}, []string{
					// --show-toplevel
					"/path/to/worktree",
					// --git-dir
					"/path/to/repo/.git",
					// --git-common-dir
					"/path/to/repo/.git",
					// --show-superproject-working-tree
				})
				runner.ExpectGitArgs(
					append(getRevParseArgs(), "--show-toplevel", "--absolute-git-dir", "--git-common-dir", "--show-superproject-working-tree"),
					strings.Join(mockOutput, "\n"),
					nil)

				// asking git to find the repo from the work tree gets us nowhere,
				// because there is no .git there
				worktree := lo.Ternary(runtime.GOOS == "windows", `C:\path\to\worktree`, "/path/to/worktree")
				runner.ExpectGitArgs(
					append([]string{"-C", worktree}, append(getRevParseArgs(), "--absolute-git-dir")...),
					"",
					errors.New("fatal: not a git repository (or any of the parent directories): .git"))
			},
			Path: "/path/to/repo",
			Expected: lo.Ternary(runtime.GOOS == "windows", &RepoPaths{
				worktreePath:       `C:\path\to\worktree`,
				worktreeGitDirPath: `C:\path\to\repo\.git`,
				repoPath:           `C:\path\to\worktree`,
				repoGitDirPath:     `C:\path\to\repo\.git`,
				repoName:           `worktree`,
				isBareRepo:         false,
				gitLocationEnvVars: []string{`GIT_DIR=C:\path\to\repo\.git`, `GIT_WORK_TREE=C:\path\to\worktree`},
			}, &RepoPaths{
				worktreePath:       "/path/to/worktree",
				worktreeGitDirPath: "/path/to/repo/.git",
				repoPath:           "/path/to/worktree",
				repoGitDirPath:     "/path/to/repo/.git",
				repoName:           "worktree",
				isBareRepo:         false,
				gitLocationEnvVars: []string{"GIT_DIR=/path/to/repo/.git", "GIT_WORK_TREE=/path/to/worktree"},
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
					// --show-superproject-working-tree
					`C:\path\to\repo`,
				}, []string{
					// --show-toplevel
					"/path/to/repo/submodule1",
					// --git-dir
					"/path/to/repo/.git/modules/submodule1",
					// --git-common-dir
					"/path/to/repo/.git/modules/submodule1",
					// --show-superproject-working-tree
					"/path/to/repo",
				})
				runner.ExpectGitArgs(
					append(getRevParseArgs(), "--show-toplevel", "--absolute-git-dir", "--git-common-dir", "--show-superproject-working-tree"),
					strings.Join(mockOutput, "\n"),
					nil)

				// git finds the submodule's git dir from its work tree, via the
				// .git file there
				worktree := lo.Ternary(runtime.GOOS == "windows", `C:\path\to\repo\submodule1`, "/path/to/repo/submodule1")
				gitDir := lo.Ternary(runtime.GOOS == "windows", `C:\path\to\repo\.git\modules\submodule1`, "/path/to/repo/.git/modules/submodule1")
				runner.ExpectGitArgs(
					append([]string{"-C", worktree}, append(getRevParseArgs(), "--absolute-git-dir")...),
					gitDir,
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
					append(getRevParseArgs(), "--show-toplevel", "--absolute-git-dir", "--git-common-dir", "--show-superproject-working-tree"),
					"",
					errors.New("fatal: invalid gitfile format: /path/to/repo/worktree2/.git"))
				// we're not in a repo at all, so asking about a bare one fails too
				runner.ExpectGitArgs(
					append(getRevParseArgs(), "--absolute-git-dir", "--git-common-dir"),
					"",
					errors.New("fatal: invalid gitfile format: /path/to/repo/worktree2/.git"))
			},
			Path:     "/path/to/repo/worktree2",
			Expected: nil,
			Err: func(getRevParseArgs argFn) error {
				args := strings.Join(getRevParseArgs(), " ")
				return fmt.Errorf("'git %v --show-toplevel --absolute-git-dir --git-common-dir --show-superproject-working-tree' failed: fatal: invalid gitfile format: /path/to/repo/worktree2/.git", args)
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
