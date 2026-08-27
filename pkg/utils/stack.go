package utils

type Stack[T any] struct {
	stack []T
}

func (self *Stack[T]) Push(item T) {
	self.stack = append(self.stack, item)
}

func (self *Stack[T]) Pop() T {
	if len(self.stack) == 0 {
		var zero T
		return zero
	}
	n := len(self.stack) - 1
	last := self.stack[n]
	self.stack = self.stack[:n]
	return last
}

func (self *Stack[T]) IsEmpty() bool {
	return len(self.stack) == 0
}

func (self *Stack[T]) Clear() {
	self.stack = nil
}
