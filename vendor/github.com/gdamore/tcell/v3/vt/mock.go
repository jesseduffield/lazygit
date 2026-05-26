// Copyright 2026 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vt

import (
	"slices"
	"sync"
	"time"

	"github.com/gdamore/tcell/v3/color"
	"github.com/gdamore/tcell/v3/tty"
)

// mockTerm implements MockTerm.
type mockTerm struct {
	mb MockBackend
	em Emulator
	ks *KeyboardState
}

// Stop the terminal.
func (mt *mockTerm) Stop() error {
	return mt.em.Stop()
}

// Start the terminal.
func (mt *mockTerm) Start() error {
	return mt.em.Start()
}

// Drain all output from the terminal, ensuring
// any queued commands are processed.
func (mt *mockTerm) Drain() error {
	return mt.em.Drain()
}

// Read data from the terminal. This is called by a terminal
// application (e.g. via tcell Tty.)  Read data will include
// key strokes, mouse events, and responses to terminal queries.
func (mt *mockTerm) Read(data []byte) (int, error) {
	return mt.em.Read(data)
}

// Write data to the terminal, typically either commands or data
// that should be displayed on the virtual screen.
func (mt *mockTerm) Write(b []byte) (n int, err error) {
	return mt.em.Write(b)
}

// WindowSize obtains the dimensions of the window.
func (mt *mockTerm) WindowSize() (tty.WindowSize, error) {
	sz := mt.mb.GetSize()
	// No pixel sizes for now
	return tty.WindowSize{Width: int(sz.X), Height: int(sz.Y)}, nil
}

// NotifyResize registers a channel to be signaled when a resize has occurred.
// In real terminal emulators this would be posted (non-blocking) by a signal handler.
func (mt *mockTerm) NotifyResize(resizeq chan<- bool) {
	if rs, ok := mt.mb.(Resizer); ok {
		rs.NotifyResize(resizeq)
	}
}

// Close closes the terminal, after which it should no longer be used. Stop is implied.
func (mt *mockTerm) Close() error {
	return mt.Stop()
}

// Pos returns the cursor position.
func (mt *mockTerm) Pos() Coord {
	return mt.mb.GetPosition()
}

// GetCell returns the contents of the cell at the given coordinates, or a zero value
// if the coordinates are out of range.
func (mt *mockTerm) GetCell(pos Coord) Cell {
	return mt.mb.GetCell(pos)
}

// Bells counts the number of times the bell has rung.
func (mt *mockTerm) Bells() int {
	return mt.mb.Bells()
}

// KeyEvent is used to inject a key event.  Call this to inject
// a synthetic, fully specified key event.  Most uses should just use
// the KeyPress, KeyRelease, or even simpler KeyTap APIs.
func (mt *mockTerm) KeyEvent(ev KeyEvent) {
	mt.em.KeyEvent(ev)
	if ev.Key == KeyEsc {
		// Inject a delay to simulate human typing.
		// Necessary to disambiguate Escape from other sequences.
		time.Sleep(time.Millisecond * 150)
	}
}

// KeyPress implements MockTerm.KeyPress.
func (mt *mockTerm) KeyPress(k Key) {
	if event := mt.ks.Pressed(k); event != nil {
		mt.KeyEvent(*event)
	}
}

// KeyRelease implements MockTerm.KeyRelease.
func (mt *mockTerm) KeyRelease(k Key) {
	if event := mt.ks.Released(k); event != nil {
		mt.KeyEvent(*event)
	}
}

// KeyTap implements MockTerm.KeyTap.
func (mt *mockTerm) KeyTap(keys ...Key) {
	for _, k := range keys {
		mt.KeyPress(k)
	}
	for _, k := range slices.Backward(keys) {
		mt.KeyRelease(k)
	}
}

// SetRepeat sets the repeat interval for the keyboard.
// Set the interval to zero to disable repeat.
func (mt *mockTerm) SetRepeat(delay, interval time.Duration) {
	mt.ks.SetRepeat(delay, interval)
}

// MouseEvent implements MockTerm.MouseEvent.
func (mt *mockTerm) MouseEvent(ev MouseEvent) {
	mt.em.MouseEvent(ev)
}

