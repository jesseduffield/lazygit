package ast

import (
	"fmt"
	"strings"

	textm "github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// A BaseInline struct implements the Node interface partialliy.
type BaseInline struct {
	BaseNode
}

// Type implements Node.Type
func (b *BaseInline) Type() NodeType {
	return TypeInline
}

// IsRaw implements Node.IsRaw
func (b *BaseInline) IsRaw() bool {
	return false
}

// HasBlankPreviousLines implements Node.HasBlankPreviousLines.
func (b *BaseInline) HasBlankPreviousLines() bool {
	panic("can not call with inline nodes.")
}

// SetBlankPreviousLines implements Node.SetBlankPreviousLines.
func (b *BaseInline) SetBlankPreviousLines(v bool) {
	panic("can not call with inline nodes.")
}

// Lines implements Node.Lines
func (b *BaseInline) Lines() *textm.Segments {
	panic("can not call with inline nodes.")
}

// SetLines implements Node.SetLines
func (b *BaseInline) SetLines(v *textm.Segments) {
	panic("can not call with inline nodes.")
}

// A Text struct represents a textual content of the Markdown text.
type Text struct {
	BaseInline
	// Segment is a position in a source text.
	Segment textm.Segment

	flags uint8
}

const (
	textSoftLineBreak = 1 << iota
	textHardLineBreak
	textRaw
	textCode
)

func textFlagsString(flags uint8) string {
	buf := []string{}
	if flags&textSoftLineBreak != 0 {
		buf = append(buf, "SoftLineBreak")
	}
	if flags&textHardLineBreak != 0 {
		buf = append(buf, "HardLineBreak")
	}
	if flags&textRaw != 0 {
		buf = append(buf, "Raw")
	}
	if flags&textCode != 0 {
		buf = append(buf, "Code")
	}
	return strings.Join(buf, ", ")
}

// Inline implements Inline.Inline.
func (n *Text) Inline() {
}

// SoftLineBreak returns true if this node ends with a new line,
// otherwise false.
func (n *Text) SoftLineBreak() bool {
	return n.flags&textSoftLineBreak != 0
}

// SetSoftLineBreak sets whether this node ends with a new line.
func (n *Text) SetSoftLineBreak(v bool) {
	if v {
		n.flags |= textSoftLineBreak
	} else {
		n.flags = n.flags &^ textSoftLineBreak
	}
}

// IsRaw returns true if this text should be rendered without unescaping
// back slash escapes and resolving references.
func (n *Text) IsRaw() bool {
	return n.flags&textRaw != 0
}

// SetRaw sets whether this text should be rendered as raw contents.
func (n *Text) SetRaw(v bool) {
	if v {
		n.flags |= textRaw
	} else {
		n.flags = n.flags &^ textRaw
	}
}

// HardLineBreak returns true if this node ends with a hard line break.
// See https://spec.commonmark.org/0.30/#hard-line-breaks for details.
func (n *Text) HardLineBreak() bool {
	return n.flags&textHardLineBreak != 0
}

// SetHardLineBreak sets whether this node ends with a hard line break.
func (n *Text) SetHardLineBreak(v bool) {
	if v {
		n.flags |= textHardLineBreak
	} else {
		n.flags = n.flags &^ textHardLineBreak
	}
}

// Merge merges a Node n into this node.
// Merge returns true if the given node has been merged, otherwise false.
func (n *Text) Merge(node Node, source []byte) bool {
	t, ok := node.(*Text)
	if !ok {
		return false
	}
	if n.Segment.Stop != t.Segment.Start || t.Segment.Padding != 0 || source[n.Segment.Stop-1] == '\n' || t.IsRaw() != n.IsRaw() {
		return false
	}
	n.Segment.Stop = t.Segment.Stop
	n.SetSoftLineBreak(t.SoftLineBreak())
	n.SetHardLineBreak(t.HardLineBreak())
	return true
}

// Text implements Node.Text.
func (n *Text) Text(source []byte) []byte {
	return n.Segment.Value(source)
}

// Dump implements Node.Dump.
func (n *Text) Dump(source []byte, level int) {
	fs := textFlagsString(n.flags)
	if len(fs) != 0 {
		fs = "(" + fs + ")"
	}
	fmt.Printf("%sText%s: \"%s\"\n", strings.Repeat("    ", level), fs, strings.TrimRight(string(n.Text(source)), "\n"))
}

// KindText is a NodeKind of the Text node.
var KindText = NewNodeKind("Text")

// Kind implements Node.Kind.
func (n *Text) Kind() NodeKind {
	return KindText
}

// NewText returns a new Text node.
func NewText() *Text {
	return &Text{
		BaseInline: BaseInline{},
	}
}

// NewTextSegment returns a new Text node with the given source position.
func NewTextSegment(v textm.Segment) *Text {
	return &Text{
		BaseInline: BaseInline{},
		Segment:    v,
	}
}

// NewRawTextSegment returns a new Text node with the given source position.
// The new node should be rendered as raw contents.
func NewRawTextSegment(v textm.Segment) *Text {
	t := &Text{
		BaseInline: BaseInline{},
		Segment:    v,
	}
	t.SetRaw(true)
	return t
}

// MergeOrAppendTextSegment merges a given s into the last child of the parent if
// it can be merged, otherwise creates a new Text node and appends it to after current
// last child.
func MergeOrAppendTextSegment(parent Node, s textm.Segment) {
	last := parent.LastChild()
	t, ok := last.(*Text)
	if ok && t.Segment.Stop == s.Start && !t.SoftLineBreak() {
		t.Segment = t.Segment.WithStop(s.Stop)
	} else {
		parent.AppendChild(parent, NewTextSegment(s))
	}
}

// MergeOrReplaceTextSegment merges a given s into a previous sibling of the node n
// if a previous sibling of the node n is *Text, otherwise replaces Node n with s.
func MergeOrReplaceTextSegment(parent Node, n Node, s textm.Segment) {
	prev := n.PreviousSibling()
	if t, ok := prev.(*Text); ok && t.Segment.Stop == s.Start && !t.SoftLineBreak() {
		t.Segment = t.Segment.WithStop(s.Stop)
		parent.RemoveChild(parent, n)
	} else {
		parent.ReplaceChild(parent, n, NewTextSegment(s))
	}
}

// A String struct is a textual content that has a concrete value
type String struct {
	BaseInline

	Value []byte
	flags uint8
}

// Inline implements Inline.Inline.
func (n *String) Inline() {
}

// IsRaw returns true if this text should be rendered without unescaping
// back slash escapes and resolving references.
func (n *String) IsRaw() bool {
	return n.flags&textRaw != 0
}

// SetRaw sets whether this text should be rendered as raw contents.
func (n *String) SetRaw(v bool) {
	if v {
		n.flags |= textRaw
	} else {
		n.flags = n.flags &^ textRaw
	}
}

// IsCode returns true if this text should be rendered without any
// modifications.
func (n *String) IsCode() bool {
	return n.flags&textCode != 0
}

// SetCode sets whether this text should be rendered without any modifications.
func (n *String) SetCode(v bool) {
	if v {
		n.flags |= textCode
	} else {
		n.flags = n.flags &^ textCode
	}
}

// Text implements Node.Text.
func (n *String) Text(source []byte) []byte {
	return n.Value
}

// Dump implements Node.Dump.
func (n *String) Dump(source []byte, level int) {
	fs := textFlagsString(n.flags)
	if len(fs) != 0 {
		fs = "(" + fs + ")"
	}
	fmt.Printf("%sString%s: \"%s\"\n", strings.Repeat("    ", level), fs, strings.TrimRight(string(n.Value), "\n"))
}

// KindString is a NodeKind of the String node.
var KindString = NewNodeKind("String")

// Kind implements Node.Kind.
func (n *String) Kind() NodeKind {
	return KindString
}

// NewString returns a new String node.
func NewString(v []byte) *String {
	return &String{
		Value: v,
	}
}

// A CodeSpan struct represents a code span of Markdown text.
type CodeSpan struct {
	BaseInline
}

// Inline implements Inline.Inline .
func (n *CodeSpan) Inline() {
}

// IsBlank returns true if this node consists of spaces, otherwise false.
func (n *CodeSpan) IsBlank(source []byte) bool {
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		text := c.(*Text).Segment
		if !util.IsBlank(text.Value(source)) {
			return false
		}
	}
	return true
}

