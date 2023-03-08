package patch

import (
	"strings"

	"github.com/jesseduffield/generics/set"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/samber/lo"
)

type patchPresenter struct {
	patch *Patch
	// if true, all following fields are ignored
	plain bool

	isFocused bool
	// first line index for selected cursor range
	firstLineIndex int
	// last line index for selected cursor range
	lastLineIndex int
	// line indices for tagged lines (e.g. lines added to a custom patch)
	incLineIndices *set.Set[int]
}

// formats the patch as a plain string
func formatPlain(patch *Patch) string {
	presenter := &patchPresenter{
		patch:          patch,
		plain:          true,
		incLineIndices: set.New[int](),
	}
	return presenter.format()
}

func formatRangePlain(patch *Patch, startIdx int, endIdx int) string {
	lines := patch.Lines()[startIdx : endIdx+1]
	return strings.Join(
		lo.Map(lines, func(line *PatchLine, _ int) string {
			return line.Content + "\n"
		}),
		"",
	)
}

type FormatViewOpts struct {
	IsFocused bool
	// first line index for selected cursor range
	FirstLineIndex int
	// last line index for selected cursor range
	LastLineIndex int
	// line indices for tagged lines (e.g. lines added to a custom patch)
	IncLineIndices *set.Set[int]
}

// formats the patch for rendering within a view, meaning it's coloured and
// highlights selected items
func formatView(patch *Patch, opts FormatViewOpts) string {
	includedLineIndices := opts.IncLineIndices
	if includedLineIndices == nil {
		includedLineIndices = set.New[int]()
	}
	presenter := &patchPresenter{
		patch:          patch,
		plain:          false,
		isFocused:      opts.IsFocused,
		firstLineIndex: opts.FirstLineIndex,
		lastLineIndex:  opts.LastLineIndex,
		incLineIndices: includedLineIndices,
	}
	return presenter.format()
}

func (self *patchPresenter) format() string {
	// if we have no changes in our patch (i.e. no additions or deletions) then
	// the patch is effectively empty and we can return an empty string
	if !self.patch.ContainsChanges() {
		return ""
	}

	stringBuilder := &strings.Builder{}
	lineIdx := 0
	appendLine := func(line string) {
		_, _ = stringBuilder.WriteString(line + "\n")

		lineIdx++
	}
	appendFormattedLine := func(line string, style style.TextStyle) {
		formattedLine := self.formatLine(
			line,
			style,
			lineIdx,
		)

		appendLine(formattedLine)
	}

	for _, line := range self.patch.header {
		appendFormattedLine(line, theme.DefaultTextColor.SetBold())
	}

	for _, hunk := range self.patch.hunks {
		appendLine(
			self.formatLine(
				hunk.formatHeaderStart(),
				style.FgCyan,
				lineIdx,
			) +
				// we're splitting the line into two parts: the diff header and the context
				// We explicitly pass 'included' as false here so that we're only tagging the
				// first half of the line as included if the line is indeed included.
				self.formatLineAux(
					hunk.headerContext,
					theme.DefaultTextColor,
					lineIdx,
					false,
				),
		)

		for _, line := range hunk.bodyLines {
			appendFormattedLine(line.Content, self.patchLineStyle(line))
		}
	}

	return stringBuilder.String()
}

func (self *patchPresenter) patchLineStyle(patchLine *PatchLine) style.TextStyle {
	switch patchLine.Kind {
	case ADDITION:
		return style.FgGreen
	case DELETION:
		return style.FgRed
	default:
		return theme.DefaultTextColor
	}
}

func (self *patchPresenter) formatLine(str string, textStyle style.TextStyle, index int) string {
	included := self.incLineIndices.Includes(index)

	return self.formatLineAux(str, textStyle, index, included)
}

// 'selected' means you've got it highlighted with your cursor
// 'included' means the line has been included in the patch (only applicable when
// building a patch)
func (self *patchPresenter) formatLineAux(str string, textStyle style.TextStyle, index int, included bool) string {
	if self.plain {
		return str
	}

	selected := self.isFocused && index >= self.firstLineIndex && index <= self.lastLineIndex

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
