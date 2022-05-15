package context

import "github.com/jesseduffield/lazygit/pkg/utils"

type FilteredListViewModel[T any] struct {
	*BasicViewModel[T]
}

func NewFilteredListViewModel[T any](
	getItems func() []T,
	getNeedle func() string,
	toString func(T) string,
) *FilteredListViewModel[T] {
	return &FilteredListViewModel[T]{
		BasicViewModel: NewBasicViewModel(getFilteredModelFn(getItems, getNeedle, toString)),
	}
}

func getFilteredModelFn[T any](
	getItems func() []T,
	getNeedle func() string,
	toString func(T) string,
) func() []T {
	return func() []T {
		needle := getNeedle()
		items := getItems()
		if needle == "" {
			return items
		}

		return utils.FuzzySearchItems(needle, items, toString)
	}
}
