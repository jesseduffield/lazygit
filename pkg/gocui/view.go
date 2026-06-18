// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

import (
	"fmt"
	"io"
	"slices"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"

	"github.com/gdamore/tcell/v3"
	"github.com/rivo/uniseg"
)

// Constants for overlapping edges
const (
	TOP    = 1 // view is overlapping at top edge
	BOTTOM = 2 // view is overlapping at bottom edge
	LEFT   = 4 // view is overlapping at left edge
	RIGHT  = 8 // view is overlapping at right edge
)

// A View is a window. It maintains its own internal buffer and cursor
// position.
// viewBuffer holds a view's content as cells, together with the cursor and
// escape-sequence decoder state used to turn incoming bytes into those cells.
// A view normally has a single buffer (the one it displays), but bundling this
// state lets a re-render build a second, off-screen buffer and swap it in
// atomically once the new content is ready, so no reader ever sees a
// half-written buffer.
type viewBuffer struct {
	// the view's content: one []cell per unwrapped line
	lines []lineType

	// write cursor into lines
	wx, wy int

	// decodes ESC sequences as bytes are written
	ei *escapeInterpreter

	// If the last character written was a newline, we don't write it but instead
	// set pendingNewline to true. If more text is written, we write the newline
	// then. This avoids an extra blank line at the end of the view.
	pendingNewline bool
}

type View struct {
	name           string
	x0, y0, x1, y1 int // left top right bottom
	ox, oy         int // view offsets
	cx, cy         int // cursor position
	rx, ry         int // Read() offsets
	outMode        OutputMode

	// buf bundles the view's cell buffer and the cursor / escape-parser state
	// used to write into it (see the viewBuffer type). It is the buffer every
	// reader sees.
	buf *viewBuffer

	// While non-nil, writes go here instead of buf, so an async re-render can
	// build its new content without disturbing what readers (draw, clicks,
	// scrolling, the diff-line readers, …) see. The task swaps it into buf once
	// it has read enough to paint (SwapInOffscreenRender), so the displayed
	// content jumps straight from the previous render to the new one with no
	// half-written frame in between. nil during normal (non-async) writes.
	offscreen *viewBuffer
	// The y position of the first line of a range selection.
	// This is not relative to the view's origin: it is relative to the first line
	// of the view's content, so you can scroll the view and this value will remain
	// the same, unlike the view's cy value.
	// A value of -1 means that there is no range selection.
	// This value can be greater than the selected line index, in the event that
	// a user starts a range select and then moves the cursor up.
	rangeSelectStartY int

	// readBuffer is used for storing unread bytes
	readBuffer []byte

	// tained is true if the viewLines must be updated
	tainted bool

	// firstDirtyLine is the index of the lowest line in `lines` that has been
	// written to or highlighted since viewLines was last refreshed, and whose
	// cached wrapping (lineType.wrappedCells) may therefore be stale. Lines
	// below it are unchanged and can reuse their cached wrapping instead of
	// being re-wrapped, which keeps refreshViewLinesIfNeeded cheap while
	// scrolling appends new lines to a long buffer.
	firstDirtyLine int

	// the last position that the mouse was hovering over; nil if the mouse is outside of
	// this view, or not hovering over a cell
	lastHoverPosition *pos

	// the location of the hyperlink that the mouse is currently hovering over; nil if none
	hoveredHyperlink *SearchPosition

	// internal representation of the view's buffer. We will keep viewLines around
	// from a previous render until we explicitly set them to nil, allowing us to
	// render the same content twice without flicker. Wherever we want to render
	// something without any chance of old content appearing (e.g. when actually
	// rendering new content or if the view is resized) we should set tainted to
	// true and viewLines to nil
	viewLines []viewLine

	// While a re-render is loading new content (see offscreen), the displayed
	// buffer is only partially filled once we've swapped the off-screen render
	// in: the task keeps appending lines after the first paint, up to the count
	// needed for an accurate scrollbar. Sizing the scrollbar from that partial
	// view-line count would make the thumb shrink and snap back as the rest
	// streams in. So while a load is in progress we hold the scrollbar's height
	// at this value — the height the view had when the load began — and let it
	// grow only if the new content turns out taller. Zero means no load is in
	// progress and the scrollbar tracks the content directly.
	scrollbarHeightFloor int

	// writeMutex protects locks the write process
	writeMutex sync.Mutex

	// Visible specifies whether the view is visible.
	Visible bool

	// BgColor and FgColor allow to configure the background and foreground
	// colors of the View.
	BgColor, FgColor Attribute

	// SelBgColor and SelFgColor are used to configure the background and
	// foreground colors of the selected line, when it is highlighted.
	SelBgColor, SelFgColor Attribute

	// InactiveViewSelBgColor is used to configure the background color of the
	// selected line, when it is highlighted but the view doesn't have the
	// focus.
	InactiveViewSelBgColor Attribute

	// If Editable is true, keystrokes will be added to the view's internal
	// buffer at the cursor position.
	Editable bool

	// Editor allows to define the editor that manages the editing mode,
	// including keybindings or cursor behaviour. DefaultEditor is used by
	// default.
	Editor Editor

	// Overwrite enables or disables the overwrite mode of the view.
	Overwrite bool

	// If Highlight is true, Sel{Bg,Fg}Colors will be used
	// for the line under the cursor position.
	Highlight bool
	// If HighlightInactive is true, InavtiveViewSel{Bg,Fg}Colors will be used
	// instead of Sel{Bg,Fg}Colors for highlighting selected lines.
	HighlightInactive bool

	// If Frame is true, a border will be drawn around the view.
	Frame bool

	// FrameColor allow to configure the color of the Frame when it is not highlighted.
	FrameColor Attribute

	// FrameRunes allows to define custom runes for the frame edges.
	// The rune slice can be defined with 3 different lengths.
	// If slice doesn't match these lengths, default runes will be used instead of missing one.
	//
	// 2 runes with only horizontal and vertical edges.
	//  []rune{'─', '│'}
	//  []rune{'═','║'}
	// 6 runes with horizontal, vertical edges and top-left, top-right, bottom-left, bottom-right cornes.
	//  []rune{'─', '│', '┌', '┐', '└', '┘'}
	//  []rune{'═','║','╔','╗','╚','╝'}
	// 11 runes which can be used with `gocui.Gui.SupportOverlaps` property.
	//  []rune{'─', '│', '┌', '┐', '└', '┘', '├', '┤', '┬', '┴', '┼'}
	//  []rune{'═','║','╔','╗','╚','╝','╠','╣','╦','╩','╬'}
	FrameRunes []rune

	// If Wrap is true, the content that is written to this View is
	// automatically wrapped when it is longer than its width. If true the
	// view's x-origin will be ignored.
	Wrap bool

	// If Autoscroll is true, the View will automatically scroll down when the
	// text overflows. If true the view's y-origin will be ignored.
	Autoscroll bool

	// If Frame is true, Title allows to configure a title for the view.
	Title string

	// If non-empty, TitlePrefix is prepended to the title of a view regardless on
	// the the currently selected tab (if any.)
	TitlePrefix string

	Tabs     []string
	TabIndex int

	// TitleColor allow to configure the color of title and subtitle for the view.
	TitleColor Attribute

	// If Frame is true, Subtitle allows to configure a subtitle for the view.
	Subtitle string

	// If Mask is true, the View will display the mask instead of the real
	// content
	Mask string

	// Overlaps describes which edges are overlapping with another view's edges
	Overlaps byte

	// ParentView is the view which catches events bubbled up from the given view if there's no matching handler
	ParentView *View

	searcher *searcher

	// KeybindOnEdit should be set to true when you want to execute keybindings even when the view is editable
	// (this is usually not the case)
	KeybindOnEdit bool

	TextArea *TextArea

	// something like '1 of 20' for a list view
	Footer string

	// if true, the user can scroll all the way past the last item until it appears at the top of the view
	CanScrollPastBottom bool

	// if true, the view will automatically recognize https: URLs in the content written to it and render
	// them as hyperlinks
	AutoRenderHyperLinks bool

	// if true, the view will underline hyperlinks only when the cursor is on
	// them; otherwise, they will always be underlined
	UnderlineHyperLinksOnlyOnHover bool

	// number of spaces per \t character, defaults to 4
	TabWidth int
}

type pos struct {
	x, y int
}

// call this in the event of a view resize, or if you want to render new content
// without the chance of old content still appearing, or if you want to remove
// a line from the existing content
func (v *View) clearViewLines() {
	v.tainted = true
	v.viewLines = nil
	v.clearHover()
}

// ClearViewLines is clearViewLines guarded by writeMutex. It's for callers on
// the UI thread (the layout pass) that touch a view whose content a task
// goroutine may be writing concurrently: viewLines/tainted/hover are all
// buffer state that writeMutex protects.
func (v *View) ClearViewLines() {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()
	v.clearViewLines()
}

type searcher struct {
	searchString       string
	searchPositions    []SearchPosition
	modelSearchResults []SearchPosition
	currentSearchIndex int
	onSelectItem       func(*View, int)
	renderSearchStatus func(*View, int, int)
}

func (v *View) setRenderSearchStatus(renderSearchStatus func(*View, int, int)) {
	v.searcher.renderSearchStatus = renderSearchStatus
}

func (v *View) setOnSelectResult(onSelectItem func(*View, int)) {
	v.searcher.onSelectItem = onSelectItem
}

func (v *View) renderSearchStatus(index int, itemCount int) {
	if v.searcher.renderSearchStatus != nil {
		v.searcher.renderSearchStatus(v, index, itemCount)
	}
}

