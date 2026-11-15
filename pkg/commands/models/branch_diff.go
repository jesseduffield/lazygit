package models

// BranchDiff is a branch's changes since it diverged from its base branch, as a pull
// request of it would show them: the diff from the merge base of the two to the
// branch's tip. Changes the branch got by merging in the base branch aren't part of it.
type BranchDiff struct {
	// Branch is a local branch or a remote one.
	Branch Ref

	// BaseBranch is the short name of the base branch, for display.
	BaseBranch string

	// MergeBase is the hash of the commit the branch diverged from its base branch at.
	MergeBase string

	// Tip is the hash of the commit the branch was at when the diff was asked for. The
	// diff stays as it was then, rather than following the branch as it moves.
	Tip string
}

func (b *BranchDiff) FullRefName() string {
	return b.Branch.FullRefName()
}

func (b *BranchDiff) RefName() string {
	return b.Tip
}

func (b *BranchDiff) ShortRefName() string {
	return b.Branch.ShortRefName()
}

// ParentRefName is the merge base, so that the diff from a ref's parent to the ref,
// which is what a ref's files are listed from, is the diff of the branch.
func (b *BranchDiff) ParentRefName() string {
	return b.MergeBase
}

func (b *BranchDiff) Description() string {
	return b.Branch.Description() + " since " + b.BaseBranch
}
