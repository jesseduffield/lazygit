package models

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/stefanhaller/git-todo-parser/todo"
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
	// "Comment" is the last one of the todo package's enum entries
	ActionConflict = todo.Comment + 1
)

type Divergence int

// For a divergence log (left/right comparison of two refs) this is set to
// either DivergenceLeft or DivergenceRight for each commit; for normal
// commit views it is always DivergenceNone.
const (
	DivergenceNone Divergence = iota
	DivergenceLeft
	DivergenceRight
)

// Commit : A git commit
type Commit struct {
	Hash          string
	Name          string
	Status        CommitStatus
	Action        todo.TodoCommand
	Tags          []string
	ExtraInfo     string // something like 'HEAD -> master, tag: v0.15.2'
	AuthorName    string // something like 'Jesse Duffield'
	AuthorEmail   string // something like 'jessedduffield@gmail.com'
	UnixTimestamp int64
	Divergence    Divergence // set to DivergenceNone unless we are showing the divergence view

	// Hashes of parent commits (will be multiple if it's a merge commit)
	Parents []string
}

func (c *Commit) ShortHash() string {
	return utils.ShortHash(c.Hash)
}

func (c *Commit) FullRefName() string {
	return c.Hash
}

func (c *Commit) RefName() string {
	return c.Hash
}

func (c *Commit) ShortRefName() string {
	return c.Hash[:7]
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
	return fmt.Sprintf("%s %s", c.Hash[:7], c.Name)
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
