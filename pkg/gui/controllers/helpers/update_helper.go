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
	// Blocks until the user has closed the intro or breaking-changes popup
	// that lazygit shows at startup. The update menu waits for this, so that
	// it doesn't replace them.
	waitForStartupPopups func()
}

func NewUpdateHelper(
	c *HelperCommon,
	updater *updates.Updater,
	appStatusHelper *AppStatusHelper,
	waitForStartupPopups func(),
) *UpdateHelper {
	return &UpdateHelper{
		c:                    c,
		updater:              updater,
		appStatusHelper:      appStatusHelper,
		waitForStartupPopups: waitForStartupPopups,
	}
}

// The link of the note that tells the user about an update
const UpdateNoteLink = "lazygit-update:"

// OnStartup checks for an update if a check is due, and shows the update menu
// if a reminder is due, once the user has closed the intro or breaking-changes
// popup. If the recent repos menu is open, the update menu replaces it.
//
// The check runs at the same time, so the menu is for an update that an
// earlier check found, unless the check finishes while a startup popup is
// open.
func (self *UpdateHelper) OnStartup() {
	self.c.OnWorkerBackground(func(gocui.Task) error {
		self.updater.CleanUpAfterUpdate()
		self.CheckForUpdateInBackground()
		return nil
	})

	go utils.Safe(func() {
		self.waitForStartupPopups()
		self.remindIfDue()
	})
}

// remindIfDue shows the update menu if a reminder is due. We only do this at
// startup, before the user starts working; a menu that appears later could
// catch a key that the user meant for something else. Call it on a worker,
// because it accesses the file system.
func (self *UpdateHelper) remindIfDue() {
	version := self.updater.DueReminder()
	if version == nil {
		return
	}

	self.showUpdateMenu(version, true)
}

// CheckForUpdateInBackground checks for an update if one is due, and shows the
// result in the note in the bottom-right corner. If the user config asks for it,
// it also installs the update. Call it on a worker.
func (self *UpdateHelper) CheckForUpdateInBackground() {
	if err := self.updater.CheckForUpdateIfDue(); err != nil {
		self.c.Log.Errorf("Error checking for an update: %v", err)
	}

	if version := self.updater.VersionToInstallInBackground(); version != nil {
		if err := self.updater.Install(version); err != nil {
			// We try again in an hour, and at the next start
			self.c.Log.Errorf("Error installing %s: %v", version, err)
		} else {
			self.c.OnUIThreadBackground(func() error {
				self.c.Toast(self.updateInstalledMessage(version))
				return nil
			})
		}
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
	if status.InstalledVersion != nil {
		return style.FgGreen.Sprint(utils.ResolvePlaceholderString(
			self.c.Tr.RestartToUseNewVersion, map[string]string{"newVersion": status.InstalledVersion.String()},
		))
	}
	if status.AvailableVersion == nil {
		return ""
	}
	return style.FgGreen.Sprint(style.PrintHyperlink(self.c.Tr.UpdateAvailableNote, UpdateNoteLink))
}

// ShowUpdateMenu shows the menu for the update that the note in the
// bottom-right corner tells the user about
func (self *UpdateHelper) ShowUpdateMenu() error {
	status := self.updater.Status()
	if status.AvailableVersion == nil || status.InstalledVersion != nil {
		return nil
	}

	self.c.OnWorker(func(gocui.Task) error {
		self.showUpdateMenu(status.AvailableVersion, false)
		return nil
	})
	return nil
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

		if installedVersion := self.updater.InstalledVersion(); installedVersion != nil {
			self.c.OnUIThread(func() error {
				self.c.Toast(self.updateInstalledMessage(installedVersion))
				return nil
			})
			return nil
		}

		self.showUpdateMenu(newVersion, false)
		return nil
	})
	return nil
}

// showUpdateMenu shows the menu that tells the user about a new version; see
// updateMenuOptions for what reminder means. Call it on a worker, because
// building the menu accesses the file system.
func (self *UpdateHelper) showUpdateMenu(newVersion *types.VersionNumber, reminder bool) {
	menuOpts := self.updateMenuOptions(newVersion, reminder)
	self.c.OnUIThread(func() error {
		return self.c.Menu(menuOpts)
	})
}

// updateMenuOptions returns the options for the menu that tells the user about
// a new version. reminder says whether lazygit shows the menu on its own; the
// menu then tells the user how to turn reminders off. Call it on a worker,
// because it accesses the file system.
func (self *UpdateHelper) updateMenuOptions(newVersion *types.VersionNumber, reminder bool) types.CreateMenuOptions {
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

	if reminder {
		paragraphs = append(paragraphs, self.c.Tr.TurnOffUpdateReminders)
	}

	// Whatever the user chooses, including closing the menu, postpones the next
	// reminder until tomorrow, unless they skip the version
	remindTomorrow := func() {
		self.c.OnWorkerBackground(func(gocui.Task) error {
			self.updater.RemindTomorrow()
			return nil
		})
	}
	for _, item := range items {
		onPress := item.OnPress
		item.OnPress = func() error {
			remindTomorrow()
			return onPress()
		}
	}
	items = append(items,
		&types.MenuItem{
			Label: self.c.Tr.RemindMeTomorrow,
			OnPress: func() error {
				remindTomorrow()
				return nil
			},
		},
		&types.MenuItem{
			Label: self.c.Tr.SkipThisVersion,
			OnPress: func() error {
				self.c.OnWorkerBackground(func(gocui.Task) error {
					self.updater.SkipVersion(newVersion)
					return nil
				})
				return nil
			},
		},
	)

	return types.CreateMenuOptions{
		Title:      self.c.Tr.UpdateAvailableTitle,
		Prompt:     strings.Join(paragraphs, "\n\n"),
		Items:      items,
		HideCancel: true,
		OnCancel: func() error {
			remindTomorrow()
			return nil
		},
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
