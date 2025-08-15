//

package pty

import (
	"os"
	"golang.org/x/sys/unix"
)

const (
	TIOCGWINSZ = 21608 // 'T' << 8 | 104
	TIOCSWINSZ = 21607 // 'T' << 8 | 103
)

// Winsize describes the terminal size.
type Winsize struct {
	Rows uint16 // ws_row: Number of rows (in cells)
	Cols uint16 // ws_col: Number of columns (in cells)
	X    uint16 // ws_xpixel: Width in pixels
	Y    uint16 // ws_ypixel: Height in pixels
}

// GetsizeFull returns the full terminal size description.
func GetsizeFull(t *os.File) (size *Winsize, err error) {
	var wsz *unix.Winsize
	wsz, err = unix.IoctlGetWinsize(int(t.Fd()), TIOCGWINSZ)

	if err != nil {
		return nil, err
	} else {
		return &Winsize{wsz.Row, wsz.Col, wsz.Xpixel, wsz.Ypixel}, nil
	}
}

// Get Windows Size
func Getsize(t *os.File) (rows, cols int, err error) {
	var wsz *unix.Winsize
	wsz, err = unix.IoctlGetWinsize(int(t.Fd()), TIOCGWINSZ)

	if err != nil {
		return 80, 25, err
	} else {
		return int(wsz.Row), int(wsz.Col), nil
	}
}

// Setsize resizes t to s.
func Setsize(t *os.File, ws *Winsize) error {
	wsz := unix.Winsize{ws.Rows, ws.Cols, ws.X, ws.Y}
	return unix.IoctlSetWinsize(int(t.Fd()), TIOCSWINSZ, &wsz)
}
