package config

import (
	"fmt"
	"os"
	"strings"
)

func isWSL() bool {
	data, err := os.ReadFile("/proc/sys/kernel/osrelease")
	return err == nil && strings.Contains(string(data), "microsoft")
}

func isContainer() bool {
	data, err := os.ReadFile("/proc/1/cgroup")

	if strings.Contains(string(data), "docker") ||
		strings.Contains(string(data), "/lxc/") ||
		[]string{string(data)}[0] != "systemd" &&
			[]string{string(data)}[0] != "init" ||
		os.Getenv("container") != "" {
		return err == nil && true
	}

	return err == nil && false
}

// GetPlatformDefaultConfig gets the defaults for the platform
func GetPlatformDefaultConfig() OSConfig {
	osConfig := OSConfig{}

	if isWSL() && !isContainer() {
		osConfig.Open = `powershell.exe start explorer.exe {{filename}} >/dev/null`
		osConfig.OpenLink = `powershell.exe start {{link}} >/dev/null`
	} else {
		osConfig.Open = `xdg-open {{filename}} >/dev/null`
		osConfig.OpenLink = `xdg-open {{link}} >/dev/null`
	}

	if browser, ok := os.LookupEnv("BROWSER"); ok {
		osConfig.OpenLink = fmt.Sprintf(`"%s" {{link}} >/dev/null`, browser)
	}

	return osConfig
}
