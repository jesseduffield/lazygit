package commands

// Commit : A git commit
type Commit struct {
	Sha       string
	Name      string
	Status    string // one of "unpushed", "pushed", "merged", "rebasing" or "selected"
	Action    string // one of "", "pick", "edit", "squash", "reword", "drop", "fixup"
	Tags      []string
	ExtraInfo string // something like 'HEAD -> master, tag: v0.15.2'
	Author    string
	Date      string
}

func (c *Commit) ShortSha() string {
	if len(c.Sha) < 8 {
		return c.Sha
	}
	return c.Sha[:8]
}
