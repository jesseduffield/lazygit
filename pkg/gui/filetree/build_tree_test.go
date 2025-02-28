package filetree

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/stretchr/testify/assert"
)

func TestBuildTreeFromFiles(t *testing.T) {
	scenarios := []struct {
		name     string
		files    []*models.File
		expected *Node[models.File]
	}{
		{
			name:  "no files",
			files: []*models.File{},
			expected: &Node[models.File]{
				Path:     "",
				Children: nil,
			},
		},
		{
			name: "files in same directory",
			files: []*models.File{
				{
					Path: "dir1/a",
				},
				{
					Path: "dir1/b",
				},
			},
			expected: &Node[models.File]{
				Path: "",
				Children: []*Node[models.File]{
					{
						Path: "dir1",
						Children: []*Node[models.File]{
							{
								File: &models.File{Path: "dir1/a"},
								Path: "dir1/a",
							},
							{
								File: &models.File{Path: "dir1/b"},
								Path: "dir1/b",
							},
						},
					},
				},
			},
		},
		{
			name: "paths that can be compressed",
			files: []*models.File{
				{
					Path: "dir1/dir3/a",
				},
				{
					Path: "dir2/dir4/b",
				},
			},
			expected: &Node[models.File]{
				Path: "",
				Children: []*Node[models.File]{
					{
						Path: "dir1/dir3",
						Children: []*Node[models.File]{
							{
								File: &models.File{Path: "dir1/dir3/a"},
								Path: "dir1/dir3/a",
							},
						},
						CompressionLevel: 1,
					},
					{
						Path: "dir2/dir4",
						Children: []*Node[models.File]{
							{
								File: &models.File{Path: "dir2/dir4/b"},
								Path: "dir2/dir4/b",
							},
						},
						CompressionLevel: 1,
					},
				},
			},
		},
		{
			name: "paths that can be sorted",
			files: []*models.File{
				{
					Path: "b",
				},
				{
					Path: "a",
				},
			},
			expected: &Node[models.File]{
				Path: "",
				Children: []*Node[models.File]{
					{
						File: &models.File{Path: "a"},
						Path: "a",
					},
					{
						File: &models.File{Path: "b"},
						Path: "b",
					},
				},
			},
		},
		{
			name: "paths that can be sorted including a merge conflict file",
			files: []*models.File{
				{
					Path: "b",
				},
				{
					Path:              "z",
					HasMergeConflicts: true,
				},
				{
					Path: "a",
				},
			},
			expected: &Node[models.File]{
				Path: "",
				// it is a little strange that we're not bubbling up our merge conflict
				// here but we are technically still in tree mode and that's the rule
				Children: []*Node[models.File]{
					{
						File: &models.File{Path: "a"},
						Path: "a",
					},
					{
						File: &models.File{Path: "b"},
						Path: "b",
					},
					{
						File: &models.File{Path: "z", HasMergeConflicts: true},
						Path: "z",
					},
				},
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			result := BuildTreeFromFiles(s.files)
			assert.EqualValues(t, s.expected, result)
		})
	}
}

func TestBuildFlatTreeFromFiles(t *testing.T) {
	scenarios := []struct {
		name     string
		files    []*models.File
		expected *Node[models.File]
	}{
		{
			name:  "no files",
			files: []*models.File{},
			expected: &Node[models.File]{
				Path:     "",
				Children: []*Node[models.File]{},
			},
		},
		{
			name: "files in same directory",
			files: []*models.File{
				{
					Path: "dir1/a",
				},
				{
					Path: "dir1/b",
				},
			},
			expected: &Node[models.File]{
				Path: "",
				Children: []*Node[models.File]{
					{
						File:             &models.File{Path: "dir1/a"},
						Path:             "dir1/a",
						CompressionLevel: 0,
					},
					{
						File:             &models.File{Path: "dir1/b"},
						Path:             "dir1/b",
						CompressionLevel: 0,
					},
				},
			},
		},
		{
			name: "paths that can be compressed",
			files: []*models.File{
				{
					Path: "dir1/a",
				},
				{
					Path: "dir2/b",
				},
			},
			expected: &Node[models.File]{
				Path: "",
				Children: []*Node[models.File]{
					{
						File:             &models.File{Path: "dir1/a"},
						Path:             "dir1/a",
						CompressionLevel: 0,
					},
					{
						File:             &models.File{Path: "dir2/b"},
						Path:             "dir2/b",
						CompressionLevel: 0,
					},
				},
			},
		},
		{
			name: "paths that can be sorted",
			files: []*models.File{
				{
					Path: "b",
				},
				{
					Path: "a",
				},
			},
			expected: &Node[models.File]{
				Path: "",
				Children: []*Node[models.File]{
					{
						File: &models.File{Path: "a"},
						Path: "a",
					},
					{
						File: &models.File{Path: "b"},
						Path: "b",
					},
				},
			},
		},
		{
			name: "tracked, untracked, and conflicted files",
			files: []*models.File{
				{
					Path:    "a2",
					Tracked: false,
				},
				{
					Path:    "a1",
					Tracked: false,
				},
				{
					Path:              "c2",
					HasMergeConflicts: true,
				},
				{
					Path:              "c1",
					HasMergeConflicts: true,
				},
				{
					Path:    "b2",
					Tracked: true,
				},
				{
					Path:    "b1",
					Tracked: true,
				},
			},
			expected: &Node[models.File]{
				Path: "",
				Children: []*Node[models.File]{
					{
						File: &models.File{Path: "c1", HasMergeConflicts: true},
						Path: "c1",
					},
					{
						File: &models.File{Path: "c2", HasMergeConflicts: true},
						Path: "c2",
					},
					{
						File: &models.File{Path: "b1", Tracked: true},
						Path: "b1",
					},
					{
						File: &models.File{Path: "b2", Tracked: true},
						Path: "b2",
					},
					{
						File: &models.File{Path: "a1", Tracked: false},
						Path: "a1",
					},
					{
						File: &models.File{Path: "a2", Tracked: false},
						Path: "a2",
					},
				},
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			result := BuildFlatTreeFromFiles(s.files)
			assert.EqualValues(t, s.expected, result)
		})
	}
}

