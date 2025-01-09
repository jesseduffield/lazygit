package icons

import (
	"testing"
)

func TestFileIcons(t *testing.T) {
	t.Run("TestFileIcons", func(t *testing.T) {
		for name, icon := range nameIconMap {
			if len([]rune(icon.Icon)) != 1 {
				t.Errorf("nameIconMap[\"%s\"] is not a single rune", name)
			}
		}

		for ext, icon := range extIconMap {
			if len([]rune(icon.Icon)) != 1 {
				t.Errorf("extIconMap[\"%s\"] is not a single rune", ext)
			}
		}
	})
}
