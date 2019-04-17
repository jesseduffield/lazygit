// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package redblacktree

import "github.com/emirpasic/gods/containers"

func assertIteratorImplementation() {
	var _ containers.ReverseIteratorWithKey = (*Iterator)(nil)
}

// Iterator holding the iterator's state
type Iterator struct {
	tree     *Tree
	node     *Node
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
	if iterator.position == end {
		goto end
	}
	if iterator.position == begin {
		left := iterator.tree.Left()
		if left == nil {
			goto end
		}
		iterator.node = left
		goto between
	}
	if iterator.node.Right != nil {
		iterator.node = iterator.node.Right
		for iterator.node.Left != nil {
			iterator.node = iterator.node.Left
		}
		goto between
	}
	if iterator.node.Parent != nil {
		node := iterator.node
		for iterator.node.Parent != nil {
			iterator.node = iterator.node.Parent
			if iterator.tree.Comparator(node.Key, iterator.node.Key) <= 0 {
				goto between
			}
		}
	}

end:
	iterator.node = nil
	iterator.position = end
	return false

between:
	iterator.position = between
	return true
}

// Prev moves the iterator to the previous element and returns true if there was a previous element in the container.
// If Prev() returns true, then previous element's key and value can be retrieved by Key() and Value().
// Modifies the state of the iterator.
func (iterator *Iterator) Prev() bool {
	if iterator.position == begin {
		goto begin
	}
	if iterator.position == end {
		right := iterator.tree.Right()
		if right == nil {
			goto begin
		}
		iterator.node = right
		goto between
	}
	if iterator.node.Left != nil {
		iterator.node = iterator.node.Left
		for iterator.node.Right != nil {
			iterator.node = iterator.node.Right
		}
		goto between
	}
	if iterator.node.Parent != nil {
		node := iterator.node
		for iterator.node.Parent != nil {
			iterator.node = iterator.node.Parent
			if iterator.tree.Comparator(node.Key, iterator.node.Key) >= 0 {
				goto between
			}
		}
	}

begin:
	iterator.node = nil
	iterator.position = begin
	return false

between:
	iterator.position = between
	return true
}

// Value returns the current element's value.
// Does not modify the state of the iterator.
func (iterator *Iterator) Value() interface{} {
	return iterator.node.Value
}

// Key returns the current element's key.
// Does not modify the state of the iterator.
func (iterator *Iterator) Key() interface{} {
	return iterator.node.Key
}

// Begin resets the iterator to its initial state (one-before-first)
// Call Next() to fetch the first element if any.
func (iterator *Iterator) Begin() {
	iterator.node = nil
	iterator.position = begin
}

// End moves the iterator past the last element (one-past-the-end).
// Call Prev() to fetch the last element if any.
func (iterator *Iterator) End() {
	iterator.node = nil
	iterator.position = end
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
