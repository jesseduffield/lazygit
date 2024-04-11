package custom_commands

import (
	"github.com/fsmiamoto/git-todo-parser/todo"
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

type Commit struct {
	Hash          string
	Sha           string
	Name          string
	Status        models.CommitStatus
	Action        todo.TodoCommand
	Tags          []string
	ExtraInfo     string
	AuthorName    string
	AuthorEmail   string
	UnixTimestamp int64
	Divergence    models.Divergence
	Parents       []string
}

func commitWrapperFromModelCommit(commit *models.Commit) *Commit {
	if commit == nil {
		return nil
	}

	return &Commit{
		Hash:          commit.Hash,
		Sha:           commit.Hash,
		Name:          commit.Name,
		Status:        commit.Status,
		Action:        commit.Action,
		Tags:          commit.Tags,
		ExtraInfo:     commit.ExtraInfo,
		AuthorName:    commit.AuthorName,
		AuthorEmail:   commit.AuthorEmail,
		UnixTimestamp: commit.UnixTimestamp,
		Divergence:    commit.Divergence,
		Parents:       commit.Parents,
	}
}

// SessionState captures the current state of the application for use in custom commands
type SessionState struct {
	SelectedLocalCommit    *Commit
	SelectedReflogCommit   *Commit
	SelectedSubCommit      *Commit
	SelectedFile           *models.File
	SelectedPath           string
	SelectedLocalBranch    *models.Branch
	SelectedRemoteBranch   *models.RemoteBranch
	SelectedRemote         *models.Remote
	SelectedTag            *models.Tag
	SelectedStashEntry     *models.StashEntry
	SelectedCommitFile     *models.CommitFile
	SelectedCommitFilePath string
	SelectedWorktree       *models.Worktree
	CheckedOutBranch       *models.Branch
}

func (self *SessionStateLoader) call() *SessionState {
	return &SessionState{
		SelectedFile:           self.c.Contexts().Files.GetSelectedFile(),
		SelectedPath:           self.c.Contexts().Files.GetSelectedPath(),
		SelectedLocalCommit:    commitWrapperFromModelCommit(self.c.Contexts().LocalCommits.GetSelected()),
		SelectedReflogCommit:   commitWrapperFromModelCommit(self.c.Contexts().ReflogCommits.GetSelected()),
		SelectedLocalBranch:    self.c.Contexts().Branches.GetSelected(),
		SelectedRemoteBranch:   self.c.Contexts().RemoteBranches.GetSelected(),
		SelectedRemote:         self.c.Contexts().Remotes.GetSelected(),
		SelectedTag:            self.c.Contexts().Tags.GetSelected(),
		SelectedStashEntry:     self.c.Contexts().Stash.GetSelected(),
		SelectedCommitFile:     self.c.Contexts().CommitFiles.GetSelectedFile(),
		SelectedCommitFilePath: self.c.Contexts().CommitFiles.GetSelectedPath(),
		SelectedSubCommit:      commitWrapperFromModelCommit(self.c.Contexts().SubCommits.GetSelected()),
		SelectedWorktree:       self.c.Contexts().Worktrees.GetSelected(),
		CheckedOutBranch:       self.refsHelper.GetCheckedOutRef(),
	}
}
