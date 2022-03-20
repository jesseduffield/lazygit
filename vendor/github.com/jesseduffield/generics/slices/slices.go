package slices

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

// This file contains the new functions that do not live in the official slices package.

func Some[T any](slice []T, test func(T) bool) bool {
	for _, value := range slice {
		if test(value) {
			return true
		}
	}

	return false
}

func Every[T any](slice []T, test func(T) bool) bool {
	for _, value := range slice {
		if !test(value) {
			return false
		}
	}

	return true
}

// Produces a new slice, leaves the input slice untouched.
func Map[T any, V any](slice []T, f func(T) V) []V {
	result := make([]V, 0, len(slice))
	for _, value := range slice {
		result = append(result, f(value))
	}

	return result
}

// Produces a new slice, leaves the input slice untouched.
func MapWithIndex[T any, V any](slice []T, f func(T, int) V) []V {
	result := make([]V, 0, len(slice))
	for i, value := range slice {
		result = append(result, f(value, i))
	}

	return result
}

func TryMap[T any, V any](slice []T, f func(T) (V, error)) ([]V, error) {
	result := make([]V, 0, len(slice))
	for _, value := range slice {
		output, err := f(value)
		if err != nil {
			return nil, err
		}
		result = append(result, output)
	}

	return result, nil
}

func TryMapWithIndex[T any, V any](slice []T, f func(T, int) (V, error)) ([]V, error) {
	result := make([]V, 0, len(slice))
	for i, value := range slice {
		output, err := f(value, i)
		if err != nil {
			return nil, err
		}
		result = append(result, output)
	}

	return result, nil
}

// Produces a new slice, leaves the input slice untouched.
func FlatMap[T any, V any](slice []T, f func(T) []V) []V {
	// impossible to know how long this slice will be in the end but the length
	// of the original slice is the lower bound
	result := make([]V, 0, len(slice))
	for _, value := range slice {
		result = append(result, f(value)...)
	}

	return result
}

func FlatMapWithIndex[T any, V any](slice []T, f func(T, int) []V) []V {
	// impossible to know how long this slice will be in the end but the length
	// of the original slice is the lower bound
	result := make([]V, 0, len(slice))
	for i, value := range slice {
		result = append(result, f(value, i)...)
	}

	return result
}

func Flatten[T any](slice [][]T) []T {
	result := make([]T, 0, len(slice))
	for _, subSlice := range slice {
		result = append(result, subSlice...)
	}
	return result
}

func MapInPlace[T any](slice []T, f func(T) T) {
	for i, value := range slice {
		slice[i] = f(value)
	}
}

// Produces a new slice, leaves the input slice untouched.
func Filter[T any](slice []T, test func(T) bool) []T {
	result := make([]T, 0)
	for _, element := range slice {
		if test(element) {
			result = append(result, element)
		}
	}
	return result
}

// Produces a new slice, leaves the input slice untouched.
func FilterWithIndex[T any](slice []T, f func(T, int) bool) []T {
	result := make([]T, 0, len(slice))
	for i, value := range slice {
		if f(value, i) {
			result = append(result, value)
		}
	}

	return result
}

func TryFilter[T any](slice []T, test func(T) (bool, error)) ([]T, error) {
	result := make([]T, 0)
	for _, element := range slice {
		ok, err := test(element)
		if err != nil {
			return nil, err
		}
		if ok {
			result = append(result, element)
		}
	}
	return result, nil
}

func TryFilterWithIndex[T any](slice []T, test func(T, int) (bool, error)) ([]T, error) {
	result := make([]T, 0)
	for i, element := range slice {
		ok, err := test(element, i)
		if err != nil {
			return nil, err
		}
		if ok {
			result = append(result, element)
		}
	}
	return result, nil
}

// Mutates original slice. Intended usage is to reassign the slice result to the input slice.
func FilterInPlace[T any](slice []T, test func(T) bool) []T {
	newLength := 0
	for _, element := range slice {
		if test(element) {
			slice[newLength] = element
			newLength++
		}
	}

	return slice[:newLength]
}

