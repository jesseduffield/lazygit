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

package tcell

import (
	"sync"

	"github.com/gdamore/tcell/v3/color"
)

// Screen represents the physical (or emulated) screen.
// This can be a terminal window or a physical console.  Platforms implement
// this differently.
type Screen interface {
	// Init initializes the screen for use.
	Init() error

	// Fini finalizes the screen also releasing resources.
	Fini()

	// Clear logically erases the screen.
	// This is effectively a short-cut for Fill(' ', StyleDefault).
	Clear()

	// Fill fills the screen with the given character and style.
	// The effect of filling the screen is not visible until Show
	// is called (or Sync).
	Fill(rune, Style)

	// Put writes the first grapheme of the given string with th
	// given style at the given coordinates. (Only the first grapheme
	// occupying either one or two cells is stored.) It returns the
	// remainder of the string, and the width displayed.
	Put(x int, y int, str string, style Style) (string, int)

	// PutStr writes a string starting at the given position, using the
	// default style. The content is clipped to the screen dimensions.
	PutStr(x int, y int, str string)

	// PutStrStyled writes a string starting at the given position, using
	// the given style. The content is clipped to the screen dimensions.
	PutStrStyled(x int, y int, str string, style Style)

	// Get the contents at the given location.  If the
	// coordinates are out of range, then the values will be 0, nil,
	// StyleDefault.  Note that the contents returned are logical contents
	// and may not actually be what is displayed, but rather are what will
	// be displayed if Show() or Sync() is called.  The width is the width
	// in screen cells; most often this will be 1, but some East Asian
	// characters and emoji require two cells.
	Get(x, y int) (str string, style Style, width int)

	// SetContent sets the contents of the given cell location.  If
	// the coordinates are out of range, then the operation is ignored.
	//
	// The first rune is the primary non-zero width rune.  The array
	// that follows is a possible list of combining characters to append,
	// and will usually be nil (no combining characters.)
	//
	// The results are not displayed until Show() or Sync() is called.
	//
	// Note that wide (East Asian full width and emoji) runes occupy two cells,
	// and attempts to place character at next cell to the right will have
	// undefined effects.  Wide runes that are printed in the
	// last column will be replaced with a single width space on output.
	SetContent(x int, y int, primary rune, combining []rune, style Style)

	// SetStyle sets the default style to use when clearing the screen
	// or when StyleDefault is specified.  If it is also StyleDefault,
	// then whatever system/terminal default is relevant will be used.
	SetStyle(style Style)

	// ShowCursor is used to display the cursor at a given location.
	// If the coordinates -1, -1 are given or are otherwise outside the
	// dimensions of the screen, the cursor will be hidden.
	ShowCursor(x int, y int)

	// HideCursor is used to hide the cursor.  It's an alias for
	// ShowCursor(-1, -1).sim
	HideCursor()

	// SetCursorStyle is used to set the cursor style.  If the style
	// is not supported (or cursor styles are not supported at all),
	// then this will have no effect.  Color will be changed if supplied,
	// and the terminal supports doing so.
	SetCursorStyle(CursorStyle, ...color.Color)

	// Size returns the screen size as width, height.  This changes in
	// response to a call to Clear or Flush.
	Size() (width, height int)

	// EventQ returns the channel of events, and is usable just like
	// any other channel.  Events can be injected by writing to
	// the channel, and they can be read by reading from it.  The
	// channel will remain open until the screen is completely shut down
	// with Fini().  Consequently, applications must not write to this
	// channel after Fini() is called.
	EventQ() chan Event

	// EnableMouse enables the mouse.  (If your terminal supports it.)
	// If no flags are specified, then all events are reported, if the
	// terminal supports them.
	EnableMouse(...MouseFlags)

	// DisableMouse disables the mouse.
	DisableMouse()

	// EnablePaste enables bracketed paste mode, if supported.
	EnablePaste()

	// DisablePaste disables bracketed paste mode.
	DisablePaste()

	// EnableFocus enables reporting of focus events, if your terminal supports it.
	EnableFocus()

	// DisableFocus disables reporting of focus events.
	DisableFocus()

	// Colors returns the number of colors.  All colors are assumed to
	// use the ANSI color map.  If a terminal is monochrome, it will
	// return 0.
	Colors() int

	// Show makes all the content changes made using SetContent() visible
	// on the display.
	//
	// It does so in the most efficient and least visually disruptive
	// manner possible.
	Show()

	// Sync works like Show(), but it updates every visible cell on the
	// physical display, assuming that it is not synchronized with any
	// internal model.  This may be both expensive and visually jarring,
	// so it should only be used when believed to actually be necessary.
	//
	// Typically, this is called as a result of a user-requested redraw
	// (e.g. to clear up on-screen corruption caused by some other program),
	// or during a resize event.
	Sync()

	// CharacterSet returns information about the character set.
	// This isn't the full locale, but it does give us the input/output
	// character set.  Note that this is just for diagnostic purposes,
	// we normally translate input/output to/from UTF-8, regardless of
	// what the user's environment is.
	CharacterSet() string

	// RegisterRuneFallback adds a fallback for runes that are not
	// part of the character set -- for example one could register
	// o as a fallback for ø.  This should be done cautiously for
	// characters that might be displayed ordinarily in language
	// specific text -- characters that could change the meaning of
	// written text would be dangerous.  The intention here is to
	// facilitate fallback characters in pseudo-graphical applications.
	//
	// If the terminal has fallbacks already in place via an alternate
	// character set, those are used in preference.  Also, standard
	// fallbacks for graphical characters in the alternate character set
	// terminfo string are registered implicitly.
	//
	// The display string should be the same width as original rune.
	// This makes it possible to register two character replacements
	// for full width East Asian characters, for example.
	//
	// It is recommended that replacement strings consist only of
	// 7-bit ASCII, since other characters may not display everywhere.
	RegisterRuneFallback(r rune, subst string)

	// UnregisterRuneFallback unmaps a replacement.  It will unmap
	// the implicit ASCII replacements for alternate characters as well.
	// When an unmapped char needs to be displayed, but no suitable
	// glyph is available, '?' is emitted instead.  It is not possible
	// to "disable" the use of alternate characters that are supported
	// by your terminal except by changing the terminal database.
	UnregisterRuneFallback(r rune)

	// Resize does nothing, since it's generally not possible to
	// ask a screen to resize, but it allows the Screen to implement
	// the View interface.
	Resize(int, int, int, int)

	// Suspend pauses input and output processing.  It also restores the
	// terminal settings to what they were when the application started.
	// This can be used to, for example, run a sub-shell.
	Suspend() error

	// Resume resumes after Suspend().
	Resume() error

	// Beep attempts to sound an OS-dependent audible alert and returns an error
	// when unsuccessful.
	Beep() error

	// SetSize attempts to resize the window.  It also invalidates the cells and
	// calls the resize function.  Note that if the window size is changed, it will
	// not be restored upon application exit.
	//
	// Many terminals cannot support this.  Perversely, the "modern" Windows Terminal
	// does not support application-initiated resizing, whereas the legacy terminal does.
	// Also, some emulators can support this but may have it disabled by default.
	SetSize(int, int)

	// LockRegion sets or unsets a lock on a region of cells. A lock on a
	// cell prevents the cell from being redrawn.
	LockRegion(x, y, width, height int, lock bool)

	// Tty returns the underlying Tty. If the screen is not a terminal, the
	// returned bool will be false
	Tty() (Tty, bool)

	// SetTitle sets a window title on the screen.
	// Terminals may be configured to ignore this, or unable to.
	// Tcell may attempt to save and restore the window title on entry and exit, but
	// the results may vary.  Use of unicode characters may not be supported.
	SetTitle(string)

	// SetClipboard is used to post arbitrary data to the system clipboard.
	// This need not be UTF-8 string data.  It's up to the recipient to decode the
	// data meaningfully.  Terminals may prevent this for security reasons.
	// An empty byte or nil can be used to clear the clipboard.
	SetClipboard([]byte)

	// GetClipboard is used to request the clipboard contents.  It may be ignored.
	// If the terminal is willing, it will be post the clipboard contents using an
	// EventPaste with the clipboard content as the Data() field.  Terminals may
	// prevent this for security reasons.
	GetClipboard()

	// HasClipboard is true if the screen claims to support the clipboard.
	// Note that GetClipboard may still not work, but SetClipboard should be functional.
	// Note that many terminals that support the clipboard don't actually report that they
	// do, so a false indication is not necessarily conclusive.
	HasClipboard() bool

	// ShowNotification is used to show a desktop notification, when the terminal
	// supports it.  Right now only terminals supporting OSC 777 support this.
	ShowNotification(title string, body string)

	// Terminal returns the terminal name and version if known.  If either of these
	// are unknown, then empty strings are returned in their place.  This is intended
	// to facilitate debug, and also applications that wish to enable very specific
	// behaviors for the terminal
	Terminal() (string, string)
}

