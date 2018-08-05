// +build !windows

package termbox

import "unicode/utf8"
import "bytes"
import "syscall"
import "unsafe"
import "strings"
import "strconv"
import "os"
import "io"

// private API

const (
	t_enter_ca = iota
	t_exit_ca
	t_show_cursor
	t_hide_cursor
	t_clear_screen
	t_sgr0
	t_underline
	t_bold
	t_blink
	t_reverse
	t_enter_keypad
	t_exit_keypad
	t_enter_mouse
	t_exit_mouse
	t_max_funcs
)

const (
	coord_invalid = -2
	attr_invalid  = Attribute(0xFFFF)
)

type input_event struct {
	data []byte
	err  error
}

type extract_event_res int

const (
	event_not_extracted extract_event_res = iota
	event_extracted
	esc_wait
)

var (
	// term specific sequences
	keys  []string
	funcs []string

	// termbox inner state
	orig_tios      syscall_Termios
	back_buffer    cellbuf
	front_buffer   cellbuf
	termw          int
	termh          int
	input_mode     = InputEsc
	output_mode    = OutputNormal
	out            *os.File
	in             int
	lastfg         = attr_invalid
	lastbg         = attr_invalid
	lastx          = coord_invalid
	lasty          = coord_invalid
	cursor_x       = cursor_hidden
	cursor_y       = cursor_hidden
	foreground     = ColorDefault
	background     = ColorDefault
	inbuf          = make([]byte, 0, 64)
	outbuf         bytes.Buffer
	sigwinch       = make(chan os.Signal, 1)
	sigio          = make(chan os.Signal, 1)
	quit           = make(chan int)
	input_comm     = make(chan input_event)
	interrupt_comm = make(chan struct{})
	intbuf         = make([]byte, 0, 16)

	// grayscale indexes
	grayscale = []Attribute{
		0, 17, 233, 234, 235, 236, 237, 238, 239, 240, 241, 242, 243, 244,
		245, 246, 247, 248, 249, 250, 251, 252, 253, 254, 255, 256, 232,
	}
)

func write_cursor(x, y int) {
	outbuf.WriteString("\033[")
	outbuf.Write(strconv.AppendUint(intbuf, uint64(y+1), 10))
	outbuf.WriteString(";")
	outbuf.Write(strconv.AppendUint(intbuf, uint64(x+1), 10))
	outbuf.WriteString("H")
}

func write_sgr_fg(a Attribute) {
	switch output_mode {
	case Output256, Output216, OutputGrayscale:
		outbuf.WriteString("\033[38;5;")
		outbuf.Write(strconv.AppendUint(intbuf, uint64(a-1), 10))
		outbuf.WriteString("m")
	default:
		outbuf.WriteString("\033[3")
		outbuf.Write(strconv.AppendUint(intbuf, uint64(a-1), 10))
		outbuf.WriteString("m")
	}
}

func write_sgr_bg(a Attribute) {
	switch output_mode {
	case Output256, Output216, OutputGrayscale:
		outbuf.WriteString("\033[48;5;")
		outbuf.Write(strconv.AppendUint(intbuf, uint64(a-1), 10))
		outbuf.WriteString("m")
	default:
		outbuf.WriteString("\033[4")
		outbuf.Write(strconv.AppendUint(intbuf, uint64(a-1), 10))
		outbuf.WriteString("m")
	}
}

func write_sgr(fg, bg Attribute) {
	switch output_mode {
	case Output256, Output216, OutputGrayscale:
		outbuf.WriteString("\033[38;5;")
		outbuf.Write(strconv.AppendUint(intbuf, uint64(fg-1), 10))
		outbuf.WriteString("m")
		outbuf.WriteString("\033[48;5;")
		outbuf.Write(strconv.AppendUint(intbuf, uint64(bg-1), 10))
		outbuf.WriteString("m")
	default:
		outbuf.WriteString("\033[3")
		outbuf.Write(strconv.AppendUint(intbuf, uint64(fg-1), 10))
		outbuf.WriteString(";4")
		outbuf.Write(strconv.AppendUint(intbuf, uint64(bg-1), 10))
		outbuf.WriteString("m")
	}
}

type winsize struct {
	rows    uint16
	cols    uint16
	xpixels uint16
	ypixels uint16
}

func get_term_size(fd uintptr) (int, int) {
	var sz winsize
	_, _, _ = syscall.Syscall(syscall.SYS_IOCTL,
		fd, uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&sz)))
	return int(sz.cols), int(sz.rows)
}

