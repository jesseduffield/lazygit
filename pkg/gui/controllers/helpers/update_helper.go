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

func (self *UpdateHelper) OnStartup() {
	self.c.OnWorkerBackground(func(gocui.Task) error {
		self.updater.CleanUpAfterUpdate()
		return nil
	})
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

	if self.updater.IsInstalling() {
		return &types.DisabledReason{Text: self.c.Tr.UpdateAlreadyInProgress}
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

		return self.showUpdatePrompt(newVersion)
	})
}

func (self *UpdateHelper) showUpdatePrompt(newVersion *types.VersionNumber) error {
	message := utils.ResolvePlaceholderString(
		self.c.Tr.UpdateAvailable, map[string]string{
			"newVersion": newVersion.String(),
		},
	)

	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.UpdateAvailableTitle,
		Prompt: message,
		HandleConfirm: func() error {
			return self.install(newVersion)
		},
	})

	return nil
}

func (self *UpdateHelper) install(newVersion *types.VersionNumber) error {
	status := utils.ResolvePlaceholderString(
		self.c.Tr.InstallingUpdateStatus, map[string]string{"newVersion": newVersion.String()},
	)
	return self.c.WithWaitingStatus(status, func(gocui.Task) error {
		err := self.updater.Install(newVersion)
		self.c.OnUIThread(func() error {
			if err != nil {
				return errors.New(utils.ResolvePlaceholderString(
					self.c.Tr.UpdateInstallFailed,
					map[string]string{"newVersion": newVersion.String(), "error": err.Error()},
				))
			}

			self.c.Toast(self.updateInstalledMessage(newVersion))
			return nil
		})
		return nil
	})
}

func (self *UpdateHelper) updateInstalledMessage(version *types.VersionNumber) string {
	return utils.ResolvePlaceholderString(
		self.c.Tr.UpdateInstalled, map[string]string{"newVersion": version.String()},
	)
}
