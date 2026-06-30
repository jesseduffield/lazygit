package oscommands

import (
	"io"
	"os"
)

// Pty is the master side of a pseudo-terminal running a subprocess. The
// concrete implementation is platform-specific: creack/pty on Unix and
// ConPTY on Windows.
type Pty interface {
	io.ReadWriteCloser
	Resize(cols, rows uint16) error
}

// StartedPty is the result of StartPty.
type StartedPty struct {
	// Pty is the master side of the pseudo-terminal; read from it to get
	// the child's combined stdout/stderr and write to it to feed stdin.
	Pty Pty
	// Process is the spawned child. Useful for signalling; on Windows the
	// original *exec.Cmd was not Start()ed (ConPTY spawns via
	// CreateProcess, not os/exec) so cmd.Process is nil and this is the
	// only handle.
	Process *os.Process
	// Wait blocks until the child exits and returns a non-nil error on a
	// nonzero exit status, matching *exec.Cmd.Wait semantics.
	Wait func() error
}

// StartPty runs cmd in a pseudo-terminal with the given initial dimensions.
// Implemented per-platform in pty_unix.go / pty_windows.go.
//
// func StartPty(cmd *exec.Cmd, cols, rows uint16) (StartedPty, error)
