//go:build !windows

package gui

import (
	"os"
	"os/exec"

	creackpty "github.com/creack/pty"
	"github.com/jesseduffield/lazygit/pkg/tasks"
)

// setPlatformPtyEnvVars is a no-op on Unix: pagers discover the terminal
// size via the pty's TIOCGWINSZ. The Windows implementation sets
// LAZYGIT_COLUMNS for backwards compat with scripts documented in
// docs/Custom_Pagers.md.
func setPlatformPtyEnvVars(_ *exec.Cmd, _ int) {}

type unixPty struct {
	master *os.File
}

func (u *unixPty) Read(p []byte) (int, error) { return u.master.Read(p) }
func (u *unixPty) Close() error               { return u.master.Close() }

func (u *unixPty) Resize(cols, rows uint16) error {
	return creackpty.Setsize(u.master, &creackpty.Winsize{Cols: cols, Rows: rows})
}

func startPty(cmd *exec.Cmd, cols, rows uint16) (pty, tasks.Cmd, error) {
	f, err := creackpty.StartWithSize(cmd, &creackpty.Winsize{Cols: cols, Rows: rows})
	if err != nil {
		return nil, nil, err
	}
	return &unixPty{master: f}, tasks.ExecCmd{Cmd: cmd}, nil
}
