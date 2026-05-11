//go:build aix || dragonfly || freebsd || (js && wasm) || nacl || linux || netbsd || openbsd || solaris

package xdg

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/adrg/xdg/internal/pathutil"
	"github.com/adrg/xdg/internal/userdirs"
)

func initDirs(home string) {
	initBaseDirs(home)
	initUserDirs(home, baseDirs.configHome)
}

func initBaseDirs(home string) {
	// Initialize standard directories.
	baseDirs.dataHome = pathutil.EnvPath(envDataHome, filepath.Join(home, ".local", "share"))
	baseDirs.data = pathutil.EnvPathList(envDataDirs, "/usr/local/share", "/usr/share")
	baseDirs.configHome = pathutil.EnvPath(envConfigHome, filepath.Join(home, ".config"))
	baseDirs.config = pathutil.EnvPathList(envConfigDirs, "/etc/xdg")
	baseDirs.stateHome = pathutil.EnvPath(envStateHome, filepath.Join(home, ".local", "state"))
	baseDirs.cacheHome = pathutil.EnvPath(envCacheHome, filepath.Join(home, ".cache"))
	baseDirs.runtime = pathutil.EnvPath(envRuntimeDir, filepath.Join("/run/user", strconv.Itoa(os.Getuid())))

	// Initialize non-standard directories.
	baseDirs.binHome = pathutil.EnvPath(envBinHome, filepath.Join(home, ".local", "bin"))

	appDirs := []string{
		filepath.Join(baseDirs.dataHome, "applications"),
		filepath.Join(home, ".local/share/applications"),
		"/usr/local/share/applications",
		"/usr/share/applications",
	}

	fontDirs := []string{
		filepath.Join(baseDirs.dataHome, "fonts"),
		filepath.Join(home, ".fonts"),
		filepath.Join(home, ".local/share/fonts"),
		"/usr/local/share/fonts",
		"/usr/share/fonts",
	}

	for _, dir := range baseDirs.data {
		appDirs = append(appDirs, filepath.Join(dir, "applications"))
		fontDirs = append(fontDirs, filepath.Join(dir, "fonts"))
	}

	baseDirs.applications = pathutil.Unique(appDirs)
	baseDirs.fonts = pathutil.Unique(fontDirs)
}

func initUserDirs(home, configHome string) {
	dirs, err := userdirs.ParseConfigFile(filepath.Join(configHome, "user-dirs.dirs"))
	if err != nil {
		dirs = &UserDirectories{}
	}

	UserDirs.Desktop = pathutil.EnvPath(userdirs.EnvDesktopDir, dirs.Desktop, filepath.Join(home, "Desktop"))
	UserDirs.Download = pathutil.EnvPath(userdirs.EnvDownloadDir, dirs.Download, filepath.Join(home, "Downloads"))
	UserDirs.Documents = pathutil.EnvPath(userdirs.EnvDocumentsDir, dirs.Documents, filepath.Join(home, "Documents"))
	UserDirs.Music = pathutil.EnvPath(userdirs.EnvMusicDir, dirs.Music, filepath.Join(home, "Music"))
	UserDirs.Pictures = pathutil.EnvPath(userdirs.EnvPicturesDir, dirs.Pictures, filepath.Join(home, "Pictures"))
	UserDirs.Videos = pathutil.EnvPath(userdirs.EnvVideosDir, dirs.Videos, filepath.Join(home, "Videos"))
	UserDirs.Templates = pathutil.EnvPath(userdirs.EnvTemplatesDir, dirs.Templates, filepath.Join(home, "Templates"))
	UserDirs.PublicShare = pathutil.EnvPath(userdirs.EnvPublicShareDir, dirs.PublicShare, filepath.Join(home, "Public"))
}
