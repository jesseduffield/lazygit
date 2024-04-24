package utils

import (
	"errors"
)

type HistoryBuffer[T any] struct {
	maxSize int
	items   []T
}

func NewHistoryBuffer[T any](maxSize int) *HistoryBuffer[T] {
	return &HistoryBuffer[T]{
		maxSize: maxSize,
		items:   make([]T, 0, maxSize),
	}
}

func (self *HistoryBuffer[T]) Push(item T) {
	if len(self.items) == self.maxSize {
		self.items = self.items[:len(self.items)-1]
	}
	self.items = append([]T{item}, self.items...)
}

func (self *HistoryBuffer[T]) PeekAt(index int) (T, error) {
	var item T
	if len(self.items) == 0 {
		return item, errors.New("Buffer is empty")
	}
	if len(self.items) <= index || index < -1 {
		return item, errors.New("Index out of range")
	}
	if index == -1 {
		return item, nil
	}
	return self.items[index], nil
}
