package gui

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) showUpdatePrompt(newVersion string) error {
	message := utils.ResolvePlaceholderString(
		gui.Tr.UpdateAvailable, map[string]string{
			"newVersion": newVersion,
		},
	)

	return gui.c.Confirm(types.ConfirmOpts{
		Title:  gui.Tr.UpdateAvailableTitle,
		Prompt: message,
		HandleConfirm: func() error {
			gui.startUpdating(newVersion)
			return nil
		},
	})
}

func (gui *Gui) onUserUpdateCheckFinish(newVersion string, err error) error {
	if err != nil {
		return gui.c.Error(err)
	}
	if newVersion == "" {
		return gui.c.ErrorMsg(gui.Tr.FailedToRetrieveLatestVersionErr)
	}
	return gui.showUpdatePrompt(newVersion)
}

func (gui *Gui) onBackgroundUpdateCheckFinish(newVersion string, err error) error {
	if err != nil {
		// ignoring the error for now so that I'm not annoying users
		gui.c.Log.Error(err.Error())
		return nil
	}
	if newVersion == "" {
		return nil
	}
	if gui.c.UserConfig.Update.Method == "background" {
		gui.startUpdating(newVersion)
		return nil
	}
	return gui.showUpdatePrompt(newVersion)
}

func (gui *Gui) startUpdating(newVersion string) {
	gui.State.Updating = true
	statusId := gui.statusManager.addWaitingStatus(gui.Tr.UpdateInProgressWaitingStatus)
	gui.Updater.Update(newVersion, func(err error) error { return gui.onUpdateFinish(statusId, err) })
}

func (gui *Gui) onUpdateFinish(statusId int, err error) error {
	gui.State.Updating = false
	gui.statusManager.removeStatus(statusId)
	gui.OnUIThread(func() error {
		_ = gui.renderString(gui.Views.AppStatus, "")
		if err != nil {
			errMessage := utils.ResolvePlaceholderString(
				gui.Tr.UpdateFailedErr, map[string]string{
					"errMessage": err.Error(),
				},
			)
			return gui.c.ErrorMsg(errMessage)
		}
		return gui.c.Alert(gui.Tr.UpdateCompletedTitle, gui.Tr.UpdateCompleted)
	})

	return nil
}

func (gui *Gui) createUpdateQuitConfirmation() error {
	return gui.c.Confirm(types.ConfirmOpts{
		Title:  gui.Tr.ConfirmQuitDuringUpdateTitle,
		Prompt: gui.Tr.ConfirmQuitDuringUpdate,
		HandleConfirm: func() error {
			return gocui.ErrQuit
		},
	})
}
