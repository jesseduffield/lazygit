package patch

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

type PatchHunk struct {
	FirstLineIdx int
	oldStart     int
	newStart     int
	heading      string
	bodyLines    []string
}

func (hunk *PatchHunk) LastLineIdx() int {
	return hunk.FirstLineIdx + len(hunk.bodyLines)
}

func newHunk(lines []string, firstLineIdx int) *PatchHunk {
	header := lines[0]
	bodyLines := lines[1:]

	oldStart, newStart, heading := headerInfo(header)

	return &PatchHunk{
		oldStart:     oldStart,
		newStart:     newStart,
		heading:      heading,
		FirstLineIdx: firstLineIdx,
		bodyLines:    bodyLines,
	}
}

func headerInfo(header string) (int, int, string) {
	match := hunkHeaderRegexp.FindStringSubmatch(header)

	oldStart := utils.MustConvertToInt(match[1])
	newStart := utils.MustConvertToInt(match[2])
	heading := match[3]

	return oldStart, newStart, heading
}

func (hunk *PatchHunk) updatedLines(lineIndices []int, reverse bool) []string {
	skippedNewlineMessageIndex := -1
	newLines := []string{}

	lineIdx := hunk.FirstLineIdx
	for _, line := range hunk.bodyLines {
		lineIdx++ // incrementing at the start to skip the header line
		if line == "" {
			break
		}
		isLineSelected := lo.Contains(lineIndices, lineIdx)

		firstChar, content := line[:1], line[1:]
		transformedFirstChar := transformedFirstChar(firstChar, reverse, isLineSelected)

		if isLineSelected || (transformedFirstChar == "\\" && skippedNewlineMessageIndex != lineIdx) || transformedFirstChar == " " {
			newLines = append(newLines, transformedFirstChar+content)
			continue
		}

		if transformedFirstChar == "+" {
			// we don't want to include the 'newline at end of file' line if it involves an addition we're not including
			skippedNewlineMessageIndex = lineIdx + 1
		}
	}

	return newLines
}

func transformedFirstChar(firstChar string, reverse bool, isLineSelected bool) string {
	if reverse {
		if !isLineSelected && firstChar == "+" {
			return " "
		} else if firstChar == "-" {
			return "+"
		} else if firstChar == "+" {
			return "-"
		} else {
			return firstChar
		}
	}

	if !isLineSelected && firstChar == "-" {
		return " "
	}

	return firstChar
}

func (hunk *PatchHunk) formatHeader(oldStart int, oldLength int, newStart int, newLength int, heading string) string {
	return fmt.Sprintf("@@ -%d,%d +%d,%d @@%s\n", oldStart, oldLength, newStart, newLength, heading)
}

func (hunk *PatchHunk) formatWithChanges(lineIndices []int, reverse bool, startOffset int) (int, string) {
	bodyLines := hunk.updatedLines(lineIndices, reverse)
	startOffset, header, ok := hunk.updatedHeader(bodyLines, startOffset, reverse)
	if !ok {
		return startOffset, ""
	}
	return startOffset, header + strings.Join(bodyLines, "")
}

func (hunk *PatchHunk) updatedHeader(newBodyLines []string, startOffset int, reverse bool) (int, string, bool) {
	changeCount := nLinesWithPrefix(newBodyLines, []string{"+", "-"})
	oldLength := nLinesWithPrefix(newBodyLines, []string{" ", "-"})
	newLength := nLinesWithPrefix(newBodyLines, []string{"+", " "})

	if changeCount == 0 {
		// if nothing has changed we just return nothing
		return startOffset, "", false
	}

	var oldStart int
	if reverse {
		oldStart = hunk.newStart
	} else {
		oldStart = hunk.oldStart
	}

	var newStartOffset int
	// if the hunk went from zero to positive length, we need to increment the starting point by one
	// if the hunk went from positive to zero length, we need to decrement the starting point by one
	if oldLength == 0 {
		newStartOffset = 1
	} else if newLength == 0 {
		newStartOffset = -1
	} else {
		newStartOffset = 0
	}

	newStart := oldStart + startOffset + newStartOffset

	newStartOffset = startOffset + newLength - oldLength
	formattedHeader := hunk.formatHeader(oldStart, oldLength, newStart, newLength, hunk.heading)
	return newStartOffset, formattedHeader, true
}
