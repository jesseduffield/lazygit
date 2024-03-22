package config

import (
	"os"
	"strings"
)

func isWSL() bool {
	data, err := os.ReadFile("/proc/sys/kernel/osrelease")
	return err == nil && strings.Contains(string(data), "microsoft")
}

func isContainer() bool {
	data, err := os.ReadFile("/proc/1/cgroup")
	return err == nil && (strings.Contains(string(data), "docker") ||
		strings.Contains(string(data), "/lxc/") ||
		os.Getenv("CONTAINER") != "")
}

// GetPlatformDefaultConfig gets the defaults for the platform
func GetPlatformDefaultConfig() OSConfig {
	if isWSL() && !isContainer() {
		return OSConfig{
			Open:     `powershell.exe start explorer.exe {{filename}} >/dev/null`,
			OpenLink: `powershell.exe start {{link}} >/dev/null`,
		}
	}

	return OSConfig{
		Open:     `xdg-open {{filename}} >/dev/null`,
		OpenLink: `xdg-open {{link}} >/dev/null`,
	}
}
