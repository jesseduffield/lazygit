package controllers

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestGithubPullRequestLineURL(t *testing.T) {
	const prURL = "https://github.com/jesseduffield/lazygit/pull/1234"
	const commitHash = "1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b"

	// The anchor names the file by the SHA-256 of its repo-relative path, taken over
	// exactly those bytes: no leading slash, no trailing newline, forward slashes.
	const fileHash = "067980d6efc4249367ceb61b0d93a00bca100a0ddb6d4a72b6dbb0eb9d3825cc" // "dir/file1"

	scenarios := []struct {
		name     string
		path     string
		info     types.DiffLineInfo
		expected string
	}{
		{
			name:     "an added line is on the right side of the diff",
			path:     "dir/file1",
			info:     types.DiffLineInfo{Type: types.DiffLineAdded, NewLine: 12},
			expected: prURL + "/changes/" + commitHash + "#diff-" + fileHash + "R12",
		},
		{
			name: "a deleted line is on the left side, at the line it sat on",
			path: "dir/file1",
			// A deletion's NewLine is only where it sits in the new version of the
			// file; the line it is, is the old one.
			info:     types.DiffLineInfo{Type: types.DiffLineDeleted, NewLine: 12, OldLine: 34},
			expected: prURL + "/changes/" + commitHash + "#diff-" + fileHash + "L34",
		},
		{
			name:     "a context line is on the right side too",
			path:     "dir/file1",
			info:     types.DiffLineInfo{Type: types.DiffLineContext, NewLine: 7, OldLine: 5},
			expected: prURL + "/changes/" + commitHash + "#diff-" + fileHash + "R7",
		},
		{
			name:     "a hunk header points at the first line of its hunk",
			path:     "dir/file1",
			info:     types.DiffLineInfo{Type: types.DiffLineHunkHeader, NewLine: 20},
			expected: prURL + "/changes/" + commitHash + "#diff-" + fileHash + "R20",
		},
		{
			name:     "the header naming a file points at the file alone",
			path:     "dir/file1",
			info:     types.DiffLineInfo{Type: types.DiffLineFileHeader},
			expected: prURL + "/changes/" + commitHash + "#diff-" + fileHash,
		},
		{
			name:     "a row that is no line of the file points at the file alone",
			path:     "dir/file1",
			info:     types.DiffLineInfo{Type: types.DiffLineOther},
			expected: prURL + "/changes/" + commitHash + "#diff-" + fileHash,
		},
		{
			name: "a file at the root of the repo",
			path: "file1",
			info: types.DiffLineInfo{Type: types.DiffLineAdded, NewLine: 1},
			expected: prURL + "/changes/" + commitHash +
				"#diff-c147efcfc2d7ea666a9e4f5187b115c90903f0fc896a56df9a6ef5d8f3fc9f31R1",
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			assert.Equal(t, s.expected, githubPullRequestLineURL(prURL, commitHash, s.path, s.info))
		})
	}
}

func TestGithubCommitRange(t *testing.T) {
	hashPool := &utils.StringPool{}
	newest := models.NewCommit(hashPool, models.NewCommitOpts{Hash: "newest"})
	oldest := models.NewCommit(hashPool, models.NewCommitOpts{Hash: "oldest"})

	scenarios := []struct {
		name     string
		commits  []*models.Commit
		baseHash string
		expected string
	}{
		{
			name:     "a single commit is named on its own",
			commits:  []*models.Commit{newest},
			baseHash: "parent",
			expected: "newest",
		},
		{
			name:     "a range is named as the commits it lies between",
			commits:  []*models.Commit{newest, oldest},
			baseHash: "parent",
			expected: "parent..newest",
		},
		{
			name:     "a range starting where the pull request does lies above BASE",
			commits:  []*models.Commit{newest, oldest},
			expected: "BASE..newest",
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			assert.Equal(t, s.expected, githubCommitRange(s.commits, s.baseHash))
		})
	}
}
