package commands

// Commit : A git commit
type Commit struct {
	Sha           string
	Name          string
	Status        string // one of "unpushed", "pushed", "merged", "rebasing" or "selected"
	DisplayString string
	Action        string // one of "", "pick", "edit", "squash", "reword", "drop", "fixup"
	Copied        bool   // to know if this commit is ready to be cherry-picked somewhere
	Tags          []string
	ExtraInfo     string // something like 'HEAD -> master, tag: v0.15.2'
	Author        string
	Date          string
}
