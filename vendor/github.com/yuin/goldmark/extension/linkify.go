package extension

import (
	"bytes"
	"regexp"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var wwwURLRegxp = regexp.MustCompile(`^www\.[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-z]+(?:[/#?][-a-zA-Z0-9@:%_\+.~#!?&/=\(\);,'">\^{}\[\]` + "`" + `]*)?`)

var urlRegexp = regexp.MustCompile(`^(?:http|https|ftp)://[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-z]+(?::\d+)?(?:[/#?][-a-zA-Z0-9@:%_+.~#$!?&/=\(\);,'">\^{}\[\]` + "`" + `]*)?`)

// An LinkifyConfig struct is a data structure that holds configuration of the
// Linkify extension.
type LinkifyConfig struct {
	AllowedProtocols [][]byte
	URLRegexp        *regexp.Regexp
	WWWRegexp        *regexp.Regexp
	EmailRegexp      *regexp.Regexp
}

const (
	optLinkifyAllowedProtocols parser.OptionName = "LinkifyAllowedProtocols"
	optLinkifyURLRegexp        parser.OptionName = "LinkifyURLRegexp"
	optLinkifyWWWRegexp        parser.OptionName = "LinkifyWWWRegexp"
	optLinkifyEmailRegexp      parser.OptionName = "LinkifyEmailRegexp"
)

// SetOption implements SetOptioner.
func (c *LinkifyConfig) SetOption(name parser.OptionName, value interface{}) {
	switch name {
	case optLinkifyAllowedProtocols:
		c.AllowedProtocols = value.([][]byte)
	case optLinkifyURLRegexp:
		c.URLRegexp = value.(*regexp.Regexp)
	case optLinkifyWWWRegexp:
		c.WWWRegexp = value.(*regexp.Regexp)
	case optLinkifyEmailRegexp:
		c.EmailRegexp = value.(*regexp.Regexp)
	}
}

// A LinkifyOption interface sets options for the LinkifyOption.
type LinkifyOption interface {
	parser.Option
	SetLinkifyOption(*LinkifyConfig)
}

type withLinkifyAllowedProtocols struct {
	value [][]byte
}

func (o *withLinkifyAllowedProtocols) SetParserOption(c *parser.Config) {
	c.Options[optLinkifyAllowedProtocols] = o.value
}

func (o *withLinkifyAllowedProtocols) SetLinkifyOption(p *LinkifyConfig) {
	p.AllowedProtocols = o.value
}

// WithLinkifyAllowedProtocols is a functional option that specify allowed
// protocols in autolinks. Each protocol must end with ':' like
// 'http:' .
func WithLinkifyAllowedProtocols(value [][]byte) LinkifyOption {
	return &withLinkifyAllowedProtocols{
		value: value,
	}
}

type withLinkifyURLRegexp struct {
	value *regexp.Regexp
}

func (o *withLinkifyURLRegexp) SetParserOption(c *parser.Config) {
	c.Options[optLinkifyURLRegexp] = o.value
}

func (o *withLinkifyURLRegexp) SetLinkifyOption(p *LinkifyConfig) {
	p.URLRegexp = o.value
}

// WithLinkifyURLRegexp is a functional option that specify
// a pattern of the URL including a protocol.
func WithLinkifyURLRegexp(value *regexp.Regexp) LinkifyOption {
	return &withLinkifyURLRegexp{
		value: value,
	}
}

// WithLinkifyWWWRegexp is a functional option that specify
// a pattern of the URL without a protocol.
// This pattern must start with 'www.' .
type withLinkifyWWWRegexp struct {
	value *regexp.Regexp
}

func (o *withLinkifyWWWRegexp) SetParserOption(c *parser.Config) {
	c.Options[optLinkifyWWWRegexp] = o.value
}

func (o *withLinkifyWWWRegexp) SetLinkifyOption(p *LinkifyConfig) {
	p.WWWRegexp = o.value
}

func WithLinkifyWWWRegexp(value *regexp.Regexp) LinkifyOption {
	return &withLinkifyWWWRegexp{
		value: value,
	}
}

// WithLinkifyWWWRegexp is a functional otpion that specify
// a pattern of the email address.
type withLinkifyEmailRegexp struct {
	value *regexp.Regexp
}

func (o *withLinkifyEmailRegexp) SetParserOption(c *parser.Config) {
	c.Options[optLinkifyEmailRegexp] = o.value
}

func (o *withLinkifyEmailRegexp) SetLinkifyOption(p *LinkifyConfig) {
	p.EmailRegexp = o.value
}

func WithLinkifyEmailRegexp(value *regexp.Regexp) LinkifyOption {
	return &withLinkifyEmailRegexp{
		value: value,
	}
}

type linkifyParser struct {
	LinkifyConfig
}

// NewLinkifyParser return a new InlineParser can parse
// text that seems like a URL.
func NewLinkifyParser(opts ...LinkifyOption) parser.InlineParser {
	p := &linkifyParser{
		LinkifyConfig: LinkifyConfig{
			AllowedProtocols: nil,
			URLRegexp:        urlRegexp,
			WWWRegexp:        wwwURLRegxp,
		},
	}
	for _, o := range opts {
		o.SetLinkifyOption(&p.LinkifyConfig)
	}
	return p
}

func (s *linkifyParser) Trigger() []byte {
	// ' ' indicates any white spaces and a line head
	return []byte{' ', '*', '_', '~', '('}
}

var (
	protoHTTP  = []byte("http:")
	protoHTTPS = []byte("https:")
	protoFTP   = []byte("ftp:")
	domainWWW  = []byte("www.")
)

func (s *linkifyParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	if pc.IsInLinkLabel() {
		return nil
	}
	line, segment := block.PeekLine()
	consumes := 0
	start := segment.Start
	c := line[0]
	// advance if current position is not a line head.
	if c == ' ' || c == '*' || c == '_' || c == '~' || c == '(' {
		consumes++
		start++
		line = line[1:]
	}

	var m []int
	var protocol []byte
	var typ ast.AutoLinkType = ast.AutoLinkURL
	if s.LinkifyConfig.AllowedProtocols == nil {
		if bytes.HasPrefix(line, protoHTTP) || bytes.HasPrefix(line, protoHTTPS) || bytes.HasPrefix(line, protoFTP) {
			m = s.LinkifyConfig.URLRegexp.FindSubmatchIndex(line)
		}
	} else {
		for _, prefix := range s.LinkifyConfig.AllowedProtocols {
			if bytes.HasPrefix(line, prefix) {
				m = s.LinkifyConfig.URLRegexp.FindSubmatchIndex(line)
				break
			}
		}
	}
	if m == nil && bytes.HasPrefix(line, domainWWW) {
		m = s.LinkifyConfig.WWWRegexp.FindSubmatchIndex(line)
		protocol = []byte("http")
	}
	if m != nil && m[0] != 0 {
		m = nil
	}
	if m != nil && m[0] == 0 {
		lastChar := line[m[1]-1]
		if lastChar == '.' {
			m[1]--
		} else if lastChar == ')' {
			closing := 0
			for i := m[1] - 1; i >= m[0]; i-- {
				if line[i] == ')' {
					closing++
				} else if line[i] == '(' {
					closing--
				}
			}
			if closing > 0 {
				m[1] -= closing
			}
		} else if lastChar == ';' {
			i := m[1] - 2
			for ; i >= m[0]; i-- {
				if util.IsAlphaNumeric(line[i]) {
					continue
				}
				break
			}
			if i != m[1]-2 {
				if line[i] == '&' {
					m[1] -= m[1] - i
				}
			}
		}
	}
	if m == nil {
		if len(line) > 0 && util.IsPunct(line[0]) {
			return nil
		}
		typ = ast.AutoLinkEmail
		stop := -1
		if s.LinkifyConfig.EmailRegexp == nil {
			stop = util.FindEmailIndex(line)
		} else {
			m := s.LinkifyConfig.EmailRegexp.FindSubmatchIndex(line)
			if m != nil && m[0] == 0 {
				stop = m[1]
			}
		}
		if stop < 0 {
			return nil
		}
		at := bytes.IndexByte(line, '@')
		m = []int{0, stop, at, stop - 1}
		if m == nil || bytes.IndexByte(line[m[2]:m[3]], '.') < 0 {
			return nil
		}
		lastChar := line[m[1]-1]
		if lastChar == '.' {
			m[1]--
		}
		if m[1] < len(line) {
			nextChar := line[m[1]]
			if nextChar == '-' || nextChar == '_' {
				return nil
			}
		}
	}
	if m == nil {
		return nil
	}
	if consumes != 0 {
		s := segment.WithStop(segment.Start + 1)
		ast.MergeOrAppendTextSegment(parent, s)
	}
	i := m[1] - 1
	for ; i > 0; i-- {
		c := line[i]
		switch c {
		case '?', '!', '.', ',', ':', '*', '_', '~':
		default:
			goto endfor
		}
	}
endfor:
	i++
	consumes += i
	block.Advance(consumes)
	n := ast.NewTextSegment(text.NewSegment(start, start+i))
	link := ast.NewAutoLink(typ, n)
	link.Protocol = protocol
	return link
}

func (s *linkifyParser) CloseBlock(parent ast.Node, pc parser.Context) {
	// nothing to do
}

type linkify struct {
	options []LinkifyOption
}

// Linkify is an extension that allow you to parse text that seems like a URL.
var Linkify = &linkify{}

func NewLinkify(opts ...LinkifyOption) goldmark.Extender {
	return &linkify{
		options: opts,
	}
}

func (e *linkify) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(NewLinkifyParser(e.options...), 999),
		),
	)
}
