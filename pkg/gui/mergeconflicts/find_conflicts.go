package mergeconflicts

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/utils"
)

// LineType tells us whether a given line is a start/middle/end marker of a conflict,
// or if it's not a marker at all
type LineType int

const (
	START LineType = iota
	ANCESTOR
	TARGET
	END
	NOT_A_MARKER
)

// The number of characters a conflict marker consists of, unless the file's
// conflict-marker-size gitattribute says otherwise.
const defaultConflictMarkerSize = 7

// The marker size that everything in here takes is the conflict-marker-size
// gitattribute of the file being examined, which is 0 for a file that doesn't
// have that attribute. Git falls back to its default size in that case, so we
// do the same.
func effectiveMarkerSize(markerSize int) int {
	if markerSize < 1 {
		return defaultConflictMarkerSize
	}

	return markerSize
}

func findConflicts(content string, markerSize int) []*mergeConflict {
	conflicts := make([]*mergeConflict, 0)

	if content == "" {
		return conflicts
	}

	var newConflict *mergeConflict
	for i, line := range utils.SplitLines(content) {
		switch determineLineType(line, markerSize) {
		case START:
			newConflict = &mergeConflict{start: i, ancestor: -1}
		case ANCESTOR:
			if newConflict != nil {
				newConflict.ancestor = i
			}
		case TARGET:
			if newConflict != nil {
				newConflict.target = i
			}
		case END:
			if newConflict != nil {
				newConflict.end = i
				conflicts = append(conflicts, newConflict)
			}
			// reset value to avoid any possible silent mutations in further iterations
			newConflict = nil
		default:
			// line isn't a merge conflict marker so we just continue
		}
	}

	return conflicts
}

func determineLineType(line string, markerSize int) LineType {
	markerSize = effectiveMarkerSize(markerSize)

	// TODO: find out whether we ever actually get this prefix
	trimmedLine := strings.TrimPrefix(line, "++")

	switch {
	case isConflictMarker(trimmedLine, '<', markerSize):
		return START
	case isConflictMarker(trimmedLine, '|', markerSize):
		return ANCESTOR
	case isTargetMarker(trimmedLine, markerSize):
		return TARGET
	case isConflictMarker(trimmedLine, '>', markerSize):
		return END
	default:
		return NOT_A_MARKER
	}
}

// Tells us whether the line begins with markerSize repetitions of markerChar.
func hasMarkerPrefix[T string | []byte](line T, markerChar byte, markerSize int) bool {
	if len(line) < markerSize {
		return false
	}

	for i := range markerSize {
		if line[i] != markerChar {
			return false
		}
	}

	return true
}

// A start, ancestor or end marker is followed by a space and a label, e.g.
// "<<<<<<< HEAD". The label can be missing though, in which case git doesn't
// write the space either; `git checkout -m` with the diff3 conflict style does
// that for the ancestor marker, for example.
func isConflictMarker[T string | []byte](line T, markerChar byte, markerSize int) bool {
	return hasMarkerPrefix(line, markerChar, markerSize) &&
		(len(line) == markerSize || line[markerSize] == ' ')
}

// The marker separating the two sides of a conflict never has a label after it.
func isTargetMarker(line string, markerSize int) bool {
	return hasMarkerPrefix(line, '=', markerSize) && len(line) == markerSize
}

// tells us whether a file actually has inline merge conflicts. We need to run this
// because git will continue showing a status of 'UU' even after the conflicts have
// been resolved in the user's editor
func FileHasConflictMarkers(path string, markerSize int) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}

	defer file.Close()

	return fileHasConflictMarkersAux(file, markerSize)
}

// Efficiently scans through a file looking for merge conflict markers. Returns true if it does
func fileHasConflictMarkersAux(file io.Reader, markerSize int) (bool, error) {
	markerSize = effectiveMarkerSize(markerSize)

	scanner := bufio.NewScanner(file)
	scanner.Split(utils.ScanLinesAndTruncateWhenLongerThanBuffer(bufio.MaxScanTokenSize))
	for scanner.Scan() {
		line := scanner.Bytes()

		// only searching for start/end markers because the others are more ambiguous
		if isConflictMarker(line, '<', markerSize) || isConflictMarker(line, '>', markerSize) {
			return true, nil
		}
	}

	return false, scanner.Err()
}
