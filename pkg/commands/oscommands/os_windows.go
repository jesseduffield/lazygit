package oscommands

import (
	"os/exec"
)

func GetPlatform() *Platform {
	return &Platform{
		OS:       "windows",
		Shell:    "cmd",
		ShellArg: "/c",
	}
}

func TerminateProcessGracefully(cmd *exec.Cmd) error {
	// Signals other than SIGKILL are not supported on Windows
	return nil
}
