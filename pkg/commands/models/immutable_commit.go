package models

// This model contains the information that is intrinsic to a commit, meaning
// that we can depend on it not changing over time. Git commits are immutable,
// but our other Commit model has extra fields added for ease of use.
type ImmutableCommit struct {
	hash string

	// hashes of parent commits (will be multiple if it's a merge commit)
	parentHashes []string
}

func NewImmutableCommit(hash string, parentHashes []string) ImmutableCommit {
	return ImmutableCommit{
		hash:         hash,
		parentHashes: parentHashes,
	}
}

func (self *ImmutableCommit) Hash() string {
	return self.hash
}

func (self *ImmutableCommit) ParentHashes() []string {
	return self.parentHashes
}

func (self *ImmutableCommit) IsRoot() bool {
	return len(self.parentHashes) == 0
}
