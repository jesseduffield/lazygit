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
		expected *FileNode
	}{
		{
			name:  "no files",
			files: []*models.File{},
			expected: &FileNode{
				Path:     "",
				Children: []*FileNode{},
			},
		},
		{
			name: "files in same directory",
			files: []*models.File{
				{
					Name: "dir1/a",
				},
				{
					Name: "dir1/b",
				},
			},
			expected: &FileNode{
				Path: "",
				Children: []*FileNode{
					{
						Path: "dir1",
						Children: []*FileNode{
							{
								File: &models.File{Name: "dir1/a"},
								Path: "dir1/a",
							},
							{
								File: &models.File{Name: "dir1/b"},
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
					Name: "dir1/dir3/a",
				},
				{
					Name: "dir2/dir4/b",
				},
			},
			expected: &FileNode{
				Path: "",
				Children: []*FileNode{
					{
						Path: "dir1/dir3",
						Children: []*FileNode{
							{
								File: &models.File{Name: "dir1/dir3/a"},
								Path: "dir1/dir3/a",
							},
						},
						CompressionLevel: 1,
					},
					{
						Path: "dir2/dir4",
						Children: []*FileNode{
							{
								File: &models.File{Name: "dir2/dir4/b"},
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
					Name: "b",
				},
				{
					Name: "a",
				},
			},
			expected: &FileNode{
				Path: "",
				Children: []*FileNode{
					{
						File: &models.File{Name: "a"},
						Path: "a",
					},
					{
						File: &models.File{Name: "b"},
						Path: "b",
					},
				},
			},
		},
		{
			name: "paths that can be sorted including a merge conflict file",
			files: []*models.File{
				{
					Name: "b",
				},
				{
					Name:              "z",
					HasMergeConflicts: true,
				},
				{
					Name: "a",
				},
			},
			expected: &FileNode{
				Path: "",
				// it is a little strange that we're not bubbling up our merge conflict
				// here but we are technically still in in tree mode and that's the rule
				Children: []*FileNode{
					{
						File: &models.File{Name: "a"},
						Path: "a",
					},
					{
						File: &models.File{Name: "b"},
						Path: "b",
					},
					{
						File: &models.File{Name: "z", HasMergeConflicts: true},
						Path: "z",
					},
				},
			},
		},
	}

	for _, s := range scenarios {
		s := s
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
		expected *FileNode
	}{
		{
			name:  "no files",
			files: []*models.File{},
			expected: &FileNode{
				Path:     "",
				Children: []*FileNode{},
			},
		},
		{
			name: "files in same directory",
			files: []*models.File{
				{
					Name: "dir1/a",
				},
				{
					Name: "dir1/b",
				},
			},
			expected: &FileNode{
				Path: "",
				Children: []*FileNode{
					{
						File:             &models.File{Name: "dir1/a"},
						Path:             "dir1/a",
						CompressionLevel: 0,
					},
					{
						File:             &models.File{Name: "dir1/b"},
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
					Name: "dir1/a",
				},
				{
					Name: "dir2/b",
				},
			},
			expected: &FileNode{
				Path: "",
				Children: []*FileNode{
					{
						File:             &models.File{Name: "dir1/a"},
						Path:             "dir1/a",
						CompressionLevel: 0,
					},
					{
						File:             &models.File{Name: "dir2/b"},
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
					Name: "b",
				},
				{
					Name: "a",
				},
			},
			expected: &FileNode{
				Path: "",
				Children: []*FileNode{
					{
						File: &models.File{Name: "a"},
						Path: "a",
					},
					{
						File: &models.File{Name: "b"},
						Path: "b",
					},
				},
			},
		},
		{
			name: "tracked, untracked, and conflicted files",
			files: []*models.File{
				{
					Name:    "a2",
					Tracked: false,
				},
				{
					Name:    "a1",
					Tracked: false,
				},
				{
					Name:              "c2",
					HasMergeConflicts: true,
				},
				{
					Name:              "c1",
					HasMergeConflicts: true,
				},
				{
					Name:    "b2",
					Tracked: true,
				},
				{
					Name:    "b1",
					Tracked: true,
				},
			},
			expected: &FileNode{
				Path: "",
				Children: []*FileNode{
					{
						File: &models.File{Name: "c1", HasMergeConflicts: true},
						Path: "c1",
					},
					{
						File: &models.File{Name: "c2", HasMergeConflicts: true},
						Path: "c2",
					},
					{
						File: &models.File{Name: "b1", Tracked: true},
						Path: "b1",
					},
					{
						File: &models.File{Name: "b2", Tracked: true},
						Path: "b2",
					},
					{
						File: &models.File{Name: "a1", Tracked: false},
						Path: "a1",
					},
					{
						File: &models.File{Name: "a2", Tracked: false},
						Path: "a2",
					},
				},
			},
		},
	}

	for _, s := range scenarios {
		s := s
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
		expected *CommitFileNode
	}{
		{
			name:  "no files",
			files: []*models.CommitFile{},
			expected: &CommitFileNode{
				Path:     "",
				Children: []*CommitFileNode{},
			},
		},
		{
			name: "files in same directory",
			files: []*models.CommitFile{
				{
					Name: "dir1/a",
				},
				{
					Name: "dir1/b",
				},
			},
			expected: &CommitFileNode{
				Path: "",
				Children: []*CommitFileNode{
					{
						Path: "dir1",
						Children: []*CommitFileNode{
							{
								File: &models.CommitFile{Name: "dir1/a"},
								Path: "dir1/a",
							},
							{
								File: &models.CommitFile{Name: "dir1/b"},
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
					Name: "dir1/dir3/a",
				},
				{
					Name: "dir2/dir4/b",
				},
			},
			expected: &CommitFileNode{
				Path: "",
				Children: []*CommitFileNode{
					{
						Path: "dir1/dir3",
						Children: []*CommitFileNode{
							{
								File: &models.CommitFile{Name: "dir1/dir3/a"},
								Path: "dir1/dir3/a",
							},
						},
						CompressionLevel: 1,
					},
					{
						Path: "dir2/dir4",
						Children: []*CommitFileNode{
							{
								File: &models.CommitFile{Name: "dir2/dir4/b"},
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
					Name: "b",
				},
				{
					Name: "a",
				},
			},
			expected: &CommitFileNode{
				Path: "",
				Children: []*CommitFileNode{
					{
						File: &models.CommitFile{Name: "a"},
						Path: "a",
					},
					{
						File: &models.CommitFile{Name: "b"},
						Path: "b",
					},
				},
			},
		},
	}

	for _, s := range scenarios {
		s := s
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
		expected *CommitFileNode
	}{
		{
			name:  "no files",
			files: []*models.CommitFile{},
			expected: &CommitFileNode{
				Path:     "",
				Children: []*CommitFileNode{},
			},
		},
		{
			name: "files in same directory",
			files: []*models.CommitFile{
				{
					Name: "dir1/a",
				},
				{
					Name: "dir1/b",
				},
			},
			expected: &CommitFileNode{
				Path: "",
				Children: []*CommitFileNode{
					{
						File:             &models.CommitFile{Name: "dir1/a"},
						Path:             "dir1/a",
						CompressionLevel: 0,
					},
					{
						File:             &models.CommitFile{Name: "dir1/b"},
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
					Name: "dir1/a",
				},
				{
					Name: "dir2/b",
				},
			},
			expected: &CommitFileNode{
				Path: "",
				Children: []*CommitFileNode{
					{
						File:             &models.CommitFile{Name: "dir1/a"},
						Path:             "dir1/a",
						CompressionLevel: 0,
					},
					{
						File:             &models.CommitFile{Name: "dir2/b"},
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
					Name: "b",
				},
				{
					Name: "a",
				},
			},
			expected: &CommitFileNode{
				Path: "",
				Children: []*CommitFileNode{
					{
						File: &models.CommitFile{Name: "a"},
						Path: "a",
					},
					{
						File: &models.CommitFile{Name: "b"},
						Path: "b",
					},
				},
			},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.name, func(t *testing.T) {
			result := BuildFlatTreeFromCommitFiles(s.files)
			assert.EqualValues(t, s.expected, result)
		})
	}
}
