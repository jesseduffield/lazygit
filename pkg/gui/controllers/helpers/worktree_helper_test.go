package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			result := worktreeParentDirCandidates(s.repoPath, s.linkedWorktreePaths, s.defaultPath)
			assert.Equal(t, s.expected, result)
		})
	}
}
