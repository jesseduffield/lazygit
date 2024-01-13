package extension

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var escapedPipeCellListKey = parser.NewContextKey()

type escapedPipeCell struct {
	Cell        *ast.TableCell
	Pos         []int
	Transformed bool
}

// TableCellAlignMethod indicates how are table cells aligned in HTML format.indicates how are table cells aligned in HTML format.
type TableCellAlignMethod int

const (
	// TableCellAlignDefault renders alignments by default method.
	// With XHTML, alignments are rendered as an align attribute.
	// With HTML5, alignments are rendered as a style attribute.
	TableCellAlignDefault TableCellAlignMethod = iota

	// TableCellAlignAttribute renders alignments as an align attribute.
	TableCellAlignAttribute

	// TableCellAlignStyle renders alignments as a style attribute.
	TableCellAlignStyle

	// TableCellAlignNone does not care about alignments.
	// If you using classes or other styles, you can add these attributes
	// in an ASTTransformer.
	TableCellAlignNone
)

// TableConfig struct holds options for the extension.
type TableConfig struct {
	html.Config

	// TableCellAlignMethod indicates how are table celss aligned.
	TableCellAlignMethod TableCellAlignMethod
}

// TableOption interface is a functional option interface for the extension.
type TableOption interface {
	renderer.Option
	// SetTableOption sets given option to the extension.
	SetTableOption(*TableConfig)
}

// NewTableConfig returns a new Config with defaults.
func NewTableConfig() TableConfig {
	return TableConfig{
		Config:               html.NewConfig(),
		TableCellAlignMethod: TableCellAlignDefault,
	}
}

// SetOption implements renderer.SetOptioner.
func (c *TableConfig) SetOption(name renderer.OptionName, value interface{}) {
	switch name {
	case optTableCellAlignMethod:
		c.TableCellAlignMethod = value.(TableCellAlignMethod)
	default:
		c.Config.SetOption(name, value)
	}
}

type withTableHTMLOptions struct {
	value []html.Option
}

func (o *withTableHTMLOptions) SetConfig(c *renderer.Config) {
	if o.value != nil {
		for _, v := range o.value {
			v.(renderer.Option).SetConfig(c)
		}
	}
}

func (o *withTableHTMLOptions) SetTableOption(c *TableConfig) {
	if o.value != nil {
		for _, v := range o.value {
			v.SetHTMLOption(&c.Config)
		}
	}
}

// WithTableHTMLOptions is functional option that wraps goldmark HTMLRenderer options.
func WithTableHTMLOptions(opts ...html.Option) TableOption {
	return &withTableHTMLOptions{opts}
}

const optTableCellAlignMethod renderer.OptionName = "TableTableCellAlignMethod"

type withTableCellAlignMethod struct {
	value TableCellAlignMethod
}

func (o *withTableCellAlignMethod) SetConfig(c *renderer.Config) {
	c.Options[optTableCellAlignMethod] = o.value
}

func (o *withTableCellAlignMethod) SetTableOption(c *TableConfig) {
	c.TableCellAlignMethod = o.value
}

// WithTableCellAlignMethod is a functional option that indicates how are table cells aligned in HTML format.
func WithTableCellAlignMethod(a TableCellAlignMethod) TableOption {
	return &withTableCellAlignMethod{a}
}

func isTableDelim(bs []byte) bool {
	for _, b := range bs {
		if !(util.IsSpace(b) || b == '-' || b == '|' || b == ':') {
			return false
		}
	}
	return true
}

var tableDelimLeft = regexp.MustCompile(`^\s*\:\-+\s*$`)
var tableDelimRight = regexp.MustCompile(`^\s*\-+\:\s*$`)
var tableDelimCenter = regexp.MustCompile(`^\s*\:\-+\:\s*$`)
var tableDelimNone = regexp.MustCompile(`^\s*\-+\s*$`)

type tableParagraphTransformer struct {
}

var defaultTableParagraphTransformer = &tableParagraphTransformer{}

// NewTableParagraphTransformer returns  a new ParagraphTransformer
// that can transform paragraphs into tables.
func NewTableParagraphTransformer() parser.ParagraphTransformer {
	return defaultTableParagraphTransformer
}

