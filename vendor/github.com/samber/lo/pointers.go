package lo

// ToPtr returns a pointer copy of value.
func ToPtr[T any](x T) *T {
	return &x
}

// ToPtr returns a slice of pointer copy of value.
func ToSlicePtr[T any](collection []T) []*T {
	return Map(collection, func (x T, _ int) *T {
		return &x
	})
}

// Empty returns an empty value.
func Empty[T any]() T {
	var t T
	return t
}
