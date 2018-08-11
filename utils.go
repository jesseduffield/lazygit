package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
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

func withPadding(str string, padding int) string {
	if padding-len(str) < 0 {
		return str
	}
	return str + strings.Repeat(" ", padding-len(str))
}

func coloredString(str string, colorAttribute color.Attribute) string {
	colour := color.New(colorAttribute)
	return coloredStringDirect(str, colour)
}

// used for aggregating a few color attributes rather than just sending a single one
func coloredStringDirect(str string, colour *color.Color) string {
	return colour.SprintFunc()(fmt.Sprint(str))
}

// used to get the project name
func getCurrentProject() string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err.Error())
	}
	return filepath.Base(pwd)
}
