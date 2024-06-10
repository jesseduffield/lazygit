package enums

type WorkingTreeState int

const (
	// this means we're neither rebasing nor merging
	WORKING_TREE_STATE_NONE WorkingTreeState = iota
	WORKING_TREE_STATE_REBASING
	WORKING_TREE_STATE_MERGING
)

func (self WorkingTreeState) IsMerging() bool {
	return self == WORKING_TREE_STATE_MERGING
}

func (self WorkingTreeState) IsRebasing() bool {
	return self == WORKING_TREE_STATE_REBASING
}
