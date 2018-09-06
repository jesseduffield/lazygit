package gui

import (
	"github.com/jesseduffield/gocui"
	"github.com/pkg/errors"
)

func (gui *Gui) handleBranchTagSwitch(g *gocui.Gui, v *gocui.View) error {

	if v.Name() == "branches" {

		err := gui.switchLayers(g, "branches", "tags")
		if err != nil {
			return err
		}

		_, err = g.SetCurrentView("tags")
		if err != nil {
			return err
		}

	} else if v.Name() == "tags" {

		err := gui.switchLayers(g, "tags", "branches")
		if err != nil {
			return err
		}

		_, err = g.SetCurrentView("branches")
		if err != nil {
			return err
		}

	} else {
		return errors.New("unsupported view")
	}

	return nil
}

func (gui *Gui) switchLayers(g *gocui.Gui, oldTop string, newTop string) error {

	_, err := g.SetViewOnBottom(oldTop)
	if err != nil {
		return err
	}

	_, err = g.SetViewOnTop(newTop)
	if err != nil {
		return err
	}

	return nil
}