func (v *View) gotoNextMatch() error {
	if len(v.searcher.searchPositions) == 0 {
		return nil
	}
	if v.Highlight && v.oy+v.cy < v.searcher.searchPositions[v.searcher.currentSearchIndex].Y {
		// If the selection is before the current match, just jump to the current match and return.
		// This can only happen if the user has moved the cursor to before the first match.
		v.SelectSearchResult(v.searcher.currentSearchIndex)
		return nil
	}
	if v.searcher.currentSearchIndex >= len(v.searcher.searchPositions)-1 {
		v.searcher.currentSearchIndex = 0
	} else {
		v.searcher.currentSearchIndex++
	}
	v.SelectSearchResult(v.searcher.currentSearchIndex)
	return nil
}

func (v *View) gotoPreviousMatch() error {
	if len(v.searcher.searchPositions) == 0 {
		return nil
	}
	if v.Highlight && v.oy+v.cy > v.searcher.searchPositions[v.searcher.currentSearchIndex].Y {
		// If the selection is after the current match, just jump to the current match and return.
		// This happens if the user has moved the cursor down from the current match.
		v.SelectSearchResult(v.searcher.currentSearchIndex)
		return nil
	}
	if v.searcher.currentSearchIndex == 0 {
		if len(v.searcher.searchPositions) > 0 {
			v.searcher.currentSearchIndex = len(v.searcher.searchPositions) - 1
		}
	} else {
		v.searcher.currentSearchIndex--
	}
	v.SelectSearchResult(v.searcher.currentSearchIndex)
	return nil
}

func (v *View) SelectSearchResult(index int) {
	itemCount := len(v.searcher.searchPositions)
	if itemCount == 0 {
		return
	}
	if index > itemCount-1 {
		index = itemCount - 1
	}

	y := v.searcher.searchPositions[index].Y

	v.FocusPoint(v.ox, y, true)
	v.renderSearchStatus(index, itemCount)
	if v.searcher.onSelectItem != nil {
		v.searcher.onSelectItem(v, y)
	}
}

// Returns <current match index>, <total matches>
func (v *View) GetSearchStatus() (int, int) {
	return v.searcher.currentSearchIndex, len(v.searcher.searchPositions)
}

// modelSearchResults is optional; pass nil to search the view. If non-nil,
// these positions will be used for highlighting search results. Even in this
// case the view will still be searched on a per-line basis, so that the caller
// doesn't have to make assumptions where in the rendered line the search result
// is. The XStart and XEnd values in the modelSearchResults are only used in
// case the search string is not found in the given line, which can happen if
// the view renders an abbreviated version of some of the model data.
//
// Mind the difference between nil and empty slice: nil means we're not
// searching the model, empty slice means we *are* searching the model but we
// didn't find any matches.
func (v *View) UpdateSearchResults(str string, modelSearchResults []SearchPosition) {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	v.searcher.search(str, modelSearchResults)
	v.updateSearchPositions()

	if len(v.searcher.searchPositions) > 0 {
		// get the first result past the current cursor
		currentIndex := 0
		if v.Highlight {
			// ...but only if we're showing the highlighted line
			adjustedY := v.oy + v.cy
			adjustedX := v.ox + v.cx
			for i, pos := range v.searcher.searchPositions {
				if pos.Y > adjustedY || (pos.Y == adjustedY && pos.XStart > adjustedX) {
					currentIndex = i
					break
				}
			}
		}
		v.searcher.currentSearchIndex = currentIndex
	}
}

func (v *View) Search(str string, modelSearchResults []SearchPosition) {
	v.UpdateSearchResults(str, modelSearchResults)

	if len(v.searcher.searchPositions) > 0 {
		v.SelectSearchResult(v.searcher.currentSearchIndex)
	} else {
		v.renderSearchStatus(0, 0)
	}
}

func (v *View) ClearSearch() {
	v.searcher.clearSearch()
}

func (v *View) IsSearching() bool {
	return v.searcher.searchString != ""
}

func (v *View) nearestSearchPosition() int {
	currentLineIndex := v.cy + v.oy
	lastSearchPos := 0
	for i, pos := range v.searcher.searchPositions {
		if pos.Y == currentLineIndex {
			return i
		}
		if pos.Y > currentLineIndex {
			break
		}
		lastSearchPos = i
	}
	return lastSearchPos
}

func (v *View) SetNearestSearchPosition() {
	if len(v.searcher.searchPositions) > 0 {
		newPos := v.nearestSearchPosition()
		if newPos != v.searcher.currentSearchIndex {
			v.searcher.currentSearchIndex = newPos
			v.renderSearchStatus(newPos, len(v.searcher.searchPositions))
		}
	}
}

func (v *View) FocusPoint(cx int, cy int, scrollIntoView bool) {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	v.refreshViewLinesIfNeeded()
	lineCount := len(v.viewLines)
	if cy < 0 || cy > lineCount {
		return
	}

	if scrollIntoView {
		height := v.InnerHeight()
		v.SetOriginY(calculateNewOrigin(cy, v.oy, lineCount, height))
	}

	v.cx = cx
	v.cy = cy - v.oy
}

func (v *View) SetRangeSelectStart(rangeSelectStartY int) {
	v.rangeSelectStartY = rangeSelectStartY
}

func (v *View) CancelRangeSelect() {
	v.rangeSelectStartY = -1
}

func calculateNewOrigin(selectedLine int, oldOrigin int, lineCount int, viewHeight int) int {
	if viewHeight >= lineCount {
		return 0
	} else if selectedLine < oldOrigin || selectedLine >= oldOrigin+viewHeight {
		// If the selected line is outside the visible area, scroll the view so
		// that the selected line is in the middle.
		newOrigin := selectedLine - viewHeight/2

		// However, take care not to overflow if the total line count is less
		// than the view height.
		maxOrigin := lineCount - viewHeight
		if newOrigin > maxOrigin {
			newOrigin = maxOrigin
		}
		if newOrigin < 0 {
			newOrigin = 0
		}

		return newOrigin
	}

	return oldOrigin
}

func (s *searcher) search(str string, modelSearchResults []SearchPosition) {
	s.searchString = str
	s.searchPositions = []SearchPosition{}
	s.modelSearchResults = modelSearchResults
	s.currentSearchIndex = 0
}

func (s *searcher) clearSearch() {
	s.searchString = ""
	s.searchPositions = []SearchPosition{}
	s.currentSearchIndex = 0
}

type SearchPosition struct {
	XStart int
	XEnd   int
	Y      int
}

type viewLine struct {
	linesX, linesY int // coordinates relative to v.buf.lines
	line           []cell

	// Colors used to extend the bg past this wrapped segment's content.
	// Derived at wrap time from the source line — see refreshViewLinesIfNeeded
	// for the per-segment rule.
	trailingFillAttributes *trailingFillAttributes
}

// lineType is one of v.lines: the cells of a source line, plus optional
// trailingFillAttributes recording the colors used to extend the bg
// past the line's content when the writer emitted '\x1b[K'.
type lineType struct {
	cells                  cells
	trailingFillAttributes *trailingFillAttributes

	// wrappedCells caches the result of wrapping `cells` to `wrappedColumns`
	// columns, so that unchanged lines don't have to be re-wrapped on every
	// refreshViewLinesIfNeeded (which runs on every scroll event, via
	// ViewLinesHeight). Wrapping measures every cell's width and allocates, so
	// for a long buffer that dominates the cost of scrolling. The cache is used
	// only for lines below View.firstDirtyLine whose wrappedColumns still
	// matches the current width; nil means nothing is cached yet.
	wrappedCells   [][]cell
	wrappedColumns int
}

// trailingFillAttributes describes the fg/bg colors that draw() should
// use for cells past the end of a wrapped segment's content. On a source
// line this records what the writer asked for via '\x1b[K' (and so opts
// the line in to trailing fill at all); the per-segment values on each
// viewLine are derived from it at wrap time.
type trailingFillAttributes struct {
	fg, bg Attribute
}

type cell struct {
	chr              string // a grapheme cluster
	width            int    // number of terminal cells occupied by chr (always 1 or 2)
	bgColor, fgColor Attribute
	hyperlink        string
	// per-line diff metadata from an OSC 1717 sequence (see
	// diff-line-metadata-notes.md); empty unless a pager emitted it
	metadata string
}

type cells []cell

func characterEquals(chr []byte, b byte) bool {
	return len(chr) == 1 && chr[0] == b
}

func isCRLF(chr []byte) bool {
	return len(chr) == 2 && chr[0] == '\r' && chr[1] == '\n'
}

// String returns a string from a given cell slice.
func (l cells) String() string {
	var str strings.Builder
	for _, c := range l {
		str.WriteString(c.chr)
	}
	return str.String()
}

// NewView returns a new View object.
func NewView(name string, x0, y0, x1, y1 int, mode OutputMode) *View {
	v := &View{
		name:              name,
		x0:                x0,
		y0:                y0,
		x1:                x1,
		y1:                y1,
		Visible:           true,
		Frame:             true,
		Editor:            DefaultEditor,
		tainted:           true,
		outMode:           mode,
		buf:               &viewBuffer{ei: newEscapeInterpreter(mode)},
		searcher:          &searcher{},
		TextArea:          &TextArea{},
		rangeSelectStartY: -1,
		TabWidth:          4,
	}

	v.FgColor, v.BgColor = ColorDefault, ColorDefault
	v.SelFgColor, v.SelBgColor = ColorDefault, ColorDefault
	v.InactiveViewSelBgColor = ColorDefault
	v.TitleColor, v.FrameColor = ColorDefault, ColorDefault
	v.buf.ei.screenColMax = v.InnerWidth()
	return v
}

