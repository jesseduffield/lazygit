package lo

import (
	"sort"

	"github.com/samber/lo/internal/constraints"
	"github.com/samber/lo/mutable"
)

// Filter iterates over elements of collection, returning a slice of all elements predicate returns true for.
// Play: https://go.dev/play/p/Apjg3WeSi7K
func Filter[T any, Slice ~[]T](collection Slice, predicate func(item T, index int) bool) Slice {
	result := make(Slice, 0, len(collection))

	for i := range collection {
		if predicate(collection[i], i) {
			result = append(result, collection[i])
		}
	}

	return result
}

// FilterErr iterates over elements of collection, returning a slice of all elements predicate returns true for.
// If the predicate returns an error, iteration stops immediately and returns the error.
// Play: https://go.dev/play/p/Apjg3WeSi7K
func FilterErr[T any, Slice ~[]T](collection Slice, predicate func(item T, index int) (bool, error)) (Slice, error) {
	result := make(Slice, 0, len(collection))

	for i := range collection {
		ok, err := predicate(collection[i], i)
		if err != nil {
			return nil, err
		}
		if ok {
			result = append(result, collection[i])
		}
	}

	return result, nil
}

// Map manipulates a slice and transforms it to a slice of another type.
// Play: https://go.dev/play/p/OkPcYAhBo0D
func Map[T, R any](collection []T, transform func(item T, index int) R) []R {
	result := make([]R, len(collection))

	for i := range collection {
		result[i] = transform(collection[i], i)
	}

	return result
}

// MapErr manipulates a slice and transforms it to a slice of another type.
// It returns the first error returned by the transform function.
func MapErr[T, R any](collection []T, transform func(item T, index int) (R, error)) ([]R, error) {
	result := make([]R, len(collection))

	for i := range collection {
		r, err := transform(collection[i], i)
		if err != nil {
			return nil, err
		}
		result[i] = r
	}

	return result, nil
}

// UniqMap manipulates a slice and transforms it to a slice of another type with unique values.
// Play: https://go.dev/play/p/fygzLBhvUdB
func UniqMap[T any, R comparable](collection []T, transform func(item T, index int) R) []R {
	seen := make(map[R]struct{}, len(collection))

	for i := range collection {
		r := transform(collection[i], i)
		if _, ok := seen[r]; !ok {
			seen[r] = struct{}{}
		}
	}

	return Keys(seen)
}

// FilterMap returns a slice obtained after both filtering and mapping using the given callback function.
// The callback function should return two values:
//   - the result of the mapping operation and
//   - whether the result element should be included or not.
//
// Play: https://go.dev/play/p/CgHYNUpOd1I
func FilterMap[T, R any](collection []T, callback func(item T, index int) (R, bool)) []R {
	result := make([]R, 0, len(collection))

	for i := range collection {
		if r, ok := callback(collection[i], i); ok {
			result = append(result, r)
		}
	}

	return result
}

// FlatMap manipulates a slice and transforms and flattens it to a slice of another type.
// The transform function can either return a slice or a `nil`, and in the `nil` case
// no value is added to the final slice.
// Play: https://go.dev/play/p/pFCF5WVB225
func FlatMap[T, R any](collection []T, transform func(item T, index int) []R) []R {
	result := make([]R, 0, len(collection))

	for i := range collection {
		result = append(result, transform(collection[i], i)...)
	}

	return result
}

// FlatMapErr manipulates a slice and transforms and flattens it to a slice of another type.
// The transform function can either return a slice or a `nil`, and in the `nil` case
// no value is added to the final slice.
// It returns the first error returned by the transform function.
func FlatMapErr[T, R any](collection []T, transform func(item T, index int) ([]R, error)) ([]R, error) {
	result := make([]R, 0, len(collection))

	for i := range collection {
		r, err := transform(collection[i], i)
		if err != nil {
			return nil, err
		}
		result = append(result, r...)
	}

	return result, nil
}

// Reduce reduces collection to a value which is the accumulated result of running each element in collection
// through accumulator, where each successive invocation is supplied the return value of the previous.
// Play: https://go.dev/play/p/CgHYNUpOd1I
func Reduce[T, R any](collection []T, accumulator func(agg R, item T, index int) R, initial R) R {
	for i := range collection {
		initial = accumulator(initial, collection[i], i)
	}

	return initial
}

// ReduceErr reduces collection to a value which is the accumulated result of running each element in collection
// through accumulator, where each successive invocation is supplied the return value of the previous.
// It returns the first error returned by the accumulator function.
func ReduceErr[T, R any](collection []T, accumulator func(agg R, item T, index int) (R, error), initial R) (R, error) {
	for i := range collection {
		result, err := accumulator(initial, collection[i], i)
		if err != nil {
			var zero R
			return zero, err
		}
		initial = result
	}

	return initial, nil
}

