package context

import (
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type FilteredList[T any] struct {
	filteredIndices []int // if nil, we are not filtering

	getList         func() []T
	getFilterFields func(T) []string
	filter          string
}

func NewFilteredList[T any](getList func() []T, getFilterFields func(T) []string) *FilteredList[T] {
	return &FilteredList[T]{
		getList:         getList,
		getFilterFields: getFilterFields,
	}
}

func (self *FilteredList[T]) GetFilter() string {
	return self.filter
}

func (self *FilteredList[T]) SetFilter(filter string) {
	self.filter = filter

	self.applyFilter()
}

func (self *FilteredList[T]) ClearFilter() {
	self.SetFilter("")
}

func (self *FilteredList[T]) IsFiltering() bool {
	return self.filter != ""
}

func (self *FilteredList[T]) GetFilteredList() []T {
	if self.filteredIndices == nil {
		return self.getList()
	}
	return utils.ValuesAtIndices(self.getList(), self.filteredIndices)
}

// TODO: update to just 'Len'
func (self *FilteredList[T]) UnfilteredLen() int {
	return len(self.getList())
}

func (self *FilteredList[T]) applyFilter() {
	if self.filter == "" {
		self.filteredIndices = nil
	} else {
		self.filteredIndices = []int{}
		for i, item := range self.getList() {
			for _, field := range self.getFilterFields(item) {
				if self.match(field, self.filter) {
					self.filteredIndices = append(self.filteredIndices, i)
					break
				}
			}
		}
	}
}

func (self *FilteredList[T]) match(haystack string, needle string) bool {
	return utils.CaseInsensitiveContains(haystack, needle)
}
