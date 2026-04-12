package pty

import "os"

// InheritSize applies the terminal size of pty to tty. This should be run
// in a signal handler for syscall.SIGWINCH to automatically resize the tty when
// the pty receives a window size change notification.
func InheritSize(pty, tty *os.File) error {
	size, err := GetsizeFull(pty)
	if err != nil {
		return err
	}
	return Setsize(tty, size)
}

// Getsize returns the number of rows (lines) and cols (positions
// in each line) in terminal t.
func Getsize(t *os.File) (rows, cols int, err error) {
	ws, err := GetsizeFull(t)
	if err != nil {
		return 0, 0, err
	}
	return int(ws.Rows), int(ws.Cols), nil
}
