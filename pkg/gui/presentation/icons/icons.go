package icons

import (
	"log"

	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/samber/lo"
)

type IconProperties struct {
	Icon  string
	Color string
}

var (
	isIconEnabled = false
	customIcons   = config.CustomIconsConfig{}
)

func IsIconEnabled() bool {
	return isIconEnabled
}

func SetNerdFontsVersion(version string) {
	if version == "" {
		isIconEnabled = false
	} else {
		if !lo.Contains([]string{"2", "3"}, version) {
			log.Fatalf("Unsupported nerdFontVersion %s", version)
		}

		if version == "2" {
			patchGitIconsForNerdFontsV2()
			patchFileIconsForNerdFontsV2()
		}

		isIconEnabled = true
	}
}

func SetCustomIcons(icons config.CustomIconsConfig) {
	customIcons = icons
}
