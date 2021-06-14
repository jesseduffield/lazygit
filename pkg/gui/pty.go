// +build !windows

package gui

import (
	"github.com/creack/pty"
	"github.com/jesseduffield/gocui"
	. "github.com/jesseduffield/lazygit/pkg/commands/types"
)

func (gui *Gui) onResize() error {
	if gui.State.Ptmx == nil {
		return nil
	}
	width, height := gui.Views.Main.Size()

	if err := pty.Setsize(gui.State.Ptmx, &pty.Winsize{Cols: uint16(width), Rows: uint16(height)}); err != nil {
		return err
	}

	// TODO: handle resizing properly

	return nil
}

// Some commands need to output for a terminal to active certain behaviour.
// For example,  git won't invoke the GIT_PAGER env var unless it thinks it's
// talking to a terminal. We typically write cmd outputs straight to a view,
// which is just an io.Reader. the pty package lets us wrap a command in a
// pseudo-terminal meaning we'll get the behaviour we want from the underlying
// command.
func (gui *Gui) newPtyTask(view *gocui.View, cmdObj ICmdObj, prefix string) error {
	width, _ := gui.Views.Main.Size()
	pager := gui.Git.GetPager(width)

	if pager == "" {
		// if we're not using a custom pager we don't need to use a pty
		return gui.newCmdTask(view, cmdObj, prefix)
	}

	cmdObj.AddEnvVars("GIT_PAGER=" + pager)

	_, height := view.Size()
	_, oy := view.Origin()

	manager := gui.getManager(view)

	cmd := cmdObj.GetCmd()
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return err
	}

	gui.State.Ptmx = ptmx
	onClose := func() {
		ptmx.Close()
		gui.State.Ptmx = nil
	}

	if err := gui.onResize(); err != nil {
		return err
	}

	if err := manager.NewTask(manager.NewCmdTask(ptmx, cmd, prefix, height+oy+10, onClose)); err != nil {
		return err
	}

	return nil
}
