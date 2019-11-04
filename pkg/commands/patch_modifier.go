package commands

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/sirupsen/logrus"
)

var hunkHeaderRegexp = regexp.MustCompile(`(?m)^@@ -(\d+)[^\+]+\+(\d+)[^@]+@@(.*)$`)
var patchHeaderRegexp = regexp.MustCompile(`(?ms)(^diff.*?)^@@`)

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

func (hunk *PatchHunk) updatedLines(lineIndices []int, reverse bool) []string {
	skippedNewlineMessageIndex := -1
	newLines := []string{}

	lineIdx := hunk.FirstLineIdx
	for _, line := range hunk.bodyLines {
		lineIdx++ // incrementing at the start to skip the header line
		if line == "" {
			break
		}
		isLineSelected := utils.IncludesInt(lineIndices, lineIdx)

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
	match := hunkHeaderRegexp.FindStringSubmatch(hunk.header)

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

func GetHeaderFromDiff(diff string) string {
	match := patchHeaderRegexp.FindStringSubmatch(diff)
	if len(match) <= 1 {
		return ""
	}
	return match[1]
}

func GetHunksFromDiff(diff string) []*PatchHunk {
	headers := hunkHeaderRegexp.FindAllString(diff, -1)
	bodies := hunkHeaderRegexp.Split(diff, -1)[1:] // discarding top bit

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
	header   string
}

func NewPatchModifier(log *logrus.Entry, filename string, diffText string) *PatchModifier {
	return &PatchModifier{
		Log:      log,
		filename: filename,
		hunks:    GetHunksFromDiff(diffText),
		header:   GetHeaderFromDiff(diffText),
	}
}

func (d *PatchModifier) ModifiedPatchForLines(lineIndices []int, reverse bool, keepOriginalHeader bool) string {
	// step one is getting only those hunks which we care about
	hunksInRange := []*PatchHunk{}
outer:
	for _, hunk := range d.hunks {
		// if there is any line in our lineIndices array that the hunk contains, we append it
		for _, lineIdx := range lineIndices {
			if lineIdx >= hunk.FirstLineIdx && lineIdx <= hunk.LastLineIdx {
				hunksInRange = append(hunksInRange, hunk)
				continue outer
			}
		}
	}

	// step 2 is collecting all the hunks with new headers
	startOffset := 0
	formattedHunks := ""
	var formattedHunk string
	for _, hunk := range hunksInRange {
		startOffset, formattedHunk = hunk.formatWithChanges(lineIndices, reverse, startOffset)
		formattedHunks += formattedHunk
	}

	if formattedHunks == "" {
		return ""
	}

	var fileHeader string
	// for staging/unstaging lines we don't want the original header because
	// it makes git confused e.g. when dealing with deleted/added files
	// but with building and applying patches the original header gives git
	// information it needs to cleanly apply patches
	if keepOriginalHeader {
		fileHeader = d.header
	} else {
		fileHeader = fmt.Sprintf("--- a/%s\n+++ b/%s\n", d.filename, d.filename)
	}

	return fileHeader + formattedHunks
}

func (d *PatchModifier) ModifiedPatchForRange(firstLineIdx int, lastLineIdx int, reverse bool, keepOriginalHeader bool) string {
	// generate array of consecutive line indices from our range
	selectedLines := []int{}
	for i := firstLineIdx; i <= lastLineIdx; i++ {
		selectedLines = append(selectedLines, i)
	}
	return d.ModifiedPatchForLines(selectedLines, reverse, keepOriginalHeader)
}

func (d *PatchModifier) OriginalPatchLength() int {
	if len(d.hunks) == 0 {
		return 0
	}

	return d.hunks[len(d.hunks)-1].LastLineIdx
}

func ModifiedPatchForRange(log *logrus.Entry, filename string, diffText string, firstLineIdx int, lastLineIdx int, reverse bool, keepOriginalHeader bool) string {
	p := NewPatchModifier(log, filename, diffText)
	return p.ModifiedPatchForRange(firstLineIdx, lastLineIdx, reverse, keepOriginalHeader)
}
