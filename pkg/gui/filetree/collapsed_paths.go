package filetree

type CollapsedPaths map[string]bool

func (cp CollapsedPaths) ExpandToPath(path string) {
	// need every directory along the way
	splitPath := split(path)
	for i := range splitPath {
		dir := join(splitPath[0 : i+1])
		cp[dir] = false
	}
}

func (cp CollapsedPaths) IsCollapsed(path string) bool {
	return cp[path]
}

func (cp CollapsedPaths) ToggleCollapsed(path string) {
	cp[path] = !cp[path]
}