// SetContentWidth tells the view the screen width that content written to it
// should count soft-wraps against (see escapeInterpreter.notifyCellsWritten).
// Callers pass the view's InnerWidth; it's a separate call, made on the UI
// thread when a render starts, so that the task goroutine that streams the
// content can consult this snapshot instead of reading the view's live
// dimensions (which the UI thread mutates during layout).
func (v *View) SetContentWidth(width int) {
	v.buf.ei.screenColMax = width
}

// Dimensions returns the dimensions of the View
func (v *View) Dimensions() (int, int, int, int) {
	return v.x0, v.y0, v.x1, v.y1
}

// Size returns the number of visible columns and rows in the View, including
// the frame if any
func (v *View) Size() (x, y int) {
	return v.Width(), v.Height()
}

// InnerSize returns the number of usable columns and rows in the View, excluding
// the frame if any
func (v *View) InnerSize() (x, y int) {
	return v.InnerWidth(), v.InnerHeight()
}

func (v *View) Width() int {
	return v.x1 - v.x0 + 1
}

func (v *View) Height() int {
	return v.y1 - v.y0 + 1
}

// The writeable area of the view is always two less then the view's size,
// because if it has a frame, we need to subtract that, but if it doesn't, the
// view is made 1 larger on all sides. I'd like to clean this up at some point,
// but for now we live with this weirdness.
func (v *View) InnerWidth() int {
	innerWidth := v.Width() - 2
	if innerWidth < 0 {
		return 0
	}

	return innerWidth
}

func (v *View) InnerHeight() int {
	innerHeight := v.Height() - 2
	if innerHeight < 0 {
		return 0
	}

	return innerHeight
}

// Name returns the name of the view.
func (v *View) Name() string {
	return v.name
}

// setCharacter sets a character (grapheme cluster) at the given point relative to the view. It applies
// the specified colors, taking into account if the cell must be highlighted. Also, it checks if the
// position is valid.
func (v *View) setCharacter(x, y int, ch string, fgColor, bgColor Attribute) {
	maxX, maxY := v.Size()
	if x < 0 || x >= maxX || y < 0 || y >= maxY {
		return
	}

	if v.Mask != "" {
		fgColor = v.FgColor
		bgColor = v.BgColor
		ch = v.Mask
	} else if v.Highlight {
		rangeSelectStart := v.cy
		rangeSelectEnd := v.cy
		if v.rangeSelectStartY != -1 {
			relativeRangeSelectStart := v.rangeSelectStartY - v.oy
			rangeSelectStart = min(relativeRangeSelectStart, v.cy)
			rangeSelectEnd = max(relativeRangeSelectStart, v.cy)
		}

		if y >= rangeSelectStart && y <= rangeSelectEnd {
			// this ensures we use the bright variant of a colour upon highlight
			fgColorComponent := fgColor & ^AttrAll
			if fgColorComponent >= AttrIsValidColor && fgColorComponent < AttrIsValidColor+8 {
				fgColor += 8
			}
			fgColor = fgColor | AttrBold
			if v.HighlightInactive {
				bgColor = (bgColor & AttrStyleBits) | v.InactiveViewSelBgColor
			} else {
				bgColor = (bgColor & AttrStyleBits) | v.SelBgColor
			}
		}
	}

	if matched, selected := v.isPatternMatchedRune(x, y); matched {
		fgColor = ColorBlack
		if selected {
			bgColor = ColorCyan
		} else {
			bgColor = ColorYellow
		}
	}

	if v.isHoveredHyperlink(x, y) {
		fgColor |= AttrUnderline
	}

	// Don't display empty characters
	if ch == "" {
		ch = " "
	}

	tcellSetCell(v.x0+x+1, v.y0+y+1, ch, fgColor, bgColor, v.outMode)
}

// SetCursor sets the cursor position of the view at the given point,
// relative to the view. It is allowed to set the position to a point outside
// the visible portion of the view, or even outside the content of the view.
// Clients are responsible for clamping to valid positions.
func (v *View) SetCursor(x, y int) {
	v.cx = x
	v.cy = y
}

func (v *View) SetCursorX(x int) {
	v.cx = x
}

func (v *View) SetCursorY(y int) {
	v.cy = y
}

// Cursor returns the cursor position of the view.
func (v *View) Cursor() (x, y int) {
	return v.cx, v.cy
}

func (v *View) CursorX() int {
	return v.cx
}

func (v *View) CursorY() int {
	return v.cy
}

// SetOrigin sets the origin position of the view's internal buffer,
// so the buffer starts to be printed from this point, which means that
// it is linked with the origin point of view. It can be used to
// implement Horizontal and Vertical scrolling with just incrementing
// or decrementing ox and oy.
func (v *View) SetOrigin(x, y int) {
	v.SetOriginX(x)
	v.SetOriginY(y)
}

func (v *View) SetOriginX(x int) {
	if x < 0 {
		x = 0
	}
	v.ox = x
}

func (v *View) SetOriginY(y int) {
	if y < 0 {
		y = 0
	}
	v.oy = y
}

// Origin returns the origin position of the view.
func (v *View) Origin() (x, y int) {
	return v.OriginX(), v.OriginY()
}

func (v *View) OriginX() int {
	return v.ox
}

func (v *View) OriginY() int {
	return v.oy
}

// SetWritePos sets the write position of the view's internal buffer.
// So the next Write call would write directly to the specified position.
func (v *View) SetWritePos(x, y int) {
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}

	v.buf.wx = x
	v.buf.wy = y

	// Changing the write position makes a pending newline obsolete
	v.buf.pendingNewline = false
}

// WritePos returns the current write position of the view's internal buffer.
func (v *View) WritePos() (x, y int) {
	return v.buf.wx, v.buf.wy
}

// SetReadPos sets the read position of the view's internal buffer.
// So the next Read call would read from the specified position.
func (v *View) SetReadPos(x, y int) {
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}

	v.readBuffer = nil
	v.rx = x
	v.ry = y
}

// ReadPos returns the current read position of the view's internal buffer.
func (v *View) ReadPos() (x, y int) {
	return v.rx, v.ry
}

// makeWriteable creates empty cells if required to make position (x, y) writeable.
func (b *viewBuffer) makeWriteable(x, y int) {
	// TODO: make this more efficient

	// line `y` must be index-able (that's why `<=`)
	for len(b.lines) <= y {
		if cap(b.lines) > len(b.lines) {
			newLen := cap(b.lines)
			if newLen > y {
				newLen = y + 1
			}
			b.lines = b.lines[:newLen]
		} else {
			b.lines = append(b.lines, lineType{})
		}
	}
	// cell `x` need not be index-able (that's why `<`)
	// append should be used by `lines[y]` user if he wants to write beyond `x`
	for len(b.lines[y].cells) < x {
		if cap(b.lines[y].cells) > len(b.lines[y].cells) {
			newLen := cap(b.lines[y].cells)
			if newLen > x {
				newLen = x
			}
			b.lines[y].cells = b.lines[y].cells[:newLen]
		} else {
			b.lines[y].cells = append(b.lines[y].cells, cell{})
		}
	}
}

// writeCells copies []cell to (b.wx, b.wy), and advances b.wx accordingly.
// !!! caller MUST ensure that specified location (x, y) is writeable by calling makeWriteable
func (b *viewBuffer) writeCells(cells []cell) {
	var newLen int
	// use maximum len available
	line := b.lines[b.wy].cells[:cap(b.lines[b.wy].cells)]
	maxCopy := len(line) - b.wx
	if maxCopy < len(cells) {
		copy(line[b.wx:], cells[:maxCopy])
		line = append(line, cells[maxCopy:]...)
		newLen = len(line)
	} else { // maxCopy >= len(cells)
		copy(line[b.wx:], cells)
		newLen = b.wx + len(cells)
		if newLen < len(b.lines[b.wy].cells) {
			newLen = len(b.lines[b.wy].cells)
		}
	}
	b.lines[b.wy].cells = line[:newLen]
	b.wx += len(cells)
}

// Write appends a byte slice into the view's internal buffer. Because
// View implements the io.Writer interface, it can be passed as parameter
// of functions like fmt.Fprintf, fmt.Fprintln, io.Copy, etc. Clear must
// be called to clear the view's buffer.
func (v *View) Write(p []byte) (n int, err error) {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	v.write(p)

	return len(p), nil
}

func (v *View) write(p []byte) {
	// An async re-render builds into the off-screen buffer (see View.offscreen)
	// until it swaps in; until then the displayed buffer, and so everything
	// readers see, is left untouched.
	if v.offscreen != nil {
		v.offscreen.write(v, p)
		return
	}

	v.tainted = true
	// write only ever touches lines from v.wy onwards, so any cached wrapping
	// below that stays valid.
	v.firstDirtyLine = min(v.firstDirtyLine, v.buf.wy)
	v.clearHover()

	v.buf.write(v, p)

	v.updateSearchPositions()
}

