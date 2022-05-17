package context

import (
	"os"
	"strings"
	"sync"

	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"gopkg.in/ozeidan/fuzzy-patricia.v3/patricia"
)

type HasID interface {
	ID() string
}

type FilteredListViewModel[T HasID] struct {
	*BasicViewModel[T]
	getItems func() []T
	trie     *patricia.Trie
}

// I'll just let my trie continue growing based on whatever it's seen previously in the same context given how often you're looking at the same stuff. So I'll need to store a set of ids separately so that I can filter down after that when we get new things. Will that be more efficient than just recreating the trie each time? I suspect so.

func NewFilteredListViewModel[T HasID](
	getItems func() []T,
	getNeedle func() string,
	toString func(T) string,
) *FilteredListViewModel[T] {
	trie := patricia.NewTrie()
	mutex := sync.RWMutex{}
	cacheKey := ""
	cachedGetItems := func() []T {
		items := getItems()

		go func() {
			newCacheKey := HashBy(items, toString)
			if newCacheKey != cacheKey {
				cacheKey = newCacheKey
				mutex.Lock()
				trie = patricia.NewTrie()
				for _, item := range items {
					trie.Set(patricia.Prefix(toString(item)), item)
				}
				mutex.Unlock()
			}
		}()

		return items
	}

	getFilteredModelFnWithTrie := func() []T {
		needle := getNeedle()
		items := cachedGetItems()

		if needle == "" {
			return items
		}

		matches := []T{}
		mutex.Lock()
		_ = trie.VisitSubstring(patricia.Prefix(needle), true, func(prefix patricia.Prefix, item patricia.Item) error {
			matches = append(matches, item.(T))
			return nil
		})
		mutex.Unlock()

		// doing another fuzzy search for good measure
		return utils.FuzzySearchItems(needle, matches, toString)
	}

	return &FilteredListViewModel[T]{
		BasicViewModel: NewBasicViewModel(getFilteredModelFnWithTrie),
		trie:           trie,
		getItems:       getItems,
	}
}

func (self *FilteredListViewModel[T]) FilteredLen() int {
	return self.BasicViewModel.Len()
}

func (self *FilteredListViewModel[T]) Len() int {
	return len(self.getItems())
}

func (self *FilteredListViewModel[T]) TrueSelectedLineIdx() int {
	idx := self.GetSelectedLineIdx()

	// how do I map from one to the other? I need to get the id of one and then find it in the unfiltered list
	id := self.GetSelected().ID()
	for i, item := range self.getItems() {
		if item.ID() == id {
			return i
		}
	}

	return idx
}

func HashBy[T any](items []T, f func(T) string) string {
	// simple implementation for now
	noIndex := func(f func(T) string) func(T, int) string {
		return func(item T, _ int) string {
			return f(item)
		}
	}

	return strings.Join(lo.Map(items, noIndex(f)), "")
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

func newLogger() *logrus.Entry {
	logPath := "/Users/jesseduffieldduffield/Library/Application Support/jesseduffield/lazygit/development.log"
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		panic("unable to log to file") // TODO: don't panic (also, remove this call to the `panic` function)
	}
	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel)
	logger.SetOutput(file)
	return logger.WithFields(logrus.Fields{})
}

var Log = newLogger()
