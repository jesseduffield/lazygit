package commands

import "fmt"

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
}

func (c *Commit) ShortSha() string {
	if len(c.Sha) < 8 {
		return c.Sha
	}
	return c.Sha[:8]
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
