package commands

// Tag : A git tag
type Tag struct {
	Name string
}

func (t *Tag) RefName() string {
	return t.Name
}
