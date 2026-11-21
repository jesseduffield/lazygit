package helpers

import (
	"errors"

	"github.com/jesseduffield/lazygit/pkg/gocui"
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

func (self *UpdateHelper) CheckForUpdateDisabledReason() *types.DisabledReason {
	if self.updater.Installation().Kind == updates.InstallKindDevBuild {
		return &types.DisabledReason{
			Text: utils.ResolvePlaceholderString(
				self.c.Tr.UpdateCheckNotPossibleForDevBuild,
				map[string]string{"version": self.c.GetConfig().GetVersion()},
			),
			// The user asked for the check, so make sure they see why it
			// doesn't happen
			ShowErrorInPanel: true,
		}
	}

	return nil
}

func (self *UpdateHelper) CheckForUpdateInForeground() error {
	return self.c.WithWaitingStatus(self.c.Tr.CheckingForUpdates, func(gocui.Task) error {
		newVersion, err := self.updater.CheckForUpdate()
		if err != nil {
			return errors.New(utils.ResolvePlaceholderString(
				self.c.Tr.UpdateCheckFailed, map[string]string{"error": err.Error()},
			))
		}

		if newVersion == nil {
			self.c.OnUIThread(func() error {
				self.c.Alert(self.c.Tr.NoUpdateAvailableTitle, utils.ResolvePlaceholderString(
					self.c.Tr.NoUpdateAvailable,
					map[string]string{"currentVersion": self.updater.CurrentVersion().String()},
				))
				return nil
			})
			return nil
		}

		return self.showUpdatePrompt(newVersion.String())
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
		self.c.Alert(self.c.Tr.UpdateCompletedTitle, self.c.Tr.UpdateCompleted)
		return nil
	})

	return nil
}

func (self *UpdateHelper) showUpdatePrompt(newVersion string) error {
	message := utils.ResolvePlaceholderString(
		self.c.Tr.UpdateAvailable, map[string]string{
			"newVersion": newVersion,
		},
	)

	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.UpdateAvailableTitle,
		Prompt: message,
		HandleConfirm: func() error {
			self.startUpdating(newVersion)
			return nil
		},
	})

	return nil
}
