package utils

import (
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
)

type Alignment int

const (
	AlignLeft Alignment = iota
	AlignRight
)

type ColumnConfig struct {
	Width     int
	Alignment Alignment
}

// WithPadding pads a string as much as you want
func WithPadding(str string, padding int, alignment Alignment) string {
	uncoloredStr := Decolorise(str)
	width := runewidth.StringWidth(uncoloredStr)
	if padding < width {
		return str
	}
	space := strings.Repeat(" ", padding-width)
	if alignment == AlignLeft {
		return str + space
	} else {
		return space + str
	}
}

// defaults to left-aligning each column. If you want to set the alignment of
// each column, pass in a slice of Alignment values.
// returns a list of strings that should be joined with "\n", and an array of
// the column positions
func RenderDisplayStrings(displayStringsArr [][]string, columnAlignments []Alignment) ([]string, []int) {
	displayStringsArr, columnAlignments, removedColumns := excludeBlankColumns(displayStringsArr, columnAlignments)
	padWidths := getPadWidths(displayStringsArr)
	columnConfigs := make([]ColumnConfig, len(padWidths))
	columnPositions := make([]int, len(padWidths)+1)
	columnPositions[0] = 0
	for i, padWidth := range padWidths {
		// gracefully handle when columnAlignments is shorter than padWidths
		alignment := AlignLeft
		if len(columnAlignments) > i {
			alignment = columnAlignments[i]
		}

		columnConfigs[i] = ColumnConfig{
			Width:     padWidth,
			Alignment: alignment,
		}
		columnPositions[i+1] = columnPositions[i] + padWidth + 1
	}
	// Add the removed columns back into columnPositions (a removed column gets
	// the same position as the following column); clients should be able to rely
	// on them all to be there
	for _, removedColumn := range removedColumns {
		if removedColumn < len(columnPositions) {
			columnPositions = slices.Insert(columnPositions, removedColumn, columnPositions[removedColumn])
		}
	}
	return getPaddedDisplayStrings(displayStringsArr, columnConfigs), columnPositions
}

// NOTE: this mutates the input slice for the sake of performance
func excludeBlankColumns(displayStringsArr [][]string, columnAlignments []Alignment) ([][]string, []Alignment, []int) {
	if len(displayStringsArr) == 0 {
		return displayStringsArr, columnAlignments, []int{}
	}

	// if all rows share a blank column, we want to remove that column
	toRemove := []int{}
outer:
	for i := range displayStringsArr[0] {
		for _, strings := range displayStringsArr {
			if strings[i] != "" {
				continue outer
			}
		}
		toRemove = append(toRemove, i)
	}

	if len(toRemove) == 0 {
		return displayStringsArr, columnAlignments, []int{}
	}

	// remove the columns
	for i, strings := range displayStringsArr {
		for j := len(toRemove) - 1; j >= 0; j-- {
			strings = slices.Delete(strings, toRemove[j], toRemove[j]+1)
		}
		displayStringsArr[i] = strings
	}

	for j := len(toRemove) - 1; j >= 0; j-- {
		if columnAlignments != nil && toRemove[j] < len(columnAlignments) {
			columnAlignments = slices.Delete(columnAlignments, toRemove[j], toRemove[j]+1)
		}
	}

	return displayStringsArr, columnAlignments, toRemove
}

func getPaddedDisplayStrings(stringArrays [][]string, columnConfigs []ColumnConfig) []string {
	result := make([]string, 0, len(stringArrays))
	for _, stringArray := range stringArrays {
		if len(stringArray) == 0 {
			continue
		}
		builder := strings.Builder{}
		for j, columnConfig := range columnConfigs {
			if len(stringArray)-1 < j {
				continue
			}
			builder.WriteString(WithPadding(stringArray[j], columnConfig.Width, columnConfig.Alignment))
			builder.WriteString(" ")
		}
		if len(stringArray)-1 < len(columnConfigs) {
			continue
		}
		builder.WriteString(stringArray[len(columnConfigs)])
		result = append(result, builder.String())
	}
	return result
}

func getPadWidths(stringArrays [][]string) []int {
	maxWidth := MaxFn(stringArrays, func(stringArray []string) int {
		return len(stringArray)
	})

	if maxWidth-1 < 0 {
		return []int{}
	}
	return lo.Map(lo.Range(maxWidth-1), func(i int, _ int) int {
		return MaxFn(stringArrays, func(stringArray []string) int {
			uncoloredStr := Decolorise(stringArray[i])

			return runewidth.StringWidth(uncoloredStr)
		})
	})
}

func MaxFn[T any](items []T, fn func(T) int) int {
	max := 0
	for _, item := range items {
		if fn(item) > max {
			max = fn(item)
		}
	}
	return max
}

// TruncateWithEllipsis returns a string, truncated to a certain length, with an ellipsis
func TruncateWithEllipsis(str string, limit int) string {
	if runewidth.StringWidth(str) > limit && limit <= 3 {
		return strings.Repeat(".", limit)
	}
	return runewidth.Truncate(str, limit, "...")
}

func SafeTruncate(str string, limit int) string {
	if len(str) > limit {
		return str[0:limit]
	} else {
		return str
	}
}

const COMMIT_HASH_SHORT_SIZE = 8

func ShortSha(sha string) string {
	if len(sha) < COMMIT_HASH_SHORT_SIZE {
		return sha
	}
	return sha[:COMMIT_HASH_SHORT_SIZE]
}

// Returns comma-separated list of paths, with ellipsis if there are more than 3
// e.g. "foo, bar, baz, [...3 more]"
func FormatPaths(paths []string) string {
	if len(paths) <= 3 {
		return strings.Join(paths, ", ")
	}
	return fmt.Sprintf("%s, %s, %s, [...%d more]", paths[0], paths[1], paths[2], len(paths)-3)
}
