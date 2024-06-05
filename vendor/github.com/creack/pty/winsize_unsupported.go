//go:build windows
// +build windows

package pty

import (
	"os"
)

// Winsize is a dummy struct to enable compilation on unsupported platforms.
type Winsize struct {
	Rows, Cols, X, Y uint16
}

// Setsize resizes t to s.
func Setsize(*os.File, *Winsize) error {
	return ErrUnsupported
}

// GetsizeFull returns the full terminal size description.
func GetsizeFull(*os.File) (*Winsize, error) {
	return nil, ErrUnsupported
}
