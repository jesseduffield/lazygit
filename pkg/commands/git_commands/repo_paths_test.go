package git_commands

import (
	"fmt"
	"strings"
	"testing"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
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
				expectedOutput := []string{
					// --show-toplevel
					"/path/to/repo",
					// --git-dir
					"/path/to/repo/.git",
					// --git-common-dir
					"/path/to/repo/.git",
					// --show-superproject-working-tree
				}
				runner.ExpectGitArgs(
					append(getRevParseArgs(), "--show-toplevel", "--absolute-git-dir", "--git-common-dir", "--show-superproject-working-tree"),
					strings.Join(expectedOutput, "\n"),
					nil)
			},
			Path: "/path/to/repo",
			Expected: &RepoPaths{
				worktreePath:       "/path/to/repo",
				worktreeGitDirPath: "/path/to/repo/.git",
				repoPath:           "/path/to/repo",
				repoGitDirPath:     "/path/to/repo/.git",
				repoName:           "repo",
			},
			Err: nil,
		},
		{
			Name: "submodule",
			BeforeFunc: func(runner *oscommands.FakeCmdObjRunner, getRevParseArgs argFn) {
				expectedOutput := []string{
					// --show-toplevel
					"/path/to/repo/submodule1",
					// --git-dir
					"/path/to/repo/.git/modules/submodule1",
					// --git-common-dir
					"/path/to/repo/.git/modules/submodule1",
					// --show-superproject-working-tree
					"/path/to/repo",
				}
				runner.ExpectGitArgs(
					append(getRevParseArgs(), "--show-toplevel", "--absolute-git-dir", "--git-common-dir", "--show-superproject-working-tree"),
					strings.Join(expectedOutput, "\n"),
					nil)
			},
			Path: "/path/to/repo/submodule1",
			Expected: &RepoPaths{
				worktreePath:       "/path/to/repo/submodule1",
				worktreeGitDirPath: "/path/to/repo/.git/modules/submodule1",
				repoPath:           "/path/to/repo/submodule1",
				repoGitDirPath:     "/path/to/repo/.git/modules/submodule1",
				repoName:           "submodule1",
			},
			Err: nil,
		},
		{
			Name: "git rev-parse returns an error",
			BeforeFunc: func(runner *oscommands.FakeCmdObjRunner, getRevParseArgs argFn) {
				runner.ExpectGitArgs(
					append(getRevParseArgs(), "--show-toplevel", "--absolute-git-dir", "--git-common-dir", "--show-superproject-working-tree"),
					"",
					errors.New("fatal: invalid gitfile format: /path/to/repo/worktree2/.git"))
			},
			Path:     "/path/to/repo/worktree2",
			Expected: nil,
			Err: func(getRevParseArgs argFn) error {
				args := strings.Join(getRevParseArgs(), " ")
				return errors.New(
					fmt.Sprintf("'git %v --show-toplevel --absolute-git-dir --git-common-dir --show-superproject-working-tree' failed: fatal: invalid gitfile format: /path/to/repo/worktree2/.git", args),
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

			repoPaths, err := GetRepoPaths(cmd, version)

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
