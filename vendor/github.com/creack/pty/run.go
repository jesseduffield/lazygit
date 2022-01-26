package pty

import (
	"os/exec"
)

// Start assigns a pseudo-terminal Tty to c.Stdin, c.Stdout,
// and c.Stderr, calls c.Start, and returns the File of the tty's
// corresponding Pty.
//
// Starts the process in a new session and sets the controlling terminal.
func Start(c *exec.Cmd) (pty Pty, err error) {
	return StartWithSize(c, nil)
}
