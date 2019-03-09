package commands

// CommitFile : A git commit file
type CommitFile struct {
	Sha           string
	Name          string
	DisplayString string
}

// GetDisplayStrings is a function.
func (f *CommitFile) GetDisplayStrings() []string {
	return []string{f.DisplayString}
}
