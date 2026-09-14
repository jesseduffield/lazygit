package helpers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadHeadInfo(t *testing.T) {
	scenarios := []struct {
		name string
		// The files to lay out below a temporary root directory, by their path
		// relative to it. "$root" in a file's content is replaced with the
		// root's path, so that a scenario can write an absolute path.
		files map[string]string
		// The repo to read, relative to the root. The directory is created
		// whether or not the scenario puts any files in it.
		repoPath   string
		expected   headInfo
		expectedOk bool
	}{
		{
			name:       "ordinary repo on a branch",
			files:      map[string]string{"repo/.git/HEAD": "ref: refs/heads/mybranch\n"},
			repoPath:   "repo",
			expected:   headInfo{branch: "mybranch"},
			expectedOk: true,
		},
		{
			name:       "ordinary repo at a detached head",
			files:      map[string]string{"repo/.git/HEAD": "d85cc9d2f5d0dc0b8f0e4d8e5b2ba0d1e7c8a3f6\n"},
			repoPath:   "repo",
			expected:   headInfo{hash: "d85cc9d2f5d0dc0b8f0e4d8e5b2ba0d1e7c8a3f6"},
			expectedOk: true,
		},
		{
			name: "worktree whose .git file names the git dir absolutely",
			files: map[string]string{
				"repo/.git/worktrees/wt/HEAD": "ref: refs/heads/mybranch\n",
				"wt/.git":                     "gitdir: $root/repo/.git/worktrees/wt\n",
			},
			repoPath:   "wt",
			expected:   headInfo{branch: "mybranch"},
			expectedOk: true,
		},
		{
			name: "worktree whose .git file names the git dir relatively",
			files: map[string]string{
				"repo/.git/worktrees/wt/HEAD": "ref: refs/heads/mybranch\n",
				"wt/.git":                     "gitdir: ../repo/.git/worktrees/wt\n",
			},
			repoPath:   "wt",
			expected:   headInfo{branch: "mybranch"},
			expectedOk: true,
		},
		{
			name: "submodule, whose .git file always names the git dir relatively",
			files: map[string]string{
				"repo/.git/modules/sub/HEAD": "ref: refs/heads/mybranch\n",
				"repo/sub/.git":              "gitdir: ../.git/modules/sub\n",
			},
			repoPath:   "repo/sub",
			expected:   headInfo{branch: "mybranch"},
			expectedOk: true,
		},
		{
			name: "repo that keeps its refs in a reftable, so HEAD holds a placeholder",
			files: map[string]string{
				"repo/.git/HEAD": "ref: refs/heads/.invalid\n",
			},
			repoPath:   "repo",
			expectedOk: false,
		},
		{
			name:       "directory without a .git entry",
			repoPath:   "notarepo",
			expectedOk: false,
		},
		{
			name:       ".git file that doesn't name a git dir",
			files:      map[string]string{"repo/.git": "not what git writes\n"},
			repoPath:   "repo",
			expectedOk: false,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			root := t.TempDir()
			for path, content := range s.files {
				fullPath := filepath.Join(root, filepath.FromSlash(path))
				assert.NoError(t, os.MkdirAll(filepath.Dir(fullPath), 0o700))
				content = strings.ReplaceAll(content, "$root", root)
				assert.NoError(t, os.WriteFile(fullPath, []byte(content), 0o600))
			}
			repoPath := filepath.Join(root, filepath.FromSlash(s.repoPath))
			assert.NoError(t, os.MkdirAll(repoPath, 0o700))

			head, ok := readHeadInfo(repoPath)

			assert.Equal(t, s.expectedOk, ok)
			assert.Equal(t, s.expected, head)
		})
	}
}
