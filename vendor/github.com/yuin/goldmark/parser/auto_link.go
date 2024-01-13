package parser

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type autoLinkParser struct {
}

var defaultAutoLinkParser = &autoLinkParser{}

// NewAutoLinkParser returns a new InlineParser that parses autolinks
// surrounded by '<' and '>' .
func NewAutoLinkParser() InlineParser {
	return defaultAutoLinkParser
}

func (s *autoLinkParser) Trigger() []byte {
	return []byte{'<'}
}

func (s *autoLinkParser) Parse(parent ast.Node, block text.Reader, pc Context) ast.Node {
	line, segment := block.PeekLine()
	stop := util.FindEmailIndex(line[1:])
	typ := ast.AutoLinkType(ast.AutoLinkEmail)
	if stop < 0 {
		stop = util.FindURLIndex(line[1:])
		typ = ast.AutoLinkURL
	}
	if stop < 0 {
		return nil
	}
	stop++
	if stop >= len(line) || line[stop] != '>' {
		return nil
	}
	value := ast.NewTextSegment(text.NewSegment(segment.Start+1, segment.Start+stop))
	block.Advance(stop + 1)
	return ast.NewAutoLink(typ, value)
}
