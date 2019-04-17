// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package btree implements a B tree.
//
// According to Knuth's definition, a B-tree of order m is a tree which satisfies the following properties:
// - Every node has at most m children.
// - Every non-leaf node (except root) has at least ⌈m/2⌉ children.
// - The root has at least two children if it is not a leaf node.
// - A non-leaf node with k children contains k−1 keys.
// - All leaves appear in the same level
//
// Structure is not thread safe.
//
// References: https://en.wikipedia.org/wiki/B-tree
package btree

import (
	"bytes"
	"fmt"
	"github.com/emirpasic/gods/trees"
	"github.com/emirpasic/gods/utils"
	"strings"
)

func assertTreeImplementation() {
	var _ trees.Tree = (*Tree)(nil)
}

// Tree holds elements of the B-tree
type Tree struct {
	Root       *Node            // Root node
	Comparator utils.Comparator // Key comparator
	size       int              // Total number of keys in the tree
	m          int              // order (maximum number of children)
}

// Node is a single element within the tree
type Node struct {
	Parent   *Node
	Entries  []*Entry // Contained keys in node
	Children []*Node  // Children nodes
}

// Entry represents the key-value pair contained within nodes
type Entry struct {
	Key   interface{}
	Value interface{}
}

// NewWith instantiates a B-tree with the order (maximum number of children) and a custom key comparator.
func NewWith(order int, comparator utils.Comparator) *Tree {
	if order < 3 {
		panic("Invalid order, should be at least 3")
	}
	return &Tree{m: order, Comparator: comparator}
}

// NewWithIntComparator instantiates a B-tree with the order (maximum number of children) and the IntComparator, i.e. keys are of type int.
func NewWithIntComparator(order int) *Tree {
	return NewWith(order, utils.IntComparator)
}

// NewWithStringComparator instantiates a B-tree with the order (maximum number of children) and the StringComparator, i.e. keys are of type string.
func NewWithStringComparator(order int) *Tree {
	return NewWith(order, utils.StringComparator)
}

// Put inserts key-value pair node into the tree.
// If key already exists, then its value is updated with the new value.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *Tree) Put(key interface{}, value interface{}) {
	entry := &Entry{Key: key, Value: value}

	if tree.Root == nil {
		tree.Root = &Node{Entries: []*Entry{entry}, Children: []*Node{}}
		tree.size++
		return
	}

	if tree.insert(tree.Root, entry) {
		tree.size++
	}
}

// Get searches the node in the tree by key and returns its value or nil if key is not found in tree.
// Second return parameter is true if key was found, otherwise false.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *Tree) Get(key interface{}) (value interface{}, found bool) {
	node, index, found := tree.searchRecursively(tree.Root, key)
	if found {
		return node.Entries[index].Value, true
	}
	return nil, false
}

// Remove remove the node from the tree by key.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *Tree) Remove(key interface{}) {
	node, index, found := tree.searchRecursively(tree.Root, key)
	if found {
		tree.delete(node, index)
		tree.size--
	}
}

// Empty returns true if tree does not contain any nodes
func (tree *Tree) Empty() bool {
	return tree.size == 0
}

// Size returns number of nodes in the tree.
func (tree *Tree) Size() int {
	return tree.size
}

// Keys returns all keys in-order
func (tree *Tree) Keys() []interface{} {
	keys := make([]interface{}, tree.size)
	it := tree.Iterator()
	for i := 0; it.Next(); i++ {
		keys[i] = it.Key()
	}
	return keys
}

// Values returns all values in-order based on the key.
func (tree *Tree) Values() []interface{} {
	values := make([]interface{}, tree.size)
	it := tree.Iterator()
	for i := 0; it.Next(); i++ {
		values[i] = it.Value()
	}
	return values
}

// Clear removes all nodes from the tree.
func (tree *Tree) Clear() {
	tree.Root = nil
	tree.size = 0
}

// Height returns the height of the tree.
func (tree *Tree) Height() int {
	return tree.Root.height()
}

// Left returns the left-most (min) node or nil if tree is empty.
func (tree *Tree) Left() *Node {
	return tree.left(tree.Root)
}

// LeftKey returns the left-most (min) key or nil if tree is empty.
func (tree *Tree) LeftKey() interface{} {
	if left := tree.Left(); left != nil {
		return left.Entries[0].Key
	}
	return nil
}

// LeftValue returns the left-most value or nil if tree is empty.
func (tree *Tree) LeftValue() interface{} {
	if left := tree.Left(); left != nil {
		return left.Entries[0].Value
	}
	return nil
}

// Right returns the right-most (max) node or nil if tree is empty.
func (tree *Tree) Right() *Node {
	return tree.right(tree.Root)
}

