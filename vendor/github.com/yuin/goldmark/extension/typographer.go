package extension

import (
	"unicode"

	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var uncloseCounterKey = parser.NewContextKey()

type unclosedCounter struct {
	Single int
	Double int
}

func (u *unclosedCounter) Reset() {
	u.Single = 0
	u.Double = 0
}

func getUnclosedCounter(pc parser.Context) *unclosedCounter {
	v := pc.Get(uncloseCounterKey)
	if v == nil {
		v = &unclosedCounter{}
		pc.Set(uncloseCounterKey, v)
	}
	return v.(*unclosedCounter)
}

// TypographicPunctuation is a key of the punctuations that can be replaced with
// typographic entities.
type TypographicPunctuation int

const (
	// LeftSingleQuote is '
	LeftSingleQuote TypographicPunctuation = iota + 1
	// RightSingleQuote is '
	RightSingleQuote
	// LeftDoubleQuote is "
	LeftDoubleQuote
	// RightDoubleQuote is "
	RightDoubleQuote
	// EnDash is --
	EnDash
	// EmDash is ---
	EmDash
	// Ellipsis is ...
	Ellipsis
	// LeftAngleQuote is <<
	LeftAngleQuote
	// RightAngleQuote is >>
	RightAngleQuote
	// Apostrophe is '
	Apostrophe

	typographicPunctuationMax
)

// An TypographerConfig struct is a data structure that holds configuration of the
// Typographer extension.
type TypographerConfig struct {
	Substitutions [][]byte
}

func newDefaultSubstitutions() [][]byte {
	replacements := make([][]byte, typographicPunctuationMax)
	replacements[LeftSingleQuote] = []byte("&lsquo;")
	replacements[RightSingleQuote] = []byte("&rsquo;")
	replacements[LeftDoubleQuote] = []byte("&ldquo;")
	replacements[RightDoubleQuote] = []byte("&rdquo;")
	replacements[EnDash] = []byte("&ndash;")
	replacements[EmDash] = []byte("&mdash;")
	replacements[Ellipsis] = []byte("&hellip;")
	replacements[LeftAngleQuote] = []byte("&laquo;")
	replacements[RightAngleQuote] = []byte("&raquo;")
	replacements[Apostrophe] = []byte("&rsquo;")

	return replacements
}

// SetOption implements SetOptioner.
func (b *TypographerConfig) SetOption(name parser.OptionName, value interface{}) {
	switch name {
	case optTypographicSubstitutions:
		b.Substitutions = value.([][]byte)
	}
}

// A TypographerOption interface sets options for the TypographerParser.
type TypographerOption interface {
	parser.Option
	SetTypographerOption(*TypographerConfig)
}

const optTypographicSubstitutions parser.OptionName = "TypographicSubstitutions"

// TypographicSubstitutions is a list of the substitutions for the Typographer extension.
type TypographicSubstitutions map[TypographicPunctuation][]byte

type withTypographicSubstitutions struct {
	value [][]byte
}

func (o *withTypographicSubstitutions) SetParserOption(c *parser.Config) {
	c.Options[optTypographicSubstitutions] = o.value
}

func (o *withTypographicSubstitutions) SetTypographerOption(p *TypographerConfig) {
	p.Substitutions = o.value
}

// WithTypographicSubstitutions is a functional otpion that specify replacement text
// for punctuations.
func WithTypographicSubstitutions(values map[TypographicPunctuation][]byte) TypographerOption {
	replacements := newDefaultSubstitutions()
	for k, v := range values {
		replacements[k] = v
	}

	return &withTypographicSubstitutions{replacements}
}

type typographerDelimiterProcessor struct {
}

func (p *typographerDelimiterProcessor) IsDelimiter(b byte) bool {
	return b == '\'' || b == '"'
}

func (p *typographerDelimiterProcessor) CanOpenCloser(opener, closer *parser.Delimiter) bool {
	return opener.Char == closer.Char
}

func (p *typographerDelimiterProcessor) OnMatch(consumes int) gast.Node {
	return nil
}

var defaultTypographerDelimiterProcessor = &typographerDelimiterProcessor{}

type typographerParser struct {
	TypographerConfig
}

// NewTypographerParser return a new InlineParser that parses
// typographer expressions.
func NewTypographerParser(opts ...TypographerOption) parser.InlineParser {
	p := &typographerParser{
		TypographerConfig: TypographerConfig{
			Substitutions: newDefaultSubstitutions(),
		},
	}
	for _, o := range opts {
		o.SetTypographerOption(&p.TypographerConfig)
	}
	return p
}

func (s *typographerParser) Trigger() []byte {
	return []byte{'\'', '"', '-', '.', ',', '<', '>', '*', '['}
}

func (s *typographerParser) Parse(parent gast.Node, block text.Reader, pc parser.Context) gast.Node {
	line, _ := block.PeekLine()
	c := line[0]
	if len(line) > 2 {
		if c == '-' {
			if s.Substitutions[EmDash] != nil && line[1] == '-' && line[2] == '-' { // ---
				node := gast.NewString(s.Substitutions[EmDash])
				node.SetCode(true)
				block.Advance(3)
				return node
			}
		} else if c == '.' {
			if s.Substitutions[Ellipsis] != nil && line[1] == '.' && line[2] == '.' { // ...
				node := gast.NewString(s.Substitutions[Ellipsis])
				node.SetCode(true)
				block.Advance(3)
				return node
			}
			return nil
		}
	}
	if len(line) > 1 {
		if c == '<' {
			if s.Substitutions[LeftAngleQuote] != nil && line[1] == '<' { // <<
				node := gast.NewString(s.Substitutions[LeftAngleQuote])
				node.SetCode(true)
				block.Advance(2)
				return node
			}
			return nil
		} else if c == '>' {
			if s.Substitutions[RightAngleQuote] != nil && line[1] == '>' { // >>
				node := gast.NewString(s.Substitutions[RightAngleQuote])
				node.SetCode(true)
				block.Advance(2)
				return node
			}
			return nil
		} else if s.Substitutions[EnDash] != nil && c == '-' && line[1] == '-' { // --
			node := gast.NewString(s.Substitutions[EnDash])
			node.SetCode(true)
			block.Advance(2)
			return node
		}
	}
	if c == '\'' || c == '"' {
		before := block.PrecendingCharacter()
		d := parser.ScanDelimiter(line, before, 1, defaultTypographerDelimiterProcessor)
		if d == nil {
			return nil
		}
		counter := getUnclosedCounter(pc)
		if c == '\'' {
			if s.Substitutions[Apostrophe] != nil {
				// Handle decade abbrevations such as '90s
				if d.CanOpen && !d.CanClose && len(line) > 3 && util.IsNumeric(line[1]) && util.IsNumeric(line[2]) && line[3] == 's' {
					after := rune(' ')
					if len(line) > 4 {
						after = util.ToRune(line, 4)
					}
					if len(line) == 3 || util.IsSpaceRune(after) || util.IsPunctRune(after) {
						node := gast.NewString(s.Substitutions[Apostrophe])
						node.SetCode(true)
						block.Advance(1)
						return node
					}
				}
				// special cases: 'twas, 'em, 'net
				if len(line) > 1 && (unicode.IsPunct(before) || unicode.IsSpace(before)) && (line[1] == 't' || line[1] == 'e' || line[1] == 'n' || line[1] == 'l') {
					node := gast.NewString(s.Substitutions[Apostrophe])
					node.SetCode(true)
					block.Advance(1)
					return node
				}
				// Convert normal apostrophes. This is probably more flexible than necessary but
				// converts any apostrophe in between two alphanumerics.
				if len(line) > 1 && (unicode.IsDigit(before) || unicode.IsLetter(before)) && (unicode.IsLetter(util.ToRune(line, 1))) {
					node := gast.NewString(s.Substitutions[Apostrophe])
					node.SetCode(true)
					block.Advance(1)
					return node
				}
			}
			if s.Substitutions[LeftSingleQuote] != nil && d.CanOpen && !d.CanClose {
				nt := LeftSingleQuote
				// special cases: Alice's, I'm, Don't, You'd
				if len(line) > 1 && (line[1] == 's' || line[1] == 'm' || line[1] == 't' || line[1] == 'd') && (len(line) < 3 || util.IsPunct(line[2]) || util.IsSpace(line[2])) {
					nt = RightSingleQuote
				}
				// special cases: I've, I'll, You're
				if len(line) > 2 && ((line[1] == 'v' && line[2] == 'e') || (line[1] == 'l' && line[2] == 'l') || (line[1] == 'r' && line[2] == 'e')) && (len(line) < 4 || util.IsPunct(line[3]) || util.IsSpace(line[3])) {
					nt = RightSingleQuote
				}
				if nt == LeftSingleQuote {
					counter.Single++
				}

				node := gast.NewString(s.Substitutions[nt])
				node.SetCode(true)
				block.Advance(1)
				return node
			}
			if s.Substitutions[RightSingleQuote] != nil {
				// plural possesives and abbreviations: Smiths', doin'
				if len(line) > 1 && unicode.IsSpace(util.ToRune(line, 0)) || unicode.IsPunct(util.ToRune(line, 0)) && (len(line) > 2 && !unicode.IsDigit(util.ToRune(line, 1))) {
					node := gast.NewString(s.Substitutions[RightSingleQuote])
					node.SetCode(true)
					block.Advance(1)
					return node
				}
			}
			if s.Substitutions[RightSingleQuote] != nil && counter.Single > 0 {
				isClose := d.CanClose && !d.CanOpen
				maybeClose := d.CanClose && d.CanOpen && len(line) > 1 && unicode.IsPunct(util.ToRune(line, 1)) && (len(line) == 2 || (len(line) > 2 && util.IsPunct(line[2]) || util.IsSpace(line[2])))
				if isClose || maybeClose {
					node := gast.NewString(s.Substitutions[RightSingleQuote])
					node.SetCode(true)
					block.Advance(1)
					counter.Single--
					return node
				}
			}
		}
		if c == '"' {
			if s.Substitutions[LeftDoubleQuote] != nil && d.CanOpen && !d.CanClose {
				node := gast.NewString(s.Substitutions[LeftDoubleQuote])
				node.SetCode(true)
				block.Advance(1)
				counter.Double++
				return node
			}
			if s.Substitutions[RightDoubleQuote] != nil && counter.Double > 0 {
				isClose := d.CanClose && !d.CanOpen
				maybeClose := d.CanClose && d.CanOpen && len(line) > 1 && (unicode.IsPunct(util.ToRune(line, 1))) && (len(line) == 2 || (len(line) > 2 && util.IsPunct(line[2]) || util.IsSpace(line[2])))
				if isClose || maybeClose {
					// special case: "Monitor 21""
					if len(line) > 1 && line[1] == '"' && unicode.IsDigit(before) {
						return nil
					}
					node := gast.NewString(s.Substitutions[RightDoubleQuote])
					node.SetCode(true)
					block.Advance(1)
					counter.Double--
					return node
				}
			}
		}
	}
	return nil
}

func (s *typographerParser) CloseBlock(parent gast.Node, pc parser.Context) {
	getUnclosedCounter(pc).Reset()
}

type typographer struct {
	options []TypographerOption
}

// Typographer is an extension that replaces punctuations with typographic entities.
var Typographer = &typographer{}

// NewTypographer returns a new Extender that replaces punctuations with typographic entities.
func NewTypographer(opts ...TypographerOption) goldmark.Extender {
	return &typographer{
		options: opts,
	}
}

func (e *typographer) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(NewTypographerParser(e.options...), 9999),
	))
}
