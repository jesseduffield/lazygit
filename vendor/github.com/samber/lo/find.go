package lo

import (
	"time"

	"github.com/samber/lo/internal/constraints"
	"github.com/samber/lo/internal/xrand"
)

// IndexOf returns the index at which the first occurrence of a value is found in a slice or -1
// if the value cannot be found.
// Play: https://go.dev/play/p/Eo7W0lvKTky
func IndexOf[T comparable](collection []T, element T) int {
	for i := range collection {
		if collection[i] == element {
			return i
		}
	}

	return -1
}

// LastIndexOf returns the index at which the last occurrence of a value is found in a slice or -1
// if the value cannot be found.
// Play: https://go.dev/play/p/Eo7W0lvKTky
func LastIndexOf[T comparable](collection []T, element T) int {
	length := len(collection)

	for i := length - 1; i >= 0; i-- {
		if collection[i] == element {
			return i
		}
	}

	return -1
}

// HasPrefix returns true if the collection has the prefix.
// Play: https://go.dev/play/p/SrljzVDpMQM
func HasPrefix[T comparable](collection, prefix []T) bool {
	if len(collection) < len(prefix) {
		return false
	}

	for i := range prefix {
		if collection[i] != prefix[i] {
			return false
		}
	}

	return true
}

// HasSuffix returns true if the collection has the suffix.
// Play: https://go.dev/play/p/bJeLetQNAON
func HasSuffix[T comparable](collection, suffix []T) bool {
	if len(collection) < len(suffix) {
		return false
	}

	for i := range suffix {
		if collection[len(collection)-len(suffix)+i] != suffix[i] {
			return false
		}
	}

	return true
}

// Find searches for an element in a slice based on a predicate. Returns element and true if element was found.
// Play: https://go.dev/play/p/Eo7W0lvKTky
func Find[T any](collection []T, predicate func(item T) bool) (T, bool) {
	for i := range collection {
		if predicate(collection[i]) {
			return collection[i], true
		}
	}

	var result T
	return result, false
}

// FindErr searches for an element in a slice based on a predicate that can return an error.
// Returns the element and nil error if the element is found.
// Returns zero value and nil error if the element is not found.
// If the predicate returns an error, iteration stops immediately and returns zero value and the error.
func FindErr[T any](collection []T, predicate func(item T) (bool, error)) (T, error) {
	for i := range collection {
		matches, err := predicate(collection[i])
		if err != nil {
			var result T
			return result, err
		}
		if matches {
			return collection[i], nil
		}
	}

	var result T
	return result, nil
}

// FindIndexOf searches for an element in a slice based on a predicate and returns the index and true.
// Returns -1 and false if the element is not found.
// Play: https://go.dev/play/p/XWSEM4Ic_t0
func FindIndexOf[T any](collection []T, predicate func(item T) bool) (T, int, bool) {
	for i := range collection {
		if predicate(collection[i]) {
			return collection[i], i, true
		}
	}

	var result T
	return result, -1, false
}

// FindLastIndexOf searches for the last element in a slice based on a predicate and returns the index and true.
// Returns -1 and false if the element is not found.
// Play: https://go.dev/play/p/dPiMRtJ6cUx
func FindLastIndexOf[T any](collection []T, predicate func(item T) bool) (T, int, bool) {
	length := len(collection)

	for i := length - 1; i >= 0; i-- {
		if predicate(collection[i]) {
			return collection[i], i, true
		}
	}

	var result T
	return result, -1, false
}

// FindOrElse searches for an element in a slice based on a predicate. Returns the element if found or a given fallback value otherwise.
// Play: https://go.dev/play/p/Eo7W0lvKTky
func FindOrElse[T any](collection []T, fallback T, predicate func(item T) bool) T {
	for i := range collection {
		if predicate(collection[i]) {
			return collection[i]
		}
	}

	return fallback
}

// FindKey returns the key of the first value matching.
// Play: https://go.dev/play/p/Bg0w1VDPYXx
func FindKey[K, V comparable](object map[K]V, value V) (K, bool) {
	for k, v := range object {
		if v == value {
			return k, true
		}
	}

	return Empty[K](), false
}