func TestBuildTreeFromCommitFiles(t *testing.T) {
	scenarios := []struct {
		name     string
		files    []*models.CommitFile
		expected *Node[models.CommitFile]
	}{
		{
			name:  "no files",
			files: []*models.CommitFile{},
			expected: &Node[models.CommitFile]{
				Path:     "",
				Children: nil,
			},
		},
		{
			name: "files in same directory",
			files: []*models.CommitFile{
				{
					Path: "dir1/a",
				},
				{
					Path: "dir1/b",
				},
			},
			expected: &Node[models.CommitFile]{
				Path: "",
				Children: []*Node[models.CommitFile]{
					{
						Path: "dir1",
						Children: []*Node[models.CommitFile]{
							{
								File: &models.CommitFile{Path: "dir1/a"},
								Path: "dir1/a",
							},
							{
								File: &models.CommitFile{Path: "dir1/b"},
								Path: "dir1/b",
							},
						},
					},
				},
			},
		},
		{
			name: "paths that can be compressed",
			files: []*models.CommitFile{
				{
					Path: "dir1/dir3/a",
				},
				{
					Path: "dir2/dir4/b",
				},
			},
			expected: &Node[models.CommitFile]{
				Path: "",
				Children: []*Node[models.CommitFile]{
					{
						Path: "dir1/dir3",
						Children: []*Node[models.CommitFile]{
							{
								File: &models.CommitFile{Path: "dir1/dir3/a"},
								Path: "dir1/dir3/a",
							},
						},
						CompressionLevel: 1,
					},
					{
						Path: "dir2/dir4",
						Children: []*Node[models.CommitFile]{
							{
								File: &models.CommitFile{Path: "dir2/dir4/b"},
								Path: "dir2/dir4/b",
							},
						},
						CompressionLevel: 1,
					},
				},
			},
		},
		{
			name: "paths that can be sorted",
			files: []*models.CommitFile{
				{
					Path: "b",
				},
				{
					Path: "a",
				},
			},
			expected: &Node[models.CommitFile]{
				Path: "",
				Children: []*Node[models.CommitFile]{
					{
						File: &models.CommitFile{Path: "a"},
						Path: "a",
					},
					{
						File: &models.CommitFile{Path: "b"},
						Path: "b",
					},
				},
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			result := BuildTreeFromCommitFiles(s.files)
			assert.EqualValues(t, s.expected, result)
		})
	}
}

func TestBuildFlatTreeFromCommitFiles(t *testing.T) {
	scenarios := []struct {
		name     string
		files    []*models.CommitFile
		expected *Node[models.CommitFile]
	}{
		{
			name:  "no files",
			files: []*models.CommitFile{},
			expected: &Node[models.CommitFile]{
				Path:     "",
				Children: []*Node[models.CommitFile]{},
			},
		},
		{
			name: "files in same directory",
			files: []*models.CommitFile{
				{
					Path: "dir1/a",
				},
				{
					Path: "dir1/b",
				},
			},
			expected: &Node[models.CommitFile]{
				Path: "",
				Children: []*Node[models.CommitFile]{
					{
						File:             &models.CommitFile{Path: "dir1/a"},
						Path:             "dir1/a",
						CompressionLevel: 0,
					},
					{
						File:             &models.CommitFile{Path: "dir1/b"},
						Path:             "dir1/b",
						CompressionLevel: 0,
					},
				},
			},
		},
		{
			name: "paths that can be compressed",
			files: []*models.CommitFile{
				{
					Path: "dir1/a",
				},
				{
					Path: "dir2/b",
				},
			},
			expected: &Node[models.CommitFile]{
				Path: "",
				Children: []*Node[models.CommitFile]{
					{
						File:             &models.CommitFile{Path: "dir1/a"},
						Path:             "dir1/a",
						CompressionLevel: 0,
					},
					{
						File:             &models.CommitFile{Path: "dir2/b"},
						Path:             "dir2/b",
						CompressionLevel: 0,
					},
				},
			},
		},
		{
			name: "paths that can be sorted",
			files: []*models.CommitFile{
				{
					Path: "b",
				},
				{
					Path: "a",
				},
			},
			expected: &Node[models.CommitFile]{
				Path: "",
				Children: []*Node[models.CommitFile]{
					{
						File: &models.CommitFile{Path: "a"},
						Path: "a",
					},
					{
						File: &models.CommitFile{Path: "b"},
						Path: "b",
					},
				},
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			result := BuildFlatTreeFromCommitFiles(s.files)
			assert.EqualValues(t, s.expected, result)
		})
	}
}
