package filetree

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/stretchr/testify/assert"
)

func TestCompress(t *testing.T) {
	scenarios := []struct {
		name     string
		root     *Node[models.File]
		expected *Node[models.File]
	}{
		{
			name:     "nil node",
			root:     nil,
			expected: nil,
		},
		{
			name: "leaf node",
			root: &Node[models.File]{
				path: "",
				Children: []*Node[models.File]{
					{File: &models.File{Path: "test", ShortStatus: " M", HasStagedChanges: true}, path: "test"},
				},
			},
			expected: &Node[models.File]{
				path: "",
				Children: []*Node[models.File]{
					{File: &models.File{Path: "test", ShortStatus: " M", HasStagedChanges: true}, path: "test"},
				},
			},
		},
		{
			name: "big example",
			root: &Node[models.File]{
				path: "",
				Children: []*Node[models.File]{
					{
						path: "dir1",
						Children: []*Node[models.File]{
							{
								File: &models.File{Path: "file2", ShortStatus: "M ", HasUnstagedChanges: true},
								path: "dir1/file2",
							},
						},
					},
					{
						path: "dir2",
						Children: []*Node[models.File]{
							{
								File: &models.File{Path: "file3", ShortStatus: " M", HasStagedChanges: true},
								path: "dir2/file3",
							},
							{
								File: &models.File{Path: "file4", ShortStatus: "M ", HasUnstagedChanges: true},
								path: "dir2/file4",
							},
						},
					},
					{
						path: "dir3",
						Children: []*Node[models.File]{
							{
								path: "dir3/dir3-1",
								Children: []*Node[models.File]{
									{
										File: &models.File{Path: "file5", ShortStatus: "M ", HasUnstagedChanges: true},
										path: "dir3/dir3-1/file5",
									},
								},
							},
						},
					},
					{
						File: &models.File{Path: "file1", ShortStatus: "M ", HasUnstagedChanges: true},
						path: "file1",
					},
				},
			},
			expected: &Node[models.File]{
				path: "",
				Children: []*Node[models.File]{
					{
						path: "dir1",
						Children: []*Node[models.File]{
							{
								File: &models.File{Path: "file2", ShortStatus: "M ", HasUnstagedChanges: true},
								path: "dir1/file2",
							},
						},
					},
					{
						path: "dir2",
						Children: []*Node[models.File]{
							{
								File: &models.File{Path: "file3", ShortStatus: " M", HasStagedChanges: true},
								path: "dir2/file3",
							},
							{
								File: &models.File{Path: "file4", ShortStatus: "M ", HasUnstagedChanges: true},
								path: "dir2/file4",
							},
						},
					},
					{
						path:             "dir3/dir3-1",
						CompressionLevel: 1,
						Children: []*Node[models.File]{
							{
								File: &models.File{Path: "file5", ShortStatus: "M ", HasUnstagedChanges: true},
								path: "dir3/dir3-1/file5",
							},
						},
					},
					{
						File: &models.File{Path: "file1", ShortStatus: "M ", HasUnstagedChanges: true},
						path: "file1",
					},
				},
			},
		},
	}

	for _, s := range scenarios {
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
			viewModel: NewFileTree(func() []*models.File { return []*models.File{{Path: "blah/one"}, {Path: "blah/two"}} }, nil, false),
			path:      "blah/two",
			expected:  &models.File{Path: "blah/two"},
		},
		{
			name:      "not found",
			viewModel: NewFileTree(func() []*models.File { return []*models.File{{Path: "blah/one"}, {Path: "blah/two"}} }, nil, false),
			path:      "blah/three",
			expected:  nil,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			assert.EqualValues(t, s.expected, s.viewModel.GetFile(s.path))
		})
	}
}
