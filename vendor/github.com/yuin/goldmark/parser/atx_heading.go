package parser

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// A HeadingConfig struct is a data structure that holds configuration of the renderers related to headings.
type HeadingConfig struct {
	AutoHeadingID bool
	Attribute     bool
}

// SetOption implements SetOptioner.
func (b *HeadingConfig) SetOption(name OptionName, value interface{}) {
	switch name {
	case optAutoHeadingID:
		b.AutoHeadingID = true
	case optAttribute:
		b.Attribute = true
	}
}

// A HeadingOption interface sets options for heading parsers.
type HeadingOption interface {
	Option
	SetHeadingOption(*HeadingConfig)
}

// AutoHeadingID is an option name that enables auto IDs for headings.
const optAutoHeadingID OptionName = "AutoHeadingID"

type withAutoHeadingID struct {
}

func (o *withAutoHeadingID) SetParserOption(c *Config) {
	c.Options[optAutoHeadingID] = true
}

func (o *withAutoHeadingID) SetHeadingOption(p *HeadingConfig) {
	p.AutoHeadingID = true
}

// WithAutoHeadingID is a functional option that enables custom heading ids and
// auto generated heading ids.
func WithAutoHeadingID() HeadingOption {
	return &withAutoHeadingID{}
}

type withHeadingAttribute struct {
	Option
}

func (o *withHeadingAttribute) SetHeadingOption(p *HeadingConfig) {
	p.Attribute = true
}

// WithHeadingAttribute is a functional option that enables custom heading attributes.
func WithHeadingAttribute() HeadingOption {
	return &withHeadingAttribute{WithAttribute()}
}

type atxHeadingParser struct {
	HeadingConfig
}

// NewATXHeadingParser return a new BlockParser that can parse ATX headings.
func NewATXHeadingParser(opts ...HeadingOption) BlockParser {
	p := &atxHeadingParser{}
	for _, o := range opts {
		o.SetHeadingOption(&p.HeadingConfig)
	}
	return p
}

func (b *atxHeadingParser) Trigger() []byte {
	return []byte{'#'}
}

func (b *atxHeadingParser) Open(parent ast.Node, reader text.Reader, pc Context) (ast.Node, State) {
	line, segment := reader.PeekLine()
	pos := pc.BlockOffset()
	if pos < 0 {
		return nil, NoChildren
	}
	i := pos
	for ; i < len(line) && line[i] == '#'; i++ {
	}
	level := i - pos
	if i == pos || level > 6 {
		return nil, NoChildren
	}
	if i == len(line) { // alone '#' (without a new line character)
		return ast.NewHeading(level), NoChildren
	}
	l := util.TrimLeftSpaceLength(line[i:])
	if l == 0 {
		return nil, NoChildren
	}
	start := i + l
	if start >= len(line) {
		start = len(line) - 1
	}
	origstart := start
	stop := len(line) - util.TrimRightSpaceLength(line)

	node := ast.NewHeading(level)
	parsed := false
	if b.Attribute { // handles special case like ### heading ### {#id}
		start--
		closureClose := -1
		closureOpen := -1
		for j := start; j < stop; {
			c := line[j]
			if util.IsEscapedPunctuation(line, j) {
				j += 2
			} else if util.IsSpace(c) && j < stop-1 && line[j+1] == '#' {
				closureOpen = j + 1
				k := j + 1
				for ; k < stop && line[k] == '#'; k++ {
				}
				closureClose = k
				break
			} else {
				j++
			}
		}
		if closureClose > 0 {
			reader.Advance(closureClose)
			attrs, ok := ParseAttributes(reader)
			rest, _ := reader.PeekLine()
			parsed = ok && util.IsBlank(rest)
			if parsed {
				for _, attr := range attrs {
					node.SetAttribute(attr.Name, attr.Value)
				}
				node.Lines().Append(text.NewSegment(segment.Start+start+1-segment.Padding, segment.Start+closureOpen-segment.Padding))
			}
		}
	}
	if !parsed {
		start = origstart
		stop := len(line) - util.TrimRightSpaceLength(line)
		if stop <= start { // empty headings like '##[space]'
			stop = start
		} else {
			i = stop - 1
			for ; line[i] == '#' && i >= start; i-- {
			}
			if i != stop-1 && !util.IsSpace(line[i]) {
				i = stop - 1
			}
			i++
			stop = i
		}

		if len(util.TrimRight(line[start:stop], []byte{'#'})) != 0 { // empty heading like '### ###'
			node.Lines().Append(text.NewSegment(segment.Start+start-segment.Padding, segment.Start+stop-segment.Padding))
		}
	}
	return node, NoChildren
}

func (b *atxHeadingParser) Continue(node ast.Node, reader text.Reader, pc Context) State {
	return Close
}

func (b *atxHeadingParser) Close(node ast.Node, reader text.Reader, pc Context) {
	if b.Attribute {
		_, ok := node.AttributeString("id")
		if !ok {
			parseLastLineAttributes(node, reader, pc)
		}
	}

	if b.AutoHeadingID {
		id, ok := node.AttributeString("id")
		if !ok {
			generateAutoHeadingID(node.(*ast.Heading), reader, pc)
		} else {
			pc.IDs().Put(id.([]byte))
		}
	}
}

func (b *atxHeadingParser) CanInterruptParagraph() bool {
	return true
}

func (b *atxHeadingParser) CanAcceptIndentedLine() bool {
	return false
}

func generateAutoHeadingID(node *ast.Heading, reader text.Reader, pc Context) {
	var line []byte
	lastIndex := node.Lines().Len() - 1
	if lastIndex > -1 {
		lastLine := node.Lines().At(lastIndex)
		line = lastLine.Value(reader.Source())
	}
	headingID := pc.IDs().Generate(line, ast.KindHeading)
	node.SetAttribute(attrNameID, headingID)
}

func parseLastLineAttributes(node ast.Node, reader text.Reader, pc Context) {
	lastIndex := node.Lines().Len() - 1
	if lastIndex < 0 { // empty headings
		return
	}
	lastLine := node.Lines().At(lastIndex)
	line := lastLine.Value(reader.Source())
	lr := text.NewReader(line)
	var attrs Attributes
	var ok bool
	var start text.Segment
	var sl int
	var end text.Segment
	for {
		c := lr.Peek()
		if c == text.EOF {
			break
		}
		if c == '\\' {
			lr.Advance(1)
			if lr.Peek() == '{' {
				lr.Advance(1)
			}
			continue
		}
		if c == '{' {
			sl, start = lr.Position()
			attrs, ok = ParseAttributes(lr)
			_, end = lr.Position()
			lr.SetPosition(sl, start)
		}
		lr.Advance(1)
	}
	if ok && util.IsBlank(line[end.Start:]) {
		for _, attr := range attrs {
			node.SetAttribute(attr.Name, attr.Value)
		}
		lastLine.Stop = lastLine.Start + start.Start
		node.Lines().Set(lastIndex, lastLine)
	}
}
