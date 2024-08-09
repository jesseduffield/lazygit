package utils

import "golang.org/x/exp/slices"

// NextIndex returns the index of the element that comes after the given number
func NextIndex(numbers []int, currentNumber int) int {
	for index, number := range numbers {
		if number > currentNumber {
			return index
		}
	}
	return len(numbers) - 1
}

// PrevIndex returns the index that comes before the given number, cycling if we reach the end
func PrevIndex(numbers []int, currentNumber int) int {
	end := len(numbers) - 1
	for i := end; i >= 0; i-- {
		if numbers[i] < currentNumber {
			return i
		}
	}
	return 0
}

// NextIntInCycle returns the next int in a slice, returning to the first index if we've reached the end
func NextIntInCycle(sl []int, current int) int {
	for i, val := range sl {
		if val == current {
			if i == len(sl)-1 {
				return sl[0]
			}
			return sl[i+1]
		}
	}
	return sl[0]
}

// PrevIntInCycle returns the prev int in a slice, returning to the first index if we've reached the end
func PrevIntInCycle(sl []int, current int) int {
	for i, val := range sl {
		if val == current {
			if i > 0 {
				return sl[i-1]
			}
			return sl[len(sl)-1]
		}
	}
	return sl[len(sl)-1]
}

func StringArraysOverlap(strArrA []string, strArrB []string) bool {
	for _, first := range strArrA {
		for _, second := range strArrB {
			if first == second {
				return true
			}
		}
	}

	return false
}

func Limit(values []string, limit int) []string {
	if len(values) > limit {
		return values[:limit]
	}
	return values
}

func LimitStr(value string, limit int) string {
	n := 0
	for i := range value {
		if n >= limit {
			return value[:i]
		}
		n++
	}
	return value
}

// Similar to a regular GroupBy, except that each item can be grouped under multiple keys,
// so the callback returns a slice of keys instead of just one key.
func MuiltiGroupBy[T any, K comparable](slice []T, f func(T) []K) map[K][]T {
	result := map[K][]T{}
	for _, item := range slice {
		for _, key := range f(item) {
			if _, ok := result[key]; !ok {
				result[key] = []T{item}
			} else {
				result[key] = append(result[key], item)
			}
		}
	}
	return result
}

// Returns a new slice with the element at index 'from' moved to index 'to'.
// Does not mutate original slice.
func MoveElement[T any](slice []T, from int, to int) []T {
	newSlice := make([]T, len(slice))
	copy(newSlice, slice)

	if from == to {
		return newSlice
	}

	if from < to {
		copy(newSlice[from:to+1], newSlice[from+1:to+1])
	} else {
		copy(newSlice[to+1:from+1], newSlice[to:from])
	}

	newSlice[to] = slice[from]

	return newSlice
}

func ValuesAtIndices[T any](slice []T, indices []int) []T {
	result := make([]T, len(indices))
	for i, index := range indices {
		// gracefully handling the situation where the index is out of bounds
		if index < len(slice) {
			result[i] = slice[index]
		}
	}
	return result
}

// returns two slices: the first is for elements that pass the test, the second for those that don't.
func Partition[T any](slice []T, test func(T) bool) ([]T, []T) {
	left := make([]T, 0, len(slice))
	right := make([]T, 0, len(slice))

	for _, value := range slice {
		if test(value) {
			left = append(left, value)
		} else {
			right = append(right, value)
		}
	}

	return left, right
}

// Prepends items to the beginning of a slice.
// E.g. Prepend([]int{1,2}, 3, 4) = []int{3,4,1,2}
// Mutates original slice. Intended usage is to reassign the slice result to the input slice.
func Prepend[T any](slice []T, values ...T) []T {
	return append(values, slice...)
}

// Removes the element at the given index. Intended usage is to reassign the result to the input slice.
func Remove[T any](slice []T, index int) []T {
	return slices.Delete(slice, index, index+1)
}

// Removes the element at the 'fromIndex' and then inserts it at 'toIndex'.
// Operates on the input slice. Expected use is to reassign the result to the input slice.
func Move[T any](slice []T, fromIndex int, toIndex int) []T {
	item := slice[fromIndex]
	slice = Remove(slice, fromIndex)
	return slices.Insert(slice, toIndex, item)
}

// Pops item from the end of the slice and returns it, along with the updated slice
// Mutates original slice. Intended usage is to reassign the slice result to the input slice.
func Pop[T any](slice []T) (T, []T) {
	index := len(slice) - 1
	value := slice[index]
	slice = slice[0:index]
	return value, slice
}

// Shifts item from the beginning of the slice and returns it, along with the updated slice.
// Mutates original slice. Intended usage is to reassign the slice result to the input slice.
func Shift[T any](slice []T) (T, []T) {
	value := slice[0]
	slice = slice[1:]
	return value, slice
}

// Compares two slices for equality
func EqualSlices[T comparable](slice1 []T, slice2 []T) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}

	return true
}
