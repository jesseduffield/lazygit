package gui

import "github.com/jesseduffield/gocui"

// showUpdatePrompt is called when there is a new update available.
// newVersion contains the string indicating the version.
// returns an error if something went wrong.
func (gui *Gui) showUpdatePrompt(newVersion string) error {
	title := "New version available!"
	message := "Download latest version? (enter/esc)"
	currentView := gui.g.CurrentView()

	err := gui.createConfirmationPanel(currentView, title, message,
		func(g *gocui.Gui, v *gocui.View) error {
			gui.startUpdating(newVersion)
			return nil
		}, nil)
	if err != nil {
		gui.Log.Errorf("Failed to createConfirmationPanel at showUpdatePrompt: %s\n", err)
		return err
	}

	return nil
}

// onUserUpdateCheckFinish is called after the update is completed.
// newVersion is the new version.
// err is the error code that was passed by the updater.
// returns an error when something goes wrong.
func (gui *Gui) onUserUpdateCheckFinish(newVersion string, err error) error {
	if err != nil {

		err = gui.createErrorPanel(err.Error())
		if err != nil {
			gui.Log.Errorf("Failed to createErrorPanel at onUserUpdateCheckFinish: %s\n", err)
			return err
		}

		return nil
	}

	if newVersion == "" {
		err = gui.createErrorPanel("New version not found")
		if err != nil {
			gui.Log.Errorf("Failed to createErrorPanel at onUserUpdateCheckFinish: %s\n", err)
			return err
		}

		return nil
	}

	err = gui.showUpdatePrompt(newVersion)
	if err != nil {
		gui.Log.Errorf("Failed to showUpdatePrompt at onUserUpdateCheckFinish: %s\n", err)
		return err
	}

	return nil
}

// onBackGroundUpdateCheckFinish is called when the background update checker
// is done.
// newVersion: the new version.
// err: the error code passed by the checker.
func (gui *Gui) onBackgroundUpdateCheckFinish(newVersion string, err error) error {

	if err != nil {
		gui.Log.Errorf("Failed to update at onBackgroundUpdateCheckFinish: %s\n", err)
		return err
	}

	if newVersion == "" {
		return nil
	}

	if gui.Config.GetUserConfig().Get("update.method") == "background" {
		gui.startUpdating(newVersion)
		return nil
	}

	err = gui.showUpdatePrompt(newVersion)
	if err != nil {
		gui.Log.Errorf("Failed to showUpdatePrompt at onBackgroundUpdateCheckFinish: %s\n", err)
		return err
	}

	return nil
}

// startUpdating actually updates the application.
// newVersion: the new version.
func (gui *Gui) startUpdating(newVersion string) {
	gui.State.Updating = true
	gui.statusManager.addWaitingStatus("updating")
	gui.Updater.Update(newVersion, gui.onUpdateFinish)
}

// onUpdaterFinish gets called when the updater finishes.
// code: the error code passed by the updater.
// returns an error when something goes wrong.
func (gui *Gui) onUpdateFinish(code error) error {
	gui.State.Updating = false
	gui.statusManager.removeStatus("updating")

	err := gui.renderString("appStatus", "")
	if err != nil {
		gui.Log.Errorf("Failed to renderString at onUpdateFinsish: %s\n", err)
		return err
	}

	if code != nil {
		err = gui.createErrorPanel("Update failed: " + code.Error())
		if err != nil {
			gui.Log.Errorf("Failed to createErrorPanel at onUpdateFinsish: %s\n", err)
			return err
		}
	}

	return nil
}

// createUpdateQuitConfirmation gets called when the user wants to quit the
// program while it is updating.
// v: the view to returns focus to.
// returns an error when something goes wrong.
func (gui *Gui) createUpdateQuitConfirmation(v *gocui.View) error {
	title := "Currently Updating"
	message := "An update is in progress. Are you sure you want to quit?"

	err := gui.createConfirmationPanel(v, title, message,
		func(g *gocui.Gui, v *gocui.View) error {
			return gocui.ErrQuit
		}, nil)
	if err != nil {
		gui.Log.Errorf("Failed to createConfirmationPanel at createUpdateQuitConfirmation: %s\n", err)
		return err
	}

	return nil
}
