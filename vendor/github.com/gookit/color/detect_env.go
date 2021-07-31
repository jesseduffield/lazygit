package color

import (
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"

	"github.com/xo/terminfo"
)

/*************************************************************
 * helper methods for detect color supports
 *************************************************************/

// DetectColorLevel for current env
//
// NOTICE: The method will detect terminal info each times,
// 	if only want get current color level, please direct call SupportColor() or TermColorLevel()
func DetectColorLevel() terminfo.ColorLevel {
	level, _ := detectTermColorLevel()
	return level
}

// detect terminal color support level
//
// refer https://github.com/Delta456/box-cli-maker
func detectTermColorLevel() (level terminfo.ColorLevel, needVTP bool) {
	// on windows WSL:
	// - runtime.GOOS == "Linux"
	// - support true-color
	// env:
	// 	WSL_DISTRO_NAME=Debian
	if val := os.Getenv("WSL_DISTRO_NAME"); val != "" {
		// detect WSL as it has True Color support
		if detectWSL() {
			debugf("True Color support on WSL environment")
			return terminfo.ColorLevelMillions, false
		}
	}

	isWin := runtime.GOOS == "windows"
	termVal := os.Getenv("TERM")

	// on TERM=screen: not support true-color
	if termVal != "screen" {
		// On JetBrains Terminal
		// - support true-color
		// env:
		// 	TERMINAL_EMULATOR=JetBrains-JediTerm
		val := os.Getenv("TERMINAL_EMULATOR")
		if val == "JetBrains-JediTerm" {
			debugf("True Color support on JetBrains-JediTerm, is win: %v", isWin)
			return terminfo.ColorLevelMillions, isWin
		}
	}

	// level, err = terminfo.ColorLevelFromEnv()
	level = detectColorLevelFromEnv(termVal, isWin)
	debugf("color level by detectColorLevelFromEnv: %s", level.String())

	// fallback: simple detect by TERM value string.
	if level == terminfo.ColorLevelNone {
		debugf("level none - fallback check special term color support")
		// on Windows: enable VTP as it has True Color support
		level, needVTP = detectSpecialTermColor(termVal)
	}
	return
}

// detectColorFromEnv returns the color level COLORTERM, FORCE_COLOR,
// TERM_PROGRAM, or determined from the TERM environment variable.
//
// refer the terminfo.ColorLevelFromEnv()
// https://en.wikipedia.org/wiki/Terminfo
func detectColorLevelFromEnv(termVal string, isWin bool) terminfo.ColorLevel {
	// check for overriding environment variables
	colorTerm, termProg, forceColor := os.Getenv("COLORTERM"), os.Getenv("TERM_PROGRAM"), os.Getenv("FORCE_COLOR")
	switch {
	case strings.Contains(colorTerm, "truecolor") || strings.Contains(colorTerm, "24bit"):
		if termVal == "screen" { // on TERM=screen: not support true-color
			return terminfo.ColorLevelHundreds
		}
		return terminfo.ColorLevelMillions
	case colorTerm != "" || forceColor != "":
		return terminfo.ColorLevelBasic
	case termProg == "Apple_Terminal":
		return terminfo.ColorLevelHundreds
	case termProg == "Terminus" || termProg == "Hyper":
		if termVal == "screen" { // on TERM=screen: not support true-color
			return terminfo.ColorLevelHundreds
		}
		return terminfo.ColorLevelMillions
	case termProg == "iTerm.app":
		if termVal == "screen" { // on TERM=screen: not support true-color
			return terminfo.ColorLevelHundreds
		}

		// check iTerm version
		ver := os.Getenv("TERM_PROGRAM_VERSION")
		if ver != "" {
			i, err := strconv.Atoi(strings.Split(ver, ".")[0])
			if err != nil {
				saveInternalError(terminfo.ErrInvalidTermProgramVersion)
				// return terminfo.ColorLevelNone
				return terminfo.ColorLevelHundreds
			}
			if i == 3 {
				return terminfo.ColorLevelMillions
			}
		}
		return terminfo.ColorLevelHundreds
	}

	// otherwise determine from TERM's max_colors capability
	if !isWin && termVal != "" {
		debugf("TERM=%s - check color level by load terminfo file", termVal)
		ti, err := terminfo.Load(termVal)
		if err != nil {
			saveInternalError(err)
			return terminfo.ColorLevelNone
		}

		debugf("the loaded term info file is: %s", ti.File)
		v, ok := ti.Nums[terminfo.MaxColors]
		switch {
		case !ok || v <= 16:
			return terminfo.ColorLevelNone
		case ok && v >= 256:
			return terminfo.ColorLevelHundreds
		}
		return terminfo.ColorLevelBasic
	}

	// no TERM env value. default return none level
	return terminfo.ColorLevelNone
	// return terminfo.ColorLevelBasic
}

