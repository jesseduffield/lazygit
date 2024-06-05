//go:build !windows
// +build !windows

package pty

import (
	"os"
	"os/exec"
	"syscall"
)

// StartWithSize assigns a pseudo-terminal tty os.File to c.Stdin, c.Stdout,
// and c.Stderr, calls c.Start, and returns the File of the tty's
// corresponding pty.
//
// This will resize the pty to the specified size before starting the command.
// Starts the process in a new session and sets the controlling terminal.
func StartWithSize(cmd *exec.Cmd, ws *Winsize) (*os.File, error) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.Setsid = true
	cmd.SysProcAttr.Setctty = true
	return StartWithAttrs(cmd, ws, cmd.SysProcAttr)
}
