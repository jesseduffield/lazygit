package configdir

import "os"

var hasVendorName = true
var systemSettingFolders = []string{"/Library/Application Support"}
var globalSettingFolder = os.Getenv("HOME") + "/Library/Application Support"
var cacheFolder = os.Getenv("HOME") + "/Library/Caches"
