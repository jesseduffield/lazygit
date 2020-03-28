package gui

import "github.com/jesseduffield/gocui"

func (gui *Gui) inFilterMode() bool {
	return gui.State.FilterPath != ""
}

func (gui *Gui) validateNotInFilterMode() (bool, error) {
	if gui.inFilterMode() {
		return false, gui.createConfirmationPanel(gui.g, gui.g.CurrentView(), true, gui.Tr.SLocalize("MustExitFilterModeTitle"), gui.Tr.SLocalize("MustExitFilterModePrompt"), func(*gocui.Gui, *gocui.View) error {
			return gui.exitFilterMode()
		}, nil)
	}
	return true, nil
}

func (gui *Gui) exitFilterMode() error {
	gui.State.FilterPath = ""
	return gui.Errors.ErrRestart
}
