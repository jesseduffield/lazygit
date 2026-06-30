package oscommands

import (
	"os/exec"
)

// StartPty is a stub on Windows for now; callers fall back to the non-pty
// path when ErrPtyUnsupported is returned. A real ConPTY implementation
// replaces this in a follow-up commit.
func StartPty(cmd *exec.Cmd, cols, rows uint16) (StartedPty, error) {
	return StartedPty{}, ErrPtyUnsupported
}
