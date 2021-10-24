package utils

import (
	"regexp"
)

// Decolorise strips a string of color
func Decolorise(str string) string {
	re := regexp.MustCompile(`\x1B\[([0-9]{1,3}(;[0-9]{1,3})*)?[mGK]`)
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
