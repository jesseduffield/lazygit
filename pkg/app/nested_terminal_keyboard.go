package app

import (
	"os"
	"runtime"
)

// shouldForceWin32KeyboardProtocol returns true when lazygit is running inside a
// nested editor terminal on Windows. tcell's multi-protocol keyboard negotiation
// can break panel navigation in Vim/Neovim :term buffers; win32-input-mode alone
// is the reliable path there.
func shouldForceWin32KeyboardProtocol(goos string, getenv func(string) string) bool {
	if goos != "windows" {
		return false
	}
	if getenv("TCELL_KEYBOARD_PROTOCOL") != "" {
		return false
	}
	return getenv("VIM_TERMINAL") != "" || getenv("NVIM") != ""
}

func applyNestedTerminalKeyboardWorkaround() {
	if shouldForceWin32KeyboardProtocol(runtime.GOOS, os.Getenv) {
		_ = os.Setenv("TCELL_KEYBOARD_PROTOCOL", "win32")
	}
}