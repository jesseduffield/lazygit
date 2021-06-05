package mergeconflicts

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/utils"
)

// LineType tells us whether a given line is a start/middle/end marker of a conflict,
// or if it's not a marker at all
type LineType int

const (
	START LineType = iota
	MIDDLE
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
			newConflict = &mergeConflict{start: i}
		case MIDDLE:
			newConflict.middle = i
		case END:
			newConflict.end = i
			conflicts = append(conflicts, newConflict)
			// reset value to avoid any possible silent mutations in further iterations
			newConflict = nil
		default:
			// line isn't a merge conflict marker so we just continue
		}
	}

	return conflicts
}

func determineLineType(line string) LineType {
	trimmedLine := strings.TrimPrefix(line, "++")

	switch {
	case strings.HasPrefix(trimmedLine, "<<<<<<< "):
		return START
	case trimmedLine == "=======":
		return MIDDLE
	case strings.HasPrefix(trimmedLine, ">>>>>>> "):
		return END
	default:
		return NOT_A_MARKER
	}
}
