package filetree

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/stretchr/testify/assert"
)

func TestFilterAction(t *testing.T) {
	scenarios := []struct {
		name     string
		filter   FileTreeDisplayFilter
		files    []*models.File
		expected []*models.File
	}{
		{
			name:   "filter files with unstaged changes",
			filter: DisplayUnstaged,
			files: []*models.File{
				{Path: "dir2/dir2/file4", ShortStatus: "M ", HasUnstagedChanges: true},
				{Path: "dir2/file5", ShortStatus: "M ", HasStagedChanges: true},
				{Path: "file1", ShortStatus: "M ", HasUnstagedChanges: true},
			},
			expected: []*models.File{
				{Path: "dir2/dir2/file4", ShortStatus: "M ", HasUnstagedChanges: true},
				{Path: "file1", ShortStatus: "M ", HasUnstagedChanges: true},
			},
		},
		{
			name:   "filter files with staged changes",
			filter: DisplayStaged,
			files: []*models.File{
				{Path: "dir2/dir2/file4", ShortStatus: "M ", HasStagedChanges: true},
				{Path: "dir2/file5", ShortStatus: "M ", HasStagedChanges: false},
				{Path: "file1", ShortStatus: "M ", HasStagedChanges: true},
			},
			expected: []*models.File{
				{Path: "dir2/dir2/file4", ShortStatus: "M ", HasStagedChanges: true},
				{Path: "file1", ShortStatus: "M ", HasStagedChanges: true},
			},
		},
		{
			name:   "filter files that are tracked",
			filter: DisplayTracked,
			files: []*models.File{
				{Path: "dir2/dir2/file4", ShortStatus: "M ", Tracked: true},
				{Path: "dir2/file5", ShortStatus: "M ", Tracked: false},
				{Path: "file1", ShortStatus: "M ", Tracked: true},
			},
			expected: []*models.File{
				{Path: "dir2/dir2/file4", ShortStatus: "M ", Tracked: true},
				{Path: "file1", ShortStatus: "M ", Tracked: true},
			},
		},
		{
			name:   "filter all files",
			filter: DisplayAll,
			files: []*models.File{
				{Path: "dir2/dir2/file4", ShortStatus: "M ", HasUnstagedChanges: true},
				{Path: "dir2/file5", ShortStatus: "M ", HasUnstagedChanges: true},
				{Path: "file1", ShortStatus: "M ", HasUnstagedChanges: true},
			},
			expected: []*models.File{
				{Path: "dir2/dir2/file4", ShortStatus: "M ", HasUnstagedChanges: true},
				{Path: "dir2/file5", ShortStatus: "M ", HasUnstagedChanges: true},
				{Path: "file1", ShortStatus: "M ", HasUnstagedChanges: true},
			},
		},
		{
			name:   "filter conflicted files",
			filter: DisplayConflicted,
			files: []*models.File{
				{Path: "dir2/dir2/file4", ShortStatus: "DU", HasMergeConflicts: true},
				{Path: "dir2/file5", ShortStatus: "M ", HasUnstagedChanges: true},
				{Path: "dir2/file6", ShortStatus: " M", HasStagedChanges: true},
				{Path: "file1", ShortStatus: "UU", HasMergeConflicts: true, HasInlineMergeConflicts: true},
			},
			expected: []*models.File{
				{Path: "dir2/dir2/file4", ShortStatus: "DU", HasMergeConflicts: true},
				{Path: "file1", ShortStatus: "UU", HasMergeConflicts: true, HasInlineMergeConflicts: true},
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			mngr := &FileTree{getFiles: func() []*models.File { return s.files }, filter: s.filter}
			result := mngr.getFilesForDisplay()
			assert.EqualValues(t, s.expected, result)
		})
	}
}
