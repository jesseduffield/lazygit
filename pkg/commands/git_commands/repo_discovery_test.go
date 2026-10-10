package git_commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiscoverRepos(t *testing.T) {
	makeRepo := func(t *testing.T, root string, rel string) {
		t.Helper()
		assert.NoError(t, os.MkdirAll(filepath.Join(root, rel, ".git"), 0o755))
	}
	makeDir := func(t *testing.T, root string, rel string) {
		t.Helper()
		assert.NoError(t, os.MkdirAll(filepath.Join(root, rel), 0o755))
	}

	scenarios := []struct {
		name     string
		setup    func(t *testing.T, root string)
		maxDepth int
		expected []string
	}{
		{
			name: "empty directory",
			setup: func(t *testing.T, root string) {
				t.Helper()
			},
			maxDepth: 2,
			expected: nil,
		},
		{
			name: "direct children only at depth 1",
			setup: func(t *testing.T, root string) {
				t.Helper()
				makeRepo(t, root, "b")
				makeRepo(t, root, "a")
				makeRepo(t, root, "org/c")
			},
			maxDepth: 1,
			expected: []string{"a", "b"},
		},
		{
			name: "nested repos within depth",
			setup: func(t *testing.T, root string) {
				t.Helper()
				makeRepo(t, root, "a")
				makeRepo(t, root, "org/c")
				makeRepo(t, root, "org/deep/d")
			},
			maxDepth: 2,
			expected: []string{"a", "org/c"},
		},
		{
			name: "stops at a repo",
			setup: func(t *testing.T, root string) {
				t.Helper()
				makeRepo(t, root, "a")
				makeRepo(t, root, "a/sub")
			},
			maxDepth: 3,
			expected: []string{"a"},
		},
		{
			name: "a .git file counts as a repo",
			setup: func(t *testing.T, root string) {
				t.Helper()
				makeDir(t, root, "wt")
				assert.NoError(t, os.WriteFile(filepath.Join(root, "wt", ".git"), []byte("gitdir: /x\n"), 0o644))
			},
			maxDepth: 1,
			expected: []string{"wt"},
		},
		{
			name: "skips hidden directories",
			setup: func(t *testing.T, root string) {
				t.Helper()
				makeRepo(t, root, ".hidden/a")
				makeRepo(t, root, ".b")
			},
			maxDepth: 2,
			expected: nil,
		},
		{
			name: "does not follow symlinks",
			setup: func(t *testing.T, root string) {
				t.Helper()
				makeRepo(t, root, "a")
				assert.NoError(t, os.Symlink(filepath.Join(root, "a"), filepath.Join(root, "link")))
				assert.NoError(t, os.Symlink(root, filepath.Join(root, "loop")))
			},
			maxDepth: 5,
			expected: []string{"a"},
		},
		{
			name: "sorted by path",
			setup: func(t *testing.T, root string) {
				t.Helper()
				makeRepo(t, root, "a-c")
				makeRepo(t, root, "a/b")
			},
			maxDepth: 2,
			expected: []string{"a-c", "a/b"},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			root := t.TempDir()
			s.setup(t, root)

			var expected []string
			for _, rel := range s.expected {
				expected = append(expected, filepath.Join(root, rel))
			}
			assert.Equal(t, expected, DiscoverRepos(root, s.maxDepth))
		})
	}
}