// write parses p into cells and appends them to the buffer at its write cursor.
// It only touches the buffer; the View wrapper above handles display-side
// effects (tainting, hover, search). v supplies render config (Editable, colors,
// width, tab width, hyperlink auto-rendering).
func (b *viewBuffer) write(v *View, p []byte) {
	// Fill with empty cells, if writing outside current view buffer
	b.makeWriteable(b.wx, b.wy)

	finishLine := func() {
		b.autoRenderHyperlinksInCurrentLine(v)
	}

	advanceToNextLine := func() {
		b.wx = 0
		b.wy++
		if b.wy >= len(b.lines) {
			b.lines = append(b.lines, lineType{})
		}
		// An OSC 1717 diff-metadata sequence applies only to the line it prefixes
		// (the pager re-emits one per line and never closes it), so drop it at the
		// line boundary rather than letting it carry onto a line with no metadata.
		b.ei.metadata.Reset()
	}

	if b.pendingNewline {
		advanceToNextLine()
		b.ei.notifyRowAdvance()
		b.pendingNewline = false
	}

	until := len(p)
	if !v.Editable && until > 0 && p[until-1] == '\n' {
		b.pendingNewline = true
		until--
	}

	state := -1
	var chr []byte
	var width int
	remaining := p[:until]

	for len(remaining) > 0 {
		chr, remaining, width, state = uniseg.FirstGraphemeCluster(remaining, state)

		switch {
		case characterEquals(chr, '\n') || isCRLF(chr):
			finishLine()
			advanceToNextLine()
			b.ei.notifyRowAdvance()
		case characterEquals(chr, '\r'):
			finishLine()
			b.wx = 0
			b.ei.notifyColumnReset()
		default:
			truncateLine, cells := b.parseInput(v, chr, width, b.wx, b.wy)
			if cd, ok := b.ei.instruction.(cursorDown); ok {
				b.ei.instructionRead()
				for range cd.n {
					b.autoRenderHyperlinksInCurrentLine(v)
					advanceToNextLine()
				}
			}
			if cells == nil {
				continue
			}
			b.writeCells(cells)
			if truncateLine {
				b.lines[b.wy].cells = b.lines[b.wy].cells[:b.wx]
			}
			// Soft-wrap tracking. truncateLine is true exactly when the
			// cells are from \x1b[K filling to end of line — ConPTY
			// doesn't advance the cursor for that, so we shouldn't count
			// it toward wraps either.
			if !truncateLine {
				totalWidth := 0
				for _, c := range cells {
					totalWidth += c.width
				}
				b.ei.notifyCellsWritten(totalWidth)
			}
		}
	}

	if b.pendingNewline {
		finishLine()
	} else {
		b.autoRenderHyperlinksInCurrentLine(v)
	}
}

// exported functions use the mutex. Non-exported functions are for internal use
// and a calling function should use a mutex
func (v *View) WriteString(s string) {
	_, _ = v.Write([]byte(s))
}

func (v *View) writeString(s string) {
	v.write([]byte(s))
}

var linkStartChars = []string{"h", "t", "t", "p", "s", ":", "/", "/"}

func findLinkStart(line []cell) int {
	for i := range len(line) - len(linkStartChars) {
		for j := range linkStartChars {
			if line[i+j].chr != linkStartChars[j] {
				break
			}
			if j == len(linkStartChars)-1 {
				return i
			}
		}
	}
	return -1
}

// We need a heuristic to find the end of a hyperlink. Searching for the
// first character that is not a valid URI character is not quite good
// enough, because in markdown it's common to have a hyperlink followed by a
// ')', so we want to stop there. Hopefully URLs containing ')' are uncommon
// enough that this is not a problem.
var lineEndCharacters = map[string]bool{
	"":   true,
	" ":  true,
	"\n": true,
	">":  true,
	"\"": true,
	")":  true,
}

func (b *viewBuffer) autoRenderHyperlinksInCurrentLine(v *View) {
	if !v.AutoRenderHyperLinks {
		return
	}

	line := b.lines[b.wy].cells
	start := 0
	for {
		linkStart := findLinkStart(line[start:])
		if linkStart == -1 {
			break
		}
		linkStart += start
		var link strings.Builder
		linkEnd := linkStart
		for ; linkEnd < len(line); linkEnd++ {
			if _, ok := lineEndCharacters[line[linkEnd].chr]; ok {
				break
			}
			link.WriteString(line[linkEnd].chr)
		}
		for i := linkStart; i < linkEnd; i++ {
			b.lines[b.wy].cells[i].hyperlink = link.String()
		}
		start = linkEnd
	}
}

// parseInput parses char by char the input written to the View. It returns nil
// while processing ESC sequences. Otherwise, it returns a cell slice that
// contains the processed data.
func (b *viewBuffer) parseInput(v *View, ch []byte, width int, x int, _ int) (bool, []cell) {
	cells := []cell{}
	truncateLine := false

	isEscape, err := b.ei.parseOne(ch)
	if err != nil {
		for _, chr := range b.ei.characters() {
			c := cell{
				fgColor: v.FgColor,
				bgColor: v.BgColor,
				chr:     chr,
				width:   uniseg.StringWidth(chr),
			}
			cells = append(cells, c)
		}
		b.ei.reset()
	} else {
		repeatCount := 1
		if _, ok := b.ei.instruction.(eraseInLineFromCursor); ok {
			// Discard any old content past the cursor and record the
			// fill colors so draw() paints the trailing area with them.
			// This extends the bg to the right edge in both the
			// content-fits and content-wraps cases — for the latter,
			// the metadata is what reaches every wrapped segment past
			// the last word.
			b.ei.instructionRead()
			truncateLine = true
			b.lines[b.wy].trailingFillAttributes = &trailingFillAttributes{
				fg: b.ei.curFgColor,
				bg: b.ei.curBgColor,
			}
			return truncateLine, []cell{}
		} else if cf, ok := b.ei.instruction.(cursorForward); ok {
			// emit `n` space cells under the parser-tracked SGR — used
			// to materialize ConPTY's compressed runs of spaces (which
			// it emits as ECH+CUF instead of literal whitespace).
			b.ei.instructionRead()
			repeatCount = cf.n
			ch = []byte{' '}
			width = 1
		} else if isEscape {
			// do not output anything
			return truncateLine, nil
		} else if characterEquals(ch, '\t') {
			// fill tab-sized space
			tabWidth := v.TabWidth
			if tabWidth < 1 {
				tabWidth = 4
			}
			ch = []byte{' '}
			width = 1
			repeatCount = tabWidth - (x % tabWidth)
		}
		c := cell{
			fgColor:   b.ei.curFgColor,
			bgColor:   b.ei.curBgColor,
			hyperlink: b.ei.hyperlink.String(),
			metadata:  b.ei.metadata.String(),
			chr:       string(ch),
			width:     width,
		}
		for range repeatCount {
			cells = append(cells, c)
		}
	}

	return truncateLine, cells
}

// Read reads data into p from the current reading position set by SetReadPos.
// It returns the number of bytes read into p.
// At EOF, err will be io.EOF.
func (v *View) Read(p []byte) (n int, err error) {
	buffer := make([]byte, utf8.UTFMax)
	offset := 0
	if v.readBuffer != nil {
		copy(p, v.readBuffer)
		if len(v.readBuffer) >= len(p) {
			if len(v.readBuffer) > len(p) {
				v.readBuffer = v.readBuffer[len(p):]
			}
			return len(p), nil
		}
		v.readBuffer = nil
	}
	for v.ry < len(v.buf.lines) {
		for v.rx < len(v.buf.lines[v.ry].cells) {
			s := v.buf.lines[v.ry].cells[v.rx].chr
			count := len(s)
			copy(p[offset:], s)
			v.rx++
			newOffset := offset + count
			if newOffset >= len(p) {
				if newOffset > len(p) {
					v.readBuffer = buffer[newOffset-len(p):]
				}
				return len(p), nil
			}
			offset += count
		}
		v.rx = 0
		v.ry++
	}
	return offset, io.EOF
}

// only use this if the calling function has a lock on writeMutex
func (v *View) clear() {
	v.rewind()
	v.buf.lines = nil
	v.clearViewLines()
	// Abandon any in-progress off-screen render: a synchronous SetContent/Clear
	// is taking over the displayed buffer, so writes must go there, not into a
	// stale off-screen buffer left by a stopped task.
	v.offscreen = nil
	// Likewise release any held scrollbar height: the new content is defined
	// synchronously (e.g. a string render superseding a still-loading diff), so
	// there's no async growth left to smooth over and the scrollbar should track
	// the new content directly.
	v.scrollbarHeightFloor = 0
}

// Clear empties the view's internal buffer.
// And resets reading and writing offsets.
func (v *View) Clear() {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	v.clear()
}

func (v *View) SetContent(str string) {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	v.clear()
	v.writeString(str)
}

func (v *View) CopyContent(from *View) {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	// A background task may be streaming output into the source view's buffer
	// via Write, so read it under its own lock. The source is always a
	// different view than the destination (see the sole caller,
	// moveMainContextToTop), and no other code holds two view write locks at
	// once, so this can't deadlock.
	from.writeMutex.Lock()
	defer from.writeMutex.Unlock()

	v.clear()

	// Clone the row slices rather than sharing them: the source view stays
	// live (its streaming task keeps appending rows, and refreshViewLinesIfNeeded
	// fills each row's wrapping cache in place via &lines[i]), so sharing the
	// backing arrays would race those writes against this view's own rendering.
	// This is a shallow clone -- the per-row cell data is immutable once written
	// and stays shared, so the cost is proportional to the number of rows, not
	// their contents.
	v.buf.lines = slices.Clone(from.buf.lines)
	v.viewLines = slices.Clone(from.viewLines)
	v.SetOriginX(from.ox)
	v.SetOriginY(from.oy)
	v.cx = from.cx
	v.cy = from.cy
}

// Rewind sets read and write pos to (0, 0).
func (v *View) Rewind() {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	v.rewind()
}