// RightKey returns the right-most (max) key or nil if tree is empty.
func (tree *Tree) RightKey() interface{} {
	if right := tree.Right(); right != nil {
		return right.Entries[len(right.Entries)-1].Key
	}
	return nil
}

// RightValue returns the right-most value or nil if tree is empty.
func (tree *Tree) RightValue() interface{} {
	if right := tree.Right(); right != nil {
		return right.Entries[len(right.Entries)-1].Value
	}
	return nil
}

// String returns a string representation of container (for debugging purposes)
func (tree *Tree) String() string {
	var buffer bytes.Buffer
	if _, err := buffer.WriteString("BTree\n"); err != nil {
	}
	if !tree.Empty() {
		tree.output(&buffer, tree.Root, 0, true)
	}
	return buffer.String()
}

func (entry *Entry) String() string {
	return fmt.Sprintf("%v", entry.Key)
}

func (tree *Tree) output(buffer *bytes.Buffer, node *Node, level int, isTail bool) {
	for e := 0; e < len(node.Entries)+1; e++ {
		if e < len(node.Children) {
			tree.output(buffer, node.Children[e], level+1, true)
		}
		if e < len(node.Entries) {
			if _, err := buffer.WriteString(strings.Repeat("    ", level)); err != nil {
			}
			if _, err := buffer.WriteString(fmt.Sprintf("%v", node.Entries[e].Key) + "\n"); err != nil {
			}
		}
	}
}

func (node *Node) height() int {
	height := 0
	for ; node != nil; node = node.Children[0] {
		height++
		if len(node.Children) == 0 {
			break
		}
	}
	return height
}

func (tree *Tree) isLeaf(node *Node) bool {
	return len(node.Children) == 0
}

func (tree *Tree) isFull(node *Node) bool {
	return len(node.Entries) == tree.maxEntries()
}

func (tree *Tree) shouldSplit(node *Node) bool {
	return len(node.Entries) > tree.maxEntries()
}

func (tree *Tree) maxChildren() int {
	return tree.m
}

func (tree *Tree) minChildren() int {
	return (tree.m + 1) / 2 // ceil(m/2)
}

func (tree *Tree) maxEntries() int {
	return tree.maxChildren() - 1
}

func (tree *Tree) minEntries() int {
	return tree.minChildren() - 1
}

func (tree *Tree) middle() int {
	return (tree.m - 1) / 2 // "-1" to favor right nodes to have more keys when splitting
}

// search searches only within the single node among its entries
func (tree *Tree) search(node *Node, key interface{}) (index int, found bool) {
	low, high := 0, len(node.Entries)-1
	var mid int
	for low <= high {
		mid = (high + low) / 2
		compare := tree.Comparator(key, node.Entries[mid].Key)
		switch {
		case compare > 0:
			low = mid + 1
		case compare < 0:
			high = mid - 1
		case compare == 0:
			return mid, true
		}
	}
	return low, false
}

// searchRecursively searches recursively down the tree starting at the startNode
func (tree *Tree) searchRecursively(startNode *Node, key interface{}) (node *Node, index int, found bool) {
	if tree.Empty() {
		return nil, -1, false
	}
	node = startNode
	for {
		index, found = tree.search(node, key)
		if found {
			return node, index, true
		}
		if tree.isLeaf(node) {
			return nil, -1, false
		}
		node = node.Children[index]
	}
}

func (tree *Tree) insert(node *Node, entry *Entry) (inserted bool) {
	if tree.isLeaf(node) {
		return tree.insertIntoLeaf(node, entry)
	}
	return tree.insertIntoInternal(node, entry)
}

func (tree *Tree) insertIntoLeaf(node *Node, entry *Entry) (inserted bool) {
	insertPosition, found := tree.search(node, entry.Key)
	if found {
		node.Entries[insertPosition] = entry
		return false
	}
	// Insert entry's key in the middle of the node
	node.Entries = append(node.Entries, nil)
	copy(node.Entries[insertPosition+1:], node.Entries[insertPosition:])
	node.Entries[insertPosition] = entry
	tree.split(node)
	return true
}

func (tree *Tree) insertIntoInternal(node *Node, entry *Entry) (inserted bool) {
	insertPosition, found := tree.search(node, entry.Key)
	if found {
		node.Entries[insertPosition] = entry
		return false
	}
	return tree.insert(node.Children[insertPosition], entry)
}

func (tree *Tree) split(node *Node) {
	if !tree.shouldSplit(node) {
		return
	}

	if node == tree.Root {
		tree.splitRoot()
		return
	}

	tree.splitNonRoot(node)
}

