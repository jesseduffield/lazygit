package oscommands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"unsafe"
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

var (
	kernel32            = syscall.NewLazyDLL("kernel32.dll")
	procSetConsoleTitle = kernel32.NewProc("SetConsoleTitleW")
)

// UpdateWindowTitle sets the console window title directly via the
// SetConsoleTitleW Win32 API instead of shelling out to `cmd /c title ...`.
//
// The repo name is attacker/user-controlled input (it's just the current
// directory's basename), and cmd.exe treats characters such as & as command
// separators. A directory named e.g. "test&aaa" would previously be split
// into two commands ("title test" and "aaa"), and cmd.exe would then try to
// run "aaa" as a program, printing an error to stderr that crashed the run
// (see #5766). Calling the Win32 API directly sidesteps cmd.exe's argument
// parsing entirely: the title string is passed as-is, verbatim, regardless
// of what characters it contains.
func (c *OSCommand) UpdateWindowTitle() error {
	path, getWdErr := os.Getwd()
	if getWdErr != nil {
		return getWdErr
	}
	title := fmt.Sprint(filepath.Base(path), " - Lazygit")

	titlePtr, err := syscall.UTF16PtrFromString(title)
	if err != nil {
		return err
	}

	r1, _, callErr := procSetConsoleTitle.Call(uintptr(unsafe.Pointer(titlePtr)))
	if r1 == 0 {
		return callErr
	}
	return nil
}

func TerminateProcessGracefully(proc *os.Process) error {
	// Signals other than SIGKILL are not supported on Windows
	return nil
}
