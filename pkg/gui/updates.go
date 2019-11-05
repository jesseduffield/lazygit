package gui

import "github.com/jesseduffield/gocui"

func (gui *Gui) showUpdatePrompt(newVersion string) error {
	title := "New version available!"
	message := "Download latest version? (enter/esc)"
	currentView := gui.g.CurrentView()
	return gui.createConfirmationPanel(gui.g, currentView, true, title, message, func(g *gocui.Gui, v *gocui.View) error {
		gui.startUpdating(newVersion)
		return nil
	}, nil)
}

func (gui *Gui) onUserUpdateCheckFinish(newVersion string, err error) error {
	if err != nil {
		return gui.createErrorPanel(gui.g, err.Error())
	}
	if newVersion == "" {
		return gui.createErrorPanel(gui.g, "New version not found")
	}
	return gui.showUpdatePrompt(newVersion)
}

func (gui *Gui) onBackgroundUpdateCheckFinish(newVersion string, err error) error {
	if err != nil {
		// ignoring the error for now so that I'm not annoying users
		gui.Log.Error(err.Error())
		return nil
	}
	if newVersion == "" {
		return nil
	}
	if gui.Config.GetUserConfig().Get("update.method") == "background" {
		gui.startUpdating(newVersion)
		return nil
	}
	return gui.showUpdatePrompt(newVersion)
}

func (gui *Gui) startUpdating(newVersion string) {
	gui.State.Updating = true
	gui.statusManager.addWaitingStatus("updating")
	gui.Updater.Update(newVersion, gui.onUpdateFinish)
}

func (gui *Gui) onUpdateFinish(err error) error {
	gui.State.Updating = false
	gui.statusManager.removeStatus("updating")
	if err := gui.renderString(gui.g, "appStatus", ""); err != nil {
		return err
	}
	if err != nil {
		return gui.createErrorPanel(gui.g, "Update failed: "+err.Error())
	}
	return nil
}

func (gui *Gui) createUpdateQuitConfirmation(g *gocui.Gui, v *gocui.View) error {
	title := "Currently Updating"
	message := "An update is in progress. Are you sure you want to quit?"
	return gui.createConfirmationPanel(gui.g, v, true, title, message, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	}, nil)
}
