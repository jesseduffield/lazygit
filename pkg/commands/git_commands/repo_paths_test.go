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
				return fmt.Errorf("'git %v --show-toplevel --absolute-git-dir --git-common-dir --is-bare-repository --show-superproject-working-tree' failed: fatal: invalid gitfile format: /path/to/repo/worktree2/.git", args)
			},
		},
		{
			Name: "bare repo with worktree setup",
			BeforeFunc: func(runner *oscommands.FakeCmdObjRunner, getRevParseArgs argFn) {
				runner.ExpectGitArgs(
					append(getRevParseArgs(), "--show-toplevel", "--absolute-git-dir", "--git-common-dir", "--is-bare-repository", "--show-superproject-working-tree"),
					"",
					errors.New("fatal: this operation must be run in a work tree"),
				)

				runner.ExpectGitArgs(
					[]string{"-C", ".bare", "rev-parse", "--is-bare-repository"},
					"true",
					nil,
				)

				runner.ExpectGitArgs(
					[]string{"-C", ".bare", "worktree", "list", "--porcelain"},
					"worktree /path/to/parent/main\nHEAD abc123\nbranch refs/heads/main\n\n",
					nil,
				)

				runner.ExpectGitArgs(
					[]string{"-C", ".bare", "symbolic-ref", "refs/remotes/origin/HEAD", "--short"},
					"origin/main",
					nil,
				)

				mockOutput := []string{
					"/path/to/parent/main",
					"/path/to/parent/.bare",
					"/path/to/parent/.bare",
					"false",
				}
				runner.ExpectGitArgs(
					append([]string{"-C", "/path/to/parent/main"}, append(getRevParseArgs(), "--show-toplevel", "--absolute-git-dir", "--git-common-dir", "--is-bare-repository", "--show-superproject-working-tree")...),
					strings.Join(mockOutput, "\n"),
					nil,
				)
			},
			Path: "",
			Expected: &RepoPaths{
				worktreePath:       "/path/to/parent/main",
				worktreeGitDirPath: "/path/to/parent/.bare",
				repoPath:           "/path/to/parent",
				repoGitDirPath:     "/path/to/parent/.bare",
				repoName:           "parent",
				isBareRepo:         false,
			},
			Err: nil,
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

func TestParseWorktreeList(t *testing.T) {
	output := `worktree /path/to/repo/main
HEAD abc123456
branch refs/heads/main

worktree /path/to/repo/feature
HEAD def789012
branch refs/heads/feature-branch

worktree /path/to/repo/detached
HEAD ghi345678
detached

`

	worktrees := parseWorktreeList(output)
	expected := []WorktreeInfo{
		{Path: "/path/to/repo/main", Head: "abc123456", Branch: "refs/heads/main"},
		{Path: "/path/to/repo/feature", Head: "def789012", Branch: "refs/heads/feature-branch"},
		{Path: "/path/to/repo/detached", Head: "ghi345678", Branch: ""},
	}

	assert.Equal(t, expected, worktrees)
}

func TestSelectBestWorktree(t *testing.T) {
	parentDir := "/path/to/parent"
	bareDir := "/path/to/parent/.bare"

	scenarios := []struct {
		Name              string
		Worktrees         []WorktreeInfo
		Expected          WorktreeInfo
		MockDefaultBranch func(runner *oscommands.FakeCmdObjRunner)
	}{
		{
			Name: "single worktree",
			Worktrees: []WorktreeInfo{
				{Path: "/path/to/parent/single", Branch: "refs/heads/feature"},
			},
			Expected: WorktreeInfo{Path: "/path/to/parent/single", Branch: "refs/heads/feature"},
			MockDefaultBranch: func(runner *oscommands.FakeCmdObjRunner) {
				runner.ExpectGitArgs(
					[]string{"-C", bareDir, "symbolic-ref", "refs/remotes/origin/HEAD", "--short"},
					"",
					errors.New("not found"),
				)
				runner.ExpectGitArgs(
					[]string{"-C", bareDir, "config", "init.defaultBranch"},
					"",
					errors.New("not found"),
				)
			},
		},
		{
			Name: "prefer main branch",
			Worktrees: []WorktreeInfo{
				{Path: "/path/to/parent/feature", Branch: "refs/heads/feature"},
				{Path: "/path/to/parent/main", Branch: "refs/heads/main"},
				{Path: "/path/to/parent/other", Branch: "refs/heads/other"},
			},
			Expected: WorktreeInfo{Path: "/path/to/parent/main", Branch: "refs/heads/main"},
			MockDefaultBranch: func(runner *oscommands.FakeCmdObjRunner) {
				runner.ExpectGitArgs(
					[]string{"-C", bareDir, "symbolic-ref", "refs/remotes/origin/HEAD", "--short"},
					"origin/main",
					nil,
				)
			},
		},
		{
			Name: "prefer master branch",
			Worktrees: []WorktreeInfo{
				{Path: "/path/to/parent/feature", Branch: "refs/heads/feature"},
				{Path: "/path/to/parent/master", Branch: "refs/heads/master"},
			},
			Expected: WorktreeInfo{Path: "/path/to/parent/master", Branch: "refs/heads/master"},
			MockDefaultBranch: func(runner *oscommands.FakeCmdObjRunner) {
				runner.ExpectGitArgs(
					[]string{"-C", bareDir, "symbolic-ref", "refs/remotes/origin/HEAD", "--short"},
					"origin/master",
					nil,
				)
			},
		},
		{
			Name: "prefer main directory name",
			Worktrees: []WorktreeInfo{
				{Path: "/path/to/parent/feature", Branch: "refs/heads/feature"},
				{Path: "/path/to/parent/main", Branch: "refs/heads/feature"},
			},
			Expected: WorktreeInfo{Path: "/path/to/parent/main", Branch: "refs/heads/feature"},
			MockDefaultBranch: func(runner *oscommands.FakeCmdObjRunner) {
				runner.ExpectGitArgs(
					[]string{"-C", bareDir, "symbolic-ref", "refs/remotes/origin/HEAD", "--short"},
					"origin/main",
					nil,
				)
			},
		},
		{
			Name: "prefer master directory name",
			Worktrees: []WorktreeInfo{
				{Path: "/path/to/parent/feature", Branch: "refs/heads/feature"},
				{Path: "/path/to/parent/master", Branch: "refs/heads/feature"},
			},
			Expected: WorktreeInfo{Path: "/path/to/parent/master", Branch: "refs/heads/feature"},
			MockDefaultBranch: func(runner *oscommands.FakeCmdObjRunner) {
				runner.ExpectGitArgs(
					[]string{"-C", bareDir, "symbolic-ref", "refs/remotes/origin/HEAD", "--short"},
					"origin/master",
					nil,
				)
			},
		},
		{
			Name: "first alphabetically when no preferences",
			Worktrees: []WorktreeInfo{
				{Path: "/path/to/parent/z-feature", Branch: "refs/heads/feature"},
				{Path: "/path/to/parent/a-feature", Branch: "refs/heads/other"},
			},
			Expected: WorktreeInfo{Path: "/path/to/parent/a-feature", Branch: "refs/heads/other"},
			MockDefaultBranch: func(runner *oscommands.FakeCmdObjRunner) {
				runner.ExpectGitArgs(
					[]string{"-C", bareDir, "symbolic-ref", "refs/remotes/origin/HEAD", "--short"},
					"",
					errors.New("not found"),
				)
				runner.ExpectGitArgs(
					[]string{"-C", bareDir, "config", "init.defaultBranch"},
					"",
					errors.New("not found"),
				)
			},
		},
		{
			Name: "prefer custom default branch",
			Worktrees: []WorktreeInfo{
				{Path: "/path/to/parent/feature", Branch: "refs/heads/feature"},
				{Path: "/path/to/parent/develop", Branch: "refs/heads/develop"},
				{Path: "/path/to/parent/master", Branch: "refs/heads/master"},
			},
			Expected: WorktreeInfo{Path: "/path/to/parent/develop", Branch: "refs/heads/develop"},
			MockDefaultBranch: func(runner *oscommands.FakeCmdObjRunner) {
				runner.ExpectGitArgs(
					[]string{"-C", bareDir, "symbolic-ref", "refs/remotes/origin/HEAD", "--short"},
					"origin/develop",
					nil,
				)
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.Name, func(t *testing.T) {
			runner := oscommands.NewFakeRunner(t)
			cmd := oscommands.NewDummyCmdObjBuilder(runner)

			if s.MockDefaultBranch != nil {
				s.MockDefaultBranch(runner)
			}

			result := selectBestWorktree(s.Worktrees, parentDir, bareDir, cmd)
			assert.Equal(t, s.Expected, result)
		})
	}
}
