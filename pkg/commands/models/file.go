package models

import (
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// File : A file from git status
// duplicating this for now
type File struct {
	Name                    string
	PreviousName            string
	HasStagedChanges        bool
	HasUnstagedChanges      bool
	Tracked                 bool
	Added                   bool
	Deleted                 bool
	HasMergeConflicts       bool
	HasInlineMergeConflicts bool
	DisplayString           string
	Type                    string // one of 'file', 'directory', and 'other'
	ShortStatus             string // e.g. 'AD', ' A', 'M ', '??'
}

// sometimes we need to deal with either a node (which contains a file) or an actual file
type IFile interface {
	GetHasUnstagedChanges() bool
	GetHasStagedChanges() bool
	GetIsTracked() bool
	GetPath() string
	GetPreviousPath() string
}

func (f *File) IsRename() bool {
	return f.PreviousName != ""
}

// Names returns an array containing just the filename, or in the case of a rename, the after filename and the before filename
func (f *File) Names() []string {
	result := []string{f.Name}
	if f.PreviousName != "" {
		result = append(result, f.PreviousName)
	}
	return result
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

func (f *File) IsSubmodule(configs []*SubmoduleConfig) bool {
	return f.SubmoduleConfig(configs) != nil
}

func (f *File) SubmoduleConfig(configs []*SubmoduleConfig) *SubmoduleConfig {
	for _, config := range configs {
		if f.Name == config.Path {
			return config
		}
	}

	return nil
}

func (f *File) GetHasUnstagedChanges() bool {
	return f.HasUnstagedChanges
}

func (f *File) GetHasStagedChanges() bool {
	return f.HasStagedChanges
}

func (f *File) GetIsTracked() bool {
	return f.Tracked
}

func (f *File) GetPath() string {
	// TODO: remove concept of name; just use path
	return f.Name
}

func (f *File) GetPreviousPath() string {
	return f.PreviousName
}
