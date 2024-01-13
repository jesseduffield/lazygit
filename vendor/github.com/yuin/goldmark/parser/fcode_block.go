package parser

import (
	"bytes"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type fencedCodeBlockParser struct {
}

var defaultFencedCodeBlockParser = &fencedCodeBlockParser{}

// NewFencedCodeBlockParser returns a new BlockParser that
// parses fenced code blocks.
func NewFencedCodeBlockParser() BlockParser {
	return defaultFencedCodeBlockParser
}

type fenceData struct {
	char   byte
	indent int
	length int
	node   ast.Node
}

var fencedCodeBlockInfoKey = NewContextKey()

func (b *fencedCodeBlockParser) Trigger() []byte {
	return []byte{'~', '`'}
}

func (b *fencedCodeBlockParser) Open(parent ast.Node, reader text.Reader, pc Context) (ast.Node, State) {
	line, segment := reader.PeekLine()
	pos := pc.BlockOffset()
	if pos < 0 || (line[pos] != '`' && line[pos] != '~') {
		return nil, NoChildren
	}
	findent := pos
	fenceChar := line[pos]
	i := pos
	for ; i < len(line) && line[i] == fenceChar; i++ {
	}
	oFenceLength := i - pos
	if oFenceLength < 3 {
		return nil, NoChildren
	}
	var info *ast.Text
	if i < len(line)-1 {
		rest := line[i:]
		left := util.TrimLeftSpaceLength(rest)
		right := util.TrimRightSpaceLength(rest)
		if left < len(rest)-right {
			infoStart, infoStop := segment.Start-segment.Padding+i+left, segment.Stop-right
			value := rest[left : len(rest)-right]
			if fenceChar == '`' && bytes.IndexByte(value, '`') > -1 {
				return nil, NoChildren
			} else if infoStart != infoStop {
				info = ast.NewTextSegment(text.NewSegment(infoStart, infoStop))
			}
		}
	}
	node := ast.NewFencedCodeBlock(info)
	pc.Set(fencedCodeBlockInfoKey, &fenceData{fenceChar, findent, oFenceLength, node})
	return node, NoChildren

}

func (b *fencedCodeBlockParser) Continue(node ast.Node, reader text.Reader, pc Context) State {
	line, segment := reader.PeekLine()
	fdata := pc.Get(fencedCodeBlockInfoKey).(*fenceData)

	w, pos := util.IndentWidth(line, reader.LineOffset())
	if w < 4 {
		i := pos
		for ; i < len(line) && line[i] == fdata.char; i++ {
		}
		length := i - pos
		if length >= fdata.length && util.IsBlank(line[i:]) {
			newline := 1
			if line[len(line)-1] != '\n' {
				newline = 0
			}
			reader.Advance(segment.Stop - segment.Start - newline + segment.Padding)
			return Close
		}
	}
	pos, padding := util.IndentPositionPadding(line, reader.LineOffset(), segment.Padding, fdata.indent)
	if pos < 0 {
		pos = util.FirstNonSpacePosition(line)
		if pos < 0 {
			pos = 0
		}
		padding = 0
	}
	seg := text.NewSegmentPadding(segment.Start+pos, segment.Stop, padding)
	// if code block line starts with a tab, keep a tab as it is.
	if padding != 0 {
		preserveLeadingTabInCodeBlock(&seg, reader, fdata.indent)
	}
	node.Lines().Append(seg)
	reader.AdvanceAndSetPadding(segment.Stop-segment.Start-pos-1, padding)
	return Continue | NoChildren
}

func (b *fencedCodeBlockParser) Close(node ast.Node, reader text.Reader, pc Context) {
	fdata := pc.Get(fencedCodeBlockInfoKey).(*fenceData)
	if fdata.node == node {
		pc.Set(fencedCodeBlockInfoKey, nil)
	}
}

func (b *fencedCodeBlockParser) CanInterruptParagraph() bool {
	return true
}

func (b *fencedCodeBlockParser) CanAcceptIndentedLine() bool {
	return false
}