// FocusEvent implements MockTerm.FocusEvent.
func (mt *mockTerm) FocusEvent(focused bool) {
	mt.em.FocusEvent(focused)
}

// GetTitle returns the current window title.
func (mt *mockTerm) GetTitle() string {
	return mt.mb.GetTitle()
}

// SetSize is used to change the terminal size.
func (mt *mockTerm) SetSize(size Coord) {
	mt.mb.SetSize(size)
	mt.em.ResizeEvent(size)
}

// Backend returns the backend for testing.
func (mt *mockTerm) Backend() MockBackend {
	return mt.mb
}

// SendRaw is used to inject raw bytes to the read stream of the app.
// Use this for fuzz testing.
func (mt *mockTerm) SendRaw(data []byte) {
	mt.em.SendRaw(data)
}

// SetLayout sets the keyboard layout.
func (mt *mockTerm) SetLayout(km *Layout) {
	mt.ks.SetLayout(km)
}

// MockTerm is a mock terminal (emulator).  It can be used to
// test the emulator itself, or to test applications (or tcell) that
// uses the terminal.  It also implements the Tty interface used
// by tcell itself.
type MockTerm interface {
	tty.Tty

	// Pos reports the current cursor position.
	Pos() Coord

	// GetCell returns the cell at the given coordinates.
	// The coordinates must be valid.
	GetCell(Coord) Cell

	// Bells returns the number of times the bell has been rung.
	Bells() int

	// Inject a keyboard event - this is a full event, and bypasses
	// the layout and keyboard state processor.
	KeyEvent(KeyEvent)

	// Inject a key press.
	KeyPress(Key)

	// Inject a key release.
	KeyRelease(Key)

	// SetRepeat configures keyboard repeating. Repeat keystrokes
	// will be assumed after the key has been held for at least delay,
	// with new keys added each interval.
	SetRepeat(delay, interval time.Duration)

	// Inject one or more key press and releases.
	// The keys are pressed in the order, and released in reverse order.
	// Thus modifiers should be listed first.  This should not be used
	// to simulate typing a sequence (e.g. a word), but if you wanted to
	// test say N-Key rollover you could do that here.
	KeyTap(...Key)

	// Inject a mouse event.
	MouseEvent(MouseEvent)

	// Inject a focus event.
	FocusEvent(bool)

	// GetTitle obtains the current window title.
	GetTitle() string

	// SetSize is used to resize the terminal.
	SetSize(Coord)

	// SendRaw is used to send raw data to the application.
	// This is mostly intended to facilitate fuzz testing the application.
	SendRaw([]byte)

	// Backend returns the backend (used for testing).
	Backend() MockBackend

	// SetLayout sets the keyboard layout to use.
	// If not specified, a US standard ANSI keyboard will be assumed.
	SetLayout(*Layout)
}

type noMockBlit struct {
	MockBackend
	Blit struct{} // prevents use as Blitter
}

// NewMockTerm gives a mock terminal emulator.
func NewMockTerm(opts ...MockOpt) MockTerm {
	mt := &mockTerm{}
	mt.mb = NewMockBackend(opts...)
	var be MockBackend = mt.mb
	emOpts := []EmulatorOpt{}
	for _, o := range opts {
		switch o.(type) {
		case MockOptNoBlit:
			be = &noMockBlit{be, struct{}{}}
		case MockOpt8BitControls:
			emOpts = append(emOpts, EmulatorOpt8BitControls{})
		}
	}
	mt.em = NewEmulator(be, emOpts...)
	mt.em.SetId("TCellMock", "1.0")
	mt.ks = &KeyboardState{}
	return mt
}

