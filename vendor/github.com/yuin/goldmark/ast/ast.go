// Package ast defines AST nodes that represent markdown elements.
package ast

import (
	"bytes"
	"fmt"
	"strings"

	textm "github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// A NodeType indicates what type a node belongs to.
type NodeType int

const (
	// TypeBlock indicates that a node is kind of block nodes.
	TypeBlock NodeType = iota + 1
	// TypeInline indicates that a node is kind of inline nodes.
	TypeInline
	// TypeDocument indicates that a node is kind of document nodes.
	TypeDocument
)

// NodeKind indicates more specific type than NodeType.
type NodeKind int

func (k NodeKind) String() string {
	return kindNames[k]
}

var kindMax NodeKind
var kindNames = []string{""}

// NewNodeKind returns a new Kind value.
func NewNodeKind(name string) NodeKind {
	kindMax++
	kindNames = append(kindNames, name)
	return kindMax
}

// An Attribute is an attribute of the Node
type Attribute struct {
	Name  []byte
	Value interface{}
}

// A Node interface defines basic AST node functionalities.
type Node interface {
	// Type returns a type of this node.
	Type() NodeType

	// Kind returns a kind of this node.
	Kind() NodeKind

	// NextSibling returns a next sibling node of this node.
	NextSibling() Node

	// PreviousSibling returns a previous sibling node of this node.
	PreviousSibling() Node

	// Parent returns a parent node of this node.
	Parent() Node

	// SetParent sets a parent node to this node.
	SetParent(Node)

	// SetPreviousSibling sets a previous sibling node to this node.
	SetPreviousSibling(Node)

	// SetNextSibling sets a next sibling node to this node.
	SetNextSibling(Node)

	// HasChildren returns true if this node has any children, otherwise false.
	HasChildren() bool

	// ChildCount returns a total number of children.
	ChildCount() int

	// FirstChild returns a first child of this node.
	FirstChild() Node

	// LastChild returns a last child of this node.
	LastChild() Node

	// AppendChild append a node child to the tail of the children.
	AppendChild(self, child Node)

	// RemoveChild removes a node child from this node.
	// If a node child is not children of this node, RemoveChild nothing to do.
	RemoveChild(self, child Node)

	// RemoveChildren removes all children from this node.
	RemoveChildren(self Node)

	// SortChildren sorts childrens by comparator.
	SortChildren(comparator func(n1, n2 Node) int)

	// ReplaceChild replace a node v1 with a node insertee.
	// If v1 is not children of this node, ReplaceChild append a insetee to the
	// tail of the children.
	ReplaceChild(self, v1, insertee Node)

	// InsertBefore inserts a node insertee before a node v1.
	// If v1 is not children of this node, InsertBefore append a insetee to the
	// tail of the children.
	InsertBefore(self, v1, insertee Node)

	// InsertAfterinserts a node insertee after a node v1.
	// If v1 is not children of this node, InsertBefore append a insetee to the
	// tail of the children.
	InsertAfter(self, v1, insertee Node)

	// OwnerDocument returns this node's owner document.
	// If this node is not a child of the Document node, OwnerDocument
	// returns nil.
	OwnerDocument() *Document

	// Dump dumps an AST tree structure to stdout.
	// This function completely aimed for debugging.
	// level is a indent level. Implementer should indent informations with
	// 2 * level spaces.
	Dump(source []byte, level int)

	// Text returns text values of this node.
	Text(source []byte) []byte

	// HasBlankPreviousLines returns true if the row before this node is blank,
	// otherwise false.
	// This method is valid only for block nodes.
	HasBlankPreviousLines() bool

	// SetBlankPreviousLines sets whether the row before this node is blank.
	// This method is valid only for block nodes.
	SetBlankPreviousLines(v bool)

	// Lines returns text segments that hold positions in a source.
	// This method is valid only for block nodes.
	Lines() *textm.Segments

	// SetLines sets text segments that hold positions in a source.
	// This method is valid only for block nodes.
	SetLines(*textm.Segments)

	// IsRaw returns true if contents should be rendered as 'raw' contents.
	IsRaw() bool

	// SetAttribute sets the given value to the attributes.
	SetAttribute(name []byte, value interface{})

	// SetAttributeString sets the given value to the attributes.
	SetAttributeString(name string, value interface{})

	// Attribute returns a (attribute value, true) if an attribute
	// associated with the given name is found, otherwise
	// (nil, false)
	Attribute(name []byte) (interface{}, bool)

	// AttributeString returns a (attribute value, true) if an attribute
	// associated with the given name is found, otherwise
	// (nil, false)
	AttributeString(name string) (interface{}, bool)

	// Attributes returns a list of attributes.
	// This may be a nil if there are no attributes.
	Attributes() []Attribute

	// RemoveAttributes removes all attributes from this node.
	RemoveAttributes()
}

// A BaseNode struct implements the Node interface partialliy.
type BaseNode struct {
	firstChild Node
	lastChild  Node
	parent     Node
	next       Node
	prev       Node
	childCount int
	attributes []Attribute
}

func ensureIsolated(v Node) {
	if p := v.Parent(); p != nil {
		p.RemoveChild(p, v)
	}
}

// HasChildren implements Node.HasChildren .
func (n *BaseNode) HasChildren() bool {
	return n.firstChild != nil
}

// SetPreviousSibling implements Node.SetPreviousSibling .
func (n *BaseNode) SetPreviousSibling(v Node) {
	n.prev = v
}

// SetNextSibling implements Node.SetNextSibling .
func (n *BaseNode) SetNextSibling(v Node) {
	n.next = v
}

// PreviousSibling implements Node.PreviousSibling .
func (n *BaseNode) PreviousSibling() Node {
	return n.prev
}

// NextSibling implements Node.NextSibling .
func (n *BaseNode) NextSibling() Node {
	return n.next
}

// RemoveChild implements Node.RemoveChild .
func (n *BaseNode) RemoveChild(self, v Node) {
	if v.Parent() != self {
		return
	}
	n.childCount--
	prev := v.PreviousSibling()
	next := v.NextSibling()
	if prev != nil {
		prev.SetNextSibling(next)
	} else {
		n.firstChild = next
	}
	if next != nil {
		next.SetPreviousSibling(prev)
	} else {
		n.lastChild = prev
	}
	v.SetParent(nil)
	v.SetPreviousSibling(nil)
	v.SetNextSibling(nil)
}

// RemoveChildren implements Node.RemoveChildren .
func (n *BaseNode) RemoveChildren(self Node) {
	for c := n.firstChild; c != nil; {
		c.SetParent(nil)
		c.SetPreviousSibling(nil)
		next := c.NextSibling()
		c.SetNextSibling(nil)
		c = next
	}
	n.firstChild = nil
	n.lastChild = nil
	n.childCount = 0
}

// SortChildren implements Node.SortChildren
func (n *BaseNode) SortChildren(comparator func(n1, n2 Node) int) {
	var sorted Node
	current := n.firstChild
	for current != nil {
		next := current.NextSibling()
		if sorted == nil || comparator(sorted, current) >= 0 {
			current.SetNextSibling(sorted)
			if sorted != nil {
				sorted.SetPreviousSibling(current)
			}
			sorted = current
			sorted.SetPreviousSibling(nil)
		} else {
			c := sorted
			for c.NextSibling() != nil && comparator(c.NextSibling(), current) < 0 {
				c = c.NextSibling()
			}
			current.SetNextSibling(c.NextSibling())
			current.SetPreviousSibling(c)
			if c.NextSibling() != nil {
				c.NextSibling().SetPreviousSibling(current)
			}
			c.SetNextSibling(current)
		}
		current = next
	}
	n.firstChild = sorted
	for c := n.firstChild; c != nil; c = c.NextSibling() {
		n.lastChild = c
	}
}

// FirstChild implements Node.FirstChild .
func (n *BaseNode) FirstChild() Node {
	return n.firstChild
}

// LastChild implements Node.LastChild .
func (n *BaseNode) LastChild() Node {
	return n.lastChild
}

// ChildCount implements Node.ChildCount .
func (n *BaseNode) ChildCount() int {
	return n.childCount
}

// Parent implements Node.Parent .
func (n *BaseNode) Parent() Node {
	return n.parent
}

// SetParent implements Node.SetParent .
func (n *BaseNode) SetParent(v Node) {
	n.parent = v
}

// AppendChild implements Node.AppendChild .
func (n *BaseNode) AppendChild(self, v Node) {
	ensureIsolated(v)
	if n.firstChild == nil {
		n.firstChild = v
		v.SetNextSibling(nil)
		v.SetPreviousSibling(nil)
	} else {
		last := n.lastChild
		last.SetNextSibling(v)
		v.SetPreviousSibling(last)
	}
	v.SetParent(self)
	n.lastChild = v
	n.childCount++
}

// ReplaceChild implements Node.ReplaceChild .
func (n *BaseNode) ReplaceChild(self, v1, insertee Node) {
	n.InsertBefore(self, v1, insertee)
	n.RemoveChild(self, v1)
}

// InsertAfter implements Node.InsertAfter .
func (n *BaseNode) InsertAfter(self, v1, insertee Node) {
	n.InsertBefore(self, v1.NextSibling(), insertee)
}

// InsertBefore implements Node.InsertBefore .
func (n *BaseNode) InsertBefore(self, v1, insertee Node) {
	n.childCount++
	if v1 == nil {
		n.AppendChild(self, insertee)
		return
	}
	ensureIsolated(insertee)
	if v1.Parent() == self {
		c := v1
		prev := c.PreviousSibling()
		if prev != nil {
			prev.SetNextSibling(insertee)
			insertee.SetPreviousSibling(prev)
		} else {
			n.firstChild = insertee
			insertee.SetPreviousSibling(nil)
		}
		insertee.SetNextSibling(c)
		c.SetPreviousSibling(insertee)
		insertee.SetParent(self)
	}
}

// OwnerDocument implements Node.OwnerDocument
func (n *BaseNode) OwnerDocument() *Document {
	d := n.Parent()
	for {
		p := d.Parent()
		if p == nil {
			if v, ok := d.(*Document); ok {
				return v
			}
			break
		}
		d = p
	}
	return nil
}

// Text implements Node.Text  .
func (n *BaseNode) Text(source []byte) []byte {
	var buf bytes.Buffer
	for c := n.firstChild; c != nil; c = c.NextSibling() {
		buf.Write(c.Text(source))
	}
	return buf.Bytes()
}

// SetAttribute implements Node.SetAttribute.
func (n *BaseNode) SetAttribute(name []byte, value interface{}) {
	if n.attributes == nil {
		n.attributes = make([]Attribute, 0, 10)
	} else {
		for i, a := range n.attributes {
			if bytes.Equal(a.Name, name) {
				n.attributes[i].Name = name
				n.attributes[i].Value = value
				return
			}
		}
	}
	n.attributes = append(n.attributes, Attribute{name, value})
}

// SetAttributeString implements Node.SetAttributeString
func (n *BaseNode) SetAttributeString(name string, value interface{}) {
	n.SetAttribute(util.StringToReadOnlyBytes(name), value)
}

// Attribute implements Node.Attribute.
func (n *BaseNode) Attribute(name []byte) (interface{}, bool) {
	if n.attributes == nil {
		return nil, false
	}
	for i, a := range n.attributes {
		if bytes.Equal(a.Name, name) {
			return n.attributes[i].Value, true
		}
	}
	return nil, false
}

// AttributeString implements Node.AttributeString.
func (n *BaseNode) AttributeString(s string) (interface{}, bool) {
	return n.Attribute(util.StringToReadOnlyBytes(s))
}

// Attributes implements Node.Attributes
func (n *BaseNode) Attributes() []Attribute {
	return n.attributes
}

// RemoveAttributes implements Node.RemoveAttributes
func (n *BaseNode) RemoveAttributes() {
	n.attributes = nil
}

// DumpHelper is a helper function to implement Node.Dump.
// kv is pairs of an attribute name and an attribute value.
// cb is a function called after wrote a name and attributes.
func DumpHelper(v Node, source []byte, level int, kv map[string]string, cb func(int)) {
	name := v.Kind().String()
	indent := strings.Repeat("    ", level)
	fmt.Printf("%s%s {\n", indent, name)
	indent2 := strings.Repeat("    ", level+1)
	if v.Type() == TypeBlock {
		fmt.Printf("%sRawText: \"", indent2)
		for i := 0; i < v.Lines().Len(); i++ {
			line := v.Lines().At(i)
			fmt.Printf("%s", line.Value(source))
		}
		fmt.Printf("\"\n")
		fmt.Printf("%sHasBlankPreviousLines: %v\n", indent2, v.HasBlankPreviousLines())
	}
	for name, value := range kv {
		fmt.Printf("%s%s: %s\n", indent2, name, value)
	}
	if cb != nil {
		cb(level + 1)
	}
	for c := v.FirstChild(); c != nil; c = c.NextSibling() {
		c.Dump(source, level+1)
	}
	fmt.Printf("%s}\n", indent)
}

// WalkStatus represents a current status of the Walk function.
type WalkStatus int

const (
	// WalkStop indicates no more walking needed.
	WalkStop WalkStatus = iota + 1

	// WalkSkipChildren indicates that Walk wont walk on children of current
	// node.
	WalkSkipChildren

	// WalkContinue indicates that Walk can continue to walk.
	WalkContinue
)

// Walker is a function that will be called when Walk find a
// new node.
// entering is set true before walks children, false after walked children.
// If Walker returns error, Walk function immediately stop walking.
type Walker func(n Node, entering bool) (WalkStatus, error)

// Walk walks a AST tree by the depth first search algorithm.
func Walk(n Node, walker Walker) error {
	_, err := walkHelper(n, walker)
	return err
}

func walkHelper(n Node, walker Walker) (WalkStatus, error) {
	status, err := walker(n, true)
	if err != nil || status == WalkStop {
		return status, err
	}
	if status != WalkSkipChildren {
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			if st, err := walkHelper(c, walker); err != nil || st == WalkStop {
				return WalkStop, err
			}
		}
	}
	status, err = walker(n, false)
	if err != nil || status == WalkStop {
		return WalkStop, err
	}
	return WalkContinue, nil
}
