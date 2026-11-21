package updates

import (
	"context"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/constants"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

const (
	checkTimeout   = 30 * time.Second
	installTimeout = 10 * time.Minute
)

// The part of config.AppConfigurer that the updater needs
type appStateStore interface {
	ReadAppState() (*config.AppState, error)
	UpdateAppState(update func(state *config.AppState)) error
}

// Updater checks for new lazygit releases, and installs them if lazygit was
// installed from one of our release archives.
type Updater struct {
	*common.Common
	appState       appStateStore
	installation   Installation
	currentVersion *types.VersionNumber
	// The modification time of the executable when lazygit started. Installing
	// an update writes a new file, so the time changes; see IsUpdateInstalled.
	executableModTime time.Time

	httpClient *http.Client
	repoURL    string
	goos       string
	goarch     string

	installing atomic.Bool

	statusMutex sync.Mutex
	status      Status

	// Replaced in tests
	now func() time.Time

	// Held while we replace the executable, so that lazygit doesn't exit in
	// the middle of it
	replaceMutex sync.Mutex
}

func NewUpdater(cmn *common.Common, config config.AppConfigurer) *Updater {
	currentVersion, _ := types.ParseVersionNumber(config.GetVersion())
	installation := DetectInstallation(config.GetBuildSource(), config.GetVersion(), executablePath())

	return &Updater{
		Common:            cmn,
		appState:          config,
		installation:      installation,
		currentVersion:    currentVersion,
		executableModTime: fileModTime(installation.ExecutablePath),
		httpClient:        &http.Client{},
		repoURL:           constants.Links.RepoUrl,
		goos:              runtime.GOOS,
		goarch:            runtime.GOARCH,
		now:               time.Now,
	}
}

func executablePath() string {
	path, err := os.Executable()
	if err != nil {
		return ""
	}
	if resolvedPath, err := filepath.EvalSymlinks(path); err == nil {
		return resolvedPath
	}
	return path
}

func (u *Updater) Installation() Installation {
	return u.installation
}

// CurrentVersion returns nil for dev builds
func (u *Updater) CurrentVersion() *types.VersionNumber {
	return u.currentVersion
}

// Status is what lazygit knows about updates
type Status struct {
	// The latest release that a check found, if it is newer than the running
	// lazygit
	AvailableVersion *types.VersionNumber
	// The version that was installed in place of the running lazygit since it
	// started, if any; see InstalledVersion
	InstalledVersion *types.VersionNumber
}

// ChecksEnabled returns whether lazygit checks for updates regularly. It
// doesn't for dev builds.
func (u *Updater) ChecksEnabled() bool {
	return u.installation.Kind != InstallKindDevBuild
}

// CheckForUpdateIfDue checks for an update if lazygit checks regularly, and
// the last successful check is long enough ago. update.days counts calendar
// days, so that with update.days set to 1, the check happens at the first start
// of each day, even if that is a bit earlier than on the day before.
func (u *Updater) CheckForUpdateIfDue() error {
	if !u.ChecksEnabled() {
		return nil
	}

	state := u.readAppState()
	now := u.now()
	lastCheck := time.Unix(state.Updates.LastCheck, 0)
	nextCheck := startOfDayAfter(lastCheck, int(u.UserConfig().Update.Days))
	// A last check in the future means that the clock was wrong at the time.
	// Check now, so that this can't keep lazygit from checking until then.
	if now.Before(nextCheck) && !lastCheck.After(now) {
		u.RefreshStatus()
		return nil
	}

	_, err := u.CheckForUpdate()
	return err
}

// startOfDayAfter returns the start of the day that is the given number of
// days after the day of t, in the time zone of t
func startOfDayAfter(t time.Time, days int) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day+days, 0, 0, 0, 0, t.Location())
}

// CheckForUpdate returns the latest version if it is newer than the current
// one, or nil if it isn't. Don't call it for dev builds.
func (u *Updater) CheckForUpdate() (*types.VersionNumber, error) {
	ctx, cancel := context.WithTimeout(context.Background(), checkTimeout)
	defer cancel()

	latestVersion, err := fetchLatestVersion(ctx, u.httpClient, u.repoURL)
	if err != nil {
		return nil, err
	}

	u.Log.Infof("Current version is %s, latest version is %s", u.currentVersion, latestVersion)
	u.updateAppState(func(state *config.AppState) {
		state.Updates.LastCheck = u.now().Unix()
		state.Updates.LatestVersion = latestVersion.String()
	})
	u.RefreshStatus()

	if !u.currentVersion.IsOlderThan(latestVersion) {
		return nil, nil
	}
	return latestVersion, nil
}

