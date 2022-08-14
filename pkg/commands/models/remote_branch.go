package models

// Remote Branch : A git remote branch
type RemoteBranch struct {
	Name       string
	RemoteName string
}

func (r *RemoteBranch) FullName() string {
	return r.RemoteName + "/" + r.Name
}

func (r *RemoteBranch) FullRefName() string {
	return "refs/remotes/" + r.FullName()
}

func (r *RemoteBranch) RefName() string {
	return r.FullName()
}

func (r *RemoteBranch) ParentRefName() string {
	return r.RefName() + "^"
}

func (r *RemoteBranch) ID() string {
	return r.RefName()
}

func (r *RemoteBranch) Description() string {
	return r.RefName()
}
