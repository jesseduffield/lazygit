package termenv

import (
	"fmt"
	"strings"
)

// Sequence definitions.
const (
	// Cursor positioning.
	CursorUpSeq              = "%dA"
	CursorDownSeq            = "%dB"
	CursorForwardSeq         = "%dC"
	CursorBackSeq            = "%dD"
	CursorNextLineSeq        = "%dE"
	CursorPreviousLineSeq    = "%dF"
	CursorHorizontalSeq      = "%dG"
	CursorPositionSeq        = "%d;%dH"
	EraseDisplaySeq          = "%dJ"
	EraseLineSeq             = "%dK"
	ScrollUpSeq              = "%dS"
	ScrollDownSeq            = "%dT"
	SaveCursorPositionSeq    = "s"
	RestoreCursorPositionSeq = "u"
	ChangeScrollingRegionSeq = "%d;%dr"
	InsertLineSeq            = "%dL"
	DeleteLineSeq            = "%dM"

	// Explicit values for EraseLineSeq.
	EraseLineRightSeq  = "0K"
	EraseLineLeftSeq   = "1K"
	EraseEntireLineSeq = "2K"

	// Mouse.
	EnableMousePressSeq       = "?9h" // press only (X10)
	DisableMousePressSeq      = "?9l"
	EnableMouseSeq            = "?1000h" // press, release, wheel
	DisableMouseSeq           = "?1000l"
	EnableMouseHiliteSeq      = "?1001h" // highlight
	DisableMouseHiliteSeq     = "?1001l"
	EnableMouseCellMotionSeq  = "?1002h" // press, release, move on pressed, wheel
	DisableMouseCellMotionSeq = "?1002l"
	EnableMouseAllMotionSeq   = "?1003h" // press, release, move, wheel
	DisableMouseAllMotionSeq  = "?1003l"

	// Screen.
	RestoreScreenSeq = "?47l"
	SaveScreenSeq    = "?47h"
	AltScreenSeq     = "?1049h"
	ExitAltScreenSeq = "?1049l"

	// Bracketed paste.
	// https://en.wikipedia.org/wiki/Bracketed-paste
	EnableBracketedPasteSeq  = "?2004h"
	DisableBracketedPasteSeq = "?2004l"
	StartBracketedPasteSeq   = "200~"
	EndBracketedPasteSeq     = "201~"

	// Session.
	SetWindowTitleSeq     = "2;%s\007"
	SetForegroundColorSeq = "10;%s\007"
	SetBackgroundColorSeq = "11;%s\007"
	SetCursorColorSeq     = "12;%s\007"
	ShowCursorSeq         = "?25h"
	HideCursorSeq         = "?25l"
)

// Reset the terminal to its default style, removing any active styles.
func (o Output) Reset() {
	fmt.Fprint(o.tty, CSI+ResetSeq+"m")
}

// SetForegroundColor sets the default foreground color.
func (o Output) SetForegroundColor(color Color) {
	fmt.Fprintf(o.tty, OSC+SetForegroundColorSeq, color)
}

// SetBackgroundColor sets the default background color.
func (o Output) SetBackgroundColor(color Color) {
	fmt.Fprintf(o.tty, OSC+SetBackgroundColorSeq, color)
}

// SetCursorColor sets the cursor color.
func (o Output) SetCursorColor(color Color) {
	fmt.Fprintf(o.tty, OSC+SetCursorColorSeq, color)
}

// RestoreScreen restores a previously saved screen state.
func (o Output) RestoreScreen() {
	fmt.Fprint(o.tty, CSI+RestoreScreenSeq)
}

// SaveScreen saves the screen state.
func (o Output) SaveScreen() {
	fmt.Fprint(o.tty, CSI+SaveScreenSeq)
}

// AltScreen switches to the alternate screen buffer. The former view can be
// restored with ExitAltScreen().
func (o Output) AltScreen() {
	fmt.Fprint(o.tty, CSI+AltScreenSeq)
}

// ExitAltScreen exits the alternate screen buffer and returns to the former
// terminal view.
func (o Output) ExitAltScreen() {
	fmt.Fprint(o.tty, CSI+ExitAltScreenSeq)
}

// ClearScreen clears the visible portion of the terminal.
func (o Output) ClearScreen() {
	fmt.Fprintf(o.tty, CSI+EraseDisplaySeq, 2)
	o.MoveCursor(1, 1)
}

// MoveCursor moves the cursor to a given position.
func (o Output) MoveCursor(row int, column int) {
	fmt.Fprintf(o.tty, CSI+CursorPositionSeq, row, column)
}

