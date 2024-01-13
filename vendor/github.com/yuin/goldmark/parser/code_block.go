package parser

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type codeBlockParser struct {
}

// CodeBlockParser is a BlockParser implementation that parses indented code blocks.
var defaultCodeBlockParser = &codeBlockParser{}

// NewCodeBlockParser returns a new BlockParser that
// parses code blocks.
func NewCodeBlockParser() BlockParser {
	return defaultCodeBlockParser
}

func (b *codeBlockParser) Trigger() []byte {
	return nil
}

func (b *codeBlockParser) Open(parent ast.Node, reader text.Reader, pc Context) (ast.Node, State) {
	line, segment := reader.PeekLine()
	pos, padding := util.IndentPosition(line, reader.LineOffset(), 4)
	if pos < 0 || util.IsBlank(line) {
		return nil, NoChildren
	}
	node := ast.NewCodeBlock()
	reader.AdvanceAndSetPadding(pos, padding)
	_, segment = reader.PeekLine()
	// if code block line starts with a tab, keep a tab as it is.
	if segment.Padding != 0 {
		preserveLeadingTabInCodeBlock(&segment, reader, 0)
	}
	node.Lines().Append(segment)
	reader.Advance(segment.Len() - 1)
	return node, NoChildren

}

func (b *codeBlockParser) Continue(node ast.Node, reader text.Reader, pc Context) State {
	line, segment := reader.PeekLine()
	if util.IsBlank(line) {
		node.Lines().Append(segment.TrimLeftSpaceWidth(4, reader.Source()))
		return Continue | NoChildren
	}
	pos, padding := util.IndentPosition(line, reader.LineOffset(), 4)
	if pos < 0 {
		return Close
	}
	reader.AdvanceAndSetPadding(pos, padding)
	_, segment = reader.PeekLine()

	// if code block line starts with a tab, keep a tab as it is.
	if segment.Padding != 0 {
		preserveLeadingTabInCodeBlock(&segment, reader, 0)
	}

	node.Lines().Append(segment)
	reader.Advance(segment.Len() - 1)
	return Continue | NoChildren
}

func (b *codeBlockParser) Close(node ast.Node, reader text.Reader, pc Context) {
	// trim trailing blank lines
	lines := node.Lines()
	length := lines.Len() - 1
	source := reader.Source()
	for length >= 0 {
		line := lines.At(length)
		if util.IsBlank(line.Value(source)) {
			length--
		} else {
			break
		}
	}
	lines.SetSliced(0, length+1)
}

func (b *codeBlockParser) CanInterruptParagraph() bool {
	return false
}

func (b *codeBlockParser) CanAcceptIndentedLine() bool {
	return true
}

func preserveLeadingTabInCodeBlock(segment *text.Segment, reader text.Reader, indent int) {
	offsetWithPadding := reader.LineOffset() + indent
	sl, ss := reader.Position()
	reader.SetPosition(sl, text.NewSegment(ss.Start-1, ss.Stop))
	if offsetWithPadding == reader.LineOffset() {
		segment.Padding = 0
		segment.Start--
	}
	reader.SetPosition(sl, ss)
}
