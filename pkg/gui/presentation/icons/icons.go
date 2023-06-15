package icons

var isIconEnabled = false
var nerdFontsVersion = 2

func IsIconEnabled() bool {
	return isIconEnabled
}

func SetIconEnabled(showIcons bool) {
	isIconEnabled = showIcons
}

func GetNerdFontsVersion() int {
    return nerdFontsVersion
}

func SetNerdFontsVersion(version int) {
    nerdFontsVersion = version
}
