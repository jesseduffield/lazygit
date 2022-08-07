//go:build !windows
// +build !windows

package oscommands

import (
	"os"
	"runtime"
)

func getShell() string {
	defaultShell := "bash"
	if shell := os.Getenv("SHELL"); shell != defaultShell {
		return shell
	}
	return defaultShell
}

func GetPlatform() *Platform {
	return &Platform{
		OS:              runtime.GOOS,
		Shell:           getShell(),
		ShellArg:        "-ic",
		OpenCommand:     "open {{filename}}",
		OpenLinkCommand: "open {{link}}",
	}
}
