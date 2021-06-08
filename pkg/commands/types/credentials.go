package types

type CredentialKind int

const (
	USERNAME CredentialKind = iota
	PASSWORD
	PASSPHRASE
)
