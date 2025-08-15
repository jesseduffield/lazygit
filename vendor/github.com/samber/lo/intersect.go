package lo

// Contains returns true if an element is present in a collection.
func Contains[T comparable](collection []T, element T) bool {
	for _, item := range collection {
		if item == element {
			return true
		}
	}

	return false
}

// ContainsBy returns true if predicate function return true.
func ContainsBy[T any](collection []T, predicate func(T) bool) bool {
	for _, item := range collection {
		if predicate(item) {
			return true
		}
	}

	return false
}

// Every returns true if all elements of a subset are contained into a collection or if the subset is empty.
func Every[T comparable](collection []T, subset []T) bool {
	for _, elem := range subset {
		if !Contains(collection, elem) {
			return false
		}
	}

	return true
}

// EveryBy returns true if the predicate returns true for all of the elements in the collection or if the collection is empty.
func EveryBy[V any](collection []V, predicate func(V) bool) bool {
	for _, v := range collection {
		if !predicate(v) {
			return false
		}
	}

	return true
}

// Some returns true if at least 1 element of a subset is contained into a collection.
// If the subset is empty Some returns false.
func Some[T comparable](collection []T, subset []T) bool {
	for _, elem := range subset {
		if Contains(collection, elem) {
			return true
		}
	}

	return false
}

// SomeBy returns true if the predicate returns true for any of the elements in the collection.
// If the collection is empty SomeBy returns false.
func SomeBy[V any](collection []V, predicate func(V) bool) bool {
	for _, v := range collection {
		if predicate(v) {
			return true
		}
	}

	return false
}

// None returns true if no element of a subset are contained into a collection or if the subset is empty.
func None[V comparable](collection []V, subset []V) bool {
	for _, elem := range subset {
		if Contains(collection, elem) {
			return false
		}
	}

	return true
}

// NoneBy returns true if the predicate returns true for none of the elements in the collection or if the collection is empty.
func NoneBy[V any](collection []V, predicate func(V) bool) bool {
	for _, v := range collection {
		if predicate(v) {
			return false
		}
	}

	return true
}

// Intersect returns the intersection between two collections.
func Intersect[T comparable](list1 []T, list2 []T) []T {
	result := []T{}
	seen := map[T]struct{}{}

	for _, elem := range list1 {
		seen[elem] = struct{}{}
	}

	for _, elem := range list2 {
		if _, ok := seen[elem]; ok {
			result = append(result, elem)
		}
	}

	return result
}

// Difference returns the difference between two collections.
// The first value is the collection of element absent of list2.
// The second value is the collection of element absent of list1.
func Difference[T comparable](list1 []T, list2 []T) ([]T, []T) {
	left := []T{}
	right := []T{}

	seenLeft := map[T]struct{}{}
	seenRight := map[T]struct{}{}

	for _, elem := range list1 {
		seenLeft[elem] = struct{}{}
	}

	for _, elem := range list2 {
		seenRight[elem] = struct{}{}
	}

	for _, elem := range list1 {
		if _, ok := seenRight[elem]; !ok {
			left = append(left, elem)
		}
	}

	for _, elem := range list2 {
		if _, ok := seenLeft[elem]; !ok {
			right = append(right, elem)
		}
	}

	return left, right
}

// Union returns all distinct elements from both collections.
// result returns will not change the order of elements relatively.
func Union[T comparable](list1 []T, list2 []T) []T {
	result := []T{}

	seen := map[T]struct{}{}
	hasAdd := map[T]struct{}{}

	for _, e := range list1 {
		seen[e] = struct{}{}
	}

	for _, e := range list2 {
		seen[e] = struct{}{}
	}

	for _, e := range list1 {
		if _, ok := seen[e]; ok {
			result = append(result, e)
			hasAdd[e] = struct{}{}
		}
	}

	for _, e := range list2 {
		if _, ok := hasAdd[e]; ok {
			continue
		}
		if _, ok := seen[e]; ok {
			result = append(result, e)
		}
	}

	return result
}

// Without returns slice excluding all given values.
func Without[T comparable](collection []T, exclude ...T) []T {
	result := make([]T, 0, len(collection))
	for _, e := range collection {
		if !Contains(exclude, e) {
			result = append(result, e)
		}
	}
	return result
}

// WithoutEmpty returns slice excluding empty values.
func WithoutEmpty[T comparable](collection []T) []T {
	var empty T

	result := make([]T, 0, len(collection))
	for _, e := range collection {
		if e != empty {
			result = append(result, e)
		}
	}

	return result
}
