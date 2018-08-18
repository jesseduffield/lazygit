package configdir

import "os"

var hasVendorName = true
var systemSettingFolders = []string{os.Getenv("PROGRAMDATA")}
var globalSettingFolder = os.Getenv("APPDATA")
var cacheFolder = os.Getenv("LOCALAPPDATA")
