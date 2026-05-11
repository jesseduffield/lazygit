package lo

import (
	"math"

	"github.com/samber/lo/internal/constraints"
)

// Range creates a slice of numbers (positive and/or negative) with given length.
// Play: https://go.dev/play/p/0r6VimXAi9H
func Range(elementNum int) []int {
	step := Ternary(elementNum < 0, -1, 1)
	length := elementNum * step
	result := make([]int, length)
	for i, j := 0, 0; i < length; i, j = i+1, j+step {
		result[i] = j
	}
	return result
}

// RangeFrom creates a slice of numbers from start with specified length.
// Play: https://go.dev/play/p/0r6VimXAi9H
func RangeFrom[T constraints.Integer | constraints.Float](start T, elementNum int) []T {
	step := Ternary(elementNum < 0, -1, 1)
	length := elementNum * step
	result := make([]T, length)
	for i, j := 0, start; i < length; i, j = i+1, j+T(step) {
		result[i] = j
	}
	return result
}

// RangeWithSteps creates a slice of numbers (positive and/or negative) progressing from start up to, but not including end.
// step set to zero will return an empty slice.
// Play: https://go.dev/play/p/0r6VimXAi9H
func RangeWithSteps[T constraints.Integer | constraints.Float](start, end, step T) []T {
	if start == end || step == 0 {
		return []T{}
	}

	capacity := func(count, delta T) int {
		// Use math.Ceil instead of (count-1)/delta+1 because integer division
		// fails for floats (e.g., 5.5/2.5=2.2 → ceil=3, not 2).
		return int(math.Ceil(float64(count) / float64(delta)))
	}

	if start < end {
		if step < 0 {
			return []T{}
		}

		result := make([]T, 0, capacity(end-start, step))
		for i := start; i < end; i += step {
			result = append(result, i)
		}
		return result
	}
	if step > 0 {
		return []T{}
	}

	result := make([]T, 0, capacity(start-end, -step))
	for i := start; i > end; i += step {
		result = append(result, i)
	}
	return result
}

// Clamp clamps number within the inclusive lower and upper bounds.
// Play: https://go.dev/play/p/RU4lJNC2hlI
func Clamp[T constraints.Ordered](value, mIn, mAx T) T {
	if value < mIn {
		return mIn
	} else if value > mAx {
		return mAx
	}
	return value
}

// Sum sums the values in a collection. If collection is empty 0 is returned.
// Play: https://go.dev/play/p/upfeJVqs4Bt
func Sum[T constraints.Float | constraints.Integer | constraints.Complex](collection []T) T {
	var sum T
	for i := range collection {
		sum += collection[i]
	}
	return sum
}

// SumBy summarizes the values in a collection using the given return value from the iteration function. If collection is empty 0 is returned.
// Play: https://go.dev/play/p/Dz_a_7jN_ca
func SumBy[T any, R constraints.Float | constraints.Integer | constraints.Complex](collection []T, iteratee func(item T) R) R {
	var sum R
	for i := range collection {
		sum += iteratee(collection[i])
	}
	return sum
}

// SumByErr summarizes the values in a collection using the given return value from the iteration function.
// If the iteratee returns an error, iteration stops and the error is returned.
// If collection is empty 0 and nil error are returned.
func SumByErr[T any, R constraints.Float | constraints.Integer | constraints.Complex](collection []T, iteratee func(item T) (R, error)) (R, error) {
	var sum R
	for i := range collection {
		v, err := iteratee(collection[i])
		if err != nil {
			return sum, err
		}
		sum += v
	}
	return sum, nil
}

// Product gets the product of the values in a collection. If collection is empty 1 is returned.
// Play: https://go.dev/play/p/2_kjM_smtAH
func Product[T constraints.Float | constraints.Integer | constraints.Complex](collection []T) T {
	var product T = 1
	for i := range collection {
		product *= collection[i]
	}
	return product
}

// ProductBy summarizes the values in a collection using the given return value from the iteration function. If collection is empty 1 is returned.
// Play: https://go.dev/play/p/wadzrWr9Aer
func ProductBy[T any, R constraints.Float | constraints.Integer | constraints.Complex](collection []T, iteratee func(item T) R) R {
	var product R = 1
	for i := range collection {
		product *= iteratee(collection[i])
	}
	return product
}

// ProductByErr summarizes the values in a collection using the given return value from the iteration function.
// If the iteratee returns an error, iteration stops and the error is returned.
// If collection is empty 1 and nil error are returned.
func ProductByErr[T any, R constraints.Float | constraints.Integer | constraints.Complex](collection []T, iteratee func(item T) (R, error)) (R, error) {
	var product R = 1
	for i := range collection {
		v, err := iteratee(collection[i])
		if err != nil {
			return product, err
		}
		product *= v
	}
	return product, nil
}

// Mean calculates the mean of a collection of numbers.
// Play: https://go.dev/play/p/tPURSuteUsP
func Mean[T constraints.Float | constraints.Integer](collection []T) T {
	length := T(len(collection))
	if length == 0 {
		return 0
	}
	sum := Sum(collection)
	return sum / length
}

// MeanBy calculates the mean of a collection of numbers using the given return value from the iteration function.
// Play: https://go.dev/play/p/j7TsVwBOZ7P
func MeanBy[T any, R constraints.Float | constraints.Integer](collection []T, iteratee func(item T) R) R {
	length := R(len(collection))
	if length == 0 {
		return 0
	}
	sum := SumBy(collection, iteratee)
	return sum / length
}

// MeanByErr calculates the mean of a collection of numbers using the given return value from the iteration function.
// If the iteratee returns an error, iteration stops and the error is returned.
// If collection is empty 0 and nil error are returned.
func MeanByErr[T any, R constraints.Float | constraints.Integer](collection []T, iteratee func(item T) (R, error)) (R, error) {
	length := R(len(collection))
	if length == 0 {
		return 0, nil
	}
	sum, err := SumByErr(collection, iteratee)
	if err != nil {
		return 0, err
	}
	return sum / length, nil
}

// Mode returns the mode (most frequent value) of a collection.
// If multiple values have the same highest frequency, then multiple values are returned.
// If the collection is empty, then the zero value of T is returned.
func Mode[T constraints.Integer | constraints.Float](collection []T) []T {
	length := T(len(collection))
	if length == 0 {
		return []T{}
	}

	mode := make([]T, 0)
	maxFreq := 0
	frequency := make(map[T]int)

	for _, item := range collection {
		frequency[item]++
		count := frequency[item]

		if count > maxFreq {
			maxFreq = count
			mode = []T{item}
		} else if count == maxFreq {
			mode = append(mode, item)
		}
	}

	return mode
}
