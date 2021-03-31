package filetree

import (
	"os"
	"strings"
)

type CollapsedPaths map[string]bool

func (cp CollapsedPaths) ExpandToPath(path string) {
	// need every directory along the way
	split := strings.Split(path, string(os.PathSeparator))
	for i := range split {
		dir := strings.Join(split[0:i+1], string(os.PathSeparator))
		cp[dir] = false
	}
}

func (cp CollapsedPaths) IsCollapsed(path string) bool {
	return cp[path]
}

func (cp CollapsedPaths) ToggleCollapsed(path string) {
	cp[path] = !cp[path]
}