// FindKeyBy returns the key of the first element predicate returns true for.
// Play: https://go.dev/play/p/9IbiPElcyo8
func FindKeyBy[K comparable, V any](object map[K]V, predicate func(key K, value V) bool) (K, bool) {
	for k, v := range object {
		if predicate(k, v) {
			return k, true
		}
	}

	return Empty[K](), false
}

// FindUniques returns a slice with all the elements that appear in the collection only once.
// The order of result values is determined by the order they occur in the collection.
func FindUniques[T comparable, Slice ~[]T](collection Slice) Slice {
	isDupl := make(map[T]bool, len(collection))

	duplicates := 0

	for i := range collection {
		duplicated, seen := isDupl[collection[i]]
		if !duplicated {
			isDupl[collection[i]] = seen

			if seen {
				duplicates++
			}
		}
	}

	result := make(Slice, 0, len(isDupl)-duplicates)

	for i := range collection {
		if duplicated := isDupl[collection[i]]; !duplicated {
			result = append(result, collection[i])
		}
	}

	return result
}

// FindUniquesBy returns a slice with all the elements that appear in the collection only once.
// The order of result values is determined by the order they occur in the slice. It accepts `iteratee` which is
// invoked for each element in the slice to generate the criterion by which uniqueness is computed.
func FindUniquesBy[T any, U comparable, Slice ~[]T](collection Slice, iteratee func(item T) U) Slice {
	isDupl := make(map[U]bool, len(collection))

	duplicates := 0

	for i := range collection {
		key := iteratee(collection[i])

		duplicated, seen := isDupl[key]
		if !duplicated {
			isDupl[key] = seen

			if seen {
				duplicates++
			}
		}
	}

	result := make(Slice, 0, len(isDupl)-duplicates)

	for i := range collection {
		key := iteratee(collection[i])

		if duplicated := isDupl[key]; !duplicated {
			result = append(result, collection[i])
		}
	}

	return result
}

// FindDuplicates returns a slice with the first occurrence of each duplicated element in the collection.
// The order of result values is determined by the order they occur in the collection.
func FindDuplicates[T comparable, Slice ~[]T](collection Slice) Slice {
	isDupl := make(map[T]bool, len(collection))

	duplicates := 0

	for i := range collection {
		duplicated, seen := isDupl[collection[i]]
		if !duplicated {
			isDupl[collection[i]] = seen

			if seen {
				duplicates++
			}
		}
	}

	result := make(Slice, 0, duplicates)

	for i := range collection {
		if duplicated := isDupl[collection[i]]; duplicated {
			result = append(result, collection[i])
			isDupl[collection[i]] = false
		}
	}

	return result
}

// FindDuplicatesBy returns a slice with the first occurrence of each duplicated element in the collection.
// The order of result values is determined by the order they occur in the slice. It accepts `iteratee` which is
// invoked for each element in the slice to generate the criterion by which uniqueness is computed.
func FindDuplicatesBy[T any, U comparable, Slice ~[]T](collection Slice, iteratee func(item T) U) Slice {
	isDupl := make(map[U]bool, len(collection))

	duplicates := 0

	for i := range collection {
		key := iteratee(collection[i])

		duplicated, seen := isDupl[key]
		if !duplicated {
			isDupl[key] = seen

			if seen {
				duplicates++
			}
		}
	}

	result := make(Slice, 0, duplicates)

	for i := range collection {
		key := iteratee(collection[i])

		if duplicated := isDupl[key]; duplicated {
			result = append(result, collection[i])
			isDupl[key] = false
		}
	}

	return result
}

// FindDuplicatesByErr returns a slice with the first occurrence of each duplicated element in the collection.
// The order of result values is determined by the order they occur in the slice. It accepts `iteratee` which is
// invoked for each element in the slice to generate the criterion by which uniqueness is computed.
// If the iteratee returns an error, iteration stops immediately and the error is returned with a nil slice.
func FindDuplicatesByErr[T any, U comparable, Slice ~[]T](collection Slice, iteratee func(item T) (U, error)) (Slice, error) {
	isDupl := make(map[U]bool, len(collection))

	duplicates := 0

	// First pass: identify duplicates
	for i := range collection {
		key, err := iteratee(collection[i])
		if err != nil {
			var result Slice
			return result, err
		}

		duplicated, seen := isDupl[key]
		if !duplicated {
			isDupl[key] = seen

			if seen {
				duplicates++
			}
		}
	}

	result := make(Slice, 0, duplicates)

	// Second pass: collect first occurrences of duplicates
	for i := range collection {
		key, err := iteratee(collection[i])
		if err != nil {
			var result Slice
			return result, err
		}

		if duplicated := isDupl[key]; duplicated {
			result = append(result, collection[i])
			isDupl[key] = false
		}
	}

	return result, nil
}

