package html

import (
	"bytes"
	"fmt"
	"strconv"
	"unicode/utf8"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// A Config struct has configurations for the HTML based renderers.
type Config struct {
	Writer              Writer
	HardWraps           bool
	EastAsianLineBreaks bool
	XHTML               bool
	Unsafe              bool
}

// NewConfig returns a new Config with defaults.
func NewConfig() Config {
	return Config{
		Writer:              DefaultWriter,
		HardWraps:           false,
		EastAsianLineBreaks: false,
		XHTML:               false,
		Unsafe:              false,
	}
}

// SetOption implements renderer.NodeRenderer.SetOption.
func (c *Config) SetOption(name renderer.OptionName, value interface{}) {
	switch name {
	case optHardWraps:
		c.HardWraps = value.(bool)
	case optEastAsianLineBreaks:
		c.EastAsianLineBreaks = value.(bool)
	case optXHTML:
		c.XHTML = value.(bool)
	case optUnsafe:
		c.Unsafe = value.(bool)
	case optTextWriter:
		c.Writer = value.(Writer)
	}
}

// An Option interface sets options for HTML based renderers.
type Option interface {
	SetHTMLOption(*Config)
}

// TextWriter is an option name used in WithWriter.
const optTextWriter renderer.OptionName = "Writer"

type withWriter struct {
	value Writer
}

func (o *withWriter) SetConfig(c *renderer.Config) {
	c.Options[optTextWriter] = o.value
}

func (o *withWriter) SetHTMLOption(c *Config) {
	c.Writer = o.value
}

// WithWriter is a functional option that allow you to set the given writer to
// the renderer.
func WithWriter(writer Writer) interface {
	renderer.Option
	Option
} {
	return &withWriter{writer}
}

// HardWraps is an option name used in WithHardWraps.
const optHardWraps renderer.OptionName = "HardWraps"

type withHardWraps struct {
}

func (o *withHardWraps) SetConfig(c *renderer.Config) {
	c.Options[optHardWraps] = true
}

func (o *withHardWraps) SetHTMLOption(c *Config) {
	c.HardWraps = true
}

// WithHardWraps is a functional option that indicates whether softline breaks
// should be rendered as '<br>'.
func WithHardWraps() interface {
	renderer.Option
	Option
} {
	return &withHardWraps{}
}

// EastAsianLineBreaks is an option name used in WithEastAsianLineBreaks.
const optEastAsianLineBreaks renderer.OptionName = "EastAsianLineBreaks"

type withEastAsianLineBreaks struct {
}

func (o *withEastAsianLineBreaks) SetConfig(c *renderer.Config) {
	c.Options[optEastAsianLineBreaks] = true
}

func (o *withEastAsianLineBreaks) SetHTMLOption(c *Config) {
	c.EastAsianLineBreaks = true
}

// WithEastAsianLineBreaks is a functional option that indicates whether softline breaks
// between east asian wide characters should be ignored.
func WithEastAsianLineBreaks() interface {
	renderer.Option
	Option
} {
	return &withEastAsianLineBreaks{}
}

// XHTML is an option name used in WithXHTML.
const optXHTML renderer.OptionName = "XHTML"

type withXHTML struct {
}

func (o *withXHTML) SetConfig(c *renderer.Config) {
	c.Options[optXHTML] = true
}

func (o *withXHTML) SetHTMLOption(c *Config) {
	c.XHTML = true
}

// WithXHTML is a functional option indicates that nodes should be rendered in
// xhtml instead of HTML5.
func WithXHTML() interface {
	Option
	renderer.Option
} {
	return &withXHTML{}
}

// Unsafe is an option name used in WithUnsafe.
const optUnsafe renderer.OptionName = "Unsafe"

type withUnsafe struct {
}

func (o *withUnsafe) SetConfig(c *renderer.Config) {
	c.Options[optUnsafe] = true
}

func (o *withUnsafe) SetHTMLOption(c *Config) {
	c.Unsafe = true
}

// WithUnsafe is a functional option that renders dangerous contents
// (raw htmls and potentially dangerous links) as it is.
func WithUnsafe() interface {
	renderer.Option
	Option
} {
	return &withUnsafe{}
}

// A Renderer struct is an implementation of renderer.NodeRenderer that renders
// nodes as (X)HTML.
type Renderer struct {
	Config
}

// NewRenderer returns a new Renderer with given options.
func NewRenderer(opts ...Option) renderer.NodeRenderer {
	r := &Renderer{
		Config: NewConfig(),
	}

	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

// RegisterFuncs implements NodeRenderer.RegisterFuncs .
func (r *Renderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	// blocks

	reg.Register(ast.KindDocument, r.renderDocument)
	reg.Register(ast.KindHeading, r.renderHeading)
	reg.Register(ast.KindBlockquote, r.renderBlockquote)
	reg.Register(ast.KindCodeBlock, r.renderCodeBlock)
	reg.Register(ast.KindFencedCodeBlock, r.renderFencedCodeBlock)
	reg.Register(ast.KindHTMLBlock, r.renderHTMLBlock)
	reg.Register(ast.KindList, r.renderList)
	reg.Register(ast.KindListItem, r.renderListItem)
	reg.Register(ast.KindParagraph, r.renderParagraph)
	reg.Register(ast.KindTextBlock, r.renderTextBlock)
	reg.Register(ast.KindThematicBreak, r.renderThematicBreak)

	// inlines

	reg.Register(ast.KindAutoLink, r.renderAutoLink)
	reg.Register(ast.KindCodeSpan, r.renderCodeSpan)
	reg.Register(ast.KindEmphasis, r.renderEmphasis)
	reg.Register(ast.KindImage, r.renderImage)
	reg.Register(ast.KindLink, r.renderLink)
	reg.Register(ast.KindRawHTML, r.renderRawHTML)
	reg.Register(ast.KindText, r.renderText)
	reg.Register(ast.KindString, r.renderString)
}

func (r *Renderer) writeLines(w util.BufWriter, source []byte, n ast.Node) {
	l := n.Lines().Len()
	for i := 0; i < l; i++ {
		line := n.Lines().At(i)
		r.Writer.RawWrite(w, line.Value(source))
	}
}

// GlobalAttributeFilter defines attribute names which any elements can have.
var GlobalAttributeFilter = util.NewBytesFilter(
	[]byte("accesskey"),
	[]byte("autocapitalize"),
	[]byte("autofocus"),
	[]byte("class"),
	[]byte("contenteditable"),
	[]byte("dir"),
	[]byte("draggable"),
	[]byte("enterkeyhint"),
	[]byte("hidden"),
	[]byte("id"),
	[]byte("inert"),
	[]byte("inputmode"),
	[]byte("is"),
	[]byte("itemid"),
	[]byte("itemprop"),
	[]byte("itemref"),
	[]byte("itemscope"),
	[]byte("itemtype"),
	[]byte("lang"),
	[]byte("part"),
	[]byte("slot"),
	[]byte("spellcheck"),
	[]byte("style"),
	[]byte("tabindex"),
	[]byte("title"),
	[]byte("translate"),
)

func (r *Renderer) renderDocument(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	// nothing to do
	return ast.WalkContinue, nil
}

// HeadingAttributeFilter defines attribute names which heading elements can have
var HeadingAttributeFilter = GlobalAttributeFilter

func (r *Renderer) renderHeading(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Heading)
	if entering {
		_, _ = w.WriteString("<h")
		_ = w.WriteByte("0123456"[n.Level])
		if n.Attributes() != nil {
			RenderAttributes(w, node, HeadingAttributeFilter)
		}
		_ = w.WriteByte('>')
	} else {
		_, _ = w.WriteString("</h")
		_ = w.WriteByte("0123456"[n.Level])
		_, _ = w.WriteString(">\n")
	}
	return ast.WalkContinue, nil
}

// BlockquoteAttributeFilter defines attribute names which blockquote elements can have
var BlockquoteAttributeFilter = GlobalAttributeFilter.Extend(
	[]byte("cite"),
)

func (r *Renderer) renderBlockquote(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		if n.Attributes() != nil {
			_, _ = w.WriteString("<blockquote")
			RenderAttributes(w, n, BlockquoteAttributeFilter)
			_ = w.WriteByte('>')
		} else {
			_, _ = w.WriteString("<blockquote>\n")
		}
	} else {
		_, _ = w.WriteString("</blockquote>\n")
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderCodeBlock(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<pre><code>")
		r.writeLines(w, source, n)
	} else {
		_, _ = w.WriteString("</code></pre>\n")
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderFencedCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.FencedCodeBlock)
	if entering {
		_, _ = w.WriteString("<pre><code")
		language := n.Language(source)
		if language != nil {
			_, _ = w.WriteString(" class=\"language-")
			r.Writer.Write(w, language)
			_, _ = w.WriteString("\"")
		}
		_ = w.WriteByte('>')
		r.writeLines(w, source, n)
	} else {
		_, _ = w.WriteString("</code></pre>\n")
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderHTMLBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.HTMLBlock)
	if entering {
		if r.Unsafe {
			l := n.Lines().Len()
			for i := 0; i < l; i++ {
				line := n.Lines().At(i)
				r.Writer.SecureWrite(w, line.Value(source))
			}
		} else {
			_, _ = w.WriteString("<!-- raw HTML omitted -->\n")
		}
	} else {
		if n.HasClosure() {
			if r.Unsafe {
				closure := n.ClosureLine
				r.Writer.SecureWrite(w, closure.Value(source))
			} else {
				_, _ = w.WriteString("<!-- raw HTML omitted -->\n")
			}
		}
	}
	return ast.WalkContinue, nil
}

// ListAttributeFilter defines attribute names which list elements can have.
var ListAttributeFilter = GlobalAttributeFilter.Extend(
	[]byte("start"),
	[]byte("reversed"),
	[]byte("type"),
)

func (r *Renderer) renderList(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.List)
	tag := "ul"
	if n.IsOrdered() {
		tag = "ol"
	}
	if entering {
		_ = w.WriteByte('<')
		_, _ = w.WriteString(tag)
		if n.IsOrdered() && n.Start != 1 {
			fmt.Fprintf(w, " start=\"%d\"", n.Start)
		}
		if n.Attributes() != nil {
			RenderAttributes(w, n, ListAttributeFilter)
		}
		_, _ = w.WriteString(">\n")
	} else {
		_, _ = w.WriteString("</")
		_, _ = w.WriteString(tag)
		_, _ = w.WriteString(">\n")
	}
	return ast.WalkContinue, nil
}

