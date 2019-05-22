package commands

// StashEntry : A git stash entry
type StashEntry struct {
	Index         int
	Name          string
	DisplayString string
}

// GetDisplayStrings returns the display string of branch
func (s *StashEntry) GetDisplayStrings(isFocused bool) []string {
	return []string{s.DisplayString}
}
