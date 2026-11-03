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
			name:           "release archive in /usr/bin",
			buildSource:    "binaryRelease",
			version:        "0.66.2",
			executablePath: "/usr/bin/lazygit",
			expected:       InstallKindOther,
		},
		{
			name:           "release archive installed by mise",
			buildSource:    "binaryRelease",
			version:        "0.66.2",
			executablePath: "/home/user/.local/share/mise/installs/lazygit/0.66.2/lazygit",
			expected:       InstallKindOther,
		},
		{
			name:           "Homebrew",
			buildSource:    "Homebrew",
			version:        "0.66.2",
			executablePath: "/opt/homebrew/Cellar/lazygit/0.66.2/bin/lazygit",
			expected:       InstallKindHomebrew,
		},
		{
			name:           "Homebrew on Linux",
			buildSource:    "Homebrew",
			version:        "0.66.2",
			executablePath: "/home/linuxbrew/.linuxbrew/Cellar/lazygit/0.66.2/bin/lazygit",
			expected:       InstallKindHomebrew,
		},
		{
			name:           "Homebrew formula with a different build source",
			buildSource:    "someone",
			version:        "0.66.2",
			executablePath: "/usr/local/Cellar/lazygit/0.66.2/bin/lazygit",
			expected:       InstallKindHomebrew,
		},
		{
			name:           "go install",
			buildSource:    "goInstall",
			version:        "0.66.2",
			executablePath: "/home/user/go/bin/lazygit",
			expected:       InstallKindGoInstall,
		},
		{
			name:           "Scoop",
			buildSource:    "binaryRelease",
			version:        "0.66.2",
			executablePath: `C:\Users\user\scoop\apps\lazygit\0.66.2\lazygit.exe`,
			expected:       InstallKindScoop,
		},
		{
			name:           "winget",
			buildSource:    "binaryRelease",
			version:        "0.66.2",
			executablePath: `C:\Users\user\AppData\Local\Microsoft\WinGet\Packages\JesseDuffield.lazygit_Microsoft.Winget.Source_8wekyb3d8bbwe\lazygit.exe`,
			expected:       InstallKindWinget,
		},
		{
			name:           "Chocolatey",
			buildSource:    "binaryRelease",
			version:        "0.66.2",
			executablePath: `C:\ProgramData\chocolatey\lib\lazygit\tools\lazygit.exe`,
			expected:       InstallKindChocolatey,
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
