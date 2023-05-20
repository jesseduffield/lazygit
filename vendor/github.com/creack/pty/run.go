package pty

import (
	"os/exec"
)

// Start assigns a pseudo-terminal tty os.File to c.Stdin, c.Stdout,
// and c.Stderr, calls c.Start, and returns the File of the tty's
// corresponding pty.
//
// Starts the process in a new session and sets the controlling terminal.
func Start(cmd *exec.Cmd) (Pty, error) {
	return StartWithSize(cmd, nil)
}