// similar to Rewind but clears lines. Also similar to Clear but doesn't reset
// viewLines
func (v *View) Reset() {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	v.rewind()
	v.buf.lines = nil
	// As in clear(): abandon any in-progress off-screen render so writes after a
	// reset go to the displayed buffer.
	v.offscreen = nil
}

// This is for when we've done a restart for the sake of avoiding a flicker and
// we've reached the end of the new content to display: we need to clear the remaining
// content from the previous round. We do this by setting v.viewLines to nil so that
// we just render the new content from v.buf.lines directly
func (v *View) FlushStaleCells() {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	v.clearViewLines()
}

// BeginOffscreenRender starts building a re-render into an off-screen buffer.
// Until SwapInOffscreenRender promotes it, writes go to that buffer and the
// displayed buffer — what every reader sees — is left as it was. This is how an
// async re-render avoids exposing a half-written buffer: it accumulates
// off-screen and swaps in once it has read enough to paint.
func (v *View) BeginOffscreenRender() {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	v.offscreen = &viewBuffer{ei: newEscapeInterpreter(v.outMode)}
}

// SwapInOffscreenRender promotes the off-screen buffer (see BeginOffscreenRender)
// to the displayed buffer in one step, so the view jumps straight from the
// previous render to the new one with no half-written frame. Writes after this
// append to the now-displayed buffer directly. It is a no-op if no off-screen
// render is in progress, so it is safe to call more than once (e.g. again at EOF
// after an earlier paint already swapped).
func (v *View) SwapInOffscreenRender() {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	if v.offscreen == nil {
		return
	}
	v.buf = v.offscreen
	v.offscreen = nil
	v.tainted = true
	v.clearHover()
}

// FreezeScrollbarHeight records the view's current content height so the
// scrollbar keeps that size while a re-render loads, instead of shrinking and
// snapping back as the partially-loaded content streams in past the first paint
// (see scrollbarHeightFloor). Call it when a load begins, while the view still
// shows the previous render; UnfreezeScrollbarHeight clears it when the load
// ends.
func (v *View) FreezeScrollbarHeight() {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	v.refreshViewLinesIfNeeded()
	v.scrollbarHeightFloor = len(v.viewLines)
}

// UnfreezeScrollbarHeight clears the height held by FreezeScrollbarHeight, so
// the scrollbar tracks the view's content directly again. Call it when a load
// ends.
func (v *View) UnfreezeScrollbarHeight() {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	v.scrollbarHeightFloor = 0
}

// scrollbarContentHeight is the view-line height the scrollbar is sized from.
// While a re-render is loading it is held at the height the view had when the
// load began (see FreezeScrollbarHeight), so the thumb doesn't shrink and jump
// as partially-loaded content streams in.
func (v *View) scrollbarContentHeight() int {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	v.refreshViewLinesIfNeeded()
	return max(len(v.viewLines), v.scrollbarHeightFloor)
}

func (v *View) rewind() {
	v.buf.ei.reset()
	v.buf.ei.resetScreenCursor()

	v.SetReadPos(0, 0)
	v.SetWritePos(0, 0)
}

func containsUpcaseChar(str string) bool {
	for _, ch := range str {
		if unicode.IsUpper(ch) {
			return true
		}
	}

	return false
}

func stringToGraphemes(s string) []string {
	var graphemes []string
	state := -1
	for s != "" {
		var chr string
		chr, s, _, state = uniseg.FirstGraphemeClusterInString(s, state)
		graphemes = append(graphemes, chr)
	}
	return graphemes
}

func (v *View) updateSearchPositions() {
	if v.searcher.searchString != "" {
		var normalizeRune func(s string) string
		var normalizedSearchStr string
		// if we have any uppercase characters we'll do a case-sensitive search
		if containsUpcaseChar(v.searcher.searchString) {
			normalizeRune = func(s string) string { return s }
			normalizedSearchStr = v.searcher.searchString
		} else {
			normalizeRune = strings.ToLower
			normalizedSearchStr = strings.ToLower(v.searcher.searchString)
		}

		searchStrGraphemes := stringToGraphemes(normalizedSearchStr)

		v.searcher.searchPositions = []SearchPosition{}

		searchPositionsForLine := func(line []cell, y int) []SearchPosition {
			var result []SearchPosition
			searchStringWidth := uniseg.StringWidth(v.searcher.searchString)
			x := 0
			for startIdx, cell := range line {
				found := true
				for i, c := range searchStrGraphemes {
					if len(line)-1 < startIdx+i {
						found = false
						break
					}
					if normalizeRune(line[startIdx+i].chr) != c {
						found = false
						break
					}
				}
				if found {
					result = append(result, SearchPosition{XStart: x, XEnd: x + searchStringWidth, Y: y})
				}
				x += cell.width
			}
			return result
		}

		if v.searcher.modelSearchResults != nil {
			for _, result := range v.searcher.modelSearchResults {
				// This code only works when v.Wrap is false.

				if result.Y >= len(v.buf.lines) {
					break
				}

				// If a view line exists for this line index:
				if v.buf.lines[result.Y].cells != nil {
					// search this view line for the search string
					positions := searchPositionsForLine(v.buf.lines[result.Y].cells, result.Y)
					if len(positions) > 0 {
						// If we found any occurrences, add them
						v.searcher.searchPositions = append(v.searcher.searchPositions, positions...)
					} else {
						// Otherwise, the search string was found in the model
						// but not in the view line; this can happen if the view
						// renders only truncated versions of the model strings.
						// In this case, add one search position with what the
						// model search function returned.
						v.searcher.searchPositions = append(v.searcher.searchPositions, result)
					}
				} else {
					// We don't have a view line for this line index. Add a
					// searchPosition anyway, just for the sake of being able to
					// show the "n of m" search status. The X positions don't
					// matter in this case.
					v.searcher.searchPositions = append(v.searcher.searchPositions, SearchPosition{XStart: -1, XEnd: -1, Y: result.Y})
				}
			}
		} else {
			v.refreshViewLinesIfNeeded()
			for y, line := range v.viewLines {
				v.searcher.searchPositions = append(v.searcher.searchPositions, searchPositionsForLine(line.line, y)...)
			}
		}
	}
}

// IsTainted tells us if the view is tainted
func (v *View) IsTainted() bool {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()
	return v.tainted
}

// draw re-draws the view's contents.
func (v *View) draw() {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	if !v.Visible {
		return
	}

	v.clearRunes()

	maxX, maxY := v.InnerSize()

	if v.Wrap {
		if maxX == 0 {
			return
		}
		v.SetOriginX(0)
	}

	v.refreshViewLinesIfNeeded()

	visibleViewLinesHeight := v.viewLineLengthIgnoringTrailingBlankLines()
	if v.Autoscroll && visibleViewLinesHeight > maxY {
		v.SetOriginY(visibleViewLinesHeight - maxY)
	}

	if len(v.viewLines) == 0 {
		return
	}

	start := v.oy
	if start > len(v.viewLines)-1 {
		start = len(v.viewLines) - 1
	}

	emptyCell := cell{chr: " ", width: 1, fgColor: ColorDefault, bgColor: ColorDefault}

	for y, vline := range v.viewLines[start:] {
		if y >= maxY {
			break
		}

		// Decide the colors used for cells past the end of vline.line:
		// the source line's trailingFillAttributes (set by '\x1b[K') if
		// any, otherwise plain defaults.
		trailingCell := emptyCell
		if attrs := vline.trailingFillAttributes; attrs != nil {
			trailingCell.fgColor = attrs.fg
			trailingCell.bgColor = attrs.bg
		}

		// x tracks the current x position in the view, and cellIdx tracks the
		// index of the cell. If we print a double-sized rune, we increment cellIdx
		// by one but x by two.
		x := -v.ox
		cellIdx := 0

		var c cell
		for x < maxX {
			if x < 0 {
				if cellIdx < len(vline.line) {
					x += uniseg.StringWidth(vline.line[cellIdx].chr)
					cellIdx++
					continue
				}

				// no more characters to write so we're only going to be printing empty cells
				// past this point
				x = 0
			}

			// if we're out of cells to write, we'll just print empty cells.
			if cellIdx > len(vline.line)-1 {
				c = trailingCell
			} else {
				c = vline.line[cellIdx]
			}

			fgColor := c.fgColor
			if fgColor == ColorDefault {
				fgColor = v.FgColor
			}
			bgColor := c.bgColor
			if bgColor == ColorDefault {
				bgColor = v.BgColor
			}
			if c.hyperlink != "" && !v.UnderlineHyperLinksOnlyOnHover {
				fgColor |= AttrUnderline
			}

			v.setCharacter(x, y, c.chr, fgColor, bgColor)

			x += c.width
			cellIdx++
		}
	}
}

