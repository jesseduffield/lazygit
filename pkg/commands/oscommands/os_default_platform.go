//go:build !windows
// +build !windows

package oscommands

import (
	"runtime"
)

func GetPlatform() *Platform {
	return &Platform{
		OS:              runtime.GOOS,
		Shell:           "bash",
		ShellArg:        "-c",
		OpenCommand:     "open {{filename}}",
		OpenLinkCommand: "open {{link}}",
	}
}
