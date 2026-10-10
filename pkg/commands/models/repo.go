package models

// Repo is one git repo in the multi-repo list.
type Repo struct {
	// Absolute path of the repo's working tree
	Path string
	// Path relative to the directory that lazygit was started in
	Name   string
	Branch string
	// True if the working tree has changes
	Dirty bool
}

func (r *Repo) ID() string {
	return r.Path
}
