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
	return u.installNewExecutable(newPath)
}

func (u *Updater) installNewExecutable(newPath string) error {
	u.replaceMutex.Lock()
	defer u.replaceMutex.Unlock()

	// Another lazygit running from the same executable may be installing the
	// update at the same time. We replace the executable while holding the
	// lock of state.yml, so that only one of them does it; the other one finds
	// it replaced already. On Windows, where replacing takes two renames, two
	// lazygit instances could get in each other's way otherwise.
	locked := false
	var err error
	lockErr := u.appState.UpdateAppState(func(*config.AppState) {
		locked = true
		if u.IsUpdateInstalled() {
			_ = os.Remove(newPath)
			return
		}
		err = replaceExecutable(u.installation.ExecutablePath, newPath, u.goos)
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

// fileModTime returns the modification time of the file at path, or the zero
// time if it can't be determined
func fileModTime(path string) time.Time {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}
	}
	return info.ModTime()
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
