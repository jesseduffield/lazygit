package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/stefanhaller/git-todo-parser/todo"
)

// We create shims for all the model classes in order to get a more stable API
// for custom commands. At the moment these are almost identical to the model
// classes, but this allows us to add "private" fields to the model classes that
// we don't want to expose to custom commands, or rename a model field to a
// better name without breaking people's custom commands. In such a case we add
// the new, better name to the shim but keep the old one for backwards
// compatibility. We already did this for Commit.Sha, which was renamed to Hash.

type Commit struct {
	Hash          string
	Sha           string // deprecated: use Hash
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
	ShortStatus             string
	IsWorktree              bool
}

type Branch struct {
	Name           string
	DisplayName    string
	Recency        string
	Pushables      string // deprecated: use AheadForPull
	Pullables      string // deprecated: use BehindForPull
	AheadForPull   string
	BehindForPull  string
	AheadForPush   string
	BehindForPush  string
	UpstreamGone   bool
	Head           bool
	DetachedHead   bool
	UpstreamRemote string
	UpstreamBranch string
	Subject        string
	CommitHash     string
}

type RemoteBranch struct {
	Name       string
	RemoteName string
}

type Remote struct {
	Name     string
	Urls     []string
	Branches []*RemoteBranch
}

type Tag struct {
	Name    string
	Message string
}

type StashEntry struct {
	Index   int
	Recency string
	Name    string
}

type CommitFile struct {
	Name         string
	ChangeStatus string
}

type Worktree struct {
	IsMain        bool
	IsCurrent     bool
	Path          string
	IsPathMissing bool
	GitDir        string
	Branch        string
	Name          string
}
