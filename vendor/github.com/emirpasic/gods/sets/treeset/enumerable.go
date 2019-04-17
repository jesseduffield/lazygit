// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package treeset

import (
	"github.com/emirpasic/gods/containers"
	rbt "github.com/emirpasic/gods/trees/redblacktree"
)

func assertEnumerableImplementation() {
	var _ containers.EnumerableWithIndex = (*Set)(nil)
}

// Each calls the given function once for each element, passing that element's index and value.
func (set *Set) Each(f func(index int, value interface{})) {
	iterator := set.Iterator()
	for iterator.Next() {
		f(iterator.Index(), iterator.Value())
	}
}

// Map invokes the given function once for each element and returns a
// container containing the values returned by the given function.
func (set *Set) Map(f func(index int, value interface{}) interface{}) *Set {
	newSet := &Set{tree: rbt.NewWith(set.tree.Comparator)}
	iterator := set.Iterator()
	for iterator.Next() {
		newSet.Add(f(iterator.Index(), iterator.Value()))
	}
	return newSet
}

// Select returns a new container containing all elements for which the given function returns a true value.
func (set *Set) Select(f func(index int, value interface{}) bool) *Set {
	newSet := &Set{tree: rbt.NewWith(set.tree.Comparator)}
	iterator := set.Iterator()
	for iterator.Next() {
		if f(iterator.Index(), iterator.Value()) {
			newSet.Add(iterator.Value())
		}
	}
	return newSet
}

// Any passes each element of the container to the given function and
// returns true if the function ever returns true for any element.
func (set *Set) Any(f func(index int, value interface{}) bool) bool {
	iterator := set.Iterator()
	for iterator.Next() {
		if f(iterator.Index(), iterator.Value()) {
			return true
		}
	}
	return false
}

// All passes each element of the container to the given function and
// returns true if the function returns true for all elements.
func (set *Set) All(f func(index int, value interface{}) bool) bool {
	iterator := set.Iterator()
	for iterator.Next() {
		if !f(iterator.Index(), iterator.Value()) {
			return false
		}
	}
	return true
}

// Find passes each element of the container to the given function and returns
// the first (index,value) for which the function is true or -1,nil otherwise
// if no element matches the criteria.
func (set *Set) Find(f func(index int, value interface{}) bool) (int, interface{}) {
	iterator := set.Iterator()
	for iterator.Next() {
		if f(iterator.Index(), iterator.Value()) {
			return iterator.Index(), iterator.Value()
		}
	}
	return -1, nil
}
