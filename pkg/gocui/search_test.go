package gocui

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// writeLines writes the given lines to the view, as a task rendering content into it
// does: one line at a time.
func writeLines(v *View, lines ...string) {
	for _, line := range lines {
		fmt.Fprintf(v, "%s\n", line)
	}
}

func TestSearchStatusAfterTheMatchesChange(t *testing.T) {
	v := NewView("name", 0, 0, 40, 10, OutputNormal)
	writeLines(v, "match", "other", "match", "other", "match")

	v.Search("match", nil)
	_ = v.gotoNextMatch()
	_ = v.gotoNextMatch()
	index, total := v.GetSearchStatus()
	assert.Equal(t, 2, index)
	assert.Equal(t, 3, total)

	// The content is re-rendered with only the first of those matches left in it.
	v.Clear()
	writeLines(v, "match", "other", "other")

	index, total = v.GetSearchStatus()
	assert.Equal(t, 0, index)
	assert.Equal(t, 1, total)
}

func TestSearchPositionsFollowStreamedContent(t *testing.T) {
	v := NewView("name", 0, 0, 40, 10, OutputNormal)
	v.Search("match", nil)

	// A render arrives a line at a time, and the status describes all of it.
	writeLines(v, "other", "match", "other", "match")

	_, total := v.GetSearchStatus()
	assert.Equal(t, 2, total)
}

func BenchmarkWriteToSearchedView(b *testing.B) {
	for b.Loop() {
		v := NewView("name", 0, 0, 100, 40, OutputNormal)
		v.Search("match", nil)
		for i := range 2000 {
			fmt.Fprintf(v, "line %d of a diff, most of which does not match\n", i)
		}
		v.GetSearchStatus()
	}
}
