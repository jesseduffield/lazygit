package context

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/stretchr/testify/assert"
)

func TestAddCommitDropIndicator(t *testing.T) {
	pendingHeader := &NonModelItem{Index: 0, Content: "pending"}
	commitsHeader := &NonModelItem{Index: 3, Content: "commits"}
	indicator := &commitDropIndicator{insertionIndex: 3}

	items := addCommitDropIndicator([]*NonModelItem{pendingHeader}, indicator, "drop here")
	items = append(items, commitsHeader)

	assert.Equal(t, []*NonModelItem{
		pendingHeader,
		{
			Index:   3,
			Content: style.FgCyan.SetBold().Sprint("━━━━━━ drop here ━━━━━━"),
			Column:  6,
		},
		commitsHeader,
	}, items)
	assert.Equal(t, 6, modelIndexToViewIndex(4, items, 3))
	assert.Equal(t, 3, viewIndexToModelIndex(4, items, 4))
}
