// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices

import "golang.org/x/exp/constraints"

// Sort sorts a slice of any ordered type in ascending order.
func Sort[E constraints.Ordered](x []E) {
	n := len(x)
	quickSortOrdered(x, 0, n, maxDepth(n))
}

// Sort sorts the slice x in ascending order as determined by the less function.
// This sort is not guaranteed to be stable.
func SortFunc[E any](x []E, less func(a, b E) bool) {
	n := len(x)
	quickSortLessFunc(x, 0, n, maxDepth(n), less)
}

// SortStable sorts the slice x while keeping the original order of equal
// elements, using less to compare elements.
func SortStableFunc[E any](x []E, less func(a, b E) bool) {
	stableLessFunc(x, len(x), less)
}

// IsSorted reports whether x is sorted in ascending order.
func IsSorted[E constraints.Ordered](x []E) bool {
	for i := len(x) - 1; i > 0; i-- {
		if x[i] < x[i-1] {
			return false
		}
	}
	return true
}

// IsSortedFunc reports whether x is sorted in ascending order, with less as the
// comparison function.
func IsSortedFunc[E any](x []E, less func(a, b E) bool) bool {
	for i := len(x) - 1; i > 0; i-- {
		if less(x[i], x[i-1]) {
			return false
		}
	}
	return true
}

// BinarySearch searches for target in a sorted slice and returns the smallest
// index at which target is found. If the target is not found, the index at
// which it could be inserted into the slice is returned; therefore, if the
// intention is to find target itself a separate check for equality with the
// element at the returned index is required.
func BinarySearch[E constraints.Ordered](x []E, target E) int {
	return search(len(x), func(i int) bool { return x[i] >= target })
}

// BinarySearchFunc uses binary search to find and return the smallest index i
// in [0, n) at which ok(i) is true, assuming that on the range [0, n),
// ok(i) == true implies ok(i+1) == true. That is, BinarySearchFunc requires
// that ok is false for some (possibly empty) prefix of the input range [0, n)
// and then true for the (possibly empty) remainder; BinarySearchFunc returns
// the first true index. If there is no such index, BinarySearchFunc returns n.
// (Note that the "not found" return value is not -1 as in, for instance,
// strings.Index.) Search calls ok(i) only for i in the range [0, n).
func BinarySearchFunc[E any](x []E, ok func(E) bool) int {
	return search(len(x), func(i int) bool { return ok(x[i]) })
}

// maxDepth returns a threshold at which quicksort should switch
// to heapsort. It returns 2*ceil(lg(n+1)).
func maxDepth(n int) int {
	var depth int
	for i := n; i > 0; i >>= 1 {
		depth++
	}
	return depth * 2
}

func search(n int, f func(int) bool) int {
	// Define f(-1) == false and f(n) == true.
	// Invariant: f(i-1) == false, f(j) == true.
	i, j := 0, n
	for i < j {
		h := int(uint(i+j) >> 1) // avoid overflow when computing h
		// i ≤ h < j
		if !f(h) {
			i = h + 1 // preserves f(i-1) == false
		} else {
			j = h // preserves f(j) == true
		}
	}
	// i == j, f(i-1) == false, and f(j) (= f(i)) == true  =>  answer is i.
	return i
}
