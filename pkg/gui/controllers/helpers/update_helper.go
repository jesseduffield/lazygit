package helpers

import (
	"errors"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/updates"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type UpdateHelper struct {
	c       *HelperCommon
	updater *updates.Updater
}

func NewUpdateHelper(c *HelperCommon, updater *updates.Updater) *UpdateHelper {
	return &UpdateHelper{
		c:       c,
		updater: updater,
	}
}

func (self *UpdateHelper) CheckForUpdateInBackground() {
	self.updater.CheckForNewUpdate(func(newVersion string, err error) error {
		if err != nil {
			// ignoring the error for now so that I'm not annoying users
			self.c.Log.Error(err.Error())
			return nil
		}
		if newVersion == "" {
			return nil
		}
		if self.c.UserConfig.Update.Method == "background" {
			self.startUpdating(newVersion)
			return nil
		}
		return self.showUpdatePrompt(newVersion)
	}, false)
}

func (self *UpdateHelper) CheckForUpdateInForeground() error {
	return self.c.WithWaitingStatus(self.c.Tr.CheckingForUpdates, func(gocui.Task) error {
		self.updater.CheckForNewUpdate(func(newVersion string, err error) error {
			if err != nil {
				return err
			}
			if newVersion == "" {
				return errors.New(self.c.Tr.FailedToRetrieveLatestVersionErr)
			}
			return self.showUpdatePrompt(newVersion)
		}, true)

		return nil
	})
}

func (self *UpdateHelper) startUpdating(newVersion string) {
	_ = self.c.WithWaitingStatus(self.c.Tr.UpdateInProgressWaitingStatus, func(gocui.Task) error {
		self.c.State().SetUpdating(true)
		err := self.updater.Update(newVersion)
		return self.onUpdateFinish(err)
	})
}

func (self *UpdateHelper) onUpdateFinish(err error) error {
	self.c.State().SetUpdating(false)
	self.c.OnUIThread(func() error {
		self.c.SetViewContent(self.c.Views().AppStatus, "")
		if err != nil {
			errMessage := utils.ResolvePlaceholderString(
				self.c.Tr.UpdateFailedErr, map[string]string{
					"errMessage": err.Error(),
				},
			)
			return errors.New(errMessage)
		}
		return self.c.Alert(self.c.Tr.UpdateCompletedTitle, self.c.Tr.UpdateCompleted)
	})

	return nil
}

func (self *UpdateHelper) showUpdatePrompt(newVersion string) error {
	message := utils.ResolvePlaceholderString(
		self.c.Tr.UpdateAvailable, map[string]string{
			"newVersion": newVersion,
		},
	)

	return self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.UpdateAvailableTitle,
		Prompt: message,
		HandleConfirm: func() error {
			self.startUpdating(newVersion)
			return nil
		},
	})
}
