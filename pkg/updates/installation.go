package updates

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

// InstallKind says how lazygit was installed. This decides what we can offer
// the user when there's an update.
type InstallKind int

const (
	// A build without a release version, e.g. one built from a git checkout. We
	// don't check for updates for these.
	InstallKindDevBuild InstallKind = iota
	// Installed from one of our release archives. We can update it in place.
	InstallKindRelease
	InstallKindHomebrew
	InstallKindGoInstall
	InstallKindScoop
	InstallKindWinget
	InstallKindChocolatey
	// Anything else with a release version, e.g. a package of a Linux
	// distribution. We can only tell the user that there's an update.
	InstallKindOther
)

type Installation struct {
	Kind InstallKind
	// The path of the lazygit executable, with symlinks resolved, or "" if we
	// don't know it
	ExecutablePath string
}

// Folders of package managers that install our release archives unchanged.
// The build source of these is BuildSourceRelease, but we must not update them
// in place, or the package manager would report the wrong version.
var packageManagerFolders = []string{
	"/mise/installs/",
	"/.asdf/installs/",
	"/aquaproj-aqua/",
	"/nix/store/",
	"/snap/",
}

func DetectInstallation(buildSource string, version string, executablePath string) Installation {
	kind := detectInstallKind(buildSource, version, executablePath)
	return Installation{Kind: kind, ExecutablePath: executablePath}
}

func detectInstallKind(buildSource string, version string, executablePath string) InstallKind {
	if _, err := types.ParseVersionNumber(version); err != nil {
		return InstallKindDevBuild
	}

	// Normalize the path so that we can match Windows paths on any platform,
	// and regardless of how they are capitalized.
	path := strings.ToLower(strings.ReplaceAll(executablePath, `\`, "/"))

	switch {
	case strings.EqualFold(buildSource, BuildSourceHomebrew) || strings.Contains(path, "/cellar/lazygit/"):
		return InstallKindHomebrew
	case buildSource == BuildSourceGoInstall:
		return InstallKindGoInstall
	case strings.Contains(path, "/scoop/apps/lazygit/"):
		return InstallKindScoop
	// For a single user, winget installs into %LOCALAPPDATA%\Microsoft\WinGet,
	// and for all users into %ProgramFiles%\WinGet
	case strings.Contains(path, "/winget/packages/jesseduffield.lazygit"):
		return InstallKindWinget
	case strings.Contains(path, "/chocolatey/lib/lazygit/"):
		return InstallKindChocolatey
	// We can't update the executable in place if we don't know where it is
	case buildSource == BuildSourceRelease && executablePath != "" && !isInPackageManagerFolder(path):
		return InstallKindRelease
	default:
		return InstallKindOther
	}
}

func isInPackageManagerFolder(path string) bool {
	// On Linux, /usr/bin belongs to the distribution's package manager
	if strings.HasPrefix(path, "/usr/bin/") {
		return true
	}

	return lo.SomeBy(packageManagerFolders, func(folder string) bool {
		return strings.Contains(path, folder)
	})
}

// PackageManagerName returns the name of the package manager that lazygit was
// installed with, or "" if we don't know it.
func (self Installation) PackageManagerName() string {
	switch self.Kind {
	case InstallKindHomebrew:
		return "Homebrew"
	case InstallKindGoInstall:
		return "go install"
	case InstallKindScoop:
		return "Scoop"
	case InstallKindWinget:
		return "winget"
	case InstallKindChocolatey:
		return "Chocolatey"
	default:
		return ""
	}
}

// UpdateCommand returns the command that updates lazygit with the package
// manager that it was installed with, or "" if we don't know it.
func (self Installation) UpdateCommand() string {
	switch self.Kind {
	case InstallKindHomebrew:
		return "brew upgrade lazygit"
	case InstallKindGoInstall:
		return "go install github.com/jesseduffield/lazygit@latest"
	case InstallKindScoop:
		return "scoop update lazygit"
	case InstallKindWinget:
		return "winget upgrade JesseDuffield.lazygit"
	case InstallKindChocolatey:
		return "choco upgrade lazygit"
	default:
		return ""
	}
}