// ListItemAttributeFilter defines attribute names which list item elements can have.
var ListItemAttributeFilter = GlobalAttributeFilter.Extend(
	[]byte("value"),
)

func (r *Renderer) renderListItem(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		if n.Attributes() != nil {
			_, _ = w.WriteString("<li")
			RenderAttributes(w, n, ListItemAttributeFilter)
			_ = w.WriteByte('>')
		} else {
			_, _ = w.WriteString("<li>")
		}
		fc := n.FirstChild()
		if fc != nil {
			if _, ok := fc.(*ast.TextBlock); !ok {
				_ = w.WriteByte('\n')
			}
		}
	} else {
		_, _ = w.WriteString("</li>\n")
	}
	return ast.WalkContinue, nil
}

// ParagraphAttributeFilter defines attribute names which paragraph elements can have.
var ParagraphAttributeFilter = GlobalAttributeFilter

func (r *Renderer) renderParagraph(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		if n.Attributes() != nil {
			_, _ = w.WriteString("<p")
			RenderAttributes(w, n, ParagraphAttributeFilter)
			_ = w.WriteByte('>')
		} else {
			_, _ = w.WriteString("<p>")
		}
	} else {
		_, _ = w.WriteString("</p>\n")
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderTextBlock(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		if _, ok := n.NextSibling().(ast.Node); ok && n.FirstChild() != nil {
			_ = w.WriteByte('\n')
		}
	}
	return ast.WalkContinue, nil
}

