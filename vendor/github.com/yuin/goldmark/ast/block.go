package ast

import (
	"fmt"
	"strings"

	textm "github.com/yuin/goldmark/text"
)

// A BaseBlock struct implements the Node interface partialliy.
type BaseBlock struct {
	BaseNode
	blankPreviousLines bool
	lines              *textm.Segments
}

// Type implements Node.Type
func (b *BaseBlock) Type() NodeType {
	return TypeBlock
}

// IsRaw implements Node.IsRaw
func (b *BaseBlock) IsRaw() bool {
	return false
}

// HasBlankPreviousLines implements Node.HasBlankPreviousLines.
func (b *BaseBlock) HasBlankPreviousLines() bool {
	return b.blankPreviousLines
}

// SetBlankPreviousLines implements Node.SetBlankPreviousLines.
func (b *BaseBlock) SetBlankPreviousLines(v bool) {
	b.blankPreviousLines = v
}

// Lines implements Node.Lines
func (b *BaseBlock) Lines() *textm.Segments {
	if b.lines == nil {
		b.lines = textm.NewSegments()
	}
	return b.lines
}

// SetLines implements Node.SetLines
func (b *BaseBlock) SetLines(v *textm.Segments) {
	b.lines = v
}

// A Document struct is a root node of Markdown text.
type Document struct {
	BaseBlock

	meta map[string]interface{}
}

// KindDocument is a NodeKind of the Document node.
var KindDocument = NewNodeKind("Document")

// Dump implements Node.Dump .
func (n *Document) Dump(source []byte, level int) {
	DumpHelper(n, source, level, nil, nil)
}

// Type implements Node.Type .
func (n *Document) Type() NodeType {
	return TypeDocument
}

// Kind implements Node.Kind.
func (n *Document) Kind() NodeKind {
	return KindDocument
}

// OwnerDocument implements Node.OwnerDocument
func (n *Document) OwnerDocument() *Document {
	return n
}

// Meta returns metadata of this document.
func (n *Document) Meta() map[string]interface{} {
	if n.meta == nil {
		n.meta = map[string]interface{}{}
	}
	return n.meta
}

// SetMeta sets given metadata to this document.
func (n *Document) SetMeta(meta map[string]interface{}) {
	if n.meta == nil {
		n.meta = map[string]interface{}{}
	}
	for k, v := range meta {
		n.meta[k] = v
	}
}

// AddMeta adds given metadata to this document.
func (n *Document) AddMeta(key string, value interface{}) {
	if n.meta == nil {
		n.meta = map[string]interface{}{}
	}
	n.meta[key] = value
}

// NewDocument returns a new Document node.
func NewDocument() *Document {
	return &Document{
		BaseBlock: BaseBlock{},
		meta:      nil,
	}
}

// A TextBlock struct is a node whose lines
// should be rendered without any containers.
type TextBlock struct {
	BaseBlock
}

// Dump implements Node.Dump .
func (n *TextBlock) Dump(source []byte, level int) {
	DumpHelper(n, source, level, nil, nil)
}

// KindTextBlock is a NodeKind of the TextBlock node.
var KindTextBlock = NewNodeKind("TextBlock")

// Kind implements Node.Kind.
func (n *TextBlock) Kind() NodeKind {
	return KindTextBlock
}

// NewTextBlock returns a new TextBlock node.
func NewTextBlock() *TextBlock {
	return &TextBlock{
		BaseBlock: BaseBlock{},
	}
}

// A Paragraph struct represents a paragraph of Markdown text.
type Paragraph struct {
	BaseBlock
}

// Dump implements Node.Dump .
func (n *Paragraph) Dump(source []byte, level int) {
	DumpHelper(n, source, level, nil, nil)
}

// KindParagraph is a NodeKind of the Paragraph node.
var KindParagraph = NewNodeKind("Paragraph")

// Kind implements Node.Kind.
func (n *Paragraph) Kind() NodeKind {
	return KindParagraph
}

// NewParagraph returns a new Paragraph node.
func NewParagraph() *Paragraph {
	return &Paragraph{
		BaseBlock: BaseBlock{},
	}
}

// IsParagraph returns true if the given node implements the Paragraph interface,
// otherwise false.
func IsParagraph(node Node) bool {
	_, ok := node.(*Paragraph)
	return ok
}

// A Heading struct represents headings like SetextHeading and ATXHeading.
type Heading struct {
	BaseBlock
	// Level returns a level of this heading.
	// This value is between 1 and 6.
	Level int
}

// Dump implements Node.Dump .
func (n *Heading) Dump(source []byte, level int) {
	m := map[string]string{
		"Level": fmt.Sprintf("%d", n.Level),
	}
	DumpHelper(n, source, level, m, nil)
}

