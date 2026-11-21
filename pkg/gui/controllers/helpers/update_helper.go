package helpers

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/constants"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/updates"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type UpdateHelper struct {
	c               *HelperCommon
	updater         *updates.Updater
	appStatusHelper *AppStatusHelper
}

func NewUpdateHelper(c *HelperCommon, updater *updates.Updater, appStatusHelper *AppStatusHelper) *UpdateHelper {
	return &UpdateHelper{
		c:               c,
		updater:         updater,
		appStatusHelper: appStatusHelper,
	}
}

func (self *UpdateHelper) OnStartup() {
	self.c.OnWorkerBackground(func(gocui.Task) error {
		self.updater.CleanUpAfterUpdate()
		self.CheckForUpdateInBackground()
		return nil
	})
}

// CheckForUpdateInBackground checks for an update if one is due, and shows the
// result in the note in the bottom-right corner. Call it on a worker.
func (self *UpdateHelper) CheckForUpdateInBackground() {
	if err := self.updater.CheckForUpdateIfDue(); err != nil {
		self.c.Log.Errorf("Error checking for an update: %v", err)
	}
	self.rerenderNote()
}

func (self *UpdateHelper) rerenderNote() {
	// The information view is rendered in every layout, so we only need to
	// trigger one
	self.c.OnUIThreadBackground(func() error { return nil })
}

// Note returns the note that tells the user about an update, for the bottom-right
// corner, or "" if there's nothing to tell
func (self *UpdateHelper) Note() string {
	if !self.updater.ChecksEnabled() {
		return ""
	}

	status := self.updater.Status()
	if status.AvailableVersion == nil {
		return ""
	}
	return style.FgGreen.Sprint(self.c.Tr.UpdateAvailableNote)
}

func (self *UpdateHelper) CheckForUpdateDisabledReason() *types.DisabledReason {
	if !self.updater.ChecksEnabled() {
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
	self.appStatusHelper.WithBackgroundWaitingStatus(self.c.Tr.CheckingForUpdates, func() error {
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

		if self.updater.IsUpdateInstalled() {
			self.c.OnUIThread(func() error {
				self.c.Toast(self.updateInstalledMessage(newVersion))
				return nil
			})
			return nil
		}

		menuOpts := self.updateMenuOptions(newVersion)
		self.c.OnUIThread(func() error {
			return self.c.Menu(menuOpts)
		})
		return nil
	})
	return nil
}

// updateMenuOptions returns the options for the menu that tells the user about
// a new version. Call it on a worker, because it accesses the file system.
func (self *UpdateHelper) updateMenuOptions(newVersion *types.VersionNumber) types.CreateMenuOptions {
	installation := self.updater.Installation()
	paragraphs := []string{utils.ResolvePlaceholderString(
		self.c.Tr.NewVersionAvailable,
		map[string]string{
			"newVersion":     newVersion.String(),
			"currentVersion": self.updater.CurrentVersion().String(),
			"releaseNotes":   style.PrintHyperlink(self.c.Tr.ReleaseNotes, constants.Links.Releases),
		},
	)}
	items := []*types.MenuItem{}

	packageManager := installation.PackageManager()
	switch {
	case installation.Kind == updates.InstallKindRelease:
		if err := self.updater.CheckCanInstall(); err != nil {
			self.c.Log.Infof("Can't install the update: %v", err)
			paragraphs = append(paragraphs, utils.ResolvePlaceholderString(
				self.c.Tr.UpdateFolderNotWritable,
				map[string]string{"folder": filepath.Dir(installation.ExecutablePath)},
			))
			items = append(items, self.openInstallationInstructionsMenuItem())
		} else {
			items = append(items, &types.MenuItem{
				Label:   self.c.Tr.InstallUpdateNow,
				OnPress: func() error { return self.install(newVersion) },
			})
		}

	case packageManager != nil:
		paragraphs = append(paragraphs, utils.ResolvePlaceholderString(
			self.c.Tr.UpdateWithPackageManager,
			map[string]string{
				"packageManager": packageManager.Name,
				"command":        packageManager.UpdateCommand,
			},
		))
		if packageManager.CanBeSlow {
			paragraphs = append(paragraphs, self.c.Tr.PackageManagersCanBeSlow)
		}
		items = append(items, &types.MenuItem{
			Label: self.c.Tr.CopyUpdateCommand,
			OnPress: func() error {
				if err := self.c.OS().CopyToClipboard(packageManager.UpdateCommand); err != nil {
					return err
				}
				self.c.Toast(self.c.Tr.UpdateCommandCopiedToClipboard)
				return nil
			},
		})

	default:
		paragraphs = append(paragraphs, self.c.Tr.UpdateFromPackageManagerOrRelease)
		items = append(items, self.openInstallationInstructionsMenuItem())
	}

	return types.CreateMenuOptions{
		Title:  self.c.Tr.UpdateAvailableTitle,
		Prompt: strings.Join(paragraphs, "\n\n"),
		Items:  items,
	}
}

func (self *UpdateHelper) openInstallationInstructionsMenuItem() *types.MenuItem {
	return &types.MenuItem{
		Label:   self.c.Tr.OpenInstallationInstructions,
		OnPress: func() error { return self.c.OS().OpenLink(constants.Links.InstallationInstructions) },
	}
}

func (self *UpdateHelper) install(newVersion *types.VersionNumber) error {
	status := utils.ResolvePlaceholderString(
		self.c.Tr.InstallingUpdateStatus, map[string]string{"newVersion": newVersion.String()},
	)
	self.appStatusHelper.WithBackgroundWaitingStatus(status, func() error {
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
	return nil
}

func (self *UpdateHelper) updateInstalledMessage(version *types.VersionNumber) string {
	return utils.ResolvePlaceholderString(
		self.c.Tr.UpdateInstalled, map[string]string{"newVersion": version.String()},
	)
}