// ThematicAttributeFilter defines attribute names which hr elements can have.
var ThematicAttributeFilter = GlobalAttributeFilter.Extend(
	[]byte("align"),   // [Deprecated]
	[]byte("color"),   // [Not Standardized]
	[]byte("noshade"), // [Deprecated]
	[]byte("size"),    // [Deprecated]
	[]byte("width"),   // [Deprecated]
)

func (r *Renderer) renderThematicBreak(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	_, _ = w.WriteString("<hr")
	if n.Attributes() != nil {
		RenderAttributes(w, n, ThematicAttributeFilter)
	}
	if r.XHTML {
		_, _ = w.WriteString(" />\n")
	} else {
		_, _ = w.WriteString(">\n")
	}
	return ast.WalkContinue, nil
}

// LinkAttributeFilter defines attribute names which link elements can have.
var LinkAttributeFilter = GlobalAttributeFilter.Extend(
	[]byte("download"),
	// []byte("href"),
	[]byte("hreflang"),
	[]byte("media"),
	[]byte("ping"),
	[]byte("referrerpolicy"),
	[]byte("rel"),
	[]byte("shape"),
	[]byte("target"),
)

func (r *Renderer) renderAutoLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.AutoLink)
	if !entering {
		return ast.WalkContinue, nil
	}
	_, _ = w.WriteString(`<a href="`)
	url := n.URL(source)
	label := n.Label(source)
	if n.AutoLinkType == ast.AutoLinkEmail && !bytes.HasPrefix(bytes.ToLower(url), []byte("mailto:")) {
		_, _ = w.WriteString("mailto:")
	}
	_, _ = w.Write(util.EscapeHTML(util.URLEscape(url, false)))
	if n.Attributes() != nil {
		_ = w.WriteByte('"')
		RenderAttributes(w, n, LinkAttributeFilter)
		_ = w.WriteByte('>')
	} else {
		_, _ = w.WriteString(`">`)
	}
	_, _ = w.Write(util.EscapeHTML(label))
	_, _ = w.WriteString(`</a>`)
	return ast.WalkContinue, nil
}

