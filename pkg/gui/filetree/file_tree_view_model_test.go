package filetree

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/config"
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
			/* EXPECTED:
			expectedPath: "dir/new.go",
			ACTUAL: */
			expectedPath: "other.go",
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
