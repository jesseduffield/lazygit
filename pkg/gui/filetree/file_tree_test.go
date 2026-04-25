package filetree

import (
	"fmt"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/config"
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

func TestFileTreeSortOrderConfig(t *testing.T) {
	// "Dir" (uppercase D), "b-file", and "Z-file" produce distinct orderings across all
	// combinations of sort order and case sensitivity:
	//   ASCII order:            D(68) < Z(90) < b(98)
	//   Case-insensitive order: b < d < z
	files := []*models.File{
		{Path: "Dir/inner"},
		{Path: "b-file"},
		{Path: "Z-file"},
	}

	scenarios := []struct {
		sortOrder     string
		caseSensitive bool
		expected      []string
	}{
		{
			sortOrder:     "mixed",
			caseSensitive: true,
			expected:      []string{"Dir", "Dir/inner", "Z-file", "b-file"},
		},
		{
			sortOrder:     "mixed",
			caseSensitive: false,
			expected:      []string{"b-file", "Dir", "Dir/inner", "Z-file"},
		},
		{
			sortOrder:     "filesFirst",
			caseSensitive: true,
			expected:      []string{"Z-file", "b-file", "Dir", "Dir/inner"},
		},
		{
			sortOrder:     "filesFirst",
			caseSensitive: false,
			expected:      []string{"b-file", "Z-file", "Dir", "Dir/inner"},
		},
		{
			sortOrder:     "foldersFirst",
			caseSensitive: true,
			expected:      []string{"Dir", "Dir/inner", "Z-file", "b-file"},
		},
		{
			sortOrder:     "foldersFirst",
			caseSensitive: false,
			expected:      []string{"Dir", "Dir/inner", "b-file", "Z-file"},
		},
	}

	for _, s := range scenarios {
		t.Run(s.sortOrder+"/caseSensitive="+fmt.Sprintf("%v", s.caseSensitive), func(t *testing.T) {
			userConfig := config.GetDefaultConfig()
			userConfig.Gui.ShowRootItemInFileTree = false
			userConfig.Gui.FileTreeSortOrder = s.sortOrder
			userConfig.Gui.FileTreeSortCaseSensitive = s.caseSensitive
			cmn := common.NewDummyCommonWithUserConfigAndAppState(userConfig, nil)
			tree := NewFileTree(func() []*models.File { return files }, cmn, true)
			tree.SetTree()

			paths := make([]string, tree.Len())
			for i := range tree.Len() {
				paths[i] = tree.Get(i).GetPath()
			}
			assert.Equal(t, s.expected, paths)
		})
	}
}