// ReduceRight is like Reduce except that it iterates over elements of collection from right to left.
// Play: https://go.dev/play/p/Fq3W70l7wXF
func ReduceRight[T, R any](collection []T, accumulator func(agg R, item T, index int) R, initial R) R {
	for i := len(collection) - 1; i >= 0; i-- {
		initial = accumulator(initial, collection[i], i)
	}

	return initial
}

// ReduceRightErr is like ReduceRight except that the accumulator function can return an error.
// It returns the first error returned by the accumulator function.
func ReduceRightErr[T, R any](collection []T, accumulator func(agg R, item T, index int) (R, error), initial R) (R, error) {
	for i := len(collection) - 1; i >= 0; i-- {
		result, err := accumulator(initial, collection[i], i)
		if err != nil {
			var zero R
			return zero, err
		}
		initial = result
	}

	return initial, nil
}

// ForEach iterates over elements of collection and invokes callback for each element.
// Play: https://go.dev/play/p/oofyiUPRf8t
func ForEach[T any](collection []T, callback func(item T, index int)) {
	for i := range collection {
		callback(collection[i], i)
	}
}

// ForEachWhile iterates over elements of collection and invokes predicate for each element
// collection return value decide to continue or break, like do while().
// Play: https://go.dev/play/p/QnLGt35tnow
func ForEachWhile[T any](collection []T, predicate func(item T, index int) bool) {
	for i := range collection {
		if !predicate(collection[i], i) {
			break
		}
	}
}

// Times invokes the iteratee n times, returning a slice of the results of each invocation.
// The iteratee is invoked with index as argument.
// Play: https://go.dev/play/p/vgQj3Glr6lT
func Times[T any](count int, iteratee func(index int) T) []T {
	result := make([]T, count)

	for i := 0; i < count; i++ {
		result[i] = iteratee(i)
	}

	return result
}

// Uniq returns a duplicate-free version of a slice, in which only the first occurrence of each element is kept.
// The order of result values is determined by the order they occur in the slice.
// Play: https://go.dev/play/p/DTzbeXZ6iEN
func Uniq[T comparable, Slice ~[]T](collection Slice) Slice {
	result := make(Slice, 0, len(collection))
	seen := make(map[T]struct{}, len(collection))

	for i := range collection {
		if _, ok := seen[collection[i]]; ok {
			continue
		}

		seen[collection[i]] = struct{}{}
		result = append(result, collection[i])
	}

	return result
}

// UniqBy returns a duplicate-free version of a slice, in which only the first occurrence of each element is kept.
// The order of result values is determined by the order they occur in the slice. It accepts `iteratee` which is
// invoked for each element in the slice to generate the criterion by which uniqueness is computed.
// Play: https://go.dev/play/p/g42Z3QSb53u
func UniqBy[T any, U comparable, Slice ~[]T](collection Slice, iteratee func(item T) U) Slice {
	result := make(Slice, 0, len(collection))
	seen := make(map[U]struct{}, len(collection))

	for i := range collection {
		key := iteratee(collection[i])

		if _, ok := seen[key]; ok {
			continue
		}

		seen[key] = struct{}{}
		result = append(result, collection[i])
	}

	return result
}

// UniqByErr returns a duplicate-free version of a slice, in which only the first occurrence of each element is kept.
// The order of result values is determined by the order they occur in the slice. It accepts `iteratee` which is
// invoked for each element in the slice to generate the criterion by which uniqueness is computed.
// It returns the first error returned by the iteratee function.
func UniqByErr[T any, U comparable, Slice ~[]T](collection Slice, iteratee func(item T) (U, error)) (Slice, error) {
	result := make(Slice, 0, len(collection))
	seen := make(map[U]struct{}, len(collection))

	for i := range collection {
		key, err := iteratee(collection[i])
		if err != nil {
			return nil, err
		}

		if _, ok := seen[key]; ok {
			continue
		}

		seen[key] = struct{}{}
		result = append(result, collection[i])
	}

	return result, nil
}

// GroupBy returns an object composed of keys generated from the results of running each element of collection through iteratee.
// Play: https://go.dev/play/p/XnQBd_v6brd
func GroupBy[T any, U comparable, Slice ~[]T](collection Slice, iteratee func(item T) U) map[U]Slice {
	result := map[U]Slice{}

	for i := range collection {
		key := iteratee(collection[i])

		result[key] = append(result[key], collection[i])
	}

	return result
}

// GroupByErr returns an object composed of keys generated from the results of running each element of collection through iteratee.
// It returns the first error returned by the iteratee function.
func GroupByErr[T any, U comparable, Slice ~[]T](collection Slice, iteratee func(item T) (U, error)) (map[U]Slice, error) {
	result := map[U]Slice{}

	for i := range collection {
		key, err := iteratee(collection[i])
		if err != nil {
			return nil, err
		}

		result[key] = append(result[key], collection[i])
	}

	return result, nil
}

