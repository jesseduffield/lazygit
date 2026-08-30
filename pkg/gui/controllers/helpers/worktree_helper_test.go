package helpers

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

// nativePath rewrites a forward-slash test path into one that is valid on the
// host OS, so the scenarios below can be written with readable Unix-style
// paths. On Windows a leading slash is not absolute (filepath.IsAbs wants a
// drive letter), so we graft one on; relative paths are left untouched.
func nativePath(p string) string {
	if runtime.GOOS == "windows" && strings.HasPrefix(p, "/") {
		p = "C:" + p
	}
	return filepath.FromSlash(p)
}

func TestWorktreeParentDirCandidates(t *testing.T) {
	scenarios := []struct {
		name                string
		repoPath            string
		linkedWorktreePaths []string
		defaultPath         string
		expected            []string
	}{
		{
			name:                "no worktrees and no default path falls back to the repo's parent",
			repoPath:            "/code/myrepo",
			linkedWorktreePaths: nil,
			defaultPath:         "",
			expected:            []string{"/code"},
		},
		{
			name:                "uses the parent of each linked worktree, in order",
			repoPath:            "/code/myrepo",
			linkedWorktreePaths: []string{"/code/worktrees/foo", "/elsewhere/bar"},
			defaultPath:         "",
			expected:            []string{"/code/worktrees", "/elsewhere"},
		},
		{
			name:                "de-duplicates parents shared by multiple worktrees",
			repoPath:            "/code/myrepo",
			linkedWorktreePaths: []string{"/code/worktrees/foo", "/code/worktrees/bar"},
			defaultPath:         "",
			expected:            []string{"/code/worktrees"},
		},
		{
			name:                "appends the default path after the worktree parents",
			repoPath:            "/code/myrepo",
			linkedWorktreePaths: []string{"/code/worktrees/foo"},
			defaultPath:         "/somewhere/else",
			expected:            []string{"/code/worktrees", "/somewhere/else"},
		},
		{
			name:                "resolves a relative default path against the repo path",
			repoPath:            "/code/myrepo",
			linkedWorktreePaths: nil,
			defaultPath:         "../worktrees",
			expected:            []string{"/code/worktrees"},
		},
		{
			name:                "resolves a dot-relative default path inside the repo",
			repoPath:            "/code/myrepo",
			linkedWorktreePaths: nil,
			defaultPath:         ".worktrees",
			expected:            []string{"/code/myrepo/.worktrees"},
		},
		{
			name:                "de-duplicates the default path against a worktree parent",
			repoPath:            "/code/myrepo",
			linkedWorktreePaths: []string{"/code/worktrees/foo"},
			defaultPath:         "/code/worktrees",
			expected:            []string{"/code/worktrees"},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			result := worktreeParentDirCandidates(
				nativePath(s.repoPath),
				lo.Map(s.linkedWorktreePaths, func(p string, _ int) string { return nativePath(p) }),
				nativePath(s.defaultPath),
			)
			expected := lo.Map(s.expected, func(p string, _ int) string { return nativePath(p) })
			assert.Equal(t, expected, result)
		})
	}
}