// Min search the minimum value of a collection.
// Returns zero value when the collection is empty.
// Play: https://go.dev/play/p/r6e-Z8JozS8
func Min[T constraints.Ordered](collection []T) T {
	var mIn T

	if len(collection) == 0 {
		return mIn
	}

	mIn = collection[0]

	for i := 1; i < len(collection); i++ {
		item := collection[i]

		if item < mIn {
			mIn = item
		}
	}

	return mIn
}

// MinIndex search the minimum value of a collection and the index of the minimum value.
// Returns (zero value, -1) when the collection is empty.
func MinIndex[T constraints.Ordered](collection []T) (T, int) {
	var (
		mIn   T
		index int
	)

	if len(collection) == 0 {
		return mIn, -1
	}

	mIn = collection[0]

	for i := 1; i < len(collection); i++ {
		item := collection[i]

		if item < mIn {
			mIn = item
			index = i
		}
	}

	return mIn, index
}

// MinBy search the minimum value of a collection using the given comparison function.
// If several values of the collection are equal to the smallest value, returns the first such value.
// Returns zero value when the collection is empty.
func MinBy[T any](collection []T, less func(a, b T) bool) T {
	var mIn T

	if len(collection) == 0 {
		return mIn
	}

	mIn = collection[0]

	for i := 1; i < len(collection); i++ {
		item := collection[i]

		if less(item, mIn) {
			mIn = item
		}
	}

	return mIn
}

// MinByErr search the minimum value of a collection using the given comparison function.
// If several values of the collection are equal to the smallest value, returns the first such value.
// Returns zero value and nil error when the collection is empty.
// If the comparison function returns an error, iteration stops and the error is returned.
func MinByErr[T any](collection []T, less func(a, b T) (bool, error)) (T, error) {
	var mIn T

	if len(collection) == 0 {
		return mIn, nil
	}

	mIn = collection[0]

	for i := 1; i < len(collection); i++ {
		item := collection[i]

		isLess, err := less(item, mIn)
		if err != nil {
			var zero T
			return zero, err
		}
		if isLess {
			mIn = item
		}
	}

	return mIn, nil
}

// MinIndexBy search the minimum value of a collection using the given comparison function and the index of the minimum value.
// If several values of the collection are equal to the smallest value, returns the first such value.
// Returns (zero value, -1) when the collection is empty.
func MinIndexBy[T any](collection []T, less func(a, b T) bool) (T, int) {
	var (
		mIn   T
		index int
	)

	if len(collection) == 0 {
		return mIn, -1
	}

	mIn = collection[0]

	for i := 1; i < len(collection); i++ {
		item := collection[i]

		if less(item, mIn) {
			mIn = item
			index = i
		}
	}

	return mIn, index
}

// MinIndexByErr search the minimum value of a collection using the given comparison function and the index of the minimum value.
// If several values of the collection are equal to the smallest value, returns the first such value.
// Returns (zero value, -1) when the collection is empty.
// Comparison function can return an error to stop iteration immediately.
func MinIndexByErr[T any](collection []T, less func(a, b T) (bool, error)) (T, int, error) {
	var (
		mIn   T
		index int
	)

	if len(collection) == 0 {
		return mIn, -1, nil
	}

	mIn = collection[0]

	for i := 1; i < len(collection); i++ {
		item := collection[i]

		isLess, err := less(item, mIn)
		if err != nil {
			var zero T
			return zero, -1, err
		}

		if isLess {
			mIn = item
			index = i
		}
	}

	return mIn, index, nil
}

