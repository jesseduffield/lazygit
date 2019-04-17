// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package hashbidimap implements a bidirectional map backed by two hashmaps.
//
// A bidirectional map, or hash bag, is an associative data structure in which the (key,value) pairs form a one-to-one correspondence.
// Thus the binary relation is functional in each direction: value can also act as a key to key.
// A pair (a,b) thus provides a unique coupling between 'a' and 'b' so that 'b' can be found when 'a' is used as a key and 'a' can be found when 'b' is used as a key.
//
// Elements are unordered in the map.
//
// Structure is not thread safe.
//
// Reference: https://en.wikipedia.org/wiki/Bidirectional_map
package hashbidimap

import (
	"fmt"
	"github.com/emirpasic/gods/maps"
	"github.com/emirpasic/gods/maps/hashmap"
)

func assertMapImplementation() {
	var _ maps.BidiMap = (*Map)(nil)
}

// Map holds the elements in two hashmaps.
type Map struct {
	forwardMap hashmap.Map
	inverseMap hashmap.Map
}

// New instantiates a bidirectional map.
func New() *Map {
	return &Map{*hashmap.New(), *hashmap.New()}
}

// Put inserts element into the map.
func (m *Map) Put(key interface{}, value interface{}) {
	if valueByKey, ok := m.forwardMap.Get(key); ok {
		m.inverseMap.Remove(valueByKey)
	}
	if keyByValue, ok := m.inverseMap.Get(value); ok {
		m.forwardMap.Remove(keyByValue)
	}
	m.forwardMap.Put(key, value)
	m.inverseMap.Put(value, key)
}

// Get searches the element in the map by key and returns its value or nil if key is not found in map.
// Second return parameter is true if key was found, otherwise false.
func (m *Map) Get(key interface{}) (value interface{}, found bool) {
	return m.forwardMap.Get(key)
}

// GetKey searches the element in the map by value and returns its key or nil if value is not found in map.
// Second return parameter is true if value was found, otherwise false.
func (m *Map) GetKey(value interface{}) (key interface{}, found bool) {
	return m.inverseMap.Get(value)
}

// Remove removes the element from the map by key.
func (m *Map) Remove(key interface{}) {
	if value, found := m.forwardMap.Get(key); found {
		m.forwardMap.Remove(key)
		m.inverseMap.Remove(value)
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

// Keys returns all keys (random order).
func (m *Map) Keys() []interface{} {
	return m.forwardMap.Keys()
}

// Values returns all values (random order).
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
	str := "HashBidiMap\n"
	str += fmt.Sprintf("%v", m.forwardMap)
	return str
}