// GroupByMap returns an object composed of keys generated from the results of running each element of collection through transform.
// Play: https://go.dev/play/p/iMeruQ3_W80
func GroupByMap[T any, K comparable, V any](collection []T, transform func(item T) (K, V)) map[K][]V {
	result := map[K][]V{}

	for i := range collection {
		k, v := transform(collection[i])

		result[k] = append(result[k], v)
	}

	return result
}

// GroupByMapErr returns an object composed of keys generated from the results of running each element of collection through transform.
// It returns the first error returned by the transform function.
func GroupByMapErr[T any, K comparable, V any](collection []T, transform func(item T) (K, V, error)) (map[K][]V, error) {
	result := map[K][]V{}

	for i := range collection {
		k, v, err := transform(collection[i])
		if err != nil {
			return nil, err
		}

		result[k] = append(result[k], v)
	}

	return result, nil
}

// Chunk returns a slice of elements split into groups of length size. If the slice can't be split evenly,
// the final chunk will be the remaining elements.
// Play: https://go.dev/play/p/kEMkFbdu85g
func Chunk[T any, Slice ~[]T](collection Slice, size int) []Slice {
	if size <= 0 {
		panic("lo.Chunk: size must be greater than 0")
	}

	chunksNum := len(collection) / size
	if len(collection)%size != 0 {
		chunksNum++
	}

	result := make([]Slice, 0, chunksNum)

	for i := 0; i < chunksNum; i++ {
		last := (i + 1) * size
		if last > len(collection) {
			last = len(collection)
		}

		// Copy chunk in a new slice, to prevent memory leak and free memory from initial collection.
		newSlice := make(Slice, last-i*size)
		copy(newSlice, collection[i*size:last])
		result = append(result, newSlice)
	}

	return result
}

// PartitionBy returns a slice of elements split into groups. The order of grouped values is
// determined by the order they occur in collection. The grouping is generated from the results
// of running each element of collection through iteratee.
// Play: https://go.dev/play/p/NfQ_nGjkgXW
func PartitionBy[T any, K comparable, Slice ~[]T](collection Slice, iteratee func(item T) K) []Slice {
	result := []Slice{}
	seen := map[K]int{}

	for i := range collection {
		key := iteratee(collection[i])

		resultIndex, ok := seen[key]
		if ok {
			result[resultIndex] = append(result[resultIndex], collection[i])
		} else {
			seen[key] = len(result)
			result = append(result, Slice{collection[i]})
		}
	}

	return result

	// unordered:
	// groups := GroupBy[T, K](collection, iteratee)
	// return Values[K, []T](groups)
}

// PartitionByErr partitions a slice into groups determined by a key computed from each element.
// The order of the partitions is determined by the order they occur in collection. The grouping
// is generated from the results of running each element of collection through iteratee.
// It returns the first error returned by the iteratee function.
func PartitionByErr[T any, K comparable, Slice ~[]T](collection Slice, iteratee func(item T) (K, error)) ([]Slice, error) {
	result := []Slice{}
	seen := map[K]int{}

	for i := range collection {
		key, err := iteratee(collection[i])
		if err != nil {
			return nil, err
		}

		resultIndex, ok := seen[key]
		if ok {
			result[resultIndex] = append(result[resultIndex], collection[i])
		} else {
			seen[key] = len(result)
			result = append(result, Slice{collection[i]})
		}
	}

	return result, nil
}

// Flatten returns a slice a single level deep.
// See also: Concat
// Play: https://go.dev/play/p/rbp9ORaMpjw
func Flatten[T any, Slice ~[]T](collection []Slice) Slice {
	totalLen := 0
	for i := range collection {
		totalLen += len(collection[i])
	}

	result := make(Slice, 0, totalLen)
	for i := range collection {
		result = append(result, collection[i]...)
	}

	return result
}

// Concat returns a new slice containing all the elements in collections. Concat conserves the order of the elements.
// See also: Flatten, Union.
func Concat[T any, Slice ~[]T](collections ...Slice) Slice {
	return Flatten(collections)
}

// Window creates a slice of sliding windows of a given size.
// Each window overlaps with the previous one by size-1 elements.
// This is equivalent to Sliding(collection, size, 1).
func Window[T any, Slice ~[]T](collection Slice, size int) []Slice {
	if size <= 0 {
		panic("lo.Window: size must be greater than 0")
	}
	return Sliding(collection, size, 1)
}

// Sliding creates a slice of sliding windows of a given size with a given step.
// If step is equal to size, windows don't overlap (similar to Chunk).
// If step is less than size, windows overlap.
func Sliding[T any, Slice ~[]T](collection Slice, size, step int) []Slice {
	if size <= 0 {
		panic("lo.Sliding: size must be greater than 0")
	}

	if step <= 0 {
		panic("lo.Sliding: step must be greater than 0")
	}

	n := len(collection) - size
	if n < 0 {
		return []Slice{}
	}

	result := make([]Slice, 0, n/step+1)

	for i := 0; i <= n; i += step {
		window := make(Slice, size)
		copy(window, collection[i:i+size])
		result = append(result, window)
	}

	return result
}