// MockBackend provides additional mock-specific capabilities on top of Backend.
// This is meant to facilitate test cases
type MockBackend interface {
	Backend

	// GetCell returns the cell at the given position, or an empty cell if the
	// position is out of the bounds of the window.
	GetCell(Coord) Cell

	// Bells counts the number of bells rung.
	Bells() int

	// GetTitle gets the current window title.
	GetTitle() string

	// SetSize is used to resize the window.
	// Newly added cells are empty, and content in old cells that out of range is lost.
	SetSize(Coord)

	// GetCursor is used to obtain the current cursor style.
	GetCursor() CursorStyle

	// SetClipboard sets the clipboard contents (copy buffer).
	SetClipboard([]byte)

	// GetClipboard returns the clipboard (copy buffer).
	GetClipboard() []byte

	// IsAdvancedKeyboard returns true, as we always support the full keyboard protocol.
	IsAdvancedKeyboard() bool
}

// mockBackend is a mock of a backend device for use with the emulator.
// It implements the following interfaces:
// vt.Backend, vt.Beeper, vt.Colorer, vt.Titler, vt.Resizer, vt.Blitter
type mockBackend struct {
	cells        []Cell // Content of cells
	size         Coord
	pos          Coord
	colors       int
	style        Style
	defaultStyle Style
	notifyQ      chan<- bool
	resized      bool
	newSize      Coord
	modes        map[PrivateMode]ModeStatus
	bells        int
	errs         int
	title        string
	clipboard    []byte
	cursor       CursorStyle
	lock         sync.Mutex
}

func (mb *mockBackend) GetSize() Coord {
	mb.lock.Lock()
	defer mb.lock.Unlock()
	mb.checkSize()
	return mb.size
}

func (mb *mockBackend) Beep() {
	mb.lock.Lock()
	mb.bells++
	mb.lock.Unlock()
}

func (mb *mockBackend) SetMouse(MouseReporting) {}

func (mb *mockBackend) GetPrivateMode(pm PrivateMode) ModeStatus {
	mb.lock.Lock()
	defer mb.lock.Unlock()
	// note default (zero) value is ModeNA
	return mb.modes[pm]
}

func (mb *mockBackend) SetPrivateMode(pm PrivateMode, status ModeStatus) error {
	mb.lock.Lock()
	defer mb.lock.Unlock()
	if old := mb.modes[pm]; old == ModeOn || old == ModeOff {
		if status == ModeOn || status == ModeOff {
			mb.modes[pm] = status
		} else {
			mb.errs++
		}
	} else {
		mb.errs++
	}
	return nil
}

func (mb *mockBackend) Put(pos Coord, cell Cell) { // grapheme string, width int, style Style) {
	mb.lock.Lock()
	defer mb.lock.Unlock()
	mb.checkSize()

	if index := mb.index(pos); index >= 0 {
		mb.cells[index] = cell

		// writing to a cell right after a wide
		// character clears that wide character (but leaves style/attributes)
		if cell.W > 0 && pos.X > 0 && mb.cells[index-1].W > 1 {
			mb.cells[index-1].C = ""
			mb.cells[index-1].W = 0
		}

		// wide characters delete the next cell
		if cell.W == 2 && pos.X < mb.size.X-1 {
			mb.cells[index+1].C = ""
			mb.cells[index+1].W = 0
			mb.cells[index+1].S = cell.S
		}
	} else {
		mb.errs++
	}
}

func (mb *mockBackend) isPositionValid(pos Coord) bool {
	mb.checkSize()

	return pos.X < mb.size.X && pos.Y < mb.size.Y && pos.X >= 0 && pos.Y >= 0
}

// index calculates the index in the cells array.  If the coordinates are invalid,
// -1 will be returned.
func (mb *mockBackend) index(pos Coord) int {
	mb.checkSize()

	if !mb.isPositionValid(pos) {
		return -1
	}
	return int(pos.X) + int(pos.Y)*int(mb.size.X)
}

func (mb *mockBackend) GetCell(pos Coord) Cell {
	mb.lock.Lock()
	defer mb.lock.Unlock()

	if index := mb.index(pos); index >= 0 {
		return mb.cells[index]
	}
	return Cell{S: BaseStyle}
}

func (mb *mockBackend) Bells() int {
	mb.lock.Lock()
	defer mb.lock.Unlock()
	return mb.bells
}

func (mb *mockBackend) GetPosition() Coord {
	mb.lock.Lock()
	defer mb.lock.Unlock()
	mb.checkSize()
	return mb.pos
}

