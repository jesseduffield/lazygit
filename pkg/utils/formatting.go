package utils

import (
	"strings"

	"github.com/mattn/go-runewidth"
)

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

			width := runewidth.StringWidth(uncoloredString)
			if width > padWidths[i] {
				padWidths[i] = width
			}
		}
	}
	return padWidths
}
