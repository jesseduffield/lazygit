package updates

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// assetName returns the name of the release archive for the given platform,
// following the name_template in .goreleaser.yml.
func assetName(version *types.VersionNumber, goos string, goarch string) string {
	arch := goarch
	switch goarch {
	case "amd64":
		arch = "x86_64"
	case "386":
		arch = "32-bit"
	case "arm":
		arch = "armv6"
	}

	extension := "tar.gz"
	if goos == "windows" {
		extension = "zip"
	}

	return fmt.Sprintf("lazygit_%d.%d.%d_%s_%s.%s", version.Major, version.Minor, version.Patch, goos, arch, extension)
}

func assetURL(repoURL string, version *types.VersionNumber, goos string, goarch string) string {
	return fmt.Sprintf("%s/releases/download/%s/%s", repoURL, version.String(), assetName(version, goos, goarch))
}
