// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package btree

import "github.com/emirpasic/gods/containers"

func assertIteratorImplementation() {
	var _ containers.ReverseIteratorWithKey = (*Iterator)(nil)
}

// Iterator holding the iterator's state
type Iterator struct {
	tree     *Tree
	node     *Node
	entry    *Entry
	position position
}

type position byte

const (
	begin, between, end position = 0, 1, 2
)

// Iterator returns a stateful iterator whose elements are key/value pairs.
func (tree *Tree) Iterator() Iterator {
	return Iterator{tree: tree, node: nil, position: begin}
}

// Next moves the iterator to the next element and returns true if there was a next element in the container.
// If Next() returns true, then next element's key and value can be retrieved by Key() and Value().
// If Next() was called for the first time, then it will point the iterator to the first element if it exists.
// Modifies the state of the iterator.
func (iterator *Iterator) Next() bool {
	// If already at end, go to end
	if iterator.position == end {
		goto end
	}
	// If at beginning, get the left-most entry in the tree
	if iterator.position == begin {
		left := iterator.tree.Left()
		if left == nil {
			goto end
		}
		iterator.node = left
		iterator.entry = left.Entries[0]
		goto between
	}
	{
		// Find current entry position in current node
		e, _ := iterator.tree.search(iterator.node, iterator.entry.Key)
		// Try to go down to the child right of the current entry
		if e+1 < len(iterator.node.Children) {
			iterator.node = iterator.node.Children[e+1]
			// Try to go down to the child left of the current node
			for len(iterator.node.Children) > 0 {
				iterator.node = iterator.node.Children[0]
			}
			// Return the left-most entry
			iterator.entry = iterator.node.Entries[0]
			goto between
		}
		// Above assures that we have reached a leaf node, so return the next entry in current node (if any)
		if e+1 < len(iterator.node.Entries) {
			iterator.entry = iterator.node.Entries[e+1]
			goto between
		}
	}
	// Reached leaf node and there are no entries to the right of the current entry, so go up to the parent
	for iterator.node.Parent != nil {
		iterator.node = iterator.node.Parent
		// Find next entry position in current node (note: search returns the first equal or bigger than entry)
		e, _ := iterator.tree.search(iterator.node, iterator.entry.Key)
		// Check that there is a next entry position in current node
		if e < len(iterator.node.Entries) {
			iterator.entry = iterator.node.Entries[e]
			goto between
		}
	}

end:
	iterator.End()
	return false

between:
	iterator.position = between
	return true
}

// Prev moves the iterator to the previous element and returns true if there was a previous element in the container.
// If Prev() returns true, then previous element's key and value can be retrieved by Key() and Value().
// Modifies the state of the iterator.
func (iterator *Iterator) Prev() bool {
	// If already at beginning, go to begin
	if iterator.position == begin {
		goto begin
	}
	// If at end, get the right-most entry in the tree
	if iterator.position == end {
		right := iterator.tree.Right()
		if right == nil {
			goto begin
		}
		iterator.node = right
		iterator.entry = right.Entries[len(right.Entries)-1]
		goto between
	}
	{
		// Find current entry position in current node
		e, _ := iterator.tree.search(iterator.node, iterator.entry.Key)
		// Try to go down to the child left of the current entry
		if e < len(iterator.node.Children) {
			iterator.node = iterator.node.Children[e]
			// Try to go down to the child right of the current node
			for len(iterator.node.Children) > 0 {
				iterator.node = iterator.node.Children[len(iterator.node.Children)-1]
			}
			// Return the right-most entry
			iterator.entry = iterator.node.Entries[len(iterator.node.Entries)-1]
			goto between
		}
		// Above assures that we have reached a leaf node, so return the previous entry in current node (if any)
		if e-1 >= 0 {
			iterator.entry = iterator.node.Entries[e-1]
			goto between
		}
	}
	// Reached leaf node and there are no entries to the left of the current entry, so go up to the parent
	for iterator.node.Parent != nil {
		iterator.node = iterator.node.Parent
		// Find previous entry position in current node (note: search returns the first equal or bigger than entry)
		e, _ := iterator.tree.search(iterator.node, iterator.entry.Key)
		// Check that there is a previous entry position in current node
		if e-1 >= 0 {
			iterator.entry = iterator.node.Entries[e-1]
			goto between
		}
	}

begin:
	iterator.Begin()
	return false

between:
	iterator.position = between
	return true
}

// Value returns the current element's value.
// Does not modify the state of the iterator.
func (iterator *Iterator) Value() interface{} {
	return iterator.entry.Value
}

// Key returns the current element's key.
// Does not modify the state of the iterator.
func (iterator *Iterator) Key() interface{} {
	return iterator.entry.Key
}

// Begin resets the iterator to its initial state (one-before-first)
// Call Next() to fetch the first element if any.
func (iterator *Iterator) Begin() {
	iterator.node = nil
	iterator.position = begin
	iterator.entry = nil
}

// End moves the iterator past the last element (one-past-the-end).
// Call Prev() to fetch the last element if any.
func (iterator *Iterator) End() {
	iterator.node = nil
	iterator.position = end
	iterator.entry = nil
}

// First moves the iterator to the first element and returns true if there was a first element in the container.
// If First() returns true, then first element's key and value can be retrieved by Key() and Value().
// Modifies the state of the iterator
func (iterator *Iterator) First() bool {
	iterator.Begin()
	return iterator.Next()
}

// Last moves the iterator to the last element and returns true if there was a last element in the container.
// If Last() returns true, then last element's key and value can be retrieved by Key() and Value().
// Modifies the state of the iterator.
func (iterator *Iterator) Last() bool {
	iterator.End()
	return iterator.Prev()
}