func (tree *Tree) splitNonRoot(node *Node) {
	middle := tree.middle()
	parent := node.Parent

	left := &Node{Entries: append([]*Entry(nil), node.Entries[:middle]...), Parent: parent}
	right := &Node{Entries: append([]*Entry(nil), node.Entries[middle+1:]...), Parent: parent}

	// Move children from the node to be split into left and right nodes
	if !tree.isLeaf(node) {
		left.Children = append([]*Node(nil), node.Children[:middle+1]...)
		right.Children = append([]*Node(nil), node.Children[middle+1:]...)
		setParent(left.Children, left)
		setParent(right.Children, right)
	}

	insertPosition, _ := tree.search(parent, node.Entries[middle].Key)

	// Insert middle key into parent
	parent.Entries = append(parent.Entries, nil)
	copy(parent.Entries[insertPosition+1:], parent.Entries[insertPosition:])
	parent.Entries[insertPosition] = node.Entries[middle]

	// Set child left of inserted key in parent to the created left node
	parent.Children[insertPosition] = left

	// Set child right of inserted key in parent to the created right node
	parent.Children = append(parent.Children, nil)
	copy(parent.Children[insertPosition+2:], parent.Children[insertPosition+1:])
	parent.Children[insertPosition+1] = right

	tree.split(parent)
}

func (tree *Tree) splitRoot() {
	middle := tree.middle()

	left := &Node{Entries: append([]*Entry(nil), tree.Root.Entries[:middle]...)}
	right := &Node{Entries: append([]*Entry(nil), tree.Root.Entries[middle+1:]...)}

	// Move children from the node to be split into left and right nodes
	if !tree.isLeaf(tree.Root) {
		left.Children = append([]*Node(nil), tree.Root.Children[:middle+1]...)
		right.Children = append([]*Node(nil), tree.Root.Children[middle+1:]...)
		setParent(left.Children, left)
		setParent(right.Children, right)
	}

	// Root is a node with one entry and two children (left and right)
	newRoot := &Node{
		Entries:  []*Entry{tree.Root.Entries[middle]},
		Children: []*Node{left, right},
	}

	left.Parent = newRoot
	right.Parent = newRoot
	tree.Root = newRoot
}

func setParent(nodes []*Node, parent *Node) {
	for _, node := range nodes {
		node.Parent = parent
	}
}

func (tree *Tree) left(node *Node) *Node {
	if tree.Empty() {
		return nil
	}
	current := node
	for {
		if tree.isLeaf(current) {
			return current
		}
		current = current.Children[0]
	}
}

func (tree *Tree) right(node *Node) *Node {
	if tree.Empty() {
		return nil
	}
	current := node
	for {
		if tree.isLeaf(current) {
			return current
		}
		current = current.Children[len(current.Children)-1]
	}
}

// leftSibling returns the node's left sibling and child index (in parent) if it exists, otherwise (nil,-1)
// key is any of keys in node (could even be deleted).
func (tree *Tree) leftSibling(node *Node, key interface{}) (*Node, int) {
	if node.Parent != nil {
		index, _ := tree.search(node.Parent, key)
		index--
		if index >= 0 && index < len(node.Parent.Children) {
			return node.Parent.Children[index], index
		}
	}
	return nil, -1
}

// rightSibling returns the node's right sibling and child index (in parent) if it exists, otherwise (nil,-1)
// key is any of keys in node (could even be deleted).
func (tree *Tree) rightSibling(node *Node, key interface{}) (*Node, int) {
	if node.Parent != nil {
		index, _ := tree.search(node.Parent, key)
		index++
		if index < len(node.Parent.Children) {
			return node.Parent.Children[index], index
		}
	}
	return nil, -1
}

// delete deletes an entry in node at entries' index
// ref.: https://en.wikipedia.org/wiki/B-tree#Deletion
func (tree *Tree) delete(node *Node, index int) {
	// deleting from a leaf node
	if tree.isLeaf(node) {
		deletedKey := node.Entries[index].Key
		tree.deleteEntry(node, index)
		tree.rebalance(node, deletedKey)
		if len(tree.Root.Entries) == 0 {
			tree.Root = nil
		}
		return
	}

	// deleting from an internal node
	leftLargestNode := tree.right(node.Children[index]) // largest node in the left sub-tree (assumed to exist)
	leftLargestEntryIndex := len(leftLargestNode.Entries) - 1
	node.Entries[index] = leftLargestNode.Entries[leftLargestEntryIndex]
	deletedKey := leftLargestNode.Entries[leftLargestEntryIndex].Key
	tree.deleteEntry(leftLargestNode, leftLargestEntryIndex)
	tree.rebalance(leftLargestNode, deletedKey)
}