// Interleave round-robin alternating input slices and sequentially appending value at index into result.
// Play: https://go.dev/play/p/-RJkTLQEDVt
func Interleave[T any, Slice ~[]T](collections ...Slice) Slice {
	if len(collections) == 0 {
		return Slice{}
	}

	maxSize := 0
	totalSize := 0
	for i := range collections {
		size := len(collections[i])
		totalSize += size
		if size > maxSize {
			maxSize = size
		}
	}

	if maxSize == 0 {
		return Slice{}
	}

	result := make(Slice, totalSize)

	resultIdx := 0
	for i := 0; i < maxSize; i++ {
		for j := range collections {
			if len(collections[j])-1 < i {
				continue
			}

			result[resultIdx] = collections[j][i]
			resultIdx++
		}
	}

	return result
}

// Shuffle returns a slice of shuffled values. Uses the Fisher-Yates shuffle algorithm.
// Play: https://go.dev/play/p/ZTGG7OUCdnp
//
// Deprecated: use mutable.Shuffle() instead.
func Shuffle[T any, Slice ~[]T](collection Slice) Slice {
	mutable.Shuffle(collection)
	return collection
}

// Reverse reverses a slice so that the first element becomes the last, the second element becomes the second to last, and so on.
// Play: https://go.dev/play/p/iv2e9jslfBM
//
// Deprecated: use mutable.Reverse() instead.
func Reverse[T any, Slice ~[]T](collection Slice) Slice {
	mutable.Reverse(collection)
	return collection
}

// Fill fills elements of a slice with `initial` value.
// Play: https://go.dev/play/p/VwR34GzqEub
func Fill[T Clonable[T], Slice ~[]T](collection Slice, initial T) Slice {
	result := make(Slice, 0, len(collection))

	for range collection {
		result = append(result, initial.Clone())
	}

	return result
}

// Repeat builds a slice with N copies of initial value.
// Play: https://go.dev/play/p/g3uHXbmc3b6
func Repeat[T Clonable[T]](count int, initial T) []T {
	result := make([]T, 0, count)

	for i := 0; i < count; i++ {
		result = append(result, initial.Clone())
	}

	return result
}

// RepeatBy builds a slice with values returned by N calls of callback.
// Play: https://go.dev/play/p/ozZLCtX_hNU
func RepeatBy[T any](count int, callback func(index int) T) []T {
	result := make([]T, 0, count)

	for i := 0; i < count; i++ {
		result = append(result, callback(i))
	}

	return result
}

// RepeatByErr builds a slice with values returned by N calls of callback.
// It returns the first error returned by the callback function.
func RepeatByErr[T any](count int, callback func(index int) (T, error)) ([]T, error) {
	result := make([]T, 0, count)

	for i := 0; i < count; i++ {
		r, err := callback(i)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}

	return result, nil
}

// KeyBy transforms a slice or a slice of structs to a map based on a pivot callback.
// Play: https://go.dev/play/p/ccUiUL_Lnel
func KeyBy[K comparable, V any](collection []V, iteratee func(item V) K) map[K]V {
	result := make(map[K]V, len(collection))

	for i := range collection {
		k := iteratee(collection[i])
		result[k] = collection[i]
	}

	return result
}

// KeyByErr transforms a slice or a slice of structs to a map based on a pivot callback to compute keys.
// Iteratee can return an error to stop iteration immediately.
// Play: https://go.dev/play/p/ccUiUL_Lnel
func KeyByErr[K comparable, V any](collection []V, iteratee func(item V) (K, error)) (map[K]V, error) {
	result := make(map[K]V, len(collection))

	for i := range collection {
		k, err := iteratee(collection[i])
		if err != nil {
			return nil, err
		}
		result[k] = collection[i]
	}

	return result, nil
}

// Associate returns a map containing key-value pairs provided by transform function applied to elements of the given slice.
// If any of two pairs have the same key the last one gets added to the map.
// The order of keys in returned map is not specified and is not guaranteed to be the same from the original slice.
// Play: https://go.dev/play/p/WHa2CfMO3Lr
func Associate[T any, K comparable, V any](collection []T, transform func(item T) (K, V)) map[K]V {
	return AssociateI(collection, func(item T, _ int) (K, V) {
		return transform(item)
	})
}

// AssociateI returns a map containing key-value pairs provided by transform function applied to elements of the given slice.
// If any of two pairs have the same key the last one gets added to the map.
// The order of keys in returned map is not specified and is not guaranteed to be the same from the original slice.
// Play: https://go.dev/play/p/Ugmz6S22rRO
func AssociateI[T any, K comparable, V any](collection []T, transform func(item T, index int) (K, V)) map[K]V {
	result := make(map[K]V, len(collection))

	for index, item := range collection {
		k, v := transform(item, index)
		result[k] = v
	}

	return result
}

