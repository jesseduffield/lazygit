package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompress(t *testing.T) {
	scenarios := []struct {
		name     string
		root     *FileChangeNode
		expected *FileChangeNode
	}{
		{
			name:     "nil node",
			root:     nil,
			expected: nil,
		},
		{
			name: "leaf node",
			root: &FileChangeNode{
				Name: "",
				Children: []*FileChangeNode{
					{File: &File{Name: "test", ShortStatus: " M", HasStagedChanges: true}, Name: "test"},
				},
			},
			expected: &FileChangeNode{
				Name: "",
				Children: []*FileChangeNode{
					{File: &File{Name: "test", ShortStatus: " M", HasStagedChanges: true}, Name: "test"},
				},
			},
		},
		{
			name: "big example",
			root: &FileChangeNode{
				Name: "",
				Children: []*FileChangeNode{
					{
						Name: "dir1",
						Path: "dir1",
						Children: []*FileChangeNode{
							{
								File: &File{Name: "file2", ShortStatus: "M ", HasUnstagedChanges: true},
								Name: "file2",
								Path: "dir1/file2",
							},
						},
					},
					{
						Name: "dir2",
						Path: "dir2",
						Children: []*FileChangeNode{
							{
								File: &File{Name: "file3", ShortStatus: " M", HasStagedChanges: true},
								Name: "file3",
								Path: "dir2/file3",
							},
							{
								File: &File{Name: "file4", ShortStatus: "M ", HasUnstagedChanges: true},
								Name: "file4",
								Path: "dir2/file4",
							},
						},
					},
					{
						Name: "dir3",
						Path: "dir3",
						Children: []*FileChangeNode{
							{
								Name: "dir3-1",
								Path: "dir3/dir3-1",
								Children: []*FileChangeNode{
									{
										File: &File{Name: "file5", ShortStatus: "M ", HasUnstagedChanges: true},
										Name: "file5",
										Path: "dir3/dir3-1/file5",
									},
								},
							},
						},
					},
					{
						File: &File{Name: "file1", ShortStatus: "M ", HasUnstagedChanges: true},
						Name: "file1",
						Path: "file1",
					},
				},
			},
			expected: &FileChangeNode{
				Name: "",
				Children: []*FileChangeNode{
					{
						Name: "dir1/file2",
						File: &File{Name: "file2", ShortStatus: "M ", HasUnstagedChanges: true},
						Path: "dir1/file2",
					},
					{
						Name: "dir2",
						Path: "dir2",
						Children: []*FileChangeNode{
							{
								File: &File{Name: "file3", ShortStatus: " M", HasStagedChanges: true},
								Name: "file3",
								Path: "dir2/file3",
							},
							{
								File: &File{Name: "file4", ShortStatus: "M ", HasUnstagedChanges: true},
								Name: "file4",
								Path: "dir2/file4",
							},
						},
					},
					{
						Name: "dir3/dir3-1/file5",
						File: &File{Name: "file5", ShortStatus: "M ", HasUnstagedChanges: true},
						Path: "dir3/dir3-1/file5",
					},
					{
						File: &File{Name: "file1", ShortStatus: "M ", HasUnstagedChanges: true},
						Name: "file1",
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
