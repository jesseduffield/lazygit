package git_commands

import (
	"testing"

	"github.com/go-errors/errors"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func mockResolveSymlinkFn(p string) (string, error) { return p, nil }

type Scenario struct {
	Name       string
	BeforeFunc func(fs afero.Fs)
	Path       string
	Expected   *RepoPaths
	Err        error
}

func TestGetRepoPathsAux(t *testing.T) {
	scenarios := []Scenario{
		{
			Name: "typical case",
			BeforeFunc: func(fs afero.Fs) {
				// setup for main worktree
				_ = fs.MkdirAll("/path/to/repo/.git", 0o755)
			},
			Path: "/path/to/repo",
			Expected: &RepoPaths{
				currentPath:        "/path/to/repo",
				worktreePath:       "/path/to/repo",
				worktreeGitDirPath: "/path/to/repo/.git",
				repoPath:           "/path/to/repo",
				repoGitDirPath:     "/path/to/repo/.git",
				repoName:           "repo",
			},
			Err: nil,
		},
		{
			Name: "linked worktree",
			BeforeFunc: func(fs afero.Fs) {
				// setup for linked worktree
				_ = fs.MkdirAll("/path/to/repo/.git/worktrees/worktree1", 0o755)
				_ = afero.WriteFile(fs, "/path/to/repo/worktree1/.git", []byte("gitdir: /path/to/repo/.git/worktrees/worktree1"), 0o644)
			},
			Path: "/path/to/repo/worktree1",
			Expected: &RepoPaths{
				currentPath:        "/path/to/repo/worktree1",
				worktreePath:       "/path/to/repo/worktree1",
				worktreeGitDirPath: "/path/to/repo/.git/worktrees/worktree1",
				repoPath:           "/path/to/repo",
				repoGitDirPath:     "/path/to/repo/.git",
				repoName:           "repo",
			},
			Err: nil,
		},
		{
			Name: "worktree .git file missing gitdir directive",
			BeforeFunc: func(fs afero.Fs) {
				_ = fs.MkdirAll("/path/to/repo/.git/worktrees/worktree2", 0o755)
				_ = afero.WriteFile(fs, "/path/to/repo/worktree2/.git", []byte("blah"), 0o644)
			},
			Path:     "/path/to/repo/worktree2",
			Expected: nil,
			Err:      errors.New("failed to get repo git dir path: could not find git dir for /path/to/repo/worktree2: /path/to/repo/worktree2/.git is a file which suggests we are in a submodule or a worktree but the file's contents do not contain a gitdir pointing to the actual .git directory"),
		},
		{
			Name: "worktree .git file gitdir directive points to a non-existing directory",
			BeforeFunc: func(fs afero.Fs) {
				_ = fs.MkdirAll("/path/to/repo/.git/worktrees/worktree2", 0o755)
				_ = afero.WriteFile(fs, "/path/to/repo/worktree2/.git", []byte("gitdir: /nonexistant"), 0o644)
			},
			Path:     "/path/to/repo/worktree2",
			Expected: nil,
			Err:      errors.New("failed to get repo git dir path: could not find git dir for /path/to/repo/worktree2. /nonexistant does not exist"),
		},
		{
			Name: "submodule",
			BeforeFunc: func(fs afero.Fs) {
				_ = fs.MkdirAll("/path/to/repo/.git/modules/submodule1", 0o755)
				_ = afero.WriteFile(fs, "/path/to/repo/submodule1/.git", []byte("gitdir: /path/to/repo/.git/modules/submodule1"), 0o644)
			},
			Path: "/path/to/repo/submodule1",
			Expected: &RepoPaths{
				currentPath:        "/path/to/repo/submodule1",
				worktreePath:       "/path/to/repo/submodule1",
				worktreeGitDirPath: "/path/to/repo/.git/modules/submodule1",
				repoPath:           "/path/to/repo/submodule1",
				repoGitDirPath:     "/path/to/repo/.git/modules/submodule1",
				repoName:           "submodule1",
			},
			Err: nil,
		},
		{
			Name: "submodule in nested directory",
			BeforeFunc: func(fs afero.Fs) {
				_ = fs.MkdirAll("/path/to/repo/.git/modules/my/submodule1", 0o755)
				_ = afero.WriteFile(fs, "/path/to/repo/my/submodule1/.git", []byte("gitdir: /path/to/repo/.git/modules/my/submodule1"), 0o644)
			},
			Path: "/path/to/repo/my/submodule1",
			Expected: &RepoPaths{
				currentPath:        "/path/to/repo/my/submodule1",
				worktreePath:       "/path/to/repo/my/submodule1",
				worktreeGitDirPath: "/path/to/repo/.git/modules/my/submodule1",
				repoPath:           "/path/to/repo/my/submodule1",
				repoGitDirPath:     "/path/to/repo/.git/modules/my/submodule1",
				repoName:           "submodule1",
			},
			Err: nil,
		},
		{
			Name: "submodule git dir not under .git/modules",
			BeforeFunc: func(fs afero.Fs) {
				_ = fs.MkdirAll("/random/submodule1", 0o755)
				_ = afero.WriteFile(fs, "/path/to/repo/my/submodule1/.git", []byte("gitdir: /random/submodule1"), 0o644)
			},
			Path:     "/path/to/repo/my/submodule1",
			Expected: nil,
			Err:      errors.New("failed to get repo git dir path: could not find git dir for /path/to/repo/my/submodule1: path is not under `worktrees` or `modules` directories"),
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.Name, func(t *testing.T) {
			fs := afero.NewMemMapFs()

			// prepare the filesystem for the scenario
			s.BeforeFunc(fs)

			// run the function with the scenario path
			repoPaths, err := getRepoPathsAux(fs, mockResolveSymlinkFn, s.Path)

			// check the error and the paths
			if s.Err != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, s.Err.Error())
			} else {
				assert.Nil(t, err)
				assert.Equal(t, s.Expected, repoPaths)
			}
		})
	}
}
