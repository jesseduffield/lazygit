package models

// Branch : A git branch
// duplicating this for now
type Branch struct {
	Name string
	// the displayname is something like '(HEAD detached at 123asdf)', whereas in that case the name would be '123asdf'
	DisplayName  string
	Recency      string
	Pushables    string
	Pullables    string
	UpstreamName string
	Head         bool
	Merged       bool
}

func (b *Branch) RefName() string {
	return b.Name
}

func (b *Branch) ID() string {
	return b.RefName()
}

func (b *Branch) Description() string {
	return b.RefName()
}
