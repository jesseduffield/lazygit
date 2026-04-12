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
	homeLibDir := filepath.Join(home, "lib")
	rootLibDir := "/lib"

	// Initialize standard directories.
	baseDirs.dataHome = pathutil.EnvPath(envDataHome, homeLibDir)
	baseDirs.data = pathutil.EnvPathList(envDataDirs, rootLibDir)
	baseDirs.configHome = pathutil.EnvPath(envConfigHome, homeLibDir)
	baseDirs.config = pathutil.EnvPathList(envConfigDirs, rootLibDir)
	baseDirs.stateHome = pathutil.EnvPath(envStateHome, filepath.Join(homeLibDir, "state"))
	baseDirs.cacheHome = pathutil.EnvPath(envCacheHome, filepath.Join(homeLibDir, "cache"))
	baseDirs.runtime = pathutil.EnvPath(envRuntimeDir, "/tmp")

	// Initialize non-standard directories.
	baseDirs.binHome = pathutil.EnvPath(envBinHome, filepath.Join(home, "bin"))

	baseDirs.applications = []string{
		filepath.Join(home, "bin"),
		"/bin",
	}

	baseDirs.fonts = []string{
		filepath.Join(homeLibDir, "font"),
		"/lib/font",
	}
}

func initUserDirs(home string) {
	UserDirs.Desktop = pathutil.EnvPath(userdirs.EnvDesktopDir, filepath.Join(home, "desktop"))
	UserDirs.Download = pathutil.EnvPath(userdirs.EnvDownloadDir, filepath.Join(home, "downloads"))
	UserDirs.Documents = pathutil.EnvPath(userdirs.EnvDocumentsDir, filepath.Join(home, "documents"))
	UserDirs.Music = pathutil.EnvPath(userdirs.EnvMusicDir, filepath.Join(home, "music"))
	UserDirs.Pictures = pathutil.EnvPath(userdirs.EnvPicturesDir, filepath.Join(home, "pictures"))
	UserDirs.Videos = pathutil.EnvPath(userdirs.EnvVideosDir, filepath.Join(home, "videos"))
	UserDirs.Templates = pathutil.EnvPath(userdirs.EnvTemplatesDir, filepath.Join(home, "templates"))
	UserDirs.PublicShare = pathutil.EnvPath(userdirs.EnvPublicShareDir, filepath.Join(home, "public"))
}
