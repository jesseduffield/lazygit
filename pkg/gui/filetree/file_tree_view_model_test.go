package filetree

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestSetStatusFilterPreservingSelection(t *testing.T) {
	files := []*models.File{
		{Path: "file1"},
		{Path: "file2", HasMergeConflicts: true},
		{Path: "file3", HasMergeConflicts: true},
	}
	viewModel := NewFileTreeViewModel(
		func() []*models.File { return files },
		common.NewDummyCommon(),
		false,
	)
	viewModel.SetTree()
	viewModel.SetStatusFilter(DisplayConflicted)
	viewModel.SetSelection(viewModel.Len() - 2)
	viewModel.ToggleStickyRange()
	viewModel.MoveSelectedLine(1)

	viewModel.SetStatusFilterPreservingSelection(DisplayAll)

	assert.Equal(t, "file3", viewModel.GetSelectedPath())
	assert.False(t, viewModel.IsSelectingRange())
}

func TestSetTreeSelectsNewFileWhenSelectedRenameSplits(t *testing.T) {
	scenarios := []struct {
		name         string
		showRootItem bool
		expectedPath string
	}{
		{
			name:         "with root item",
			showRootItem: true,
			expectedPath: "dir/new.go",
		},
		{
			name:         "without root item",
			showRootItem: false,
			expectedPath: "dir/new.go",
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			userConfig := config.GetDefaultConfig()
			userConfig.Gui.ShowRootItemInFileTree = s.showRootItem
			cmn := common.NewDummyCommonWithUserConfigAndAppState(userConfig, nil)

			files := []*models.File{
				{Path: "dir/new.go", PreviousPath: "dir/old.go"},
				{Path: "other.go"},
			}
			viewModel := NewFileTreeViewModel(func() []*models.File { return files }, cmn, true)
			viewModel.SetTree()
			idx, found := viewModel.GetIndexForPath(InternalTreePathForFilePath("dir/new.go", s.showRootItem))
			assert.True(t, found)
			viewModel.SetSelection(idx)

			// the rename is split into its two halves, e.g. because it was unstaged
			files = []*models.File{
				{Path: "dir/new.go"},
				{Path: "dir/old.go"},
				{Path: "other.go"},
			}
			viewModel.SetTree()

			assert.Equal(t, s.expectedPath, viewModel.GetSelectedPath())
		})
	}
}

func TestSetTreeFollowsRenameIntoCollapsedDir(t *testing.T) {
	scenarios := []struct {
		name         string
		showRootItem bool
		expectedPath string
	}{
		{
			name:         "with root item",
			showRootItem: true,
			expectedPath: "a/new.go",
		},
		{
			name:         "without root item",
			showRootItem: false,
			expectedPath: "a/new.go",
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			userConfig := config.GetDefaultConfig()
			userConfig.Gui.ShowRootItemInFileTree = s.showRootItem
			cmn := common.NewDummyCommonWithUserConfigAndAppState(userConfig, nil)

			files := []*models.File{
				{Path: "a/b.go"},
				{Path: "a/new.go"},
				{Path: "old.go"},
			}
			viewModel := NewFileTreeViewModel(func() []*models.File { return files }, cmn, true)
			viewModel.SetTree()
			viewModel.ToggleCollapsed(InternalTreePathForFilePath("a", s.showRootItem))
			idx, found := viewModel.GetIndexForPath(InternalTreePathForFilePath("old.go", s.showRootItem))
			assert.True(t, found)
			viewModel.SetSelection(idx)

			// staging the deletion of old.go turns it into the old half of a rename
			files = []*models.File{
				{Path: "a/b.go"},
				{Path: "a/new.go", PreviousPath: "old.go"},
			}
			viewModel.SetTree()

			assert.Equal(t, s.expectedPath, viewModel.GetSelectedPath())
		})
	}
}

func TestSetTreeKeepsSelectionAcrossCompressionChanges(t *testing.T) {
	scenarios := []struct {
		name         string
		filesBefore  []string
		selectedPath string
		filesAfter   []string
		expectedPath string
	}{
		{
			name:         "compressed root directory splits",
			filesBefore:  []string{"pkg/gui/controllers/helpers/refresh_helper.go"},
			selectedPath: "pkg/gui/controllers/helpers",
			filesAfter: []string{
				"pkg/gui/context/base_context.go",
				"pkg/gui/controllers/helpers/refresh_helper.go",
			},
			expectedPath: "pkg/gui",
		},
		{
			name:         "compressed subdirectory splits",
			filesBefore:  []string{"a/b/c/file1", "file2"},
			selectedPath: "a/b/c",
			filesAfter:   []string{"a/b/c/file1", "a/b/d/file3", "file2"},
			expectedPath: "a/b",
		},
		{
			name:         "file inside a compressed directory that splits",
			filesBefore:  []string{"pkg/gui/controllers/helpers/refresh_helper.go"},
			selectedPath: "pkg/gui/controllers/helpers/refresh_helper.go",
			filesAfter: []string{
				"pkg/gui/context/base_context.go",
				"pkg/gui/controllers/helpers/refresh_helper.go",
			},
			expectedPath: "pkg/gui/controllers/helpers/refresh_helper.go",
		},
		{
			name: "directories merge into one compressed node",
			filesBefore: []string{
				"pkg/gui/context/base_context.go",
				"pkg/gui/controllers/helpers/refresh_helper.go",
			},
			selectedPath: "pkg/gui",
			filesAfter:   []string{"pkg/gui/controllers/helpers/refresh_helper.go"},
			expectedPath: "pkg/gui/controllers/helpers",
		},
	}

	toFiles := func(paths []string) []*models.File {
		return lo.Map(paths, func(path string, _ int) *models.File {
			return &models.File{Path: path}
		})
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			files := toFiles(s.filesBefore)
			cmn := common.NewDummyCommon()
			viewModel := NewFileTreeViewModel(
				func() []*models.File { return files },
				cmn,
				true,
			)
			viewModel.SetTree()
			showRootItem := cmn.UserConfig().Gui.ShowRootItemInFileTree
			idx, found := viewModel.GetIndexForPath(InternalTreePathForFilePath(s.selectedPath, showRootItem))
			assert.True(t, found)
			viewModel.SetSelection(idx)

			files = toFiles(s.filesAfter)
			viewModel.SetTree()

			assert.Equal(t, s.expectedPath, viewModel.GetSelectedPath())
		})
	}
}
