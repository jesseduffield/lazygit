package ast

import (
	"fmt"

	gast "github.com/yuin/goldmark/ast"
)

// A FootnoteLink struct represents a link to a footnote of Markdown
// (PHP Markdown Extra) text.
type FootnoteLink struct {
	gast.BaseInline
	Index    int
	RefCount int
	RefIndex int
}

// Dump implements Node.Dump.
func (n *FootnoteLink) Dump(source []byte, level int) {
	m := map[string]string{}
	m["Index"] = fmt.Sprintf("%v", n.Index)
	m["RefCount"] = fmt.Sprintf("%v", n.RefCount)
	m["RefIndex"] = fmt.Sprintf("%v", n.RefIndex)
	gast.DumpHelper(n, source, level, m, nil)
}

// KindFootnoteLink is a NodeKind of the FootnoteLink node.
var KindFootnoteLink = gast.NewNodeKind("FootnoteLink")

// Kind implements Node.Kind.
func (n *FootnoteLink) Kind() gast.NodeKind {
	return KindFootnoteLink
}

// NewFootnoteLink returns a new FootnoteLink node.
func NewFootnoteLink(index int) *FootnoteLink {
	return &FootnoteLink{
		Index:    index,
		RefCount: 0,
		RefIndex: 0,
	}
}

// A FootnoteBacklink struct represents a link to a footnote of Markdown
// (PHP Markdown Extra) text.
type FootnoteBacklink struct {
	gast.BaseInline
	Index    int
	RefCount int
	RefIndex int
}

// Dump implements Node.Dump.
func (n *FootnoteBacklink) Dump(source []byte, level int) {
	m := map[string]string{}
	m["Index"] = fmt.Sprintf("%v", n.Index)
	m["RefCount"] = fmt.Sprintf("%v", n.RefCount)
	m["RefIndex"] = fmt.Sprintf("%v", n.RefIndex)
	gast.DumpHelper(n, source, level, m, nil)
}

// KindFootnoteBacklink is a NodeKind of the FootnoteBacklink node.
var KindFootnoteBacklink = gast.NewNodeKind("FootnoteBacklink")

// Kind implements Node.Kind.
func (n *FootnoteBacklink) Kind() gast.NodeKind {
	return KindFootnoteBacklink
}

// NewFootnoteBacklink returns a new FootnoteBacklink node.
func NewFootnoteBacklink(index int) *FootnoteBacklink {
	return &FootnoteBacklink{
		Index:    index,
		RefCount: 0,
		RefIndex: 0,
	}
}

// A Footnote struct represents a footnote of Markdown
// (PHP Markdown Extra) text.
type Footnote struct {
	gast.BaseBlock
	Ref   []byte
	Index int
}

// Dump implements Node.Dump.
func (n *Footnote) Dump(source []byte, level int) {
	m := map[string]string{}
	m["Index"] = fmt.Sprintf("%v", n.Index)
	m["Ref"] = fmt.Sprintf("%s", n.Ref)
	gast.DumpHelper(n, source, level, m, nil)
}

// KindFootnote is a NodeKind of the Footnote node.
var KindFootnote = gast.NewNodeKind("Footnote")

// Kind implements Node.Kind.
func (n *Footnote) Kind() gast.NodeKind {
	return KindFootnote
}

// NewFootnote returns a new Footnote node.
func NewFootnote(ref []byte) *Footnote {
	return &Footnote{
		Ref:   ref,
		Index: -1,
	}
}

// A FootnoteList struct represents footnotes of Markdown
// (PHP Markdown Extra) text.
type FootnoteList struct {
	gast.BaseBlock
	Count int
}

// Dump implements Node.Dump.
func (n *FootnoteList) Dump(source []byte, level int) {
	m := map[string]string{}
	m["Count"] = fmt.Sprintf("%v", n.Count)
	gast.DumpHelper(n, source, level, m, nil)
}

// KindFootnoteList is a NodeKind of the FootnoteList node.
var KindFootnoteList = gast.NewNodeKind("FootnoteList")

// Kind implements Node.Kind.
func (n *FootnoteList) Kind() gast.NodeKind {
	return KindFootnoteList
}

// NewFootnoteList returns a new FootnoteList node.
func NewFootnoteList() *FootnoteList {
	return &FootnoteList{
		Count: 0,
	}
}
