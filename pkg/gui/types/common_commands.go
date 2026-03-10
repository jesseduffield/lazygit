package types

type CheckoutRefOptions struct {
	WaitingStatus string
	EnvVars       []string
	OnRefNotFound func(ref string) error
	// If set, this function is called right before the checkout command.
	// Used e.g. to detach a worktree before checking out the branch.
	PreCheckoutCommand func() error
}
