package utils

import (
	"strings"

	"github.com/jesseduffield/generics/slices"
	"github.com/mattn/go-runewidth"
	"github.com/samber/lo"
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
	displayStringsArr = excludeBlankColumns(displayStringsArr)
	padWidths := getPadWidths(displayStringsArr)
	output := getPaddedDisplayStrings(displayStringsArr, padWidths)

	return output
}

// NOTE: this mutates the input slice for the sake of performance
func excludeBlankColumns(displayStringsArr [][]string) [][]string {
	if len(displayStringsArr) == 0 {
		return displayStringsArr
	}

	// if all rows share a blank column, we want to remove that column
	toRemove := []int{}
outer:
	for i := range displayStringsArr[0] {
		for _, strings := range displayStringsArr {
			if strings[i] != "" {
				continue outer
			}
		}
		toRemove = append(toRemove, i)
	}

	if len(toRemove) == 0 {
		return displayStringsArr
	}

	// remove the columns
	for i, strings := range displayStringsArr {
		for j := len(toRemove) - 1; j >= 0; j-- {
			strings = append(strings[:toRemove[j]], strings[toRemove[j]+1:]...)
		}
		displayStringsArr[i] = strings
	}

	return displayStringsArr
}

func getPaddedDisplayStrings(stringArrays [][]string, padWidths []int) string {
	builder := strings.Builder{}
	for i, stringArray := range stringArrays {
		if len(stringArray) == 0 {
			continue
		}
		for j, padWidth := range padWidths {
			if len(stringArray)-1 < j {
				continue
			}
			builder.WriteString(WithPadding(stringArray[j], padWidth))
			builder.WriteString(" ")
		}
		if len(stringArray)-1 < len(padWidths) {
			continue
		}
		builder.WriteString(stringArray[len(padWidths)])

		if i < len(stringArrays)-1 {
			builder.WriteString("\n")
		}
	}
	return builder.String()
}

func getPadWidths(stringArrays [][]string) []int {
	maxWidth := slices.MaxBy(stringArrays, func(stringArray []string) int {
		return len(stringArray)
	})

	if maxWidth-1 < 0 {
		return []int{}
	}
	return slices.Map(lo.Range(maxWidth-1), func(i int) int {
		return slices.MaxBy(stringArrays, func(stringArray []string) int {
			uncoloredStr := Decolorise(stringArray[i])

			return runewidth.StringWidth(uncoloredStr)
		})
	})
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

func ShortSha(sha string) string {
	if len(sha) < 8 {
		return sha
	}
	return sha[:8]
}
