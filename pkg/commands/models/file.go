package models

import (
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

// File : A file from git status
// duplicating this for now
type File struct {
	Path                    string
	PreviousPath            string
	HasStagedChanges        bool
	HasUnstagedChanges      bool
	Tracked                 bool
	Added                   bool
	Deleted                 bool
	HasMergeConflicts       bool
	HasInlineMergeConflicts bool
	DisplayString           string
	ShortStatus             string // e.g. 'AD', ' A', 'M ', '??'
	LinesDeleted            int
	LinesAdded              int

	// If true, this must be a worktree folder
	IsWorktree bool
}

// sometimes we need to deal with either a node (which contains a file) or an actual file
type IFile interface {
	GetHasUnstagedChanges() bool
	GetHasStagedChanges() bool
	GetIsTracked() bool
	GetPath() string
	GetPreviousPath() string
	GetIsFile() bool
}

func (f *File) IsRename() bool {
	return f.PreviousPath != ""
}

// Names returns an array containing just the filename, or in the case of a rename, the after filename and the before filename
func (f *File) Names() []string {
	result := []string{f.Path}
	if f.PreviousPath != "" {
		result = append(result, f.PreviousPath)
	}
	return result
}

// returns true if the file names are the same or if a file rename includes the filename of the other
func (f *File) Matches(f2 *File) bool {
	return utils.StringArraysOverlap(f.Names(), f2.Names())
}

func (f *File) ID() string {
	return f.Path
}

func (f *File) Description() string {
	return f.Path
}

func (f *File) IsSubmodule(configs []*SubmoduleConfig) bool {
	return f.SubmoduleConfig(configs) != nil
}

func (f *File) SubmoduleConfig(configs []*SubmoduleConfig) *SubmoduleConfig {
	for _, config := range configs {
		if f.Path == config.Path {
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
	return f.Path
}

func (f *File) GetPreviousPath() string {
	return f.PreviousPath
}

func (f *File) GetIsFile() bool {
	return true
}

func (f *File) GetMergeStateDescription(tr *i18n.TranslationSet) string {
	m := map[string]string{
		"DD": tr.MergeConflictDescription_DD,
		"AU": tr.MergeConflictDescription_AU,
		"UA": tr.MergeConflictDescription_UA,
		"DU": tr.MergeConflictDescription_DU,
		"UD": tr.MergeConflictDescription_UD,
	}

	if description, ok := m[f.ShortStatus]; ok {
		return description
	}

	panic("should only be called if there's a merge conflict")
}

type StatusFields struct {
	HasStagedChanges        bool
	HasUnstagedChanges      bool
	Tracked                 bool
	Deleted                 bool
	Added                   bool
	HasMergeConflicts       bool
	HasInlineMergeConflicts bool
	ShortStatus             string
}

func SetStatusFields(file *File, shortStatus string) {
	derived := deriveStatusFields(shortStatus)

	file.HasStagedChanges = derived.HasStagedChanges
	file.HasUnstagedChanges = derived.HasUnstagedChanges
	file.Tracked = derived.Tracked
	file.Deleted = derived.Deleted
	file.Added = derived.Added
	file.HasMergeConflicts = derived.HasMergeConflicts
	file.HasInlineMergeConflicts = derived.HasInlineMergeConflicts
	file.ShortStatus = derived.ShortStatus
}

// shortStatus is something like '??' or 'A '
func deriveStatusFields(shortStatus string) StatusFields {
	stagedChange := shortStatus[0:1]
	unstagedChange := shortStatus[1:2]
	tracked := !lo.Contains([]string{"??", "A ", "AM"}, shortStatus)
	hasStagedChanges := !lo.Contains([]string{" ", "U", "?"}, stagedChange)
	hasInlineMergeConflicts := lo.Contains([]string{"UU", "AA"}, shortStatus)
	hasMergeConflicts := hasInlineMergeConflicts || lo.Contains([]string{"DD", "AU", "UA", "UD", "DU"}, shortStatus)

	return StatusFields{
		HasStagedChanges:        hasStagedChanges,
		HasUnstagedChanges:      unstagedChange != " ",
		Tracked:                 tracked,
		Deleted:                 unstagedChange == "D" || stagedChange == "D",
		Added:                   unstagedChange == "A" || !tracked,
		HasMergeConflicts:       hasMergeConflicts,
		HasInlineMergeConflicts: hasInlineMergeConflicts,
		ShortStatus:             shortStatus,
	}
}
