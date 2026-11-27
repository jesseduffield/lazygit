package updates

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/go-errors/errors"

	"github.com/kardianos/osext"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/constants"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// Updater checks for updates and does updates
type Updater struct {
	*common.Common
	Config    config.AppConfigurer
	OSCommand *oscommands.OSCommand
}

// Updaterer implements the check and update methods
type Updaterer interface {
	CheckForNewUpdate()
	Update()
}

// NewUpdater creates a new updater
func NewUpdater(cmn *common.Common, config config.AppConfigurer, osCommand *oscommands.OSCommand) (*Updater, error) {
	return &Updater{
		Common:    cmn,
		Config:    config,
		OSCommand: osCommand,
	}, nil
}

func (u *Updater) getLatestVersionNumber() (string, error) {
	req, err := http.NewRequest("GET", constants.Links.RepoUrl+"/releases/latest", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	data := struct {
		TagName string `json:"tag_name"`
	}{}
	if err := dec.Decode(&data); err != nil {
		return "", err
	}

	return data.TagName, nil
}

// RecordLastUpdateCheck records last time an update check was performed
func (u *Updater) RecordLastUpdateCheck() error {
	return u.Config.UpdateAppState(func(state *config.AppState) {
		state.LastUpdateCheck = time.Now().Unix()
	})
}

// expecting version to be of the form `v12.34.56`
func (u *Updater) majorVersionDiffers(oldVersion, newVersion string) bool {
	if oldVersion == "unversioned" {
		return false
	}
	oldVersion = strings.TrimPrefix(oldVersion, "v")
	newVersion = strings.TrimPrefix(newVersion, "v")
	return strings.Split(oldVersion, ".")[0] != strings.Split(newVersion, ".")[0]
}

func (u *Updater) currentVersion() string {
	versionNumber := u.Config.GetVersion()
	if versionNumber == "unversioned" {
		return versionNumber
	}

	return fmt.Sprintf("v%s", u.Config.GetVersion())
}

func (u *Updater) checkForNewUpdate() (string, error) {
	u.Log.Info("Checking for an updated version")
	currentVersion := u.currentVersion()
	if err := u.RecordLastUpdateCheck(); err != nil {
		return "", err
	}

	newVersion, err := u.getLatestVersionNumber()
	if err != nil {
		return "", err
	}
	u.Log.Info("Current version is " + currentVersion)
	u.Log.Info("New version is " + newVersion)

	if newVersion == currentVersion {
		return "", errors.New(u.Tr.OnLatestVersionErr)
	}

	if u.majorVersionDiffers(currentVersion, newVersion) {
		errMessage := utils.ResolvePlaceholderString(
			u.Tr.MajorVersionErr, map[string]string{
				"newVersion":     newVersion,
				"currentVersion": currentVersion,
			},
		)
		return "", errors.New(errMessage)
	}

	rawUrl := u.getBinaryUrl(newVersion)

	u.Log.Info("Checking for resource at url " + rawUrl)
	if !u.verifyResourceFound(rawUrl) {
		errMessage := utils.ResolvePlaceholderString(
			u.Tr.CouldNotFindBinaryErr, map[string]string{
				"url": rawUrl,
			},
		)

		return "", errors.New(errMessage)
	}
	u.Log.Info("Verified resource is available, ready to update")

	return newVersion, nil
}

// CheckForNewUpdate checks if there is an available update
func (u *Updater) CheckForNewUpdate(onFinish func(string, error) error, userRequested bool) {
	if !userRequested && u.skipUpdateCheck() {
		return
	}

	newVersion, err := u.checkForNewUpdate()
	if err = onFinish(newVersion, err); err != nil {
		u.Log.Error(err)
	}
}

func (u *Updater) skipUpdateCheck() bool {
	// will remove the check for windows after adding a manifest file asking for
	// the required permissions
	if runtime.GOOS == "windows" {
		u.Log.Info("Updating is currently not supported for windows until we can fix permission issues")
		return true
	}

	if u.Config.GetVersion() == "unversioned" {
		u.Log.Info("Current version is not built from an official release so we won't check for an update")
		return true
	}

	if u.Config.GetBuildSource() != "buildBinary" {
		u.Log.Info("Binary is not built with the buildBinary flag so we won't check for an update")
		return true
	}

	userConfig := u.UserConfig()
	if userConfig.Update.Method == "never" {
		u.Log.Info("Update method is set to never so we won't check for an update")
		return true
	}

	currentTimestamp := time.Now().Unix()
	appState, err := u.Config.ReadAppState()
	if err != nil {
		u.Log.Errorf("error when reading app state: %v", err)
	}
	lastUpdateCheck := appState.LastUpdateCheck
	days := userConfig.Update.Days

	if (currentTimestamp-lastUpdateCheck)/(60*60*24) < days {
		u.Log.Info("Last update was too recent so we won't check for an update")
		return true
	}

	return false
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
	// One temp dir for download + extraction, cleaned up on exit
	tempDir, err := os.MkdirTemp("", "lazygit-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)
	u.Log.Info("Temp directory is " + tempDir)

	zipPath := filepath.Join(tempDir, "temp_lazygit."+u.zipExtension())
	u.Log.Info("Temp path to tarball/zip file is " + zipPath)

	if err := u.download(rawUrl, zipPath); err != nil {
		return err
	}

	u.Log.Info("untarring tarball/unzipping zip file")
	err = u.OSCommand.Cmd.New([]string{"tar", "-zxf", zipPath, "-C", tempDir, "lazygit"}).Run()
	if err != nil {
		return err
	}
	tempLazygitFilePath := filepath.Join(tempDir, "lazygit")
	u.Log.Infof("Path to temp binary is %s", tempLazygitFilePath)

	binaryPath, err := osext.Executable()
	if err != nil {
		return err
	}
	u.Log.Info("Binary path is " + binaryPath)

	// swap out the old binary for the new one
	if err := moveWithFallback(tempLazygitFilePath, binaryPath); err != nil {
		return err
	}
	u.Log.Info("Update complete!")
	return nil
}

func (u *Updater) download(rawUrl, dstPath string) error {
	client := &http.Client{Timeout: 3 * time.Minute}
	resp, err := client.Get(rawUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error while trying to download latest lazygit: %s", resp.Status)
	}

	out, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, resp.Body); err != nil {
		out.Close()
		return err
	}
	// Close explicitly so flush errors surface and the file is ready for tar
	return out.Close()
}

func moveWithFallback(src, dst string) error {
	err := os.Rename(src, dst)
	if err == nil {
		return nil
	}

	// Only fall back for cross-device moves; surface any other error as-is.
	// Windows reports cross-drive renames with a different errno, so there we
	// fall back on any LinkError.
	var linkErr *os.LinkError
	if !errors.As(err, &linkErr) {
		return err
	}
	if runtime.GOOS != "windows" && !errors.Is(linkErr.Err, syscall.EXDEV) {
		return err
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	info, err := srcFile.Stat()
	if err != nil {
		return err
	}

	// Temp file next to dst so the final rename is same-filesystem and atomic
	tmp, err := os.CreateTemp(filepath.Dir(dst), ".lazygit-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	cleanup := func() { os.Remove(tmpName) }

	if _, err := io.Copy(tmp, srcFile); err != nil {
		tmp.Close()
		cleanup()
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		cleanup()
		return err
	}
	if err := tmp.Close(); err != nil {
		cleanup()
		return err
	}
	// CreateTemp uses 0600, so restore the original (executable) permissions
	if err := os.Chmod(tmpName, info.Mode().Perm()); err != nil {
		cleanup()
		return err
	}
	if err := os.Rename(tmpName, dst); err != nil {
		cleanup()
		return err
	}

	// release the handle before removing (matters on Windows)
	srcFile.Close()
	return os.Remove(src)
}

func (u *Updater) verifyResourceFound(rawUrl string) bool {
	resp, err := http.Head(rawUrl)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	u.Log.Info("Received status code ", resp.StatusCode)
	// OK (200) indicates that the resource is present.
	return resp.StatusCode == http.StatusOK
}
