package extension

import (
	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type definitionListParser struct {
}

var defaultDefinitionListParser = &definitionListParser{}

// NewDefinitionListParser return a new parser.BlockParser that
// can parse PHP Markdown Extra Definition lists.
func NewDefinitionListParser() parser.BlockParser {
	return defaultDefinitionListParser
}

func (b *definitionListParser) Trigger() []byte {
	return []byte{':'}
}

func (b *definitionListParser) Open(parent gast.Node, reader text.Reader, pc parser.Context) (gast.Node, parser.State) {
	if _, ok := parent.(*ast.DefinitionList); ok {
		return nil, parser.NoChildren
	}
	line, _ := reader.PeekLine()
	pos := pc.BlockOffset()
	indent := pc.BlockIndent()
	if pos < 0 || line[pos] != ':' || indent != 0 {
		return nil, parser.NoChildren
	}

	last := parent.LastChild()
	// need 1 or more spaces after ':'
	w, _ := util.IndentWidth(line[pos+1:], pos+1)
	if w < 1 {
		return nil, parser.NoChildren
	}
	if w >= 8 { // starts with indented code
		w = 5
	}
	w += pos + 1 /* 1 = ':' */

	para, lastIsParagraph := last.(*gast.Paragraph)
	var list *ast.DefinitionList
	status := parser.HasChildren
	var ok bool
	if lastIsParagraph {
		list, ok = last.PreviousSibling().(*ast.DefinitionList)
		if ok { // is not first item
			list.Offset = w
			list.TemporaryParagraph = para
		} else { // is first item
			list = ast.NewDefinitionList(w, para)
			status |= parser.RequireParagraph
		}
	} else if list, ok = last.(*ast.DefinitionList); ok { // multiple description
		list.Offset = w
		list.TemporaryParagraph = nil
	} else {
		return nil, parser.NoChildren
	}

	return list, status
}

func (b *definitionListParser) Continue(node gast.Node, reader text.Reader, pc parser.Context) parser.State {
	line, _ := reader.PeekLine()
	if util.IsBlank(line) {
		return parser.Continue | parser.HasChildren
	}
	list, _ := node.(*ast.DefinitionList)
	w, _ := util.IndentWidth(line, reader.LineOffset())
	if w < list.Offset {
		return parser.Close
	}
	pos, padding := util.IndentPosition(line, reader.LineOffset(), list.Offset)
	reader.AdvanceAndSetPadding(pos, padding)
	return parser.Continue | parser.HasChildren
}

func (b *definitionListParser) Close(node gast.Node, reader text.Reader, pc parser.Context) {
	// nothing to do
}

func (b *definitionListParser) CanInterruptParagraph() bool {
	return true
}

func (b *definitionListParser) CanAcceptIndentedLine() bool {
	return false
}

type definitionDescriptionParser struct {
}

var defaultDefinitionDescriptionParser = &definitionDescriptionParser{}

// NewDefinitionDescriptionParser return a new parser.BlockParser that
// can parse definition description starts with ':'.
func NewDefinitionDescriptionParser() parser.BlockParser {
	return defaultDefinitionDescriptionParser
}

func (b *definitionDescriptionParser) Trigger() []byte {
	return []byte{':'}
}

func (b *definitionDescriptionParser) Open(parent gast.Node, reader text.Reader, pc parser.Context) (gast.Node, parser.State) {
	line, _ := reader.PeekLine()
	pos := pc.BlockOffset()
	indent := pc.BlockIndent()
	if pos < 0 || line[pos] != ':' || indent != 0 {
		return nil, parser.NoChildren
	}
	list, _ := parent.(*ast.DefinitionList)
	if list == nil {
		return nil, parser.NoChildren
	}
	para := list.TemporaryParagraph
	list.TemporaryParagraph = nil
	if para != nil {
		lines := para.Lines()
		l := lines.Len()
		for i := 0; i < l; i++ {
			term := ast.NewDefinitionTerm()
			segment := lines.At(i)
			term.Lines().Append(segment.TrimRightSpace(reader.Source()))
			list.AppendChild(list, term)
		}
		para.Parent().RemoveChild(para.Parent(), para)
	}
	cpos, padding := util.IndentPosition(line[pos+1:], pos+1, list.Offset-pos-1)
	reader.AdvanceAndSetPadding(cpos+1, padding)

	return ast.NewDefinitionDescription(), parser.HasChildren
}

func (b *definitionDescriptionParser) Continue(node gast.Node, reader text.Reader, pc parser.Context) parser.State {
	// definitionListParser detects end of the description.
	// so this method will never be called.
	return parser.Continue | parser.HasChildren
}

func (b *definitionDescriptionParser) Close(node gast.Node, reader text.Reader, pc parser.Context) {
	desc := node.(*ast.DefinitionDescription)
	desc.IsTight = !desc.HasBlankPreviousLines()
	if desc.IsTight {
		for gc := desc.FirstChild(); gc != nil; gc = gc.NextSibling() {
			paragraph, ok := gc.(*gast.Paragraph)
			if ok {
				textBlock := gast.NewTextBlock()
				textBlock.SetLines(paragraph.Lines())
				desc.ReplaceChild(desc, paragraph, textBlock)
			}
		}
	}
}

func (b *definitionDescriptionParser) CanInterruptParagraph() bool {
	return true
}

func (b *definitionDescriptionParser) CanAcceptIndentedLine() bool {
	return false
}

// DefinitionListHTMLRenderer is a renderer.NodeRenderer implementation that
// renders DefinitionList nodes.
type DefinitionListHTMLRenderer struct {
	html.Config
}

// NewDefinitionListHTMLRenderer returns a new DefinitionListHTMLRenderer.
func NewDefinitionListHTMLRenderer(opts ...html.Option) renderer.NodeRenderer {
	r := &DefinitionListHTMLRenderer{
		Config: html.NewConfig(),
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

// RegisterFuncs implements renderer.NodeRenderer.RegisterFuncs.
func (r *DefinitionListHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindDefinitionList, r.renderDefinitionList)
	reg.Register(ast.KindDefinitionTerm, r.renderDefinitionTerm)
	reg.Register(ast.KindDefinitionDescription, r.renderDefinitionDescription)
}

// DefinitionListAttributeFilter defines attribute names which dl elements can have.
var DefinitionListAttributeFilter = html.GlobalAttributeFilter

func (r *DefinitionListHTMLRenderer) renderDefinitionList(w util.BufWriter, source []byte, n gast.Node, entering bool) (gast.WalkStatus, error) {
	if entering {
		if n.Attributes() != nil {
			_, _ = w.WriteString("<dl")
			html.RenderAttributes(w, n, DefinitionListAttributeFilter)
			_, _ = w.WriteString(">\n")
		} else {
			_, _ = w.WriteString("<dl>\n")
		}
	} else {
		_, _ = w.WriteString("</dl>\n")
	}
	return gast.WalkContinue, nil
}

// DefinitionTermAttributeFilter defines attribute names which dd elements can have.
var DefinitionTermAttributeFilter = html.GlobalAttributeFilter

func (r *DefinitionListHTMLRenderer) renderDefinitionTerm(w util.BufWriter, source []byte, n gast.Node, entering bool) (gast.WalkStatus, error) {
	if entering {
		if n.Attributes() != nil {
			_, _ = w.WriteString("<dt")
			html.RenderAttributes(w, n, DefinitionTermAttributeFilter)
			_ = w.WriteByte('>')
		} else {
			_, _ = w.WriteString("<dt>")
		}
	} else {
		_, _ = w.WriteString("</dt>\n")
	}
	return gast.WalkContinue, nil
}

// DefinitionDescriptionAttributeFilter defines attribute names which dd elements can have.
var DefinitionDescriptionAttributeFilter = html.GlobalAttributeFilter

func (r *DefinitionListHTMLRenderer) renderDefinitionDescription(w util.BufWriter, source []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	if entering {
		n := node.(*ast.DefinitionDescription)
		_, _ = w.WriteString("<dd")
		if n.Attributes() != nil {
			html.RenderAttributes(w, n, DefinitionDescriptionAttributeFilter)
		}
		if n.IsTight {
			_, _ = w.WriteString(">")
		} else {
			_, _ = w.WriteString(">\n")
		}
	} else {
		_, _ = w.WriteString("</dd>\n")
	}
	return gast.WalkContinue, nil
}

type definitionList struct {
}

// DefinitionList is an extension that allow you to use PHP Markdown Extra Definition lists.
var DefinitionList = &definitionList{}

func (e *definitionList) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithBlockParsers(
		util.Prioritized(NewDefinitionListParser(), 101),
		util.Prioritized(NewDefinitionDescriptionParser(), 102),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewDefinitionListHTMLRenderer(), 500),
	))
}
