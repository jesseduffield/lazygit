package updates

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/stretchr/testify/assert"
)

func TestAssetName(t *testing.T) {
	// The assets of the v0.66.2 release
	expectedNames := []string{
		"lazygit_0.66.2_darwin_arm64.tar.gz",
		"lazygit_0.66.2_darwin_x86_64.tar.gz",
		"lazygit_0.66.2_freebsd_32-bit.tar.gz",
		"lazygit_0.66.2_freebsd_arm64.tar.gz",
		"lazygit_0.66.2_freebsd_armv6.tar.gz",
		"lazygit_0.66.2_freebsd_x86_64.tar.gz",
		"lazygit_0.66.2_linux_32-bit.tar.gz",
		"lazygit_0.66.2_linux_arm64.tar.gz",
		"lazygit_0.66.2_linux_armv6.tar.gz",
		"lazygit_0.66.2_linux_x86_64.tar.gz",
		"lazygit_0.66.2_windows_32-bit.zip",
		"lazygit_0.66.2_windows_arm64.zip",
		"lazygit_0.66.2_windows_x86_64.zip",
	}

	version := &types.VersionNumber{Major: 0, Minor: 66, Patch: 2}
	names := []string{}
	for _, goos := range []string{"darwin", "freebsd", "linux", "windows"} {
		for _, goarch := range []string{"arm64", "amd64", "386", "arm"} {
			// goreleaser doesn't build these
			if (goos == "darwin" && (goarch == "386" || goarch == "arm")) || (goos == "windows" && goarch == "arm") {
				continue
			}
			names = append(names, assetName(version, goos, goarch))
		}
	}

	assert.ElementsMatch(t, expectedNames, names)
}

func TestAssetURL(t *testing.T) {
	version := &types.VersionNumber{Major: 0, Minor: 67, Patch: 0}
	assert.Equal(t,
		"https://github.com/jesseduffield/lazygit/releases/download/v0.67.0/lazygit_0.67.0_linux_x86_64.tar.gz",
		assetURL("https://github.com/jesseduffield/lazygit", version, "linux", "amd64"))
}
