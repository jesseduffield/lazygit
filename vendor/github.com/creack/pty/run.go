//go:build !windows
//+build !windows

package pty

import (
	"os"
	"os/exec"
	"syscall"
)

// Start assigns a pseudo-terminal tty os.File to c.Stdin, c.Stdout,
// and c.Stderr, calls c.Start, and returns the File of the tty's
// corresponding pty.
//
// Starts the process in a new session and sets the controlling terminal.
func Start(cmd *exec.Cmd) (*os.File, error) {
	return StartWithSize(cmd, nil)
}

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

// StartWithAttrs assigns a pseudo-terminal tty os.File to c.Stdin, c.Stdout,
// and c.Stderr, calls c.Start, and returns the File of the tty's
// corresponding pty.
//
// This will resize the pty to the specified size before starting the command if a size is provided.
// The `attrs` parameter overrides the one set in c.SysProcAttr.
//
// This should generally not be needed. Used in some edge cases where it is needed to create a pty
// without a controlling terminal.
func StartWithAttrs(c *exec.Cmd, sz *Winsize, attrs *syscall.SysProcAttr) (*os.File, error) {
	pty, tty, err := Open()
	if err != nil {
		return nil, err
	}
	defer func() { _ = tty.Close() }() // Best effort.

	if sz != nil {
		if err := Setsize(pty, sz); err != nil {
			_ = pty.Close() // Best effort.
			return nil, err
		}
	}
	if c.Stdout == nil {
		c.Stdout = tty
	}
	if c.Stderr == nil {
		c.Stderr = tty
	}
	if c.Stdin == nil {
		c.Stdin = tty
	}

	c.SysProcAttr = attrs

	if err := c.Start(); err != nil {
		_ = pty.Close() // Best effort.
		return nil, err
	}
	return pty, err
}
