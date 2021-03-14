package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompress(t *testing.T) {
	scenarios := []struct {
		name     string
		root     *StatusLineNode
		expected *StatusLineNode
	}{
		{
			name:     "nil node",
			root:     nil,
			expected: nil,
		},
		{
			name: "leaf node",
			root: &StatusLineNode{
				Name: "",
				Children: []*StatusLineNode{
					{File: &File{Name: "test", ShortStatus: " M", HasStagedChanges: true}, Name: "test"},
				},
			},
			expected: &StatusLineNode{
				Name: "",
				Children: []*StatusLineNode{
					{File: &File{Name: "test", ShortStatus: " M", HasStagedChanges: true}, Name: "test"},
				},
			},
		},
		{
			name: "big example",
			root: &StatusLineNode{
				Name: "",
				Children: []*StatusLineNode{
					{
						Name: "dir1",
						Path: "dir1",
						Children: []*StatusLineNode{
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
						Children: []*StatusLineNode{
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
						File: &File{Name: "file1", ShortStatus: "M ", HasUnstagedChanges: true},
						Name: "file1",
						Path: "file1",
					},
				},
			},
			expected: &StatusLineNode{
				Name: "",
				Children: []*StatusLineNode{
					{
						Name: "dir1/file2",
						File: &File{Name: "file2", ShortStatus: "M ", HasUnstagedChanges: true},
						Path: "dir1/file2",
					},
					{
						Name: "dir2",
						Path: "dir2",
						Children: []*StatusLineNode{
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
