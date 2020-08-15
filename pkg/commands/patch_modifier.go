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

	oldStart := mustConvertToInt(match[1])
	newStart := mustConvertToInt(match[2])
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
	hunks := []*PatchHunk{}
	firstLineIdx := -1
	var hunkLines []string
	pastDiffHeader := false

	for lineIdx, line := range strings.SplitAfter(diff, "\n") {
		isHunkHeader := strings.HasPrefix(line, "@@ -")

		if isHunkHeader {
			if pastDiffHeader { // we need to persist the current hunk
				hunks = append(hunks, newHunk(hunkLines, firstLineIdx))
			}
			pastDiffHeader = true
			firstLineIdx = lineIdx
			hunkLines = []string{line}
			continue
		}

		if !pastDiffHeader { // skip through the stuff that precedes the first hunk
			continue
		}

		hunkLines = append(hunkLines, line)
	}

	if pastDiffHeader {
		hunks = append(hunks, newHunk(hunkLines, firstLineIdx))
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
			if lineIdx >= hunk.FirstLineIdx && lineIdx <= hunk.LastLineIdx() {
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

	return d.hunks[len(d.hunks)-1].LastLineIdx()
}

func ModifiedPatchForRange(log *logrus.Entry, filename string, diffText string, firstLineIdx int, lastLineIdx int, reverse bool, keepOriginalHeader bool) string {
	p := NewPatchModifier(log, filename, diffText)
	return p.ModifiedPatchForRange(firstLineIdx, lastLineIdx, reverse, keepOriginalHeader)
}

func ModifiedPatchForLines(log *logrus.Entry, filename string, diffText string, includedLineIndices []int, reverse bool, keepOriginalHeader bool) string {
	p := NewPatchModifier(log, filename, diffText)
	return p.ModifiedPatchForLines(includedLineIndices, reverse, keepOriginalHeader)
}

// I want to know, given a hunk, what line a given index is on
func (hunk *PatchHunk) LineNumberOfLine(idx int) int {
	lines := hunk.bodyLines[0 : idx-hunk.FirstLineIdx-1]

	offset := nLinesWithPrefix(lines, []string{"+", " "})

	return hunk.newStart + offset
}

func nLinesWithPrefix(lines []string, chars []string) int {
	result := 0
	for _, line := range lines {
		for _, char := range chars {
			if line[:1] == char {
				result++
			}
		}
	}
	return result
}
