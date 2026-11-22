package updates

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
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

func DetectInstallation(buildSource string, version string, executablePath string) Installation {
	kind := detectInstallKind(buildSource, version, executablePath)
	return Installation{Kind: kind, ExecutablePath: executablePath}
}

func detectInstallKind(buildSource string, version string, executablePath string) InstallKind {
	if _, err := types.ParseVersionNumber(version); err != nil {
		return InstallKindDevBuild
	}

	// We can't update the executable in place if we don't know where it is
	if buildSource == BuildSourceRelease && executablePath != "" {
		return InstallKindRelease
	}
	return InstallKindOther
}
