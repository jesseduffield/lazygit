package oscommands

import (
	"os/exec"
)

func (self *cmdObjRunner) getCmdHandlerPty(cmd *exec.Cmd) (*cmdHandler, error) {
	// We don't have PTY support on Windows yet, so we just return a non-PTY handler.
	return self.getCmdHandlerNonPty(cmd)
}