// Earliest search the minimum time.Time of a collection.
// Returns zero value when the collection is empty.
func Earliest(times ...time.Time) time.Time {
	var mIn time.Time

	if len(times) == 0 {
		return mIn
	}

	mIn = times[0]

	for i := 1; i < len(times); i++ {
		item := times[i]

		if item.Before(mIn) {
			mIn = item
		}
	}

	return mIn
}

// EarliestBy search the minimum time.Time of a collection using the given iteratee function.
// Returns zero value when the collection is empty.
func EarliestBy[T any](collection []T, iteratee func(item T) time.Time) T {
	var earliest T

	if len(collection) == 0 {
		return earliest
	}

	earliest = collection[0]
	earliestTime := iteratee(collection[0])

	for i := 1; i < len(collection); i++ {
		itemTime := iteratee(collection[i])

		if itemTime.Before(earliestTime) {
			earliest = collection[i]
			earliestTime = itemTime
		}
	}

	return earliest
}

// EarliestByErr search the minimum time.Time of a collection using the given iteratee function.
// Returns zero value and nil error when the collection is empty.
// If the iteratee returns an error, iteration stops and the error is returned.
func EarliestByErr[T any](collection []T, iteratee func(item T) (time.Time, error)) (T, error) {
	var earliest T

	if len(collection) == 0 {
		return earliest, nil
	}

	earliestTime, err := iteratee(collection[0])
	if err != nil {
		return earliest, err
	}
	earliest = collection[0]

	for i := 1; i < len(collection); i++ {
		itemTime, err := iteratee(collection[i])
		if err != nil {
			return earliest, err
		}

		if itemTime.Before(earliestTime) {
			earliest = collection[i]
			earliestTime = itemTime
		}
	}

	return earliest, nil
}

// Max searches the maximum value of a collection.
// Returns zero value when the collection is empty.
// Play: https://go.dev/play/p/r6e-Z8JozS8
func Max[T constraints.Ordered](collection []T) T {
	var mAx T

	if len(collection) == 0 {
		return mAx
	}

	mAx = collection[0]

	for i := 1; i < len(collection); i++ {
		item := collection[i]

		if item > mAx {
			mAx = item
		}
	}

	return mAx
}

// MaxIndex searches the maximum value of a collection and the index of the maximum value.
// Returns (zero value, -1) when the collection is empty.
func MaxIndex[T constraints.Ordered](collection []T) (T, int) {
	var (
		mAx   T
		index int
	)

	if len(collection) == 0 {
		return mAx, -1
	}

	mAx = collection[0]

	for i := 1; i < len(collection); i++ {
		item := collection[i]

		if item > mAx {
			mAx = item
			index = i
		}
	}

	return mAx, index
}

// MaxBy search the maximum value of a collection using the given comparison function.
// If several values of the collection are equal to the greatest value, returns the first such value.
// Returns zero value when the collection is empty.
//
// Note: the comparison function is inconsistent with most languages, since we use the opposite of the usual convention.
// See https://github.com/samber/lo/issues/129
//
// Play: https://go.dev/play/p/JW1qu-ECwF7
func MaxBy[T any](collection []T, greater func(a, b T) bool) T {
	var mAx T

	if len(collection) == 0 {
		return mAx
	}

	mAx = collection[0]

	for i := 1; i < len(collection); i++ {
		item := collection[i]

		if greater(item, mAx) {
			mAx = item
		}
	}

	return mAx
}

// MaxByErr search the maximum value of a collection using the given comparison function.
// If several values of the collection are equal to the greatest value, returns the first such value.
// Returns zero value and nil error when the collection is empty.
// If the comparison function returns an error, iteration stops and the error is returned.
//
// Note: the comparison function is inconsistent with most languages, since we use the opposite of the usual convention.
// See https://github.com/samber/lo/issues/129
func MaxByErr[T any](collection []T, greater func(a, b T) (bool, error)) (T, error) {
	var mAx T

	if len(collection) == 0 {
		return mAx, nil
	}

	mAx = collection[0]

	for i := 1; i < len(collection); i++ {
		item := collection[i]

		isGreater, err := greater(item, mAx)
		if err != nil {
			return mAx, err
		}
		if isGreater {
			mAx = item
		}
	}

	return mAx, nil
}

