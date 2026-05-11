package models

import (
	"fmt"
	"sync/atomic"
)

// Branch : A git branch
// duplicating this for now
type Branch struct {
	Name string
	// the displayname is something like '(HEAD detached at 123asdf)', whereas in that case the name would be '123asdf'
	DisplayName string
	// indicator of when the branch was last checked out e.g. '2d', '3m'
	Recency string
	// how many commits ahead we are from the remote branch (how many commits we can push, assuming we push to our tracked remote branch)
	AheadForPull string
	// how many commits behind we are from the remote branch (how many commits we can pull)
	BehindForPull string
	// how many commits ahead we are from the branch we're pushing to (which might not be the same as our upstream branch in a triangular workflow)
	AheadForPush string
	// how many commits behind we are from the branch we're pushing to (which might not be the same as our upstream branch in a triangular workflow)
	BehindForPush string
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

	// How far we have fallen behind our base branch. 0 means either not
	// determined yet, or up to date with base branch. (We don't need to
	// distinguish the two, as we don't draw anything in both cases.)
	BehindBaseBranch atomic.Int32
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

func (b *Branch) ShortRefName() string {
	return b.RefName()
}

func (b *Branch) ParentRefName() string {
	return b.RefName() + "^"
}

func (b *Branch) FullUpstreamRefName() string {
	if b.UpstreamRemote == "" || b.UpstreamBranch == "" {
		return ""
	}

	return fmt.Sprintf("refs/remotes/%s/%s", b.UpstreamRemote, b.UpstreamBranch)
}

func (b *Branch) ShortUpstreamRefName() string {
	if b.UpstreamRemote == "" || b.UpstreamBranch == "" {
		return ""
	}

	return fmt.Sprintf("%s/%s", b.UpstreamRemote, b.UpstreamBranch)
}

func (b *Branch) ID() string {
	return b.RefName()
}

func (b *Branch) URN() string {
	return "branch-" + b.ID()
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
	return b.IsTrackingRemote() && b.AheadForPull != "?" && b.BehindForPull != "?"
}

func (b *Branch) RemoteBranchNotStoredLocally() bool {
	return b.IsTrackingRemote() && b.AheadForPull == "?" && b.BehindForPull == "?"
}

func (b *Branch) MatchesUpstream() bool {
	return b.RemoteBranchStoredLocally() && b.AheadForPull == "0" && b.BehindForPull == "0"
}

func (b *Branch) IsAheadForPull() bool {
	return b.RemoteBranchStoredLocally() && b.AheadForPull != "0"
}

func (b *Branch) IsBehindForPull() bool {
	return b.RemoteBranchStoredLocally() && b.BehindForPull != "0"
}

func (b *Branch) IsBehindForPush() bool {
	return b.RemoteBranchStoredLocally() && b.BehindForPush != "0"
}

// for when we're in a detached head state
func (b *Branch) IsRealBranch() bool {
	return b.AheadForPull != "" && b.BehindForPull != ""
}
