package set

import "github.com/jesseduffield/generics/maps"

type Set[T comparable] struct {
	hashMap map[T]bool
}

func New[T comparable]() *Set[T] {
	return &Set[T]{hashMap: make(map[T]bool)}
}

func NewFromSlice[T comparable](slice []T) *Set[T] {
	hashMap := make(map[T]bool)
	for _, value := range slice {
		hashMap[value] = true
	}

	return &Set[T]{hashMap: hashMap}
}

func (s *Set[T]) Add(values ...T) {
	for _, value := range values {
		s.hashMap[value] = true
	}
}

func (s *Set[T]) Remove(value T) {
	delete(s.hashMap, value)
}

func (s *Set[T]) RemoveSlice(slice []T) {
	for _, value := range slice {
		s.Remove(value)
	}
}

func (s *Set[T]) Includes(value T) bool {
	return s.hashMap[value]
}

// output slice is not necessarily in the same order that items were added
func (s *Set[T]) ToSlice() []T {
	return maps.Keys(s.hashMap)
}
