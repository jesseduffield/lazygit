package models

// Branch : A git branch
// duplicating this for now
type Branch struct {
	Name string
	// the displayname is something like '(HEAD detached at 123asdf)', whereas in that case the name would be '123asdf'
	DisplayName  string
	Recency      string
	Pushables    string
	Pullables    string
	UpstreamName string
	Head         bool
}

func (b *Branch) RefName() string {
	return b.Name
}

func (b *Branch) ID() string {
	return b.RefName()
}

func (b *Branch) Description() string {
	return b.RefName()
}

// this method does not consider the case where the git config states that a branch is tracking the config.
// The Pullables value here is based on whether or not we saw an upstream when doing `git branch`
func (b *Branch) IsTrackingRemote() bool {
	return b.IsRealBranch() && b.Pullables != "?"
}

func (b *Branch) MatchesUpstream() bool {
	return b.IsRealBranch() && b.Pushables == "0" && b.Pullables == "0"
}

func (b *Branch) HasCommitsToPush() bool {
	return b.IsRealBranch() && b.Pushables != "0"
}

func (b *Branch) HasCommitsToPull() bool {
	return b.IsRealBranch() && b.Pullables != "0"
}

// for when we're in a detached head state
func (b *Branch) IsRealBranch() bool {
	return b.Pushables != "" && b.Pullables != ""
}
