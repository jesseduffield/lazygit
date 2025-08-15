package lo

import "golang.org/x/exp/constraints"

// Range creates an array of numbers (positive and/or negative) with given length.
// Play: https://go.dev/play/p/0r6VimXAi9H
func Range(elementNum int) []int {
	length := If(elementNum < 0, -elementNum).Else(elementNum)
	result := make([]int, length)
	step := If(elementNum < 0, -1).Else(1)
	for i, j := 0, 0; i < length; i, j = i+1, j+step {
		result[i] = j
	}
	return result
}

// RangeFrom creates an array of numbers from start with specified length.
// Play: https://go.dev/play/p/0r6VimXAi9H
func RangeFrom[T constraints.Integer | constraints.Float](start T, elementNum int) []T {
	length := If(elementNum < 0, -elementNum).Else(elementNum)
	result := make([]T, length)
	step := If(elementNum < 0, -1).Else(1)
	for i, j := 0, start; i < length; i, j = i+1, j+T(step) {
		result[i] = j
	}
	return result
}

// RangeWithSteps creates an array of numbers (positive and/or negative) progressing from start up to, but not including end.
// step set to zero will return empty array.
// Play: https://go.dev/play/p/0r6VimXAi9H
func RangeWithSteps[T constraints.Integer | constraints.Float](start, end, step T) []T {
	result := []T{}
	if start == end || step == 0 {
		return result
	}
	if start < end {
		if step < 0 {
			return result
		}
		for i := start; i < end; i += step {
			result = append(result, i)
		}
		return result
	}
	if step > 0 {
		return result
	}
	for i := start; i > end; i += step {
		result = append(result, i)
	}
	return result
}

// Clamp clamps number within the inclusive lower and upper bounds.
// Play: https://go.dev/play/p/RU4lJNC2hlI
func Clamp[T constraints.Ordered](value T, min T, max T) T {
	if value < min {
		return min
	} else if value > max {
		return max
	}
	return value
}

// SumBy summarizes the values in a collection using the given return value from the iteration function. If collection is empty 0 is returned.
// Play: https://go.dev/play/p/Dz_a_7jN_ca
func SumBy[T any, R constraints.Float | constraints.Integer](collection []T, iteratee func(T) R) R {
	var sum R = 0
	for _, item := range collection {
		sum = sum + iteratee(item)
	}
	return sum
}