// HideCursor hides the cursor.
func (o Output) HideCursor() {
	fmt.Fprint(o.tty, CSI+HideCursorSeq)
}

// ShowCursor shows the cursor.
func (o Output) ShowCursor() {
	fmt.Fprint(o.tty, CSI+ShowCursorSeq)
}

// SaveCursorPosition saves the cursor position.
func (o Output) SaveCursorPosition() {
	fmt.Fprint(o.tty, CSI+SaveCursorPositionSeq)
}

// RestoreCursorPosition restores a saved cursor position.
func (o Output) RestoreCursorPosition() {
	fmt.Fprint(o.tty, CSI+RestoreCursorPositionSeq)
}

// CursorUp moves the cursor up a given number of lines.
func (o Output) CursorUp(n int) {
	fmt.Fprintf(o.tty, CSI+CursorUpSeq, n)
}

// CursorDown moves the cursor down a given number of lines.
func (o Output) CursorDown(n int) {
	fmt.Fprintf(o.tty, CSI+CursorDownSeq, n)
}

// CursorForward moves the cursor up a given number of lines.
func (o Output) CursorForward(n int) {
	fmt.Fprintf(o.tty, CSI+CursorForwardSeq, n)
}

// CursorBack moves the cursor backwards a given number of cells.
func (o Output) CursorBack(n int) {
	fmt.Fprintf(o.tty, CSI+CursorBackSeq, n)
}

// CursorNextLine moves the cursor down a given number of lines and places it at
// the beginning of the line.
func (o Output) CursorNextLine(n int) {
	fmt.Fprintf(o.tty, CSI+CursorNextLineSeq, n)
}

// CursorPrevLine moves the cursor up a given number of lines and places it at
// the beginning of the line.
func (o Output) CursorPrevLine(n int) {
	fmt.Fprintf(o.tty, CSI+CursorPreviousLineSeq, n)
}

// ClearLine clears the current line.
func (o Output) ClearLine() {
	fmt.Fprint(o.tty, CSI+EraseEntireLineSeq)
}

// ClearLineLeft clears the line to the left of the cursor.
func (o Output) ClearLineLeft() {
	fmt.Fprint(o.tty, CSI+EraseLineLeftSeq)
}

// ClearLineRight clears the line to the right of the cursor.
func (o Output) ClearLineRight() {
	fmt.Fprint(o.tty, CSI+EraseLineRightSeq)
}

// ClearLines clears a given number of lines.
func (o Output) ClearLines(n int) {
	clearLine := fmt.Sprintf(CSI+EraseLineSeq, 2)
	cursorUp := fmt.Sprintf(CSI+CursorUpSeq, 1)
	fmt.Fprint(o.tty, clearLine+strings.Repeat(cursorUp+clearLine, n))
}

// ChangeScrollingRegion sets the scrolling region of the terminal.
func (o Output) ChangeScrollingRegion(top, bottom int) {
	fmt.Fprintf(o.tty, CSI+ChangeScrollingRegionSeq, top, bottom)
}

// InsertLines inserts the given number of lines at the top of the scrollable
// region, pushing lines below down.
func (o Output) InsertLines(n int) {
	fmt.Fprintf(o.tty, CSI+InsertLineSeq, n)
}

// DeleteLines deletes the given number of lines, pulling any lines in
// the scrollable region below up.
func (o Output) DeleteLines(n int) {
	fmt.Fprintf(o.tty, CSI+DeleteLineSeq, n)
}

// EnableMousePress enables X10 mouse mode. Button press events are sent only.
func (o Output) EnableMousePress() {
	fmt.Fprint(o.tty, CSI+EnableMousePressSeq)
}

// DisableMousePress disables X10 mouse mode.
func (o Output) DisableMousePress() {
	fmt.Fprint(o.tty, CSI+DisableMousePressSeq)
}

// EnableMouse enables Mouse Tracking mode.
func (o Output) EnableMouse() {
	fmt.Fprint(o.tty, CSI+EnableMouseSeq)
}

// DisableMouse disables Mouse Tracking mode.
func (o Output) DisableMouse() {
	fmt.Fprint(o.tty, CSI+DisableMouseSeq)
}

// EnableMouseHilite enables Hilite Mouse Tracking mode.
func (o Output) EnableMouseHilite() {
	fmt.Fprint(o.tty, CSI+EnableMouseHiliteSeq)
}

