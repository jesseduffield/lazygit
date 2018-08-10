package main

import (
	"strings"

	"github.com/jesseduffield/gocui"
)

func splitLines(multilineString string) []string {
	multilineString = strings.Replace(multilineString, "\r", "", -1)
	if multilineString == "" || multilineString == "\n" {
		return make([]string, 0)
	}
	lines := strings.Split(multilineString, "\n")
	if lines[len(lines)-1] == "" {
		return lines[:len(lines)-1]
	}
	return lines
}

func trimmedContent(v *gocui.View) string {
	return strings.TrimSpace(v.Buffer())
}
