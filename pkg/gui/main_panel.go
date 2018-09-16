package gui

import (
	"github.com/jesseduffield/gocui"
)

// scrollUpMain scrolls up the main view.
// returns error if something goes wrong
func (gui *Gui) scrollUpMain(g *gocui.Gui, v *gocui.View) error {
	mainView, _ := gui.g.View("main")
	ox, oy := mainView.Origin()

	if oy >= 1 {
		err := mainView.SetOrigin(ox, oy-gui.Config.GetUserConfig().GetInt("gui.scrollHeight"))
		if err != nil {
			gui.Log.Errorf("Error while scrolling up main: %s\n", err)
			return err
		}
	}

	return nil
}

// scrollDownMain scrolls down the main view.
// returns error if something goes wrong
func (gui *Gui) scrollDownMain(g *gocui.Gui, v *gocui.View) error {
	mainView, _ := gui.g.View("main")
	ox, oy := mainView.Origin()
	if oy < len(mainView.BufferLines()) {
		err := mainView.SetOrigin(ox, oy+gui.Config.GetUserConfig().GetInt("gui.scrollHeight"))
		if err != nil {
			gui.Log.Errorf("Error while scrolling down main: %s\n", err)
			return err
		}
		return nil
	}

	return nil
}
