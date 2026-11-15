package gocui

// FrameEdge names one of the edges of a view's frame.
type FrameEdge int

const (
	EdgeLeft FrameEdge = iota
	EdgeRight
	EdgeTop
	EdgeBottom
)

// EdgeHighlight draws one edge of a view's frame in a different color than the
// rest of the frame, and optionally with different runes.
type EdgeHighlight struct {
	Edge FrameEdge
	// The rune for the edge; 0 keeps the frame's rune.
	Rune rune
	// The runes for the corners at the two ends of the edge, the top or left
	// one first; 0 keeps the frame's rune.
	StartCorner, EndCorner rune
	Color                  Attribute
}

func (g *Gui) drawEdgeHighlight(v *View, bgColor Attribute) {
	h := v.EdgeHighlight

	runeH, runeV := '─', '│'
	if len(v.FrameRunes) >= 2 {
		runeH, runeV = v.FrameRunes[0], v.FrameRunes[1]
	}
	runeTL, runeTR, runeBL, runeBR := '┌', '┐', '└', '┘'
	if len(v.FrameRunes) >= 6 {
		runeTL, runeTR, runeBL, runeBR = v.FrameRunes[2], v.FrameRunes[3], v.FrameRunes[4], v.FrameRunes[5]
	}
	orDefault := func(r rune, defaultRune rune) rune {
		if r == 0 {
			return defaultRune
		}
		return r
	}
	set := func(x, y int, ch rune) {
		g.SetRune(x, y, ch, h.Color, bgColor)
	}

	switch h.Edge {
	case EdgeLeft, EdgeRight:
		x, startCorner, endCorner := v.x0, runeTL, runeBL
		if h.Edge == EdgeRight {
			x, startCorner, endCorner = v.x1, runeTR, runeBR
		}
		edgeRune := orDefault(h.Rune, runeV)
		showScrollbar, scrollbarStart, scrollbarEnd := calcRealScrollbarStartEnd(v)
		for y := v.y0 + 1; y < v.y1; y++ {
			ch := edgeRune
			if h.Edge == EdgeRight {
				ch = calcScrollbarRune(showScrollbar, scrollbarStart, scrollbarEnd, y, edgeRune)
			}
			set(x, y, ch)
		}
		set(x, v.y0, orDefault(h.StartCorner, startCorner))
		set(x, v.y1, orDefault(h.EndCorner, endCorner))
	case EdgeTop, EdgeBottom:
		y, startCorner, endCorner := v.y0, runeTL, runeTR
		if h.Edge == EdgeBottom {
			y, startCorner, endCorner = v.y1, runeBL, runeBR
		}
		edgeRune := orDefault(h.Rune, runeH)
		for x := v.x0 + 1; x < v.x1; x++ {
			set(x, y, edgeRune)
		}
		set(v.x0, y, orDefault(h.StartCorner, startCorner))
		set(v.x1, y, orDefault(h.EndCorner, endCorner))
	}
}

// titlePrefixSeparator returns the rune between a view's title prefix and its
// title, and the color to draw it in.
func (v *View) titlePrefixSeparator(fgColor Attribute) (rune, Attribute) {
	separator := '─'
	if len(v.FrameRunes) > 0 {
		separator = v.FrameRunes[0]
	}
	if h := v.EdgeHighlight; h != nil && h.Edge == EdgeTop {
		if h.Rune != 0 {
			separator = h.Rune
		}
		return separator, h.Color
	}
	return separator, fgColor
}