// Dump implements Node.Dump
func (n *CodeSpan) Dump(source []byte, level int) {
	DumpHelper(n, source, level, nil, nil)
}

// KindCodeSpan is a NodeKind of the CodeSpan node.
var KindCodeSpan = NewNodeKind("CodeSpan")

// Kind implements Node.Kind.
func (n *CodeSpan) Kind() NodeKind {
	return KindCodeSpan
}

// NewCodeSpan returns a new CodeSpan node.
func NewCodeSpan() *CodeSpan {
	return &CodeSpan{
		BaseInline: BaseInline{},
	}
}

// An Emphasis struct represents an emphasis of Markdown text.
type Emphasis struct {
	BaseInline

	// Level is a level of the emphasis.
	Level int
}

// Dump implements Node.Dump.
func (n *Emphasis) Dump(source []byte, level int) {
	m := map[string]string{
		"Level": fmt.Sprintf("%v", n.Level),
	}
	DumpHelper(n, source, level, m, nil)
}

// KindEmphasis is a NodeKind of the Emphasis node.
var KindEmphasis = NewNodeKind("Emphasis")

// Kind implements Node.Kind.
func (n *Emphasis) Kind() NodeKind {
	return KindEmphasis
}

// NewEmphasis returns a new Emphasis node with the given level.
func NewEmphasis(level int) *Emphasis {
	return &Emphasis{
		BaseInline: BaseInline{},
		Level:      level,
	}
}

