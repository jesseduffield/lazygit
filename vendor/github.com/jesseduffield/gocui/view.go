// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"

	"github.com/gdamore/tcell/v2"
	"github.com/go-errors/errors"
	"github.com/mattn/go-runewidth"
)

// Constants for overlapping edges
const (
	TOP    = 1 // view is overlapping at top edge
	BOTTOM = 2 // view is overlapping at bottom edge
	LEFT   = 4 // view is overlapping at left edge
	RIGHT  = 8 // view is overlapping at right edge
)

// ErrInvalidPoint is returned when client passed invalid coordinates of a cell.
// Most likely client has passed negative coordinates of a cell.
var ErrInvalidPoint = errors.New("invalid point")

// A View is a window. It maintains its own internal buffer and cursor
// position.
type View struct {
	name           string
	x0, y0, x1, y1 int      // left top right bottom
	ox, oy         int      // view offsets
	cx, cy         int      // cursor position
	rx, ry         int      // Read() offsets
	wx, wy         int      // Write() offsets
	lines          [][]cell // All the data
	outMode        OutputMode
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

	// writeMutex protects locks the write process
	writeMutex sync.Mutex

	// ei is used to decode ESC sequences on Write
	ei *escapeInterpreter

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
	Mask rune

	// Overlaps describes which edges are overlapping with another view's edges
	Overlaps byte

	// If HasLoader is true, the message will be appended with a spinning loader animation
	HasLoader bool

	// IgnoreCarriageReturns tells us whether to ignore '\r' characters
	IgnoreCarriageReturns bool

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

	// if true, the view will underline hyperlinks only when the cursor is on
	// them; otherwise, they will always be underlined
	UnderlineHyperLinksOnlyOnHover bool
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

type searcher struct {
	searchString       string
	searchPositions    []SearchPosition
	modelSearchResults []SearchPosition
	currentSearchIndex int
	onSelectItem       func(int, int, int) error
}

func (v *View) SetOnSelectItem(onSelectItem func(int, int, int) error) {
	v.searcher.onSelectItem = onSelectItem
}

func (v *View) gotoNextMatch() error {
	if len(v.searcher.searchPositions) == 0 {
		return nil
	}
	if v.searcher.currentSearchIndex >= len(v.searcher.searchPositions)-1 {
		v.searcher.currentSearchIndex = 0
	} else {
		v.searcher.currentSearchIndex++
	}
	return v.SelectSearchResult(v.searcher.currentSearchIndex)
}

func (v *View) gotoPreviousMatch() error {
	if len(v.searcher.searchPositions) == 0 {
		return nil
	}
	if v.searcher.currentSearchIndex == 0 {
		if len(v.searcher.searchPositions) > 0 {
			v.searcher.currentSearchIndex = len(v.searcher.searchPositions) - 1
		}
	} else {
		v.searcher.currentSearchIndex--
	}
	return v.SelectSearchResult(v.searcher.currentSearchIndex)
}

func (v *View) SelectSearchResult(index int) error {
	itemCount := len(v.searcher.searchPositions)
	if itemCount == 0 {
		return nil
	}
	if index > itemCount-1 {
		index = itemCount - 1
	}

	y := v.searcher.searchPositions[index].Y

	v.FocusPoint(v.ox, y)
	if v.searcher.onSelectItem != nil {
		return v.searcher.onSelectItem(y, index, itemCount)
	}
	return nil
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
		adjustedY := v.oy + v.cy
		adjustedX := v.ox + v.cx
		for i, pos := range v.searcher.searchPositions {
			if pos.Y > adjustedY || (pos.Y == adjustedY && pos.XStart > adjustedX) {
				currentIndex = i
				break
			}
		}
		v.searcher.currentSearchIndex = currentIndex
	}
}

func (v *View) Search(str string, modelSearchResults []SearchPosition) error {
	v.UpdateSearchResults(str, modelSearchResults)

	if len(v.searcher.searchPositions) > 0 {
		return v.SelectSearchResult(v.searcher.currentSearchIndex)
	}

	return v.searcher.onSelectItem(-1, -1, 0)
}

func (v *View) ClearSearch() {
	v.searcher.clearSearch()
}