func (mb *mockBackend) SetPosition(pos Coord) {
	mb.lock.Lock()
	defer mb.lock.Unlock()
	mb.checkSize()
	pos.X = min(mb.size.X-1, max(0, pos.X))
	pos.Y = min(mb.size.Y-1, max(0, pos.Y))
	mb.pos = pos
}

func (mb *mockBackend) Colors() int {
	return mb.colors
}

func (mb *mockBackend) SetStyle(style Style) {
	mb.lock.Lock()
	defer mb.lock.Unlock()

	mb.style = style
}

// SetWindowTitle implements the Titler interface.
func (mb *mockBackend) SetWindowTitle(title string) {
	mb.lock.Lock()
	defer mb.lock.Unlock()

	mb.title = title
}

// GetTitle allows test code to observe what was set with SetWindowTitle.
func (mb *mockBackend) GetTitle() string {
	mb.lock.Lock()
	defer mb.lock.Unlock()

	return mb.title
}

// NotifyResize registers a channel to be written to (non-blocking) if the
// backend changes size.
func (mb *mockBackend) NotifyResize(rq chan<- bool) {
	mb.lock.Lock()
	defer mb.lock.Unlock()

	mb.notifyQ = rq
}

// checkSize performs a possible terminal resize. Cells that are
// added are treated as empty, while cells that are removed are just lost.
// (Note that at least one other emulator erases content on a resize.  There is no
// standard for what to do here.) This is done inline when calculating the index.
// The caller is expected to hold mb.lock.
func (mb *mockBackend) checkSize() {
	if !mb.resized {
		return
	}
	size := mb.newSize
	old := mb.cells
	ox := int(mb.size.X)
	oy := int(mb.size.Y)
	nx := int(size.X)
	ny := int(size.Y)
	cells := make([]Cell, int(size.Y)*int(size.X))
	for i := range cells {
		cells[i].S = BaseStyle
	}
	for y := range min(ny, oy) {
		for x := range min(nx, ox) {
			cells[y*nx+x] = old[y*ox+x]
		}
	}
	mb.cells = cells
	mb.size = size
	mb.pos.X = min(mb.pos.X, size.X-1)
	mb.pos.Y = min(mb.pos.Y, size.Y-1)
	mb.resized = false
}

func (mb *mockBackend) RaiseResize() {
	mb.lock.Lock()
	defer mb.lock.Unlock()

	if rq := mb.notifyQ; rq != nil {
		select {
		case rq <- true:
		default:
		}
	}
}

// SetSize is used to change the size of the virtual terminal.
func (mb *mockBackend) SetSize(size Coord) {
	mb.lock.Lock()
	mb.resized = true
	mb.newSize = size
	mb.lock.Unlock()
}

// Reset the terminal to startup defaults.
func (mb *mockBackend) Reset() {
	mb.lock.Lock()
	defer mb.lock.Unlock()

	mb.style = mb.defaultStyle

	mb.title = ""
	mb.errs = 0
	mb.bells = 0
	mb.pos = Coord{X: 0, Y: 0}
	mb.modes[PmShowCursor] = ModeOn
	mb.modes[PmBlinkCursor] = ModeOn
	mb.modes[PmGraphemeClusters] = ModeOff
}

func (mb *mockBackend) Blit(src, dst, dim Coord) {
	mb.lock.Lock()
	defer mb.lock.Unlock()

	mb.checkSize()

	// clip to visible source
	if dim.X+src.X > mb.size.X {
		dim.X = mb.size.X - src.X
	}
	if dim.Y+src.Y > mb.size.Y {
		dim.Y = mb.size.Y - src.Y
	}
	// and clip to final destination
	if dim.X+dst.X > mb.size.X {
		dim.X = mb.size.X - dst.X
	}
	if dim.Y+dst.Y > mb.size.Y {
		dim.Y = mb.size.Y - dst.Y
	}

	// gap represents decrement when shifting to the next row --
	// skipping over the irrelevant cells. (The increment in the
	// index when going from last cell of row to first cell of next row,
	// or vice versa.)
	gap := int(mb.size.X - dim.X)

	// the following logic is carefully constructed to avoid expensive
	// operations in the loops (only addition or subtraction)
	if mb.index(src) > mb.index(dst) { // source appears later, so we can forward copy
		si := mb.index(src)
		di := mb.index(dst)
		for range dim.Y {
			for range dim.X {
				mb.cells[di] = mb.cells[si]
				di++
				si++
			}
			// advance to next row
			si += gap
			di += gap
		}
	} else { // source appears earlier, so we have to reverse copy
		src.Y += dim.Y - 1
		dst.Y += dim.Y - 1
		src.X += dim.X - 1
		dst.X += dim.X - 1
		si := mb.index(src)
		di := mb.index(dst)

		for range dim.Y {
			for range dim.X {
				mb.cells[di] = mb.cells[si]
				si--
				di--
			}
			si -= gap
			di -= gap
		}
	}
}

