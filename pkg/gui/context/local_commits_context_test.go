package context

import (
	"testing"
	"time"

	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/stretchr/testify/assert"
)

func TestAddCommitDropIndicator(t *testing.T) {
	pendingHeader := &NonModelItem{Index: 0, Content: "pending"}
	commitsHeader := &NonModelItem{Index: 3, Content: "commits"}
	indicator := &commitDropIndicator{insertionIndex: 3}
	spinnerConfig := config.SpinnerConfig{Frames: []string{"one", "two"}, Rate: 100}

	items := addCommitDropIndicator(
		[]*NonModelItem{pendingHeader}, indicator, "drop here", "moving commits here", spinnerConfig, time.UnixMilli(0),
	)
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

func TestAddMovingCommitsIndicator(t *testing.T) {
	items := addCommitDropIndicator(
		nil,
		&commitDropIndicator{insertionIndex: 2, moving: true},
		"drop here",
		"moving commits here",
		config.SpinnerConfig{Frames: []string{"one", "two"}, Rate: 100},
		time.UnixMilli(100),
	)

	assert.Equal(t, []*NonModelItem{
		{
			Index:   2,
			Content: style.FgCyan.SetBold().Sprint("━━━━━━ moving commits here two ━━━━━━"),
			Column:  6,
		},
	}, items)
}
