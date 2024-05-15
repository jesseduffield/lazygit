package mergeconflicts

import (
	"bufio"
	"bytes"
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

func findConflicts(content string) []*mergeConflict {
	conflicts := make([]*mergeConflict, 0)

	if content == "" {
		return conflicts
	}

	var newConflict *mergeConflict
	for i, line := range utils.SplitLines(content) {
		switch determineLineType(line) {
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

var (
	CONFLICT_START       = "<<<<<<< "
	CONFLICT_END         = ">>>>>>> "
	CONFLICT_START_BYTES = []byte(CONFLICT_START)
	CONFLICT_END_BYTES   = []byte(CONFLICT_END)
)

func determineLineType(line string) LineType {
	// TODO: find out whether we ever actually get this prefix
	trimmedLine := strings.TrimPrefix(line, "++")

	switch {
	case strings.HasPrefix(trimmedLine, CONFLICT_START):
		return START
	case strings.HasPrefix(trimmedLine, "||||||| "):
		return ANCESTOR
	case trimmedLine == "=======":
		return TARGET
	case strings.HasPrefix(trimmedLine, CONFLICT_END):
		return END
	default:
		return NOT_A_MARKER
	}
}

// tells us whether a file actually has inline merge conflicts. We need to run this
// because git will continue showing a status of 'UU' even after the conflicts have
// been resolved in the user's editor
func FileHasConflictMarkers(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}

	defer file.Close()

	return fileHasConflictMarkersAux(file), nil
}

// Efficiently scans through a file looking for merge conflict markers. Returns true if it does
func fileHasConflictMarkersAux(file io.Reader) bool {
	scanner := bufio.NewScanner(file)
	scanner.Split(utils.ScanLinesAndTruncateWhenLongerThanBuffer(bufio.MaxScanTokenSize))
	for scanner.Scan() {
		line := scanner.Bytes()

		// only searching for start/end markers because the others are more ambiguous
		if bytes.HasPrefix(line, CONFLICT_START_BYTES) {
			return true
		}

		if bytes.HasPrefix(line, CONFLICT_END_BYTES) {
			return true
		}
	}

	return false
}