// DisableMouseHilite disables Hilite Mouse Tracking mode.
func (o Output) DisableMouseHilite() {
	fmt.Fprint(o.tty, CSI+DisableMouseHiliteSeq)
}

// EnableMouseCellMotion enables Cell Motion Mouse Tracking mode.
func (o Output) EnableMouseCellMotion() {
	fmt.Fprint(o.tty, CSI+EnableMouseCellMotionSeq)
}

// DisableMouseCellMotion disables Cell Motion Mouse Tracking mode.
func (o Output) DisableMouseCellMotion() {
	fmt.Fprint(o.tty, CSI+DisableMouseCellMotionSeq)
}

// EnableMouseAllMotion enables All Motion Mouse mode.
func (o Output) EnableMouseAllMotion() {
	fmt.Fprint(o.tty, CSI+EnableMouseAllMotionSeq)
}

// DisableMouseAllMotion disables All Motion Mouse mode.
func (o Output) DisableMouseAllMotion() {
	fmt.Fprint(o.tty, CSI+DisableMouseAllMotionSeq)
}

// SetWindowTitle sets the terminal window title.
func (o Output) SetWindowTitle(title string) {
	fmt.Fprintf(o.tty, OSC+SetWindowTitleSeq, title)
}

// EnableBracketedPaste enables bracketed paste.
func (o Output) EnableBracketedPaste() {
	fmt.Fprintf(o.tty, CSI+EnableBracketedPasteSeq)
}

// DisableBracketedPaste disables bracketed paste.
func (o Output) DisableBracketedPaste() {
	fmt.Fprintf(o.tty, CSI+DisableBracketedPasteSeq)
}

// Legacy functions.

// Reset the terminal to its default style, removing any active styles.
//
// Deprecated: please use termenv.Output instead.
func Reset() {
	output.Reset()
}

// SetForegroundColor sets the default foreground color.
//
// Deprecated: please use termenv.Output instead.
func SetForegroundColor(color Color) {
	output.SetForegroundColor(color)
}

// SetBackgroundColor sets the default background color.
//
// Deprecated: please use termenv.Output instead.
func SetBackgroundColor(color Color) {
	output.SetBackgroundColor(color)
}

// SetCursorColor sets the cursor color.
//
// Deprecated: please use termenv.Output instead.
func SetCursorColor(color Color) {
	output.SetCursorColor(color)
}

// RestoreScreen restores a previously saved screen state.
//
// Deprecated: please use termenv.Output instead.
func RestoreScreen() {
	output.RestoreScreen()
}

// SaveScreen saves the screen state.
//
// Deprecated: please use termenv.Output instead.
func SaveScreen() {
	output.SaveScreen()
}

// AltScreen switches to the alternate screen buffer. The former view can be
// restored with ExitAltScreen().
//
// Deprecated: please use termenv.Output instead.
func AltScreen() {
	output.AltScreen()
}

// ExitAltScreen exits the alternate screen buffer and returns to the former
// terminal view.
//
// Deprecated: please use termenv.Output instead.
func ExitAltScreen() {
	output.ExitAltScreen()
}

// ClearScreen clears the visible portion of the terminal.
//
// Deprecated: please use termenv.Output instead.
func ClearScreen() {
	output.ClearScreen()
}

// MoveCursor moves the cursor to a given position.
//
// Deprecated: please use termenv.Output instead.
func MoveCursor(row int, column int) {
	output.MoveCursor(row, column)
}

// HideCursor hides the cursor.
//
// Deprecated: please use termenv.Output instead.
func HideCursor() {
	output.HideCursor()
}

// ShowCursor shows the cursor.
//
// Deprecated: please use termenv.Output instead.
func ShowCursor() {
	output.ShowCursor()
}

// SaveCursorPosition saves the cursor position.
//
// Deprecated: please use termenv.Output instead.
func SaveCursorPosition() {
	output.SaveCursorPosition()
}

// RestoreCursorPosition restores a saved cursor position.
//
// Deprecated: please use termenv.Output instead.
func RestoreCursorPosition() {
	output.RestoreCursorPosition()
}

// CursorUp moves the cursor up a given number of lines.
//
// Deprecated: please use termenv.Output instead.
func CursorUp(n int) {
	output.CursorUp(n)
}

// CursorDown moves the cursor down a given number of lines.
//
// Deprecated: please use termenv.Output instead.
func CursorDown(n int) {
	output.CursorDown(n)
}

