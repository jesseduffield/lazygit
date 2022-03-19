package filetree

import "github.com/jesseduffield/generics/set"

type CollapsedPaths struct {
	collapsedPaths *set.Set[string]
}

func NewCollapsedPaths() *CollapsedPaths {
	return &CollapsedPaths{
		collapsedPaths: set.New[string](),
	}
}

func (self *CollapsedPaths) ExpandToPath(path string) {
	// need every directory along the way
	splitPath := split(path)
	for i := range splitPath {
		dir := join(splitPath[0 : i+1])
		self.collapsedPaths.Remove(dir)
	}
}

func (self *CollapsedPaths) IsCollapsed(path string) bool {
	return self.collapsedPaths.Includes(path)
}

func (self *CollapsedPaths) Collapse(path string) {
	self.collapsedPaths.Add(path)
}

func (self *CollapsedPaths) ToggleCollapsed(path string) {
	if self.collapsedPaths.Includes(path) {
		self.collapsedPaths.Remove(path)
	} else {
		self.collapsedPaths.Add(path)
	}
}
