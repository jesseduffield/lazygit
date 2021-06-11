// +build windows

package gui

import (
	"github.com/jesseduffield/gocui"
	. "github.com/jesseduffield/lazygit/pkg/commands/types"
)

func (gui *Gui) onResize() error {
	return nil
}

func (gui *Gui) newPtyTask(view *gocui.View, cmdObj ICmdObj, prefix string) error {
	return gui.newCmdTask(view, cmdObj, prefix)
}
