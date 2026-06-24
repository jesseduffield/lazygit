// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

import (
	"strings"
	"testing"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
	"github.com/rivo/uniseg"
	"github.com/stretchr/testify/assert"
)

// WithSimulationScreen swaps the package-level Screen for a tcell
// terminfo-backed mock terminal so tests can call view.draw() and
// inspect rendered cells via Screen.Get(). The previous Screen is
// restored on test cleanup.
func WithSimulationScreen(t *testing.T, width, height int) {
	t.Helper()
	saved := Screen
	if err := (&Gui{}).tcellInitSimulation(width, height); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		Screen.Fini()
		Screen = saved
	})
}

func TestWriteString(t *testing.T) {
	tests := []struct {
		existingLines  []string
		stringsToWrite []string
		expectedLines  [][]string
	}{
		{
			[]string{},
			[]string{""},
			[][]string{{}},
		},
		{
			[]string{},
			[]string{"1\n"},
			[][]string{{"1"}},
		},
		{
			[]string{},
			[]string{"1\n", "2\n"},
			[][]string{{"1"}, {"2"}},
		},
		{
			[]string{"a"},
			[]string{"1\n"},
			[][]string{{"1"}},
		},
		{
			[]string{"a\x00"},
			[]string{"1\n"},
			[][]string{{"1", "\x00"}},
		},
		{
			[]string{"ab"},
			[]string{"1\n"},
			[][]string{{"1", "b"}},
		},
		{
			[]string{"abc"},
			[]string{"1\n"},
			[][]string{{"1", "b", "c"}},
		},
		{
			[]string{},
			[]string{"1\r"},
			[][]string{{"1"}},
		},
		{
			[]string{"a"},
			[]string{"1\r"},
			[][]string{{"1"}},
		},
		{
			[]string{"a\x00"},
			[]string{"1\r"},
			[][]string{{"1", "\x00"}},
		},
		{
			[]string{"ab"},
			[]string{"1\r"},
			[][]string{{"1", "b"}},
		},
		{
			[]string{"abc"},
			[]string{"1\r"},
			[][]string{{"1", "b", "c"}},
		},
	}

	for _, test := range tests {
		v := NewView("name", 0, 0, 10, 10, OutputNormal)
		for _, l := range test.existingLines {
			v.buf.lines = append(v.buf.lines, lineType{cells: stringToCells(l)})
		}
		for _, s := range test.stringsToWrite {
			v.writeString(s)
		}
		var resultingLines [][]string
		for _, l := range v.buf.lines {
			resultingLines = append(resultingLines, cellsToStrings(l.cells))
		}
		assert.Equal(t, test.expectedLines, resultingLines)
	}
}

func TestUpdatedCursorAndOrigin(t *testing.T) {
	tests := []struct {
		prevOrigin     int
		size           int
		cursor         int
		expectedCursor int
		expectedOrigin int
	}{
		{0, 10, 0, 0, 0},
		{0, 10, 9, 9, 0},
		{0, 10, 10, 9, 1},
		{0, 10, 19, 9, 10},
		{0, 10, 20, 9, 11},
		{20, 10, 19, 0, 19},
		{20, 10, 25, 5, 20},
	}

	for _, test := range tests {
		cursor, origin := updatedCursorAndOrigin(test.prevOrigin, test.size, test.cursor)
		assert.EqualValues(t, test.expectedCursor, cursor, "Cursor is wrong")
		assert.EqualValues(t, test.expectedOrigin, origin, "Origin in wrong")
	}
}

func TestAutoRenderingHyperlinks(t *testing.T) {
	v := NewView("name", 0, 0, 10, 10, OutputNormal)
	v.AutoRenderHyperLinks = true

	v.writeString("htt")
	// No hyperlinks are generated for incomplete URLs
	assert.Equal(t, "", v.buf.lines[0].cells[0].hyperlink)
	// Writing more characters to the same line makes the link complete (even
	// though we didn't see a newline yet)
	v.writeString("ps://example.com")
	assert.Equal(t, "https://example.com", v.buf.lines[0].cells[0].hyperlink)

	v.Clear()
	// Valid but incomplete URL
	v.writeString("https://exa")
	assert.Equal(t, "https://exa", v.buf.lines[0].cells[0].hyperlink)
	// Writing more characters to the same fixes the link
	v.writeString("mple.com")
	assert.Equal(t, "https://example.com", v.buf.lines[0].cells[0].hyperlink)
}

