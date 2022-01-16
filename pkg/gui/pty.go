//go:build !windows
// +build !windows

package gui

import (
	"io"
	"os/exec"
	"strings"

	"github.com/creack/pty"
	"github.com/jesseduffield/gocui"
)

func (gui *Gui) desiredPtySize() *pty.Winsize {
	width, height := gui.Views.Main.Size()

	return &pty.Winsize{Cols: uint16(width), Rows: uint16(height)}
}

func (gui *Gui) onResize() error {
	if gui.State.Ptmx == nil {
		return nil
	}

	// TODO: handle resizing properly: we need to actually clear the main view
	// and re-read the output from our pty. Or we could just re-run the original
	// command from scratch
	if err := pty.Setsize(gui.State.Ptmx, gui.desiredPtySize()); err != nil {
		return err
	}

	return nil
}

// Some commands need to output for a terminal to active certain behaviour.
// For example,  git won't invoke the GIT_PAGER env var unless it thinks it's
// talking to a terminal. We typically write cmd outputs straight to a view,
// which is just an io.Reader. the pty package lets us wrap a command in a
// pseudo-terminal meaning we'll get the behaviour we want from the underlying
// command.
func (gui *Gui) newPtyTask(view *gocui.View, cmd *exec.Cmd, prefix string) error {
	width, _ := gui.Views.Main.Size()
	pager := gui.git.Config.GetPager(width)

	if pager == "" {
		// if we're not using a custom pager we don't need to use a pty
		return gui.newCmdTask(view, cmd, prefix)
	}

	cmdStr := strings.Join(cmd.Args, " ")

	cmd.Env = append(cmd.Env, "GIT_PAGER="+pager)

	_, height := view.Size()
	_, oy := view.Origin()

	manager := gui.getManager(view)

	start := func() (*exec.Cmd, io.Reader) {
		ptmx, err := pty.StartWithSize(cmd, gui.desiredPtySize())
		if err != nil {
			gui.c.Log.Error(err)
		}

		gui.State.Ptmx = ptmx

		return cmd, ptmx
	}

	onClose := func() {
		gui.State.Ptmx.Close()
		gui.State.Ptmx = nil
	}

	if err := manager.NewTask(manager.NewCmdTask(start, prefix, height+oy+10, onClose), cmdStr); err != nil {
		return err
	}

	return nil
}
