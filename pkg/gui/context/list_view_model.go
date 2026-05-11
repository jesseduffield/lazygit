package context

import (
	"github.com/jesseduffield/lazygit/pkg/gui/context/traits"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

type HasID interface {
	ID() string
}

type ListViewModel[T HasID] struct {
	*traits.ListCursor
	getModel func() []T
}

func NewListViewModel[T HasID](getModel func() []T) *ListViewModel[T] {
	self := &ListViewModel[T]{
		getModel: getModel,
	}

	self.ListCursor = traits.NewListCursor(func() int { return len(getModel()) })

	return self
}

func (self *ListViewModel[T]) GetSelected() T {
	if self.Len() == 0 {
		return Zero[T]()
	}

	return self.getModel()[self.GetSelectedLineIdx()]
}

func (self *ListViewModel[T]) GetSelectedItemId() string {
	if self.Len() == 0 {
		return ""
	}

	return self.GetSelected().ID()
}

func (self *ListViewModel[T]) GetSelectedItems() ([]T, int, int) {
	if self.Len() == 0 {
		return nil, -1, -1
	}

	startIdx, endIdx := self.GetSelectionRange()

	return self.getModel()[startIdx : endIdx+1], startIdx, endIdx
}

func (self *ListViewModel[T]) GetSelectedItemIds() ([]string, int, int) {
	selectedItems, startIdx, endIdx := self.GetSelectedItems()

	ids := lo.Map(selectedItems, func(item T, _ int) string {
		return item.ID()
	})

	return ids, startIdx, endIdx
}

func (self *ListViewModel[T]) GetItems() []T {
	return self.getModel()
}

func Zero[T any]() T {
	return *new(T)
}

func (self *ListViewModel[T]) GetItem(index int) types.HasUrn {
	item := self.getModel()[index]
	return any(item).(types.HasUrn)
}
