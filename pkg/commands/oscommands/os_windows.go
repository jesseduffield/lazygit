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

// SetWindowTitle sets the terminal window title using Windows cmd.exe's title command.
// This is a fallback for the legacy Windows Console which doesn't support ANSI escape sequences.
func (c *OSCommand) SetWindowTitle(title string) error {
	return c.Cmd.NewShell("title "+title, c.UserConfig().OS.ShellFunctionsFile).Run()
}

func TerminateProcessGracefully(cmd *exec.Cmd) error {
	// Signals other than SIGKILL are not supported on Windows
	return nil
}
