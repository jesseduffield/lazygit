package parser

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type paragraphParser struct {
}

var defaultParagraphParser = &paragraphParser{}

// NewParagraphParser returns a new BlockParser that
// parses paragraphs.
func NewParagraphParser() BlockParser {
	return defaultParagraphParser
}

func (b *paragraphParser) Trigger() []byte {
	return nil
}

func (b *paragraphParser) Open(parent ast.Node, reader text.Reader, pc Context) (ast.Node, State) {
	_, segment := reader.PeekLine()
	segment = segment.TrimLeftSpace(reader.Source())
	if segment.IsEmpty() {
		return nil, NoChildren
	}
	node := ast.NewParagraph()
	node.Lines().Append(segment)
	reader.Advance(segment.Len() - 1)
	return node, NoChildren
}

func (b *paragraphParser) Continue(node ast.Node, reader text.Reader, pc Context) State {
	_, segment := reader.PeekLine()
	segment = segment.TrimLeftSpace(reader.Source())
	if segment.IsEmpty() {
		return Close
	}
	node.Lines().Append(segment)
	reader.Advance(segment.Len() - 1)
	return Continue | NoChildren
}

func (b *paragraphParser) Close(node ast.Node, reader text.Reader, pc Context) {
	parent := node.Parent()
	if parent == nil {
		// paragraph has been transformed
		return
	}
	lines := node.Lines()
	if lines.Len() != 0 {
		// trim trailing spaces
		length := lines.Len()
		lastLine := node.Lines().At(length - 1)
		node.Lines().Set(length-1, lastLine.TrimRightSpace(reader.Source()))
	}
	if lines.Len() == 0 {
		node.Parent().RemoveChild(node.Parent(), node)
		return
	}
}

func (b *paragraphParser) CanInterruptParagraph() bool {
	return false
}

func (b *paragraphParser) CanAcceptIndentedLine() bool {
	return false
}
