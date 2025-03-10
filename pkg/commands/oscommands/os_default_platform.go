//go:build !windows
// +build !windows

package oscommands

import (
	"os"
	"runtime"
	"strings"
)

func GetPlatform() *Platform {
	shell := getUserShell()

	prefixForShellAliasesFile := ""
	if strings.HasSuffix(shell, "bash") {
		prefixForShellAliasesFile = "shopt -s expand_aliases\n"
	}

	return &Platform{
		OS:                        runtime.GOOS,
		Shell:                     shell,
		ShellArg:                  "-c",
		PrefixForShellAliasesFile: prefixForShellAliasesFile,
		OpenCommand:               "open {{filename}}",
		OpenLinkCommand:           "open {{link}}",
	}
}

func getUserShell() string {
	if shell := os.Getenv("SHELL"); shell != "" {
		return shell
	}

	return "bash"
}
