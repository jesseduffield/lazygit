package models

// Branch : A git branch
// duplicating this for now
type Branch struct {
	Name string
	// the displayname is something like '(HEAD detached at 123asdf)', whereas in that case the name would be '123asdf'
	DisplayName string
	// indicator of when the branch was last checked out e.g. '2d', '3m'
	Recency string
	// how many commits ahead we are from the remote branch (how many commits we can push)
	Pushables string
	// how many commits behind we are from the remote branch (how many commits we can pull)
	Pullables string
	// whether the remote branch is 'gone' i.e. we're tracking a remote branch that has been deleted
	UpstreamGone bool
	// whether this is the current branch. Exactly one branch should have this be true
	Head         bool
	DetachedHead bool
	// if we have a named remote locally this will be the name of that remote e.g.
	// 'origin' or 'tiwood'. If we don't have the remote locally it'll look like
	// 'git@github.com:tiwood/lazygit.git'
	UpstreamRemote string
	UpstreamBranch string
	// subject line in commit message
	Subject string
	// commit hash
	CommitHash string
}

func (b *Branch) FullRefName() string {
	if b.DetachedHead {
		return b.Name
	}
	return "refs/heads/" + b.Name
}

func (b *Branch) RefName() string {
	return b.Name
}

func (b *Branch) ParentRefName() string {
	return b.RefName() + "^"
}

func (b *Branch) ID() string {
	return b.RefName()
}

func (b *Branch) Description() string {
	return b.RefName()
}

func (b *Branch) IsTrackingRemote() bool {
	return b.UpstreamRemote != ""
}

// we know that the remote branch is not stored locally based on our pushable/pullable
// count being question marks.
func (b *Branch) RemoteBranchStoredLocally() bool {
	return b.IsTrackingRemote() && b.Pushables != "?" && b.Pullables != "?"
}

func (b *Branch) RemoteBranchNotStoredLocally() bool {
	return b.IsTrackingRemote() && b.Pushables == "?" && b.Pullables == "?"
}

func (b *Branch) MatchesUpstream() bool {
	return b.RemoteBranchStoredLocally() && b.Pushables == "0" && b.Pullables == "0"
}

func (b *Branch) HasCommitsToPush() bool {
	return b.RemoteBranchStoredLocally() && b.Pushables != "0"
}

func (b *Branch) HasCommitsToPull() bool {
	return b.RemoteBranchStoredLocally() && b.Pullables != "0"
}

// for when we're in a detached head state
func (b *Branch) IsRealBranch() bool {
	return b.Pushables != "" && b.Pullables != ""
}
