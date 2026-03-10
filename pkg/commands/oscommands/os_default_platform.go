//go:build !windows

package oscommands

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
)

func GetPlatform() *Platform {
	shell := getUserShell()

	prefixForShellFunctionsFile := ""
	if strings.HasSuffix(shell, "bash") {
		prefixForShellFunctionsFile = "shopt -s expand_aliases\n"
	}

	return &Platform{
		OS:                          runtime.GOOS,
		Shell:                       shell,
		ShellArg:                    "-c",
		PrefixForShellFunctionsFile: prefixForShellFunctionsFile,
		OpenCommand:                 "open {{filename}}",
		OpenLinkCommand:             "open {{link}}",
	}
}

func getUserShell() string {
	if shell := os.Getenv("SHELL"); shell != "" {
		return shell
	}

	return "bash"
}

// SetWindowTitle is a no-op on non-Windows platforms.
// Terminal title is set via ANSI escape sequences through gocui.Screen.SetTitle().
func (c *OSCommand) SetWindowTitle(title string) error {
	return nil
}

func TerminateProcessGracefully(cmd *exec.Cmd) error {
	if cmd.Process == nil {
		return nil
	}

	return cmd.Process.Signal(syscall.SIGTERM)
}
