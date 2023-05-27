package context

import "github.com/jesseduffield/lazygit/pkg/gui/context/traits"

type ListViewModel[T any] struct {
	*traits.ListCursor
	getModel func() []T
}

func NewListViewModel[T any](getModel func() []T) *ListViewModel[T] {
	self := &ListViewModel[T]{
		getModel: getModel,
	}

	self.ListCursor = traits.NewListCursor(self)

	return self
}

func (self *ListViewModel[T]) Len() int {
	return len(self.getModel())
}

func (self *ListViewModel[T]) GetSelected() T {
	if self.Len() == 0 {
		return Zero[T]()
	}

	return self.getModel()[self.GetSelectedLineIdx()]
}

func (self *ListViewModel[T]) GetItems() []T {
	return self.getModel()
}

func Zero[T any]() T {
	return *new(T)
}
