package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
)

func (gui *Gui) refreshCommandStatus(g *gocui.Gui) error {
	v, err := g.View("commandStatus")
	if err != nil {
		return err
	}

	g.Update(func(*gocui.Gui) error {
		v.Clear()

		fmt.Fprintf(v, gui.CmdStatus)
		return nil

	})
	return nil
}
