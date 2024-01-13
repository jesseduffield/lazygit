package ast

import (
	"fmt"
	gast "github.com/yuin/goldmark/ast"
	"strings"
)

// Alignment is a text alignment of table cells.
type Alignment int

const (
	// AlignLeft indicates text should be left justified.
	AlignLeft Alignment = iota + 1

	// AlignRight indicates text should be right justified.
	AlignRight

	// AlignCenter indicates text should be centered.
	AlignCenter

	// AlignNone indicates text should be aligned by default manner.
	AlignNone
)

func (a Alignment) String() string {
	switch a {
	case AlignLeft:
		return "left"
	case AlignRight:
		return "right"
	case AlignCenter:
		return "center"
	case AlignNone:
		return "none"
	}
	return ""
}

// A Table struct represents a table of Markdown(GFM) text.
type Table struct {
	gast.BaseBlock

	// Alignments returns alignments of the columns.
	Alignments []Alignment
}

// Dump implements Node.Dump
func (n *Table) Dump(source []byte, level int) {
	gast.DumpHelper(n, source, level, nil, func(level int) {
		indent := strings.Repeat("    ", level)
		fmt.Printf("%sAlignments {\n", indent)
		for i, alignment := range n.Alignments {
			indent2 := strings.Repeat("    ", level+1)
			fmt.Printf("%s%s", indent2, alignment.String())
			if i != len(n.Alignments)-1 {
				fmt.Println("")
			}
		}
		fmt.Printf("\n%s}\n", indent)
	})
}

// KindTable is a NodeKind of the Table node.
var KindTable = gast.NewNodeKind("Table")

// Kind implements Node.Kind.
func (n *Table) Kind() gast.NodeKind {
	return KindTable
}

// NewTable returns a new Table node.
func NewTable() *Table {
	return &Table{
		Alignments: []Alignment{},
	}
}

// A TableRow struct represents a table row of Markdown(GFM) text.
type TableRow struct {
	gast.BaseBlock
	Alignments []Alignment
}

// Dump implements Node.Dump.
func (n *TableRow) Dump(source []byte, level int) {
	gast.DumpHelper(n, source, level, nil, nil)
}

// KindTableRow is a NodeKind of the TableRow node.
var KindTableRow = gast.NewNodeKind("TableRow")

// Kind implements Node.Kind.
func (n *TableRow) Kind() gast.NodeKind {
	return KindTableRow
}

// NewTableRow returns a new TableRow node.
func NewTableRow(alignments []Alignment) *TableRow {
	return &TableRow{Alignments: alignments}
}

// A TableHeader struct represents a table header of Markdown(GFM) text.
type TableHeader struct {
	gast.BaseBlock
	Alignments []Alignment
}

// KindTableHeader is a NodeKind of the TableHeader node.
var KindTableHeader = gast.NewNodeKind("TableHeader")

// Kind implements Node.Kind.
func (n *TableHeader) Kind() gast.NodeKind {
	return KindTableHeader
}

// Dump implements Node.Dump.
func (n *TableHeader) Dump(source []byte, level int) {
	gast.DumpHelper(n, source, level, nil, nil)
}

// NewTableHeader returns a new TableHeader node.
func NewTableHeader(row *TableRow) *TableHeader {
	n := &TableHeader{}
	for c := row.FirstChild(); c != nil; {
		next := c.NextSibling()
		n.AppendChild(n, c)
		c = next
	}
	return n
}

// A TableCell struct represents a table cell of a Markdown(GFM) text.
type TableCell struct {
	gast.BaseBlock
	Alignment Alignment
}

// Dump implements Node.Dump.
func (n *TableCell) Dump(source []byte, level int) {
	gast.DumpHelper(n, source, level, nil, nil)
}

// KindTableCell is a NodeKind of the TableCell node.
var KindTableCell = gast.NewNodeKind("TableCell")

// Kind implements Node.Kind.
func (n *TableCell) Kind() gast.NodeKind {
	return KindTableCell
}

// NewTableCell returns a new TableCell node.
func NewTableCell() *TableCell {
	return &TableCell{
		Alignment: AlignNone,
	}
}