// MaxIndexBy search the maximum value of a collection using the given comparison function and the index of the maximum value.
// If several values of the collection are equal to the greatest value, returns the first such value.
// Returns (zero value, -1) when the collection is empty.
//
// Note: the comparison function is inconsistent with most languages, since we use the opposite of the usual convention.
// See https://github.com/samber/lo/issues/129
//
// Play: https://go.dev/play/p/uaUszc-c9QK
func MaxIndexBy[T any](collection []T, greater func(a, b T) bool) (T, int) {
	var (
		mAx   T
		index int
	)

	if len(collection) == 0 {
		return mAx, -1
	}

	mAx = collection[0]

	for i := 1; i < len(collection); i++ {
		item := collection[i]

		if greater(item, mAx) {
			mAx = item
			index = i
		}
	}

	return mAx, index
}

// MaxIndexByErr search the maximum value of a collection using the given comparison function and the index of the maximum value.
// If several values of the collection are equal to the greatest value, returns the first such value.
// Returns (zero value, -1, nil) when the collection is empty.
// If the comparison function returns an error, iteration stops and the error is returned.
//
// Note: the comparison function is inconsistent with most languages, since we use the opposite of the usual convention.
// See https://github.com/samber/lo/issues/129
func MaxIndexByErr[T any](collection []T, greater func(a, b T) (bool, error)) (T, int, error) {
	var (
		mAx   T
		index int
	)

	if len(collection) == 0 {
		return mAx, -1, nil
	}

	mAx = collection[0]

	for i := 1; i < len(collection); i++ {
		item := collection[i]

		isGreater, err := greater(item, mAx)
		if err != nil {
			var zero T
			return zero, -1, err
		}
		if isGreater {
			mAx = item
			index = i
		}
	}

	return mAx, index, nil
}

// Latest search the maximum time.Time of a collection.
// Returns zero value when the collection is empty.
func Latest(times ...time.Time) time.Time {
	var mAx time.Time

	if len(times) == 0 {
		return mAx
	}

	mAx = times[0]

	for i := 1; i < len(times); i++ {
		item := times[i]

		if item.After(mAx) {
			mAx = item
		}
	}

	return mAx
}

// LatestBy search the maximum time.Time of a collection using the given iteratee function.
// Returns zero value when the collection is empty.
func LatestBy[T any](collection []T, iteratee func(item T) time.Time) T {
	var latest T

	if len(collection) == 0 {
		return latest
	}

	latest = collection[0]
	latestTime := iteratee(collection[0])

	for i := 1; i < len(collection); i++ {
		itemTime := iteratee(collection[i])

		if itemTime.After(latestTime) {
			latest = collection[i]
			latestTime = itemTime
		}
	}

	return latest
}

// LatestByErr search the maximum time.Time of a collection using the given iteratee function.
// Returns zero value and nil error when the collection is empty.
// If the iteratee returns an error, iteration stops and the error is returned.
func LatestByErr[T any](collection []T, iteratee func(item T) (time.Time, error)) (T, error) {
	var latest T

	if len(collection) == 0 {
		return latest, nil
	}

	latestTime, err := iteratee(collection[0])
	if err != nil {
		return latest, err
	}
	latest = collection[0]

	for i := 1; i < len(collection); i++ {
		itemTime, err := iteratee(collection[i])
		if err != nil {
			return latest, err
		}

		if itemTime.After(latestTime) {
			latest = collection[i]
			latestTime = itemTime
		}
	}

	return latest, nil
}

// First returns the first element of a collection and check for availability of the first element.
// Play: https://go.dev/play/p/ul45Z0y2EFO
func First[T any](collection []T) (T, bool) {
	length := len(collection)

	if length == 0 {
		var t T
		return t, false
	}

	return collection[0], true
}

// FirstOrEmpty returns the first element of a collection or zero value if empty.
// Play: https://go.dev/play/p/ul45Z0y2EFO
func FirstOrEmpty[T any](collection []T) T {
	i, _ := First(collection)
	return i
}

// FirstOr returns the first element of a collection or the fallback value if empty.
// Play: https://go.dev/play/p/ul45Z0y2EFO
func FirstOr[T any](collection []T, fallback T) T {
	i, ok := First(collection)
	if !ok {
		return fallback
	}

	return i
}