func (b *tableParagraphTransformer) Transform(node *gast.Paragraph, reader text.Reader, pc parser.Context) {
	lines := node.Lines()
	if lines.Len() < 2 {
		return
	}
	for i := 1; i < lines.Len(); i++ {
		alignments := b.parseDelimiter(lines.At(i), reader)
		if alignments == nil {
			continue
		}
		header := b.parseRow(lines.At(i-1), alignments, true, reader, pc)
		if header == nil || len(alignments) != header.ChildCount() {
			return
		}
		table := ast.NewTable()
		table.Alignments = alignments
		table.AppendChild(table, ast.NewTableHeader(header))
		for j := i + 1; j < lines.Len(); j++ {
			table.AppendChild(table, b.parseRow(lines.At(j), alignments, false, reader, pc))
		}
		node.Lines().SetSliced(0, i-1)
		node.Parent().InsertAfter(node.Parent(), node, table)
		if node.Lines().Len() == 0 {
			node.Parent().RemoveChild(node.Parent(), node)
		} else {
			last := node.Lines().At(i - 2)
			last.Stop = last.Stop - 1 // trim last newline(\n)
			node.Lines().Set(i-2, last)
		}
	}
}

func (b *tableParagraphTransformer) parseRow(segment text.Segment, alignments []ast.Alignment, isHeader bool, reader text.Reader, pc parser.Context) *ast.TableRow {
	source := reader.Source()
	line := segment.Value(source)
	pos := 0
	pos += util.TrimLeftSpaceLength(line)
	limit := len(line)
	limit -= util.TrimRightSpaceLength(line)
	row := ast.NewTableRow(alignments)
	if len(line) > 0 && line[pos] == '|' {
		pos++
	}
	if len(line) > 0 && line[limit-1] == '|' {
		limit--
	}
	i := 0
	for ; pos < limit; i++ {
		alignment := ast.AlignNone
		if i >= len(alignments) {
			if !isHeader {
				return row
			}
		} else {
			alignment = alignments[i]
		}

		var escapedCell *escapedPipeCell
		node := ast.NewTableCell()
		node.Alignment = alignment
		hasBacktick := false
		closure := pos
		for ; closure < limit; closure++ {
			if line[closure] == '`' {
				hasBacktick = true
			}
			if line[closure] == '|' {
				if closure == 0 || line[closure-1] != '\\' {
					break
				} else if hasBacktick {
					if escapedCell == nil {
						escapedCell = &escapedPipeCell{node, []int{}, false}
						escapedList := pc.ComputeIfAbsent(escapedPipeCellListKey,
							func() interface{} {
								return []*escapedPipeCell{}
							}).([]*escapedPipeCell)
						escapedList = append(escapedList, escapedCell)
						pc.Set(escapedPipeCellListKey, escapedList)
					}
					escapedCell.Pos = append(escapedCell.Pos, segment.Start+closure-1)
				}
			}
		}
		seg := text.NewSegment(segment.Start+pos, segment.Start+closure)
		seg = seg.TrimLeftSpace(source)
		seg = seg.TrimRightSpace(source)
		node.Lines().Append(seg)
		row.AppendChild(row, node)
		pos = closure + 1
	}
	for ; i < len(alignments); i++ {
		row.AppendChild(row, ast.NewTableCell())
	}
	return row
}

func (b *tableParagraphTransformer) parseDelimiter(segment text.Segment, reader text.Reader) []ast.Alignment {
	line := segment.Value(reader.Source())
	if !isTableDelim(line) {
		return nil
	}
	cols := bytes.Split(line, []byte{'|'})
	if util.IsBlank(cols[0]) {
		cols = cols[1:]
	}
	if len(cols) > 0 && util.IsBlank(cols[len(cols)-1]) {
		cols = cols[:len(cols)-1]
	}

	var alignments []ast.Alignment
	for _, col := range cols {
		if tableDelimLeft.Match(col) {
			alignments = append(alignments, ast.AlignLeft)
		} else if tableDelimRight.Match(col) {
			alignments = append(alignments, ast.AlignRight)
		} else if tableDelimCenter.Match(col) {
			alignments = append(alignments, ast.AlignCenter)
		} else if tableDelimNone.Match(col) {
			alignments = append(alignments, ast.AlignNone)
		} else {
			return nil
		}
	}
	return alignments
}

type tableASTTransformer struct {
}

var defaultTableASTTransformer = &tableASTTransformer{}

// NewTableASTTransformer returns a parser.ASTTransformer for tables.
func NewTableASTTransformer() parser.ASTTransformer {
	return defaultTableASTTransformer
}

func (a *tableASTTransformer) Transform(node *gast.Document, reader text.Reader, pc parser.Context) {
	lst := pc.Get(escapedPipeCellListKey)
	if lst == nil {
		return
	}
	pc.Set(escapedPipeCellListKey, nil)
	for _, v := range lst.([]*escapedPipeCell) {
		if v.Transformed {
			continue
		}
		_ = gast.Walk(v.Cell, func(n gast.Node, entering bool) (gast.WalkStatus, error) {
			if !entering || n.Kind() != gast.KindCodeSpan {
				return gast.WalkContinue, nil
			}

			for c := n.FirstChild(); c != nil; {
				next := c.NextSibling()
				if c.Kind() != gast.KindText {
					c = next
					continue
				}
				parent := c.Parent()
				ts := &c.(*gast.Text).Segment
				n := c
				for _, v := range lst.([]*escapedPipeCell) {
					for _, pos := range v.Pos {
						if ts.Start <= pos && pos < ts.Stop {
							segment := n.(*gast.Text).Segment
							n1 := gast.NewRawTextSegment(segment.WithStop(pos))
							n2 := gast.NewRawTextSegment(segment.WithStart(pos + 1))
							parent.InsertAfter(parent, n, n1)
							parent.InsertAfter(parent, n1, n2)
							parent.RemoveChild(parent, n)
							n = n2
							v.Transformed = true
						}
					}
				}
				c = next
			}
			return gast.WalkContinue, nil
		})
	}
}

