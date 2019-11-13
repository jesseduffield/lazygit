package commands

// Remote : A git remote
type Remote struct {
	Name     string
	Urls     []string
	Selected bool
}

// GetDisplayStrings returns the display string of a remote
func (r *Remote) GetDisplayStrings(isFocused bool) []string {

	return []string{r.Name}
}
