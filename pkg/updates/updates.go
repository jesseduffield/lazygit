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

	err = u.appState.UpdateAppState(func(state *config.AppState) {
		state.Updates.InstalledVersion = version.String()
		state.Updates.InstalledPath = u.installation.ExecutablePath
	})
	if err != nil {
		// The update is installed; only other lazygit instances won't know
		u.Log.Errorf("Error saving the app state: %v", err)
	}
	return nil
}

func (u *Updater) replaceExecutable(binary []byte) error {
	u.replaceMutex.Lock()
	defer u.replaceMutex.Unlock()
	return replaceExecutable(u.installation.ExecutablePath, binary, u.goos)
}

// InstalledVersion returns the version that an update installed in place of
// the running lazygit, or nil if there's none. This lazygit, or another one
// running from the same executable, can have installed it. The new version
// takes effect when lazygit is restarted.
func (u *Updater) InstalledVersion() *types.VersionNumber {
	state, err := u.appState.ReadAppState()
	if err != nil {
		u.Log.Errorf("Error reading the app state: %v", err)
	}

	if state.Updates.InstalledPath != u.installation.ExecutablePath {
		return nil
	}
	version, err := types.ParseVersionNumber(state.Updates.InstalledVersion)
	if err != nil || !u.currentVersion.IsOlderThan(version) {
		return nil
	}
	return version
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