var detectedWSL bool
var wslContents string

// https://github.com/Microsoft/WSL/issues/423#issuecomment-221627364
func detectWSL() bool {
	if !detectedWSL {
		b := make([]byte, 1024)
		// `cat /proc/version`
		// on mac:
		// 	!not the file!
		// on linux(debian,ubuntu,alpine):
		//	Linux version 4.19.121-linuxkit (root@18b3f92ade35) (gcc version 9.2.0 (Alpine 9.2.0)) #1 SMP Thu Jan 21 15:36:34 UTC 2021
		// on win git bash, conEmu:
		// 	MINGW64_NT-10.0-19042 version 3.1.7-340.x86_64 (@WIN-N0G619FD3UK) (gcc version 9.3.0 (GCC) ) 2020-10-23 13:08 UTC
		// on WSL:
		//  Linux version 4.4.0-19041-Microsoft (Microsoft@Microsoft.com) (gcc version 5.4.0 (GCC) ) #488-Microsoft Mon Sep 01 13:43:00 PST 2020
		f, err := os.Open("/proc/version")
		if err == nil {
			_, _ = f.Read(b) // ignore error
			if err = f.Close(); err != nil {
				saveInternalError(err)
			}

			wslContents = string(b)
		}
		detectedWSL = true
	}
	return strings.Contains(wslContents, "Microsoft")
}

// refer
//  https://github.com/Delta456/box-cli-maker/blob/7b5a1ad8a016ce181e7d8b05e24b54ff60b4b38a/detect_unix.go#L27-L45
// detect WSL as it has True Color support
func isWSL() bool {
	// on windows WSL:
	// - runtime.GOOS == "Linux"
	// - support true-color
	// 	WSL_DISTRO_NAME=Debian
	if val := os.Getenv("WSL_DISTRO_NAME"); val == "" {
		return false
	}

	// `cat /proc/sys/kernel/osrelease`
	// on mac:
	//	!not the file!
	// on linux:
	// 	4.19.121-linuxkit
	// on WSL Output:
	//  4.4.0-19041-Microsoft
	wsl, err := ioutil.ReadFile("/proc/sys/kernel/osrelease")
	if err != nil {
		saveInternalError(err)
		return false
	}

	// it gives "Microsoft" for WSL and "microsoft" for WSL 2
	// it support True-color
	content := strings.ToLower(string(wsl))
	return strings.Contains(content, "microsoft")
}

/*************************************************************
 * helper methods for check env
 *************************************************************/

// IsWindows OS env
func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// IsConsole Determine whether w is one of stderr, stdout, stdin
func IsConsole(w io.Writer) bool {
	o, ok := w.(*os.File)
	if !ok {
		return false
	}

	fd := o.Fd()

	// fix: cannot use 'o == os.Stdout' to compare
	return fd == uintptr(syscall.Stdout) || fd == uintptr(syscall.Stdin) || fd == uintptr(syscall.Stderr)
}

// IsMSys msys(MINGW64) environment, does not necessarily support color
func IsMSys() bool {
	// like "MSYSTEM=MINGW64"
	if len(os.Getenv("MSYSTEM")) > 0 {
		return true
	}

	return false
}

// IsSupportColor check current console is support color.
//
// NOTICE: The method will detect terminal info each times,
// 	if only want get current color level, please direct call SupportColor() or TermColorLevel()
func IsSupportColor() bool {
	return IsSupport16Color()
}

// IsSupportColor check current console is support color.
//
// NOTICE: The method will detect terminal info each times,
// 	if only want get current color level, please direct call SupportColor() or TermColorLevel()
func IsSupport16Color() bool {
	level, _ := detectTermColorLevel()
	return level > terminfo.ColorLevelNone
}

// IsSupport256Color render check
//
// NOTICE: The method will detect terminal info each times,
// 	if only want get current color level, please direct call SupportColor() or TermColorLevel()
func IsSupport256Color() bool {
	level, _ := detectTermColorLevel()
	return level > terminfo.ColorLevelBasic
}

// IsSupportRGBColor check. alias of the IsSupportTrueColor()
//
// NOTICE: The method will detect terminal info each times,
// 	if only want get current color level, please direct call SupportColor() or TermColorLevel()
func IsSupportRGBColor() bool {
	return IsSupportTrueColor()
}

// IsSupportTrueColor render check.
//
// NOTICE: The method will detect terminal info each times,
// 	if only want get current color level, please direct call SupportColor() or TermColorLevel()
//
// ENV:
// "COLORTERM=truecolor"
// "COLORTERM=24bit"
func IsSupportTrueColor() bool {
	level, _ := detectTermColorLevel()
	return level > terminfo.ColorLevelHundreds
}
