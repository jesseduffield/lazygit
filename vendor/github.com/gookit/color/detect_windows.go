// +build windows

// Display color on windows
// refer:
//  golang.org/x/sys/windows
// 	golang.org/x/crypto/ssh/terminal
// 	https://docs.microsoft.com/en-us/windows/console
package color

import (
	"os"
	"syscall"
	"unsafe"

	"github.com/xo/terminfo"
	"golang.org/x/sys/windows"
)

// related docs
// https://docs.microsoft.com/zh-cn/windows/console/console-virtual-terminal-sequences
// https://docs.microsoft.com/zh-cn/windows/console/console-virtual-terminal-sequences#samples
var (
	// isMSys bool
	kernel32 *syscall.LazyDLL

	procGetConsoleMode *syscall.LazyProc
	procSetConsoleMode *syscall.LazyProc
)

func init() {
	if !SupportColor() {
		isLikeInCmd = true
		return
	}

	// if disabled.
	if !Enable {
		return
	}

	// if at windows's ConEmu, Cmder, putty ... terminals not need VTP

	// -------- try force enable colors on windows terminal -------
	tryEnableVTP(needVTP)

	// fetch console screen buffer info
	// err := getConsoleScreenBufferInfo(uintptr(syscall.Stdout), &defScreenInfo)
}

// try force enable colors on windows terminal
func tryEnableVTP(enable bool) bool {
	if !enable {
		return false
	}

	debugf("True-Color by enable VirtualTerminalProcessing on windows")

	initKernel32Proc()

	// enable colors on windows terminal
	if tryEnableOnCONOUT() {
		return true
	}

	return tryEnableOnStdout()
}

func initKernel32Proc() {
	if kernel32 != nil {
		return
	}

	// load related windows dll
	// https://docs.microsoft.com/en-us/windows/console/setconsolemode
	kernel32 = syscall.NewLazyDLL("kernel32.dll")

	procGetConsoleMode = kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode = kernel32.NewProc("SetConsoleMode")
}

func tryEnableOnCONOUT() bool {
	outHandle, err := syscall.Open("CONOUT$", syscall.O_RDWR, 0)
	if err != nil {
		saveInternalError(err)
		return false
	}

	err = EnableVirtualTerminalProcessing(outHandle, true)
	if err != nil {
		saveInternalError(err)
		return false
	}

	return true
}

func tryEnableOnStdout() bool {
	// try direct open syscall.Stdout
	err := EnableVirtualTerminalProcessing(syscall.Stdout, true)
	if err != nil {
		saveInternalError(err)
		return false
	}

	return true
}

// Get the Windows Version and Build Number
var (
	winVersion, _, buildNumber = windows.RtlGetNtVersionNumbers()
)

// refer
//  https://github.com/Delta456/box-cli-maker/blob/7b5a1ad8a016ce181e7d8b05e24b54ff60b4b38a/detect_windows.go#L30-L57
//  https://github.com/gookit/color/issues/25#issuecomment-738727917
// detects the Color Level Supported on windows: cmd, powerShell
func detectSpecialTermColor(termVal string) (tl terminfo.ColorLevel, needVTP bool) {
	if os.Getenv("ConEmuANSI") == "ON" {
		debugf("support True Color by ConEmuANSI=ON")
		// ConEmuANSI is "ON" for generic ANSI support
		// but True Color option is enabled by default
		// I am just assuming that people wouldn't have disabled it
		// Even if it is not enabled then ConEmu will auto round off
		// accordingly
		return terminfo.ColorLevelMillions, false
	}

	// Before Windows 10 Build Number 10586, console never supported ANSI Colors
	if buildNumber < 10586 || winVersion < 10 {
		// Detect if using ANSICON on older systems
		if os.Getenv("ANSICON") != "" {
			conVersion := os.Getenv("ANSICON_VER")
			// 8 bit Colors were only supported after v1.81 release
			if conVersion >= "181" {
				return terminfo.ColorLevelHundreds, false
			}
			return terminfo.ColorLevelBasic, false
		}

		return terminfo.ColorLevelNone, false
	}

	// True Color is not available before build 14931 so fallback to 8 bit color.
	if buildNumber < 14931 {
		return terminfo.ColorLevelHundreds, true
	}

	// Windows 10 build 14931 is the first release that supports 16m/TrueColor
	debugf("support True Color on windows version is >= build 14931")
	return terminfo.ColorLevelMillions, true
}

/*************************************************************
 * render full color code on windows(8,16,24bit color)
 *************************************************************/

// docs https://docs.microsoft.com/zh-cn/windows/console/getconsolemode#parameters
const (
	// equals to docs page's ENABLE_VIRTUAL_TERMINAL_PROCESSING 0x0004
	EnableVirtualTerminalProcessingMode uint32 = 0x4
)

// EnableVirtualTerminalProcessing Enable virtual terminal processing
//
// ref from github.com/konsorten/go-windows-terminal-sequences
// doc https://docs.microsoft.com/zh-cn/windows/console/console-virtual-terminal-sequences#samples
//
// Usage:
// 	err := EnableVirtualTerminalProcessing(syscall.Stdout, true)
// 	// support print color text
// 	err = EnableVirtualTerminalProcessing(syscall.Stdout, false)
func EnableVirtualTerminalProcessing(stream syscall.Handle, enable bool) error {
	var mode uint32
	// Check if it is currently in the terminal
	// err := syscall.GetConsoleMode(syscall.Stdout, &mode)
	err := syscall.GetConsoleMode(stream, &mode)
	if err != nil {
		// fmt.Println("EnableVirtualTerminalProcessing", err)
		return err
	}

	if enable {
		mode |= EnableVirtualTerminalProcessingMode
	} else {
		mode &^= EnableVirtualTerminalProcessingMode
	}

	ret, _, err := procSetConsoleMode.Call(uintptr(stream), uintptr(mode))
	if ret == 0 {
		return err
	}

	return nil
}

// renderColorCodeOnCmd enable cmd color render.
// func renderColorCodeOnCmd(fn func()) {
// 	err := EnableVirtualTerminalProcessing(syscall.Stdout, true)
// 	// if is not in terminal, will clear color tag.
// 	if err != nil {
// 		// panic(err)
// 		fn()
// 		return
// 	}
//
// 	// force open color render
// 	old := ForceOpenColor()
// 	fn()
// 	// revert color setting
// 	supportColor = old
//
// 	err = EnableVirtualTerminalProcessing(syscall.Stdout, false)
// 	if err != nil {
// 		panic(err)
// 	}
// }

/*************************************************************
 * render simple color code on windows
 *************************************************************/

// IsTty returns true if the given file descriptor is a terminal.
func IsTty(fd uintptr) bool {
	initKernel32Proc()

	var st uint32
	r, _, e := syscall.Syscall(procGetConsoleMode.Addr(), 2, fd, uintptr(unsafe.Pointer(&st)), 0)
	return r != 0 && e == 0
}

// IsTerminal returns true if the given file descriptor is a terminal.
//
// Usage:
// 	fd := os.Stdout.Fd()
// 	fd := uintptr(syscall.Stdout) // for windows
// 	IsTerminal(fd)
func IsTerminal(fd uintptr) bool {
	initKernel32Proc()

	var st uint32
	r, _, e := syscall.Syscall(procGetConsoleMode.Addr(), 2, fd, uintptr(unsafe.Pointer(&st)), 0)
	return r != 0 && e == 0
}
