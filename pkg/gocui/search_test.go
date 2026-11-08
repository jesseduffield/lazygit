package gocui

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gdamore/tcell/v3/color"
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

func TestSearchHighlights(t *testing.T) {
	for _, test := range []struct {
		name       string
		content    string
		search     string
		selected   int
		originX    int
		originY    int
		wrap       bool
		model      []SearchPosition
		wantColors []string
	}{
		{name: "no matches", content: "some text", search: "match", wantColors: []string{"........", "........"}},
		{name: "several matches", content: "xx xx\nxx", search: "xx", selected: 1, wantColors: []string{"yy.cc...", "yy......"}},
		{name: "overlapping matches", content: "aaa", search: "aa", selected: 1, wantColors: []string{"yyc.....", "........"}},
		{name: "scrolling", content: "other\nxx match\nmatch", search: "match", selected: 1, originX: 3, originY: 1, wantColors: []string{"yyyyy...", "cc......"}},
		{name: "wrapping", content: "12345678match", search: "match", wrap: true, wantColors: []string{"........", "ccccc..."}},
		{name: "Unicode columns", content: "界ématch", search: "match", wantColors: []string{"...ccccc", "........"}},
		{name: "abbreviated model match", content: "abc\nother", search: "abcdef", model: []SearchPosition{{Y: 0, XStart: 0, XEnd: 3}}, wantColors: []string{"ccc.....", "........"}},
		{name: "unrendered model match", content: "abc", search: "abcdef", model: []SearchPosition{{Y: 0, XStart: -1, XEnd: -1}}, wantColors: []string{"........", "........"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			WithSimulationScreen(t, 10, 4)
			v := NewView("diff", 0, 0, 9, 3, OutputNormal)
			v.Wrap = test.wrap
			v.SetContent(test.content)
			v.UpdateSearchResults(test.search, test.model)
			v.searcher.currentSearchIndex = test.selected
			v.SetOrigin(test.originX, test.originY)
			v.draw(true)

			for y, want := range test.wantColors {
				var colors strings.Builder
				for x := range 8 {
					_, style, _ := Screen.Get(x+1, y+1)
					switch style.GetBackground() {
					case color.Olive:
						colors.WriteByte('y')
					case color.Teal:
						colors.WriteByte('c')
					default:
						colors.WriteByte('.')
					}
				}
				assert.Equal(t, want, colors.String())
			}
		})
	}
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