func (v *View) IsSearching() bool {
	return v.searcher.searchString != ""
}

func (v *View) FocusPoint(cx int, cy int) {
	lineCount := len(v.lines)
	if cy < 0 || cy > lineCount {
		return
	}
	_, height := v.Size()

	ly := height - 1
	if ly < 0 {
		ly = 0
	}

	v.oy = calculateNewOrigin(cy, v.oy, lineCount, ly)
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
	if viewHeight > lineCount {
		return 0
	} else if selectedLine < oldOrigin || selectedLine > oldOrigin+viewHeight {
		// If the selected line is outside the visible area, scroll the view so
		// that the selected line is in the middle.
		newOrigin := selectedLine - viewHeight/2

		// However, take care not to overflow if the total line count is less
		// than the view height.
		maxOrigin := lineCount - viewHeight - 1
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
	linesX, linesY int // coordinates relative to v.lines
	line           []cell
}

type cell struct {
	chr              rune
	bgColor, fgColor Attribute
	hyperlink        string
}

type lineType []cell

// String returns a string from a given cell slice.
func (l lineType) String() string {
	str := ""
	for _, c := range l {
		str += string(c.chr)
	}
	return str
}

// newView returns a new View object.
func newView(name string, x0, y0, x1, y1 int, mode OutputMode) *View {
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
		ei:                newEscapeInterpreter(mode),
		searcher:          &searcher{},
		TextArea:          &TextArea{},
		rangeSelectStartY: -1,
	}

	v.FgColor, v.BgColor = ColorDefault, ColorDefault
	v.SelFgColor, v.SelBgColor = ColorDefault, ColorDefault
	v.InactiveViewSelBgColor = ColorDefault
	v.TitleColor, v.FrameColor = ColorDefault, ColorDefault
	return v
}

// Dimensions returns the dimensions of the View
func (v *View) Dimensions() (int, int, int, int) {
	return v.x0, v.y0, v.x1, v.y1
}

// Size returns the number of visible columns and rows in the View.
func (v *View) Size() (x, y int) {
	return v.Width(), v.Height()
}

func (v *View) Width() int {
	return v.x1 - v.x0 - 1
}

func (v *View) Height() int {
	return v.y1 - v.y0 - 1
}

// if a view has a frame, that leaves less space for its writeable area
func (v *View) InnerWidth() int {
	innerWidth := v.Width() - v.frameOffset()
	if innerWidth < 0 {
		return 0
	}

	return innerWidth
}

func (v *View) InnerHeight() int {
	innerHeight := v.Height() - v.frameOffset()
	if innerHeight < 0 {
		return 0
	}

	return innerHeight
}

func (v *View) frameOffset() int {
	if v.Frame {
		return 1
	} else {
		return 0
	}
}

// Name returns the name of the view.
func (v *View) Name() string {
	return v.name
}

