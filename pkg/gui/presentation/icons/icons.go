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
	isIconEnabled     = false
	customNameIconMap = map[string]IconProperties{}
	customExtIconMap  = map[string]IconProperties{}
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

func SetCustomIcons(customIcons config.CustomIconsConfig) {
	for name, icon := range customIcons.Filenames {
		customNameIconMap[name] = IconProperties{
			Icon:  icon.Icon,
			Color: icon.Color,
		}
	}
	for ext, icon := range customIcons.Extensions {
		customExtIconMap[ext] = IconProperties{
			Icon:  icon.Icon,
			Color: icon.Color,
		}
	}
}
