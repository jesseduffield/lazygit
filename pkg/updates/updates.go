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

	httpClient *http.Client
	repoURL    string
	goos       string
	goarch     string

	installing atomic.Bool

	statusMutex sync.Mutex
	status      Status

	// Held while we replace the executable, so that lazygit doesn't exit in
	// the middle of it
	replaceMutex sync.Mutex
}

func NewUpdater(cmn *common.Common, config config.AppConfigurer) *Updater {
	currentVersion, _ := types.ParseVersionNumber(config.GetVersion())

	return &Updater{
		Common:         cmn,
		appState:       config,
		installation:   DetectInstallation(config.GetBuildSource(), config.GetVersion(), executablePath()),
		currentVersion: currentVersion,
		httpClient:     &http.Client{},
		repoURL:        constants.Links.RepoUrl,
		goos:           runtime.GOOS,
		goarch:         runtime.GOARCH,
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
	// The version that an update installed in place of the running lazygit.
	// This lazygit, or another one running from the same executable, can have
	// installed it. The new version takes effect when lazygit is restarted.
	InstalledVersion *types.VersionNumber
}

// ChecksEnabled returns whether lazygit checks for updates regularly. It
// doesn't for dev builds.
func (u *Updater) ChecksEnabled() bool {
	return u.installation.Kind != InstallKindDevBuild
}

// CheckForUpdateIfDue checks for an update if lazygit checks regularly, and
// the last successful check is long enough ago.
func (u *Updater) CheckForUpdateIfDue() error {
	if !u.ChecksEnabled() {
		return nil
	}

	state := u.readAppState()
	interval := time.Duration(u.UserConfig().Update.Days) * 24 * time.Hour
	if time.Since(time.Unix(state.Updates.LastCheck, 0)) < interval {
		u.RefreshStatus()
		return nil
	}

	_, err := u.CheckForUpdate()
	return err
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
		state.Updates.LastCheck = time.Now().Unix()
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
// lazygit. Only call it for InstallKindRelease.
func (u *Updater) Install(version *types.VersionNumber) error {
	if !u.installing.CompareAndSwap(false, true) {
		return errors.New("an update is already in progress")
	}
	defer u.installing.Store(false)

	ctx, cancel := context.WithTimeout(context.Background(), installTimeout)
	defer cancel()

	url := assetURL(u.repoURL, version, u.goos, u.goarch)
	u.Log.Infof("Installing %s", url)

	binary, err := downloadBinary(ctx, u.httpClient, url, u.goos)
	if err != nil {
		return err
	}

	if err := u.replaceExecutable(binary); err != nil {
		return err
	}

	u.updateAppState(func(state *config.AppState) {
		state.Updates.InstalledVersion = version.String()
		state.Updates.InstalledPath = u.installation.ExecutablePath
	})
	u.RefreshStatus()
	return nil
}

func (u *Updater) replaceExecutable(binary []byte) error {
	u.replaceMutex.Lock()
	defer u.replaceMutex.Unlock()
	return replaceExecutable(u.installation.ExecutablePath, binary, u.goos)
}

// RefreshStatus reads the status from the app state, where other lazygit
// instances can have changed it, and returns it. Don't call it on the UI
// thread, because it accesses the file system.
func (u *Updater) RefreshStatus() Status {
	state := u.readAppState()

	status := Status{
		AvailableVersion: u.newerVersion(state.Updates.LatestVersion),
	}
	if state.Updates.InstalledPath == u.installation.ExecutablePath {
		status.InstalledVersion = u.newerVersion(state.Updates.InstalledVersion)
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
