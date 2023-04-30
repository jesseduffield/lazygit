package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
)

// loads the session state at the time that a custom command is invoked, for use
// in the custom command's template strings
type SessionStateLoader struct {
	c          *helpers.HelperCommon
	refsHelper *helpers.RefsHelper
}

func NewSessionStateLoader(c *helpers.HelperCommon, refsHelper *helpers.RefsHelper) *SessionStateLoader {
	return &SessionStateLoader{
		c:          c,
		refsHelper: refsHelper,
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
		SelectedFile:           self.c.Contexts().Files.GetSelectedFile(),
		SelectedPath:           self.c.Contexts().Files.GetSelectedPath(),
		SelectedLocalCommit:    self.c.Contexts().LocalCommits.GetSelected(),
		SelectedReflogCommit:   self.c.Contexts().ReflogCommits.GetSelected(),
		SelectedLocalBranch:    self.c.Contexts().Branches.GetSelected(),
		SelectedRemoteBranch:   self.c.Contexts().RemoteBranches.GetSelected(),
		SelectedRemote:         self.c.Contexts().Remotes.GetSelected(),
		SelectedTag:            self.c.Contexts().Tags.GetSelected(),
		SelectedStashEntry:     self.c.Contexts().Stash.GetSelected(),
		SelectedCommitFile:     self.c.Contexts().CommitFiles.GetSelectedFile(),
		SelectedCommitFilePath: self.c.Contexts().CommitFiles.GetSelectedPath(),
		SelectedSubCommit:      self.c.Contexts().SubCommits.GetSelected(),
		CheckedOutBranch:       self.refsHelper.GetCheckedOutRef(),
	}
}
