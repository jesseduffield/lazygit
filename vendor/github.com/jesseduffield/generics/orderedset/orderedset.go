package orderedset

import (
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type OrderedSet[T comparable] struct {
	om *orderedmap.OrderedMap[T, bool]
}

func New[T comparable]() *OrderedSet[T] {
	return &OrderedSet[T]{om: orderedmap.New[T, bool]()}
}

func NewFromSlice[T comparable](slice []T) *OrderedSet[T] {
	result := &OrderedSet[T]{om: orderedmap.New[T, bool](len(slice))}
	result.Add(slice...)
	return result
}

func (os *OrderedSet[T]) Add(values ...T) {
	for _, value := range values {
		os.om.Set(value, true)
	}
}

func (os *OrderedSet[T]) Remove(value T) {
	os.om.Delete(value)
}

func (os *OrderedSet[T]) RemoveSlice(slice []T) {
	for _, value := range slice {
		os.Remove(value)
	}
}

func (os *OrderedSet[T]) Includes(value T) bool {
	return os.om.Value(value)
}

func (os *OrderedSet[T]) Len() int {
	return os.om.Len()
}

func (os *OrderedSet[T]) ToSliceFromOldest() []T {
	// TODO: can be simplified to
	//   return os.om.KeysFromOldest()
	// when we update to a newer version of go-ordered-map
	result := make([]T, 0, os.Len())
	for pair := os.om.Oldest(); pair != nil; pair = pair.Next() {
		result = append(result, pair.Key)
	}
	return result
}

func (os *OrderedSet[T]) ToSliceFromNewest() []T {
	// TODO: can be simplified to
	//   return os.om.KeysFromNewest()
	// when we update to a newer version of go-ordered-map
	result := make([]T, 0, os.Len())
	for pair := os.om.Newest(); pair != nil; pair = pair.Prev() {
		result = append(result, pair.Key)
	}
	return result
}
