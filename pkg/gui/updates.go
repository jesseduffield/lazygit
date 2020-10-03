package gui

import "github.com/jesseduffield/gocui"

func (gui *Gui) showUpdatePrompt(newVersion string) error {
	return gui.ask(askOpts{
		title:  "New version available!",
		prompt: "Download latest version? (enter/esc)",
		handleConfirm: func() error {
			gui.startUpdating(newVersion)
			return nil
		},
	})
}

func (gui *Gui) onUserUpdateCheckFinish(newVersion string, err error) error {
	if err != nil {
		return gui.surfaceError(err)
	}
	if newVersion == "" {
		return gui.createErrorPanel("New version not found")
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
	if gui.Config.GetUserConfig().Update.Method == "background" {
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
	gui.renderString("appStatus", "")
	if err != nil {
		return gui.createErrorPanel("Update failed: " + err.Error())
	}
	return nil
}

func (gui *Gui) createUpdateQuitConfirmation() error {
	return gui.ask(askOpts{
		title:  "Currently Updating",
		prompt: "An update is in progress. Are you sure you want to quit?",
		handleConfirm: func() error {
			return gocui.ErrQuit
		},
	})
}