func (v *View) refreshViewLinesIfNeeded() {
	if !v.tainted {
		return
	}

	maxX := v.InnerWidth()
	wrap := 0
	if v.Wrap {
		wrap = maxX
	}

	lineIdx := 0
	lines := v.buf.lines
	for i := range lines {
		line := &lines[i]

		// Reuse the previously wrapped result for lines that haven't changed
		// since the last refresh (i.e. below firstDirtyLine) and were wrapped at
		// the current width. Wrapping is expensive and this loop runs on every
		// scroll event, so only the lines that were actually just read (or
		// re-highlighted) should be wrapped afresh.
		if line.wrappedCells == nil || line.wrappedColumns != wrap || i >= v.firstDirtyLine {
			line.wrappedCells = lineWrap(line.cells, wrap)
			line.wrappedColumns = wrap
		}
		ls := line.wrappedCells

		for j := range ls {
			// Per-segment trailing fill. When the source line opted in
			// via '\x1b[K', the LAST wrapped segment uses those colors
			// directly; earlier segments use the colors of their own
			// last cell, so the trailing area matches the bg active
			// where that segment ended rather than bleeding the
			// '\x1b[K' bg back across color changes in the line.
			var attrs *trailingFillAttributes
			if line.trailingFillAttributes != nil {
				if j == len(ls)-1 {
					attrs = line.trailingFillAttributes
				} else if len(ls[j]) > 0 {
					last := ls[j][len(ls[j])-1]
					attrs = &trailingFillAttributes{fg: last.fgColor, bg: last.bgColor}
				}
			}
			vline := viewLine{
				linesX: j, linesY: i, line: ls[j],
				trailingFillAttributes: attrs,
			}

			if lineIdx > len(v.viewLines)-1 {
				v.viewLines = append(v.viewLines, vline)
			} else {
				v.viewLines[lineIdx] = vline
			}
			lineIdx++
		}
	}

	v.firstDirtyLine = len(lines)
	// Truncate any entries left over from a previous, longer render. An async
	// re-render builds its content off-screen and swaps it in whole (see
	// View.offscreen), so the buffer this rebuilds from is always a complete
	// render — there is no half-loaded shorter buffer whose tail we'd need to
	// keep showing to avoid a flicker, and a leftover tail would just be stale
	// lines mapping to the wrong buffer rows.
	v.viewLines = v.viewLines[:lineIdx]
	v.tainted = false
}

// if autoscroll is enabled but we only have a single row of cells shown to the
// user, we don't want to scroll to the final line if it contains no text. So
// this tells us the view lines height when we ignore any trailing blank lines
func (v *View) viewLineLengthIgnoringTrailingBlankLines() int {
	for i := len(v.viewLines) - 1; i >= 0; i-- {
		if len(v.viewLines[i].line) > 0 {
			return i + 1
		}
	}
	return 0
}

func (v *View) isPatternMatchedRune(x, y int) (bool, bool) {
	for i, pos := range v.searcher.searchPositions {
		adjustedY := y + v.oy
		adjustedX := x + v.ox
		if adjustedY == pos.Y && adjustedX >= pos.XStart && adjustedX < pos.XEnd {
			return true, i == v.searcher.currentSearchIndex
		}
	}
	return false, false
}

func (v *View) isHoveredHyperlink(x, y int) bool {
	if v.UnderlineHyperLinksOnlyOnHover && v.hoveredHyperlink != nil {
		adjustedY := y + v.oy
		adjustedX := x + v.ox
		return adjustedY == v.hoveredHyperlink.Y && adjustedX >= v.hoveredHyperlink.XStart && adjustedX < v.hoveredHyperlink.XEnd
	}
	return false
}

// realPosition returns the position in the internal buffer corresponding to the
// point (x, y) of the view.
func (v *View) realPosition(vx, vy int) (x, y int, ok bool) {
	vx = v.ox + vx
	vy = v.oy + vy

	if vx < 0 || vy < 0 {
		return 0, 0, false
	}

	if len(v.viewLines) == 0 {
		return vx, vy, true
	}

	if vy < len(v.viewLines) {
		vline := v.viewLines[vy]
		x = vline.linesX + vx
		y = vline.linesY
	} else {
		vline := v.viewLines[len(v.viewLines)-1]
		x = vx
		y = vline.linesY + vy - len(v.viewLines) + 1
	}

	return x, y, true
}

// clearRunes erases all the cells in the view.
func (v *View) clearRunes() {
	maxX, maxY := v.InnerSize()
	for x := range maxX {
		for y := range maxY {
			tcellSetCell(v.x0+x+1, v.y0+y+1, " ", v.FgColor, v.BgColor, v.outMode)
		}
	}
}

// BufferLines returns the lines in the view's internal
// buffer.
func (v *View) BufferLines() []string {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	lines := make([]string, len(v.buf.lines))
	for i, l := range v.buf.lines {
		lines[i] = l.cells.String()
	}
	return lines
}

// Buffer returns a string with the contents of the view's internal
// buffer.
func (v *View) Buffer() string {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	return linesToString(v.buf.lines)
}

// ViewBufferLines returns the lines in the view's internal
// buffer that is shown to the user.
func (v *View) ViewBufferLines() []string {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	v.refreshViewLinesIfNeeded()

	lines := make([]string, len(v.viewLines))
	for i, l := range v.viewLines {
		lines[i] = cells(l.line).String()
	}
	return lines
}

// LinesHeight is the count of view lines (i.e. lines excluding wrapping)
func (v *View) LinesHeight() int {
	return len(v.buf.lines)
}

// ViewLinesHeight is the count of view lines (i.e. lines including wrapping)
func (v *View) ViewLinesHeight() int {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	v.refreshViewLinesIfNeeded()
	return len(v.viewLines)
}

// ViewBuffer returns a string with the contents of the view's buffer that is
// shown to the user.
func (v *View) ViewBuffer() string {
	strs := make([]string, len(v.viewLines))
	for i := range v.viewLines {
		strs[i] = cells(v.viewLines[i].line).String()
	}

	return strings.Join(strs, "\n")
}

// Line returns a string with the line of the view's internal buffer
// at the position corresponding to the point (x, y).
func (v *View) Line(y int) (string, bool) {
	_, y, ok := v.realPosition(0, y)
	if !ok {
		return "", false
	}

	if y < 0 || y >= len(v.buf.lines) {
		return "", false
	}

	return v.buf.lines[y].cells.String(), true
}

// Word returns a string with the word of the view's internal buffer
// at the position corresponding to the point (x, y).
func (v *View) Word(x, y int) (string, bool) {
	x, y, ok := v.realPosition(x, y)
	if !ok {
		return "", false
	}

	if x < 0 || y < 0 || y >= len(v.buf.lines) || x >= len(v.buf.lines[y].cells) {
		return "", false
	}

	str := v.buf.lines[y].cells.String()

	nl := strings.LastIndexFunc(str[:x], indexFunc)
	if nl == -1 {
		nl = 0
	} else {
		nl = nl + 1
	}
	nr := strings.IndexFunc(str[x:], indexFunc)
	if nr == -1 {
		nr = len(str)
	} else {
		nr = nr + x
	}
	return str[nl:nr], true
}

func (v *View) HyperLinkInLine(y int, urlScheme string) (string, bool) {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	linesY, ok := v.bufferLineForViewLine(y)
	if !ok {
		return "", false
	}

	for _, c := range v.buf.lines[linesY].cells {
		if strings.HasPrefix(c.hyperlink, urlScheme) {
			return c.hyperlink, true
		}
	}

	return "", false
}

// DiffLineMetadataInLine returns the OSC 1717 per-line diff metadata payload
// attached to the given (wrapped) view line, if a pager emitted one. In the
// single-column case every cell of the line carries the same payload, so the
// first non-empty one is the answer. See diff-line-metadata-notes.md.
func (v *View) DiffLineMetadataInLine(y int) (string, bool) {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	linesY, ok := v.bufferLineForViewLine(y)
	if !ok {
		return "", false
	}

	for _, c := range v.buf.lines[linesY].cells {
		if c.metadata != "" {
			return c.metadata, true
		}
	}

	return "", false
}

// DiffLineMetadataPayloads returns, per unwrapped buffer line, the distinct
// OSC-1717 metadata payloads carried by that line's cells, in left-to-right order.
// A single-column rendering tags every cell of a line with the same payload (one
// entry); a side-by-side rendering tags each side differently, so a changed row
// yields one payload per side (and a context row, where both sides match, still
// one). It is the multi-record counterpart of DiffLineContent.Metadata, which keeps
// only the first payload — enough to identify a single-column row, but it drops the
// other side of a side-by-side row. Staging a selection uses this to act on every
// change a row covers. Taken under the write lock in one pass so the payloads stay
// consistent with the buffer even if a concurrent re-render is rebuilding it.
func (v *View) DiffLineMetadataPayloads() [][]string {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	result := make([][]string, len(v.buf.lines))
	for i, line := range v.buf.lines {
		var payloads []string
		for _, c := range line.cells {
			if c.metadata != "" && !slices.Contains(payloads, c.metadata) {
				payloads = append(payloads, c.metadata)
			}
		}
		result[i] = payloads
	}
	return result
}

// DiffLineContent is the raw per-line material the diff-line backends parse to
// recover a rendered row's patch-space identity (see diff-line-metadata-notes.md):
// the decolorized text (for host-side parsing, mechanism #1), the OSC-1717
// metadata payload a pager emitted (#2), and the line's hyperlink (delta's
// lazygit-edit fallback). It is indexed by unwrapped buffer line, so one entry
// covers all the (wrapped) view lines that line maps to.
type DiffLineContent struct {
	Text      string
	Metadata  string
	Hyperlink string
}

// DiffLineContents returns the per-line diff material (see DiffLineContent) for
// every line of the displayed buffer. It is the snapshot the diff-line backends
// scan: taken under the write lock in one call, so the text, metadata and
// hyperlink of a given line stay consistent with each other even if a concurrent
// re-render is rebuilding the buffer.
func (v *View) DiffLineContents() []DiffLineContent {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	return diffLineContents(v.buf)
}

// OffscreenDiffLineContents returns the per-line diff material (see
// DiffLineContent) for the lines read so far into an in-progress off-screen
// re-render (see BeginOffscreenRender), or nil when no off-screen render is
// underway. While a focused main view re-renders, its displayed buffer still
// holds the previous render; this is how the escape restore scans the *incoming*
// content as it loads, to find the row matching a target patch identity and
// decide when it has read far enough to swap in and scroll there.
func (v *View) OffscreenDiffLineContents() []DiffLineContent {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	if v.offscreen == nil {
		return nil
	}
	return diffLineContents(v.offscreen)
}

