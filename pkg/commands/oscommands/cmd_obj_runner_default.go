//go:build !windows
// +build !windows

package oscommands

import (
	"fmt"
	"os/exec"

	"github.com/creack/pty"
)

// we define this separately for windows and non-windows given that windows does
// not have great PTY support and we need a PTY to handle a credential request
func (self *cmdObjRunner) getCmdHandlerPty(cmd *exec.Cmd) (*cmdHandler, error) {
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return nil, err
	}

	LogCmd(fmt.Sprintf("Started cmd: %s, pid: %d", cmd.Args, cmd.Process.Pid))

	return &cmdHandler{
		stdoutPipe: ptmx,
		stdinPipe:  ptmx,
		close:      ptmx.Close,
	}, nil
}
