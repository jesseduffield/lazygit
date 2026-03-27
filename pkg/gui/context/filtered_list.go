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

	getList          func() []T
	getFilterFields  func(T) []string
	preprocessFilter func(string) string
	filter           string
	lastFilterMode   string // last non-empty gui.filterMode passed to SetFilter; default substring when empty
	lastRegexpPrefix string // last gui.regexpFilterPrefix (OrDefault) passed to SetFilter

	mutex deadlock.Mutex
}

func NewFilteredList[T any](getList func() []T, getFilterFields func(T) []string) *FilteredList[T] {
	return &FilteredList[T]{
		getList:         getList,
		getFilterFields: getFilterFields,
	}
}

func (self *FilteredList[T]) SetPreprocessFilterFunc(preprocessFilter func(string) string) {
	self.preprocessFilter = preprocessFilter
}

func (self *FilteredList[T]) GetFilter() string {
	return self.filter
}

func (self *FilteredList[T]) SetFilter(filter string, useFuzzySearch bool, filterMode string, regexpPrefix string) {
	self.filter = filter
	if filterMode != "" {
		self.lastFilterMode = filterMode
	}
	if regexpPrefix != "" {
		self.lastRegexpPrefix = regexpPrefix
	}

	self.applyFilter(useFuzzySearch, filterMode, regexpPrefix)
}

func (self *FilteredList[T]) ClearFilter() {
	mode := self.lastFilterMode
	if mode == "" {
		mode = "substring"
	}
	prefix := self.lastRegexpPrefix
	if prefix == "" {
		prefix = "re:"
	}
	self.SetFilter("", false, mode, prefix)
}

func (self *FilteredList[T]) ReApplyFilter(useFuzzySearch bool, filterMode string, regexpPrefix string) {
	mode := filterMode
	if mode == "" {
		mode = self.lastFilterMode
	}
	if mode == "" {
		mode = "substring"
	}
	prefix := regexpPrefix
	if prefix == "" {
		prefix = self.lastRegexpPrefix
	}
	if prefix == "" {
		prefix = "re:"
	}
	self.applyFilter(useFuzzySearch, mode, prefix)
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

func (self *FilteredList[T]) applyFilter(useFuzzySearch bool, filterMode string, regexpPrefix string) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	processed := self.filter
	if self.preprocessFilter != nil {
		processed = self.preprocessFilter(self.filter)
	}

	if processed == "" {
		self.filteredIndices = nil
	} else {
		source := &fuzzySource[T]{
			list:            self.getList(),
			getFilterFields: self.getFilterFields,
		}

		prefix := regexpPrefix
		if prefix == "" {
			prefix = "re:"
		}
		pattern, useRegexp := utils.ViewFilterPattern(filterMode, processed, prefix)
		matches := utils.FindFrom(pattern, source, useFuzzySearch, useRegexp)
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