var overrideScreen chan Screen
var overrideOnce sync.Once

// NewScreen returns a default Screen suitable for the user's terminal environment.
func NewScreen() (Screen, error) {

	// Allow an application (presumably test code) to inject a replacement default
	// screen.  This could also be used to create shims for things like nesting screens.
	select {
	case s := <-overrideScreen:
		return s, nil
	default:
	}

	if s, e := NewTerminfoScreen(); s != nil {
		return s, nil
	} else {
		return nil, e
	}
}

// ShimScreen allows an application to override the screen that will
// be returned by NewScreen.  Typically this  is used for testing,
// where the test code calls this once before running an example.
// It could also be used to intercept a regular Screen.
func ShimScreen(s Screen) {
	overrideOnce.Do(func() {
		overrideScreen = make(chan Screen, 8) // normally would only be one anyway
	})
	overrideScreen <- s
}

// MouseFlags are options to modify the handling of mouse events.
// Actual events can be ORed together.
type MouseFlags int

const (
	MouseButtonEvents = MouseFlags(1) // Click events only
	MouseDragEvents   = MouseFlags(2) // Click-drag events (includes button events)
	MouseMotionEvents = MouseFlags(4) // All mouse events (includes click and drag events)
)

// CursorStyle represents a given cursor style, which can include the shape and
// whether the cursor blinks or is solid.  Support for changing this is not universal.
type CursorStyle int

