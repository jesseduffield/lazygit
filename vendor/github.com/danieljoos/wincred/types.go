package wincred

import (
	"time"
)

// CredentialPersistence describes one of three persistence modes of a credential.
// A detailed description of the available modes can be found on
// Docs: https://docs.microsoft.com/en-us/windows/desktop/api/wincred/ns-wincred-_credentialw
type CredentialPersistence uint32

const (
	// PersistSession indicates that the credential only persists for the life
	// of the current Windows login session. Such a credential is not visible in
	// any other logon session, even from the same user.
	PersistSession CredentialPersistence = 0x1

	// PersistLocalMachine indicates that the credential persists for this and
	// all subsequent logon sessions on this local machine/computer. It is
	// however not visible for logon sessions of this user on a different
	// machine.
	PersistLocalMachine CredentialPersistence = 0x2

	// PersistEnterprise indicates that the credential persists for this and all
	// subsequent logon sessions for this user. It is also visible for logon
	// sessions on different computers.
	PersistEnterprise CredentialPersistence = 0x3
)

// CredentialAttribute represents an application-specific attribute of a credential.
type CredentialAttribute struct {
	Keyword string
	Value   []byte
}

// Credential is the basic credential structure.
// A credential is identified by its target name.
// The actual credential secret is available in the CredentialBlob field.
type Credential struct {
	TargetName     string
	Comment        string
	LastWritten    time.Time
	CredentialBlob []byte
	Attributes     []CredentialAttribute
	TargetAlias    string
	UserName       string
	Persist        CredentialPersistence
}

// GenericCredential holds a credential for generic usage.
// It is typically defined and used by applications that need to manage user
// secrets.
//
// More information about the available kinds of credentials of the Windows
// Credential Management API can be found on Docs:
// https://docs.microsoft.com/en-us/windows/desktop/SecAuthN/kinds-of-credentials
type GenericCredential struct {
	Credential
}

// DomainPassword holds a domain credential that is typically used by the
// operating system for user logon.
//
// More information about the available kinds of credentials of the Windows
// Credential Management API can be found on Docs:
// https://docs.microsoft.com/en-us/windows/desktop/SecAuthN/kinds-of-credentials
type DomainPassword struct {
	Credential
}
