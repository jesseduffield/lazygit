package gocui

import (
	"errors"
	"strings"
	"testing"

	"github.com/gdamore/tcell/v3"
	"github.com/stretchr/testify/assert"
)

func screenRow(y, x0, x1 int) string {
	var sb strings.Builder
	for x := x0; x <= x1; x++ {
		str, _, _ := Screen.Get(x, y)
		sb.WriteString(str)
	}
	return sb.String()
}

func screenColumn(x, y0, y1 int) string {
	var sb strings.Builder
	for y := y0; y <= y1; y++ {
		str, _, _ := Screen.Get(x, y)
		sb.WriteString(str)
	}
	return sb.String()
}

// setFramedView creates a view with a rounded frame in rows 0 to 4.
func setFramedView(t *testing.T, g *Gui, name string, x0, x1 int) *View {
	t.Helper()
	v, err := g.SetView(name, x0, 0, x1, 4, 0)
	if err != nil && !errors.Is(err, ErrUnknownView) {
		assert.NoError(t, err)
	}
	v.FrameRunes = []rune{'─', '│', '╭', '╮', '╰', '╯'}
	return v
}

func TestEdgeHighlightOnVerticalEdges(t *testing.T) {
	g := newTestGui(t)
	left := setFramedView(t, g, "left", 0, 10)
	right := setFramedView(t, g, "right", 11, 30)
	left.EdgeHighlight = &EdgeHighlight{Edge: EdgeRight, Rune: '┃', StartCorner: '┒', EndCorner: '┚', Color: ColorGreen}
	right.EdgeHighlight = &EdgeHighlight{Edge: EdgeLeft, Rune: '┃', StartCorner: '┎', EndCorner: '┖', Color: ColorGreen}

	assert.NoError(t, g.flush())

	assert.Equal(t, "┒┃┃┃┚", screenColumn(10, 0, 4))
	assert.Equal(t, "┎┃┃┃┖", screenColumn(11, 0, 4))
	assert.Equal(t, "╭─────────┒", screenRow(0, 0, 10))
	assert.Equal(t, "│", screenColumn(0, 2, 2))
	_, style, _ := Screen.Get(10, 2)
	assert.Equal(t, tcell.ColorGreen, style.GetForeground())
	_, style, _ = Screen.Get(0, 2)
	assert.NotEqual(t, tcell.ColorGreen, style.GetForeground())
}

func TestEdgeHighlightKeepsFrameRunesIfNoneGiven(t *testing.T) {
	g := newTestGui(t)
	v := setFramedView(t, g, "view", 0, 10)
	v.EdgeHighlight = &EdgeHighlight{Edge: EdgeBottom, Color: ColorGreen}

	assert.NoError(t, g.flush())

	assert.Equal(t, "╰─────────╯", screenRow(4, 0, 10))
	_, style, _ := Screen.Get(5, 4)
	assert.Equal(t, tcell.ColorGreen, style.GetForeground())
}

func TestEdgeHighlightOnTopEdgeAppliesToTitlePrefixSeparator(t *testing.T) {
	g := newTestGui(t)
	v := setFramedView(t, g, "view", 0, 20)
	v.TitlePrefix = "[0]"
	v.Title = "Patch"
	v.EdgeHighlight = &EdgeHighlight{Edge: EdgeTop, Rune: '━', StartCorner: '┍', EndCorner: '┑', Color: ColorGreen}

	assert.NoError(t, g.flush())

	assert.Equal(t, "┍━[0]━Patch━━━━━━━━━┑", screenRow(0, 0, 20))
	_, style, _ := Screen.Get(5, 0)
	assert.Equal(t, tcell.ColorGreen, style.GetForeground())
}
