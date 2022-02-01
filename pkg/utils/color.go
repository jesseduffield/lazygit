package utils

import (
	"regexp"
	"sync"

	"github.com/gookit/color"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
)

var decoloriseCache = make(map[string]string)
var decoloriseMutex sync.RWMutex

// Decolorise strips a string of color
func Decolorise(str string) string {
	decoloriseMutex.RLock()
	val := decoloriseCache[str]
	decoloriseMutex.RUnlock()

	if val != "" {
		return val
	}

	re := regexp.MustCompile(`\x1B\[([0-9]{1,3}(;[0-9]{1,3})*)?[mGK]`)
	ret := re.ReplaceAllString(str, "")

	decoloriseMutex.Lock()
	decoloriseCache[str] = ret
	decoloriseMutex.Unlock()

	return ret
}

func IsValidHexValue(v string) bool {
	if len(v) != 4 && len(v) != 7 {
		return false
	}

	if v[0] != '#' {
		return false
	}

	for _, char := range v[1:] {
		switch char {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f', 'A', 'B', 'C', 'D', 'E', 'F':
			continue
		default:
			return false
		}
	}

	return true
}

func SetCustomColors(customColors map[string]string) map[string]style.TextStyle {
	colors := make(map[string]style.TextStyle)
	for key, colorSequence := range customColors {
		style := style.New().SetFg(style.NewRGBColor(color.HEX(colorSequence, false)))
		colors[key] = style
	}
	return colors
}
