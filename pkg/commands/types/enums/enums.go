package enums

type RebaseMode int

const (
	// this means we're neither rebasing nor merging
	REBASE_MODE_NONE RebaseMode = iota
	REBASE_MODE_REBASING
	REBASE_MODE_MERGING
)

func (self RebaseMode) IsMerging() bool {
	return self == REBASE_MODE_MERGING
}

func (self RebaseMode) IsRebasing() bool {
	return self == REBASE_MODE_REBASING
}
