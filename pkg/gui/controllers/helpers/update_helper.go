package helpers

import (
	"errors"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/jesseduffield/lazygit/pkg/constants"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/updates"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type UpdateHelper struct {
	c                  *HelperCommon
	updater            *updates.Updater
	confirmationHelper *ConfirmationHelper

	// When we showed the update menu (Unix time), or 0 if it isn't showing
	menuShownAt atomic.Int64
	// The prompt of the update menu, to recognize it. Only accessed on the UI
	// thread.
	menuPrompt string
}

func NewUpdateHelper(c *HelperCommon, updater *updates.Updater, confirmationHelper *ConfirmationHelper) *UpdateHelper {
	return &UpdateHelper{
		c:                  c,
		updater:            updater,
		confirmationHelper: confirmationHelper,
	}
}

// The link of the note that tells the user about an update
const UpdateNoteLink = "lazygit-update:"

// How long the user must not have used lazygit before we remind them of an
// update
const idleTimeBeforeReminder = time.Minute

func (self *UpdateHelper) OnStartup() {
	self.c.OnWorkerBackground(func(gocui.Task) error {
		self.updater.CleanUpAfterUpdate()
		self.CheckForUpdateInBackground()
		self.RemindIfDue()
		return nil
	})
}

// RemindIfDue shows the update menu if the user hasn't used lazygit for a
// while, and a reminder is due. Every lazygit instance that is idle shows it,
// so that one that the user forgot about in some terminal tab doesn't keep the
// others from reminding them. When the user responds to the menu in one of
// them, the others close theirs. Call it on a worker.
func (self *UpdateHelper) RemindIfDue() {
	if shownAt := self.menuShownAt.Load(); shownAt != 0 {
		if self.updater.IsHandledAfter(time.Unix(shownAt, 0)) {
			self.c.OnUIThreadBackground(self.closeUpdateMenu)
		}
		return
	}

	reminder := updates.ReminderMenuWhenIdle
	if self.updater.IsStartingUp() {
		reminder = updates.ReminderMenuAtStartup
	} else if self.c.GocuiGui().TimeSinceLastInput() < idleTimeBeforeReminder {
		return
	}

	version := self.updater.DueReminder(reminder)
	if version == nil {
		return
	}

	menuOpts := self.updateMenuOptions(version)
	self.c.OnUIThreadBackground(func() error {
		// Don't get in the way of whatever the user is doing in a popup; we'll
		// try again later
		if self.confirmationHelper.IsPopupPanelFocused() {
			return nil
		}
		return self.showUpdateMenu(menuOpts)
	})
}

// showUpdateMenu shows the update menu, and remembers that it is showing. Call
// it on the UI thread.
func (self *UpdateHelper) showUpdateMenu(menuOpts types.CreateMenuOptions) error {
	self.menuPrompt = menuOpts.Prompt
	self.menuShownAt.Store(time.Now().Unix())
	return self.c.Menu(menuOpts)
}

// closeUpdateMenu closes the update menu without counting this as a response,
// because the user responded in another lazygit. Call it on the UI thread.
func (self *UpdateHelper) closeUpdateMenu() error {
	if self.menuShownAt.Swap(0) == 0 {
		return nil
	}

	// Leave the menu alone if something else is on top of it; the user can
	// close it themselves then
	menu := self.c.Contexts().Menu
	if self.c.Context().Current().GetKey() == context.MENU_CONTEXT_KEY && menu.GetPrompt() == self.menuPrompt {
		self.c.Context().Pop()
	}
	return nil
}

// recordResponse records that the user responded to the update menu. Call it
// on the UI thread.
func (self *UpdateHelper) recordResponse(version *types.VersionNumber) {
	self.menuShownAt.Store(0)
	self.c.OnWorkerBackground(func(gocui.Task) error {
		self.updater.RecordResponse(version)
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
	if status.InstalledVersion != nil {
		return style.FgGreen.Sprint(utils.ResolvePlaceholderString(
			self.c.Tr.RestartToUseNewVersion, map[string]string{"newVersion": status.InstalledVersion.String()},
		))
	}
	if status.AvailableVersion != nil {
		return style.FgGreen.Sprint(style.PrintHyperlink(self.c.Tr.UpdateAvailableNote, UpdateNoteLink))
	}
	return ""
}

// ShowUpdateMenu shows the menu for the update that the note in the
// bottom-right corner tells the user about
func (self *UpdateHelper) ShowUpdateMenu() error {
	status := self.updater.Status()
	if status.AvailableVersion == nil || status.InstalledVersion != nil {
		return nil
	}

	self.c.OnWorker(func(gocui.Task) error {
		menuOpts := self.updateMenuOptions(status.AvailableVersion)
		self.c.OnUIThread(func() error {
			return self.showUpdateMenu(menuOpts)
		})
		return nil
	})
	return nil
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

		if installedVersion := self.updater.Status().InstalledVersion; installedVersion != nil && !installedVersion.IsOlderThan(newVersion) {
			self.c.OnUIThread(func() error {
				self.c.Toast(self.updateInstalledMessage(installedVersion))
				return nil
			})
			return nil
		}

		menuOpts := self.updateMenuOptions(newVersion)
		self.c.OnUIThread(func() error {
			return self.showUpdateMenu(menuOpts)
		})
		return nil
	})
}

// updateMenuOptions returns the options for the menu that tells the user about
// a new version. Call it on a worker, because it accesses the file system.
func (self *UpdateHelper) updateMenuOptions(newVersion *types.VersionNumber) types.CreateMenuOptions {
	installation := self.updater.Installation()
	// The release notes are a link rather than a menu item, because choosing a
	// menu item closes the menu
	paragraphs := []string{utils.ResolvePlaceholderString(
		self.c.Tr.NewVersionAvailable,
		map[string]string{
			"newVersion":     newVersion.String(),
			"currentVersion": self.updater.CurrentVersion().String(),
			"releaseNotes":   style.PrintHyperlink(self.c.Tr.ReleaseNotes, constants.Links.Releases),
		},
	)}
	items := []*types.MenuItem{}

	switch installation.Kind {
	case updates.InstallKindRelease:
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

	case updates.InstallKindHomebrew, updates.InstallKindGoInstall, updates.InstallKindScoop,
		updates.InstallKindWinget, updates.InstallKindChocolatey:
		paragraphs = append(paragraphs, utils.ResolvePlaceholderString(
			self.c.Tr.UpdateWithPackageManager,
			map[string]string{
				"packageManager": installation.PackageManagerName(),
				"command":        installation.UpdateCommand(),
			},
		))
		if installation.Kind != updates.InstallKindGoInstall {
			paragraphs = append(paragraphs, self.c.Tr.PackageManagersCanBeSlow)
		}
		items = append(items, &types.MenuItem{
			Label: self.c.Tr.CopyUpdateCommand,
			OnPress: func() error {
				if err := self.c.OS().CopyToClipboard(installation.UpdateCommand()); err != nil {
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

	items = append(items,
		&types.MenuItem{
			Label: self.c.Tr.RemindMeTomorrow,
			// Recording the response below is all it takes
			OnPress: func() error { return nil },
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
	paragraphs = append(paragraphs, self.c.Tr.TurnOffUpdateReminders)

	// Whatever the user chooses counts as a response, and so does closing the
	// menu
	for _, item := range items {
		onPress := item.OnPress
		item.OnPress = func() error {
			self.recordResponse(newVersion)
			return onPress()
		}
	}

	return types.CreateMenuOptions{
		Title:      self.c.Tr.UpdateAvailableTitle,
		Prompt:     strings.Join(paragraphs, "\n\n"),
		Items:      items,
		HideCancel: true,
		OnCancel: func() error {
			self.recordResponse(newVersion)
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
