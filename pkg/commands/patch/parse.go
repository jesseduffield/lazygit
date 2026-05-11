package patch

import (
	"regexp"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/utils"
)

var hunkHeaderRegexp = regexp.MustCompile(`(?m)^@@ -(\d+)[^\+]+\+(\d+)[^@]+@@(.*)$`)

func Parse(patchStr string) *Patch {
	// ignore trailing newline.
	lines := strings.Split(strings.TrimSuffix(patchStr, "\n"), "\n")

	hunks := []*Hunk{}
	patchHeader := []string{}

	var currentHunk *Hunk
	for _, line := range lines {
		if strings.HasPrefix(line, "@@") {
			oldStart, newStart, headerContext := headerInfo(line)

			currentHunk = &Hunk{
				oldStart:      oldStart,
				newStart:      newStart,
				headerContext: headerContext,
				bodyLines:     []*PatchLine{},
			}
			hunks = append(hunks, currentHunk)
		} else if currentHunk != nil {
			currentHunk.bodyLines = append(currentHunk.bodyLines, newHunkLine(line))
		} else {
			patchHeader = append(patchHeader, line)
		}
	}

	return &Patch{
		hunks:  hunks,
		header: patchHeader,
	}
}

func headerInfo(header string) (int, int, string) {
	match := hunkHeaderRegexp.FindStringSubmatch(header)

	oldStart := utils.MustConvertToInt(match[1])
	newStart := utils.MustConvertToInt(match[2])
	headerContext := match[3]

	return oldStart, newStart, headerContext
}

func newHunkLine(line string) *PatchLine {
	if line == "" {
		return &PatchLine{
			Kind:    CONTEXT,
			Content: "",
		}
	}

	firstChar := line[:1]

	kind := parseFirstChar(firstChar)

	return &PatchLine{
		Kind:    kind,
		Content: line,
	}
}

func parseFirstChar(firstChar string) PatchLineKind {
	switch firstChar {
	case " ":
		return CONTEXT
	case "+":
		return ADDITION
	case "-":
		return DELETION
	case "\\":
		return NEWLINE_MESSAGE
	}

	return CONTEXT
}