func send_attr(fg, bg Attribute) {
	if fg == lastfg && bg == lastbg {
		return
	}

	outbuf.WriteString(funcs[t_sgr0])

	var fgcol, bgcol Attribute

	switch output_mode {
	case Output256:
		fgcol = fg & 0x1FF
		bgcol = bg & 0x1FF
	case Output216:
		fgcol = fg & 0xFF
		bgcol = bg & 0xFF
		if fgcol > 216 {
			fgcol = ColorDefault
		}
		if bgcol > 216 {
			bgcol = ColorDefault
		}
		if fgcol != ColorDefault {
			fgcol += 0x10
		}
		if bgcol != ColorDefault {
			bgcol += 0x10
		}
	case OutputGrayscale:
		fgcol = fg & 0x1F
		bgcol = bg & 0x1F
		if fgcol > 26 {
			fgcol = ColorDefault
		}
		if bgcol > 26 {
			bgcol = ColorDefault
		}
		if fgcol != ColorDefault {
			fgcol = grayscale[fgcol]
		}
		if bgcol != ColorDefault {
			bgcol = grayscale[bgcol]
		}
	default:
		fgcol = fg & 0x0F
		bgcol = bg & 0x0F
	}

	if fgcol != ColorDefault {
		if bgcol != ColorDefault {
			write_sgr(fgcol, bgcol)
		} else {
			write_sgr_fg(fgcol)
		}
	} else if bgcol != ColorDefault {
		write_sgr_bg(bgcol)
	}

	if fg&AttrBold != 0 {
		outbuf.WriteString(funcs[t_bold])
	}
	if bg&AttrBold != 0 {
		outbuf.WriteString(funcs[t_blink])
	}
	if fg&AttrUnderline != 0 {
		outbuf.WriteString(funcs[t_underline])
	}
	if fg&AttrReverse|bg&AttrReverse != 0 {
		outbuf.WriteString(funcs[t_reverse])
	}

	lastfg, lastbg = fg, bg
}

func send_char(x, y int, ch rune) {
	var buf [8]byte
	n := utf8.EncodeRune(buf[:], ch)
	if x-1 != lastx || y != lasty {
		write_cursor(x, y)
	}
	lastx, lasty = x, y
	outbuf.Write(buf[:n])
}

func flush() error {
	_, err := io.Copy(out, &outbuf)
	outbuf.Reset()
	return err
}

func send_clear() error {
	send_attr(foreground, background)
	outbuf.WriteString(funcs[t_clear_screen])
	if !is_cursor_hidden(cursor_x, cursor_y) {
		write_cursor(cursor_x, cursor_y)
	}

	// we need to invalidate cursor position too and these two vars are
	// used only for simple cursor positioning optimization, cursor
	// actually may be in the correct place, but we simply discard
	// optimization once and it gives us simple solution for the case when
	// cursor moved
	lastx = coord_invalid
	lasty = coord_invalid

	return flush()
}

func update_size_maybe() error {
	w, h := get_term_size(out.Fd())
	if w != termw || h != termh {
		termw, termh = w, h
		back_buffer.resize(termw, termh)
		front_buffer.resize(termw, termh)
		front_buffer.clear()
		return send_clear()
	}
	return nil
}

func tcsetattr(fd uintptr, termios *syscall_Termios) error {
	r, _, e := syscall.Syscall(syscall.SYS_IOCTL,
		fd, uintptr(syscall_TCSETS), uintptr(unsafe.Pointer(termios)))
	if r != 0 {
		return os.NewSyscallError("SYS_IOCTL", e)
	}
	return nil
}

func tcgetattr(fd uintptr, termios *syscall_Termios) error {
	r, _, e := syscall.Syscall(syscall.SYS_IOCTL,
		fd, uintptr(syscall_TCGETS), uintptr(unsafe.Pointer(termios)))
	if r != 0 {
		return os.NewSyscallError("SYS_IOCTL", e)
	}
	return nil
}