// SliceToMap returns a map containing key-value pairs provided by transform function applied to elements of the given slice.
// If any of two pairs have the same key the last one gets added to the map.
// The order of keys in returned map is not specified and is not guaranteed to be the same from the original slice.
// Alias of Associate().
// Play: https://go.dev/play/p/WHa2CfMO3Lr
func SliceToMap[T any, K comparable, V any](collection []T, transform func(item T) (K, V)) map[K]V {
	return Associate(collection, transform)
}

// SliceToMapI returns a map containing key-value pairs provided by transform function applied to elements of the given slice.
// If any of two pairs have the same key the last one gets added to the map.
// The order of keys in returned map is not specified and is not guaranteed to be the same from the original slice.
// Alias of AssociateI().
// Play: https://go.dev/play/p/mMBm5GV3_eq
func SliceToMapI[T any, K comparable, V any](collection []T, transform func(item T, index int) (K, V)) map[K]V {
	return AssociateI(collection, transform)
}

// FilterSliceToMap returns a map containing key-value pairs provided by transform function applied to elements of the given slice.
// If any of two pairs have the same key the last one gets added to the map.
// The order of keys in returned map is not specified and is not guaranteed to be the same from the original slice.
// The third return value of the transform function is a boolean that indicates whether the key-value pair should be included in the map.
// Play: https://go.dev/play/p/2z0rDz2ZSGU
func FilterSliceToMap[T any, K comparable, V any](collection []T, transform func(item T) (K, V, bool)) map[K]V {
	return FilterSliceToMapI(collection, func(item T, _ int) (K, V, bool) {
		return transform(item)
	})
}

// FilterSliceToMapI returns a map containing key-value pairs provided by transform function applied to elements of the given slice.
// If any of two pairs have the same key the last one gets added to the map.
// The order of keys in returned map is not specified and is not guaranteed to be the same from the original slice.
// The third return value of the transform function is a boolean that indicates whether the key-value pair should be included in the map.
// Play: https://go.dev/play/p/mSz_bUIk9aJ
func FilterSliceToMapI[T any, K comparable, V any](collection []T, transform func(item T, index int) (K, V, bool)) map[K]V {
	result := make(map[K]V, len(collection))

	for index, item := range collection {
		k, v, ok := transform(item, index)
		if ok {
			result[k] = v
		}
	}

	return result
}

// Keyify returns a map with each unique element of the slice as a key.
// Play: https://go.dev/play/p/RYhhM_csqIG
func Keyify[T comparable, Slice ~[]T](collection Slice) map[T]struct{} {
	result := make(map[T]struct{}, len(collection))

	for i := range collection {
		result[collection[i]] = struct{}{}
	}

	return result
}

// Drop drops n elements from the beginning of a slice.
// Play: https://go.dev/play/p/JswS7vXRJP2
func Drop[T any, Slice ~[]T](collection Slice, n int) Slice {
	if n < 0 {
		panic("lo.Drop: n must not be negative")
	}

	if len(collection) <= n {
		return make(Slice, 0)
	}

	result := make(Slice, 0, len(collection)-n)

	return append(result, collection[n:]...)
}

// DropRight drops n elements from the end of a slice.
// Play: https://go.dev/play/p/GG0nXkSJJa3
func DropRight[T any, Slice ~[]T](collection Slice, n int) Slice {
	if n < 0 {
		panic("lo.DropRight: n must not be negative")
	}

	if len(collection) <= n {
		return Slice{}
	}

	result := make(Slice, 0, len(collection)-n)
	return append(result, collection[:len(collection)-n]...)
}

// DropWhile drops elements from the beginning of a slice while the predicate returns true.
// Play: https://go.dev/play/p/7gBPYw2IK16
func DropWhile[T any, Slice ~[]T](collection Slice, predicate func(item T) bool) Slice {
	i := 0
	for ; i < len(collection); i++ {
		if !predicate(collection[i]) {
			break
		}
	}

	result := make(Slice, 0, len(collection)-i)
	return append(result, collection[i:]...)
}

// DropRightWhile drops elements from the end of a slice while the predicate returns true.
// Play: https://go.dev/play/p/3-n71oEC0Hz
func DropRightWhile[T any, Slice ~[]T](collection Slice, predicate func(item T) bool) Slice {
	i := len(collection) - 1
	for ; i >= 0; i-- {
		if !predicate(collection[i]) {
			break
		}
	}

	result := make(Slice, 0, i+1)
	return append(result, collection[:i+1]...)
}

