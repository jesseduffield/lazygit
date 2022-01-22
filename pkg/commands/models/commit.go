package models

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/utils"
)

// Commit : A git commit
type Commit struct {
	Sha           string
	Name          string
	Status        string // one of "unpushed", "pushed", "merged", "rebasing" or "selected"
	Action        string // one of "", "pick", "edit", "squash", "reword", "drop", "fixup"
	Tags          []string
	ExtraInfo     string // something like 'HEAD -> master, tag: v0.15.2'
	Author        string
	UnixTimestamp int64

	// SHAs of parent commits (will be multiple if it's a merge commit)
	Parents []string
}

func (c *Commit) ShortSha() string {
	return utils.ShortSha(c.Sha)
}

func (c *Commit) RefName() string {
	return c.Sha
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
	return c.Action != ""
}
