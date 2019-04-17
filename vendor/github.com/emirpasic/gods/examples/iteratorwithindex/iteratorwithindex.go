// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/emirpasic/gods/sets/treeset"
)

// IteratorWithIndexExample to demonstrate basic usage of IteratorWithIndex
func main() {
	set := treeset.NewWithStringComparator()
	set.Add("a", "b", "c")
	it := set.Iterator()

	fmt.Print("\nForward iteration\n")
	for it.Next() {
		index, value := it.Index(), it.Value()
		fmt.Print("[", index, ":", value, "]") // [0:a][1:b][2:c]
	}

	fmt.Print("\nForward iteration (again)\n")
	for it.Begin(); it.Next(); {
		index, value := it.Index(), it.Value()
		fmt.Print("[", index, ":", value, "]") // [0:a][1:b][2:c]
	}

	fmt.Print("\nBackward iteration\n")
	for it.Prev() {
		index, value := it.Index(), it.Value()
		fmt.Print("[", index, ":", value, "]") // [2:c][1:b][0:a]
	}

	fmt.Print("\nBackward iteration (again)\n")
	for it.End(); it.Prev(); {
		index, value := it.Index(), it.Value()
		fmt.Print("[", index, ":", value, "]") // [2:c][1:b][0:a]
	}

	if it.First() {
		fmt.Print("\nFirst index: ", it.Index()) // First index: 0
		fmt.Print("\nFirst value: ", it.Value()) // First value: a
	}

	if it.Last() {
		fmt.Print("\nLast index: ", it.Index()) // Last index: 3
		fmt.Print("\nLast value: ", it.Value()) // Last value: c
	}
}
