//go:build !windows
// +build !windows

package oscommands

import (
	"os/exec"

	"github.com/creack/pty"
)

// we define this separately for windows and non-windows given that windows does
// not have great PTY support and we need a PTY to handle a credential request
func (self *cmdObjRunner) RunCommandWithOutputLive(
	cmdObj ICmdObj,
	output func(string) string,
) error {
	return self.RunCommandWithOutputLiveAux(
		cmdObj,
		output,
		func(cmd *exec.Cmd) (*cmdHandler, error) {
			ptmx, err := pty.Start(cmd)
			if err != nil {
				return nil, err
			}

			return &cmdHandler{
				stdoutPipe: ptmx,
				stdinPipe:  ptmx,
				close:      ptmx.Close,
			}, nil
		},
	)
}
