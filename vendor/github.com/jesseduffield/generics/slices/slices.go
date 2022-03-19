package slices

import (
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
	result := make([]V, len(slice))
	for i, value := range slice {
		result[i] = f(value)
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
func FilterMap[T any, E any](slice []T, test func(T) (bool, E)) []E {
	result := make([]E, 0, len(slice))
	for _, element := range slice {
		ok, mapped := test(element)
		if ok {
			result = append(result, mapped)
		}
	}

	return result
}

// Produces a new slice, leaves the input slice untouched.
func FilterThenMap[T any, E any](slice []T, test func(T) bool, mapFn func(T) E) []E {
	result := make([]E, 0, len(slice))
	for _, element := range slice {
		if test(element) {
			result = append(result, mapFn(element))
		}
	}
	return result
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

// Operates on the input slice. Expected use is to reassign the result to the input slice.
func Move[T any](slice []T, fromIndex int, toIndex int) []T {
	item := slice[fromIndex]
	slice = Remove(slice, fromIndex)
	return slices.Insert(slice, toIndex, item)
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
