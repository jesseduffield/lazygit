package context

import "github.com/jesseduffield/lazygit/pkg/gui/context/traits"

type BasicViewModel[T any] struct {
	*traits.ListCursor
	getModel func() []T
}

func NewBasicViewModel[T any](getModel func() []T) *BasicViewModel[T] {
	self := &BasicViewModel[T]{
		getModel: getModel,
	}

	self.ListCursor = traits.NewListCursor(self)

	return self
}

func (self *BasicViewModel[T]) Len() int {
	return len(self.getModel())
}

func (self *BasicViewModel[T]) GetSelected() T {
	if self.Len() == 0 {
		return Zero[T]()
	}

	return self.getModel()[self.GetSelectedLineIdx()]
}

func Zero[T any]() T {
	return *new(T)
}