const (
	CursorStyleDefault = CursorStyle(iota) // The default
	CursorStyleBlinkingBlock
	CursorStyleSteadyBlock
	CursorStyleBlinkingUnderline
	CursorStyleSteadyUnderline
	CursorStyleBlinkingBar
	CursorStyleSteadyBar
)

// screenImpl is a subset of Screen that can be used with baseScreen to formulate
// a complete implementation of Screen.  See Screen for doc comments about methods.
type screenImpl interface {
	Init() error
	Fini()
	SetStyle(style Style)
	ShowCursor(x int, y int)
	HideCursor()
	SetCursor(CursorStyle, color.Color)
	Size() (width, height int)
	EnableMouse(...MouseFlags)
	DisableMouse()
	EnablePaste()
	DisablePaste()
	EnableFocus()
	DisableFocus()
	Colors() int
	Show()
	Sync()
	CharacterSet() string
	RegisterRuneFallback(r rune, subst string)
	UnregisterRuneFallback(r rune)
	Resize(int, int, int, int)
	Suspend() error
	Resume() error
	Beep() error
	SetSize(int, int)
	SetTitle(string)
	Tty() (Tty, bool)
	SetClipboard([]byte)
	GetClipboard()
	HasClipboard() bool
	ShowNotification(string, string)
	Terminal() (string, string)

	// Following methods are not part of the Screen api, but are used for interaction with
	// the common layer code.

	// Locker locks the underlying data structures so that we can access them
	// in a thread-safe way.
	sync.Locker

	// GetCells returns a pointer to the underlying CellBuffer that the implementation uses.
	// Various methods will write to these for performance, but will use the lock to do so.
	GetCells() *CellBuffer

	// StopQ is closed when the screen is shut down via Fini.  It remains open if the screen
	// is merely suspended.
	StopQ() <-chan struct{}

	// EventQ delivers events.  Events are posted to this by the screen in response to
	// key presses, resizes, etc.  Application code receives events from this via the
	// Screen.PollEvent, Screen.ChannelEvents APIs.
	EventQ() chan Event
}

type baseScreen struct {
	screenImpl
}

func (b *baseScreen) Put(x int, y int, str string, style Style) (remain string, width int) {
	cells := b.GetCells()
	b.Lock()
	defer b.Unlock()
	return cells.Put(x, y, str, style)
}

func (b *baseScreen) PutStrStyled(x int, y int, str string, style Style) {
	cells := b.GetCells()
	b.Lock()
	cols, rows := cells.Size()
	width := 0
	for str != "" && x < cols && y < rows {
		str, width = cells.Put(x, y, str, style)
		if width == 0 {
			break
		}
		x += width
	}
	defer b.Unlock()
}

func (b *baseScreen) PutStr(x, y int, str string) {
	b.PutStrStyled(x, y, str, StyleDefault)
}

func (b *baseScreen) Clear() {
	b.Fill(' ', StyleDefault)
}

func (b *baseScreen) Fill(r rune, style Style) {
	cb := b.GetCells()
	b.Lock()
	cb.Fill(r, style)
	b.Unlock()
}

func (b *baseScreen) SetContent(x, y int, mainc rune, combc []rune, style Style) {
	b.Put(x, y, string(append([]rune{mainc}, combc...)), style)
}

func (b *baseScreen) Get(x, y int) (string, Style, int) {
	cells := b.GetCells()
	b.Lock()
	defer b.Unlock()
	return cells.Get(x, y)
}

func (b *baseScreen) LockRegion(x, y, width, height int, lock bool) {
	cells := b.GetCells()
	b.Lock()
	for j := y; j < (y + height); j += 1 {
		for i := x; i < (x + width); i += 1 {
			switch lock {
			case true:
				cells.LockCell(i, j)
			case false:
				cells.UnlockCell(i, j)
			}
		}
	}
	b.Unlock()
}

func (b *baseScreen) SetCursorStyle(cs CursorStyle, ccs ...color.Color) {
	if len(ccs) > 0 {
		b.SetCursor(cs, ccs[0])
	} else {
		b.SetCursor(cs, ColorNone)
	}
}
