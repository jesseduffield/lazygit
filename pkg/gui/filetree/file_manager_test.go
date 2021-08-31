package filetree

import (
	"testing"

	"github.com/gookit/color"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/stretchr/testify/assert"
	"github.com/xo/terminfo"
)

func TestRender(t *testing.T) {
	color.ForceSetColorLevel(terminfo.ColorLevelNone)

	scenarios := []struct {
		name           string
		root           *FileNode
		collapsedPaths map[string]bool
		expected       []string
	}{
		{
			name:     "nil node",
			root:     nil,
			expected: []string{},
		},
		{
			name: "leaf node",
			root: &FileNode{
				Path: "",
				Children: []*FileNode{
					{File: &models.File{Name: "test", ShortStatus: " M", HasStagedChanges: true}, Path: "test"},
				},
			},
			expected: []string{" M test"},
		},
		{
			name: "big example",
			root: &FileNode{
				Path: "",
				Children: []*FileNode{
					{
						Path: "dir1",
						Children: []*FileNode{
							{
								File: &models.File{Name: "dir1/file2", ShortStatus: "M ", HasUnstagedChanges: true},
								Path: "dir1/file2",
							},
							{
								File: &models.File{Name: "dir1/file3", ShortStatus: "M ", HasUnstagedChanges: true},
								Path: "dir1/file3",
							},
						},
					},
					{
						Path: "dir2",
						Children: []*FileNode{
							{
								Path: "dir2/dir2",
								Children: []*FileNode{
									{
										File: &models.File{Name: "dir2/dir2/file3", ShortStatus: " M", HasStagedChanges: true},
										Path: "dir2/dir2/file3",
									},
									{
										File: &models.File{Name: "dir2/dir2/file4", ShortStatus: "M ", HasUnstagedChanges: true},
										Path: "dir2/dir2/file4",
									},
								},
							},
							{
								File: &models.File{Name: "dir2/file5", ShortStatus: "M ", HasUnstagedChanges: true},
								Path: "dir2/file5",
							},
						},
					},
					{
						File: &models.File{Name: "file1", ShortStatus: "M ", HasUnstagedChanges: true},
						Path: "file1",
					},
				},
			},
			expected:       []string{"dir1 ►", "dir2 ▼", "├─ dir2 ▼", "│  ├─  M file3", "│  └─ M  file4", "└─ M  file5", "M  file1"},
			collapsedPaths: map[string]bool{"dir1": true},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.name, func(t *testing.T) {
			mngr := &FileManager{tree: s.root, collapsedPaths: s.collapsedPaths}
			result := mngr.Render("", nil)
			assert.EqualValues(t, s.expected, result)
		})
	}
}

func TestFilterAction(t *testing.T) {
	scenarios := []struct {
		name     string
		filters  []bool
		files    []*models.File
		expected []*models.File
	}{
		{
			name: "filter untracked files",
			filters: []bool{
				DisplayAll:        false,
				DisplayStaged:     false,
				DisplayModified:   false,
				DisplayConflicted: false,
				DisplayUntracked:  true,
			},
			files: []*models.File{
				{Name: "dir2/dir2/file4", ShortStatus: "M ", HasUnstagedChanges: true, Tracked: true},
				{Name: "dir2/file5", ShortStatus: "M ", HasStagedChanges: true, Tracked: true},
				{Name: "file1", ShortStatus: "??", Tracked: false},
			},
			expected: []*models.File{
				{Name: "file1", ShortStatus: "??", Tracked: false},
			},
		},
		{
			name: "filter modified files",
			filters: []bool{
				DisplayAll:        false,
				DisplayStaged:     false,
				DisplayModified:   true,
				DisplayConflicted: false,
				DisplayUntracked:  false,
			},
			files: []*models.File{
				{Name: "dir2/dir2/file4", ShortStatus: "M ", HasUnstagedChanges: true, Tracked: true},
				{Name: "dir2/file5", ShortStatus: "M ", HasStagedChanges: true},
				{Name: "file1", ShortStatus: "M ", HasUnstagedChanges: true, Tracked: false},
			},
			expected: []*models.File{
				{Name: "dir2/dir2/file4", ShortStatus: "M ", HasUnstagedChanges: true, Tracked: true},
			},
		},
		{
			name: "filter files with conflicts",
			filters: []bool{
				DisplayAll:        false,
				DisplayStaged:     false,
				DisplayModified:   false,
				DisplayConflicted: true,
				DisplayUntracked:  false,
			},
			files: []*models.File{
				{Name: "dir2/dir2/file4", ShortStatus: "UU", HasUnstagedChanges: true, HasMergeConflicts: true},
				{Name: "dir2/file5", ShortStatus: "UU", HasMergeConflicts: false, HasInlineMergeConflicts: true},
				{Name: "file1", ShortStatus: "M ", HasStagedChanges: true, HasMergeConflicts: false},
			},
			expected: []*models.File{
				{Name: "dir2/dir2/file4", ShortStatus: "UU", HasUnstagedChanges: true, HasMergeConflicts: true},
			},
		},
		{
			name: "filter files with staged changes",
			filters: []bool{
				DisplayAll:        false,
				DisplayStaged:     true,
				DisplayModified:   false,
				DisplayConflicted: false,
				DisplayUntracked:  false,
			},
			files: []*models.File{
				{Name: "dir2/dir2/file4", ShortStatus: "M ", HasStagedChanges: true},
				{Name: "dir2/file5", ShortStatus: "M ", HasStagedChanges: false},
				{Name: "file1", ShortStatus: "M ", HasStagedChanges: true},
			},
			expected: []*models.File{
				{Name: "dir2/dir2/file4", ShortStatus: "M ", HasStagedChanges: true},
				{Name: "file1", ShortStatus: "M ", HasStagedChanges: true},
			},
		},
		{
			name: "filter all files",
			filters: []bool{
				DisplayAll:        true,
				DisplayStaged:     false,
				DisplayModified:   false,
				DisplayConflicted: false,
				DisplayUntracked:  false,
			},
			files: []*models.File{
				{Name: "dir2/dir2/file4", ShortStatus: "M ", HasUnstagedChanges: true},
				{Name: "dir2/file5", ShortStatus: "M ", HasUnstagedChanges: true},
				{Name: "file1", ShortStatus: "M ", HasUnstagedChanges: true},
			},
			expected: []*models.File{
				{Name: "dir2/dir2/file4", ShortStatus: "M ", HasUnstagedChanges: true},
				{Name: "dir2/file5", ShortStatus: "M ", HasUnstagedChanges: true},
				{Name: "file1", ShortStatus: "M ", HasUnstagedChanges: true},
			},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.name, func(t *testing.T) {
			mngr := &FileManager{files: s.files, filters: s.filters}
			result := mngr.GetFilesForDisplay()
			assert.EqualValues(t, s.expected, result)
		})
	}
}
