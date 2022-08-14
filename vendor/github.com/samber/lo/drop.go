package lo

//Drop drops n elements from the beginning of a slice or array.
func Drop[T any](collection []T, n int) []T {
	if len(collection) <= n {
		return make([]T, 0)
	}

	result := make([]T, len(collection)-n)
	for i := n; i < len(collection); i++ {
		result[i-n] = collection[i]
	}

	return result
}

//DropWhile drops elements from the beginning of a slice or array while the predicate returns true.
func DropWhile[T any](collection []T, predicate func(T) bool) []T {
	i := 0
	for ; i < len(collection); i++ {
		if !predicate(collection[i]) {
			break
		}
	}

	result := make([]T, len(collection)-i)

	for j := 0; i < len(collection); i, j = i+1, j+1 {
		result[j] = collection[i]
	}

	return result
}

//DropRight drops n elements from the end of a slice or array.
func DropRight[T any](collection []T, n int) []T {
	if len(collection) <= n {
		return make([]T, 0)
	}

	result := make([]T, len(collection)-n)
	for i := len(collection) - 1 - n; i != 0; i-- {
		result[i] = collection[i]
	}

	return result
}

//DropRightWhile drops elements from the end of a slice or array while the predicate returns true.
func DropRightWhile[T any](collection []T, predicate func(T) bool) []T {
	i := len(collection) - 1
	for ; i >= 0; i-- {
		if !predicate(collection[i]) {
			break
		}
	}

	result := make([]T, i+1)

	for ; i >= 0; i-- {
		result[i] = collection[i]
	}

	return result
}
