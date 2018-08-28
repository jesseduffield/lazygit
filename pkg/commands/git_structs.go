package commands

// File : A staged/unstaged file
type File struct {
	Name               string
	HasStagedChanges   bool
	HasUnstagedChanges bool
	Tracked            bool
	Deleted            bool
	HasMergeConflicts  bool
	DisplayString      string
	Type               string // one of 'file', 'directory', and 'other'
}

// Commit : A git commit
type Commit struct {
	Sha           string
	Name          string
	Pushed        bool
	DisplayString string
}

// StashEntry : A git stash entry
type StashEntry struct {
	Index         int
	Name          string
	DisplayString string
}

// Conflict : A git conflict with a start middle and end corresponding to line
// numbers in the file where the conflict bars appear
type Conflict struct {
	Start  int
	Middle int
	End    int
}
