package commands

// Remote : A git remote
type Remote struct {
	Name     string
	Urls     []string
	Branches []*RemoteBranch
}
