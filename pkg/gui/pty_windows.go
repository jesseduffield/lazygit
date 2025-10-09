package gui

import (
	"fmt"
	"os/exec"

	"github.com/jesseduffield/gocui"
)

func (gui *Gui) onResize() error {
	return nil
}

func (gui *Gui) newPtyTask(view *gocui.View, cmd *exec.Cmd, prefix string) error {
	cmd.Env = append(cmd.Env, fmt.Sprintf("LAZYGIT_COLUMNS=%d", view.InnerWidth()))
	return gui.newCmdTask(view, cmd, prefix)
}