// OffscreenDiffLineContentsFrom is OffscreenDiffLineContents restricted to the
// rows from index `from` onward (relative to the off-screen buffer's start, so
// result[0] is buffer line `from`). It lets a scan that tracks how far it has read
// process only the lines that have arrived since, instead of re-snapshotting the
// whole off-screen buffer on every line — the difference between an O(n) and an
// O(n²) restore scan on a large diff. Returns nil when no off-screen render is
// underway or `from` is past the lines read so far.
func (v *View) OffscreenDiffLineContentsFrom(from int) []DiffLineContent {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	if v.offscreen == nil || from < 0 || from >= len(v.offscreen.lines) {
		return nil
	}
	return diffLineContentsFrom(v.offscreen, from)
}

// OffscreenLineCount returns the number of unwrapped buffer lines read so far
// into an in-progress off-screen re-render (see BeginOffscreenRender), or 0 if
// none is underway. The escape restore uses it to tell, cheaply, once it has
// found its target line, when a screenful below it has loaded too — so the swap
// shows the target with context rather than at the very bottom of a part-filled
// view.
func (v *View) OffscreenLineCount() int {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	if v.offscreen == nil {
		return 0
	}
	return len(v.offscreen.lines)
}

func diffLineContents(buf *viewBuffer) []DiffLineContent {
	return diffLineContentsFrom(buf, 0)
}

func diffLineContentsFrom(buf *viewBuffer, from int) []DiffLineContent {
	lines := buf.lines[from:]
	contents := make([]DiffLineContent, len(lines))
	for i, line := range lines {
		text := strings.ReplaceAll(line.cells.String(), "\x00", "")
		var metadata, hyperlink string
		for _, c := range line.cells {
			if metadata == "" {
				metadata = c.metadata
			}
			if hyperlink == "" {
				hyperlink = c.hyperlink
			}
		}
		contents[i] = DiffLineContent{Text: text, Metadata: metadata, Hyperlink: hyperlink}
	}
	return contents
}

// BufferLineForViewLine maps a view line index (which counts wrapped lines) to
// the index of the corresponding line in the unwrapped internal buffer (as
// returned by BufferLines). Several view lines can map to the same buffer line
// when wrapping is on. Returns false if the view line is out of range.
func (v *View) BufferLineForViewLine(y int) (int, bool) {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	return v.bufferLineForViewLine(y)
}

// bufferLineForViewLine maps a (wrapped) view line index to the index of the
// corresponding line in the unwrapped internal buffer (v.buf.lines). It is the
// shared core of the public readers that look up information about the buffer
// line under a given view line (its buffer index, its hyperlink, its diff
// metadata); they all need the same view-line→buffer-line mapping to stay
// consistent with the buffer they then read. The caller must hold writeMutex,
// so that the mapping and the subsequent read of v.buf.lines see the same buffer
// even if a concurrent re-render is rebuilding it.
func (v *View) bufferLineForViewLine(y int) (int, bool) {
	v.refreshViewLinesIfNeeded()

	if y < 0 || y >= len(v.viewLines) {
		return 0, false
	}

	// refreshViewLinesIfNeeded overwrites viewLines in place without truncating,
	// so while a shorter re-render is loading, the tail of viewLines can still
	// hold stale entries pointing past the (shrunk) v.buf.lines. Guard against that.
	linesY := v.viewLines[y].linesY
	if linesY >= len(v.buf.lines) {
		return 0, false
	}

	return linesY, true
}

// ViewLineForBufferLine maps an unwrapped buffer line index to the index of the
// first (wrapped) view line that renders it — the inverse of BufferLineForViewLine.
// The escape restore uses it to turn the buffer line it matched against a target
// patch identity into the view line to scroll to and select. Returns false if the
// buffer line isn't currently rendered into any view line.
func (v *View) ViewLineForBufferLine(bufferLineIdx int) (int, bool) {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	v.refreshViewLinesIfNeeded()
	for i, vl := range v.viewLines {
		if vl.linesY == bufferLineIdx {
			return i, true
		}
	}
	return 0, false
}

// indexFunc allows to split lines by words taking into account spaces
// and 0.
func indexFunc(r rune) bool {
	return r == ' ' || r == 0
}

// SetHighlight toggles highlighting of separate lines, for custom lists
// or multiple selection in views.
func (v *View) SetHighlight(y int, on bool) {
	if y < 0 || y >= len(v.buf.lines) {
		return
	}

	cells := make([]cell, 0, len(v.buf.lines[y].cells))
	for _, c := range v.buf.lines[y].cells {
		if on {
			c.bgColor = v.SelBgColor
			c.fgColor = v.SelFgColor
		} else {
			c.bgColor = v.BgColor
			c.fgColor = v.FgColor
		}
		cells = append(cells, c)
	}
	v.tainted = true
	v.firstDirtyLine = min(v.firstDirtyLine, y)
	v.buf.lines[y].cells = cells
	v.clearHover()
}

func lineWrap(line []cell, columns int) [][]cell {
	if columns == 0 {
		return [][]cell{line}
	}

	var n int
	var offset int
	lastWhitespaceIndex := -1
	lines := make([][]cell, 0, 1)
	for i := range line {
		currChr := line[i].chr
		rw := uniseg.StringWidth(currChr)
		n += rw
		// if currChr == 'g' {
		// 	panic(n)
		// }
		if n > columns {
			// This code is convoluted but we've got comprehensive tests so feel free to do whatever you want
			// to the code to simplify it so long as our tests still pass.
			if currChr == " " {
				// if the line ends in a space, we'll omit it. This means there'll be no
				// way to distinguish between a clean break and a mid-word break, but
				// I think it's worth it.
				lines = append(lines, line[offset:i])
				offset = i + 1
				n = 0
			} else if currChr == "-" {
				// if the last character is hyphen and the width of line is equal to the columns
				lines = append(lines, line[offset:i])
				offset = i
				n = rw
			} else if lastWhitespaceIndex != -1 {
				// if there is a space in the line and the line is not breaking at a space/hyphen
				if line[lastWhitespaceIndex].chr == "-" {
					// if break occurs at hyphen, we'll retain the hyphen
					lines = append(lines, line[offset:lastWhitespaceIndex+1])
				} else {
					// if break occurs at space, we'll omit the space
					lines = append(lines, line[offset:lastWhitespaceIndex])
				}
				// Either way, continue *after* the break
				offset = lastWhitespaceIndex + 1
				n = 0
				for _, c := range line[offset : i+1] {
					n += c.width
				}
			} else {
				// in this case we're breaking mid-word
				lines = append(lines, line[offset:i])
				offset = i
				n = rw
			}
			lastWhitespaceIndex = -1
		} else if line[i].chr == " " || line[i].chr == "-" {
			lastWhitespaceIndex = i
		}
	}

	lines = append(lines, line[offset:])
	return lines
}

func linesToString(lines []lineType) string {
	str := make([]string, len(lines))
	for i := range lines {
		str[i] = lines[i].cells.String()
	}

	return strings.Join(str, "\n")
}

// GetClickedTabIndex tells us which tab was clicked
func (v *View) GetClickedTabIndex(x int) int {
	if len(v.Tabs) <= 1 {
		return 0
	}

	charX := len(v.TitlePrefix) + 1
	if v.TitlePrefix != "" {
		charX += 1
	}
	if x <= charX {
		return -1
	}
	for i, tab := range v.Tabs {
		charX += uniseg.StringWidth(tab)
		if x <= charX {
			return i
		}
		charX += uniseg.StringWidth(" - ")
		if x <= charX {
			return -1
		}
	}

	return -1
}

func (v *View) SelectedLineIdx() int {
	_, seletedLineIdx := v.SelectedPoint()
	return seletedLineIdx
}

// MiddleVisibleLineIdx returns the index of the view line at the middle of the
// content currently on screen. When the content is taller than the viewport this is
// the middle row of the viewport; when it's shorter, it's the middle of the content,
// so the result lands within the content rather than in the empty space below it.
func (v *View) MiddleVisibleLineIdx() int {
	top := v.OriginY()
	bottom := min(top+v.InnerHeight(), v.ViewLinesHeight())
	return (top + bottom) / 2
}

// expected to only be used in tests
func (v *View) SelectedLine() string {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	if len(v.buf.lines) == 0 {
		return ""
	}

	return v.lineContentAtIdx(v.SelectedLineIdx())
}

// expected to only be used in tests
func (v *View) SelectedLines() []string {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	if len(v.buf.lines) == 0 {
		return nil
	}

	startIdx, endIdx := v.SelectedLineRange()

	lines := make([]string, 0, endIdx-startIdx+1)
	for i := startIdx; i <= endIdx; i++ {
		lines = append(lines, v.lineContentAtIdx(i))
	}

	return lines
}

func (v *View) lineContentAtIdx(idx int) string {
	return v.buf.lines[idx].cells.String()
}

func (v *View) SelectedPoint() (int, int) {
	cx, cy := v.Cursor()
	ox, oy := v.Origin()
	return cx + ox, cy + oy
}

func (v *View) SelectedLineRange() (int, int) {
	_, cy := v.Cursor()
	_, oy := v.Origin()

	start := cy + oy

	if v.rangeSelectStartY == -1 {
		return start, start
	}

	end := v.rangeSelectStartY

	if start > end {
		return end, start
	}

	return start, end
}

