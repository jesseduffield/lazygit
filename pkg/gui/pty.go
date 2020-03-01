package gui

import (
	"github.com/jesseduffield/pty"
)

func (gui *Gui) onResize() error {
	if gui.State.Ptmx == nil {
		return nil
	}
	mainView := gui.getMainView()
	width, height := mainView.Size()

	if err := pty.Setsize(gui.State.Ptmx, &pty.Winsize{Cols: uint16(width), Rows: uint16(height)}); err != nil {
		return err
	}

	// TODO: handle resizing properly

	return nil
}
