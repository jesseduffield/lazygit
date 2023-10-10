package models

// Tag : A git tag
type Tag struct {
	Name string
	// this is either the first line of the message of an annotated tag, or the
	// first line of a commit message for a lightweight tag
	Message string
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

func (t *Tag) URN() string {
	return "tag-" + t.ID()
}

func (t *Tag) Description() string {
	return t.Message
}
