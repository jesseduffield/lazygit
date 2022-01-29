package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type IRefsHelper interface {
	CheckoutRef(ref string, options types.CheckoutRefOptions) error
	CreateGitResetMenu(ref string) error
	ResetToRef(ref string, strength string, envVars []string) error
}

type ISuggestionsHelper interface {
	GetRemoteSuggestionsFunc() func(string) []*types.Suggestion
	GetBranchNameSuggestionsFunc() func(string) []*types.Suggestion
	GetFilePathSuggestionsFunc() func(string) []*types.Suggestion
	GetRemoteBranchesSuggestionsFunc(separator string) func(string) []*types.Suggestion
	GetRefsSuggestionsFunc() func(string) []*types.Suggestion
	GetCustomCommandsHistorySuggestionsFunc() func(string) []*types.Suggestion
}

type IFileHelper interface {
	EditFile(filename string) error
	EditFileAtLine(filename string, lineNumber int) error
	OpenFile(filename string) error
}

type IWorkingTreeHelper interface {
	AnyStagedFiles() bool
	AnyTrackedFiles() bool
	IsWorkingTreeDirty() bool
	FileForSubmodule(submodule *models.SubmoduleConfig) *models.File
}

// all fields mandatory (except `CanRebase` because it's boolean)
type SwitchToCommitFilesContextOpts struct {
	RefName    string
	CanRebase  bool
	Context    types.Context
	WindowName string
}
