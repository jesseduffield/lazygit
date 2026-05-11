//go:build windows
// +build windows

package pty

import (
	"os"
	"os/exec"
)

// StartWithSize assigns a pseudo-terminal tty os.File to c.Stdin, c.Stdout,
// and c.Stderr, calls c.Start, and returns the File of the tty's
// corresponding pty.
//
// This will resize the pty to the specified size before starting the command.
// Starts the process in a new session and sets the controlling terminal.
func StartWithSize(cmd *exec.Cmd, ws *Winsize) (*os.File, error) {
	return nil, ErrUnsupported
}
