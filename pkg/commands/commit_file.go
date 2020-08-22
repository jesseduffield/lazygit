package commands

// CommitFile : A git commit file
type CommitFile struct {
	// Parent is the identifier of the parent object e.g. a commit SHA if this commit file is for a commit, or a stash entry ref like 'stash@{1}'
	Parent        string
	Name          string
	DisplayString string
	Status        int // one of 'WHOLE' 'PART' 'NONE'
}

func (f *CommitFile) ID() string {
	return f.Name
}

func (f *CommitFile) Description() string {
	return f.Name
}
