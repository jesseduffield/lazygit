package filetree

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/common"
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
