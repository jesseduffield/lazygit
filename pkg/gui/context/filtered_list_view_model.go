package context

type FilteredListViewModel[T any] struct {
	*FilteredList[T]
	*ListViewModel[T]
}

func NewFilteredListViewModel[T any](getList func() []T, getFilterFields func(T) []string) *FilteredListViewModel[T] {
	filteredList := &FilteredList[T]{
		getList:         getList,
		getFilterFields: getFilterFields,
	}

	self := &FilteredListViewModel[T]{
		FilteredList: filteredList,
	}

	listViewModel := NewListViewModel(filteredList.GetList)

	self.ListViewModel = listViewModel

	return self
}

// used for type switch
func (self *FilteredListViewModel[T]) IsFilterableContext() {}
