//go:build !windows
// +build !windows

package oscommands

import (
	"runtime"

	"github.com/jesseduffield/lazygit/pkg/config"
)

func GetPlatform(osConfig config.OSConfig) *Platform {
	platform := Platform{
		OS:              runtime.GOOS,
		Shell:           "bash",
		ShellArg:        "-c",
		OpenCommand:     "open {{filename}}",
		OpenLinkCommand: "open {{link}}",
	}
	if osConfig.Shell != "" {
		platform.Shell = osConfig.Shell
		platform.ShellArg = osConfig.ShellArg
	}
	return &platform
}
