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
			v.lines = append(v.lines, lineType{cells: stringToCells(l)})
		}
		for _, s := range test.stringsToWrite {
			v.writeString(s)
		}
		var resultingLines [][]string
		for _, l := range v.lines {
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
	assert.Equal(t, "", v.lines[0].cells[0].hyperlink)
	// Writing more characters to the same line makes the link complete (even
	// though we didn't see a newline yet)
	v.writeString("ps://example.com")
	assert.Equal(t, "https://example.com", v.lines[0].cells[0].hyperlink)

	v.Clear()
	// Valid but incomplete URL
	v.writeString("https://exa")
	assert.Equal(t, "https://exa", v.lines[0].cells[0].hyperlink)
	// Writing more characters to the same fixes the link
	v.writeString("mple.com")
	assert.Equal(t, "https://example.com", v.lines[0].cells[0].hyperlink)
}

func TestDiffLineMetadata(t *testing.T) {
	v := NewView("name", 0, 0, 80, 10, OutputNormal)

	// Synthetic delta-style output: each content line is prefixed with an
	// OSC 456 sequence carrying version;type;new;old;file (old empty unless a
	// deletion), and the OSC bytes themselves must not become visible cells. The
	// final line is a header with no OSC, to prove the metadata doesn't bleed.
	osc := func(payload string) string { return "\x1b]456;" + payload + "\x1b\\" }
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
	got := make([]result, len(v.lines))
	for y := range v.lines {
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

// When a re-render produces fewer view lines than the previous one,
// refreshViewLinesIfNeeded overwrites viewLines in place without truncating, so
// the tail keeps the previous render's entries (deliberately, so the view keeps
// showing old content until the new content catches up). A reader must not map
// a view line in that stale tail to a buffer line. With wrapping, the stale
// entry's buffer index can still be in range of the new (shorter, less-wrapped)
// buffer, so the in-range guard alone lets it through and maps a view line that
// no longer exists onto the wrong buffer line. See diff-line-metadata-notes.md
// §8.
func TestBufferLineForViewLineStaleTail(t *testing.T) {
	v := NewView("name", 0, 0, 10, 10, OutputNormal) // InnerWidth is 9
	v.Wrap = true

	// First render: two lines that each wrap into three view lines, so the
	// buffer has 2 lines but there are 6 view lines.
	v.writeString(strings.Repeat("a", 27) + "\n" + strings.Repeat("b", 27))
	assert.Equal(t, 6, v.ViewLinesHeight())

	// Re-render with shorter content the flicker-avoidance way: rewind (which
	// keeps the old view lines) and overwrite from the top with three short,
	// unwrapped lines. There are now only 3 real view lines, but the previous
	// render's view lines 3..5 linger in the tail.
	v.Reset()
	v.writeString("aaa\nbbb\nccc")

	// A real view line maps to its buffer line as usual.
	bufferLine, ok := v.BufferLineForViewLine(1)
	assert.True(t, ok)
	assert.Equal(t, 1, bufferLine)

	// View line 4 is in the stale tail: it no longer exists in the current
	// buffer, so the mapping must fail. (On the buggy code it instead maps to
	// buffer line 1, the stale entry's lingering index.)
	_, ok = v.BufferLineForViewLine(4)
	assert.False(t, ok)
}

// While holding the view lines, the view keeps drawing the previous render even
// as the buffer is overwritten, so a re-render that is restoring a scroll
// position can keep showing the coherent placeholder until its first paint
// reveals the loaded content in one step. See View.holdViewLines.
func TestHoldViewLines(t *testing.T) {
	v := NewView("name", 0, 0, 80, 10, OutputNormal)

	v.writeString("a\nb\nc")
	assert.Equal(t, []string{"a", "b", "c"}, v.ViewBufferLines())

	v.SetHoldViewLines(true)

	// Re-render the flicker-avoidance way (rewind, then overwrite from the top).
	v.Reset()
	v.writeString("w\nx\ny\nz")

	// The new content is in the buffer, but while held the view keeps showing the
	// previous render, and the view-line→buffer-line mapping reports no result.
	assert.Equal(t, []string{"a", "b", "c"}, v.ViewBufferLines())
	_, ok := v.BufferLineForViewLine(0)
	assert.False(t, ok)

	// Releasing the hold reveals the loaded content.
	v.SetHoldViewLines(false)
	assert.Equal(t, []string{"w", "x", "y", "z"}, v.ViewBufferLines())
	bufferLine, ok := v.BufferLineForViewLine(0)
	assert.True(t, ok)
	assert.Equal(t, 0, bufferLine)
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
		v := &View{lines: lines}
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

	got := make([][]string, 0, len(v.lines))
	for _, l := range v.lines {
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

	got := make([][]string, 0, len(v.lines))
	for _, l := range v.lines {
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

	got := make([][]string, 0, len(v.lines))
	for _, l := range v.lines {
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

	got := make([][]string, 0, len(v.lines))
	for _, l := range v.lines {
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
