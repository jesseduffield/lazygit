//go:build windows
// +build windows

package gui

import (
	"os/exec"

	"github.com/jesseduffield/gocui"
)

func (gui *Gui) onResize() error {
	return nil
}

func (gui *Gui) newPtyTask(view *gocui.View, cmd *exec.Cmd, prefix string) error {
	return gui.newCmdTask(view, cmd, prefix)
}