// setRune sets a rune at the given point relative to the view. It applies the
// specified colors, taking into account if the cell must be highlighted. Also,
// it checks if the position is valid.
func (v *View) setRune(x, y int, ch rune, fgColor, bgColor Attribute) error {
	maxX, maxY := v.Size()
	if x < 0 || x >= maxX || y < 0 || y >= maxY {
		return ErrInvalidPoint
	}

	if v.Mask != 0 {
		fgColor = v.FgColor
		bgColor = v.BgColor
		ch = v.Mask
	} else if v.Highlight {
		var (
			ry, rcy int
			err     error
		)

		_, ry, err = v.realPosition(x, y)
		if err != nil {
			return err
		}
		_, rrcy, err := v.realPosition(v.cx, v.cy)
		// if error is not nil, then the cursor is out of bounds, which is fine
		if err == nil {
			rcy = rrcy
		}

		rangeSelectStart := rcy
		rangeSelectEnd := rcy
		if v.rangeSelectStartY != -1 {
			_, realRangeSelectStart, err := v.realPosition(0, v.rangeSelectStartY-v.oy)
			if err != nil {
				return err
			}

			rangeSelectStart = min(realRangeSelectStart, rcy)
			rangeSelectEnd = max(realRangeSelectStart, rcy)
		}

		if ry >= rangeSelectStart && ry <= rangeSelectEnd {
			// this ensures we use the bright variant of a colour upon highlight
			fgColorComponent := fgColor & ^AttrAll
			if fgColorComponent >= AttrIsValidColor && fgColorComponent < AttrIsValidColor+8 {
				fgColor += 8
			}
			fgColor = fgColor | AttrBold
			if v.HighlightInactive {
				bgColor = bgColor | v.InactiveViewSelBgColor
			} else {
				bgColor = bgColor | v.SelBgColor
			}
		}
	}

	if matched, selected := v.isPatternMatchedRune(x, y); matched {
		if selected {
			bgColor = ColorCyan
		} else {
			bgColor = ColorYellow
		}
	}

	if v.isHoveredHyperlink(x, y) {
		fgColor |= AttrUnderline
	}

	// Don't display NUL characters
	if ch == 0 {
		ch = ' '
	}

	tcellSetCell(v.x0+x+1, v.y0+y+1, ch, fgColor, bgColor, v.outMode)

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// SetCursor sets the cursor position of the view at the given point,
// relative to the view. It checks if the position is valid.
func (v *View) SetCursor(x, y int) error {
	maxX, maxY := v.Size()
	if x < 0 || x >= maxX || y < 0 || y >= maxY {
		return nil
	}
	v.cx = x
	v.cy = y
	return nil
}

func (v *View) SetCursorX(x int) {
	maxX, _ := v.Size()
	if x < 0 || x >= maxX {
		return
	}
	v.cx = x
}

func (v *View) SetCursorY(y int) {
	_, maxY := v.Size()
	if y < 0 || y >= maxY {
		return
	}
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
func (v *View) SetOrigin(x, y int) error {
	if x < 0 || y < 0 {
		return ErrInvalidPoint
	}
	v.ox = x
	v.oy = y
	return nil
}

func (v *View) SetOriginX(x int) error {
	if x < 0 {
		return ErrInvalidPoint
	}
	v.ox = x
	return nil
}

func (v *View) SetOriginY(y int) error {
	if y < 0 {
		return ErrInvalidPoint
	}
	v.oy = y
	return nil
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
func (v *View) SetWritePos(x, y int) error {
	if x < 0 || y < 0 {
		return ErrInvalidPoint
	}
	v.wx = x
	v.wy = y
	return nil
}

// WritePos returns the current write position of the view's internal buffer.
func (v *View) WritePos() (x, y int) {
	return v.wx, v.wy
}

// SetReadPos sets the read position of the view's internal buffer.
// So the next Read call would read from the specified position.
func (v *View) SetReadPos(x, y int) error {
	if x < 0 || y < 0 {
		return ErrInvalidPoint
	}
	v.readBuffer = nil
	v.rx = x
	v.ry = y
	return nil
}

// ReadPos returns the current read position of the view's internal buffer.
func (v *View) ReadPos() (x, y int) {
	return v.rx, v.ry
}

// makeWriteable creates empty cells if required to make position (x, y) writeable.
func (v *View) makeWriteable(x, y int) {
	// TODO: make this more efficient

	// line `y` must be index-able (that's why `<=`)
	for len(v.lines) <= y {
		if cap(v.lines) > len(v.lines) {
			newLen := cap(v.lines)
			if newLen > y {
				newLen = y + 1
			}
			v.lines = v.lines[:newLen]
		} else {
			v.lines = append(v.lines, nil)
		}
	}
	// cell `x` must not be index-able (that's why `<`)
	// append should be used by `lines[y]` user if he wants to write beyond `x`
	for len(v.lines[y]) < x {
		if cap(v.lines[y]) > len(v.lines[y]) {
			newLen := cap(v.lines[y])
			if newLen > x {
				newLen = x
			}
			v.lines[y] = v.lines[y][:newLen]
		} else {
			v.lines[y] = append(v.lines[y], cell{})
		}
	}
}

// writeCells copies []cell to specified location (x, y)
// !!! caller MUST ensure that specified location (x, y) is writeable by calling makeWriteable
func (v *View) writeCells(x, y int, cells []cell) {
	var newLen int
	// use maximum len available
	line := v.lines[y][:cap(v.lines[y])]
	maxCopy := len(line) - x
	if maxCopy < len(cells) {
		copy(line[x:], cells[:maxCopy])
		line = append(line, cells[maxCopy:]...)
		newLen = len(line)
	} else { // maxCopy >= len(cells)
		copy(line[x:], cells)
		newLen = x + len(cells)
		if newLen < len(v.lines[y]) {
			newLen = len(v.lines[y])
		}
	}
	v.lines[y] = line[:newLen]
}

// readCell gets cell at specified location (x, y)
func (v *View) readCell(x, y int) (cell, bool) {
	if y < 0 || y >= len(v.lines) || x < 0 || x >= len(v.lines[y]) {
		return cell{}, false
	}
	return v.lines[y][x], true
}

// Write appends a byte slice into the view's internal buffer. Because
// View implements the io.Writer interface, it can be passed as parameter
// of functions like fmt.Fprintf, fmt.Fprintln, io.Copy, etc. Clear must
// be called to clear the view's buffer.
func (v *View) Write(p []byte) (n int, err error) {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	v.writeRunes(bytes.Runes(p))

	return len(p), nil
}

func (v *View) WriteRunes(p []rune) {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	v.writeRunes(p)
}

// writeRunes copies slice of runes into internal lines buffer.
func (v *View) writeRunes(p []rune) {
	v.tainted = true
	v.clearHover()

	// Fill with empty cells, if writing outside current view buffer
	v.makeWriteable(v.wx, v.wy)

	for _, r := range p {
		switch r {
		case '\n':
			if c, ok := v.readCell(v.wx+1, v.wy); !ok || c.chr == 0 {
				v.writeCells(v.wx, v.wy, []cell{{
					chr:     0,
					fgColor: 0,
					bgColor: 0,
				}})
			}
			v.wx = 0
			v.wy++
			if v.wy >= len(v.lines) {
				v.lines = append(v.lines, nil)
			}
		case '\r':
			if c, ok := v.readCell(v.wx, v.wy); !ok || c.chr == 0 {
				v.writeCells(v.wx, v.wy, []cell{{
					chr:     0,
					fgColor: 0,
					bgColor: 0,
				}})
			}
			v.wx = 0
		default:
			truncateLine, cells := v.parseInput(r, v.wx, v.wy)
			if cells == nil {
				continue
			}
			v.writeCells(v.wx, v.wy, cells)
			if truncateLine {
				length := v.wx + len(cells)
				v.lines[v.wy] = v.lines[v.wy][:length]
			} else {
				v.wx += len(cells)
			}
		}
	}

	v.updateSearchPositions()
}

// exported functions use the mutex. Non-exported functions are for internal use
// and a calling function should use a mutex
func (v *View) WriteString(s string) {
	v.WriteRunes([]rune(s))
}

func (v *View) writeString(s string) {
	v.writeRunes([]rune(s))
}

// parseInput parses char by char the input written to the View. It returns nil
// while processing ESC sequences. Otherwise, it returns a cell slice that
// contains the processed data.
func (v *View) parseInput(ch rune, x int, _ int) (bool, []cell) {
	cells := []cell{}
	truncateLine := false

	isEscape, err := v.ei.parseOne(ch)
	if err != nil {
		for _, r := range v.ei.runes() {
			c := cell{
				fgColor: v.FgColor,
				bgColor: v.BgColor,
				chr:     r,
			}
			cells = append(cells, c)
		}
		v.ei.reset()
	} else {
		repeatCount := 1
		if _, ok := v.ei.instruction.(eraseInLineFromCursor); ok {
			// fill rest of line
			v.ei.instructionRead()
			cx := 0
			for _, cell := range v.lines[v.wy][0:v.wx] {
				cx += runewidth.RuneWidth(cell.chr)
			}
			repeatCount = v.InnerWidth() - cx + 1
			ch = ' '
			truncateLine = true
		} else if isEscape {
			// do not output anything
			return truncateLine, nil
		} else if ch == '\t' {
			// fill tab-sized space
			const tabStop = 4
			ch = ' '
			repeatCount = tabStop - (x % tabStop)
		}
		c := cell{
			fgColor:   v.ei.curFgColor,
			bgColor:   v.ei.curBgColor,
			hyperlink: v.ei.hyperlink,
			chr:       ch,
		}
		for i := 0; i < repeatCount; i++ {
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
	for v.ry < len(v.lines) {
		for v.rx < len(v.lines[v.ry]) {
			count := utf8.EncodeRune(buffer, v.lines[v.ry][v.rx].chr)
			copy(p[offset:], buffer[:count])
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
	v.lines = nil
	v.clearViewLines()
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

	v.clear()

	v.lines = from.lines
	v.viewLines = from.viewLines
	v.ox = from.ox
	v.oy = from.oy
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
	v.lines = nil
}

// This is for when we've done a restart for the sake of avoiding a flicker and
// we've reached the end of the new content to display: we need to clear the remaining
// content from the previous round. We do this by setting v.viewLines to nil so that
// we just render the new content from v.lines directly
func (v *View) FlushStaleCells() {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	v.clearViewLines()
}

func (v *View) rewind() {
	v.ei.reset()

	if err := v.SetReadPos(0, 0); err != nil {
		// SetReadPos returns error only if x and y are negative
		// we are passing 0, 0, thus no error should occur.
		panic(err)
	}
	if err := v.SetWritePos(0, 0); err != nil {
		// SetWritePos returns error only if x and y are negative
		// we are passing 0, 0, thus no error should occur.
		panic(err)
	}
}

func containsUpcaseChar(str string) bool {
	for _, ch := range str {
		if unicode.IsUpper(ch) {
			return true
		}
	}

	return false
}

func (v *View) updateSearchPositions() {
	if v.searcher.searchString != "" {
		var normalizeRune func(r rune) rune
		var normalizedSearchStr string
		// if we have any uppercase characters we'll do a case-sensitive search
		if containsUpcaseChar(v.searcher.searchString) {
			normalizeRune = func(r rune) rune { return r }
			normalizedSearchStr = v.searcher.searchString
		} else {
			normalizeRune = unicode.ToLower
			normalizedSearchStr = strings.ToLower(v.searcher.searchString)
		}

		v.searcher.searchPositions = []SearchPosition{}

		searchPositionsForLine := func(line []cell, y int) []SearchPosition {
			var result []SearchPosition
			searchStringWidth := runewidth.StringWidth(v.searcher.searchString)
			x := 0
			for startIdx, c := range line {
				found := true
				offset := 0
				for _, c := range normalizedSearchStr {
					if len(line)-1 < startIdx+offset {
						found = false
						break
					}
					if normalizeRune(line[startIdx+offset].chr) != c {
						found = false
						break
					}
					offset += 1
				}
				if found {
					result = append(result, SearchPosition{XStart: x, XEnd: x + searchStringWidth, Y: y})
				}
				x += runewidth.RuneWidth(c.chr)
			}
			return result
		}

		if v.searcher.modelSearchResults != nil {
			for _, result := range v.searcher.modelSearchResults {
				if result.Y >= len(v.lines) {
					break
				}

				// If a view line exists for this line index:
				if v.lines[result.Y] != nil {
					// search this view line for the search string
					positions := searchPositionsForLine(v.lines[result.Y], result.Y)
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
			for y, line := range v.lines {
				v.searcher.searchPositions = append(v.searcher.searchPositions, searchPositionsForLine(line, y)...)
			}
		}
	}
}

// IsTainted tells us if the view is tainted
func (v *View) IsTainted() bool {
	return v.tainted
}

// draw re-draws the view's contents.
func (v *View) draw() error {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	if !v.Visible {
		return nil
	}

	v.clearRunes()

	maxX, maxY := v.Size()

	if v.Wrap {
		if maxX == 0 {
			return nil
		}
		v.ox = 0
	}

	v.refreshViewLinesIfNeeded()

	visibleViewLinesHeight := v.viewLineLengthIgnoringTrailingBlankLines()
	if v.Autoscroll && visibleViewLinesHeight > maxY {
		v.oy = visibleViewLinesHeight - maxY
	}

	if len(v.viewLines) == 0 {
		return nil
	}

	start := v.oy
	if start > len(v.viewLines)-1 {
		start = len(v.viewLines) - 1
	}

	emptyCell := cell{chr: ' ', fgColor: ColorDefault, bgColor: ColorDefault}
	var prevFgColor Attribute

	for y, vline := range v.viewLines[start:] {
		if y >= maxY {
			break
		}

		// x tracks the current x position in the view, and cellIdx tracks the
		// index of the cell. If we print a double-sized rune, we increment cellIdx
		// by one but x by two.
		x := -v.ox
		cellIdx := 0

		var c cell
		for {
			if x >= maxX {
				break
			}

			if x < 0 {
				if cellIdx < len(vline.line) {
					x += runewidth.RuneWidth(vline.line[cellIdx].chr)
					cellIdx++
					continue
				} else {
					// no more characters to write so we're only going to be printing empty cells
					// past this point
					x = 0
				}
			}

			// if we're out of cells to write, we'll just print empty cells.
			if cellIdx > len(vline.line)-1 {
				c = emptyCell
				c.fgColor = prevFgColor
			} else {
				c = vline.line[cellIdx]
				// capturing previous foreground colour so that if we're using the reverse
				// attribute we honour the final character's colour and don't awkwardly switch
				// to a new background colour for the remainder of the line
				prevFgColor = c.fgColor
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

			if err := v.setRune(x, y, c.chr, fgColor, bgColor); err != nil {
				return err
			}

			// Not sure why the previous code was here but it caused problems
			// when typing wide characters in an editor
			x += runewidth.RuneWidth(c.chr)
			cellIdx++
		}
	}
	return nil
}

func (v *View) refreshViewLinesIfNeeded() {
	if v.tainted {
		maxX := v.Width()
		lineIdx := 0
		lines := v.lines
		if v.HasLoader {
			lines = v.loaderLines()
		}
		for i, line := range lines {
			wrap := 0
			if v.Wrap {
				wrap = maxX
			}

			ls := lineWrap(line, wrap)
			for j := range ls {
				vline := viewLine{linesX: j, linesY: i, line: ls[j]}

				if lineIdx > len(v.viewLines)-1 {
					v.viewLines = append(v.viewLines, vline)
				} else {
					v.viewLines[lineIdx] = vline
				}
				lineIdx++
			}
		}
		if !v.HasLoader {
			v.tainted = false
		}
	}
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
func (v *View) realPosition(vx, vy int) (x, y int, err error) {
	vx = v.ox + vx
	vy = v.oy + vy

	if vx < 0 || vy < 0 {
		return 0, 0, ErrInvalidPoint
	}

	if len(v.viewLines) == 0 {
		return vx, vy, nil
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

	return x, y, nil
}

// clearRunes erases all the cells in the view.
func (v *View) clearRunes() {
	maxX, maxY := v.Size()
	for x := 0; x < maxX; x++ {
		for y := 0; y < maxY; y++ {
			tcellSetCell(v.x0+x+1, v.y0+y+1, ' ', v.FgColor, v.BgColor, v.outMode)
		}
	}
}

// BufferLines returns the lines in the view's internal
// buffer.
func (v *View) BufferLines() []string {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	lines := make([]string, len(v.lines))
	for i, l := range v.lines {
		str := lineType(l).String()
		str = strings.Replace(str, "\x00", "", -1)
		lines[i] = str
	}
	return lines
}

// Buffer returns a string with the contents of the view's internal
// buffer.
func (v *View) Buffer() string {
	return linesToString(v.lines)
}

// ViewBufferLines returns the lines in the view's internal
// buffer that is shown to the user.
func (v *View) ViewBufferLines() []string {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	lines := make([]string, len(v.viewLines))
	for i, l := range v.viewLines {
		str := lineType(l.line).String()
		str = strings.Replace(str, "\x00", "", -1)
		lines[i] = str
	}
	return lines
}

// LinesHeight is the count of view lines (i.e. lines excluding wrapping)
func (v *View) LinesHeight() int {
	return len(v.lines)
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
	lines := make([][]cell, len(v.viewLines))
	for i := range v.viewLines {
		lines[i] = v.viewLines[i].line
	}

	return linesToString(lines)
}

// Line returns a string with the line of the view's internal buffer
// at the position corresponding to the point (x, y).
func (v *View) Line(y int) (string, error) {
	_, y, err := v.realPosition(0, y)
	if err != nil {
		return "", err
	}

	if y < 0 || y >= len(v.lines) {
		return "", ErrInvalidPoint
	}

	return lineType(v.lines[y]).String(), nil
}

// Word returns a string with the word of the view's internal buffer
// at the position corresponding to the point (x, y).
func (v *View) Word(x, y int) (string, error) {
	x, y, err := v.realPosition(x, y)
	if err != nil {
		return "", err
	}

	if x < 0 || y < 0 || y >= len(v.lines) || x >= len(v.lines[y]) {
		return "", ErrInvalidPoint
	}

	str := lineType(v.lines[y]).String()

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
	return str[nl:nr], nil
}

// indexFunc allows to split lines by words taking into account spaces
// and 0.
func indexFunc(r rune) bool {
	return r == ' ' || r == 0
}

// SetHighlight toggles highlighting of separate lines, for custom lists
// or multiple selection in views.
func (v *View) SetHighlight(y int, on bool) error {
	if y < 0 || y >= len(v.lines) {
		err := ErrInvalidPoint
		return err
	}

	line := v.lines[y]
	cells := make([]cell, 0)
	for _, c := range line {
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
	v.lines[y] = cells
	v.clearHover()
	return nil
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
		rw := runewidth.RuneWidth(currChr)
		n += rw
		// if currChr == 'g' {
		// 	panic(n)
		// }
		if n > columns {
			// This code is convoluted but we've got comprehensive tests so feel free to do whatever you want
			// to the code to simplify it so long as our tests still pass.
			if currChr == ' ' {
				// if the line ends in a space, we'll omit it. This means there'll be no
				// way to distinguish between a clean break and a mid-word break, but
				// I think it's worth it.
				lines = append(lines, line[offset:i])
				offset = i + 1
				n = 0
			} else if currChr == '-' {
				// if the last character is hyphen and the width of line is equal to the columns
				lines = append(lines, line[offset:i])
				offset = i
				n = rw
			} else if lastWhitespaceIndex != -1 && lastWhitespaceIndex+1 != i {
				// if there is a space in the line and the line is not breaking at a space/hyphen
				if line[lastWhitespaceIndex].chr == '-' {
					// if break occurs at hyphen, we'll retain the hyphen
					lines = append(lines, line[offset:lastWhitespaceIndex+1])
					offset = lastWhitespaceIndex + 1
					n = i - offset
				} else {
					// if break occurs at space, we'll omit the space
					lines = append(lines, line[offset:lastWhitespaceIndex])
					offset = lastWhitespaceIndex + 1
					n = i - offset + 1
				}
			} else {
				// in this case we're breaking mid-word
				lines = append(lines, line[offset:i])
				offset = i
				n = rw
			}
			lastWhitespaceIndex = -1
		} else if line[i].chr == ' ' || line[i].chr == '-' {
			lastWhitespaceIndex = i
		}
	}

	lines = append(lines, line[offset:])
	return lines
}

func linesToString(lines [][]cell) string {
	str := make([]string, len(lines))
	for i := range lines {
		rns := make([]rune, 0, len(lines[i]))
		line := lineType(lines[i]).String()
		for _, c := range line {
			if c != '\x00' {
				rns = append(rns, c)
			}
		}
		str[i] = string(rns)
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
		charX += runewidth.StringWidth(tab)
		if x <= charX {
			return i
		}
		charX += runewidth.StringWidth(" - ")
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

// expected to only be used in tests
func (v *View) SelectedLine() string {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	if len(v.lines) == 0 {
		return ""
	}

	return v.lineContentAtIdx(v.SelectedLineIdx())
}

// expected to only be used in tests
func (v *View) SelectedLines() []string {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	if len(v.lines) == 0 {
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
	line := v.lines[idx]
	str := lineType(line).String()
	return strings.Replace(str, "\x00", "", -1)
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
	} else {
		return start, end
	}
}

func (v *View) RenderTextArea() {
	v.Clear()
	fmt.Fprint(v, v.TextArea.GetContent())
	cursorX, cursorY := v.TextArea.GetCursorXY()
	prevOriginX, prevOriginY := v.Origin()
	width, height := v.InnerWidth(), v.InnerHeight()

	newViewCursorX, newOriginX := updatedCursorAndOrigin(prevOriginX, width, cursorX)
	newViewCursorY, newOriginY := updatedCursorAndOrigin(prevOriginY, height, cursorY)

	_ = v.SetCursor(newViewCursorX, newViewCursorY)
	_ = v.SetOrigin(newOriginX, newOriginY)
}

func updatedCursorAndOrigin(prevOrigin int, size int, cursor int) (int, int) {
	var newViewCursor int
	newOrigin := prevOrigin

	if cursor > prevOrigin+size {
		newOrigin = cursor - size
		newViewCursor = size
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
	_ = v.SetOrigin(0, 0)
	_ = v.SetCursor(0, 0)
}

func (v *View) overwriteLines(y int, content string) {
	// break by newline, then for each line, write it, then add that erase command
	v.wx = 0
	v.wy = y
	v.clearViewLines()

	lines := strings.Replace(content, "\n", "\x1b[K\n", -1)
	// If the last line doesn't end with a linefeed, add the erase command at
	// the end too
	if !strings.HasSuffix(lines, "\n") {
		lines += "\x1b[K"
	}
	v.writeString(lines)
}

// only call this function if you don't care where v.wx and v.wy end up
func (v *View) OverwriteLines(y int, content string) {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	v.overwriteLines(y, content)
}

// only call this function if you don't care where v.wx and v.wy end up
func (v *View) OverwriteLinesAndClearEverythingElse(y int, content string) {
	v.writeMutex.Lock()
	defer v.writeMutex.Unlock()

	v.overwriteLines(y, content)

	for i := 0; i < y; i += 1 {
		v.lines[i] = nil
	}

	for i := v.wy + 1; i < len(v.lines); i += 1 {
		v.lines[i] = nil
	}
}

func (v *View) SetContentLineCount(lineCount int) {
	if lineCount > 0 {
		v.makeWriteable(0, lineCount-1)
	}
	v.lines = v.lines[:lineCount]
}

func (v *View) ScrollUp(amount int) {
	if amount > v.oy {
		amount = v.oy
	}

	if amount != 0 {
		v.oy -= amount
		v.cy += amount

		v.clearHover()
	}
}

// ensures we don't scroll past the end of the view's content
func (v *View) ScrollDown(amount int) {
	adjustedAmount := v.adjustDownwardScrollAmount(amount)
	if adjustedAmount > 0 {
		v.oy += adjustedAmount
		v.cy -= adjustedAmount

		v.clearHover()
	}
}

func (v *View) ScrollLeft(amount int) {
	newOx := v.ox - amount
	if newOx < 0 {
		newOx = 0
	}
	if newOx != v.ox {
		v.ox = newOx

		v.clearHover()
	}
}

// not applying any limits to this
func (v *View) ScrollRight(amount int) {
	v.ox += amount

	v.clearHover()
}

func (v *View) adjustDownwardScrollAmount(scrollHeight int) int {
	_, oy := v.Origin()
	y := oy
	if !v.CanScrollPastBottom {
		_, sy := v.Size()
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
	} else {
		return scrollHeight
	}
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
	} else {
		return 0
	}
}

// Returns true if the view contains a line containing the given text with the given
// foreground color
func (v *View) ContainsColoredText(fgColor string, text string) bool {
	for _, line := range v.lines {
		if containsColoredTextInLine(fgColor, text, line) {
			return true
		}
	}

	return false
}

func containsColoredTextInLine(fgColorStr string, text string, line []cell) bool {
	fgColor := tcell.GetColor(fgColorStr)

	currentMatch := ""
	for i := 0; i < len(line); i++ {
		cell := line[i]

		// stripping attributes by converting to and from hex
		cellColor := tcell.NewHexColor(cell.fgColor.Hex())

		if cellColor == fgColor {
			currentMatch += string(cell.chr)
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
