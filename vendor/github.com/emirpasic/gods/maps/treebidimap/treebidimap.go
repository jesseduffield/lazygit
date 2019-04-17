// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package treebidimap implements a bidirectional map backed by two red-black tree.
//
// This structure guarantees that the map will be in both ascending key and value order.
//
// Other than key and value ordering, the goal with this structure is to avoid duplication of elements, which can be significant if contained elements are large.
//
// A bidirectional map, or hash bag, is an associative data structure in which the (key,value) pairs form a one-to-one correspondence.
// Thus the binary relation is functional in each direction: value can also act as a key to key.
// A pair (a,b) thus provides a unique coupling between 'a' and 'b' so that 'b' can be found when 'a' is used as a key and 'a' can be found when 'b' is used as a key.
//
// Structure is not thread safe.
//
// Reference: https://en.wikipedia.org/wiki/Bidirectional_map
package treebidimap

import (
	"fmt"
	"github.com/emirpasic/gods/maps"
	"github.com/emirpasic/gods/trees/redblacktree"
	"github.com/emirpasic/gods/utils"
	"strings"
)

func assertMapImplementation() {
	var _ maps.BidiMap = (*Map)(nil)
}

// Map holds the elements in two red-black trees.
type Map struct {
	forwardMap      redblacktree.Tree
	inverseMap      redblacktree.Tree
	keyComparator   utils.Comparator
	valueComparator utils.Comparator
}

type data struct {
	key   interface{}
	value interface{}
}

// NewWith instantiates a bidirectional map.
func NewWith(keyComparator utils.Comparator, valueComparator utils.Comparator) *Map {
	return &Map{
		forwardMap:      *redblacktree.NewWith(keyComparator),
		inverseMap:      *redblacktree.NewWith(valueComparator),
		keyComparator:   keyComparator,
		valueComparator: valueComparator,
	}
}

// NewWithIntComparators instantiates a bidirectional map with the IntComparator for key and value, i.e. keys and values are of type int.
func NewWithIntComparators() *Map {
	return NewWith(utils.IntComparator, utils.IntComparator)
}

// NewWithStringComparators instantiates a bidirectional map with the StringComparator for key and value, i.e. keys and values are of type string.
func NewWithStringComparators() *Map {
	return NewWith(utils.StringComparator, utils.StringComparator)
}

// Put inserts element into the map.
func (m *Map) Put(key interface{}, value interface{}) {
	if d, ok := m.forwardMap.Get(key); ok {
		m.inverseMap.Remove(d.(*data).value)
	}
	if d, ok := m.inverseMap.Get(value); ok {
		m.forwardMap.Remove(d.(*data).key)
	}
	d := &data{key: key, value: value}
	m.forwardMap.Put(key, d)
	m.inverseMap.Put(value, d)
}

// Get searches the element in the map by key and returns its value or nil if key is not found in map.
// Second return parameter is true if key was found, otherwise false.
func (m *Map) Get(key interface{}) (value interface{}, found bool) {
	if d, ok := m.forwardMap.Get(key); ok {
		return d.(*data).value, true
	}
	return nil, false
}

// GetKey searches the element in the map by value and returns its key or nil if value is not found in map.
// Second return parameter is true if value was found, otherwise false.
func (m *Map) GetKey(value interface{}) (key interface{}, found bool) {
	if d, ok := m.inverseMap.Get(value); ok {
		return d.(*data).key, true
	}
	return nil, false
}

// Remove removes the element from the map by key.
func (m *Map) Remove(key interface{}) {
	if d, found := m.forwardMap.Get(key); found {
		m.forwardMap.Remove(key)
		m.inverseMap.Remove(d.(*data).value)
	}
}

// Empty returns true if map does not contain any elements
func (m *Map) Empty() bool {
	return m.Size() == 0
}

// Size returns number of elements in the map.
func (m *Map) Size() int {
	return m.forwardMap.Size()
}

// Keys returns all keys (ordered).
func (m *Map) Keys() []interface{} {
	return m.forwardMap.Keys()
}

// Values returns all values (ordered).
func (m *Map) Values() []interface{} {
	return m.inverseMap.Keys()
}

// Clear removes all elements from the map.
func (m *Map) Clear() {
	m.forwardMap.Clear()
	m.inverseMap.Clear()
}

// String returns a string representation of container
func (m *Map) String() string {
	str := "TreeBidiMap\nmap["
	it := m.Iterator()
	for it.Next() {
		str += fmt.Sprintf("%v:%v ", it.Key(), it.Value())
	}
	return strings.TrimRight(str, " ") + "]"
}
