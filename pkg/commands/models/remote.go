package models

// Remote : A git remote
type Remote struct {
	Name     string
	Urls     []string
	Branches []*RemoteBranch
}

func (r *Remote) RefName() string {
	return r.Name
}

func (r *Remote) ID() string {
	return r.RefName()
}

func (r *Remote) URN() string {
	return "remote-" + r.ID()
}

func (r *Remote) Description() string {
	return r.RefName()
}
