package main

import (
	"strings"

	"github.com/fatih/color"
)

// Branch : A git branch
type Branch struct {
	Name    string
	Recency string
}

func (b *Branch) getDisplayString() string {
	return withPadding(b.Recency, 4) + coloredString(b.Name, b.getColor())
}

func (b *Branch) getColor() color.Attribute {
	switch b.getType() {
	case "feature":
		return color.FgGreen
	case "bugfix":
		return color.FgYellow
	case "hotfix":
		return color.FgRed
	default:
		return color.FgWhite
	}
}

// expected to return feature/bugfix/hotfix or blank string
func (b *Branch) getType() string {
	return strings.Split(b.Name, "/")[0]
}

func withPadding(str string, padding int) string {
	if padding-len(str) < 0 {
		return str
	}
	return str + strings.Repeat(" ", padding-len(str))
}
