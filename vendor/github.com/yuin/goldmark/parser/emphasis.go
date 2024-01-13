package parser

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type emphasisDelimiterProcessor struct {
}

func (p *emphasisDelimiterProcessor) IsDelimiter(b byte) bool {
	return b == '*' || b == '_'
}

func (p *emphasisDelimiterProcessor) CanOpenCloser(opener, closer *Delimiter) bool {
	return opener.Char == closer.Char
}

func (p *emphasisDelimiterProcessor) OnMatch(consumes int) ast.Node {
	return ast.NewEmphasis(consumes)
}

var defaultEmphasisDelimiterProcessor = &emphasisDelimiterProcessor{}

type emphasisParser struct {
}

var defaultEmphasisParser = &emphasisParser{}

// NewEmphasisParser return a new InlineParser that parses emphasises.
func NewEmphasisParser() InlineParser {
	return defaultEmphasisParser
}

func (s *emphasisParser) Trigger() []byte {
	return []byte{'*', '_'}
}

func (s *emphasisParser) Parse(parent ast.Node, block text.Reader, pc Context) ast.Node {
	before := block.PrecendingCharacter()
	line, segment := block.PeekLine()
	node := ScanDelimiter(line, before, 1, defaultEmphasisDelimiterProcessor)
	if node == nil {
		return nil
	}
	node.Segment = segment.WithStop(segment.Start + node.OriginalLength)
	block.Advance(node.OriginalLength)
	pc.PushDelimiter(node)
	return node
}