type baseLink struct {
	BaseInline

	// Destination is a destination(URL) of this link.
	Destination []byte

	// Title is a title of this link.
	Title []byte
}

// Inline implements Inline.Inline.
func (n *baseLink) Inline() {
}

// A Link struct represents a link of the Markdown text.
type Link struct {
	baseLink
}

// Dump implements Node.Dump.
func (n *Link) Dump(source []byte, level int) {
	m := map[string]string{}
	m["Destination"] = string(n.Destination)
	m["Title"] = string(n.Title)
	DumpHelper(n, source, level, m, nil)
}

// KindLink is a NodeKind of the Link node.
var KindLink = NewNodeKind("Link")

// Kind implements Node.Kind.
func (n *Link) Kind() NodeKind {
	return KindLink
}

// NewLink returns a new Link node.
func NewLink() *Link {
	c := &Link{
		baseLink: baseLink{
			BaseInline: BaseInline{},
		},
	}
	return c
}

// An Image struct represents an image of the Markdown text.
type Image struct {
	baseLink
}

// Dump implements Node.Dump.
func (n *Image) Dump(source []byte, level int) {
	m := map[string]string{}
	m["Destination"] = string(n.Destination)
	m["Title"] = string(n.Title)
	DumpHelper(n, source, level, m, nil)
}

// KindImage is a NodeKind of the Image node.
var KindImage = NewNodeKind("Image")

// Kind implements Node.Kind.
func (n *Image) Kind() NodeKind {
	return KindImage
}

// NewImage returns a new Image node.
func NewImage(link *Link) *Image {
	c := &Image{
		baseLink: baseLink{
			BaseInline: BaseInline{},
		},
	}
	c.Destination = link.Destination
	c.Title = link.Title
	for n := link.FirstChild(); n != nil; {
		next := n.NextSibling()
		link.RemoveChild(link, n)
		c.AppendChild(c, n)
		n = next
	}

	return c
}

// AutoLinkType defines kind of auto links.
type AutoLinkType int

const (
	// AutoLinkEmail indicates that an autolink is an email address.
	AutoLinkEmail AutoLinkType = iota + 1
	// AutoLinkURL indicates that an autolink is a generic URL.
	AutoLinkURL
)

// An AutoLink struct represents an autolink of the Markdown text.
type AutoLink struct {
	BaseInline
	// Type is a type of this autolink.
	AutoLinkType AutoLinkType

	// Protocol specified a protocol of the link.
	Protocol []byte

	value *Text
}

// Inline implements Inline.Inline.
func (n *AutoLink) Inline() {}

// Dump implements Node.Dump
func (n *AutoLink) Dump(source []byte, level int) {
	segment := n.value.Segment
	m := map[string]string{
		"Value": string(segment.Value(source)),
	}
	DumpHelper(n, source, level, m, nil)
}

// KindAutoLink is a NodeKind of the AutoLink node.
var KindAutoLink = NewNodeKind("AutoLink")

// Kind implements Node.Kind.
func (n *AutoLink) Kind() NodeKind {
	return KindAutoLink
}

// URL returns an url of this node.
func (n *AutoLink) URL(source []byte) []byte {
	if n.Protocol != nil {
		s := n.value.Segment
		ret := make([]byte, 0, len(n.Protocol)+s.Len()+3)
		ret = append(ret, n.Protocol...)
		ret = append(ret, ':', '/', '/')
		ret = append(ret, n.value.Text(source)...)
		return ret
	}
	return n.value.Text(source)
}

// Label returns a label of this node.
func (n *AutoLink) Label(source []byte) []byte {
	return n.value.Text(source)
}

// NewAutoLink returns a new AutoLink node.
func NewAutoLink(typ AutoLinkType, value *Text) *AutoLink {
	return &AutoLink{
		BaseInline:   BaseInline{},
		value:        value,
		AutoLinkType: typ,
	}
}

// A RawHTML struct represents an inline raw HTML of the Markdown text.
type RawHTML struct {
	BaseInline
	Segments *textm.Segments
}

// Inline implements Inline.Inline.
func (n *RawHTML) Inline() {}

// Dump implements Node.Dump.
func (n *RawHTML) Dump(source []byte, level int) {
	m := map[string]string{}
	t := []string{}
	for i := 0; i < n.Segments.Len(); i++ {
		segment := n.Segments.At(i)
		t = append(t, string(segment.Value(source)))
	}
	m["RawText"] = strings.Join(t, "")
	DumpHelper(n, source, level, m, nil)
}

// KindRawHTML is a NodeKind of the RawHTML node.
var KindRawHTML = NewNodeKind("RawHTML")

// Kind implements Node.Kind.
func (n *RawHTML) Kind() NodeKind {
	return KindRawHTML
}

// NewRawHTML returns a new RawHTML node.
func NewRawHTML() *RawHTML {
	return &RawHTML{
		Segments: textm.NewSegments(),
	}
}