func TestSelectedLineBgColorWidth(t *testing.T) {
	tests := []struct {
		name         string
		bgColorWidth int
		selected     func(screenX int) bool
	}{
		{
			name:         "zero uses full-width selection background",
			bgColorWidth: 0,
			selected:     func(_ int) bool { return true },
		},
		{
			name:         "non-zero uses left-edge selection background",
			bgColorWidth: 2,
			selected:     func(screenX int) bool { return screenX <= 2 },
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			WithSimulationScreen(t, 14, 4)

			v := NewView("name", 0, 0, 11, 3, OutputNormal) // InnerWidth=10
			v.Highlight = true
			v.SelBgColor = ColorBlue
			v.SelectedLineBgColorWidth = test.bgColorWidth
			v.writeString("0123456789\n")
			v.draw()

			selectedBg := getTcellColor(ColorBlue, OutputNormal)
			for screenX := 1; screenX <= 10; screenX++ {
				_, style, _ := Screen.Get(screenX, 1)
				if test.selected(screenX) {
					assert.Equal(t, selectedBg, style.GetBackground(), "cell %d should show the selection background", screenX)
				} else {
					assert.Equal(t, tcell.ColorDefault, style.GetBackground(), "cell %d should keep its original background", screenX)
				}
			}
		})
	}
}

func TestDiffLineMetadata(t *testing.T) {
	v := NewView("name", 0, 0, 80, 10, OutputNormal)

	// Synthetic delta-style output: each content line is prefixed with an
	// OSC 1717 sequence carrying version;type;new;old;file (old empty unless a
	// deletion), and the OSC bytes themselves must not become visible cells. The
	// final line is a header with no OSC, to prove the metadata doesn't bleed.
	osc := func(payload string) string { return "\x1b]1717;" + payload + "\x1b\\" }
	v.writeString(strings.Join([]string{
		osc("1;c;1;;foo.txt") + "line1",
		osc("1;d;2;2;foo.txt") + "old2",
		osc("1;a;2;;foo.txt") + "new2",
		"@@ a header line with no metadata @@",
	}, "\n"))

	type result struct {
		payload string
		ok      bool
	}
	got := make([]result, len(v.buf.lines))
	for y := range v.buf.lines {
		payload, ok := v.DiffLineMetadataInLine(y)
		got[y] = result{payload, ok}
	}

	assert.Equal(t, []result{
		{"1;c;1;;foo.txt", true},
		{"1;d;2;2;foo.txt", true},
		{"1;a;2;;foo.txt", true},
		{"", false}, // the header line carries no metadata (no bleed)
	}, got)

	// The OSC sequence is consumed as an escape, so the visible text is intact.
	assert.Equal(t, "line1", v.BufferLines()[0])
	assert.Equal(t, "@@ a header line with no metadata @@", v.BufferLines()[3])
}

func TestDiffLineMetadataPayloads(t *testing.T) {
	v := NewView("name", 0, 0, 80, 10, OutputNormal)

	osc := func(payload string) string { return "\x1b]1717;" + payload + "\x1b\\" }
	v.writeString(strings.Join([]string{
		// A single-column row: one payload tags the whole line.
		osc("1;c;1;;foo.txt") + "context",
		// A side-by-side change row: the deletion tags the left half and the
		// addition replacing it tags the right half of the same rendered line.
		osc("1;d;2;2;foo.txt") + "old2  " + osc("1;a;2;;foo.txt") + "new2",
		// A header line with no metadata.
		"@@ header @@",
	}, "\n"))

	assert.Equal(t, [][]string{
		{"1;c;1;;foo.txt"},
		{"1;d;2;2;foo.txt", "1;a;2;;foo.txt"},
		nil,
	}, v.DiffLineMetadataPayloads())
}

func TestDiffLineMetadataHandshakeSwallowed(t *testing.T) {
	v := NewView("name", 0, 0, 80, 10, OutputNormal)

	// A metadata-aware pager emits a version-only handshake (an OSC 1717 with no
	// fields) as its first output, immediately before the diff, to announce it speaks
	// the protocol. It must be swallowed whole: no visible bytes, no phantom line, and
	// crucially it must not attach as metadata to the diff header that follows it.
	osc := func(payload string) string { return "\x1b]1717;" + payload + "\x1b\\" }
	v.writeString(osc("1") + strings.Join([]string{
		"diff --git a/foo.txt b/foo.txt",
		osc("1;a;1;;foo.txt") + "added",
	}, "\n"))

	// The handshake produced no phantom line and no visible bytes.
	assert.Equal(t, []string{
		"diff --git a/foo.txt b/foo.txt",
		"added",
	}, v.BufferLines())

	// The handshake didn't bleed onto the header line, and the real per-line metadata
	// after it still applies.
	type result struct {
		payload string
		ok      bool
	}
	got := make([]result, len(v.buf.lines))
	for y := range v.buf.lines {
		payload, ok := v.DiffLineMetadataInLine(y)
		got[y] = result{payload, ok}
	}
	assert.Equal(t, []result{
		{"", false},
		{"1;a;1;;foo.txt", true},
	}, got)
}

