package termbox

import (
	"syscall"
)

// public API

// Initializes termbox library. This function should be called before any other functions.
// After successful initialization, the library must be finalized using 'Close' function.
//
// Example usage:
//      err := termbox.Init()
//      if err != nil {
//              panic(err)
//      }
//      defer termbox.Close()
func Init() error {
	var err error

	interrupt, err = create_event()
	if err != nil {
		return err
	}

	in, err = syscall.Open("CONIN$", syscall.O_RDWR, 0)
	if err != nil {
		return err
	}
	out, err = syscall.Open("CONOUT$", syscall.O_RDWR, 0)
	if err != nil {
		return err
	}

	err = get_console_mode(in, &orig_mode)
	if err != nil {
		return err
	}

	err = set_console_mode(in, enable_window_input)
	if err != nil {
		return err
	}

	orig_size = get_term_size(out)
	win_size := get_win_size(out)

	err = set_console_screen_buffer_size(out, win_size)
	if err != nil {
		return err
	}

	err = get_console_cursor_info(out, &orig_cursor_info)
	if err != nil {
		return err
	}

	show_cursor(false)
	term_size = get_term_size(out)
	back_buffer.init(int(term_size.x), int(term_size.y))
	front_buffer.init(int(term_size.x), int(term_size.y))
	back_buffer.clear()
	front_buffer.clear()
	clear()

	diffbuf = make([]diff_msg, 0, 32)

	go input_event_producer()
	IsInit = true
	return nil
}

// Finalizes termbox library, should be called after successful initialization
// when termbox's functionality isn't required anymore.
func Close() {
	// we ignore errors here, because we can't really do anything about them
	Clear(0, 0)
	Flush()

	// stop event producer
	cancel_comm <- true
	set_event(interrupt)
	select {
	case <-input_comm:
	default:
	}
	<-cancel_done_comm

	set_console_cursor_info(out, &orig_cursor_info)
	set_console_cursor_position(out, coord{})
	set_console_screen_buffer_size(out, orig_size)
	set_console_mode(in, orig_mode)
	syscall.Close(in)
	syscall.Close(out)
	syscall.Close(interrupt)
	IsInit = false
}

// Interrupt an in-progress call to PollEvent by causing it to return
// EventInterrupt.  Note that this function will block until the PollEvent
// function has successfully been interrupted.
func Interrupt() {
	interrupt_comm <- struct{}{}
}

// Synchronizes the internal back buffer with the terminal.
func Flush() error {
	update_size_maybe()
	prepare_diff_messages()
	for _, diff := range diffbuf {
		r := small_rect{
			left:   0,
			top:    diff.pos,
			right:  term_size.x - 1,
			bottom: diff.pos + diff.lines - 1,
		}
		write_console_output(out, diff.chars, r)
	}
	if !is_cursor_hidden(cursor_x, cursor_y) {
		move_cursor(cursor_x, cursor_y)
	}
	return nil
}

// Sets the position of the cursor. See also HideCursor().
func SetCursor(x, y int) {
	if is_cursor_hidden(cursor_x, cursor_y) && !is_cursor_hidden(x, y) {
		show_cursor(true)
	}

	if !is_cursor_hidden(cursor_x, cursor_y) && is_cursor_hidden(x, y) {
		show_cursor(false)
	}

	cursor_x, cursor_y = x, y
	if !is_cursor_hidden(cursor_x, cursor_y) {
		move_cursor(cursor_x, cursor_y)
	}
}

// The shortcut for SetCursor(-1, -1).
func HideCursor() {
	SetCursor(cursor_hidden, cursor_hidden)
}

// Changes cell's parameters in the internal back buffer at the specified
// position.
func SetCell(x, y int, ch rune, fg, bg Attribute) {
	if x < 0 || x >= back_buffer.width {
		return
	}
	if y < 0 || y >= back_buffer.height {
		return
	}

	back_buffer.cells[y*back_buffer.width+x] = Cell{ch, fg, bg}
}

// Returns a slice into the termbox's back buffer. You can get its dimensions
// using 'Size' function. The slice remains valid as long as no 'Clear' or
// 'Flush' function calls were made after call to this function.
func CellBuffer() []Cell {
	return back_buffer.cells
}

// Wait for an event and return it. This is a blocking function call.
func PollEvent() Event {
	select {
	case ev := <-input_comm:
		return ev
	case <-interrupt_comm:
		return Event{Type: EventInterrupt}
	}
}

// Returns the size of the internal back buffer (which is mostly the same as
// console's window size in characters). But it doesn't always match the size
// of the console window, after the console size has changed, the internal back
// buffer will get in sync only after Clear or Flush function calls.
func Size() (int, int) {
	return int(term_size.x), int(term_size.y)
}

// Clears the internal back buffer.
func Clear(fg, bg Attribute) error {
	foreground, background = fg, bg
	update_size_maybe()
	back_buffer.clear()
	return nil
}

// Sets termbox input mode. Termbox has two input modes:
//
// 1. Esc input mode. When ESC sequence is in the buffer and it doesn't match
// any known sequence. ESC means KeyEsc. This is the default input mode.
//
// 2. Alt input mode. When ESC sequence is in the buffer and it doesn't match
// any known sequence. ESC enables ModAlt modifier for the next keyboard event.
//
// Both input modes can be OR'ed with Mouse mode. Setting Mouse mode bit up will
// enable mouse button press/release and drag events.
//
// If 'mode' is InputCurrent, returns the current input mode. See also Input*
// constants.
func SetInputMode(mode InputMode) InputMode {
	if mode == InputCurrent {
		return input_mode
	}
	if mode&InputMouse != 0 {
		err := set_console_mode(in, enable_window_input|enable_mouse_input|enable_extended_flags)
		if err != nil {
			panic(err)
		}
	} else {
		err := set_console_mode(in, enable_window_input)
		if err != nil {
			panic(err)
		}
	}

	input_mode = mode
	return input_mode
}

// Sets the termbox output mode.
//
// Windows console does not support extra colour modes,
// so this will always set and return OutputNormal.
func SetOutputMode(mode OutputMode) OutputMode {
	return OutputNormal
}

// Sync comes handy when something causes desync between termbox's understanding
// of a terminal buffer and the reality. Such as a third party process. Sync
// forces a complete resync between the termbox and a terminal, it may not be
// visually pretty though. At the moment on Windows it does nothing.
func Sync() error {
	return nil
}
