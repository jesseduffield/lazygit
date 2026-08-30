package xdg

import (
	"path/filepath"

	"github.com/adrg/xdg/internal/pathutil"
	"github.com/adrg/xdg/internal/userdirs"
	"golang.org/x/sys/windows"
)

func initDirs(home string) {
	kf := initKnownFolders(home)
	initBaseDirs(home, kf)
	initUserDirs(home, kf)
}

func initBaseDirs(home string, kf *knownFolders) {
	// Initialize standard directories.
	baseDirs.dataHome = pathutil.EnvPath(envDataHome, kf.localAppData)
	baseDirs.data = pathutil.EnvPathList(envDataDirs, kf.roamingAppData, kf.programData)
	baseDirs.configHome = pathutil.EnvPath(envConfigHome, kf.localAppData)
	baseDirs.config = pathutil.EnvPathList(envConfigDirs, kf.programData, kf.roamingAppData)
	baseDirs.stateHome = pathutil.EnvPath(envStateHome, kf.localAppData)
	baseDirs.cacheHome = pathutil.EnvPath(envCacheHome, filepath.Join(kf.localAppData, "cache"))
	baseDirs.runtime = pathutil.EnvPath(envRuntimeDir, kf.localAppData)

	// Initialize non-standard directories.
	baseDirs.binHome = pathutil.EnvPath(envBinHome, kf.userProgramFiles)

	baseDirs.applications = []string{
		kf.programs,
		kf.commonPrograms,
		kf.programFiles,
		kf.programFilesCommon,
		kf.userProgramFiles,
		kf.userProgramFilesCommon,
	}

	baseDirs.fonts = []string{
		kf.fonts,
		filepath.Join(kf.localAppData, "Microsoft", "Windows", "Fonts"),
	}
}

func initUserDirs(home string, kf *knownFolders) {
	UserDirs.Desktop = pathutil.EnvPath(userdirs.EnvDesktopDir, kf.desktop)
	UserDirs.Download = pathutil.EnvPath(userdirs.EnvDownloadDir, kf.downloads)
	UserDirs.Documents = pathutil.EnvPath(userdirs.EnvDocumentsDir, kf.documents)
	UserDirs.Music = pathutil.EnvPath(userdirs.EnvMusicDir, kf.music)
	UserDirs.Pictures = pathutil.EnvPath(userdirs.EnvPicturesDir, kf.pictures)
	UserDirs.Videos = pathutil.EnvPath(userdirs.EnvVideosDir, kf.videos)
	UserDirs.Templates = pathutil.EnvPath(userdirs.EnvTemplatesDir, kf.templates)
	UserDirs.PublicShare = pathutil.EnvPath(userdirs.EnvPublicShareDir, kf.public)
}

type knownFolders struct {
	systemDrive            string
	systemRoot             string
	programData            string
	userProfile            string
	userProfiles           string
	roamingAppData         string
	localAppData           string
	desktop                string
	downloads              string
	documents              string
	music                  string
	pictures               string
	videos                 string
	templates              string
	public                 string
	fonts                  string
	programs               string
	commonPrograms         string
	programFiles           string
	programFilesCommon     string
	userProgramFiles       string
	userProgramFilesCommon string
}

func initKnownFolders(home string) *knownFolders {
	kf := &knownFolders{
		userProfile: home,
	}
	kf.systemDrive = filepath.VolumeName(pathutil.KnownFolder(
		windows.FOLDERID_Windows,
		[]string{"SystemDrive", "SystemRoot", "windir"},
		[]string{home, `C:`},
	)) + string(filepath.Separator)
	kf.systemRoot = pathutil.KnownFolder(
		windows.FOLDERID_Windows,
		[]string{"SystemRoot", "windir"},
		[]string{filepath.Join(kf.systemDrive, "Windows")},
	)
	kf.programData = pathutil.KnownFolder(
		windows.FOLDERID_ProgramData,
		[]string{"ProgramData", "ALLUSERSPROFILE"},
		[]string{filepath.Join(kf.systemDrive, "ProgramData")},
	)
	kf.userProfiles = pathutil.KnownFolder(
		windows.FOLDERID_UserProfiles,
		nil,
		[]string{filepath.Join(kf.systemDrive, "Users")},
	)
	kf.roamingAppData = pathutil.KnownFolder(
		windows.FOLDERID_RoamingAppData,
		[]string{"APPDATA"},
		[]string{filepath.Join(home, "AppData", "Roaming")},
	)
	kf.localAppData = pathutil.KnownFolder(
		windows.FOLDERID_LocalAppData,
		[]string{"LOCALAPPDATA"},
		[]string{filepath.Join(home, "AppData", "Local")},
	)
	kf.desktop = pathutil.KnownFolder(
		windows.FOLDERID_Desktop,
		nil,
		[]string{filepath.Join(home, "Desktop")},
	)
	kf.downloads = pathutil.KnownFolder(
		windows.FOLDERID_Downloads,
		nil,
		[]string{filepath.Join(home, "Downloads")},
	)
	kf.documents = pathutil.KnownFolder(
		windows.FOLDERID_Documents,
		nil,
		[]string{filepath.Join(home, "Documents")},
	)
	kf.music = pathutil.KnownFolder(
		windows.FOLDERID_Music,
		nil,
		[]string{filepath.Join(home, "Music")},
	)
	kf.pictures = pathutil.KnownFolder(
		windows.FOLDERID_Pictures,
		nil,
		[]string{filepath.Join(home, "Pictures")},
	)
	kf.videos = pathutil.KnownFolder(
		windows.FOLDERID_Videos,
		nil,
		[]string{filepath.Join(home, "Videos")},
	)
	kf.templates = pathutil.KnownFolder(
		windows.FOLDERID_Templates,
		nil,
		[]string{filepath.Join(kf.roamingAppData, "Microsoft", "Windows", "Templates")},
	)
	kf.public = pathutil.KnownFolder(
		windows.FOLDERID_Public,
		[]string{"PUBLIC"},
		[]string{filepath.Join(kf.userProfiles, "Public")},
	)
	kf.fonts = pathutil.KnownFolder(
		windows.FOLDERID_Fonts,
		nil,
		[]string{filepath.Join(kf.systemRoot, "Fonts")},
	)
	kf.programs = pathutil.KnownFolder(
		windows.FOLDERID_Programs,
		nil,
		[]string{filepath.Join(kf.roamingAppData, "Microsoft", "Windows", "Start Menu", "Programs")},
	)
	kf.commonPrograms = pathutil.KnownFolder(
		windows.FOLDERID_CommonPrograms,
		nil,
		[]string{filepath.Join(kf.programData, "Microsoft", "Windows", "Start Menu", "Programs")},
	)
	kf.programFiles = pathutil.KnownFolder(
		windows.FOLDERID_ProgramFiles,
		[]string{"ProgramFiles"},
		[]string{filepath.Join(kf.systemDrive, "Program Files")},
	)
	kf.programFilesCommon = pathutil.KnownFolder(
		windows.FOLDERID_ProgramFilesCommon,
		nil,
		[]string{filepath.Join(kf.programFiles, "Common Files")},
	)
	kf.userProgramFiles = pathutil.KnownFolder(
		windows.FOLDERID_UserProgramFiles,
		nil,
		[]string{
			filepath.Join(kf.localAppData, "Programs"),
		},
	)
	kf.userProgramFilesCommon = pathutil.KnownFolder(
		windows.FOLDERID_UserProgramFilesCommon,
		nil,
		[]string{
			filepath.Join(kf.userProgramFiles, "Common"),
		},
	)

	return kf
}
