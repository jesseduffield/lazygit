package icons

var isIconEnabled = false

func IsIconEnabled() bool {
	return isIconEnabled
}

func SetIconEnabled(showIcons bool) {
	isIconEnabled = showIcons
}
