package ast

import (
	gast "github.com/yuin/goldmark/ast"
)

// A DefinitionList struct represents a definition list of Markdown
// (PHPMarkdownExtra) text.
type DefinitionList struct {
	gast.BaseBlock
	Offset             int
	TemporaryParagraph *gast.Paragraph
}

// Dump implements Node.Dump.
func (n *DefinitionList) Dump(source []byte, level int) {
	gast.DumpHelper(n, source, level, nil, nil)
}

// KindDefinitionList is a NodeKind of the DefinitionList node.
var KindDefinitionList = gast.NewNodeKind("DefinitionList")

// Kind implements Node.Kind.
func (n *DefinitionList) Kind() gast.NodeKind {
	return KindDefinitionList
}

// NewDefinitionList returns a new DefinitionList node.
func NewDefinitionList(offset int, para *gast.Paragraph) *DefinitionList {
	return &DefinitionList{
		Offset:             offset,
		TemporaryParagraph: para,
	}
}

// A DefinitionTerm struct represents a definition list term of Markdown
// (PHPMarkdownExtra) text.
type DefinitionTerm struct {
	gast.BaseBlock
}

// Dump implements Node.Dump.
func (n *DefinitionTerm) Dump(source []byte, level int) {
	gast.DumpHelper(n, source, level, nil, nil)
}

// KindDefinitionTerm is a NodeKind of the DefinitionTerm node.
var KindDefinitionTerm = gast.NewNodeKind("DefinitionTerm")

// Kind implements Node.Kind.
func (n *DefinitionTerm) Kind() gast.NodeKind {
	return KindDefinitionTerm
}

// NewDefinitionTerm returns a new DefinitionTerm node.
func NewDefinitionTerm() *DefinitionTerm {
	return &DefinitionTerm{}
}

// A DefinitionDescription struct represents a definition list description of Markdown
// (PHPMarkdownExtra) text.
type DefinitionDescription struct {
	gast.BaseBlock
	IsTight bool
}

// Dump implements Node.Dump.
func (n *DefinitionDescription) Dump(source []byte, level int) {
	gast.DumpHelper(n, source, level, nil, nil)
}

// KindDefinitionDescription is a NodeKind of the DefinitionDescription node.
var KindDefinitionDescription = gast.NewNodeKind("DefinitionDescription")

// Kind implements Node.Kind.
func (n *DefinitionDescription) Kind() gast.NodeKind {
	return KindDefinitionDescription
}

// NewDefinitionDescription returns a new DefinitionDescription node.
func NewDefinitionDescription() *DefinitionDescription {
	return &DefinitionDescription{}
}