// rebalance rebalances the tree after deletion if necessary and returns true, otherwise false.
// Note that we first delete the entry and then call rebalance, thus the passed deleted key as reference.
func (tree *Tree) rebalance(node *Node, deletedKey interface{}) {
	// check if rebalancing is needed
	if node == nil || len(node.Entries) >= tree.minEntries() {
		return
	}

	// try to borrow from left sibling
	leftSibling, leftSiblingIndex := tree.leftSibling(node, deletedKey)
	if leftSibling != nil && len(leftSibling.Entries) > tree.minEntries() {
		// rotate right
		node.Entries = append([]*Entry{node.Parent.Entries[leftSiblingIndex]}, node.Entries...) // prepend parent's separator entry to node's entries
		node.Parent.Entries[leftSiblingIndex] = leftSibling.Entries[len(leftSibling.Entries)-1]
		tree.deleteEntry(leftSibling, len(leftSibling.Entries)-1)
		if !tree.isLeaf(leftSibling) {
			leftSiblingRightMostChild := leftSibling.Children[len(leftSibling.Children)-1]
			leftSiblingRightMostChild.Parent = node
			node.Children = append([]*Node{leftSiblingRightMostChild}, node.Children...)
			tree.deleteChild(leftSibling, len(leftSibling.Children)-1)
		}
		return
	}

	// try to borrow from right sibling
	rightSibling, rightSiblingIndex := tree.rightSibling(node, deletedKey)
	if rightSibling != nil && len(rightSibling.Entries) > tree.minEntries() {
		// rotate left
		node.Entries = append(node.Entries, node.Parent.Entries[rightSiblingIndex-1]) // append parent's separator entry to node's entries
		node.Parent.Entries[rightSiblingIndex-1] = rightSibling.Entries[0]
		tree.deleteEntry(rightSibling, 0)
		if !tree.isLeaf(rightSibling) {
			rightSiblingLeftMostChild := rightSibling.Children[0]
			rightSiblingLeftMostChild.Parent = node
			node.Children = append(node.Children, rightSiblingLeftMostChild)
			tree.deleteChild(rightSibling, 0)
		}
		return
	}

	// merge with siblings
	if rightSibling != nil {
		// merge with right sibling
		node.Entries = append(node.Entries, node.Parent.Entries[rightSiblingIndex-1])
		node.Entries = append(node.Entries, rightSibling.Entries...)
		deletedKey = node.Parent.Entries[rightSiblingIndex-1].Key
		tree.deleteEntry(node.Parent, rightSiblingIndex-1)
		tree.appendChildren(node.Parent.Children[rightSiblingIndex], node)
		tree.deleteChild(node.Parent, rightSiblingIndex)
	} else if leftSibling != nil {
		// merge with left sibling
		entries := append([]*Entry(nil), leftSibling.Entries...)
		entries = append(entries, node.Parent.Entries[leftSiblingIndex])
		node.Entries = append(entries, node.Entries...)
		deletedKey = node.Parent.Entries[leftSiblingIndex].Key
		tree.deleteEntry(node.Parent, leftSiblingIndex)
		tree.prependChildren(node.Parent.Children[leftSiblingIndex], node)
		tree.deleteChild(node.Parent, leftSiblingIndex)
	}

	// make the merged node the root if its parent was the root and the root is empty
	if node.Parent == tree.Root && len(tree.Root.Entries) == 0 {
		tree.Root = node
		node.Parent = nil
		return
	}

	// parent might underflow, so try to rebalance if necessary
	tree.rebalance(node.Parent, deletedKey)
}

func (tree *Tree) prependChildren(fromNode *Node, toNode *Node) {
	children := append([]*Node(nil), fromNode.Children...)
	toNode.Children = append(children, toNode.Children...)
	setParent(fromNode.Children, toNode)
}

func (tree *Tree) appendChildren(fromNode *Node, toNode *Node) {
	toNode.Children = append(toNode.Children, fromNode.Children...)
	setParent(fromNode.Children, toNode)
}

func (tree *Tree) deleteEntry(node *Node, index int) {
	copy(node.Entries[index:], node.Entries[index+1:])
	node.Entries[len(node.Entries)-1] = nil
	node.Entries = node.Entries[:len(node.Entries)-1]
}

func (tree *Tree) deleteChild(node *Node, index int) {
	if index >= len(node.Children) {
		return
	}
	copy(node.Children[index:], node.Children[index+1:])
	node.Children[len(node.Children)-1] = nil
	node.Children = node.Children[:len(node.Children)-1]
}