// When a re-render produces fewer view lines than the previous one,
// refreshViewLinesIfNeeded must truncate viewLines to the new content. If it
// didn't (it used to overwrite in place and keep the tail), a reader could map a
// view line that no longer exists onto the wrong buffer line — and with wrapping
// the stale entry's buffer index can still be in range of the new, shorter,
// less-wrapped buffer, so an in-range guard alone wouldn't catch it. See
// diff-line-metadata-notes.md §8.
func TestBufferLineForViewLineStaleTail(t *testing.T) {
	v := NewView("name", 0, 0, 10, 10, OutputNormal) // InnerWidth is 9
	v.Wrap = true

	// First render: two lines that each wrap into three view lines, so the
	// buffer has 2 lines but there are 6 view lines.
	v.writeString(strings.Repeat("a", 27) + "\n" + strings.Repeat("b", 27))
	assert.Equal(t, 6, v.ViewLinesHeight())

	// Re-render with shorter content (rewind, then overwrite from the top with
	// three short, unwrapped lines). There are now only 3 view lines.
	v.Reset()
	v.writeString("aaa\nbbb\nccc")
	assert.Equal(t, 3, v.ViewLinesHeight())

	// A real view line maps to its buffer line as usual.
	bufferLine, ok := v.BufferLineForViewLine(1)
	assert.True(t, ok)
	assert.Equal(t, 1, bufferLine)

	// View line 4 no longer exists in the current buffer, so the mapping must
	// fail rather than land on a stale entry from the previous render.
	_, ok = v.BufferLineForViewLine(4)
	assert.False(t, ok)
}

// An async re-render builds into an off-screen buffer and swaps it in once it
// has enough to paint, so readers keep seeing the previous render — coherent and
// consistent — until the new content appears in one step. See View.offscreen.
func TestOffscreenRender(t *testing.T) {
	v := NewView("name", 0, 0, 80, 10, OutputNormal)

	v.writeString("a\nb\nc")
	assert.Equal(t, []string{"a", "b", "c"}, v.ViewBufferLines())

	// Render new, longer content off-screen.
	v.BeginOffscreenRender()
	v.writeString("w\nx\ny\nz")

	// The displayed buffer is untouched: readers still see the previous render,
	// and the view-line→buffer-line mapping stays consistent with it.
	assert.Equal(t, []string{"a", "b", "c"}, v.ViewBufferLines())
	bufferLine, ok := v.BufferLineForViewLine(1)
	assert.True(t, ok)
	assert.Equal(t, 1, bufferLine)

	// Swapping in reveals the new content in one step.
	v.SwapInOffscreenRender()
	assert.Equal(t, []string{"w", "x", "y", "z"}, v.ViewBufferLines())

	// A further write now appends to the displayed buffer directly.
	v.writeString("\nmore")
	assert.Equal(t, []string{"w", "x", "y", "z", "more"}, v.ViewBufferLines())
}

// The escape restore scans the *incoming* content of a re-render as it loads,
// before it is swapped in, so it can find the row matching a target identity and
// decide when to swap. That means reading the off-screen buffer's loaded rows
// (text, metadata, hyperlink) while the displayed buffer still shows the old
// render. See View.offscreen / OffscreenDiffLineContents.
func TestOffscreenDiffLineContents(t *testing.T) {
	v := NewView("name", 0, 0, 80, 10, OutputNormal)

	// No off-screen render in progress: nothing to scan.
	assert.Nil(t, v.OffscreenDiffLineContents())

	osc := func(payload string) string { return "\x1b]1717;" + payload + "\x1b\\" }
	v.BeginOffscreenRender()
	v.writeString(strings.Join([]string{
		osc("1;c;1;;foo.txt") + "context",
		osc("1;a;2;;foo.txt") + "added",
	}, "\n"))

	contents := v.OffscreenDiffLineContents()
	assert.Equal(t, []DiffLineContent{
		{Text: "context", Metadata: "1;c;1;;foo.txt"},
		{Text: "added", Metadata: "1;a;2;;foo.txt"},
	}, contents)

	// The displayed buffer is still empty; the scan reads the off-screen render.
	assert.Empty(t, v.BufferLines())
}

