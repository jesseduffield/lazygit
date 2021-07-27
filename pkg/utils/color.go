package utils

import (
	"encoding/hex"
	"regexp"
)

// Decolorise strips a string of color
func Decolorise(str string) string {
	re := regexp.MustCompile(`\x1B\[([0-9]{1,2}(;[0-9]{1,2})?)?[m|K]`)
	return re.ReplaceAllString(str, "")
}

func getPadWidths(stringArrays [][]string) []int {
	maxWidth := 0
	for _, stringArray := range stringArrays {
		if len(stringArray) > maxWidth {
			maxWidth = len(stringArray)
		}
	}
	if maxWidth-1 < 0 {
		return []int{}
	}
	padWidths := make([]int, maxWidth-1)
	for i := range padWidths {
		for _, strings := range stringArrays {
			uncoloredString := Decolorise(strings[i])
			if len(uncoloredString) > padWidths[i] {
				padWidths[i] = len(uncoloredString)
			}
		}
	}
	return padWidths
}

// GetHexColorValues returns the rgb values of a hex color
func GetHexColorValues(v string) (r uint8, g uint8, b uint8, valid bool) {
	if len(v) == 4 {
		v = string([]byte{v[0], v[1], v[1], v[2], v[2], v[3], v[3]})
	} else if len(v) != 7 {
		return
	}

	if v[0] != '#' {
		return
	}

	rgb, err := hex.DecodeString(v[1:])
	if err != nil {
		return
	}

	return rgb[0], rgb[1], rgb[2], true
}