// TableHTMLRenderer is a renderer.NodeRenderer implementation that
// renders Table nodes.
type TableHTMLRenderer struct {
	TableConfig
}

// NewTableHTMLRenderer returns a new TableHTMLRenderer.
func NewTableHTMLRenderer(opts ...TableOption) renderer.NodeRenderer {
	r := &TableHTMLRenderer{
		TableConfig: NewTableConfig(),
	}
	for _, opt := range opts {
		opt.SetTableOption(&r.TableConfig)
	}
	return r
}

// RegisterFuncs implements renderer.NodeRenderer.RegisterFuncs.
func (r *TableHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindTable, r.renderTable)
	reg.Register(ast.KindTableHeader, r.renderTableHeader)
	reg.Register(ast.KindTableRow, r.renderTableRow)
	reg.Register(ast.KindTableCell, r.renderTableCell)
}

// TableAttributeFilter defines attribute names which table elements can have.
var TableAttributeFilter = html.GlobalAttributeFilter.Extend(
	[]byte("align"),       // [Deprecated]
	[]byte("bgcolor"),     // [Deprecated]
	[]byte("border"),      // [Deprecated]
	[]byte("cellpadding"), // [Deprecated]
	[]byte("cellspacing"), // [Deprecated]
	[]byte("frame"),       // [Deprecated]
	[]byte("rules"),       // [Deprecated]
	[]byte("summary"),     // [Deprecated]
	[]byte("width"),       // [Deprecated]
)

func (r *TableHTMLRenderer) renderTable(w util.BufWriter, source []byte, n gast.Node, entering bool) (gast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<table")
		if n.Attributes() != nil {
			html.RenderAttributes(w, n, TableAttributeFilter)
		}
		_, _ = w.WriteString(">\n")
	} else {
		_, _ = w.WriteString("</table>\n")
	}
	return gast.WalkContinue, nil
}

// TableHeaderAttributeFilter defines attribute names which <thead> elements can have.
var TableHeaderAttributeFilter = html.GlobalAttributeFilter.Extend(
	[]byte("align"),   // [Deprecated since HTML4] [Obsolete since HTML5]
	[]byte("bgcolor"), // [Not Standardized]
	[]byte("char"),    // [Deprecated since HTML4] [Obsolete since HTML5]
	[]byte("charoff"), // [Deprecated since HTML4] [Obsolete since HTML5]
	[]byte("valign"),  // [Deprecated since HTML4] [Obsolete since HTML5]
)

func (r *TableHTMLRenderer) renderTableHeader(w util.BufWriter, source []byte, n gast.Node, entering bool) (gast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<thead")
		if n.Attributes() != nil {
			html.RenderAttributes(w, n, TableHeaderAttributeFilter)
		}
		_, _ = w.WriteString(">\n")
		_, _ = w.WriteString("<tr>\n") // Header <tr> has no separate handle
	} else {
		_, _ = w.WriteString("</tr>\n")
		_, _ = w.WriteString("</thead>\n")
		if n.NextSibling() != nil {
			_, _ = w.WriteString("<tbody>\n")
		}
	}
	return gast.WalkContinue, nil
}

// TableRowAttributeFilter defines attribute names which <tr> elements can have.
var TableRowAttributeFilter = html.GlobalAttributeFilter.Extend(
	[]byte("align"),   // [Obsolete since HTML5]
	[]byte("bgcolor"), // [Obsolete since HTML5]
	[]byte("char"),    // [Obsolete since HTML5]
	[]byte("charoff"), // [Obsolete since HTML5]
	[]byte("valign"),  // [Obsolete since HTML5]
)

func (r *TableHTMLRenderer) renderTableRow(w util.BufWriter, source []byte, n gast.Node, entering bool) (gast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<tr")
		if n.Attributes() != nil {
			html.RenderAttributes(w, n, TableRowAttributeFilter)
		}
		_, _ = w.WriteString(">\n")
	} else {
		_, _ = w.WriteString("</tr>\n")
		if n.Parent().LastChild() == n {
			_, _ = w.WriteString("</tbody>\n")
		}
	}
	return gast.WalkContinue, nil
}

