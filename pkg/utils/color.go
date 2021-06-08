package utils

import (
	"fmt"
	"regexp"

	"github.com/fatih/color"
)

// ColoredString takes a string and a colour attribute and returns a colored
// string with that attribute
func ColoredString(str string, colorAttributes ...color.Attribute) string {
	colour := color.New(colorAttributes...)
	return ColoredStringDirect(str, colour)
}

// ColoredStringDirect used for aggregating a few color attributes rather than
// just sending a single one
func ColoredStringDirect(str string, colour *color.Color) string {
	return colour.SprintFunc()(fmt.Sprint(str))
}

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
