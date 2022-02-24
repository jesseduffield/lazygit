package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
)

// loads the session state at the time that a custom command is invoked, for use
// in the custom command's template strings
type SessionStateLoader struct {
	contexts *context.ContextTree
	helpers  *helpers.Helpers
}

func NewSessionStateLoader(contexts *context.ContextTree, helpers *helpers.Helpers) *SessionStateLoader {
	return &SessionStateLoader{
		contexts: contexts,
		helpers:  helpers,
	}
}

// SessionState captures the current state of the application for use in custom commands
type SessionState struct {
	SelectedLocalCommit    *models.Commit
	SelectedReflogCommit   *models.Commit
	SelectedSubCommit      *models.Commit
	SelectedFile           *models.File
	SelectedPath           string
	SelectedLocalBranch    *models.Branch
	SelectedRemoteBranch   *models.RemoteBranch
	SelectedRemote         *models.Remote
	SelectedTag            *models.Tag
	SelectedStashEntry     *models.StashEntry
	SelectedCommitFile     *models.CommitFile
	SelectedCommitFilePath string
	CheckedOutBranch       *models.Branch
}

func (self *SessionStateLoader) call() *SessionState {
	return &SessionState{
		SelectedFile:           self.contexts.Files.GetSelectedFile(),
		SelectedPath:           self.contexts.Files.GetSelectedPath(),
		SelectedLocalCommit:    self.contexts.LocalCommits.GetSelected(),
		SelectedReflogCommit:   self.contexts.ReflogCommits.GetSelected(),
		SelectedLocalBranch:    self.contexts.Branches.GetSelected(),
		SelectedRemoteBranch:   self.contexts.RemoteBranches.GetSelected(),
		SelectedRemote:         self.contexts.Remotes.GetSelected(),
		SelectedTag:            self.contexts.Tags.GetSelected(),
		SelectedStashEntry:     self.contexts.Stash.GetSelected(),
		SelectedCommitFile:     self.contexts.CommitFiles.GetSelectedFile(),
		SelectedCommitFilePath: self.contexts.CommitFiles.GetSelectedPath(),
		SelectedSubCommit:      self.contexts.SubCommits.GetSelected(),
		CheckedOutBranch:       self.helpers.Refs.GetCheckedOutRef(),
	}
}