// Produces a new slice, leaves the input slice untouched
func Reverse[T any](slice []T) []T {
	result := make([]T, len(slice))
	for i := range slice {
		result[i] = slice[len(slice)-1-i]
	}
	return result
}

func ReverseInPlace[T any](slice []T) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// Produces a new slice, leaves the input slice untouched.
func FilterMap[T any, E any](slice []T, test func(T) (E, bool)) []E {
	result := make([]E, 0, len(slice))
	for _, element := range slice {
		mapped, ok := test(element)
		if ok {
			result = append(result, mapped)
		}
	}

	return result
}

func FilterMapWithIndex[T any, E any](slice []T, test func(T, int) (E, bool)) []E {
	result := make([]E, 0, len(slice))
	for i, element := range slice {
		mapped, ok := test(element, i)
		if ok {
			result = append(result, mapped)
		}
	}

	return result
}

func TryFilterMap[T any, E any](slice []T, test func(T) (E, bool, error)) ([]E, error) {
	result := make([]E, 0, len(slice))
	for _, element := range slice {
		mapped, ok, err := test(element)
		if err != nil {
			return nil, err
		}
		if ok {
			result = append(result, mapped)
		}
	}

	return result, nil
}

func TryFilterMapWithIndex[T any, E any](slice []T, test func(T, int) (E, bool, error)) ([]E, error) {
	result := make([]E, 0, len(slice))
	for i, element := range slice {
		mapped, ok, err := test(element, i)
		if err != nil {
			return nil, err
		}
		if ok {
			result = append(result, mapped)
		}
	}

	return result, nil
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

// Swaps two elements at the given indices.
// Operates on the input slice.
func Swap[T any](slice []T, index1 int, index2 int) {
	slice[index1], slice[index2] = slice[index2], slice[index1]
}

// Similar to Append but we leave the original slice untouched and return a new slice
func Concat[T any](slice []T, values ...T) []T {
	newSlice := make([]T, 0, len(slice)+len(values))
	newSlice = append(newSlice, slice...)
	newSlice = append(newSlice, values...)
	return newSlice
}

func ContainsFunc[T any](slice []T, f func(T) bool) bool {
	return IndexFunc(slice, f) != -1
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

func MaxBy[T any, V constraints.Ordered](slice []T, f func(T) V) V {
	if len(slice) == 0 {
		return zero[V]()
	}

	max := f(slice[0])
	for _, element := range slice[1:] {
		value := f(element)
		if value > max {
			max = value
		}
	}
	return max
}

func MinBy[T any, V constraints.Ordered](slice []T, f func(T) V) V {
	if len(slice) == 0 {
		return zero[V]()
	}

	min := f(slice[0])
	for _, element := range slice[1:] {
		value := f(element)
		if value < min {
			min = value
		}
	}
	return min
}

func Find[T any](slice []T, f func(T) bool) (T, bool) {
	for _, element := range slice {
		if f(element) {
			return element, true
		}
	}
	return zero[T](), false
}

// Sometimes you need to find an element and then map it to some other value based on
// information you obtained while finding it. This function lets you do that
func FindMap[T any, V any](slice []T, f func(T) (V, bool)) (V, bool) {
	for _, element := range slice {
		if value, ok := f(element); ok {
			return value, true
		}
	}
	return zero[V](), false
}

func ForEach[T any](slice []T, f func(T)) {
	for _, element := range slice {
		f(element)
	}
}

func ForEachWithIndex[T any](slice []T, f func(T, int)) {
	for i, element := range slice {
		f(element, i)
	}
}

func TryForEach[T any](slice []T, f func(T) error) error {
	for _, element := range slice {
		if err := f(element); err != nil {
			return err
		}
	}
	return nil
}

func TryForEachWithIndex[T any](slice []T, f func(T, int) error) error {
	for i, element := range slice {
		if err := f(element, i); err != nil {
			return err
		}
	}
	return nil
}

func Sum[T constraints.Ordered](i []T) T {
	sum := zero[T]()
	for _, value := range i {
		sum += value
	}
	return sum
}

func zero[T any]() T {
	var value T
	return value
}
