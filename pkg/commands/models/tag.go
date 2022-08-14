package models

// Tag : A git tag
type Tag struct {
	Name string
}

func (t *Tag) FullRefName() string {
	return "refs/tags/" + t.RefName()
}

func (t *Tag) RefName() string {
	return t.Name
}

func (t *Tag) ParentRefName() string {
	return t.RefName() + "^"
}

func (t *Tag) ID() string {
	return t.RefName()
}

func (t *Tag) Description() string {
	return "tag " + t.Name
}
