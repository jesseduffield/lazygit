package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {
	scenarios := []struct {
		name     string
		root     *StatusLineNode
		expected []string
	}{
		{
			name:     "nil node",
			root:     nil,
			expected: []string{},
		},
		{
			name: "leaf node",
			root: &StatusLineNode{
				Name: "",
				Children: []*StatusLineNode{
					{File: &File{Name: "test", ShortStatus: " M", HasStagedChanges: true}, Name: "test"},
				},
			},
			expected: []string{" M test"},
		},
		{
			name: "big example",
			root: &StatusLineNode{
				Name: "",
				Children: []*StatusLineNode{
					{
						Name:      "dir1",
						Collapsed: true,
						Children: []*StatusLineNode{
							{
								File: &File{Name: "file2", ShortStatus: "M ", HasUnstagedChanges: true},
								Name: "file2",
							},
						},
					},
					{
						Name: "dir2",
						Children: []*StatusLineNode{
							{
								File: &File{Name: "file3", ShortStatus: " M", HasStagedChanges: true},
								Name: "file3",
							},
							{
								File: &File{Name: "file4", ShortStatus: "M ", HasUnstagedChanges: true},
								Name: "file4",
							},
						},
					},
					{
						File: &File{Name: "file1", ShortStatus: "M ", HasUnstagedChanges: true},
						Name: "file1",
					},
				},
			},

			expected: []string{"M  dir1 ►", "MM dir2 ▼", "   M file3", "  M  file4", "M  file1"},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.name, func(t *testing.T) {
			result := s.root.Render()[1:]
			assert.EqualValues(t, s.expected, result)
		})
	}
}
