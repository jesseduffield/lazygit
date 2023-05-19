package presentation

import (
	"strings"
	"testing"

	"github.com/gookit/color"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/xo/terminfo"
)

func init() {
	color.ForceSetColorLevel(terminfo.ColorLevelNone)
}

func toStringSlice(str string) []string {
	return strings.Split(strings.TrimSpace(str), "\n")
}

func TestRenderFileTree(t *testing.T) {
	scenarios := []struct {
		name           string
		root           *filetree.FileNode
		files          []*models.File
		collapsedPaths []string
		expected       []string
	}{
		{
			name:     "nil node",
			files:    nil,
			expected: []string{},
		},
		{
			name: "leaf node",
			files: []*models.File{
				{Name: "test", ShortStatus: " M", HasStagedChanges: true},
			},
			expected: []string{" M test"},
		},
		{
			name: "big example",
			files: []*models.File{
				{Name: "dir1/file2", ShortStatus: "M ", HasUnstagedChanges: true},
				{Name: "dir1/file3", ShortStatus: "M ", HasUnstagedChanges: true},
				{Name: "dir2/dir2/file3", ShortStatus: " M", HasStagedChanges: true},
				{Name: "dir2/dir2/file4", ShortStatus: "M ", HasUnstagedChanges: true},
				{Name: "dir2/file5", ShortStatus: "M ", HasUnstagedChanges: true},
				{Name: "file1", ShortStatus: "M ", HasUnstagedChanges: true},
			},
			expected: toStringSlice(
				`
▶ dir1
▼ dir2
  ▼ dir2
     M file3
    M  file4
  M  file5
M  file1
`,
			),
			collapsedPaths: []string{"dir1"},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.name, func(t *testing.T) {
			viewModel := filetree.NewFileTree(func() []*models.File { return s.files }, utils.NewDummyLog(), true)
			viewModel.SetTree()
			for _, path := range s.collapsedPaths {
				viewModel.ToggleCollapsed(path)
			}
			result := RenderFileTree(viewModel, "", nil)
			assert.EqualValues(t, s.expected, result)
		})
	}
}

func TestRenderCommitFileTree(t *testing.T) {
	scenarios := []struct {
		name           string
		root           *filetree.FileNode
		files          []*models.CommitFile
		collapsedPaths []string
		expected       []string
	}{
		{
			name:     "nil node",
			files:    nil,
			expected: []string{},
		},
		{
			name: "leaf node",
			files: []*models.CommitFile{
				{Name: "test", ChangeStatus: "A"},
			},
			expected: []string{"A test"},
		},
		{
			name: "big example",
			files: []*models.CommitFile{
				{Name: "dir1/file2", ChangeStatus: "M"},
				{Name: "dir1/file3", ChangeStatus: "A"},
				{Name: "dir2/dir2/file3", ChangeStatus: "D"},
				{Name: "dir2/dir2/file4", ChangeStatus: "M"},
				{Name: "dir2/file5", ChangeStatus: "M"},
				{Name: "file1", ChangeStatus: "M"},
			},
			expected: toStringSlice(
				`
▶ dir1
▼ dir2
  ▼ dir2
    D file3
    M file4
  M file5
M file1
`,
			),
			collapsedPaths: []string{"dir1"},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.name, func(t *testing.T) {
			viewModel := filetree.NewCommitFileTreeViewModel(func() []*models.CommitFile { return s.files }, utils.NewDummyLog(), true)
			viewModel.SetRef(&models.Commit{})
			viewModel.SetTree()
			for _, path := range s.collapsedPaths {
				viewModel.ToggleCollapsed(path)
			}
			patchBuilder := patch.NewPatchBuilder(
				utils.NewDummyLog(),
				func(from string, to string, reverse bool, filename string, plain bool) (string, error) {
					return "", nil
				},
			)
			patchBuilder.Start("from", "to", false, false)
			result := RenderCommitFileTree(viewModel, "", patchBuilder)
			assert.EqualValues(t, s.expected, result)
		})
	}
}