// CodeAttributeFilter defines attribute names which code elements can have.
var CodeAttributeFilter = GlobalAttributeFilter

func (r *Renderer) renderCodeSpan(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		if n.Attributes() != nil {
			_, _ = w.WriteString("<code")
			RenderAttributes(w, n, CodeAttributeFilter)
			_ = w.WriteByte('>')
		} else {
			_, _ = w.WriteString("<code>")
		}
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			segment := c.(*ast.Text).Segment
			value := segment.Value(source)
			if bytes.HasSuffix(value, []byte("\n")) {
				r.Writer.RawWrite(w, value[:len(value)-1])
				r.Writer.RawWrite(w, []byte(" "))
			} else {
				r.Writer.RawWrite(w, value)
			}
		}
		return ast.WalkSkipChildren, nil
	}
	_, _ = w.WriteString("</code>")
	return ast.WalkContinue, nil
}

// EmphasisAttributeFilter defines attribute names which emphasis elements can have.
var EmphasisAttributeFilter = GlobalAttributeFilter

func (r *Renderer) renderEmphasis(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Emphasis)
	tag := "em"
	if n.Level == 2 {
		tag = "strong"
	}
	if entering {
		_ = w.WriteByte('<')
		_, _ = w.WriteString(tag)
		if n.Attributes() != nil {
			RenderAttributes(w, n, EmphasisAttributeFilter)
		}
		_ = w.WriteByte('>')
	} else {
		_, _ = w.WriteString("</")
		_, _ = w.WriteString(tag)
		_ = w.WriteByte('>')
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Link)
	if entering {
		_, _ = w.WriteString("<a href=\"")
		if r.Unsafe || !IsDangerousURL(n.Destination) {
			_, _ = w.Write(util.EscapeHTML(util.URLEscape(n.Destination, true)))
		}
		_ = w.WriteByte('"')
		if n.Title != nil {
			_, _ = w.WriteString(` title="`)
			r.Writer.Write(w, n.Title)
			_ = w.WriteByte('"')
		}
		if n.Attributes() != nil {
			RenderAttributes(w, n, LinkAttributeFilter)
		}
		_ = w.WriteByte('>')
	} else {
		_, _ = w.WriteString("</a>")
	}
	return ast.WalkContinue, nil
}

