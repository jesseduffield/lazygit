package utils

import (
	"regexp"
	"sync"
)

var decoloriseCache = make(map[string]string)
var decoloriseMutex sync.Mutex

// Decolorise strips a string of color
func Decolorise(str string) string {
	decoloriseMutex.Lock()
	defer decoloriseMutex.Unlock()

	if decoloriseCache[str] != "" {
		return decoloriseCache[str]
	}

	re := regexp.MustCompile(`\x1B\[([0-9]{1,3}(;[0-9]{1,3})*)?[mGK]`)
	ret := re.ReplaceAllString(str, "")

	decoloriseCache[str] = ret

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
