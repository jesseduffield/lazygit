package commands

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/utils"
)

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
	IsSubmodule             bool
}

const RENAME_SEPARATOR = " -> "

func (f *File) IsRename() bool {
	return strings.Contains(f.Name, RENAME_SEPARATOR)
}

// Names returns an array containing just the filename, or in the case of a rename, the after filename and the before filename
func (f *File) Names() []string {
	return strings.Split(f.Name, RENAME_SEPARATOR)
}

// returns true if the file names are the same or if a a file rename includes the filename of the other
func (f *File) Matches(f2 *File) bool {
	return utils.StringArraysOverlap(f.Names(), f2.Names())
}

func (f *File) ID() string {
	return f.Name
}

func (f *File) Description() string {
	return f.Name
}
