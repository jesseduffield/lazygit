package controllers

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestDiffPathsForNode(t *testing.T) {
	files := []*models.CommitFile{
		{Path: "dir/file1", PreviousPath: "file1", ChangeStatus: "R"},
		{Path: "dir/file2-renamed", PreviousPath: "dir/file2", ChangeStatus: "R"},
		{Path: "dir/sub/file3", ChangeStatus: "M"},
		{Path: "file4", PreviousPath: "dir/sub/file4", ChangeStatus: "R"},
		{Path: "file5", ChangeStatus: "M"},
	}

	scenarios := []struct {
		testName      string
		files         []*models.CommitFile // defaults to the files above
		selectedPath  string
		isFiltering   bool
		expectedPaths []string
	}{
		{
			testName:      "file",
			selectedPath:  "dir/sub/file3",
			expectedPaths: []string{"dir/sub/file3"},
		},
		{
			testName:      "renamed file",
			selectedPath:  "dir/file1",
			expectedPaths: []string{"dir/file1", "file1"},
		},
		{
			testName:     "directory: pass the other end of each rename that crosses its boundary",
			selectedPath: "dir",
			// dir/file2-renamed was renamed within the directory, so both of its
			// paths are covered by it already
			expectedPaths: []string{"dir", "file1", "file4"},
		},
		{
			testName:      "directory without renames crossing its boundary",
			selectedPath:  "dir/sub",
			expectedPaths: []string{"dir/sub", "file4"},
		},
		{
			testName:      "root",
			selectedPath:  ".",
			expectedPaths: []string{"."},
		},
		{
			testName: "a whole directory moved into the selected one collapses to that directory",
			files: []*models.CommitFile{
				{Path: "dir/a", PreviousPath: "src/a", ChangeStatus: "R"},
				{Path: "dir/b", PreviousPath: "src/nested/b", ChangeStatus: "R"},
				{Path: "dir/c", PreviousPath: "src/nested/c", ChangeStatus: "R"},
				{Path: "unrelated", ChangeStatus: "M"},
			},
			selectedPath:  "dir",
			expectedPaths: []string{"dir", "src"},
		},
		{
			testName: "a directory that stands in for the selected one as well",
			files: []*models.CommitFile{
				{Path: "a/b/c", PreviousPath: "a/c", ChangeStatus: "R"},
				{Path: "a/b/d", ChangeStatus: "M"},
				{Path: "unrelated", ChangeStatus: "M"},
			},
			selectedPath:  "a/b",
			expectedPaths: []string{"a"},
		},
		{
			testName: "a directory with changes of its own doesn't collapse",
			files: []*models.CommitFile{
				{Path: "dir/a", PreviousPath: "src/a", ChangeStatus: "R"},
				{Path: "dir/b", PreviousPath: "src/nested/b", ChangeStatus: "R"},
				{Path: "src/nested/c", ChangeStatus: "M"},
			},
			selectedPath: "dir",
			// src/nested is left out of it, so that only src/a stays behind
			expectedPaths: []string{"dir", "src/a", "src/nested/b"},
		},
		{
			testName:     "directory while filtering",
			selectedPath: "dir",
			isFiltering:  true,
			expectedPaths: []string{
				"dir/file1", "file1",
				"dir/file2-renamed", "dir/file2",
				"dir/sub/file3",
				"file4", "dir/sub/file4",
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			files := lo.Ternary(s.files != nil, s.files, files)
			cmp := filetree.NodeSortComparator[models.CommitFile]("mixed", true)
			root := filetree.BuildTreeFromCommitFiles(files, true, cmp)
			node, found := lo.Find(root.Flatten(filetree.NewCollapsedPaths()), func(node *filetree.Node[models.CommitFile]) bool {
				return node.GetPath() == s.selectedPath
			})
			assert.True(t, found, "no node for path %s", s.selectedPath)

			assert.Equal(t, s.expectedPaths, diffPathsForNode(node, root, files, s.isFiltering))
		})
	}
}