// The escape restore matches a target identity against a buffer line, then needs
// the view line that renders it to scroll there and select it. ViewLineForBufferLine
// is that inverse mapping, and it must point at the *first* of the (wrapped) view
// lines a buffer line spans.
func TestViewLineForBufferLine(t *testing.T) {
	v := NewView("name", 0, 0, 10, 10, OutputNormal) // InnerWidth is 9
	v.Wrap = true

	// Buffer line 0 is short (one view line); buffer line 1 wraps into three view
	// lines (view lines 1, 2, 3); buffer line 2 is short again (view line 4).
	v.writeString("short\n" + strings.Repeat("b", 27) + "\nlast")

	for bufferLine, wantViewLine := range map[int]int{0: 0, 1: 1, 2: 4} {
		viewLine, ok := v.ViewLineForBufferLine(bufferLine)
		assert.True(t, ok)
		assert.Equal(t, wantViewLine, viewLine)
	}

	_, ok := v.ViewLineForBufferLine(3)
	assert.False(t, ok)
}

// While an async re-render loads, it swaps in only a partially-filled buffer at
// its first paint and keeps appending lines afterwards. The scrollbar must keep
// using the pre-load height until the load ends, so the thumb doesn't shrink and
// snap back as the rest streams in. See View.scrollbarHeightFloor.
func TestScrollbarHeightHeldWhileLoading(t *testing.T) {
	v := NewView("name", 0, 0, 80, 12, OutputNormal)

	// Initial render: 100 lines, scrolled well down.
	v.writeString(strings.Repeat("x\n", 100))
	v.SetOrigin(0, 80)
	assert.Equal(t, 100, v.scrollbarContentHeight())

	// A re-render begins while the previous render is still shown: hold the
	// scrollbar height at the current value.
	v.FreezeScrollbarHeight()

	// The off-screen render swaps in only a screenful at its first paint.
	v.BeginOffscreenRender()
	v.writeString(strings.Repeat("y\n", 30))
	v.SwapInOffscreenRender()

	// The displayed buffer is now short, but the scrollbar height stays held, so
	// the thumb keeps its position instead of jumping.
	assert.Equal(t, 30, v.ViewLinesHeight())
	assert.Equal(t, 100, v.scrollbarContentHeight())

	// The rest of the content streams in.
	v.writeString(strings.Repeat("y\n", 70))
	assert.Equal(t, 100, v.scrollbarContentHeight())

	// Once the load ends, the scrollbar tracks the real content directly again.
	v.UnfreezeScrollbarHeight()
	assert.Equal(t, 100, v.scrollbarContentHeight())
}

// If a synchronous render (e.g. a string render) supersedes a still-loading diff
// before it reaches its end, the held scrollbar height must be released, so the
// scrollbar reflects the new content rather than the abandoned load's height.
func TestScrollbarHeightReleasedWhenContentReplaced(t *testing.T) {
	v := NewView("name", 0, 0, 80, 12, OutputNormal)

	v.writeString(strings.Repeat("x\n", 100))
	v.FreezeScrollbarHeight()
	assert.Equal(t, 100, v.scrollbarContentHeight())

	// A synchronous render replaces the content before the (notional) load ends.
	v.SetContent("just a few\nshort lines\nhere")
	assert.Equal(t, 3, v.scrollbarContentHeight())
}

func TestContainsColoredText(t *testing.T) {
	hexColor := func(text string, hexStr string) []cell {
		cells := make([]cell, len(text))
		hex := GetColor(hexStr)
		for i, chr := range text {
			cells[i] = cell{fgColor: hex, chr: string(chr)}
		}
		return cells
	}
	red := "#ff0000"
	green := "#00ff00"
	redStr := func(text string) []cell { return hexColor(text, red) }
	greenStr := func(text string) []cell { return hexColor(text, green) }

	concat := func(lines ...[]cell) []cell {
		var cells []cell
		for _, line := range lines {
			cells = append(cells, line...)
		}
		return cells
	}

	tests := []struct {
		lines      [][]cell
		fgColorStr string
		text       string
		expected   bool
	}{
		{
			lines:      [][]cell{concat(redStr("a"))},
			fgColorStr: red,
			text:       "a",
			expected:   true,
		},
		{
			lines:      [][]cell{concat(redStr("a"))},
			fgColorStr: red,
			text:       "b",
			expected:   false,
		},
		{
			lines:      [][]cell{concat(redStr("a"))},
			fgColorStr: green,
			text:       "b",
			expected:   false,
		},
		{
			lines:      [][]cell{concat(redStr("hel"), greenStr("lo"), redStr(" World!"))},
			fgColorStr: red,
			text:       "hello",
			expected:   false,
		},
		{
			lines:      [][]cell{concat(redStr("hel"), greenStr("lo"), redStr(" World!"))},
			fgColorStr: green,
			text:       "lo",
			expected:   true,
		},
		{
			lines: [][]cell{
				redStr("hel"),
				redStr("lo"),
			},
			fgColorStr: red,
			text:       "hello",
			expected:   false,
		},
	}

	for i, test := range tests {
		lines := make([]lineType, len(test.lines))
		for j, cells := range test.lines {
			lines[j] = lineType{cells: cells}
		}
		v := &View{buf: &viewBuffer{lines: lines}}
		assert.Equal(t, test.expected, v.ContainsColoredText(test.fgColorStr, test.text), "Test %d failed", i)
	}
}