// TableThCellAttributeFilter defines attribute names which table <th> cells can have.
var TableThCellAttributeFilter = html.GlobalAttributeFilter.Extend(
	[]byte("abbr"), // [OK] Contains a short abbreviated description of the cell's content [NOT OK in <td>]

	[]byte("align"),   // [Obsolete since HTML5]
	[]byte("axis"),    // [Obsolete since HTML5]
	[]byte("bgcolor"), // [Not Standardized]
	[]byte("char"),    // [Obsolete since HTML5]
	[]byte("charoff"), // [Obsolete since HTML5]

	[]byte("colspan"), // [OK] Number of columns that the cell is to span
	[]byte("headers"), // [OK] This attribute contains a list of space-separated strings, each corresponding to the id attribute of the <th> elements that apply to this element

	[]byte("height"), // [Deprecated since HTML4] [Obsolete since HTML5]

	[]byte("rowspan"), // [OK] Number of rows that the cell is to span
	[]byte("scope"),   // [OK] This enumerated attribute defines the cells that the header (defined in the <th>) element relates to [NOT OK in <td>]

	[]byte("valign"), // [Obsolete since HTML5]
	[]byte("width"),  // [Deprecated since HTML4] [Obsolete since HTML5]
)

// TableTdCellAttributeFilter defines attribute names which table <td> cells can have.
var TableTdCellAttributeFilter = html.GlobalAttributeFilter.Extend(
	[]byte("abbr"),    // [Obsolete since HTML5] [OK in <th>]
	[]byte("align"),   // [Obsolete since HTML5]
	[]byte("axis"),    // [Obsolete since HTML5]
	[]byte("bgcolor"), // [Not Standardized]
	[]byte("char"),    // [Obsolete since HTML5]
	[]byte("charoff"), // [Obsolete since HTML5]

	[]byte("colspan"), // [OK] Number of columns that the cell is to span
	[]byte("headers"), // [OK] This attribute contains a list of space-separated strings, each corresponding to the id attribute of the <th> elements that apply to this element

	[]byte("height"), // [Deprecated since HTML4] [Obsolete since HTML5]

	[]byte("rowspan"), // [OK] Number of rows that the cell is to span

	[]byte("scope"),  // [Obsolete since HTML5] [OK in <th>]
	[]byte("valign"), // [Obsolete since HTML5]
	[]byte("width"),  // [Deprecated since HTML4] [Obsolete since HTML5]
)

func (r *TableHTMLRenderer) renderTableCell(w util.BufWriter, source []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	n := node.(*ast.TableCell)
	tag := "td"
	if n.Parent().Kind() == ast.KindTableHeader {
		tag = "th"
	}
	if entering {
		fmt.Fprintf(w, "<%s", tag)
		if n.Alignment != ast.AlignNone {
			amethod := r.TableConfig.TableCellAlignMethod
			if amethod == TableCellAlignDefault {
				if r.Config.XHTML {
					amethod = TableCellAlignAttribute
				} else {
					amethod = TableCellAlignStyle
				}
			}
			switch amethod {
			case TableCellAlignAttribute:
				if _, ok := n.AttributeString("align"); !ok { // Skip align render if overridden
					fmt.Fprintf(w, ` align="%s"`, n.Alignment.String())
				}
			case TableCellAlignStyle:
				v, ok := n.AttributeString("style")
				var cob util.CopyOnWriteBuffer
				if ok {
					cob = util.NewCopyOnWriteBuffer(v.([]byte))
					cob.AppendByte(';')
				}
				style := fmt.Sprintf("text-align:%s", n.Alignment.String())
				cob.AppendString(style)
				n.SetAttributeString("style", cob.Bytes())
			}
		}
		if n.Attributes() != nil {
			if tag == "td" {
				html.RenderAttributes(w, n, TableTdCellAttributeFilter) // <td>
			} else {
				html.RenderAttributes(w, n, TableThCellAttributeFilter) // <th>
			}
		}
		_ = w.WriteByte('>')
	} else {
		fmt.Fprintf(w, "</%s>\n", tag)
	}
	return gast.WalkContinue, nil
}

type table struct {
	options []TableOption
}

// Table is an extension that allow you to use GFM tables .
var Table = &table{
	options: []TableOption{},
}

// NewTable returns a new extension with given options.
func NewTable(opts ...TableOption) goldmark.Extender {
	return &table{
		options: opts,
	}
}

func (e *table) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithParagraphTransformers(
			util.Prioritized(NewTableParagraphTransformer(), 200),
		),
		parser.WithASTTransformers(
			util.Prioritized(defaultTableASTTransformer, 0),
		),
	)
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewTableHTMLRenderer(e.options...), 500),
	))
}
