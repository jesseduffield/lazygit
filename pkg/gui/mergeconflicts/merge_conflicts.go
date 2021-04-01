package mergeconflicts

import (
	"bytes"
	"strings"

	"github.com/fatih/color"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type Selection int

const (
	TOP Selection = iota
	BOTTOM
	BOTH
)

func FindConflicts(content string) []commands.Conflict {
	conflicts := make([]commands.Conflict, 0)

	if content == "" {
		return conflicts
	}

	var newConflict commands.Conflict
	for i, line := range utils.SplitLines(content) {
		trimmedLine := strings.TrimPrefix(line, "++")
		switch trimmedLine {
		case "<<<<<<< HEAD", "<<<<<<< MERGE_HEAD", "<<<<<<< Updated upstream", "<<<<<<< ours":
			newConflict = commands.Conflict{Start: i}
		case "=======":
			newConflict.Middle = i
		default:
			if strings.HasPrefix(trimmedLine, ">>>>>>> ") {
				newConflict.End = i
				conflicts = append(conflicts, newConflict)
			}
		}

	}
	return conflicts
}

func ColoredConflictFile(content string, conflicts []commands.Conflict, conflictIndex int, conflictTop, hasFocus bool) string {
	if len(conflicts) == 0 {
		return content
	}
	conflict, remainingConflicts := shiftConflict(conflicts)
	var outputBuffer bytes.Buffer
	for i, line := range utils.SplitLines(content) {
		colourAttr := theme.DefaultTextColor
		if i == conflict.Start || i == conflict.Middle || i == conflict.End {
			colourAttr = color.FgRed
		}
		colour := color.New(colourAttr)
		if hasFocus && conflictIndex < len(conflicts) && conflicts[conflictIndex] == conflict && shouldHighlightLine(i, conflict, conflictTop) {
			colour.Add(color.Bold)
			colour.Add(theme.SelectedRangeBgColor)
		}
		if i == conflict.End && len(remainingConflicts) > 0 {
			conflict, remainingConflicts = shiftConflict(remainingConflicts)
		}
		outputBuffer.WriteString(utils.ColoredStringDirect(line, colour) + "\n")
	}
	return outputBuffer.String()
}

func IsIndexToDelete(i int, conflict commands.Conflict, selection Selection) bool {
	return i == conflict.Middle ||
		i == conflict.Start ||
		i == conflict.End ||
		selection != BOTH &&
			(selection == BOTTOM && i > conflict.Start && i < conflict.Middle) ||
		(selection == TOP && i > conflict.Middle && i < conflict.End)
}

func shiftConflict(conflicts []commands.Conflict) (commands.Conflict, []commands.Conflict) {
	return conflicts[0], conflicts[1:]
}

func shouldHighlightLine(index int, conflict commands.Conflict, top bool) bool {
	return (index >= conflict.Start && index <= conflict.Middle && top) || (index >= conflict.Middle && index <= conflict.End && !top)
}
