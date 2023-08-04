package utils

import "fmt"

type CircularBuffer[T any] struct {
	maxSize int
	Items   []T
}

func NewCircularBuffer[T any](maxSize int) *CircularBuffer[T] {
	return &CircularBuffer[T]{
		maxSize: maxSize,
		Items:   make([]T, 0, maxSize),
	}
}

func (self *CircularBuffer[T]) Push(item T) {
	if len(self.Items) == self.maxSize {
		self.Items = self.Items[:len(self.Items)-1]
	}
	self.Items = append([]T{item}, self.Items...)
}

func (self *CircularBuffer[T]) Pop() (T, error) {
	var item T
	if len(self.Items) == 0 {
		return item, fmt.Errorf("Queue is empty")
	}
	item = self.Items[0]
	self.Items = self.Items[1:]
	return item, nil
}

func (self *CircularBuffer[T]) PeekAt(index int) (T, error) {
	var item T
	if len(self.Items) == 0 {
		return item, fmt.Errorf("Queue is empty")
	}
	length := len(self.Items)
	index = index % length
	if index < 0 {
		index += length
	}
	return self.Items[index], nil
}
