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

// drawSearchHighlights draws the view and returns how each of its rows is highlighted:
// "#" for a search match, "*" for the current one, and "." for neither.
func drawSearchHighlights(v *View) []string {
	v.draw(true)

	width, height := v.InnerSize()
	rows := make([]string, 0, height)
	for y := range height {
		row := ""
		for x := range width {
			_, style, _ := Screen.Get(v.x0+x+1, v.y0+y+1)
			switch style.GetBackground() {
			case getTcellColor(ColorYellow, OutputNormal):
				row += "#"
			case getTcellColor(ColorCyan, OutputNormal):
				row += "*"
			default:
				row += "."
			}
		}
		rows = append(rows, row)
	}
	return rows
}

// The search matches on screen are highlighted, and the current one in a color of its
// own.
func TestDrawSearchMatches(t *testing.T) {
	WithSimulationScreen(t, 14, 6)

	v := NewView("name", 0, 0, 11, 5, OutputNormal) // InnerWidth 10, InnerHeight 4
	writeLines(v,
		"ab",
		"abab xab",
		"no match",
		"xxxxxxxxxab xxxxx ab",
		"ab ab",
		"ab",
	)
	v.Search("ab", nil)
	_ = v.gotoNextMatch()
	_ = v.gotoNextMatch()
	index, total := v.GetSearchStatus()
	assert.Equal(t, 2, index, "the current match is the second one in the second line")
	assert.Equal(t, 9, total)

	// Scrolled down past the first line, so that the last line is off screen too.
	v.SetOriginY(1)
	assert.Equal(t,
		[]string{
			"##**..##..", // abab xab
			"..........", // no match
			".........#", // xxxxxxxxxa
			"##.##.....", // ab ab
		},
		drawSearchHighlights(v))

	// Scrolled sideways, matches show where they are partly on screen.
	v.SetOriginX(3)
	assert.Equal(t,
		[]string{
			"*..##.....", // b xab
			"..........", // match
			"......##..", // xxxxxxab x
			"##........", // ab
		},
		drawSearchHighlights(v))
}

// Where matches overlap, a cell is highlighted as part of the first match it is in.
func TestDrawOverlappingSearchMatches(t *testing.T) {
	WithSimulationScreen(t, 14, 6)

	v := NewView("name", 0, 0, 11, 2, OutputNormal) // InnerWidth 10, InnerHeight 1
	writeLines(v, "aaaa")
	v.Search("aa", nil)
	_ = v.gotoNextMatch()
	index, total := v.GetSearchStatus()
	assert.Equal(t, 1, index, "the current match is the one in the middle")
	assert.Equal(t, 3, total)

	assert.Equal(t,
		[]string{"##*#......"},
		drawSearchHighlights(v))
}

// The model can find a match that the view doesn't show, for example in a commit hash
// that the view abbreviates. Such a match is highlighted where the model says it is.
// The other matches are highlighted where the view shows them.
func TestDrawSearchMatchesFromTheModel(t *testing.T) {
	WithSimulationScreen(t, 14, 6)

	v := NewView("name", 0, 0, 11, 4, OutputNormal) // InnerWidth 10, InnerHeight 3
	writeLines(v, "ab ab", "1234567", "ab")
	v.Search("ab", []SearchPosition{
		{XStart: 2, XEnd: 6, Y: 0},
		{XStart: 2, XEnd: 6, Y: 1},
		{XStart: 2, XEnd: 6, Y: 2},
	})
	_ = v.gotoNextMatch()
	_ = v.gotoNextMatch()
	index, total := v.GetSearchStatus()
	assert.Equal(t, 2, index, "the current match is the one the view doesn't show")
	assert.Equal(t, 4, total)

	assert.Equal(t,
		[]string{
			"##.##.....", // ab ab
			"..****....", // 1234567
			"##........", // ab
		},
		drawSearchHighlights(v))
}

func BenchmarkDrawViewWithManySearchMatches(b *testing.B) {
	WithSimulationScreen(b, 102, 42)
	v := NewView("name", 0, 0, 101, 41, OutputNormal)
	// 44000 matches, as when searching a long log for a single letter.
	for i := range 4000 {
		fmt.Fprintf(v, "line %d of a branch log, where we searched for every letter e\n", i)
	}
	v.Search("e", nil)
	v.SetOriginY(2000)

	for b.Loop() {
		v.draw(true)
	}
}