// Take takes the first n elements from a slice.
func Take[T any, Slice ~[]T](collection Slice, n int) Slice {
	if n < 0 {
		panic("lo.Take: n must not be negative")
	}

	if n == 0 {
		return make(Slice, 0)
	}

	size := len(collection)
	if size == 0 {
		return make(Slice, 0)
	}

	if n >= size {
		result := make(Slice, size)
		copy(result, collection)
		return result
	}

	result := make(Slice, n)
	copy(result, collection)
	return result
}

// TakeWhile takes elements from the beginning of a slice while the predicate returns true.
func TakeWhile[T any, Slice ~[]T](collection Slice, predicate func(item T) bool) Slice {
	i := 0
	for ; i < len(collection); i++ {
		if !predicate(collection[i]) {
			break
		}
	}

	result := make(Slice, i)
	copy(result, collection[:i])
	return result
}

// DropByIndex drops elements from a slice by the index.
// A negative index will drop elements from the end of the slice.
// Play: https://go.dev/play/p/bPIH4npZRxS
func DropByIndex[T any, Slice ~[]T](collection Slice, indexes ...int) Slice {
	initialSize := len(collection)
	if initialSize == 0 {
		return Slice{}
	}

	// do not change the input
	indexes = append(make([]int, 0, len(indexes)), indexes...)

	for i, index := range indexes {
		if index < 0 {
			indexes[i] += initialSize
		}
	}

	sort.Ints(indexes)

	prev := -1
	indexes = mutable.Filter(indexes, func(index int) bool {
		ok := index != prev && // uniq
			uint(index) < uint(initialSize) // in range

		prev = index
		return ok
	})

	result := make(Slice, 0, initialSize-len(indexes))

	i := 0
	for _, index := range indexes {
		result = append(result, collection[i:index]...)
		i = index + 1
	}

	return append(result, collection[i:]...)
}

// TakeFilter filters elements and takes the first n elements that match the predicate.
// Equivalent to calling Take(Filter(...)), but more efficient as it stops after finding n matches.
func TakeFilter[T any, Slice ~[]T](collection Slice, n int, predicate func(item T, index int) bool) Slice {
	if n < 0 {
		panic("lo.TakeFilter: n must not be negative")
	}

	if n == 0 {
		return make(Slice, 0)
	}

	result := make(Slice, 0, n)
	count := 0

	for i := range collection {
		if predicate(collection[i], i) {
			result = append(result, collection[i])
			count++
			if count >= n {
				break
			}
		}
	}

	return result
}

// Reject is the opposite of Filter, this method returns the elements of collection that predicate does not return true for.
// Play: https://go.dev/play/p/pFCF5WVB225
func Reject[T any, Slice ~[]T](collection Slice, predicate func(item T, index int) bool) Slice {
	result := Slice{}

	for i := range collection {
		if !predicate(collection[i], i) {
			result = append(result, collection[i])
		}
	}

	return result
}

// RejectErr is the opposite of FilterErr, this method returns the elements of collection that predicate does not return true for.
// If the predicate returns an error, iteration stops immediately and returns the error.
// Play: https://go.dev/play/p/pFCF5WVB225
func RejectErr[T any, Slice ~[]T](collection Slice, predicate func(item T, index int) (bool, error)) (Slice, error) {
	result := Slice{}

	for i := range collection {
		match, err := predicate(collection[i], i)
		if err != nil {
			return nil, err
		}
		if !match {
			result = append(result, collection[i])
		}
	}

	return result, nil
}

// RejectMap is the opposite of FilterMap, this method returns a slice obtained after both filtering and mapping using the given callback function.
// The callback function should return two values:
//   - the result of the mapping operation and
//   - whether the result element should be included or not.
//
// Play: https://go.dev/play/p/W9Ug9r0QFkL
func RejectMap[T, R any](collection []T, callback func(item T, index int) (R, bool)) []R {
	result := []R{}

	for i := range collection {
		if r, ok := callback(collection[i], i); !ok {
			result = append(result, r)
		}
	}

	return result
}

// FilterReject mixes Filter and Reject, this method returns two slices, one for the elements of collection that
// predicate returns true for and one for the elements that predicate does not return true for.
// Play: https://go.dev/play/p/lHSEGSznJjB
func FilterReject[T any, Slice ~[]T](collection Slice, predicate func(T, int) bool) (kept, rejected Slice) {
	kept = make(Slice, 0, len(collection))
	rejected = make(Slice, 0, len(collection))

	for i := range collection {
		if predicate(collection[i], i) {
			kept = append(kept, collection[i])
		} else {
			rejected = append(rejected, collection[i])
		}
	}

	return kept, rejected
}

// Count counts the number of elements in the collection that equal value.
// Play: https://go.dev/play/p/Y3FlK54yveC
func Count[T comparable](collection []T, value T) int {
	var count int

	for i := range collection {
		if collection[i] == value {
			count++
		}
	}

	return count
}

// CountBy counts the number of elements in the collection for which predicate is true.
// Play: https://go.dev/play/p/ByQbNYQQi4X
func CountBy[T any](collection []T, predicate func(item T) bool) int {
	var count int

	for i := range collection {
		if predicate(collection[i]) {
			count++
		}
	}

	return count
}

