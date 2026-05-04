package models

// Remote : A git remote
type Remote struct {
	Name string
	Urls []string
	// PushUrls is empty unless the remote has explicit `remote.<name>.pushurl`
	// entries; when empty, pushes go to Urls.
	PushUrls []string
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