func TestWriteCursorPositionEscape(t *testing.T) {
	// ConPTY presents its child's output as a screen buffer and uses cursor
	// positioning escapes (CUP, `\x1b[<row>;<col>H`) to skip over blank rows
	// rather than emitting empty LFs for them. The escape interpreter must
	// synthesize the row advances those CUPs imply; otherwise non-blank rows
	// that the child separated with blank lines end up adjacent in the view.
	v := NewView("name", 0, 0, 20, 10, OutputNormal)
	// "a", then "skip to row 3" (i.e. one blank row), then "b".
	v.writeString("a\r\n\x1b[3;1Hb\r\n")

	got := make([][]string, 0, len(v.buf.lines))
	for _, l := range v.buf.lines {
		got = append(got, cellsToStrings(l.cells))
	}

	assert.Equal(t, [][]string{{"a"}, {}, {"b"}}, got)
}

func TestWriteCursorPositionEscapeAcrossWrites(t *testing.T) {
	// Mirrors the production flow: bufio.Scanner splits the pty output on
	// LF and feeds the view one Write per line (each with a trailing \n
	// appended). The parser's screen-row counter must keep ticking across
	// writes; otherwise CUPs are evaluated against a stale row and
	// overshoot, producing too many blank lines instead of the right
	// number.
	v := NewView("name", 0, 0, 30, 30, OutputNormal)
	v.writeString("a\n")
	v.writeString("b\n")
	// ConPTY is on row 3 here; CUP to row 5 should skip exactly one row.
	v.writeString("c\x1b[5;1Hd\n")

	got := make([][]string, 0, len(v.buf.lines))
	for _, l := range v.buf.lines {
		got = append(got, cellsToStrings(l.cells))
	}
	assert.Equal(t, [][]string{
		{"a"},
		{"b"},
		{"c"},
		{},
		{"d"},
	}, got)
}

func TestWriteCursorForwardEscape(t *testing.T) {
	// ConPTY compresses runs of default-colored spaces into ECH (\x1b[NX,
	// "clear N cells, cursor stationary") + CUF (\x1b[NC, "cursor forward
	// N") rather than emitting them literally. The interpreter has to
	// materialize CUF as N visible spaces; otherwise the gap collapses
	// and content that followed the indentation slides left.
	v := NewView("name", 0, 0, 20, 10, OutputNormal)
	// "a" + ECH 5 + CUF 5 + "b" — visually "a     b".
	v.writeString("a\x1b[5X\x1b[5Cb\n")

	got := make([][]string, 0, len(v.buf.lines))
	for _, l := range v.buf.lines {
		got = append(got, cellsToStrings(l.cells))
	}

	assert.Equal(t, [][]string{{"a", " ", " ", " ", " ", " ", "b"}}, got)
}