func parse_mouse_event(event *Event, buf string) (int, bool) {
	if strings.HasPrefix(buf, "\033[M") && len(buf) >= 6 {
		// X10 mouse encoding, the simplest one
		// \033 [ M Cb Cx Cy
		b := buf[3] - 32
		switch b & 3 {
		case 0:
			if b&64 != 0 {
				event.Key = MouseWheelUp
			} else {
				event.Key = MouseLeft
			}
		case 1:
			if b&64 != 0 {
				event.Key = MouseWheelDown
			} else {
				event.Key = MouseMiddle
			}
		case 2:
			event.Key = MouseRight
		case 3:
			event.Key = MouseRelease
		default:
			return 6, false
		}
		event.Type = EventMouse // KeyEvent by default
		if b&32 != 0 {
			event.Mod |= ModMotion
		}

		// the coord is 1,1 for upper left
		event.MouseX = int(buf[4]) - 1 - 32
		event.MouseY = int(buf[5]) - 1 - 32
		return 6, true
	} else if strings.HasPrefix(buf, "\033[<") || strings.HasPrefix(buf, "\033[") {
		// xterm 1006 extended mode or urxvt 1015 extended mode
		// xterm: \033 [ < Cb ; Cx ; Cy (M or m)
		// urxvt: \033 [ Cb ; Cx ; Cy M

		// find the first M or m, that's where we stop
		mi := strings.IndexAny(buf, "Mm")
		if mi == -1 {
			return 0, false
		}

		// whether it's a capital M or not
		isM := buf[mi] == 'M'

		// whether it's urxvt or not
		isU := false

		// buf[2] is safe here, because having M or m found means we have at
		// least 3 bytes in a string
		if buf[2] == '<' {
			buf = buf[3:mi]
		} else {
			isU = true
			buf = buf[2:mi]
		}

		s1 := strings.Index(buf, ";")
		s2 := strings.LastIndex(buf, ";")
		// not found or only one ';'
		if s1 == -1 || s2 == -1 || s1 == s2 {
			return 0, false
		}

		n1, err := strconv.ParseInt(buf[0:s1], 10, 64)
		if err != nil {
			return 0, false
		}
		n2, err := strconv.ParseInt(buf[s1+1:s2], 10, 64)
		if err != nil {
			return 0, false
		}
		n3, err := strconv.ParseInt(buf[s2+1:], 10, 64)
		if err != nil {
			return 0, false
		}

		// on urxvt, first number is encoded exactly as in X10, but we need to
		// make it zero-based, on xterm it is zero-based already
		if isU {
			n1 -= 32
		}
		switch n1 & 3 {
		case 0:
			if n1&64 != 0 {
				event.Key = MouseWheelUp
			} else {
				event.Key = MouseLeft
			}
		case 1:
			if n1&64 != 0 {
				event.Key = MouseWheelDown
			} else {
				event.Key = MouseMiddle
			}
		case 2:
			event.Key = MouseRight
		case 3:
			event.Key = MouseRelease
		default:
			return mi + 1, false
		}
		if !isM {
			// on xterm mouse release is signaled by lowercase m
			event.Key = MouseRelease
		}

		event.Type = EventMouse // KeyEvent by default
		if n1&32 != 0 {
			event.Mod |= ModMotion
		}

		event.MouseX = int(n2) - 1
		event.MouseY = int(n3) - 1
		return mi + 1, true
	}

	return 0, false
}

func parse_escape_sequence(event *Event, buf []byte) (int, bool) {
	bufstr := string(buf)
	for i, key := range keys {
		if strings.HasPrefix(bufstr, key) {
			event.Ch = 0
			event.Key = Key(0xFFFF - i)
			return len(key), true
		}
	}

	// if none of the keys match, let's try mouse sequences
	return parse_mouse_event(event, bufstr)
}

func extract_raw_event(data []byte, event *Event) bool {
	if len(inbuf) == 0 {
		return false
	}

	n := len(data)
	if n == 0 {
		return false
	}

	n = copy(data, inbuf)
	copy(inbuf, inbuf[n:])
	inbuf = inbuf[:len(inbuf)-n]

	event.N = n
	event.Type = EventRaw
	return true
}

func extract_event(inbuf []byte, event *Event, allow_esc_wait bool) extract_event_res {
	if len(inbuf) == 0 {
		event.N = 0
		return event_not_extracted
	}

	if inbuf[0] == '\033' {
		// possible escape sequence
		if n, ok := parse_escape_sequence(event, inbuf); n != 0 {
			event.N = n
			if ok {
				return event_extracted
			} else {
				return event_not_extracted
			}
		}

		// possible partially read escape sequence; trigger a wait if appropriate
		if enable_wait_for_escape_sequence() && allow_esc_wait {
			event.N = 0
			return esc_wait
		}

		// it's not escape sequence, then it's Alt or Esc, check input_mode
		switch {
		case input_mode&InputEsc != 0:
			// if we're in escape mode, fill Esc event, pop buffer, return success
			event.Ch = 0
			event.Key = KeyEsc
			event.Mod = 0
			event.N = 1
			return event_extracted
		case input_mode&InputAlt != 0:
			// if we're in alt mode, set Alt modifier to event and redo parsing
			event.Mod = ModAlt
			status := extract_event(inbuf[1:], event, false)
			if status == event_extracted {
				event.N++
			} else {
				event.N = 0
			}
			return status
		default:
			panic("unreachable")
		}
	}

	// if we're here, this is not an escape sequence and not an alt sequence
	// so, it's a FUNCTIONAL KEY or a UNICODE character

	// first of all check if it's a functional key
	if Key(inbuf[0]) <= KeySpace || Key(inbuf[0]) == KeyBackspace2 {
		// fill event, pop buffer, return success
		event.Ch = 0
		event.Key = Key(inbuf[0])
		event.N = 1
		return event_extracted
	}

	// the only possible option is utf8 rune
	if r, n := utf8.DecodeRune(inbuf); r != utf8.RuneError {
		event.Ch = r
		event.Key = 0
		event.N = n
		return event_extracted
	}

	return event_not_extracted
}

func fcntl(fd int, cmd int, arg int) (val int, err error) {
	r, _, e := syscall.Syscall(syscall.SYS_FCNTL, uintptr(fd), uintptr(cmd),
		uintptr(arg))
	val = int(r)
	if e != 0 {
		err = e
	}
	return
}
