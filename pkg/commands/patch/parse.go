package patch

import (
	"regexp"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/utils"
)

// Captures, in order: old start, old length (omitted when 1), new start, new
// length (omitted when 1), and the trailing context.
var hunkHeaderRegexp = regexp.MustCompile(`(?m)^@@ -(\d+)(?:,(\d+))? \+(\d+)(?:,(\d+))? @@(.*)$`)

func Parse(patchStr string) *Patch {
	// ignore trailing newline.
	lines := strings.Split(strings.TrimSuffix(patchStr, "\n"), "\n")

	hunks := []*Hunk{}
	patchHeader := []string{}

	var currentHunk *Hunk
	for _, line := range lines {
		if strings.HasPrefix(line, "@@") {
			oldStart, oldLength, newStart, newLength, headerContext := headerInfo(line)

			currentHunk = &Hunk{
				oldStart:          oldStart,
				newStart:          newStart,
				declaredOldLength: oldLength,
				declaredNewLength: newLength,
				headerContext:     headerContext,
				bodyLines:         []*PatchLine{},
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

func headerInfo(header string) (oldStart int, oldLength int, newStart int, newLength int, headerContext string) {
	match := hunkHeaderRegexp.FindStringSubmatch(header)

	oldStart = utils.MustConvertToInt(match[1])
	oldLength = declaredLength(match[2])
	newStart = utils.MustConvertToInt(match[3])
	newLength = declaredLength(match[4])
	headerContext = match[5]

	return oldStart, oldLength, newStart, newLength, headerContext
}

// declaredLength parses a hunk header length capture, which git omits when it
// is 1 (e.g. "@@ -0,0 +1 @@").
func declaredLength(match string) int {
	if match == "" {
		return 1
	}
	return utils.MustConvertToInt(match)
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
