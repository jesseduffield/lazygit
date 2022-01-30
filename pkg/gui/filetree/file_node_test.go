package filetree

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/stretchr/testify/assert"
)

func TestCompress(t *testing.T) {
	scenarios := []struct {
		name     string
		root     *FileNode
		expected *FileNode
	}{
		{
			name:     "nil node",
			root:     nil,
			expected: nil,
		},
		{
			name: "leaf node",
			root: &FileNode{
				Path: "",
				Children: []*FileNode{
					{File: &models.File{Name: "test", ShortStatus: " M", HasStagedChanges: true}, Path: "test"},
				},
			},
			expected: &FileNode{
				Path: "",
				Children: []*FileNode{
					{File: &models.File{Name: "test", ShortStatus: " M", HasStagedChanges: true}, Path: "test"},
				},
			},
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
								File: &models.File{Name: "file2", ShortStatus: "M ", HasUnstagedChanges: true},
								Path: "dir1/file2",
							},
						},
					},
					{
						Path: "dir2",
						Children: []*FileNode{
							{
								File: &models.File{Name: "file3", ShortStatus: " M", HasStagedChanges: true},
								Path: "dir2/file3",
							},
							{
								File: &models.File{Name: "file4", ShortStatus: "M ", HasUnstagedChanges: true},
								Path: "dir2/file4",
							},
						},
					},
					{
						Path: "dir3",
						Children: []*FileNode{
							{
								Path: "dir3/dir3-1",
								Children: []*FileNode{
									{
										File: &models.File{Name: "file5", ShortStatus: "M ", HasUnstagedChanges: true},
										Path: "dir3/dir3-1/file5",
									},
								},
							},
						},
					},
					{
						File: &models.File{Name: "file1", ShortStatus: "M ", HasUnstagedChanges: true},
						Path: "file1",
					},
				},
			},
			expected: &FileNode{
				Path: "",
				Children: []*FileNode{
					{
						Path: "dir1",
						Children: []*FileNode{
							{
								File: &models.File{Name: "file2", ShortStatus: "M ", HasUnstagedChanges: true},
								Path: "dir1/file2",
							},
						},
					},
					{
						Path: "dir2",
						Children: []*FileNode{
							{
								File: &models.File{Name: "file3", ShortStatus: " M", HasStagedChanges: true},
								Path: "dir2/file3",
							},
							{
								File: &models.File{Name: "file4", ShortStatus: "M ", HasUnstagedChanges: true},
								Path: "dir2/file4",
							},
						},
					},
					{
						Path:             "dir3/dir3-1",
						CompressionLevel: 1,
						Children: []*FileNode{
							{
								File: &models.File{Name: "file5", ShortStatus: "M ", HasUnstagedChanges: true},
								Path: "dir3/dir3-1/file5",
							},
						},
					},
					{
						File: &models.File{Name: "file1", ShortStatus: "M ", HasUnstagedChanges: true},
						Path: "file1",
					},
				},
			},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.name, func(t *testing.T) {
			s.root.Compress()
			assert.EqualValues(t, s.expected, s.root)
		})
	}
}

func TestGetFile(t *testing.T) {
	scenarios := []struct {
		name      string
		viewModel *FileTree
		path      string
		expected  *models.File
	}{
		{
			name:      "valid case",
			viewModel: NewFileTree(func() []*models.File { return []*models.File{{Name: "blah/one"}, {Name: "blah/two"}} }, nil, false),
			path:      "blah/two",
			expected:  &models.File{Name: "blah/two"},
		},
		{
			name:      "not found",
			viewModel: NewFileTree(func() []*models.File { return []*models.File{{Name: "blah/one"}, {Name: "blah/two"}} }, nil, false),
			path:      "blah/three",
			expected:  nil,
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.name, func(t *testing.T) {
			assert.EqualValues(t, s.expected, s.viewModel.GetFile(s.path))
		})
	}
}
