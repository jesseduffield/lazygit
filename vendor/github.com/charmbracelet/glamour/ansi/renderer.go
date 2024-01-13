package ansi

import (
	"io"
	"net/url"
	"strings"

	"github.com/muesli/termenv"
	east "github.com/yuin/goldmark-emoji/ast"
	"github.com/yuin/goldmark/ast"
	astext "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// Options is used to configure an ANSIRenderer.
type Options struct {
	BaseURL          string
	WordWrap         int
	PreserveNewLines bool
	ColorProfile     termenv.Profile
	Styles           StyleConfig
}

// ANSIRenderer renders markdown content as ANSI escaped sequences.
type ANSIRenderer struct {
	context RenderContext
}

// NewRenderer returns a new ANSIRenderer with style and options set.
func NewRenderer(options Options) *ANSIRenderer {
	return &ANSIRenderer{
		context: NewRenderContext(options),
	}
}

// RegisterFuncs implements NodeRenderer.RegisterFuncs.
func (r *ANSIRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	// blocks
	reg.Register(ast.KindDocument, r.renderNode)
	reg.Register(ast.KindHeading, r.renderNode)
	reg.Register(ast.KindBlockquote, r.renderNode)
	reg.Register(ast.KindCodeBlock, r.renderNode)
	reg.Register(ast.KindFencedCodeBlock, r.renderNode)
	reg.Register(ast.KindHTMLBlock, r.renderNode)
	reg.Register(ast.KindList, r.renderNode)
	reg.Register(ast.KindListItem, r.renderNode)
	reg.Register(ast.KindParagraph, r.renderNode)
	reg.Register(ast.KindTextBlock, r.renderNode)
	reg.Register(ast.KindThematicBreak, r.renderNode)

	// inlines
	reg.Register(ast.KindAutoLink, r.renderNode)
	reg.Register(ast.KindCodeSpan, r.renderNode)
	reg.Register(ast.KindEmphasis, r.renderNode)
	reg.Register(ast.KindImage, r.renderNode)
	reg.Register(ast.KindLink, r.renderNode)
	reg.Register(ast.KindRawHTML, r.renderNode)
	reg.Register(ast.KindText, r.renderNode)
	reg.Register(ast.KindString, r.renderNode)

	// tables
	reg.Register(astext.KindTable, r.renderNode)
	reg.Register(astext.KindTableHeader, r.renderNode)
	reg.Register(astext.KindTableRow, r.renderNode)
	reg.Register(astext.KindTableCell, r.renderNode)

	// definitions
	reg.Register(astext.KindDefinitionList, r.renderNode)
	reg.Register(astext.KindDefinitionTerm, r.renderNode)
	reg.Register(astext.KindDefinitionDescription, r.renderNode)

	// footnotes
	reg.Register(astext.KindFootnote, r.renderNode)
	reg.Register(astext.KindFootnoteList, r.renderNode)
	reg.Register(astext.KindFootnoteLink, r.renderNode)
	reg.Register(astext.KindFootnoteBacklink, r.renderNode)

	// checkboxes
	reg.Register(astext.KindTaskCheckBox, r.renderNode)

	// strikethrough
	reg.Register(astext.KindStrikethrough, r.renderNode)

	// emoji
	reg.Register(east.KindEmoji, r.renderNode)
}

func (r *ANSIRenderer) renderNode(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	// _, _ = w.Write([]byte(node.Type.String()))
	writeTo := io.Writer(w)
	bs := r.context.blockStack

	// children get rendered by their parent
	if isChild(node) {
		return ast.WalkContinue, nil
	}

	e := r.NewElement(node, source)
	if entering {
		// everything below the Document element gets rendered into a block buffer
		if bs.Len() > 0 {
			writeTo = io.Writer(bs.Current().Block)
		}

		_, _ = writeTo.Write([]byte(e.Entering))
		if e.Renderer != nil {
			err := e.Renderer.Render(writeTo, r.context)
			if err != nil {
				return ast.WalkStop, err
			}
		}
	} else {
		// everything below the Document element gets rendered into a block buffer
		if bs.Len() > 0 {
			writeTo = io.Writer(bs.Parent().Block)
		}

		// if we're finished rendering the entire document,
		// flush to the real writer
		if node.Type() == ast.TypeDocument {
			writeTo = w
		}

		if e.Finisher != nil {
			err := e.Finisher.Finish(writeTo, r.context)
			if err != nil {
				return ast.WalkStop, err
			}
		}
		_, _ = bs.Current().Block.Write([]byte(e.Exiting))
	}

	return ast.WalkContinue, nil
}

func isChild(node ast.Node) bool {
	if node.Parent() != nil && node.Parent().Kind() == ast.KindBlockquote {
		// skip paragraph within blockquote to avoid reflowing text
		return true
	}
	for n := node.Parent(); n != nil; n = n.Parent() {
		// These types are already rendered by their parent
		switch n.Kind() {
		case ast.KindLink, ast.KindImage, ast.KindEmphasis, astext.KindStrikethrough, astext.KindTableCell:
			return true
		}
	}

	return false
}

func resolveRelativeURL(baseURL string, rel string) string {
	u, err := url.Parse(rel)
	if err != nil {
		return rel
	}
	if u.IsAbs() {
		return rel
	}
	u.Path = strings.TrimPrefix(u.Path, "/")

	base, err := url.Parse(baseURL)
	if err != nil {
		return rel
	}
	return base.ResolveReference(u).String()
}
