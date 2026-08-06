package utils

type Stack[T any] struct {
	stack []T
}

func (self *Stack[T]) Push(v T) {
	self.stack = append(self.stack, v)
}

// Pop returns the zero value of T if the stack is empty.
func (self *Stack[T]) Pop() T {
	var v T
	if len(self.stack) == 0 {
		return v
	}
	n := len(self.stack) - 1
	v = self.stack[n]
	self.stack = self.stack[:n]
	return v
}

func (self *Stack[T]) IsEmpty() bool {
	return len(self.stack) == 0
}

func (self *Stack[T]) Clear() {
	self.stack = nil
}