// CountByErr counts the number of elements in the collection for which predicate is true.
// It returns the first error returned by the predicate.
func CountByErr[T any](collection []T, predicate func(item T) (bool, error)) (int, error) {
	var count int

	for i := range collection {
		ok, err := predicate(collection[i])
		if err != nil {
			return 0, err
		}
		if ok {
			count++
		}
	}

	return count, nil
}

// CountValues counts the number of each element in the collection.
// Play: https://go.dev/play/p/-p-PyLT4dfy
func CountValues[T comparable](collection []T) map[T]int {
	result := make(map[T]int)

	for i := range collection {
		result[collection[i]]++
	}

	return result
}

// CountValuesBy counts the number of each element returned from transform function.
// Is equivalent to chaining lo.Map and lo.CountValues.
// Play: https://go.dev/play/p/2U0dG1SnOmS
func CountValuesBy[T any, U comparable](collection []T, transform func(item T) U) map[U]int {
	result := make(map[U]int)

	for i := range collection {
		result[transform(collection[i])]++
	}

	return result
}

// Subset returns a copy of a slice from `offset` up to `length` elements. Like `slice[start:start+length]`, but does not panic on overflow.
// Play: https://go.dev/play/p/tOQu1GhFcog
func Subset[T any, Slice ~[]T](collection Slice, offset int, length uint) Slice {
	size := len(collection)

	if offset < 0 {
		offset = size + offset
		if offset < 0 {
			offset = 0
		}
	}

	if offset > size {
		return Slice{}
	}

	if length > uint(size)-uint(offset) {
		length = uint(size - offset)
	}

	return collection[offset : offset+int(length)]
}

// Slice returns a slice from `start` up to, but not including `end`. Like `slice[start:end]`, but does not panic on overflow.
// Play: https://go.dev/play/p/8XWYhfMMA1h
func Slice[T any, Slice ~[]T](collection Slice, start, end int) Slice {
	if start >= end {
		return Slice{}
	}

	size := len(collection)
	if start < 0 {
		start = 0
	} else if start > size {
		start = size
	}

	if end < 0 {
		end = 0
	} else if end > size {
		end = size
	}

	return collection[start:end]
}

// Replace returns a copy of the slice with the first n non-overlapping instances of old replaced by new.
// Play: https://go.dev/play/p/XfPzmf9gql6
func Replace[T comparable, Slice ~[]T](collection Slice, old, nEw T, n int) Slice {
	result := make(Slice, len(collection))
	copy(result, collection)

	for i := range result {
		if result[i] == old && n != 0 {
			result[i] = nEw
			n--
		}
	}

	return result
}

// ReplaceAll returns a copy of the slice with all non-overlapping instances of old replaced by new.
// Play: https://go.dev/play/p/a9xZFUHfYcV
func ReplaceAll[T comparable, Slice ~[]T](collection Slice, old, nEw T) Slice {
	return Replace(collection, old, nEw, -1)
}

// Clone returns a shallow copy of the collection.
func Clone[T any, Slice ~[]T](collection Slice) Slice {
	// backporting from slices.Clone in Go 1.21
	// when we drop support for Go 1.20, this can be replaced with: return slices.Clone(collection)

	// Preserve nilness in case it matters.
	if collection == nil {
		return nil
	}
	// Avoid s[:0:0] as it leads to unwanted liveness when cloning a
	// zero-length slice of a large array; see https://go.dev/issue/68488.
	return append(Slice{}, collection...)
}

// Compact returns a slice of all non-zero elements.
// Play: https://go.dev/play/p/tXiy-iK6PAc
func Compact[T comparable, Slice ~[]T](collection Slice) Slice {
	var zero T

	result := make(Slice, 0, len(collection))

	for i := range collection {
		if collection[i] != zero {
			result = append(result, collection[i])
		}
	}

	return result
}

// IsSorted checks if a slice is sorted.
// Play: https://go.dev/play/p/mc3qR-t4mcx
func IsSorted[T constraints.Ordered](collection []T) bool {
	for i := 1; i < len(collection); i++ {
		if collection[i-1] > collection[i] {
			return false
		}
	}

	return true
}

// IsSortedBy checks if a slice is sorted by iteratee.
func IsSortedBy[T any, K constraints.Ordered](collection []T, iteratee func(item T) K) bool {
	size := len(collection)

	for i := 0; i < size-1; i++ {
		if iteratee(collection[i]) > iteratee(collection[i+1]) {
			return false
		}
	}

	return true
}

// IsSortedByKey checks if a slice is sorted by iteratee.
//
// Deprecated: Use lo.IsSortedBy instead.
func IsSortedByKey[T any, K constraints.Ordered](collection []T, iteratee func(item T) K) bool {
	return IsSortedBy(collection, iteratee)
}

