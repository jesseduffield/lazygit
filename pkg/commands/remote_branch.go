package commands

// Remote Branch : A git remote branch
type RemoteBranch struct {
	Name       string
	RemoteName string
}

func (r *RemoteBranch) FullName() string {
	return r.RemoteName + "/" + r.Name
}
