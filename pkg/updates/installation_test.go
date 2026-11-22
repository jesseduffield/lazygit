package updates

import (
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectInstallKind(t *testing.T) {
	scenarios := []struct {
		name           string
		buildSource    string
		version        string
		executablePath string
		expected       InstallKind
	}{
		{
			name:           "built from a git checkout",
			buildSource:    "unknown",
			version:        "0123abcd",
			executablePath: "/home/user/go/bin/lazygit",
			expected:       InstallKindDevBuild,
		},
		{
			name:           "built from a git checkout based on a release",
			buildSource:    "unknown",
			version:        "v0.66.2+0123abcd-dirty",
			executablePath: "/home/user/go/bin/lazygit",
			expected:       InstallKindDevBuild,
		},
		{
			name:           "built without version",
			buildSource:    "unknown",
			version:        "unversioned",
			executablePath: "/home/user/lazygit/lazygit",
			expected:       InstallKindDevBuild,
		},
		{
			name:           "dev build with release build source",
			buildSource:    "binaryRelease",
			version:        "0.66.2-next",
			executablePath: "/home/user/.local/bin/lazygit",
			expected:       InstallKindDevBuild,
		},
		{
			name:           "release archive",
			buildSource:    "binaryRelease",
			version:        "0.66.2",
			executablePath: "/home/user/.local/bin/lazygit",
			expected:       InstallKindRelease,
		},
		{
			name:           "release archive on Windows",
			buildSource:    "binaryRelease",
			version:        "0.66.2",
			executablePath: `C:\Users\user\AppData\Local\Programs\lazygit\lazygit.exe`,
			expected:       InstallKindRelease,
		},
		{
			name:           "release archive in /usr/local/bin",
			buildSource:    "binaryRelease",
			version:        "0.66.2",
			executablePath: "/usr/local/bin/lazygit",
			expected:       InstallKindRelease,
		},
		{
			name:           "release archive in an unknown folder",
			buildSource:    "binaryRelease",
			version:        "0.66.2",
			executablePath: "",
			expected:       InstallKindOther,
		},
		{
			name:           "distribution package",
			buildSource:    "unknown",
			version:        "0.66.2",
			executablePath: "/usr/bin/lazygit",
			expected:       InstallKindOther,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			installation := DetectInstallation(s.buildSource, s.version, s.executablePath)
			assert.Equal(t, s.expected, installation.Kind)
			assert.Equal(t, s.executablePath, installation.ExecutablePath)
		})
	}
}

func TestReleaseBuildSourceMatchesGoreleaserConfig(t *testing.T) {
	content, err := os.ReadFile("../../.goreleaser.yml")
	if !assert.NoError(t, err) {
		return
	}

	matches := regexp.MustCompile(`-X main\.buildSource=(\S+)`).FindAllStringSubmatch(string(content), -1)
	if assert.Len(t, matches, 1) {
		assert.Equal(t, BuildSourceRelease, matches[0][1])
	}
}
