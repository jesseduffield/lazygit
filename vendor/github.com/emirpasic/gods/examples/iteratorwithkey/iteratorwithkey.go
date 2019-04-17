// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/emirpasic/gods/maps/treemap"
)

// IteratorWithKeyExample to demonstrate basic usage of IteratorWithKey
func main() {
	m := treemap.NewWithIntComparator()
	m.Put(1, "a")
	m.Put(2, "b")
	m.Put(3, "a")
	it := m.Iterator()

	fmt.Print("\nForward iteration\n")
	for it.Next() {
		key, value := it.Key(), it.Value()
		fmt.Print("[", key, ":", value, "]") // [0:a][1:b][2:c]
	}

	fmt.Print("\nForward iteration (again)\n")
	for it.Begin(); it.Next(); {
		key, value := it.Key(), it.Value()
		fmt.Print("[", key, ":", value, "]") // [0:a][1:b][2:c]
	}

	fmt.Print("\nBackward iteration\n")
	for it.Prev() {
		key, value := it.Key(), it.Value()
		fmt.Print("[", key, ":", value, "]") // [2:c][1:b][0:a]
	}

	fmt.Print("\nBackward iteration (again)\n")
	for it.End(); it.Prev(); {
		key, value := it.Key(), it.Value()
		fmt.Print("[", key, ":", value, "]") // [2:c][1:b][0:a]
	}

	if it.First() {
		fmt.Print("\nFirst key: ", it.Key())     // First key: 0
		fmt.Print("\nFirst value: ", it.Value()) // First value: a
	}

	if it.Last() {
		fmt.Print("\nLast key: ", it.Key())     // Last key: 3
		fmt.Print("\nLast value: ", it.Value()) // Last value: c
	}
}
