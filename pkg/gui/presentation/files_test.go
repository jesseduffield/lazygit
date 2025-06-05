package presentation

import (
	"strings"
	"testing"

	"github.com/gookit/color"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/xo/terminfo"
)

func toStringSlice(str string) []string {
	return strings.Split(strings.TrimSpace(str), "\n")
}

func TestRenderFileTree(t *testing.T) {
	scenarios := []struct {
		name            string
		root            *filetree.FileNode
		files           []*models.File
		collapsedPaths  []string
		showLineChanges bool
		showRootItem    bool
		expected        []string
	}{
		{
			name:     "nil node",
			files:    nil,
			expected: []string{},
		},
		{
			name: "leaf node",
			files: []*models.File{
				{Path: "test", ShortStatus: " M", HasStagedChanges: true},
			},
			showRootItem: true,
			expected:     []string{" M test"},
		},
		{
			name: "numstat",
			files: []*models.File{
				{Path: "test", ShortStatus: " M", HasStagedChanges: true, LinesAdded: 1, LinesDeleted: 1},
				{Path: "test2", ShortStatus: " M", HasStagedChanges: true, LinesAdded: 1},
				{Path: "test3", ShortStatus: " M", HasStagedChanges: true, LinesDeleted: 1},
				{Path: "test4", ShortStatus: " M", HasStagedChanges: true, LinesAdded: 0, LinesDeleted: 0},
			},
			showLineChanges: true,
			showRootItem:    true,
			expected: []string{
				"▼ /",
				"   M test +1 -1",
				"   M test2 +1",
				"   M test3 -1",
				"   M test4",
			},
		},
		{
			name: "big example",
			files: []*models.File{
				{Path: "dir1/file2", ShortStatus: "M ", HasUnstagedChanges: true},
				{Path: "dir1/file3", ShortStatus: "M ", HasUnstagedChanges: true},
				{Path: "dir2/dir2/file3", ShortStatus: " M", HasStagedChanges: true},
				{Path: "dir2/dir2/file4", ShortStatus: "M ", HasUnstagedChanges: true},
				{Path: "dir2/file5", ShortStatus: "M ", HasUnstagedChanges: true},
				{Path: "file1", ShortStatus: "M ", HasUnstagedChanges: true},
			},
			showRootItem: true,
			expected: toStringSlice(
				`
▼ /
  ▶ dir1
  ▼ dir2
    ▼ dir2
       M file3
      M  file4
    M  file5
  M  file1
`,
			),
			collapsedPaths: []string{"./dir1"},
		},
		{
			name: "big example without root item",
			files: []*models.File{
				{Path: "dir1/file2", ShortStatus: "M ", HasUnstagedChanges: true},
				{Path: "dir1/file3", ShortStatus: "M ", HasUnstagedChanges: true},
				{Path: "dir2/dir2/file3", ShortStatus: " M", HasStagedChanges: true},
				{Path: "dir2/dir2/file4", ShortStatus: "M ", HasUnstagedChanges: true},
				{Path: "dir2/file5", ShortStatus: "M ", HasUnstagedChanges: true},
				{Path: "file1", ShortStatus: "M ", HasUnstagedChanges: true},
			},
			showRootItem: false,
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

	oldColorLevel := color.ForceSetColorLevel(terminfo.ColorLevelNone)
	defer color.ForceSetColorLevel(oldColorLevel)

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			common := common.NewDummyCommon()
			common.UserConfig().Gui.ShowRootItemInFileTree = s.showRootItem
			viewModel := filetree.NewFileTree(func() []*models.File { return s.files }, common, true)
			viewModel.SetTree()
			for _, path := range s.collapsedPaths {
				viewModel.ToggleCollapsed(path)
			}
			result := RenderFileTree(viewModel, nil, false, s.showLineChanges, &config.CustomIconsConfig{}, s.showRootItem)
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
		showRootItem   bool
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
				{Path: "test", ChangeStatus: "A"},
			},
			showRootItem: true,
			expected:     []string{"A test"},
		},
		{
			name: "big example",
			files: []*models.CommitFile{
				{Path: "dir1/file2", ChangeStatus: "M"},
				{Path: "dir1/file3", ChangeStatus: "A"},
				{Path: "dir2/dir2/file3", ChangeStatus: "D"},
				{Path: "dir2/dir2/file4", ChangeStatus: "M"},
				{Path: "dir2/file5", ChangeStatus: "M"},
				{Path: "file1", ChangeStatus: "M"},
			},
			showRootItem: true,
			expected: toStringSlice(
				`
▼ /
  ▶ dir1
  ▼ dir2
    ▼ dir2
      D file3
      M file4
    M file5
  M file1
`,
			),
			collapsedPaths: []string{"./dir1"},
		},
		{
			name: "big example without root item",
			files: []*models.CommitFile{
				{Path: "dir1/file2", ChangeStatus: "M"},
				{Path: "dir1/file3", ChangeStatus: "A"},
				{Path: "dir2/dir2/file3", ChangeStatus: "D"},
				{Path: "dir2/dir2/file4", ChangeStatus: "M"},
				{Path: "dir2/file5", ChangeStatus: "M"},
				{Path: "file1", ChangeStatus: "M"},
			},
			showRootItem: false,
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

	oldColorLevel := color.ForceSetColorLevel(terminfo.ColorLevelNone)
	defer color.ForceSetColorLevel(oldColorLevel)

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			hashPool := &utils.StringPool{}

			common := common.NewDummyCommon()
			common.UserConfig().Gui.ShowRootItemInFileTree = s.showRootItem
			viewModel := filetree.NewCommitFileTreeViewModel(func() []*models.CommitFile { return s.files }, common, true)
			viewModel.SetRef(models.NewCommit(hashPool, models.NewCommitOpts{Hash: "1234"}))
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
			result := RenderCommitFileTree(viewModel, patchBuilder, false, &config.CustomIconsConfig{})
			assert.EqualValues(t, s.expected, result)
		})
	}
}
