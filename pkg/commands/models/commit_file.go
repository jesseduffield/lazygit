package models

// CommitFile : A git commit file
type CommitFile struct {
	Path string

	// For a renamed file, the path it was renamed from; empty otherwise.
	PreviousPath string

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

func (f *CommitFile) IsRename() bool {
	return f.PreviousPath != ""
}

// Names returns an array containing just the path, or in the case of a rename,
// the after path and the before path.
func (f *CommitFile) Names() []string {
	result := []string{f.Path}
	if f.PreviousPath != "" {
		result = append(result, f.PreviousPath)
	}
	return result
}

func (f *CommitFile) GetPath() string {
	return f.Path
}

func (f *CommitFile) GetPreviousPath() string {
	return f.PreviousPath
}
