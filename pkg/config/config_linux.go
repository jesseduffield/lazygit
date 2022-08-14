package config

import (
	"io/ioutil"
	"strings"
)

func isWSL() bool {
	data, err := ioutil.ReadFile("/proc/sys/kernel/osrelease")
	return err == nil && strings.Contains(string(data), "microsoft")
}

// GetPlatformDefaultConfig gets the defaults for the platform
func GetPlatformDefaultConfig() OSConfig {
	if isWSL() {
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