// KindHeading is a NodeKind of the Heading node.
var KindHeading = NewNodeKind("Heading")

// Kind implements Node.Kind.
func (n *Heading) Kind() NodeKind {
	return KindHeading
}

// NewHeading returns a new Heading node.
func NewHeading(level int) *Heading {
	return &Heading{
		BaseBlock: BaseBlock{},
		Level:     level,
	}
}

// A ThematicBreak struct represents a thematic break of Markdown text.
type ThematicBreak struct {
	BaseBlock
}

// Dump implements Node.Dump .
func (n *ThematicBreak) Dump(source []byte, level int) {
	DumpHelper(n, source, level, nil, nil)
}

// KindThematicBreak is a NodeKind of the ThematicBreak node.
var KindThematicBreak = NewNodeKind("ThematicBreak")

// Kind implements Node.Kind.
func (n *ThematicBreak) Kind() NodeKind {
	return KindThematicBreak
}

// NewThematicBreak returns a new ThematicBreak node.
func NewThematicBreak() *ThematicBreak {
	return &ThematicBreak{
		BaseBlock: BaseBlock{},
	}
}

// A CodeBlock interface represents an indented code block of Markdown text.
type CodeBlock struct {
	BaseBlock
}

// IsRaw implements Node.IsRaw.
func (n *CodeBlock) IsRaw() bool {
	return true
}

// Dump implements Node.Dump .
func (n *CodeBlock) Dump(source []byte, level int) {
	DumpHelper(n, source, level, nil, nil)
}

// KindCodeBlock is a NodeKind of the CodeBlock node.
var KindCodeBlock = NewNodeKind("CodeBlock")

// Kind implements Node.Kind.
func (n *CodeBlock) Kind() NodeKind {
	return KindCodeBlock
}

// NewCodeBlock returns a new CodeBlock node.
func NewCodeBlock() *CodeBlock {
	return &CodeBlock{
		BaseBlock: BaseBlock{},
	}
}

// A FencedCodeBlock struct represents a fenced code block of Markdown text.
type FencedCodeBlock struct {
	BaseBlock
	// Info returns a info text of this fenced code block.
	Info *Text

	language []byte
}

// Language returns an language in an info string.
// Language returns nil if this node does not have an info string.
func (n *FencedCodeBlock) Language(source []byte) []byte {
	if n.language == nil && n.Info != nil {
		segment := n.Info.Segment
		info := segment.Value(source)
		i := 0
		for ; i < len(info); i++ {
			if info[i] == ' ' {
				break
			}
		}
		n.language = info[:i]
	}
	return n.language
}

// IsRaw implements Node.IsRaw.
func (n *FencedCodeBlock) IsRaw() bool {
	return true
}

// Dump implements Node.Dump .
func (n *FencedCodeBlock) Dump(source []byte, level int) {
	m := map[string]string{}
	if n.Info != nil {
		m["Info"] = fmt.Sprintf("\"%s\"", n.Info.Text(source))
	}
	DumpHelper(n, source, level, m, nil)
}

// KindFencedCodeBlock is a NodeKind of the FencedCodeBlock node.
var KindFencedCodeBlock = NewNodeKind("FencedCodeBlock")

// Kind implements Node.Kind.
func (n *FencedCodeBlock) Kind() NodeKind {
	return KindFencedCodeBlock
}

// NewFencedCodeBlock return a new FencedCodeBlock node.
func NewFencedCodeBlock(info *Text) *FencedCodeBlock {
	return &FencedCodeBlock{
		BaseBlock: BaseBlock{},
		Info:      info,
	}
}

// A Blockquote struct represents an blockquote block of Markdown text.
type Blockquote struct {
	BaseBlock
}

// Dump implements Node.Dump .
func (n *Blockquote) Dump(source []byte, level int) {
	DumpHelper(n, source, level, nil, nil)
}

// KindBlockquote is a NodeKind of the Blockquote node.
var KindBlockquote = NewNodeKind("Blockquote")

// Kind implements Node.Kind.
func (n *Blockquote) Kind() NodeKind {
	return KindBlockquote
}

// NewBlockquote returns a new Blockquote node.
func NewBlockquote() *Blockquote {
	return &Blockquote{
		BaseBlock: BaseBlock{},
	}
}

// A List struct represents a list of Markdown text.
type List struct {
	BaseBlock

	// Marker is a marker character like '-', '+', ')' and '.'.
	Marker byte

	// IsTight is a true if this list is a 'tight' list.
	// See https://spec.commonmark.org/0.30/#loose for details.
	IsTight bool

	// Start is an initial number of this ordered list.
	// If this list is not an ordered list, Start is 0.
	Start int
}

// IsOrdered returns true if this list is an ordered list, otherwise false.
func (l *List) IsOrdered() bool {
	return l.Marker == '.' || l.Marker == ')'
}

