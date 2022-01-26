//go:build !windows
//+build !windows

package pty

import (
	"os/exec"
	"syscall"
)

// StartWithSize assigns a pseudo-terminal Tty to c.Stdin, c.Stdout,
// and c.Stderr, calls c.Start, and returns the File of the tty's
// corresponding Pty.
//
// This will resize the Pty to the specified size before starting the command.
// Starts the process in a new session and sets the controlling terminal.
func StartWithSize(c *exec.Cmd, sz *Winsize) (Pty, error) {
	if c.SysProcAttr == nil {
		c.SysProcAttr = &syscall.SysProcAttr{}
	}
	c.SysProcAttr.Setsid = true
	c.SysProcAttr.Setctty = true
	return StartWithAttrs(c, sz, c.SysProcAttr)
}

// StartWithAttrs assigns a pseudo-terminal Tty to c.Stdin, c.Stdout,
// and c.Stderr, calls c.Start, and returns the File of the tty's
// corresponding Pty.
//
// This will resize the Pty to the specified size before starting the command if a size is provided.
// The `attrs` parameter overrides the one set in c.SysProcAttr.
//
// This should generally not be needed. Used in some edge cases where it is needed to create a pty
// without a controlling terminal.
func StartWithAttrs(c *exec.Cmd, sz *Winsize, attrs *syscall.SysProcAttr) (Pty, error) {
	pty, tty, err := open()
	if err != nil {
		return nil, err
	}
	defer func() {
		// always close tty fds since it's being used in another process
		// but pty is kept to resize tty
		_ = tty.Close()
	}()

	if sz != nil {
		if err := Setsize(pty, sz); err != nil {
			_ = pty.Close()
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
		_ = pty.Close()
		return nil, err
	}
	return pty, err
}
