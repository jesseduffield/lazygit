package icons

import (
	"log"

	"github.com/samber/lo"
)

type iconProperties struct {
    icon string
    color uint8 
}

var isIconEnabled = false

func IsIconEnabled() bool {
	return isIconEnabled
}

func SetNerdFontsVersion(version string) {
	if !lo.Contains([]string{"2", "3"}, version) {
		log.Fatalf("Unsupported nerdFontVersion %s", version)
	}

	if version == "2" {
		patchGitIconsForNerdFontsV2()
		patchFileIconsForNerdFontsV2()
	}

	isIconEnabled = true
}