// CanContinue returns true if this list can continue with
// the given mark and a list type, otherwise false.
func (l *List) CanContinue(marker byte, isOrdered bool) bool {
	return marker == l.Marker && isOrdered == l.IsOrdered()
}

// Dump implements Node.Dump.
func (l *List) Dump(source []byte, level int) {
	m := map[string]string{
		"Ordered": fmt.Sprintf("%v", l.IsOrdered()),
		"Marker":  fmt.Sprintf("%c", l.Marker),
		"Tight":   fmt.Sprintf("%v", l.IsTight),
	}
	if l.IsOrdered() {
		m["Start"] = fmt.Sprintf("%d", l.Start)
	}
	DumpHelper(l, source, level, m, nil)
}

// KindList is a NodeKind of the List node.
var KindList = NewNodeKind("List")

// Kind implements Node.Kind.
func (l *List) Kind() NodeKind {
	return KindList
}

// NewList returns a new List node.
func NewList(marker byte) *List {
	return &List{
		BaseBlock: BaseBlock{},
		Marker:    marker,
		IsTight:   true,
	}
}

// A ListItem struct represents a list item of Markdown text.
type ListItem struct {
	BaseBlock

	// Offset is an offset position of this item.
	Offset int
}

// Dump implements Node.Dump.
func (n *ListItem) Dump(source []byte, level int) {
	m := map[string]string{
		"Offset": fmt.Sprintf("%d", n.Offset),
	}
	DumpHelper(n, source, level, m, nil)
}

// KindListItem is a NodeKind of the ListItem node.
var KindListItem = NewNodeKind("ListItem")

// Kind implements Node.Kind.
func (n *ListItem) Kind() NodeKind {
	return KindListItem
}

// NewListItem returns a new ListItem node.
func NewListItem(offset int) *ListItem {
	return &ListItem{
		BaseBlock: BaseBlock{},
		Offset:    offset,
	}
}

// HTMLBlockType represents kinds of an html blocks.
// See https://spec.commonmark.org/0.30/#html-blocks
type HTMLBlockType int

const (
	// HTMLBlockType1 represents type 1 html blocks
	HTMLBlockType1 HTMLBlockType = iota + 1
	// HTMLBlockType2 represents type 2 html blocks
	HTMLBlockType2
	// HTMLBlockType3 represents type 3 html blocks
	HTMLBlockType3
	// HTMLBlockType4 represents type 4 html blocks
	HTMLBlockType4
	// HTMLBlockType5 represents type 5 html blocks
	HTMLBlockType5
	// HTMLBlockType6 represents type 6 html blocks
	HTMLBlockType6
	// HTMLBlockType7 represents type 7 html blocks
	HTMLBlockType7
)

// An HTMLBlock struct represents an html block of Markdown text.
type HTMLBlock struct {
	BaseBlock

	// Type is a type of this html block.
	HTMLBlockType HTMLBlockType

	// ClosureLine is a line that closes this html block.
	ClosureLine textm.Segment
}

// IsRaw implements Node.IsRaw.
func (n *HTMLBlock) IsRaw() bool {
	return true
}

// HasClosure returns true if this html block has a closure line,
// otherwise false.
func (n *HTMLBlock) HasClosure() bool {
	return n.ClosureLine.Start >= 0
}

// Dump implements Node.Dump.
func (n *HTMLBlock) Dump(source []byte, level int) {
	indent := strings.Repeat("    ", level)
	fmt.Printf("%s%s {\n", indent, "HTMLBlock")
	indent2 := strings.Repeat("    ", level+1)
	fmt.Printf("%sRawText: \"", indent2)
	for i := 0; i < n.Lines().Len(); i++ {
		s := n.Lines().At(i)
		fmt.Print(string(source[s.Start:s.Stop]))
	}
	fmt.Printf("\"\n")
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		c.Dump(source, level+1)
	}
	if n.HasClosure() {
		cl := n.ClosureLine
		fmt.Printf("%sClosure: \"%s\"\n", indent2, string(cl.Value(source)))
	}
	fmt.Printf("%s}\n", indent)
}

// KindHTMLBlock is a NodeKind of the HTMLBlock node.
var KindHTMLBlock = NewNodeKind("HTMLBlock")

// Kind implements Node.Kind.
func (n *HTMLBlock) Kind() NodeKind {
	return KindHTMLBlock
}

// NewHTMLBlock returns a new HTMLBlock node.
func NewHTMLBlock(typ HTMLBlockType) *HTMLBlock {
	return &HTMLBlock{
		BaseBlock:     BaseBlock{},
		HTMLBlockType: typ,
		ClosureLine:   textm.NewSegment(-1, -1),
	}
}