// CheckCanInstall returns an error if we can't install an update in place
// because we don't have write access to the folder of the executable. Only call
// it for InstallKindRelease.
func (u *Updater) CheckCanInstall() error {
	return checkCanReplaceExecutable(u.installation.ExecutablePath)
}

// Install downloads the given release and installs it in place of the running
// lazygit. If an update is installed already, it does nothing. Only call it for
// InstallKindRelease.
func (u *Updater) Install(version *types.VersionNumber) error {
	if !u.installing.CompareAndSwap(false, true) {
		return errors.New("an update is already in progress")
	}
	defer u.installing.Store(false)

	if u.IsUpdateInstalled() {
		u.RefreshStatus()
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), installTimeout)
	defer cancel()

	url := assetURL(u.repoURL, version, u.goos, u.goarch)
	u.Log.Infof("Installing %s", url)

	binary, err := downloadBinary(ctx, u.httpClient, url, u.goos)
	if err != nil {
		return err
	}

	newPath, err := writeNewExecutable(u.installation.ExecutablePath, binary)
	if err != nil {
		return err
	}
	if err := u.installNewExecutable(newPath, version); err != nil {
		return err
	}

	u.RefreshStatus()
	return nil
}

func (u *Updater) installNewExecutable(newPath string, version *types.VersionNumber) error {
	u.replaceMutex.Lock()
	defer u.replaceMutex.Unlock()

	// Another lazygit running from the same executable may be installing the
	// update at the same time. We replace the executable while holding the
	// lock of state.yml, so that only one of them does it; the other one finds
	// it replaced already. On Windows, where replacing takes two renames, two
	// lazygit instances could get in each other's way otherwise.
	locked := false
	var err error
	lockErr := u.appState.UpdateAppState(func(state *config.AppState) {
		locked = true
		if u.IsUpdateInstalled() {
			_ = os.Remove(newPath)
			return
		}
		err = replaceExecutable(u.installation.ExecutablePath, newPath, u.goos)
		if err == nil {
			state.Updates.InstalledVersion = version.String()
			state.Updates.InstalledModTime = fileModTime(u.installation.ExecutablePath).UnixNano()
		}
	})
	if !locked {
		_ = os.Remove(newPath)
		return lockErr
	}
	return err
}

// IsUpdateInstalled returns whether an update was installed in place of the
// running lazygit since it started. This lazygit, or another one running from
// the same executable, can have installed it. The new version takes effect when
// lazygit is restarted. Don't call it on the UI thread, because it accesses the
// file system.
func (u *Updater) IsUpdateInstalled() bool {
	if u.installation.Kind != InstallKindRelease || u.executableModTime.IsZero() {
		return false
	}
	modTime := fileModTime(u.installation.ExecutablePath)
	return !modTime.IsZero() && !modTime.Equal(u.executableModTime)
}

// InstalledVersion returns the version that was installed in place of the
// running lazygit since it started, or nil if none was; see IsUpdateInstalled.
// Don't call it on the UI thread, because it accesses the file system.
func (u *Updater) InstalledVersion() *types.VersionNumber {
	return u.installedVersion(u.readAppState())
}

func (u *Updater) installedVersion(state *config.AppState) *types.VersionNumber {
	if !u.IsUpdateInstalled() {
		return nil
	}

	if fileModTime(u.installation.ExecutablePath).UnixNano() == state.Updates.InstalledModTime {
		if version, err := types.ParseVersionNumber(state.Updates.InstalledVersion); err == nil {
			return version
		}
	}

	// Something other than lazygit replaced the executable, so we don't know
	// its version. Most likely it's the latest one.
	return u.newerVersion(state.Updates.LatestVersion)
}

// fileModTime returns the modification time of the file at path, or the zero
// time if it can't be determined
func fileModTime(path string) time.Time {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}
	}
	return info.ModTime()
}

