package patch

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	hunkHeaderRegexp  = regexp.MustCompile(`(?m)^@@ -(\d+)[^\+]+\+(\d+)[^@]+@@(.*)$`)
	patchHeaderRegexp = regexp.MustCompile(`(?ms)(^diff.*?)^@@`)
)

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
	var hunkLines []string //nolint:prealloc
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
	n := idx - hunk.FirstLineIdx - 1
	if n < 0 {
		n = 0
	} else if n >= len(hunk.bodyLines) {
		n = len(hunk.bodyLines) - 1
	}

	lines := hunk.bodyLines[0:n]

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
