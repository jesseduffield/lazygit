package models

import (
	"fmt"

	"github.com/fsmiamoto/git-todo-parser/todo"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// Special commit hash for empty tree object
const EmptyTreeCommitHash = "4b825dc642cb6eb9a060e54bf8d69288fbee4904"

type CommitStatus int

const (
	StatusNone CommitStatus = iota
	StatusUnpushed
	StatusPushed
	StatusMerged
	StatusRebasing
	StatusSelected
	StatusReflog
)

const (
	// Conveniently for us, the todo package starts the enum at 1, and given
	// that it doesn't have a "none" value, we're setting ours to 0
	ActionNone todo.TodoCommand = 0
)

// Commit : A git commit
type Commit struct {
	Sha           string
	Name          string
	Status        CommitStatus
	Action        todo.TodoCommand
	Tags          []string
	ExtraInfo     string // something like 'HEAD -> master, tag: v0.15.2'
	AuthorName    string // something like 'Jesse Duffield'
	AuthorEmail   string // something like 'jessedduffield@gmail.com'
	UnixTimestamp int64

	// SHAs of parent commits (will be multiple if it's a merge commit)
	Parents []string
}

func (c *Commit) ShortSha() string {
	return utils.ShortSha(c.Sha)
}

func (c *Commit) FullRefName() string {
	return c.Sha
}

func (c *Commit) RefName() string {
	return c.Sha
}

func (c *Commit) ParentRefName() string {
	if c.IsFirstCommit() {
		return EmptyTreeCommitHash
	}
	return c.RefName() + "^"
}

func (c *Commit) IsFirstCommit() bool {
	return len(c.Parents) == 0
}

func (c *Commit) ID() string {
	return c.RefName()
}

func (c *Commit) Description() string {
	return fmt.Sprintf("%s %s", c.Sha[:7], c.Name)
}

func (c *Commit) IsMerge() bool {
	return len(c.Parents) > 1
}

// returns true if this commit is not actually in the git log but instead
// is from a TODO file for an interactive rebase.
func (c *Commit) IsTODO() bool {
	return c.Action != ActionNone
}

func IsHeadCommit(commits []*Commit, index int) bool {
	return !commits[index].IsTODO() && (index == 0 || commits[index-1].IsTODO())
}
