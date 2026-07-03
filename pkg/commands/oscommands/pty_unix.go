//go:build !windows

package oscommands

import (
	"os"
	"os/exec"

	creackpty "github.com/creack/pty"
)

type unixPty struct {
	master *os.File
}

func (u *unixPty) Read(p []byte) (int, error)  { return u.master.Read(p) }
func (u *unixPty) Write(p []byte) (int, error) { return u.master.Write(p) }
func (u *unixPty) Close() error                { return u.master.Close() }

func (u *unixPty) Resize(cols, rows uint16) error {
	return creackpty.Setsize(u.master, &creackpty.Winsize{Cols: cols, Rows: rows})
}

func StartPty(cmd *exec.Cmd, cols, rows uint16) (StartedPty, error) {
	f, err := creackpty.StartWithSize(cmd, &creackpty.Winsize{Cols: cols, Rows: rows})
	if err != nil {
		return StartedPty{}, err
	}
	return StartedPty{
		Pty:     &unixPty{master: f},
		Process: cmd.Process,
		Wait:    cmd.Wait,
	}, nil
}
