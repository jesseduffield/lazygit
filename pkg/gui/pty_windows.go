package gui

import (
	"os/exec"

	"github.com/jesseduffield/lazygit/pkg/tasks"
)

const ptySupported = false

func startPty(cmd *exec.Cmd, cols, rows uint16) (pty, tasks.Cmd, error) {
	return nil, nil, errPtyUnsupported
}
