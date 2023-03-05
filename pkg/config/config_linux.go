package config

import (
	"io/ioutil"
	"os"
	"strings"
)

func isWSL() bool {
	data, err := ioutil.ReadFile("/proc/sys/kernel/osrelease")
	return err == nil && strings.Contains(string(data), "microsoft")
}

func isContainer() bool {
	data, err := ioutil.ReadFile("/proc/1/cgroup")

	if
	strings.Contains(string(data), "docker")   ||
	strings.Contains(string(data), "/lxc/")    ||
	[]string{string(data)}[0] != "systemd"     &&
	[]string{string(data)}[0] != "init"        ||
    os.Getenv("container") != "" {
		return err == nil && true
	}

	return err == nil && false
}

// GetPlatformDefaultConfig gets the defaults for the platform
func GetPlatformDefaultConfig() OSConfig {
	if isWSL() && !isContainer() {
		return OSConfig{
			EditCommand:         ``,
			EditCommandTemplate: "",
			OpenCommand:         `powershell.exe start explorer.exe {{filename}} >/dev/null`,
			OpenLinkCommand:     `powershell.exe start {{link}} >/dev/null`,
		}
	}

	return OSConfig{
		EditCommand:         ``,
		EditCommandTemplate: "",
		OpenCommand:         `xdg-open {{filename}} >/dev/null`,
		OpenLinkCommand:     `xdg-open {{link}} >/dev/null`,
	}
}