// CursorForward moves the cursor up a given number of lines.
//
// Deprecated: please use termenv.Output instead.
func CursorForward(n int) {
	output.CursorForward(n)
}

// CursorBack moves the cursor backwards a given number of cells.
//
// Deprecated: please use termenv.Output instead.
func CursorBack(n int) {
	output.CursorBack(n)
}

// CursorNextLine moves the cursor down a given number of lines and places it at
// the beginning of the line.
//
// Deprecated: please use termenv.Output instead.
func CursorNextLine(n int) {
	output.CursorNextLine(n)
}

// CursorPrevLine moves the cursor up a given number of lines and places it at
// the beginning of the line.
//
// Deprecated: please use termenv.Output instead.
func CursorPrevLine(n int) {
	output.CursorPrevLine(n)
}

// ClearLine clears the current line.
//
// Deprecated: please use termenv.Output instead.
func ClearLine() {
	output.ClearLine()
}

// ClearLineLeft clears the line to the left of the cursor.
//
// Deprecated: please use termenv.Output instead.
func ClearLineLeft() {
	output.ClearLineLeft()
}

// ClearLineRight clears the line to the right of the cursor.
//
// Deprecated: please use termenv.Output instead.
func ClearLineRight() {
	output.ClearLineRight()
}

// ClearLines clears a given number of lines.
//
// Deprecated: please use termenv.Output instead.
func ClearLines(n int) {
	output.ClearLines(n)
}

// ChangeScrollingRegion sets the scrolling region of the terminal.
//
// Deprecated: please use termenv.Output instead.
func ChangeScrollingRegion(top, bottom int) {
	output.ChangeScrollingRegion(top, bottom)
}

// InsertLines inserts the given number of lines at the top of the scrollable
// region, pushing lines below down.
//
// Deprecated: please use termenv.Output instead.
func InsertLines(n int) {
	output.InsertLines(n)
}

// DeleteLines deletes the given number of lines, pulling any lines in
// the scrollable region below up.
//
// Deprecated: please use termenv.Output instead.
func DeleteLines(n int) {
	output.DeleteLines(n)
}

// EnableMousePress enables X10 mouse mode. Button press events are sent only.
//
// Deprecated: please use termenv.Output instead.
func EnableMousePress() {
	output.EnableMousePress()
}

// DisableMousePress disables X10 mouse mode.
//
// Deprecated: please use termenv.Output instead.
func DisableMousePress() {
	output.DisableMousePress()
}

// EnableMouse enables Mouse Tracking mode.
//
// Deprecated: please use termenv.Output instead.
func EnableMouse() {
	output.EnableMouse()
}

// DisableMouse disables Mouse Tracking mode.
//
// Deprecated: please use termenv.Output instead.
func DisableMouse() {
	output.DisableMouse()
}

// EnableMouseHilite enables Hilite Mouse Tracking mode.
//
// Deprecated: please use termenv.Output instead.
func EnableMouseHilite() {
	output.EnableMouseHilite()
}

// DisableMouseHilite disables Hilite Mouse Tracking mode.
//
// Deprecated: please use termenv.Output instead.
func DisableMouseHilite() {
	output.DisableMouseHilite()
}

// EnableMouseCellMotion enables Cell Motion Mouse Tracking mode.
//
// Deprecated: please use termenv.Output instead.
func EnableMouseCellMotion() {
	output.EnableMouseCellMotion()
}

// DisableMouseCellMotion disables Cell Motion Mouse Tracking mode.
//
// Deprecated: please use termenv.Output instead.
func DisableMouseCellMotion() {
	output.DisableMouseCellMotion()
}

// EnableMouseAllMotion enables All Motion Mouse mode.
//
// Deprecated: please use termenv.Output instead.
func EnableMouseAllMotion() {
	output.EnableMouseAllMotion()
}

// DisableMouseAllMotion disables All Motion Mouse mode.
//
// Deprecated: please use termenv.Output instead.
func DisableMouseAllMotion() {
	output.DisableMouseAllMotion()
}

// SetWindowTitle sets the terminal window title.
//
// Deprecated: please use termenv.Output instead.
func SetWindowTitle(title string) {
	output.SetWindowTitle(title)
}

// EnableBracketedPaste enables bracketed paste.
//
// Deprecated: please use termenv.Output instead.
func EnableBracketedPaste() {
	output.EnableBracketedPaste()
}

// DisableBracketedPaste disables bracketed paste.
//
// Deprecated: please use termenv.Output instead.
func DisableBracketedPaste() {
	output.DisableBracketedPaste()
}