// Splice inserts multiple elements at index i. A negative index counts back
// from the end of the slice. The helper is protected against overflow errors.
// Play: https://go.dev/play/p/G5_GhkeSUBA
func Splice[T any, Slice ~[]T](collection Slice, i int, elements ...T) Slice {
	sizeCollection := len(collection)
	sizeElements := len(elements)
	output := make(Slice, 0, sizeCollection+sizeElements) // preallocate memory for the output slice

	switch {
	case sizeElements == 0:
		return append(output, collection...) // simple copy
	case i > sizeCollection:
		// positive overflow
		return append(append(output, collection...), elements...)
	case i < -sizeCollection:
		// negative overflow
		return append(append(output, elements...), collection...)
	case i < 0:
		// backward
		i = sizeCollection + i
	}

	return append(append(append(output, collection[:i]...), elements...), collection[i:]...)
}

// Cut slices collection around the first instance of separator, returning the part of collection
// before and after separator. The found result reports whether separator appears in collection.
// If separator does not appear in s, cut returns collection, empty slice of []T, false.
// Play: https://go.dev/play/p/GiL3qhpIP3f
func Cut[T comparable, Slice ~[]T](collection, separator Slice) (before, after Slice, found bool) {
	if len(separator) == 0 {
		return make(Slice, 0), collection, true
	}

	for i := 0; i+len(separator) <= len(collection); i++ {
		match := true
		for j := 0; j < len(separator); j++ {
			if collection[i+j] != separator[j] {
				match = false
				break
			}
		}
		if match {
			return collection[:i], collection[i+len(separator):], true
		}
	}

	return collection, make(Slice, 0), false
}

// CutPrefix returns collection without the provided leading prefix []T
// and reports whether it found the prefix.
// If s doesn't start with prefix, CutPrefix returns collection, false.
// If prefix is the empty []T, CutPrefix returns collection, true.
// Play: https://go.dev/play/p/7Plak4a1ICl
func CutPrefix[T comparable, Slice ~[]T](collection, separator Slice) (after Slice, found bool) {
	if HasPrefix(collection, separator) {
		return collection[len(separator):], true
	}
	return collection, false
}

// CutSuffix returns collection without the provided ending suffix []T and reports
// whether it found the suffix. If s doesn't end with suffix, CutSuffix returns collection, false.
// If suffix is the empty []T, CutSuffix returns collection, true.
// Play: https://go.dev/play/p/7FKfBFvPTaT
func CutSuffix[T comparable, Slice ~[]T](collection, separator Slice) (before Slice, found bool) {
	if HasSuffix(collection, separator) {
		return collection[:len(collection)-len(separator)], true
	}
	return collection, false
}

// Trim removes all the leading and trailing cutset from the collection.
// Play: https://go.dev/play/p/1an9mxLdRG5
func Trim[T comparable, Slice ~[]T](collection, cutset Slice) Slice {
	set := Keyify(cutset)

	i := 0
	for ; i < len(collection); i++ {
		if _, ok := set[collection[i]]; !ok {
			break
		}
	}

	if i >= len(collection) {
		return Slice{}
	}

	j := len(collection) - 1
	for ; j >= 0; j-- {
		if _, ok := set[collection[j]]; !ok {
			break
		}
	}

	result := make(Slice, 0, j+1-i)
	return append(result, collection[i:j+1]...)
}

// TrimLeft removes all the leading cutset from the collection.
// Play: https://go.dev/play/p/74aqfAYLmyi
func TrimLeft[T comparable, Slice ~[]T](collection, cutset Slice) Slice {
	set := Keyify(cutset)

	return DropWhile(collection, func(item T) bool {
		_, ok := set[item]
		return ok
	})
}

// TrimPrefix removes all the leading prefix from the collection.
// Play: https://go.dev/play/p/SHO6X-YegPg
func TrimPrefix[T comparable, Slice ~[]T](collection, prefix Slice) Slice {
	if len(prefix) == 0 {
		return collection
	}

	for HasPrefix(collection, prefix) {
		collection = collection[len(prefix):]
	}

	return collection
}

// TrimRight removes all the trailing cutset from the collection.
// Play: https://go.dev/play/p/MRpAfR6sf0g
func TrimRight[T comparable, Slice ~[]T](collection, cutset Slice) Slice {
	set := Keyify(cutset)

	return DropRightWhile(collection, func(item T) bool {
		_, ok := set[item]
		return ok
	})
}

// TrimSuffix removes all the trailing suffix from the collection.
// Play: https://go.dev/play/p/IjEUrV0iofq
func TrimSuffix[T comparable, Slice ~[]T](collection, suffix Slice) Slice {
	if len(suffix) == 0 {
		return collection
	}

	for HasSuffix(collection, suffix) {
		collection = collection[:len(collection)-len(suffix)]
	}

	return collection
}
