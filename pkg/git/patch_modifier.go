package git

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

var headerRegexp = regexp.MustCompile(`(?m)^@@ -(\d+)[^\+]+\+(\d+)[^@]+@@(.*)$`)

type PatchHunk struct {
	header       string
	FirstLineIdx int
	LastLineIdx  int
	bodyLines    []string
}

func newHunk(header string, body string, firstLineIdx int) *PatchHunk {
	bodyLines := strings.SplitAfter(header+body, "\n")[1:] // dropping the header line

	return &PatchHunk{
		header:       header,
		FirstLineIdx: firstLineIdx,
		LastLineIdx:  firstLineIdx + len(bodyLines),
		bodyLines:    bodyLines,
	}
}

func (hunk *PatchHunk) updatedLinesForRange(firstLineIdx int, lastLineIdx int, reverse bool) []string {
	skippedNewlineMessageIndex := -1
	newLines := []string{}

	lineIdx := hunk.FirstLineIdx
	for _, line := range hunk.bodyLines {
		lineIdx++ // incrementing at the start to skip the header line
		if line == "" {
			break
		}
		isLineInsideRange := (firstLineIdx <= lineIdx && lineIdx <= lastLineIdx)

		firstChar, content := line[:1], line[1:]
		transformedFirstChar := transformedFirstChar(firstChar, reverse, isLineInsideRange)

		if isLineInsideRange || (transformedFirstChar == "\\" && skippedNewlineMessageIndex != lineIdx) || transformedFirstChar == " " {
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

func transformedFirstChar(firstChar string, reverse bool, isLineInsideRange bool) string {
	if reverse {
		if !isLineInsideRange && firstChar == "+" {
			return " "
		} else if firstChar == "-" {
			return "+"
		} else if firstChar == "+" {
			return "-"
		} else {
			return firstChar
		}
	}

	if !isLineInsideRange && firstChar == "-" {
		return " "
	}

	return firstChar
}

func (hunk *PatchHunk) formatHeader(oldStart int, oldLength int, newStart int, newLength int, heading string) string {
	return fmt.Sprintf("@@ -%d,%d +%d,%d @@%s\n", oldStart, oldLength, newStart, newLength, heading)
}

func (hunk *PatchHunk) formatWithChanges(firstLineIdx int, lastLineIdx int, reverse bool, startOffset int) (int, string) {
	bodyLines := hunk.updatedLinesForRange(firstLineIdx, lastLineIdx, reverse)
	startOffset, header, ok := hunk.updatedHeader(bodyLines, startOffset, reverse)
	if !ok {
		return startOffset, ""
	}
	return startOffset, header + strings.Join(bodyLines, "")
}

func (hunk *PatchHunk) updatedHeader(newBodyLines []string, startOffset int, reverse bool) (int, string, bool) {
	changeCount := 0
	oldLength := 0
	newLength := 0
	for _, line := range newBodyLines {
		switch line[:1] {
		case "+":
			newLength++
			changeCount++
		case "-":
			oldLength++
			changeCount++
		case " ":
			oldLength++
			newLength++
		}
	}

	if changeCount == 0 {
		// if nothing has changed we just return nothing
		return startOffset, "", false
	}

	// get oldstart, newstart, and heading from header
	match := headerRegexp.FindStringSubmatch(hunk.header)

	var oldStart int
	if reverse {
		oldStart = mustConvertToInt(match[2])
	} else {
		oldStart = mustConvertToInt(match[1])
	}
	heading := match[3]

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
	formattedHeader := hunk.formatHeader(oldStart, oldLength, newStart, newLength, heading)
	return newStartOffset, formattedHeader, true
}

func mustConvertToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func GetHunksFromDiff(diff string) []*PatchHunk {
	headers := headerRegexp.FindAllString(diff, -1)
	bodies := headerRegexp.Split(diff, -1)[1:] // discarding top bit

	headerFirstLineIndices := []int{}
	for lineIdx, line := range strings.Split(diff, "\n") {
		if strings.HasPrefix(line, "@@ -") {
			headerFirstLineIndices = append(headerFirstLineIndices, lineIdx)
		}
	}

	hunks := make([]*PatchHunk, len(headers))
	for index, header := range headers {
		hunks[index] = newHunk(header, bodies[index], headerFirstLineIndices[index])
	}

	return hunks
}

type PatchModifier struct {
	Log      *logrus.Entry
	filename string
	hunks    []*PatchHunk
}

func NewPatchModifier(log *logrus.Entry, filename string, diffText string) *PatchModifier {
	return &PatchModifier{
		Log:      log,
		filename: filename,
		hunks:    GetHunksFromDiff(diffText),
	}
}

func (d *PatchModifier) ModifiedPatchForRange(firstLineIdx int, lastLineIdx int, reverse bool) string {
	// step one is getting only those hunks which we care about
	hunksInRange := []*PatchHunk{}
	for _, hunk := range d.hunks {
		if hunk.LastLineIdx >= firstLineIdx && hunk.FirstLineIdx <= lastLineIdx {
			hunksInRange = append(hunksInRange, hunk)
		}
	}

	// step 2 is collecting all the hunks with new headers
	startOffset := 0
	formattedHunks := ""
	var formattedHunk string
	for _, hunk := range hunksInRange {
		startOffset, formattedHunk = hunk.formatWithChanges(firstLineIdx, lastLineIdx, reverse, startOffset)
		formattedHunks += formattedHunk
	}

	if formattedHunks == "" {
		return ""
	}

	fileHeader := fmt.Sprintf("--- a/%s\n+++ b/%s\n", d.filename, d.filename)

	return fileHeader + formattedHunks
}

func ModifiedPatch(log *logrus.Entry, filename string, diffText string, firstLineIdx int, lastLineIdx int, reverse bool) string {
	p := NewPatchModifier(log, filename, diffText)
	return p.ModifiedPatchForRange(firstLineIdx, lastLineIdx, reverse)
}
