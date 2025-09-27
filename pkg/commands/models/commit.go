package models

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
	"github.com/stefanhaller/git-todo-parser/todo"
)

// Special commit hash for empty tree object
const EmptyTreeCommitHash = "4b825dc642cb6eb9a060e54bf8d69288fbee4904"

type CommitStatus uint8

const (
	StatusNone CommitStatus = iota
	StatusUnpushed
	StatusPushed
	StatusMerged
	StatusRebasing
	StatusCherryPickingOrReverting
	StatusConflicted
	StatusReflog
)

const (
	// Conveniently for us, the todo package starts the enum at 1, and given
	// that it doesn't have a "none" value, we're setting ours to 0
	ActionNone todo.TodoCommand = 0
)

type Divergence uint8

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
	hash          *string
	Name          string
	Tags          []string
	ExtraInfo     string // something like 'HEAD -> master, tag: v0.15.2'
	AuthorName    string // something like 'Jesse Duffield'
	AuthorEmail   string // something like 'jessedduffield@gmail.com'
	UnixTimestamp int64

	// Hashes of parent commits (will be multiple if it's a merge commit)
	parents []*string

	// When filtering by path, this contains the paths that were changed in this
	// commit; nil when not filtering by path.
	FilterPaths []string

	Status     CommitStatus
	Action     todo.TodoCommand
	Divergence Divergence // set to DivergenceNone unless we are showing the divergence view
}

type NewCommitOpts struct {
	Hash          string
	Name          string
	Status        CommitStatus
	Action        todo.TodoCommand
	Tags          []string
	ExtraInfo     string
	AuthorName    string
	AuthorEmail   string
	UnixTimestamp int64
	Divergence    Divergence
	Parents       []string
}

func NewCommit(hashPool *utils.StringPool, opts NewCommitOpts) *Commit {
	return &Commit{
		hash:          hashPool.Add(opts.Hash),
		Name:          opts.Name,
		Status:        opts.Status,
		Action:        opts.Action,
		Tags:          opts.Tags,
		ExtraInfo:     opts.ExtraInfo,
		AuthorName:    opts.AuthorName,
		AuthorEmail:   opts.AuthorEmail,
		UnixTimestamp: opts.UnixTimestamp,
		Divergence:    opts.Divergence,
		parents:       lo.Map(opts.Parents, func(s string, _ int) *string { return hashPool.Add(s) }),
	}
}

func (c *Commit) Hash() string {
	return *c.hash
}

func (c *Commit) HashPtr() *string {
	return c.hash
}

func (c *Commit) ShortHash() string {
	return utils.ShortHash(c.Hash())
}

func (c *Commit) FullRefName() string {
	return c.Hash()
}

func (c *Commit) RefName() string {
	return c.Hash()
}

func (c *Commit) ShortRefName() string {
	return c.Hash()[:7]
}

func (c *Commit) ParentRefName() string {
	if c.IsFirstCommit() {
		return EmptyTreeCommitHash
	}
	return c.RefName() + "^"
}

func (c *Commit) Parents() []string {
	return lo.Map(c.parents, func(s *string, _ int) string { return *s })
}

func (c *Commit) ParentPtrs() []*string {
	return c.parents
}

func (c *Commit) IsFirstCommit() bool {
	return len(c.parents) == 0
}

func (c *Commit) ID() string {
	return c.RefName()
}

func (c *Commit) Description() string {
	return fmt.Sprintf("%s %s", c.Hash()[:7], c.Name)
}

func (c *Commit) IsMerge() bool {
	return len(c.parents) > 1
}

// returns true if this commit is not actually in the git log but instead
// is from a TODO file for an interactive rebase.
func (c *Commit) IsTODO() bool {
	return c.Action != ActionNone
}

func IsHeadCommit(commits []*Commit, index int) bool {
	return !commits[index].IsTODO() && (index == 0 || commits[index-1].IsTODO())
}
