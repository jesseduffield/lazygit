package lo

import (
	"fmt"
	"math/rand"

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

// LastIndexOf returns the index at which the last occurrence of a value is found in an array or return -1
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

// FindIndexOf searches an element in a slice based on a predicate and returns the index and true.
// It returns -1 and false if the element is not found.
func FindIndexOf[T any](collection []T, predicate func(T) bool) (T, int, bool) {
	for i, item := range collection {
		if predicate(item) {
			return item, i, true
		}
	}

	var result T
	return result, -1, false
}

// FindLastIndexOf searches last element in a slice based on a predicate and returns the index and true.
// It returns -1 and false if the element is not found.
func FindLastIndexOf[T any](collection []T, predicate func(T) bool) (T, int, bool) {
	length := len(collection)

	for i := length - 1; i >= 0; i-- {
		if predicate(collection[i]) {
			return collection[i], i, true
		}
	}

	var result T
	return result, -1, false
}

// FindOrElse search an element in a slice based on a predicate. It returns the element if found or a given fallback value otherwise.
func FindOrElse[T any](collection []T, fallback T, predicate func(T) bool) T {
	for _, item := range collection {
		if predicate(item) {
			return item
		}
	}

	return fallback
}

// FindKey returns the key of the first value matching.
func FindKey[K comparable, V comparable](object map[K]V, value V) (K, bool) {
	for k, v := range object {
		if v == value {
			return k, true
		}
	}

	return Empty[K](), false
}

// FindKeyBy returns the key of the first element predicate returns truthy for.
func FindKeyBy[K comparable, V any](object map[K]V, predicate func(K, V) bool) (K, bool) {
	for k, v := range object {
		if predicate(k, v) {
			return k, true
		}
	}

	return Empty[K](), false
}

// FindUniques returns a slice with all the unique elements of the collection.
// The order of result values is determined by the order they occur in the collection.
func FindUniques[T comparable](collection []T) []T {
	isDupl := make(map[T]bool, len(collection))

	for _, item := range collection {
		duplicated, ok := isDupl[item]
		if !ok {
			isDupl[item] = false
		} else if !duplicated {
			isDupl[item] = true
		}
	}

	result := make([]T, 0, len(collection)-len(isDupl))

	for _, item := range collection {
		if duplicated := isDupl[item]; !duplicated {
			result = append(result, item)
		}
	}

	return result
}

// FindUniquesBy returns a slice with all the unique elements of the collection.
// The order of result values is determined by the order they occur in the array. It accepts `iteratee` which is
// invoked for each element in array to generate the criterion by which uniqueness is computed.
func FindUniquesBy[T any, U comparable](collection []T, iteratee func(T) U) []T {
	isDupl := make(map[U]bool, len(collection))

	for _, item := range collection {
		key := iteratee(item)

		duplicated, ok := isDupl[key]
		if !ok {
			isDupl[key] = false
		} else if !duplicated {
			isDupl[key] = true
		}
	}

	result := make([]T, 0, len(collection)-len(isDupl))

	for _, item := range collection {
		key := iteratee(item)

		if duplicated := isDupl[key]; !duplicated {
			result = append(result, item)
		}
	}

	return result
}

// FindDuplicates returns a slice with the first occurence of each duplicated elements of the collection.
// The order of result values is determined by the order they occur in the collection.
func FindDuplicates[T comparable](collection []T) []T {
	isDupl := make(map[T]bool, len(collection))

	for _, item := range collection {
		duplicated, ok := isDupl[item]
		if !ok {
			isDupl[item] = false
		} else if !duplicated {
			isDupl[item] = true
		}
	}

	result := make([]T, 0, len(collection)-len(isDupl))

	for _, item := range collection {
		if duplicated := isDupl[item]; duplicated {
			result = append(result, item)
			isDupl[item] = false
		}
	}

	return result
}

// FindDuplicatesBy returns a slice with the first occurence of each duplicated elements of the collection.
// The order of result values is determined by the order they occur in the array. It accepts `iteratee` which is
// invoked for each element in array to generate the criterion by which uniqueness is computed.
func FindDuplicatesBy[T any, U comparable](collection []T, iteratee func(T) U) []T {
	isDupl := make(map[U]bool, len(collection))

	for _, item := range collection {
		key := iteratee(item)

		duplicated, ok := isDupl[key]
		if !ok {
			isDupl[key] = false
		} else if !duplicated {
			isDupl[key] = true
		}
	}

	result := make([]T, 0, len(collection)-len(isDupl))

	for _, item := range collection {
		key := iteratee(item)

		if duplicated := isDupl[key]; duplicated {
			result = append(result, item)
			isDupl[key] = false
		}
	}

	return result
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

		if item < min {
			min = item
		}
	}

	return min
}

// MinBy search the minimum value of a collection using the given comparison function.
// If several values of the collection are equal to the smallest value, returns the first such value.
func MinBy[T any](collection []T, comparison func(T, T) bool) T {
	var min T

	if len(collection) == 0 {
		return min
	}

	min = collection[0]

	for i := 1; i < len(collection); i++ {
		item := collection[i]

		if comparison(item, min) {
			min = item
		}
	}

	return min
}

// Max searches the maximum value of a collection.
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

// MaxBy search the maximum value of a collection using the given comparison function.
// If several values of the collection are equal to the greatest value, returns the first such value.
func MaxBy[T any](collection []T, comparison func(T, T) bool) T {
	var max T

	if len(collection) == 0 {
		return max
	}

	max = collection[0]

	for i := 1; i < len(collection); i++ {
		item := collection[i]

		if comparison(item, max) {
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
func Nth[T any, N constraints.Integer](collection []T, nth N) (T, error) {
	n := int(nth)
	l := len(collection)
	if n >= l || -n > l {
		var t T
		return t, fmt.Errorf("nth: %d out of slice bounds", n)
	}

	if n >= 0 {
		return collection[n], nil
	}
	return collection[l+n], nil
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

	cOpy := append([]T{}, collection...)

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
