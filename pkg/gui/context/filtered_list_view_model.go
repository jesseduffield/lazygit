package context

type FilteredListViewModel[T any] struct {
	*FilteredList[T]
	*ListViewModel[T]
}

func NewFilteredListViewModel[T any](getList func() []T, getFilterFields func(T) []string) *FilteredListViewModel[T] {
	filteredList := NewFilteredList(getList, getFilterFields)

	self := &FilteredListViewModel[T]{
		FilteredList: filteredList,
	}

	listViewModel := NewListViewModel(filteredList.GetFilteredList)

	self.ListViewModel = listViewModel

	return self
}

// used for type switch
func (self *FilteredListViewModel[T]) IsFilterableContext() {}

func (self *FilteredListViewModel[T]) ClearFilter() {
	// Set the selected line index to the unfiltered index of the currently selected line,
	// so that the current item is still selected after the filter is cleared.
	unfilteredIndex := self.FilteredList.UnfilteredIndex(self.GetSelectedLineIdx())

	self.FilteredList.ClearFilter()

	self.SetSelectedLineIdx(unfilteredIndex)
}
