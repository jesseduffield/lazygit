package commands

import (
	"testing"

	"github.com/go-errors/errors"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestFindWorktreeRoot(t *testing.T) {
	type scenario struct {
		testName     string
		currentPath  string
		before       func(fs afero.Fs)
		expectedPath string
		expectedErr  string
	}

	scenarios := []scenario{
		{
			testName:    "at root of worktree",
			currentPath: "/path/to/repo",
			before: func(fs afero.Fs) {
				_ = fs.MkdirAll("/path/to/repo/.git", 0o755)
			},
			expectedPath: "/path/to/repo",
			expectedErr:  "",
		},
		{
			testName:    "inside worktree",
			currentPath: "/path/to/repo/subdir",
			before: func(fs afero.Fs) {
				_ = fs.MkdirAll("/path/to/repo/.git", 0o755)
				_ = fs.MkdirAll("/path/to/repo/subdir", 0o755)
			},
			expectedPath: "/path/to/repo",
			expectedErr:  "",
		},
		{
			testName:     "not in a git repo",
			currentPath:  "/path/to/dir",
			before:       func(fs afero.Fs) {},
			expectedPath: "",
			expectedErr:  "Must open lazygit in a git repository",
		},
		{
			testName:    "In linked worktree",
			currentPath: "/path/to/worktree",
			before: func(fs afero.Fs) {
				_ = fs.MkdirAll("/path/to/worktree", 0o755)
				_ = afero.WriteFile(fs, "/path/to/worktree/.git", []byte("blah"), 0o755)
			},
			expectedPath: "/path/to/worktree",
			expectedErr:  "",
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			s.before(fs)

			root, err := findWorktreeRoot(fs, s.currentPath)
			if s.expectedErr != "" {
				assert.EqualError(t, errors.New(s.expectedErr), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, s.expectedPath, root)
			}
		})
	}
}
