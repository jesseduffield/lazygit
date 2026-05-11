package types

type SearchType int

const (
	SearchTypeNone SearchType = iota
	// searching is where matches are highlighted but the content is not filtered down
	SearchTypeSearch
	// filter is where the list is filtered down to only matches
	SearchTypeFilter
)

// TODO: could we remove this entirely?
type SearchState struct {
	Context         Context
	PrevSearchIndex int
}

func NewSearchState() *SearchState {
	return &SearchState{PrevSearchIndex: -1}
}

func (self *SearchState) SearchType() SearchType {
	switch self.Context.(type) {
	case IFilterableContext:
		return SearchTypeFilter
	case ISearchableContext:
		return SearchTypeSearch
	default:
		return SearchTypeNone
	}
}
