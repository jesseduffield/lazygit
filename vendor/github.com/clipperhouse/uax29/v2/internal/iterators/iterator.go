package iterators

import "github.com/clipperhouse/stringish"

type SplitFunc[T stringish.Interface] func(T, bool) (int, T, error)

// Iterator is a generic iterator for words that are either []byte or string.
// Iterate while Next() is true, and access the word via Value().
type Iterator[T stringish.Interface] struct {
	split SplitFunc[T]
	data  T
	start int
	pos   int
}

// New creates a new Iterator for the given data and SplitFunc.
func New[T stringish.Interface](split SplitFunc[T], data T) *Iterator[T] {
	return &Iterator[T]{
		split: split,
		data:  data,
	}
}

// SetText sets the text for the iterator to operate on, and resets all state.
func (iter *Iterator[T]) SetText(data T) {
	iter.data = data
	iter.start = 0
	iter.pos = 0
}

// Split sets the SplitFunc for the Iterator.
func (iter *Iterator[T]) Split(split SplitFunc[T]) {
	iter.split = split
}

// Next advances the iterator to the next token. It returns false when there
// are no remaining tokens or an error occurred.
func (iter *Iterator[T]) Next() bool {
	if iter.pos == len(iter.data) {
		return false
	}
	if iter.pos > len(iter.data) {
		panic("SplitFunc advanced beyond the end of the data")
	}

	iter.start = iter.pos

	advance, _, err := iter.split(iter.data[iter.pos:], true)
	if err != nil {
		panic(err)
	}
	if advance <= 0 {
		panic("SplitFunc returned a zero or negative advance")
	}

	iter.pos += advance
	if iter.pos > len(iter.data) {
		panic("SplitFunc advanced beyond the end of the data")
	}

	return true
}

// Value returns the current token.
func (iter *Iterator[T]) Value() T {
	return iter.data[iter.start:iter.pos]
}

// Start returns the byte position of the current token in the original data.
func (iter *Iterator[T]) Start() int {
	return iter.start
}

// End returns the byte position after the current token in the original data.
func (iter *Iterator[T]) End() int {
	return iter.pos
}

// Reset resets the iterator to the beginning of the data.
func (iter *Iterator[T]) Reset() {
	iter.start = 0
	iter.pos = 0
}

func (iter *Iterator[T]) First() T {
	if len(iter.data) == 0 {
		return iter.data
	}
	advance, _, err := iter.split(iter.data, true)
	if err != nil {
		panic(err)
	}
	if advance <= 0 {
		panic("SplitFunc returned a zero or negative advance")
	}
	if advance > len(iter.data) {
		panic("SplitFunc advanced beyond the end of the data")
	}
	return iter.data[:advance]
}
