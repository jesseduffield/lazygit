package commands

// Tag : A git tag
type Tag struct {
	Name string
}

// GetDisplayStrings returns the display string of a remote
func (r *Tag) GetDisplayStrings(isFocused bool) []string {
	return []string{r.Name}
}
