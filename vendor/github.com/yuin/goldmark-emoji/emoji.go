// package emoji is a extension for the goldmark(http://github.com/yuin/goldmark).
package emoji

import (
	"fmt"
	"strings"

	"github.com/yuin/goldmark"
	east "github.com/yuin/goldmark-emoji/ast"
	"github.com/yuin/goldmark-emoji/definition"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// Option interface sets options for this extension.
type Option interface {
	emojiOption()
}

// ParserConfig struct is a data structure that holds configuration of
// the Emoji extension.
type ParserConfig struct {
	Emojis definition.Emojis
}

const optEmojis parser.OptionName = "EmojiEmojis"

// SetOption implements parser.SetOptioner
func (c *ParserConfig) SetOption(name parser.OptionName, value interface{}) {
	switch name {
	case optEmojis:
		c.Emojis = value.(definition.Emojis)
	}
}

// A ParserOption interface sets options for the emoji parser.
type ParserOption interface {
	Option
	parser.Option

	SetEmojiOption(*ParserConfig)
}

var _ ParserOption = &withEmojis{}

type withEmojis struct {
	value definition.Emojis
}

func (o *withEmojis) emojiOption() {}

func (o *withEmojis) SetParserOption(c *parser.Config) {
	c.Options[optEmojis] = o.value
}

func (o *withEmojis) SetEmojiOption(c *ParserConfig) {
	c.Emojis = o.value
}

// WithMaping is a functional option that defines links names to unicode emojis.
func WithEmojis(value definition.Emojis) Option {
	return &withEmojis{
		value: value,
	}
}

// RenderingMethod indicates how emojis are rendered.
type RenderingMethod int

// RendererFunc will be used for rendering emojis.
type RendererFunc func(w util.BufWriter, source []byte, n *east.Emoji, config *RendererConfig)

const (
	// Entity renders an emoji as an html entity.
	Entity RenderingMethod = iota

	// Unicode renders an emoji as unicode character.
	Unicode

	// Twemoji renders an emoji as an img tag with [twemoji](https://github.com/twitter/twemoji).
	Twemoji

	// Func renders an emoji using RendererFunc.
	Func
)

// RendererConfig struct holds options for the emoji renderer.
type RendererConfig struct {
	html.Config

	// Method indicates how emojis are rendered.
	Method RenderingMethod

	// TwemojiTemplate is a printf template for twemoji. This value is valid only when Method is set to Twemoji.
	// `printf` arguments are:
	//
	//     1: name (e.g. "face with tears of joy")
	//     2: file name without an extension (e.g. 1f646-2642)
	//     3: '/' if XHTML, otherwise ''
	//
	TwemojiTemplate string

	// RendererFunc is a RendererFunc that renders emojis. This value is valid only when Method is set to Func.
	RendererFunc RendererFunc
}

// DefaultTwemojiTemplate is a default value for RendererConfig.TwemojiTemplate.
const DefaultTwemojiTemplate = `<img class="emoji" draggable="false" alt="%[1]s" src="https://twemoji.maxcdn.com/v/latest/72x72/%[2]s.png"%[3]s>`

// SetOption implements renderer.SetOptioner.
func (c *RendererConfig) SetOption(name renderer.OptionName, value interface{}) {
	switch name {
	case optRenderingMethod:
		c.Method = value.(RenderingMethod)
	case optTwemojiTemplate:
		c.TwemojiTemplate = value.(string)
	case optRendererFunc:
		c.RendererFunc = value.(RendererFunc)
	default:
		c.Config.SetOption(name, value)
	}
}

// A RendererOption interface sets options for the emoji renderer.
type RendererOption interface {
	Option
	renderer.Option

	SetEmojiOption(*RendererConfig)
}

var _ RendererOption = &withRenderingMethod{}

type withRenderingMethod struct {
	value RenderingMethod
}

func (o *withRenderingMethod) emojiOption() {
}

// SetConfig implements renderer.Option#SetConfig.
func (o *withRenderingMethod) SetConfig(c *renderer.Config) {
	c.Options[optRenderingMethod] = o.value
}

// SetEmojiOption implements RendererOption#SetEmojiOption
func (o *withRenderingMethod) SetEmojiOption(c *RendererConfig) {
	c.Method = o.value
}

const optRenderingMethod renderer.OptionName = "EmojiRenderingMethod"

// WithRenderingMethod is a functional option that indicates how emojis are rendered.
func WithRenderingMethod(a RenderingMethod) Option {
	return &withRenderingMethod{a}
}

type withTwemojiTemplate struct {
	value string
}

func (o *withTwemojiTemplate) emojiOption() {
}

// SetConfig implements renderer.Option#SetConfig.
func (o *withTwemojiTemplate) SetConfig(c *renderer.Config) {
	c.Options[optTwemojiTemplate] = o.value
}

// SetEmojiOption implements RendererOption#SetEmojiOption
func (o *withTwemojiTemplate) SetEmojiOption(c *RendererConfig) {
	c.TwemojiTemplate = o.value
}

const optTwemojiTemplate renderer.OptionName = "EmojiTwemojiTemplate"

// WithTwemojiTemplate is a functional option that changes a twemoji img tag.
func WithTwemojiTemplate(s string) Option {
	return &withTwemojiTemplate{s}
}

var _ RendererOption = &withRendererFunc{}

type withRendererFunc struct {
	value RendererFunc
}

func (o *withRendererFunc) emojiOption() {
}

// SetConfig implements renderer.Option#SetConfig.
func (o *withRendererFunc) SetConfig(c *renderer.Config) {
	c.Options[optRendererFunc] = o.value
}

// SetEmojiOption implements RendererOption#SetEmojiOption
func (o *withRendererFunc) SetEmojiOption(c *RendererConfig) {
	c.RendererFunc = o.value
}

const optRendererFunc renderer.OptionName = "EmojiRendererFunc"

// WithRendererFunc is a functional option that changes a renderer func.
func WithRendererFunc(f RendererFunc) Option {
	return &withRendererFunc{f}
}

type emojiParser struct {
	ParserConfig
}

// NewParser returns a new parser.InlineParser that can parse emoji expressions.
func NewParser(opts ...ParserOption) parser.InlineParser {
	p := &emojiParser{
		ParserConfig: ParserConfig{
			Emojis: definition.Github(),
		},
	}
	for _, o := range opts {
		o.SetEmojiOption(&p.ParserConfig)
	}
	return p
}

func (s *emojiParser) Trigger() []byte {
	return []byte{':'}
}

func (s *emojiParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, _ := block.PeekLine()
	if len(line) < 1 {
		return nil
	}
	i := 1
	for ; i < len(line); i++ {
		c := line[i]
		if !(util.IsAlphaNumeric(c) || c == '_' || c == '-' || c == '+') {
			break
		}
	}
	if i >= len(line) || line[i] != ':' {
		return nil
	}
	block.Advance(i + 1)
	shortName := line[1:i]
	emoji, ok := s.Emojis.Get(util.BytesToReadOnlyString(shortName))
	if !ok {
		return nil
	}
	return east.NewEmoji(shortName, emoji)
}

type emojiHTMLRenderer struct {
	RendererConfig
}

// NewHTMLRenderer returns a new HTMLRenderer.
func NewHTMLRenderer(opts ...RendererOption) renderer.NodeRenderer {
	r := &emojiHTMLRenderer{
		RendererConfig: RendererConfig{
			Config:          html.NewConfig(),
			Method:          Entity,
			TwemojiTemplate: DefaultTwemojiTemplate,
			RendererFunc:    nil,
		},
	}
	for _, opt := range opts {
		opt.SetEmojiOption(&r.RendererConfig)
	}
	return r
}

// RegisterFuncs implements renderer.NodeRenderer.RegisterFuncs.
func (r *emojiHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(east.KindEmoji, r.renderEmoji)
}

const slash = " /"
const empty = ""

func (r *emojiHTMLRenderer) renderEmoji(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	node := n.(*east.Emoji)
	if !node.Value.IsUnicode() && r.Method != Func {
		fmt.Fprintf(w, `<span title="%s">:%s:</span>`, util.EscapeHTML(util.StringToReadOnlyBytes(node.Value.Name)), node.ShortName)
		return ast.WalkContinue, nil
	}

	switch r.Method {
	case Entity:
		for _, r := range node.Value.Unicode {
			if r == 0x200D {
				_, _ = w.WriteString("&zwj;")
				continue
			}
			fmt.Fprintf(w, "&#x%x;", r)
		}
	case Unicode:
		fmt.Fprintf(w, "%s", string(node.Value.Unicode))
	case Twemoji:
		s := slash
		if !r.XHTML {
			s = empty
		}
		values := []string{}
		for _, r := range node.Value.Unicode {
			values = append(values, fmt.Sprintf("%x", r))
		}
		fmt.Fprintf(w, r.TwemojiTemplate, util.EscapeHTML(util.StringToReadOnlyBytes(node.Value.Name)), strings.Join(values, "-"), s)
	case Func:
		r.RendererFunc(w, source, node, &r.RendererConfig)
	}
	return ast.WalkContinue, nil
}

type emoji struct {
	options []Option
}

// Emoji is a goldmark.Extender implementation.
var Emoji = &emoji{
	options: []Option{},
}

// New returns a new extension with given options.
func New(opts ...Option) goldmark.Extender {
	return &emoji{
		options: opts,
	}
}

// Extend implements goldmark.Extender.
func (e *emoji) Extend(m goldmark.Markdown) {
	pOpts := []ParserOption{}
	rOpts := []RendererOption{}
	for _, o := range e.options {
		if po, ok := o.(ParserOption); ok {
			pOpts = append(pOpts, po)
			continue
		}
		if ro, ok := o.(RendererOption); ok {
			rOpts = append(rOpts, ro)
		}
	}

	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewHTMLRenderer(rOpts...), 200),
	))

	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(NewParser(pOpts...), 999),
	))

}
