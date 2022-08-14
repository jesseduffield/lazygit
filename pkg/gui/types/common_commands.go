package types

type CheckoutRefOptions struct {
	WaitingStatus string
	EnvVars       []string
	OnRefNotFound func(ref string) error
}
