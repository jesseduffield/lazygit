package updates

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/constants"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// Updater checks for updates and does updates
type Updater struct {
	*common.Common
	Config         config.AppConfigurer
	OSCommand      *oscommands.OSCommand
	installation   Installation
	currentVersion *types.VersionNumber

	httpClient *http.Client
	repoURL    string
}

const checkTimeout = 30 * time.Second

// NewUpdater creates a new updater
func NewUpdater(cmn *common.Common, config config.AppConfigurer, osCommand *oscommands.OSCommand) (*Updater, error) {
	currentVersion, _ := types.ParseVersionNumber(config.GetVersion())

	return &Updater{
		Common:         cmn,
		Config:         config,
		OSCommand:      osCommand,
		installation:   DetectInstallation(config.GetBuildSource(), config.GetVersion(), executablePath()),
		currentVersion: currentVersion,
		httpClient:     &http.Client{},
		repoURL:        constants.Links.RepoUrl,
	}, nil
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

func (u *Updater) mappedOs(os string) string {
	osMap := map[string]string{
		"darwin":  "Darwin",
		"linux":   "Linux",
		"windows": "Windows",
	}
	result, found := osMap[os]
	if found {
		return result
	}
	return os
}

func (u *Updater) mappedArch(arch string) string {
	archMap := map[string]string{
		"386":   "32-bit",
		"amd64": "x86_64",
	}
	result, found := archMap[arch]
	if found {
		return result
	}
	return arch
}

func (u *Updater) zipExtension() string {
	if runtime.GOOS == "windows" {
		return "zip"
	}

	return "tar.gz"
}

// example: https://github.com/jesseduffield/lazygit/releases/download/v0.1.73/lazygit_0.1.73_Darwin_x86_64.tar.gz
func (u *Updater) getBinaryUrl(newVersion string) string {
	url := fmt.Sprintf(
		"%s/releases/download/%s/lazygit_%s_%s_%s.%s",
		constants.Links.RepoUrl,
		newVersion,
		newVersion[1:],
		u.mappedOs(runtime.GOOS),
		u.mappedArch(runtime.GOARCH),
		u.zipExtension(),
	)
	u.Log.Info("Url for latest release is " + url)
	return url
}

// Update downloads the latest binary and replaces the current binary with it
func (u *Updater) Update(newVersion string) error {
	return u.update(newVersion)
}

func (u *Updater) update(newVersion string) error {
	rawUrl := u.getBinaryUrl(newVersion)
	u.Log.Info("Updating with url " + rawUrl)
	return u.downloadAndInstall(rawUrl)
}

func (u *Updater) downloadAndInstall(rawUrl string) error {
	configDir := u.Config.GetUserConfigDir()
	u.Log.Info("Download directory is " + configDir)

	zipPath := filepath.Join(configDir, "temp_lazygit."+u.zipExtension())
	u.Log.Info("Temp path to tarball/zip file is " + zipPath)

	// remove existing zip file
	if err := os.RemoveAll(zipPath); err != nil && !os.IsNotExist(err) {
		return err
	}

	// Create the zip file
	out, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(rawUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error while trying to download latest lazygit: %s", resp.Status)
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	u.Log.Info("untarring tarball/unzipping zip file")
	err = u.OSCommand.Cmd.New([]string{"tar", "-zxf", zipPath, "lazygit"}).Run()
	if err != nil {
		return err
	}

	// the `tar` terminal cannot store things in a new location without permission
	// so it creates it in the current directory. As such our path is fairly simple.
	// You won't see it because it's gitignored.
	tempLazygitFilePath := "lazygit"

	u.Log.Infof("Path to temp binary is %s", tempLazygitFilePath)

	// get the path of the current binary
	binaryPath, err := os.Executable()
	if err != nil {
		return err
	}
	u.Log.Info("Binary path is " + binaryPath)

	// Verify the main file exists
	if _, err := os.Stat(zipPath); err != nil {
		return err
	}

	// swap out the old binary for the new one
	err = os.Rename(tempLazygitFilePath, binaryPath)
	if err != nil {
		return err
	}
	u.Log.Info("Update complete!")

	return nil
}
