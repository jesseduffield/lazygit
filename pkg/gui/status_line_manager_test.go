package gui

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {
	scenarios := []struct {
		name     string
		root     *models.StatusLineNode
		expected []string
	}{
		{
			name:     "nil node",
			root:     nil,
			expected: []string{},
		},
		{
			name: "leaf node",
			root: &models.StatusLineNode{
				Name: "",
				Children: []*models.StatusLineNode{
					{File: &models.File{Name: "test", ShortStatus: " M", HasStagedChanges: true}, Name: "test"},
				},
			},
			expected: []string{" M test"},
		},
		{
			name: "big example",
			root: &models.StatusLineNode{
				Name: "",
				Children: []*models.StatusLineNode{
					{
						Name:      "dir1",
						Collapsed: true,
						Children: []*models.StatusLineNode{
							{
								File: &models.File{Name: "file2", ShortStatus: "M ", HasUnstagedChanges: true},
								Name: "file2",
							},
						},
					},
					{
						Name: "dir2",
						Children: []*models.StatusLineNode{
							{
								Name: "dir2",
								Children: []*models.StatusLineNode{
									{
										File: &models.File{Name: "file3", ShortStatus: " M", HasStagedChanges: true},
										Name: "file3",
									},
									{
										File: &models.File{Name: "file4", ShortStatus: "M ", HasUnstagedChanges: true},
										Name: "file4",
									},
								},
							},
							{
								File: &models.File{Name: "file5", ShortStatus: "M ", HasUnstagedChanges: true},
								Name: "file5",
							},
						},
					},
					{
						File: &models.File{Name: "file1", ShortStatus: "M ", HasUnstagedChanges: true},
						Name: "file1",
					},
				},
			},

			expected: []string{" M dir1 ►", "MM dir2 ▼", "├─ MM dir2 ▼", "│   ├─  M file3", "│   └─ M  file4", "└─ M  file5", "M  file1"},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.name, func(t *testing.T) {
			mngr := &StatusLineManager{Tree: s.root}
			result := mngr.Render("", nil)
			assert.EqualValues(t, s.expected, result)
		})
	}
}