// remindersEnabled returns whether lazygit reminds the user of available
// updates by showing the update menu at startup. With update.method set to
// never, it only shows the note in the bottom-right corner.
func (u *Updater) remindersEnabled() bool {
	return u.ChecksEnabled() && u.UserConfig().Update.Method != "never"
}

// DueReminder returns the version to remind the user of with the update menu,
// or nil if no reminder is due. Don't call it on the UI thread, because it
// accesses the file system.
func (u *Updater) DueReminder() *types.VersionNumber {
	if !u.remindersEnabled() {
		return nil
	}

	state := u.readAppState()
	version := u.newerVersion(state.Updates.LatestVersion)
	if version == nil || u.IsUpdateInstalled() ||
		state.Updates.SkippedVersion == state.Updates.LatestVersion ||
		u.reminderPostponed(state) {
		return nil
	}
	return version
}

// reminderPostponed returns whether the user asked to be reminded tomorrow.
// RemindTomorrow postpones the reminder to the start of the next day; a later
// time means that the clock was wrong when the user responded. We ignore it, so
// that it can't keep lazygit from reminding the user until then.
func (u *Updater) reminderPostponed(state *config.AppState) bool {
	now := u.now()
	remindAfter := state.Updates.RemindAfter
	return now.Unix() < remindAfter && remindAfter <= startOfDayAfter(now, 1).Unix()
}

// RemindTomorrow keeps lazygit from showing the update menu at startup until
// the next day. It also undoes skipping a version; if the user responded to the
// menu in several lazygit instances, the last response counts.
func (u *Updater) RemindTomorrow() {
	remindAfter := startOfDayAfter(u.now(), 1).Unix()
	u.updateAppState(func(state *config.AppState) {
		state.Updates.RemindAfter = remindAfter
		state.Updates.SkippedVersion = ""
	})
}

// SkipVersion keeps lazygit from reminding the user of the given version
func (u *Updater) SkipVersion(version *types.VersionNumber) {
	u.updateAppState(func(state *config.AppState) {
		state.Updates.SkippedVersion = version.String()
	})
}

// RefreshStatus reads the status from the app state, where other lazygit
// instances can have changed it, and returns it. Don't call it on the UI
// thread, because it accesses the file system.
func (u *Updater) RefreshStatus() Status {
	state := u.readAppState()

	status := Status{
		AvailableVersion: u.newerVersion(state.Updates.LatestVersion),
		InstalledVersion: u.installedVersion(state),
	}

	u.statusMutex.Lock()
	defer u.statusMutex.Unlock()
	u.status = status
	return status
}

// Status returns the status as of the last refresh
func (u *Updater) Status() Status {
	u.statusMutex.Lock()
	defer u.statusMutex.Unlock()
	return u.status
}

// newerVersion parses the given version and returns it if it is newer than the
// running lazygit, or nil otherwise
func (u *Updater) newerVersion(versionStr string) *types.VersionNumber {
	if u.currentVersion == nil {
		return nil
	}
	version, err := types.ParseVersionNumber(versionStr)
	if err != nil || !u.currentVersion.IsOlderThan(version) {
		return nil
	}
	return version
}

func (u *Updater) readAppState() *config.AppState {
	state, err := u.appState.ReadAppState()
	if err != nil {
		u.Log.Errorf("Error reading the app state: %v", err)
	}
	return state
}

// Errors from saving the update state aren't worth bothering the user with; at
// worst, we check or remind again.
func (u *Updater) updateAppState(update func(state *config.AppState)) {
	if err := u.appState.UpdateAppState(update); err != nil {
		u.Log.Errorf("Error saving the app state: %v", err)
	}
}

func (u *Updater) IsInstalling() bool {
	return u.installing.Load()
}

// WaitForInstallToFinish waits until a running install is done replacing the
// executable, and keeps later installs from starting. Call it before exiting.
func (u *Updater) WaitForInstallToFinish() {
	u.replaceMutex.Lock()
}

// CleanUpAfterUpdate removes files that an earlier update left behind
func (u *Updater) CleanUpAfterUpdate() {
	if u.installation.Kind == InstallKindRelease {
		cleanUpAfterUpdate(u.installation.ExecutablePath)
	}
}
