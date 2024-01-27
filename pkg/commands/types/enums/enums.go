package enums

type RebaseMode int

const (
	// this means we're neither rebasing nor merging
	REBASE_MODE_NONE RebaseMode = iota
	// this means normal rebase as opposed to interactive rebase
	REBASE_MODE_NORMAL
	REBASE_MODE_INTERACTIVE
	// REBASE_MODE_REBASING is a general state that captures both REBASE_MODE_NORMAL and REBASE_MODE_INTERACTIVE
	REBASE_MODE_REBASING
	REBASE_MODE_MERGING
)

func (self RebaseMode) IsMerging() bool {
	return self == REBASE_MODE_MERGING
}

func (self RebaseMode) IsRebasing() bool {
	return self == REBASE_MODE_INTERACTIVE || self == REBASE_MODE_NORMAL || self == REBASE_MODE_REBASING
}