// Last returns the last element of a collection or error if empty.
// Play: https://go.dev/play/p/ul45Z0y2EFO
func Last[T any](collection []T) (T, bool) {
	length := len(collection)

	if length == 0 {
		var t T
		return t, false
	}

	return collection[length-1], true
}

// LastOrEmpty returns the last element of a collection or zero value if empty.
// Play: https://go.dev/play/p/ul45Z0y2EFO
func LastOrEmpty[T any](collection []T) T {
	i, _ := Last(collection)
	return i
}

// LastOr returns the last element of a collection or the fallback value if empty.
// Play: https://go.dev/play/p/ul45Z0y2EFO
func LastOr[T any](collection []T, fallback T) T {
	i, ok := Last(collection)
	if !ok {
		return fallback
	}

	return i
}

// Nth returns the element at index `nth` of collection. If `nth` is negative, the nth element
// from the end is returned. An error is returned when nth is out of slice bounds.
// Play: https://go.dev/play/p/sHoh88KWt6B
func Nth[T any, N constraints.Integer](collection []T, nth N) (T, error) {
	value, ok := sliceNth(collection, nth)

	return value, Validate(ok, "nth: %d out of slice bounds", nth)
}

func sliceNth[T any, N constraints.Integer](collection []T, nth N) (T, bool) {
	n := int(nth)
	l := len(collection)
	if n >= l || -n > l {
		return Empty[T](), false
	}

	if n >= 0 {
		return collection[n], true
	}
	return collection[l+n], true
}

// NthOr returns the element at index `nth` of collection.
// If `nth` is negative, it returns the nth element from the end.
// If `nth` is out of slice bounds, it returns the fallback value instead of an error.
// Play: https://go.dev/play/p/sHoh88KWt6B
func NthOr[T any, N constraints.Integer](collection []T, nth N, fallback T) T {
	value, ok := sliceNth(collection, nth)
	if !ok {
		return fallback
	}
	return value
}

// NthOrEmpty returns the element at index `nth` of collection.
// If `nth` is negative, it returns the nth element from the end.
// If `nth` is out of slice bounds, it returns the zero value (empty value) for that type.
// Play: https://go.dev/play/p/sHoh88KWt6B
func NthOrEmpty[T any, N constraints.Integer](collection []T, nth N) T {
	value, _ := sliceNth(collection, nth)
	return value
}

// randomIntGenerator is a function that should return a random integer in the range [0, n)
// where n is the argument passed to the randomIntGenerator.
type randomIntGenerator func(n int) int

// Sample returns a random item from collection.
// Play: https://go.dev/play/p/vCcSJbh5s6l
func Sample[T any](collection []T) T {
	return SampleBy(collection, xrand.IntN)
}

// SampleBy returns a random item from collection, using randomIntGenerator as the random index generator.
// Play: https://go.dev/play/p/HDmKmMgq0XN
func SampleBy[T any](collection []T, randomIntGenerator randomIntGenerator) T {
	size := len(collection)
	if size == 0 {
		return Empty[T]()
	}
	return collection[randomIntGenerator(size)]
}

// Samples returns N random unique items from collection.
// Play: https://go.dev/play/p/vCcSJbh5s6l
func Samples[T any, Slice ~[]T](collection Slice, count int) Slice {
	return SamplesBy(collection, count, xrand.IntN)
}

// SamplesBy returns N random unique items from collection, using randomIntGenerator as the random index generator.
// Play: https://go.dev/play/p/HDmKmMgq0XN
func SamplesBy[T any, Slice ~[]T](collection Slice, count int, randomIntGenerator randomIntGenerator) Slice {
	if count <= 0 {
		return Slice{}
	}

	size := len(collection)

	if size < count {
		count = size
	}

	indexes := Range(size)
	results := make(Slice, count)

	for i := range results {
		n := len(indexes)

		index := randomIntGenerator(n)
		results[i] = collection[indexes[index]]

		// Removes index.
		// It is faster to swap with last element and remove it.
		indexes[index] = indexes[n-1]
		indexes = indexes[:n-1]
	}

	return results
}
