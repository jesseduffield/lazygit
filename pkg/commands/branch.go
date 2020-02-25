package commands

// Branch : A git branch
// duplicating this for now
type Branch struct {
	Name      string
	Recency   string
	Pushables string
	Pullables string
}
