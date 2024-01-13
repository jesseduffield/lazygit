// Package ast defines AST nodes that represents extension's elements
package ast

import (
	gast "github.com/yuin/goldmark/ast"
)

// A Strikethrough struct represents a strikethrough of GFM text.
type Strikethrough struct {
	gast.BaseInline
}

// Dump implements Node.Dump.
func (n *Strikethrough) Dump(source []byte, level int) {
	gast.DumpHelper(n, source, level, nil, nil)
}

// KindStrikethrough is a NodeKind of the Strikethrough node.
var KindStrikethrough = gast.NewNodeKind("Strikethrough")

// Kind implements Node.Kind.
func (n *Strikethrough) Kind() gast.NodeKind {
	return KindStrikethrough
}

// NewStrikethrough returns a new Strikethrough node.
func NewStrikethrough() *Strikethrough {
	return &Strikethrough{}
}
