// +build !windows,!darwin

package configdir

import (
	"os"
	"path/filepath"
	"strings"
)

// https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html

var hasVendorName = true
var systemSettingFolders []string
var globalSettingFolder string
var cacheFolder string

func init() {
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		globalSettingFolder = os.Getenv("XDG_CONFIG_HOME")
	} else {
		globalSettingFolder = filepath.Join(os.Getenv("HOME"), ".config")
	}
	if os.Getenv("XDG_CONFIG_DIRS") != "" {
		systemSettingFolders = strings.Split(os.Getenv("XDG_CONFIG_DIRS"), ":")
	} else {
		systemSettingFolders = []string{"/etc/xdg"}
	}
	if os.Getenv("XDG_CACHE_HOME") != "" {
		cacheFolder = os.Getenv("XDG_CACHE_HOME")
	} else {
		cacheFolder = filepath.Join(os.Getenv("HOME"), ".cache")
	}
}