func (v *View) RenderTextArea() {
	v.Clear()
	fmt.Fprint(v, v.TextArea.GetContent())
	cursorX, cursorY := v.TextArea.GetCursorXY()
	prevOriginX, prevOriginY := v.Origin()
	width, height := v.InnerWidth(), v.InnerHeight()

	newViewCursorX, newOriginX := updatedCursorAndOrigin(prevOriginX, width, cursorX)
	newViewCursorY, newOriginY := updatedCursorAndOrigin(prevOriginY, height, cursorY)

	v.SetCursor(newViewCursorX, newViewCursorY)
	v.SetOrigin(newOriginX, newOriginY)
}

func updatedCursorAndOrigin(prevOrigin int, size int, cursor int) (int, int) {
	var newViewCursor int
	newOrigin := prevOrigin
	usableSize := size - 1

	if cursor > prevOrigin+usableSize {
		newOrigin = cursor - usableSize
		newViewCursor = usableSize
	} else if cursor < prevOrigin {
		newOrigin = cursor
		newViewCursor = 0
	} else {
		newViewCursor = cursor - prevOrigin
	}

	return newViewCursor, newOrigin
}

func (v *View) ClearTextArea() {
	v.Clear()

	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	v.TextArea.Clear()
	v.SetOrigin(0, 0)
	v.SetCursor(0, 0)
}

func (v *View) overwriteLines(y int, content string) {
	// break by newline, then for each line, write it, then add that erase command
	v.buf.wx = 0
	v.buf.wy = y
	v.clearViewLines()

	lines := strings.ReplaceAll(content, "\n", "\x1b[K\n")
	// If the last line doesn't end with a linefeed, add the erase command at
	// the end too
	if !strings.HasSuffix(lines, "\n") {
		lines += "\x1b[K"
	}
	v.writeString(lines)
}

// only call this function if you don't care where v.buf.wx and v.buf.wy end up
func (v *View) OverwriteLines(y int, content string) {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	v.overwriteLines(y, content)
}

// only call this function if you don't care where v.buf.wx and v.buf.wy end up
func (v *View) OverwriteLinesAndClearEverythingElse(lineCount int, y int, content string) {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	v.setContentLineCount(lineCount)

	v.overwriteLines(y, content)

	for i := range y {
		v.buf.lines[i] = lineType{}
	}

	for i := v.buf.wy + 1; i < len(v.buf.lines); i += 1 {
		v.buf.lines[i] = lineType{}
	}
}

func (v *View) setContentLineCount(lineCount int) {
	if lineCount > 0 {
		v.buf.makeWriteable(0, lineCount-1)
	}
	v.buf.lines = v.buf.lines[:lineCount]
}

// If the current search result is no longer visible after a scroll up, select the last search
// result that is visible in the view, if any, or the first one that is below the view if none is
// visible.
func (v *View) selectVisibleSearchResultAfterScrollUp() {
	if !v.Highlight && len(v.searcher.searchPositions) != 0 {
		windowBottom := v.oy + v.InnerHeight()
		if v.searcher.searchPositions[v.searcher.currentSearchIndex].Y >= windowBottom {
			newSearchIndex := v.searcher.currentSearchIndex
			for newSearchIndex > 0 &&
				v.searcher.searchPositions[newSearchIndex-1].Y >= v.oy {
				newSearchIndex--
				if v.searcher.searchPositions[newSearchIndex].Y < windowBottom {
					break
				}
			}
			if v.searcher.currentSearchIndex != newSearchIndex {
				v.searcher.currentSearchIndex = newSearchIndex
				v.renderSearchStatus(newSearchIndex, len(v.searcher.searchPositions))
			}
		}
	}
}

// If the current search result is no longer visible after a scroll down, select the first search
// result that is visible in the view, if any, or the last one that is above the view if none is
// visible.
func (v *View) selectVisibleSearchResultAfterScrollDown() {
	if !v.Highlight && len(v.searcher.searchPositions) != 0 {
		if v.searcher.searchPositions[v.searcher.currentSearchIndex].Y < v.oy {
			newSearchIndex := v.searcher.currentSearchIndex
			windowBottom := v.oy + v.InnerHeight()
			for newSearchIndex+1 < len(v.searcher.searchPositions) &&
				v.searcher.searchPositions[newSearchIndex+1].Y < windowBottom {
				newSearchIndex++
				if v.searcher.searchPositions[newSearchIndex].Y >= v.oy {
					break
				}
			}
			if v.searcher.currentSearchIndex != newSearchIndex {
				v.searcher.currentSearchIndex = newSearchIndex
				v.renderSearchStatus(newSearchIndex, len(v.searcher.searchPositions))
			}
		}
	}
}

func (v *View) ScrollUp(amount int) {
	if amount > v.oy {
		amount = v.oy
	}

	if amount != 0 {
		v.SetOriginY(v.oy - amount)
		v.cy += amount

		v.clearHover()
		v.selectVisibleSearchResultAfterScrollUp()
	}
}

// ensures we don't scroll past the end of the view's content
func (v *View) ScrollDown(amount int) {
	adjustedAmount := v.adjustDownwardScrollAmount(amount)
	if adjustedAmount > 0 {
		v.SetOriginY(v.oy + adjustedAmount)
		v.cy -= adjustedAmount

		v.clearHover()
		v.selectVisibleSearchResultAfterScrollDown()
	}
}

func (v *View) ScrollLeft(amount int) {
	newOx := v.ox - amount
	if newOx < 0 {
		newOx = 0
	}
	if newOx != v.ox {
		v.SetOriginX(newOx)

		v.clearHover()
	}
}

// not applying any limits to this
func (v *View) ScrollRight(amount int) {
	v.SetOriginX(v.ox + amount)

	v.clearHover()
}

func (v *View) adjustDownwardScrollAmount(scrollHeight int) int {
	_, oy := v.Origin()
	y := oy
	if !v.CanScrollPastBottom {
		sy := v.InnerHeight()
		y += sy
	}
	scrollableLines := v.ViewLinesHeight() - y
	if scrollableLines < 0 {
		return 0
	}

	margin := v.scrollMargin()
	if scrollableLines-margin < scrollHeight {
		scrollHeight = scrollableLines - margin
	}
	if oy+scrollHeight < 0 {
		return 0
	}

	return scrollHeight
}

// scrollMargin is about how many lines must still appear if you scroll
// all the way down. We'll subtract this from the total amount of scrollable lines
func (v *View) scrollMargin() int {
	if v.CanScrollPastBottom {
		// Setting to 2 because of the newline at the end of the file that we're likely showing.
		// If we want to scroll past bottom outside the context of reading a file's contents,
		// we should make this into a field on the view to be configured by the client.
		// For now we're hardcoding it.
		return 2
	}

	return 0
}

// Returns true if the view contains a line containing the given text with the given
// foreground color
func (v *View) ContainsColoredText(fgColor string, text string) bool {
	for _, line := range v.buf.lines {
		if containsColoredTextInLine(fgColor, text, line.cells) {
			return true
		}
	}

	return false
}

func containsColoredTextInLine(fgColorStr string, text string, line []cell) bool {
	fgColor := tcell.GetColor(fgColorStr)

	currentMatch := ""
	for i := range line {
		cell := line[i]

		// stripping attributes by converting to and from hex
		cellColor := tcell.NewHexColor(cell.fgColor.Hex())

		if cellColor == fgColor {
			currentMatch += cell.chr
		} else if currentMatch != "" {
			if strings.Contains(currentMatch, text) {
				return true
			}
			currentMatch = ""
		}
	}

	return strings.Contains(currentMatch, text)
}

func (v *View) onMouseMove(x int, y int) {
	if v.Editable || !v.UnderlineHyperLinksOnlyOnHover {
		return
	}

	// Reading v.viewLines (here and in findHyperlinkAt) must hold writeMutex like
	// every other reader: this runs on the event-handling goroutine, and a
	// concurrent re-render on the task goroutine can rebuild or shrink viewLines
	// between the bounds check below and the indexing in findHyperlinkAt — which
	// panicked with an out-of-range index.
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	// newCx and newCy are relative to the view port, i.e. to the visible area of the view
	newCx := x - v.x0 - 1
	newCy := y - v.y0 - 1
	// newX and newY are relative to the view's content, independent of its scroll position
	newX := newCx + v.ox
	newY := newCy + v.oy

	if newY >= 0 && newY <= len(v.viewLines)-1 && newX >= 0 && newX <= len(v.viewLines[newY].line)-1 {
		if v.lastHoverPosition == nil || v.lastHoverPosition.x != newX || v.lastHoverPosition.y != newY {
			v.hoveredHyperlink = v.findHyperlinkAt(newX, newY)
		}
		v.lastHoverPosition = &pos{x: newX, y: newY}
	} else {
		v.lastHoverPosition = nil
		v.hoveredHyperlink = nil
	}
}

func (v *View) findHyperlinkAt(x, y int) *SearchPosition {
	linkStr := v.viewLines[y].line[x].hyperlink
	if linkStr == "" {
		return nil
	}

	xStart := x
	for xStart > 0 && v.viewLines[y].line[xStart-1].hyperlink == linkStr {
		xStart--
	}
	xEnd := x + 1
	for xEnd < len(v.viewLines[y].line) && v.viewLines[y].line[xEnd].hyperlink == linkStr {
		xEnd++
	}

	return &SearchPosition{XStart: xStart, XEnd: xEnd, Y: y}
}

func (v *View) clearHover() {
	v.hoveredHyperlink = nil
	v.lastHoverPosition = nil
}
