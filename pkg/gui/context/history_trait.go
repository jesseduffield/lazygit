package context

import (
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// Maintains a list of strings that have previously been searched/filtered for
type SearchHistory struct {
	history *utils.HistoryBuffer[string]
}

func NewSearchHistory() *SearchHistory {
	return &SearchHistory{
		history: utils.NewHistoryBuffer[string](1000),
	}
}

func (self *SearchHistory) GetSearchHistory() *utils.HistoryBuffer[string] {
	return self.history
}
