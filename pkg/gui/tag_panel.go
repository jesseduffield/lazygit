package gui

import (
	"github.com/jesseduffield/gocui"
	"github.com/pkg/errors"
)

func (gui *Gui) handleBranchTagSwitch(g *gocui.Gui, v *gocui.View) error {

	if v.Name() == "branches" {
		g.SetViewOnTop("tags")
		g.SetViewOnBottom("branches")
		g.SetCurrentView("tags")
	} else if v.Name() == "tags" {
		g.SetViewOnTop("branches")
		g.SetViewOnBottom("tags")
		g.SetCurrentView("branches")
	} else {
		return errors.New("unsupported view")
	}

	return nil
}
