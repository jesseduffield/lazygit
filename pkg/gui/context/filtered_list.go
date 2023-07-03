package context

import (
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/sasha-s/go-deadlock"
)

type FilteredList[T any] struct {
	filteredIndices []int // if nil, we are not filtering

	getList         func() []T
	getFilterFields func(T) []string
	filter          string

	mutex *deadlock.Mutex
}

func NewFilteredList[T any](getList func() []T, getFilterFields func(T) []string) *FilteredList[T] {
	return &FilteredList[T]{
		getList:         getList,
		getFilterFields: getFilterFields,
		mutex:           &deadlock.Mutex{},
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
	self.mutex.Lock()
	defer self.mutex.Unlock()

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
	return utils.CaseAwareContains(haystack, needle)
}

func (self *FilteredList[T]) UnfilteredIndex(index int) int {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if self.filteredIndices == nil {
		return index
	}

	// we use -1 when there are no items
	if index == -1 {
		return -1
	}

	return self.filteredIndices[index]
}
