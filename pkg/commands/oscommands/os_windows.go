package oscommands

import (
	"os/exec"

	"github.com/jesseduffield/kill"
)

func GetPlatform() *Platform {
	return &Platform{
		OS:       "windows",
		Shell:    "cmd",
		ShellArg: "/c",
	}
}

// Kill kills a process.
func Kill(cmd *exec.Cmd) error {
	return kill.Kill(cmd)
}
