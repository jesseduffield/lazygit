package main

import (
	"fmt"

	"github.com/fatih/color"
)

func coloredString(str string, colorAttribute color.Attribute) string {
	colour := color.New(colorAttribute)
	return coloredStringDirect(str, colour)
}

// used for aggregating a few color attributes rather than just sending a single one
func coloredStringDirect(str string, colour *color.Color) string {
	return colour.SprintFunc()(fmt.Sprint(str))
}
