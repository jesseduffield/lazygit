//go:build !windows
// +build !windows

package oscommands

import (
	"io"
	"os/exec"

	"github.com/creack/pty"
)

func RunCommandWithOutputLiveWrapper(
	c *OSCommand,
	cmdObj ICmdObj,
	writer io.Writer,
	output func(string) string,
) error {
	return RunCommandWithOutputLiveAux(
		c,
		cmdObj,
		writer,
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
