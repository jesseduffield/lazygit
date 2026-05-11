package types

type CheckoutRefOptions struct {
	WaitingStatus string
	EnvVars       []string
	OnRefNotFound func(ref string) error

	// Refreshing pull requests is necessary when checking out a branch that doesn't exist locally
	// (e.g. checking out a remote branch), but it not needed when checking out an existing local
	// branch or a detached head (e.g. a tag).
	RefreshPullRequests bool
}
