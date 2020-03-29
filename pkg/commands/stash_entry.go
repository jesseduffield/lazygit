package commands

import "fmt"

// StashEntry : A git stash entry
type StashEntry struct {
	Index int
	Name  string
}

func (s *StashEntry) RefName() string {
	return fmt.Sprintf("stash@{%d}", s.Index)
}
