package searching

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/sirupsen/logrus"
)

type searchMode int

const (
	inactive searchMode = iota
	inPrompt
	inView
)

type Searching struct {
	view       *gocui.View
	mode       searchMode
	needle     string
	contextKey types.ContextKey
	log        *logrus.Entry
}

func New(log *logrus.Entry) *Searching {
	return &Searching{log: log}
}

func (self *Searching) SearchingInContext(contextKey types.ContextKey) bool {
	if self.mode == inactive || self.needle == "" {
		return false
	}

	return self.contextKey != "" && self.contextKey == contextKey
}

func (self *Searching) NewSearchingState() bool {
	return self.mode == inView
}

func (self *Searching) Escape() {
	self.mode = inactive
	self.contextKey = ""
	self.needle = ""

	if self.view != nil {
		self.view.ClearSearch()
		self.view = nil
	}
}

func (self *Searching) OnSearch(needle string) error {
	self.needle = needle
	self.mode = inView

	view := self.view
	if view == nil {
		return nil
	}

	if err := view.Search(self.needle); err != nil {
		return err
	}

	return nil
}

func (self *Searching) OnSearchPrompt(view *gocui.View, contextKey types.ContextKey) {
	self.mode = inPrompt
	self.view = view
	self.contextKey = contextKey
}

func (self *Searching) GetSearchString() string {
	return self.needle
}

func (self *Searching) InPrompt() bool {
	return self.mode == inPrompt
}

func (self *Searching) Active() bool {
	return self.mode == inView || self.InPrompt()
}
