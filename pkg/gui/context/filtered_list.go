package context

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/sahilm/fuzzy"
	"github.com/samber/lo"
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

type fuzzySource[T any] struct {
	list            []T
	getFilterFields func(T) []string
}

var _ fuzzy.Source = &fuzzySource[string]{}

func (self *fuzzySource[T]) String(i int) string {
	return strings.Join(self.getFilterFields(self.list[i]), " ")
}

func (self *fuzzySource[T]) Len() int {
	return len(self.list)
}

func (self *FilteredList[T]) applyFilter() {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if self.filter == "" {
		self.filteredIndices = nil
	} else {
		source := &fuzzySource[T]{
			list:            self.getList(),
			getFilterFields: self.getFilterFields,
		}

		matches := fuzzy.FindFrom(self.filter, source)
		self.filteredIndices = lo.Map(matches, func(match fuzzy.Match, _ int) int {
			return match.Index
		})
	}
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
