package controllers

import (
	"strings"

	"github.com/samber/lo"
)

// Removes '+' or '-' from the beginning of each line in the diff string, except
// when both '+' and '-' lines are present, or diff header lines, in which case
// the diff is returned unchanged. This is useful for copying parts of diffs to
// the clipboard in order to paste them into code.
func dropDiffPrefix(diff string) string {
	lines := strings.Split(strings.TrimRight(diff, "\n"), "\n")

	const (
		PLUS int = iota
		MINUS
		CONTEXT
		OTHER
	)

	linesByType := lo.GroupBy(lines, func(line string) int {
		switch {
		case strings.HasPrefix(line, "+"):
			return PLUS
		case strings.HasPrefix(line, "-"):
			return MINUS
		case strings.HasPrefix(line, " "):
			return CONTEXT
		}
		return OTHER
	})

	hasLinesOfType := func(lineType int) bool { return len(linesByType[lineType]) > 0 }

	keepPrefix := hasLinesOfType(OTHER) || (hasLinesOfType(PLUS) && hasLinesOfType(MINUS))
	if keepPrefix {
		return diff
	}

	return strings.Join(lo.Map(lines, func(line string, _ int) string { return line[1:] + "\n" }), "")
}
