package models

// CommitFile : A git commit file
// TODO: this should really be renamed to 'DiffFile' because it need not be a file contained within a commit: it could be from any diff
type CommitFile struct {
	// TODO: rename this to Path
	Name string

	ChangeStatus string // e.g. 'A' for added or 'M' for modified. This is based on the result from git diff --name-status
}

func (f *CommitFile) ID() string {
	return f.Name
}

func (f *CommitFile) Description() string {
	return f.Name
}
