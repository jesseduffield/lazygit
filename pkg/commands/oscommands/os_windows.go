package oscommands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

// setRawCmdLine hands cmd.exe the exact command line we built, bypassing
// os/exec's default composition (which quotes args with the
// CommandLineToArgvW `\"` convention that cmd.exe doesn't understand).
//
// The shell-building logic in NewShell is portable and dispatches on
// platform.OS, which keeps it (and its quoting) unit-testable on any host.
// Assigning SysProcAttr.CmdLine is the only step that needs a Windows-only
// field, so it's the single piece split out behind a build tag; every other
// platform gets the no-op in os_default_platform.go.
func setRawCmdLine(cmd *exec.Cmd, cmdLine string) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.CmdLine = cmdLine
}

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

func TerminateProcessGracefully(cmd *exec.Cmd) error {
	// Signals other than SIGKILL are not supported on Windows
	return nil
}
