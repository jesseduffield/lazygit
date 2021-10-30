package utils

import (
	"strings"

	"github.com/mattn/go-runewidth"
)

// WithPadding pads a string as much as you want
func WithPadding(str string, padding int) string {
	uncoloredStr := Decolorise(str)
	width := runewidth.StringWidth(uncoloredStr)
	if padding < width {
		return str
	}
	return str + strings.Repeat(" ", padding-width)
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
			uncoloredStr := Decolorise(strings[i])

			width := runewidth.StringWidth(uncoloredStr)
			if width > padWidths[i] {
				padWidths[i] = width
			}
		}
	}
	return padWidths
}

// TruncateWithEllipsis returns a string, truncated to a certain length, with an ellipsis
func TruncateWithEllipsis(str string, limit int) string {
	if runewidth.StringWidth(str) > limit && limit <= 3 {
		return strings.Repeat(".", limit)
	}
	return runewidth.Truncate(str, limit, "...")
}

func SafeTruncate(str string, limit int) string {
	if len(str) > limit {
		return str[0:limit]
	} else {
		return str
	}
}