// Buffering is not supported by the mockBackend, and there is little point in it.
func (mb *mockBackend) Buffering(bool) {}

// SetCursor is used to set how the cursor is displayed.
func (mb *mockBackend) SetCursor(cs CursorStyle) {
	mb.lock.Lock()
	defer mb.lock.Unlock()
	mb.cursor = cs
}

// GetCursor returns the current cursor style.
func (mb *mockBackend) GetCursor() CursorStyle {
	mb.lock.Lock()
	defer mb.lock.Unlock()
	return mb.cursor
}

// SetClipboard sets the current clipboard contents.
func (mb *mockBackend) SetClipboard(data []byte) {
	mb.lock.Lock()
	defer mb.lock.Unlock()
	mb.clipboard = data
}

// GetClipboard gets the current clipboard contents.
func (mb *mockBackend) GetClipboard() []byte {
	mb.lock.Lock()
	defer mb.lock.Unlock()
	return mb.clipboard
}

// IsAdvancedKeyboard returns true - we always implement
// the raw keyboard protocol.
func (mb *mockBackend) IsAdvancedKeyboard() bool { return true }

// MockOpt is an interface by which options can change the behavior of the mocked terminal.
// This is intended to permit easier testing.
type MockOpt interface{ SetMockOpt(mb *mockBackend) }

// MockOptSize changes the default terminal size, which is normally 80x24.
type MockOptSize Coord

func (o MockOptSize) SetMockOpt(mb *mockBackend) { mb.size = Coord(o) }

// MockOptColors changes the number of colors the terminal supports.
type MockOptColors int

func (o MockOptColors) SetMockOpt(mb *mockBackend) { mb.colors = int(o) }

// MockOptNoBlit suppresses the blitter interface.
type MockOptNoBlit struct{}

func (MockOptNoBlit) SetMockOpt(mb *mockBackend) {}

// MockOpt8BitControls enables raw 8-bit and UTF-8 encoded C1 controls in the
// emulator. The default is to accept only 7-bit ESC-prefixed controls.
type MockOpt8BitControls struct{}

func (MockOpt8BitControls) SetMockOpt(mb *mockBackend) {}

// NewMockBackend returns a MockBackend modified by the given options.
// The default is a fully featured 256-color backend with initial size 80x24.
func NewMockBackend(options ...MockOpt) MockBackend {
	mb := &mockBackend{
		size:         Coord{X: 80, Y: 24},
		colors:       256,
		style:        BaseStyle,
		defaultStyle: BaseStyle.WithFg(color.Silver).WithBg(color.Black),
		cursor:       BlinkingBlock,
	}

	for _, opt := range options {
		opt.SetMockOpt(mb)
	}

	if mb.colors > 0 {
		mb.style = mb.defaultStyle
	}
	mb.cells = make([]Cell, int(mb.size.X)*int(mb.size.Y))
	for i := range mb.cells {
		mb.cells[i].S = BaseStyle
	}

	mb.modes = make(map[PrivateMode]ModeStatus)
	mb.modes[PmShowCursor] = ModeOn
	mb.modes[PmBlinkCursor] = ModeOn
	mb.modes[PmGraphemeClusters] = ModeOff
	mb.modes[PmSyncOutput] = ModeOff
	return mb
}