// ImageAttributeFilter defines attribute names which image elements can have.
var ImageAttributeFilter = GlobalAttributeFilter.Extend(
	[]byte("align"),
	[]byte("border"),
	[]byte("crossorigin"),
	[]byte("decoding"),
	[]byte("height"),
	[]byte("importance"),
	[]byte("intrinsicsize"),
	[]byte("ismap"),
	[]byte("loading"),
	[]byte("referrerpolicy"),
	[]byte("sizes"),
	[]byte("srcset"),
	[]byte("usemap"),
	[]byte("width"),
)

func (r *Renderer) renderImage(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*ast.Image)
	_, _ = w.WriteString("<img src=\"")
	if r.Unsafe || !IsDangerousURL(n.Destination) {
		_, _ = w.Write(util.EscapeHTML(util.URLEscape(n.Destination, true)))
	}
	_, _ = w.WriteString(`" alt="`)
	_, _ = w.Write(nodeToHTMLText(n, source))
	_ = w.WriteByte('"')
	if n.Title != nil {
		_, _ = w.WriteString(` title="`)
		r.Writer.Write(w, n.Title)
		_ = w.WriteByte('"')
	}
	if n.Attributes() != nil {
		RenderAttributes(w, n, ImageAttributeFilter)
	}
	if r.XHTML {
		_, _ = w.WriteString(" />")
	} else {
		_, _ = w.WriteString(">")
	}
	return ast.WalkSkipChildren, nil
}

func (r *Renderer) renderRawHTML(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkSkipChildren, nil
	}
	if r.Unsafe {
		n := node.(*ast.RawHTML)
		l := n.Segments.Len()
		for i := 0; i < l; i++ {
			segment := n.Segments.At(i)
			_, _ = w.Write(segment.Value(source))
		}
		return ast.WalkSkipChildren, nil
	}
	_, _ = w.WriteString("<!-- raw HTML omitted -->")
	return ast.WalkSkipChildren, nil
}

func (r *Renderer) renderText(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*ast.Text)
	segment := n.Segment
	if n.IsRaw() {
		r.Writer.RawWrite(w, segment.Value(source))
	} else {
		value := segment.Value(source)
		r.Writer.Write(w, value)
		if n.HardLineBreak() || (n.SoftLineBreak() && r.HardWraps) {
			if r.XHTML {
				_, _ = w.WriteString("<br />\n")
			} else {
				_, _ = w.WriteString("<br>\n")
			}
		} else if n.SoftLineBreak() {
			if r.EastAsianLineBreaks && len(value) != 0 {
				sibling := node.NextSibling()
				if sibling != nil && sibling.Kind() == ast.KindText {
					if siblingText := sibling.(*ast.Text).Text(source); len(siblingText) != 0 {
						thisLastRune := util.ToRune(value, len(value)-1)
						siblingFirstRune, _ := utf8.DecodeRune(siblingText)
						if !(util.IsEastAsianWideRune(thisLastRune) &&
							util.IsEastAsianWideRune(siblingFirstRune)) {
							_ = w.WriteByte('\n')
						}
					}
				}
			} else {
				_ = w.WriteByte('\n')
			}
		}
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderString(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*ast.String)
	if n.IsCode() {
		_, _ = w.Write(n.Value)
	} else {
		if n.IsRaw() {
			r.Writer.RawWrite(w, n.Value)
		} else {
			r.Writer.Write(w, n.Value)
		}
	}
	return ast.WalkContinue, nil
}

