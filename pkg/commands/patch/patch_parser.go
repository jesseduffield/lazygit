package patch

import (
	"regexp"
	"strings"

	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

type PatchLineKind int

const (
	PATCH_HEADER PatchLineKind = iota
	COMMIT_SHA
	COMMIT_DESCRIPTION
	HUNK_HEADER
	ADDITION
	DELETION
	CONTEXT
	NEWLINE_MESSAGE
)

// the job of this file is to parse a diff, find out where the hunks begin and end, which lines are stageable, and how to find the next hunk from the current position or the next stageable line from the current position.

type PatchLine struct {
	Kind    PatchLineKind
	Content string // something like '+ hello' (note the first character is not removed)
}

type PatchParser struct {
	Log            *logrus.Entry
	PatchLines     []*PatchLine
	PatchHunks     []*PatchHunk
	HunkStarts     []int
	StageableLines []int // rename to mention we're talking about indexes
}

// NewPatchParser builds a new branch list builder
func NewPatchParser(log *logrus.Entry, patch string) *PatchParser {
	hunkStarts, stageableLines, patchLines := parsePatch(patch)

	patchHunks := GetHunksFromDiff(patch)

	return &PatchParser{
		Log:            log,
		HunkStarts:     hunkStarts, // deprecated
		StageableLines: stageableLines,
		PatchLines:     patchLines,
		PatchHunks:     patchHunks,
	}
}

// GetHunkContainingLine takes a line index and an offset and finds the hunk
// which contains the line index, then returns the hunk considering the offset.
// e.g. if the offset is 1 it will return the next hunk.
func (p *PatchParser) GetHunkContainingLine(lineIndex int, offset int) *PatchHunk {
	if len(p.PatchHunks) == 0 {
		return nil
	}

	for index, hunk := range p.PatchHunks {
		if lineIndex >= hunk.FirstLineIdx && lineIndex <= hunk.LastLineIdx() {
			resultIndex := index + offset
			if resultIndex < 0 {
				resultIndex = 0
			} else if resultIndex > len(p.PatchHunks)-1 {
				resultIndex = len(p.PatchHunks) - 1
			}
			return p.PatchHunks[resultIndex]
		}
	}

	// if your cursor is past the last hunk, select the last hunk
	if lineIndex > p.PatchHunks[len(p.PatchHunks)-1].LastLineIdx() {
		return p.PatchHunks[len(p.PatchHunks)-1]
	}

	// otherwise select the first
	return p.PatchHunks[0]
}

// selected means you've got it highlighted with your cursor
// included means the line has been included in the patch (only applicable when
// building a patch)
func (l *PatchLine) render(selected bool, included bool) string {
	content := l.Content
	if len(content) == 0 {
		content = " " // using the space so that we can still highlight if necessary
	}

	// for hunk headers we need to start off cyan and then use white for the message
	if l.Kind == HUNK_HEADER {
		re := regexp.MustCompile("(@@.*?@@)(.*)")
		match := re.FindStringSubmatch(content)
		return coloredString(style.FgCyan, match[1], selected, included) + coloredString(theme.DefaultTextColor, match[2], selected, false)
	}

	var textStyle style.TextStyle
	switch l.Kind {
	case PATCH_HEADER:
		textStyle = textStyle.SetBold()
	case ADDITION:
		textStyle = style.FgGreen
	case DELETION:
		textStyle = style.FgRed
	case COMMIT_SHA:
		textStyle = style.FgYellow
	default:
		textStyle = theme.DefaultTextColor
	}

	return coloredString(textStyle, content, selected, included)
}

func coloredString(textStyle style.TextStyle, str string, selected bool, included bool) string {
	if selected {
		textStyle = textStyle.MergeStyle(theme.SelectedRangeBgColor)
	}

	firstCharStyle := textStyle
	if included {
		firstCharStyle = firstCharStyle.MergeStyle(style.BgGreen)
	}

	if len(str) < 2 {
		return firstCharStyle.Sprint(str)
	}

	return firstCharStyle.Sprint(str[:1]) + textStyle.Sprint(str[1:])
}

func parsePatch(patch string) ([]int, []int, []*PatchLine) {
	lines := strings.Split(patch, "\n")
	hunkStarts := []int{}
	stageableLines := []int{}
	pastFirstHunkHeader := false
	pastCommitDescription := true
	patchLines := make([]*PatchLine, len(lines))
	var lineKind PatchLineKind
	var firstChar string
	for index, line := range lines {
		firstChar = " "
		if len(line) > 0 {
			firstChar = line[:1]
		}
		if index == 0 && strings.HasPrefix(line, "commit") {
			lineKind = COMMIT_SHA
			pastCommitDescription = false
		} else if !pastCommitDescription {
			if strings.HasPrefix(line, "diff") || strings.HasPrefix(line, "---") {
				pastCommitDescription = true
				lineKind = PATCH_HEADER
			} else {
				lineKind = COMMIT_DESCRIPTION
			}
		} else if firstChar == "@" {
			pastFirstHunkHeader = true
			hunkStarts = append(hunkStarts, index)
			lineKind = HUNK_HEADER
		} else if pastFirstHunkHeader {
			switch firstChar {
			case "-":
				lineKind = DELETION
				stageableLines = append(stageableLines, index)
			case "+":
				lineKind = ADDITION
				stageableLines = append(stageableLines, index)
			case "\\":
				lineKind = NEWLINE_MESSAGE
			case " ":
				lineKind = CONTEXT
			}
		} else {
			lineKind = PATCH_HEADER
		}
		patchLines[index] = &PatchLine{Kind: lineKind, Content: line}
	}
	return hunkStarts, stageableLines, patchLines
}

// Render returns the coloured string of the diff with any selected lines highlighted
func (p *PatchParser) Render(isFocused bool, firstLineIndex int, lastLineIndex int, incLineIndices []int) string {
	contentToDisplay := slices.Some(p.PatchLines, func(line *PatchLine) bool {
		return line.Content != ""
	})
	if !contentToDisplay {
		return ""
	}

	renderedLines := slices.MapWithIndex(p.PatchLines, func(patchLine *PatchLine, index int) string {
		selected := isFocused && index >= firstLineIndex && index <= lastLineIndex
		included := lo.Contains(incLineIndices, index)
		return patchLine.render(selected, included)
	})

	result := strings.Join(renderedLines, "\n")

	return result
}

// RenderLinesPlain returns the non-coloured string of diff part from firstLineIndex to
// lastLineIndex
func (p *PatchParser) RenderLinesPlain(firstLineIndex, lastLineIndex int) string {
	return renderLinesPlain(p.PatchLines[firstLineIndex : lastLineIndex+1])
}

func renderLinesPlain(lines []*PatchLine) string {
	renderedLines := slices.Map(lines, func(line *PatchLine) string {
		return line.Content + "\n"
	})

	return strings.Join(renderedLines, "")
}

// GetNextStageableLineIndex takes a line index and returns the line index of the next stageable line
// note this will actually include the current index if it is stageable
func (p *PatchParser) GetNextStageableLineIndex(currentIndex int) int {
	for _, lineIndex := range p.StageableLines {
		if lineIndex >= currentIndex {
			return lineIndex
		}
	}
	return p.StageableLines[len(p.StageableLines)-1]
}