func TestWriteCursorPositionEscapeWithSoftWraps(t *testing.T) {
	// If a logical line is longer than ConPTY's terminal width, ConPTY
	// soft-wraps it onto multiple physical rows in its screen, and any
	// subsequent CUP is addressed against the post-wrap row count. To
	// keep our screen-row counter accurate we have to count those wraps
	// as we write the cells. InnerWidth here is 5; "abcdefghij" (10
	// cells) wraps onto 2 rows, so ConPTY is on row 3 after the LF and a
	// CUP to row 4 should skip exactly one row.
	v := NewView("name", 0, 0, 6, 30, OutputNormal) // Width=7, InnerWidth=5
	v.writeString("abcdefghij\n")
	v.writeString("\x1b[4;1Hxyz\n")

	got := make([][]string, 0, len(v.buf.lines))
	for _, l := range v.buf.lines {
		got = append(got, cellsToStrings(l.cells))
	}
	assert.Equal(t, [][]string{
		{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"},
		{},
		{"x", "y", "z"},
	}, got)
}

func stringToCells(s string) []cell {
	var cells []cell
	state := -1
	for len(s) > 0 {
		var c string
		var w int
		c, s, w, state = uniseg.FirstGraphemeClusterInString(s, state)
		cells = append(cells, cell{chr: c, width: w})
	}
	return cells
}

func cellsToString(cells []cell) string {
	var s strings.Builder
	for _, c := range cells {
		s.WriteString(c.chr)
	}
	return s.String()
}

func cellsToStrings(cells []cell) []string {
	s := []string{}
	for _, c := range cells {
		s = append(s, c.chr)
	}
	return s
}

func TestLineWrap(t *testing.T) {
	testCases := []struct {
		name     string
		line     string
		columns  int
		expected []string
	}{
		{
			name:    "Wrap on space",
			line:    "Hello World",
			columns: 5,
			expected: []string{
				"Hello",
				"World",
			},
		},
		{
			name:    "Wrap on hyphen",
			line:    "Hello-World",
			columns: 6,
			expected: []string{
				"Hello-",
				"World",
			},
		},
		{
			name:    "Wrap on hyphen 2",
			line:    "Blah Hello-World",
			columns: 12,
			expected: []string{
				"Blah Hello-",
				"World",
			},
		},
		{
			name:    "Wrap on hyphen 3",
			line:    "Blah Hello-World",
			columns: 11,
			expected: []string{
				"Blah Hello-",
				"World",
			},
		},
		{
			name:    "Wrap on hyphen 4",
			line:    "Blah Hello-World",
			columns: 10,
			expected: []string{
				"Blah Hello",
				"-World",
			},
		},
		{
			name:    "Wrap on space 2",
			line:    "Blah Hello World",
			columns: 10,
			expected: []string{
				"Blah Hello",
				"World",
			},
		},
		{
			name:    "Wrap on space with more words",
			line:    "Longer word here",
			columns: 10,
			expected: []string{
				"Longer",
				"word here",
			},
		},
		{
			name:    "Split word that's too long",
			line:    "ThisWordIsWayTooLong",
			columns: 10,
			expected: []string{
				"ThisWordIs",
				"WayTooLong",
			},
		},
		{
			name:    "Split word that's too long over multiple lines",
			line:    "ThisWordIsWayTooLong",
			columns: 5,
			expected: []string{
				"ThisW",
				"ordIs",
				"WayTo",
				"oLong",
			},
		},
		{
			name:    "Lots of hyphens",
			line:    "one-two-three-four-five",
			columns: 8,
			expected: []string{
				"one-two-",
				"three-",
				"four-",
				"five",
			},
		},
		{
			name:    "Several lines using all the available width",
			line:    "aaa bb cc ddd-ee ff",
			columns: 5,
			expected: []string{
				"aaa",
				"bb cc",
				"ddd-",
				"ee ff",
			},
		},
		{
			name:    "Multi-cell runes",
			line:    "🐤🐤🐤 🐝🐝 🙉 🦊🦊🦊-🐬🐬 🦢🦢",
			columns: 9,
			expected: []string{
				"🐤🐤🐤",
				"🐝🐝 🙉",
				"🦊🦊🦊-",
				"🐬🐬 🦢🦢",
			},
		},
		{
			name:    "Space in last column",
			line:    "hello world",
			columns: 6,
			expected: []string{
				"hello",
				"world",
			},
		},
		{
			name:    "Hyphen in last column",
			line:    "hello-world",
			columns: 6,
			expected: []string{
				"hello-",
				"world",
			},
		},
		{
			name:    "English text",
			line:    "+The sea reach of the Thames stretched before us like the bedinnind of an interminable waterway. In the offind the sea and the sky were welded todether without a joint, and in the luminous space the tanned sails of the bardes drifting blah blah",
			columns: 81,
			expected: []string{
				"+The sea reach of the Thames stretched before us like the bedinnind of an",
				"interminable waterway. In the offind the sea and the sky were welded todether",
				"without a joint, and in the luminous space the tanned sails of the bardes",
				"drifting blah blah",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lineCells := stringToCells(tc.line)

			result := lineWrap(lineCells, tc.columns)

			resultStrings := make([]string, len(result))
			for i, line := range result {
				resultStrings[i] = cellsToString(line)
			}

			assert.EqualValues(t, tc.expected, resultStrings)
		})
	}
}

// TestNewlineTerminatedLineClearsTrailingBg verifies that a '\n' resets
// any attributes (e.g. AttrReverse-driven background) past the line's
// content, so a reversed cell at the end doesn't bleed into the empty
// area to the right.
func TestNewlineTerminatedLineClearsTrailingBg(t *testing.T) {
	WithSimulationScreen(t, 14, 5)

	v := NewView("name", 0, 0, 11, 4, OutputNormal)

	// \x1b[7m sets reverse; \x1b[31m sets fg=red. With reverse the cell
	// renders with bg=red. The trailing area past "foo" must NOT extend
	// the red bg because '\n' marks the line as cleanly terminated.
	v.writeString("\x1b[7m\x1b[31mfoo\x1b[0m\n")
	v.draw()

	// First row: cells 1..3 are "foo" (render with red bg via reverse),
	// cells 4..10 are trailing and should be plain default.
	for x := 4; x <= 10; x++ {
		_, style, _ := Screen.Get(x, 1)
		assert.Equal(t, tcell.ColorDefault, style.GetForeground(),
			"trailing cell at (%d, 1) should have default fg", x)
		assert.False(t, style.HasReverse(),
			"trailing cell at (%d, 1) should not have reverse attribute", x)
	}
}

// TestUnterminatedReverseLineDoesNotExtend verifies that an unterminated
// line ending with an AttrReverse cell does NOT propagate the reversed
// background past the line's content — matching real terminal behavior
// (try `print '\x1b[7m\x1b[31mfoo'` in a shell). The trailing area
// is rendered as plain default.
func TestUnterminatedReverseLineDoesNotExtend(t *testing.T) {
	WithSimulationScreen(t, 14, 5)

	v := NewView("name", 0, 0, 11, 4, OutputNormal)

	// Reverse + red fg, "foo", no termination. The trailing cells past
	// "foo" should be plain default, NOT a continuation of the red bg.
	v.writeString("\x1b[7m\x1b[31mfoo")
	v.draw()

	// Cells 4..10 are trailing and should be default with no reverse.
	for x := 4; x <= 10; x++ {
		_, style, _ := Screen.Get(x, 1)
		assert.Equal(t, tcell.ColorDefault, style.GetForeground(),
			"trailing cell at (%d, 1) should have default fg", x)
		assert.False(t, style.HasReverse(),
			"trailing cell at (%d, 1) should not have reverse attribute", x)
	}
}

// TestShortFilledLineExtendsBgWithoutWrap verifies that '\x1b[K' fills
// the rest of the line with the current bg color for a line that's
// short enough to fit within the view's inner width.
func TestShortFilledLineExtendsBgWithoutWrap(t *testing.T) {
	WithSimulationScreen(t, 14, 5)

	v := NewView("name", 0, 0, 11, 4, OutputNormal)

	// \x1b[41m sets bg=red. "hi" fits within InnerWidth=10; \x1b[K should
	// fill the remaining 8 cells with red.
	v.writeString("\x1b[41mhi\x1b[K\x1b[0m\n")
	v.draw()

	// All ten cells at (1..10, 1) should have red bg.
	for x := 1; x <= 10; x++ {
		_, style, _ := Screen.Get(x, 1)
		assert.Equal(t, color.Maroon, style.GetBackground(),
			"cell at (%d, 1) should have red bg", x)
	}
}

// TestWrappedFilledLineExtendsBgToEdge verifies that when a line is
// filled to the edge with \x1b[K (the pattern used by `delta` for diff
// lines) but exceeds the view's inner width, every wrapped segment
// extends the fill background past its content to the right edge.
func TestWrappedFilledLineExtendsBgToEdge(t *testing.T) {
	WithSimulationScreen(t, 14, 6)

	// View dimensions: Width=12 (x0=0..x1=11), Height=6; InnerWidth=10,
	// InnerHeight=4. Frame inset of 1 places content cells at screen
	// (1..10, 1..4).
	v := NewView("name", 0, 0, 11, 5, OutputNormal)
	v.Wrap = true

	// Content with spaces so word wrap ends each segment before the
	// right edge: "aaa bbb ccc ddd eee" wraps at InnerWidth=10 to three
	// segments — "aaa bbb" / "ccc ddd" / "eee". Each row's trailing area
	// must pick up the red fill from \x1b[K.
	v.writeString("\x1b[41m" + "aaa bbb ccc ddd eee" + "\x1b[0m\x1b[41m\x1b[K\x1b[0m\n")
	v.draw()

	// All three wrapped rows should have the red fill background across
	// the full InnerWidth, including the trailing cells past each row's
	// last word.
	for y := 1; y <= 3; y++ {
		for x := 1; x <= 10; x++ {
			_, style, _ := Screen.Get(x, y)
			assert.Equal(t, color.Maroon, style.GetBackground(),
				"cell at (%d, %d) should have red bg", x, y)
		}
	}
}

// TestMulticolorWrappedFillUsesLastCellOfEachSegment demonstrates that
// when a wrapped line switches bg color part-way through and ends with
// \x1b[K, the trailing area on each wrapped row should match the bg
// that was active where that row's content ended — not the \x1b[K bg,
// which would bleed the color from the end of the logical line back
// into the earlier wrapped rows.
func TestMulticolorWrappedFillUsesLastCellOfEachSegment(t *testing.T) {
	WithSimulationScreen(t, 14, 6)

	// View dimensions: Width=12 (x0=0..x1=11), Height=6; InnerWidth=10,
	// InnerHeight=4. Frame inset of 1 places content cells at screen
	// (1..10, 1..4).
	v := NewView("name", 0, 0, 11, 5, OutputNormal)
	v.Wrap = true

	// Content "aaa bbb ccc" is 11 cells; lineWrap breaks at the space
	// between "bbb" and "ccc" (index 7) so segment 1 is "aaa bbb" (red,
	// last cell red) and segment 2 is "ccc" (green, last cell green).
	// \x1b[K records the green bg on the source line.
	v.writeString("\x1b[41maaa bbb\x1b[42m ccc\x1b[K\x1b[0m\n")
	v.draw()

	// Row 1's content ends with a red cell at x=7, so trailing columns
	// 8..10 should pick up red rather than the \x1b[K's green.
	for x := 8; x <= 10; x++ {
		_, style, _ := Screen.Get(x, 1)
		assert.Equal(t, color.Maroon, style.GetBackground(),
			"trailing cell at (%d, 1) should have red bg (matching segment's last cell)", x)
	}

	// Row 2's content ends with a green cell at x=3, so trailing
	// columns 4..10 should pick up green (matching both the segment's
	// last cell and the \x1b[K bg — these happen to agree here).
	for x := 4; x <= 10; x++ {
		_, style, _ := Screen.Get(x, 2)
		assert.Equal(t, color.Green, style.GetBackground(),
			"trailing cell at (%d, 2) should have green bg", x)
	}
}

// TestInclusionGutter verifies the on-demand inclusion gutter reserves a
// left-hand column, draws the marker glyph on marked lines only, and shifts the
// content right past it.
func TestInclusionGutter(t *testing.T) {
	WithSimulationScreen(t, 14, 6)

	// InnerWidth=10; the frame inset of 1 places view x=0 at screen x=1.
	v := NewView("name", 0, 0, 11, 5, OutputNormal)
	v.Wrap = true
	v.InclusionGutterMarker = "✓"

	v.writeString("aaa\nbbb\nccc\n")

	// The gutter is 2 columns wide (marker + separator); mark the middle line.
	v.SetInclusionGutter(true, []bool{false, true, false})
	v.draw()

	// The marker appears at the gutter's first column (view x=0 → screen x=1) on
	// the marked line only.
	chr, _, _ := Screen.Get(1, 1)
	assert.Equal(t, " ", chr, "unmarked line has no gutter marker")
	chr, _, _ = Screen.Get(1, 2)
	assert.Equal(t, "✓", chr, "marked line shows the gutter marker")
	chr, _, _ = Screen.Get(1, 3)
	assert.Equal(t, " ", chr, "unmarked line has no gutter marker")

	// The content is shifted right past the 2-column gutter (view x=2 → screen x=3).
	chr, _, _ = Screen.Get(3, 1)
	assert.Equal(t, "a", chr, "content is shifted past the gutter")
	chr, _, _ = Screen.Get(3, 2)
	assert.Equal(t, "b", chr)
	chr, _, _ = Screen.Get(3, 3)
	assert.Equal(t, "c", chr)

	// Hiding the gutter again returns the content flush left (view x=0 → screen x=1).
	v.SetInclusionGutter(false, nil)
	v.draw()
	chr, _, _ = Screen.Get(1, 1)
	assert.Equal(t, "a", chr, "content is flush left with no gutter")
}

// TestInclusionGutterMarkerOnEverySegment verifies that a marked buffer line that
// wraps shows the marker on every wrapped segment, and that the gutter narrows the
// content wrap width.
func TestInclusionGutterMarkerOnEverySegment(t *testing.T) {
	WithSimulationScreen(t, 14, 6)

	v := NewView("name", 0, 0, 11, 5, OutputNormal) // InnerWidth=10
	v.Wrap = true
	v.InclusionGutterMarker = "✓"

	// 10 cells; with a 2-column gutter the content wrap width is 8, so this wraps
	// to "01234567" / "89".
	v.writeString("0123456789\n")
	v.SetInclusionGutter(true, []bool{true})
	v.draw()

	// First segment: marker present, content starts at screen x=3 and the eighth
	// content cell ("7") sits at the right edge (screen x=10).
	chr, _, _ := Screen.Get(1, 1)
	assert.Equal(t, "✓", chr)
	chr, _, _ = Screen.Get(3, 1)
	assert.Equal(t, "0", chr)
	chr, _, _ = Screen.Get(10, 1)
	assert.Equal(t, "7", chr)

	// Continuation segment: marker too, content resumes at screen x=3.
	chr, _, _ = Screen.Get(1, 2)
	assert.Equal(t, "✓", chr, "continuation segment also shows the marker")
	chr, _, _ = Screen.Get(3, 2)
	assert.Equal(t, "8", chr)
}
