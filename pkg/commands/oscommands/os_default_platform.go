//go:build !windows
// +build !windows

package oscommands

import (
	"os"
	"runtime"
)

func GetPlatform() *Platform {
	return &Platform{
		OS:                  runtime.GOOS,
		Shell:               "bash",
		InteractiveShell:    getUserShell(),
		ShellArg:            "-c",
		InteractiveShellArg: "-i",
		OpenCommand:         "open {{filename}}",
		OpenLinkCommand:     "open {{link}}",
	}
}

func getUserShell() string {
	if shell := os.Getenv("SHELL"); shell != "" {
		return shell
	}

	return "bash"
}
