package models

// CommitFile : A git commit file
type CommitFile struct {
	Path string

	ChangeStatus string // e.g. 'A' for added or 'M' for modified. This is based on the result from git diff --name-status
}

func (f *CommitFile) ID() string {
	return f.Path
}

func (f *CommitFile) Description() string {
	return f.Path
}

func (f *CommitFile) Added() bool {
	return f.ChangeStatus == "A"
}

func (f *CommitFile) Deleted() bool {
	return f.ChangeStatus == "D"
}

func (f *CommitFile) GetPath() string {
	return f.Path
}
