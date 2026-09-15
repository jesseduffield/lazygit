package xdg

import (
	"path/filepath"

	"github.com/adrg/xdg/internal/pathutil"
	"github.com/adrg/xdg/internal/userdirs"
)

func initDirs(home string) {
	initBaseDirs(home)
	initUserDirs(home)
}

func initBaseDirs(home string) {
	homeAppSupport := filepath.Join(home, "Library", "Application Support")
	rootAppSupport := "/Library/Application Support"

	// Initialize standard directories.
	baseDirs.dataHome = pathutil.EnvPath(envDataHome, homeAppSupport)
	baseDirs.data = pathutil.EnvPathList(envDataDirs,
		rootAppSupport,
		filepath.Join(home, ".local", "share"),
	)
	baseDirs.configHome = pathutil.EnvPath(envConfigHome, homeAppSupport)
	baseDirs.config = pathutil.EnvPathList(envConfigDirs,
		filepath.Join(home, "Library", "Preferences"),
		rootAppSupport,
		"/Library/Preferences",
		filepath.Join(home, ".config"),
	)
	baseDirs.stateHome = pathutil.EnvPath(envStateHome, homeAppSupport)
	baseDirs.cacheHome = pathutil.EnvPath(envCacheHome, filepath.Join(home, "Library", "Caches"))
	baseDirs.runtime = pathutil.EnvPath(envRuntimeDir, homeAppSupport)

	// Initialize non-standard directories.
	baseDirs.binHome = pathutil.EnvPath(envBinHome, filepath.Join(home, ".local", "bin"))

	baseDirs.applications = []string{
		"/Applications",
	}

	baseDirs.fonts = []string{
		filepath.Join(home, "Library/Fonts"),
		"/Library/Fonts",
		"/System/Library/Fonts",
		"/Network/Library/Fonts",
	}
}

func initUserDirs(home string) {
	UserDirs.Desktop = pathutil.EnvPath(userdirs.EnvDesktopDir, filepath.Join(home, "Desktop"))
	UserDirs.Download = pathutil.EnvPath(userdirs.EnvDownloadDir, filepath.Join(home, "Downloads"))
	UserDirs.Documents = pathutil.EnvPath(userdirs.EnvDocumentsDir, filepath.Join(home, "Documents"))
	UserDirs.Music = pathutil.EnvPath(userdirs.EnvMusicDir, filepath.Join(home, "Music"))
	UserDirs.Pictures = pathutil.EnvPath(userdirs.EnvPicturesDir, filepath.Join(home, "Pictures"))
	UserDirs.Videos = pathutil.EnvPath(userdirs.EnvVideosDir, filepath.Join(home, "Movies"))
	UserDirs.Templates = pathutil.EnvPath(userdirs.EnvTemplatesDir, filepath.Join(home, "Templates"))
	UserDirs.PublicShare = pathutil.EnvPath(userdirs.EnvPublicShareDir, filepath.Join(home, "Public"))
}