var dataPrefix = []byte("data-")

// RenderAttributes renders given node's attributes.
// You can specify attribute names to render by the filter.
// If filter is nil, RenderAttributes renders all attributes.
func RenderAttributes(w util.BufWriter, node ast.Node, filter util.BytesFilter) {
	for _, attr := range node.Attributes() {
		if filter != nil && !filter.Contains(attr.Name) {
			if !bytes.HasPrefix(attr.Name, dataPrefix) {
				continue
			}
		}
		_, _ = w.WriteString(" ")
		_, _ = w.Write(attr.Name)
		_, _ = w.WriteString(`="`)
		// TODO: convert numeric values to strings
		_, _ = w.Write(util.EscapeHTML(attr.Value.([]byte)))
		_ = w.WriteByte('"')
	}
}

// A Writer interface writes textual contents to a writer.
type Writer interface {
	// Write writes the given source to writer with resolving references and unescaping
	// backslash escaped characters.
	Write(writer util.BufWriter, source []byte)

	// RawWrite writes the given source to writer without resolving references and
	// unescaping backslash escaped characters.
	RawWrite(writer util.BufWriter, source []byte)

	// SecureWrite writes the given source to writer with replacing insecure characters.
	SecureWrite(writer util.BufWriter, source []byte)
}

var replacementCharacter = []byte("\ufffd")

// A WriterConfig struct has configurations for the HTML based writers.
type WriterConfig struct {
	// EscapedSpace is an option that indicates that a '\' escaped half-space(0x20) should not be rendered.
	EscapedSpace bool
}

// A WriterOption interface sets options for HTML based writers.
type WriterOption func(*WriterConfig)

// WithEscapedSpace is a WriterOption indicates that a '\' escaped half-space(0x20) should not be rendered.
func WithEscapedSpace() WriterOption {
	return func(c *WriterConfig) {
		c.EscapedSpace = true
	}
}

type defaultWriter struct {
	WriterConfig
}

// NewWriter returns a new Writer.
func NewWriter(opts ...WriterOption) Writer {
	w := &defaultWriter{}
	for _, opt := range opts {
		opt(&w.WriterConfig)
	}
	return w
}

func escapeRune(writer util.BufWriter, r rune) {
	if r < 256 {
		v := util.EscapeHTMLByte(byte(r))
		if v != nil {
			_, _ = writer.Write(v)
			return
		}
	}
	_, _ = writer.WriteRune(util.ToValidRune(r))
}

func (d *defaultWriter) SecureWrite(writer util.BufWriter, source []byte) {
	n := 0
	l := len(source)
	for i := 0; i < l; i++ {
		if source[i] == '\u0000' {
			_, _ = writer.Write(source[i-n : i])
			n = 0
			_, _ = writer.Write(replacementCharacter)
			continue
		}
		n++
	}
	if n != 0 {
		_, _ = writer.Write(source[l-n:])
	}
}

func (d *defaultWriter) RawWrite(writer util.BufWriter, source []byte) {
	n := 0
	l := len(source)
	for i := 0; i < l; i++ {
		v := util.EscapeHTMLByte(source[i])
		if v != nil {
			_, _ = writer.Write(source[i-n : i])
			n = 0
			_, _ = writer.Write(v)
			continue
		}
		n++
	}
	if n != 0 {
		_, _ = writer.Write(source[l-n:])
	}
}

