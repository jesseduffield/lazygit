package commands

// File : A file from git status
// duplicating this for now
type File struct {
	Name                    string
	HasStagedChanges        bool
	HasUnstagedChanges      bool
	Tracked                 bool
	Deleted                 bool
	HasMergeConflicts       bool
	HasInlineMergeConflicts bool
	DisplayString           string
	Type                    string // one of 'file', 'directory', and 'other'
	ShortStatus             string // e.g. 'AD', ' A', 'M ', '??'
}
