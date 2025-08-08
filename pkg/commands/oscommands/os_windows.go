package oscommands

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetPlatform() *Platform {
	return &Platform{
		OS:       "windows",
		Shell:    "cmd",
		ShellArg: "/c",
	}
}

func (c *OSCommand) UpdateWindowTitle() error {
	path, getWdErr := os.Getwd()
	if getWdErr != nil {
		return getWdErr
	}
	argString := fmt.Sprint("title ", filepath.Base(path), " - Lazygit")
	return c.Cmd.NewShell(argString, c.UserConfig().OS.ShellFunctionsFile).Run()
}
