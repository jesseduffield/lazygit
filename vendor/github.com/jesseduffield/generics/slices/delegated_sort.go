package slices

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

// This file delegates to the official slices package, so that we end up with a superset of the official API.

// Sort sorts a slice of any ordered type in ascending order.
func Sort[E constraints.Ordered](x []E) {
	slices.Sort(x)
}

// Sort sorts the slice x in ascending order as determined by the less function.
// This sort is not guaranteed to be stable.
func SortFunc[E any](x []E, less func(a, b E) bool) {
	slices.SortFunc(x, less)
}

// SortStable sorts the slice x while keeping the original order of equal
// elements, using less to compare elements.
func SortStableFunc[E any](x []E, less func(a, b E) bool) {
	slices.SortStableFunc(x, less)
}

// IsSorted reports whether x is sorted in ascending order.
func IsSorted[E constraints.Ordered](x []E) bool {
	return slices.IsSorted(x)
}

// IsSortedFunc reports whether x is sorted in ascending order, with less as the
// comparison function.
func IsSortedFunc[E any](x []E, less func(a, b E) bool) bool {
	return slices.IsSortedFunc(x, less)
}

// BinarySearch searches for target in a sorted slice and returns the smallest
// index at which target is found. If the target is not found, the index at
// which it could be inserted into the slice is returned; therefore, if the
// intention is to find target itself a separate check for equality with the
// element at the returned index is required.
func BinarySearch[E constraints.Ordered](x []E, target E) int {
	return slices.BinarySearch(x, target)
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
	return slices.BinarySearchFunc(x, ok)
}
