package utils

import "sync"

// A simple string pool implementation that can help reduce memory usage for
// cases where the same string is used multiple times.
type StringPool struct {
	sync.Map
}

func (self *StringPool) Add(s string) *string {
	poolEntry, _ := self.LoadOrStore(s, &s)
	return poolEntry.(*string)
}
