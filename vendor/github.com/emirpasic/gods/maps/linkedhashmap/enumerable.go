// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package linkedhashmap

import "github.com/emirpasic/gods/containers"

func assertEnumerableImplementation() {
	var _ containers.EnumerableWithKey = (*Map)(nil)
}

// Each calls the given function once for each element, passing that element's key and value.
func (m *Map) Each(f func(key interface{}, value interface{})) {
	iterator := m.Iterator()
	for iterator.Next() {
		f(iterator.Key(), iterator.Value())
	}
}

// Map invokes the given function once for each element and returns a container
// containing the values returned by the given function as key/value pairs.
func (m *Map) Map(f func(key1 interface{}, value1 interface{}) (interface{}, interface{})) *Map {
	newMap := New()
	iterator := m.Iterator()
	for iterator.Next() {
		key2, value2 := f(iterator.Key(), iterator.Value())
		newMap.Put(key2, value2)
	}
	return newMap
}

// Select returns a new container containing all elements for which the given function returns a true value.
func (m *Map) Select(f func(key interface{}, value interface{}) bool) *Map {
	newMap := New()
	iterator := m.Iterator()
	for iterator.Next() {
		if f(iterator.Key(), iterator.Value()) {
			newMap.Put(iterator.Key(), iterator.Value())
		}
	}
	return newMap
}

// Any passes each element of the container to the given function and
// returns true if the function ever returns true for any element.
func (m *Map) Any(f func(key interface{}, value interface{}) bool) bool {
	iterator := m.Iterator()
	for iterator.Next() {
		if f(iterator.Key(), iterator.Value()) {
			return true
		}
	}
	return false
}

// All passes each element of the container to the given function and
// returns true if the function returns true for all elements.
func (m *Map) All(f func(key interface{}, value interface{}) bool) bool {
	iterator := m.Iterator()
	for iterator.Next() {
		if !f(iterator.Key(), iterator.Value()) {
			return false
		}
	}
	return true
}

// Find passes each element of the container to the given function and returns
// the first (key,value) for which the function is true or nil,nil otherwise if no element
// matches the criteria.
func (m *Map) Find(f func(key interface{}, value interface{}) bool) (interface{}, interface{}) {
	iterator := m.Iterator()
	for iterator.Next() {
		if f(iterator.Key(), iterator.Value()) {
			return iterator.Key(), iterator.Value()
		}
	}
	return nil, nil
}
