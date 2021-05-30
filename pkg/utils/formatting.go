package utils

import "strings"

// WithPadding pads a string as much as you want
func WithPadding(str string, padding int) string {
	uncoloredStr := Decolorise(str)
	if padding < len(uncoloredStr) {
		return str
	}
	return str + strings.Repeat(" ", padding-len(uncoloredStr))
}

func RenderDisplayStrings(displayStringsArr [][]string) string {
	padWidths := getPadWidths(displayStringsArr)
	paddedDisplayStrings := getPaddedDisplayStrings(displayStringsArr, padWidths)

	return strings.Join(paddedDisplayStrings, "\n")
}

func getPaddedDisplayStrings(stringArrays [][]string, padWidths []int) []string {
	paddedDisplayStrings := make([]string, len(stringArrays))
	for i, stringArray := range stringArrays {
		if len(stringArray) == 0 {
			continue
		}
		for j, padWidth := range padWidths {
			if len(stringArray)-1 < j {
				continue
			}
			paddedDisplayStrings[i] += WithPadding(stringArray[j], padWidth) + " "
		}
		if len(stringArray)-1 < len(padWidths) {
			continue
		}
		paddedDisplayStrings[i] += stringArray[len(padWidths)]
	}
	return paddedDisplayStrings
}