func (d *defaultWriter) Write(writer util.BufWriter, source []byte) {
	escaped := false
	var ok bool
	limit := len(source)
	n := 0
	for i := 0; i < limit; i++ {
		c := source[i]
		if escaped {
			if util.IsPunct(c) {
				d.RawWrite(writer, source[n:i-1])
				n = i
				escaped = false
				continue
			}
			if d.EscapedSpace && c == ' ' {
				d.RawWrite(writer, source[n:i-1])
				n = i + 1
				escaped = false
				continue
			}
		}
		if c == '\x00' {
			d.RawWrite(writer, source[n:i])
			d.RawWrite(writer, replacementCharacter)
			n = i + 1
			escaped = false
			continue
		}
		if c == '&' {
			pos := i
			next := i + 1
			if next < limit && source[next] == '#' {
				nnext := next + 1
				if nnext < limit {
					nc := source[nnext]
					// code point like #x22;
					if nnext < limit && nc == 'x' || nc == 'X' {
						start := nnext + 1
						i, ok = util.ReadWhile(source, [2]int{start, limit}, util.IsHexDecimal)
						if ok && i < limit && source[i] == ';' && i-start < 7 {
							v, _ := strconv.ParseUint(util.BytesToReadOnlyString(source[start:i]), 16, 32)
							d.RawWrite(writer, source[n:pos])
							n = i + 1
							escapeRune(writer, rune(v))
							continue
						}
						// code point like #1234;
					} else if nc >= '0' && nc <= '9' {
						start := nnext
						i, ok = util.ReadWhile(source, [2]int{start, limit}, util.IsNumeric)
						if ok && i < limit && i-start < 8 && source[i] == ';' {
							v, _ := strconv.ParseUint(util.BytesToReadOnlyString(source[start:i]), 10, 32)
							d.RawWrite(writer, source[n:pos])
							n = i + 1
							escapeRune(writer, rune(v))
							continue
						}
					}
				}
			} else {
				start := next
				i, ok = util.ReadWhile(source, [2]int{start, limit}, util.IsAlphaNumeric)
				// entity reference
				if ok && i < limit && source[i] == ';' {
					name := util.BytesToReadOnlyString(source[start:i])
					entity, ok := util.LookUpHTML5EntityByName(name)
					if ok {
						d.RawWrite(writer, source[n:pos])
						n = i + 1
						d.RawWrite(writer, entity.Characters)
						continue
					}
				}
			}
			i = next - 1
		}
		if c == '\\' {
			escaped = true
			continue
		}
		escaped = false
	}
	d.RawWrite(writer, source[n:])
}

// DefaultWriter is a default instance of the Writer.
var DefaultWriter = NewWriter()

var bDataImage = []byte("data:image/")
var bPng = []byte("png;")
var bGif = []byte("gif;")
var bJpeg = []byte("jpeg;")
var bWebp = []byte("webp;")
var bSvg = []byte("svg+xml;")
var bJs = []byte("javascript:")
var bVb = []byte("vbscript:")
var bFile = []byte("file:")
var bData = []byte("data:")

// IsDangerousURL returns true if the given url seems a potentially dangerous url,
// otherwise false.
func IsDangerousURL(url []byte) bool {
	if bytes.HasPrefix(url, bDataImage) && len(url) >= 11 {
		v := url[11:]
		if bytes.HasPrefix(v, bPng) || bytes.HasPrefix(v, bGif) ||
			bytes.HasPrefix(v, bJpeg) || bytes.HasPrefix(v, bWebp) ||
			bytes.HasPrefix(v, bSvg) {
			return false
		}
		return true
	}
	return bytes.HasPrefix(url, bJs) || bytes.HasPrefix(url, bVb) ||
		bytes.HasPrefix(url, bFile) || bytes.HasPrefix(url, bData)
}

func nodeToHTMLText(n ast.Node, source []byte) []byte {
	var buf bytes.Buffer
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		if s, ok := c.(*ast.String); ok && s.IsCode() {
			buf.Write(s.Text(source))
		} else if !c.HasChildren() {
			buf.Write(util.EscapeHTML(c.Text(source)))
		} else {
			buf.Write(nodeToHTMLText(c, source))
		}
	}
	return buf.Bytes()
}
