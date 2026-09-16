package filetree

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/stretchr/testify/assert"
)

// When the tree shrinks under the selection - e.g. moving a patch out into the
// index removes a file - SetTree must keep the selection in range. Otherwise
// GetSelectedItems returns a nil node, which crashes callers such as
// canEditFiles when the options map is rendered during layout.
func TestCommitFileTreeViewModelSetTreeClampsSelectionOnShrink(t *testing.T) {
	files := []*models.CommitFile{
		{Path: "file1"},
		{Path: "file2"},
		{Path: "file3"},
	}
	viewModel := NewCommitFileTreeViewModel(
		func() []*models.CommitFile { return files },
		common.NewDummyCommon(),
		false, // flat list
	)
	viewModel.SetTree()
	viewModel.SetSelectedLineIdx(viewModel.Len() - 1)

	// The file under the cursor goes away and the tree shrinks.
	files = []*models.CommitFile{{Path: "file1"}}
	viewModel.SetTree()

	assert.Less(t, viewModel.GetSelectedLineIdx(), viewModel.Len())
	assert.NotNil(t, viewModel.GetSelected())
	items, _, _ := viewModel.GetSelectedItems()
	assert.NotEmpty(t, items)
	for _, item := range items {
		assert.NotNil(t, item)
	}
}
