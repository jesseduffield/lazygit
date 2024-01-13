package extension

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// A CJKOption sets options for CJK support mostly for HTML based renderers.
type CJKOption func(*cjk)

// WithEastAsianLineBreaks is a functional option that indicates whether softline breaks
// between east asian wide characters should be ignored.
func WithEastAsianLineBreaks() CJKOption {
	return func(c *cjk) {
		c.EastAsianLineBreaks = true
	}
}

// WithEscapedSpace is a functional option that indicates that a '\' escaped half-space(0x20) should not be rendered.
func WithEscapedSpace() CJKOption {
	return func(c *cjk) {
		c.EscapedSpace = true
	}
}

type cjk struct {
	EastAsianLineBreaks bool
	EscapedSpace        bool
}

var CJK = NewCJK(WithEastAsianLineBreaks(), WithEscapedSpace())

// NewCJK returns a new extension with given options.
func NewCJK(opts ...CJKOption) goldmark.Extender {
	e := &cjk{}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

func (e *cjk) Extend(m goldmark.Markdown) {
	if e.EastAsianLineBreaks {
		m.Renderer().AddOptions(html.WithEastAsianLineBreaks())
	}
	if e.EscapedSpace {
		m.Renderer().AddOptions(html.WithWriter(html.NewWriter(html.WithEscapedSpace())))
		m.Parser().AddOptions(parser.WithEscapedSpace())
	}
}
