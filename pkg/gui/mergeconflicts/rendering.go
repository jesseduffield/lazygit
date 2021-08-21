package mergeconflicts

import (
	"bytes"

	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func ColoredConflictFile(content string, state *State, hasFocus bool) string {
	if len(state.conflicts) == 0 {
		return content
	}
	conflict, remainingConflicts := shiftConflict(state.conflicts)
	var outputBuffer bytes.Buffer
	for i, line := range utils.SplitLines(content) {
		textStyle := theme.DefaultTextColor
		if i == conflict.start || i == conflict.ancestor || i == conflict.target || i == conflict.end {
			textStyle = style.FgRed
		}

		if hasFocus && state.conflictIndex < len(state.conflicts) && *state.conflicts[state.conflictIndex] == *conflict && shouldHighlightLine(i, conflict, state.conflictSelection) {
			textStyle = textStyle.MergeStyle(theme.SelectedRangeBgColor).SetBold()
		}
		if i == conflict.end && len(remainingConflicts) > 0 {
			conflict, remainingConflicts = shiftConflict(remainingConflicts)
		}
		outputBuffer.WriteString(textStyle.Sprint(line) + "\n")
	}
	return outputBuffer.String()
}

func shiftConflict(conflicts []*mergeConflict) (*mergeConflict, []*mergeConflict) {
	return conflicts[0], conflicts[1:]
}

func shouldHighlightLine(index int, conflict *mergeConflict, selection Selection) bool {
	switch selection {
	case TOP:
		if conflict.ancestor >= 0 {
			return index >= conflict.start && index <= conflict.ancestor
		} else {
			return index >= conflict.start && index <= conflict.target
		}
	case MIDDLE:
		return index >= conflict.ancestor && index <= conflict.target
	case BOTTOM:
		return index >= conflict.target && index <= conflict.end
	default:
		return false
	}
}
