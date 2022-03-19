package lo

import (
	"fmt"
	"math/rand"
	"math"
	"golang.org/x/exp/constraints"
)

// import "golang.org/x/exp/constraints"

// IndexOf returns the index at which the first occurrence of a value is found in an array or return -1
// if the value cannot be found.
func IndexOf[T comparable](collection []T, element T) int {
	for i, item := range collection {
		if item == element {
			return i
		}
	}

	return -1
}

// IndexOf returns the index at which the last occurrence of a value is found in an array or return -1
// if the value cannot be found.
func LastIndexOf[T comparable](collection []T, element T) int {
	length := len(collection)

	for i := length - 1; i >= 0; i-- {
		if collection[i] == element {
			return i
		}
	}

	return -1
}

// Find search an element in a slice based on a predicate. It returns element and true if element was found.
func Find[T any](collection []T, predicate func(T) bool) (T, bool) {
	for _, item := range collection {
		if predicate(item) {
			return item, true
		}
	}

	var result T
	return result, false
}

// Min search the minimum value of a collection.
func Min[T constraints.Ordered](collection []T) T {
	var min T

	if len(collection) == 0 {
		return min
	}

	min = collection[0]

	for i := 1; i < len(collection); i++ {
		item := collection[i]

		// if item.Less(min) {
		if item < min {
			min = item
		}
	}

	return min
}

// Max search the maximum value of a collection.
func Max[T constraints.Ordered](collection []T) T {
	var max T

	if len(collection) == 0 {
		return max
	}

	max = collection[0]

	for i := 1; i < len(collection); i++ {
		item := collection[i]

		if item > max {
			max = item
		}
	}

	return max
}

// Last returns the last element of a collection or error if empty.
func Last[T any](collection []T) (T, error) {
	length := len(collection)

	if length == 0 {
		var t T
		return t, fmt.Errorf("last: cannot extract the last element of an empty slice")
	}

	return collection[length-1], nil
}

// Nth returns the element at index `nth` of collection. If `nth` is negative, the nth element
// from the end is returned. An error is returned when nth is out of slice bounds.
func Nth[T any](collection []T, nth int) (T, error) {
	if int(math.Abs(float64(nth))) >= len(collection) {
		var t T
		return t, fmt.Errorf("nth: %d out of slice bounds", nth)
	}

	length := len(collection)

	if nth >= 0 {
		return collection[nth], nil
	}

	return collection[length+nth], nil
}

// Sample returns a random item from collection.
func Sample[T any](collection []T) T {
	size := len(collection)
	if size == 0 {
		return Empty[T]()
	}

	return collection[rand.Intn(size)]
}

// Samples returns N random unique items from collection.
func Samples[T any](collection []T, count int) []T {
	size := len(collection)

	// put values into a map, for faster deletion
	cOpy := make([]T, 0, size)
	for _, v := range collection {
		cOpy = append(cOpy, v)
	}

	results := []T{}

	for i := 0; i < size && i < count; i++ {
		copyLength := size - i

		index := rand.Intn(size - i)
		results = append(results, cOpy[index])

		// Removes element.
		// It is faster to swap with last element and remove it.
		cOpy[index] = cOpy[copyLength-1]
		cOpy = cOpy[:copyLength-1]
	}

	return results
}
