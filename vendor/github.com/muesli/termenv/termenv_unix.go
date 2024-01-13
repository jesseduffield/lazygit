//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris
// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package termenv

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sys/unix"
)

const (
	// timeout for OSC queries
	OSCTimeout = 5 * time.Second
)

// ColorProfile returns the supported color profile:
// Ascii, ANSI, ANSI256, or TrueColor.
func (o *Output) ColorProfile() Profile {
	if !o.isTTY() {
		return Ascii
	}

	term := o.environ.Getenv("TERM")
	colorTerm := o.environ.Getenv("COLORTERM")

	switch strings.ToLower(colorTerm) {
	case "24bit":
		fallthrough
	case "truecolor":
		if strings.HasPrefix(term, "screen") {
			// tmux supports TrueColor, screen only ANSI256
			if o.environ.Getenv("TERM_PROGRAM") != "tmux" {
				return ANSI256
			}
		}
		return TrueColor
	case "yes":
		fallthrough
	case "true":
		return ANSI256
	}

	switch term {
	case "xterm-kitty":
		return TrueColor
	case "linux":
		return ANSI
	}

	if strings.Contains(term, "256color") {
		return ANSI256
	}
	if strings.Contains(term, "color") {
		return ANSI
	}
	if strings.Contains(term, "ansi") {
		return ANSI
	}

	return Ascii
}

func (o Output) foregroundColor() Color {
	s, err := o.termStatusReport(10)
	if err == nil {
		c, err := xTermColor(s)
		if err == nil {
			return c
		}
	}

	colorFGBG := o.environ.Getenv("COLORFGBG")
	if strings.Contains(colorFGBG, ";") {
		c := strings.Split(colorFGBG, ";")
		i, err := strconv.Atoi(c[0])
		if err == nil {
			return ANSIColor(i)
		}
	}

	// default gray
	return ANSIColor(7)
}

func (o Output) backgroundColor() Color {
	s, err := o.termStatusReport(11)
	if err == nil {
		c, err := xTermColor(s)
		if err == nil {
			return c
		}
	}

	colorFGBG := o.environ.Getenv("COLORFGBG")
	if strings.Contains(colorFGBG, ";") {
		c := strings.Split(colorFGBG, ";")
		i, err := strconv.Atoi(c[len(c)-1])
		if err == nil {
			return ANSIColor(i)
		}
	}

	// default black
	return ANSIColor(0)
}

func waitForData(fd uintptr, timeout time.Duration) error {
	tv := unix.NsecToTimeval(int64(timeout))
	var readfds unix.FdSet
	readfds.Set(int(fd))

	for {
		n, err := unix.Select(int(fd)+1, &readfds, nil, nil, &tv)
		if err == unix.EINTR {
			continue
		}
		if err != nil {
			return err
		}
		if n == 0 {
			return fmt.Errorf("timeout")
		}

		break
	}

	return nil
}

func readNextByte(f File) (byte, error) {
	if err := waitForData(f.Fd(), OSCTimeout); err != nil {
		return 0, err
	}

	var b [1]byte
	n, err := f.Read(b[:])
	if err != nil {
		return 0, err
	}

	if n == 0 {
		panic("read returned no data")
	}

	return b[0], nil
}

// readNextResponse reads either an OSC response or a cursor position response:
//   - OSC response: "\x1b]11;rgb:1111/1111/1111\x1b\\"
//   - cursor position response: "\x1b[42;1R"
func readNextResponse(fd File) (response string, isOSC bool, err error) {
	start, err := readNextByte(fd)
	if err != nil {
		return "", false, err
	}

	// first byte must be ESC
	for start != '\033' {
		start, err = readNextByte(fd)
		if err != nil {
			return "", false, err
		}
	}

	response += string(start)

	// next byte is either '[' (cursor position response) or ']' (OSC response)
	tpe, err := readNextByte(fd)
	if err != nil {
		return "", false, err
	}

	response += string(tpe)

	var oscResponse bool
	switch tpe {
	case '[':
		oscResponse = false
	case ']':
		oscResponse = true
	default:
		return "", false, ErrStatusReport
	}

	for {
		b, err := readNextByte(fd)
		if err != nil {
			return "", false, err
		}

		response += string(b)

		if oscResponse {
			// OSC can be terminated by BEL (\a) or ST (ESC)
			if b == '\a' || strings.HasSuffix(response, "\033") {
				return response, true, nil
			}
		} else {
			// cursor position response is terminated by 'R'
			if b == 'R' {
				return response, false, nil
			}
		}

		// both responses have less than 25 bytes, so if we read more, that's an error
		if len(response) > 25 {
			break
		}
	}

	return "", false, ErrStatusReport
}

func (o Output) termStatusReport(sequence int) (string, error) {
	// screen/tmux can't support OSC, because they can be connected to multiple
	// terminals concurrently.
	term := o.environ.Getenv("TERM")
	if strings.HasPrefix(term, "screen") || strings.HasPrefix(term, "tmux") {
		return "", ErrStatusReport
	}

	tty := o.TTY()
	if tty == nil {
		return "", ErrStatusReport
	}

	fd := int(tty.Fd())
	// if in background, we can't control the terminal
	if !isForeground(fd) {
		return "", ErrStatusReport
	}

	t, err := unix.IoctlGetTermios(fd, tcgetattr)
	if err != nil {
		return "", fmt.Errorf("%s: %s", ErrStatusReport, err)
	}
	defer unix.IoctlSetTermios(fd, tcsetattr, t) //nolint:errcheck

	noecho := *t
	noecho.Lflag = noecho.Lflag &^ unix.ECHO
	noecho.Lflag = noecho.Lflag &^ unix.ICANON
	if err := unix.IoctlSetTermios(fd, tcsetattr, &noecho); err != nil {
		return "", fmt.Errorf("%s: %s", ErrStatusReport, err)
	}

	// first, send OSC query, which is ignored by terminal which do not support it
	fmt.Fprintf(tty, "\033]%d;?\033\\", sequence)

	// then, query cursor position, should be supported by all terminals
	fmt.Fprintf(tty, "\033[6n")

	// read the next response
	res, isOSC, err := readNextResponse(tty)
	if err != nil {
		return "", fmt.Errorf("%s: %s", ErrStatusReport, err)
	}

	// if this is not OSC response, then the terminal does not support it
	if !isOSC {
		return "", ErrStatusReport
	}

	// read the cursor query response next and discard the result
	_, _, err = readNextResponse(tty)
	if err != nil {
		return "", err
	}

	// fmt.Println("Rcvd", res[1:])
	return res, nil
}

// EnableVirtualTerminalProcessing enables virtual terminal processing on
// Windows for w and returns a function that restores w to its previous state.
// On non-Windows platforms, or if w does not refer to a terminal, then it
// returns a non-nil no-op function and no error.
func EnableVirtualTerminalProcessing(w io.Writer) (func() error, error) {
	return func() error { return nil }, nil
}
