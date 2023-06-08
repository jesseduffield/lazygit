//go:build windows
// +build windows

package oscommands

import (
	"os/exec"

	"github.com/creack/pty"
)

func (self *cmdObjRunner) getCmdHandler(cmd *exec.Cmd) (*cmdHandler, error) {
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return nil, err
	}

	return &cmdHandler{
		stdoutPipe: ptmx,
		stdinPipe:  ptmx,
		close:      ptmx.Close,
	}, nil
}
