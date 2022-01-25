package mergeconflicts

import (
	"bytes"

	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func ColoredConflictFile(state *State, hasFocus bool) string {
	content := state.GetContent()
	if len(state.conflicts) == 0 {
		return content
	}
	conflict, remainingConflicts := shiftConflict(state.conflicts)
	var outputBuffer bytes.Buffer
	for i, line := range utils.SplitLines(content) {
		textStyle := theme.DefaultTextColor
		if conflict.isMarkerLine(i) {
			textStyle = style.FgRed
		}

		if hasFocus && state.conflictIndex < len(state.conflicts) && *state.conflicts[state.conflictIndex] == *conflict && shouldHighlightLine(i, conflict, state.Selection()) {
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
	selectionStart, selectionEnd := selection.bounds(conflict)
	return index >= selectionStart && index <= selectionEnd
}
