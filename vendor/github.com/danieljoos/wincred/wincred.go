// Package wincred provides primitives for accessing the Windows Credentials Management API.
// This includes functions for retrieval, listing and storage of credentials as well as Go structures for convenient access to the credential data.
//
// A more detailed description of Windows Credentials Management can be found on
// Docs: https://docs.microsoft.com/en-us/windows/desktop/SecAuthN/credentials-management
package wincred

import "errors"

const (
	// ErrElementNotFound is the error that is returned if a requested element cannot be found.
	// This error constant can be used to check if a credential could not be found.
	ErrElementNotFound = sysERROR_NOT_FOUND

	// ErrInvalidParameter is the error that is returned for invalid parameters.
	// This error constant can be used to check if the given function parameters were invalid.
	// For example when trying to create a new generic credential with an empty target name.
	ErrInvalidParameter = sysERROR_INVALID_PARAMETER

	// ErrBadUsername is returned when the credential's username is invalid.
	ErrBadUsername = sysERROR_BAD_USERNAME
)

// GetGenericCredential fetches the generic credential with the given name from Windows credential manager.
// It returns nil and an error if the credential could not be found or an error occurred.
func GetGenericCredential(targetName string) (*GenericCredential, error) {
	cred, err := sysCredRead(targetName, sysCRED_TYPE_GENERIC)
	if cred != nil {
		return &GenericCredential{Credential: *cred}, err
	}
	return nil, err
}

// NewGenericCredential creates a new generic credential object with the given name.
// The persist mode of the newly created object is set to a default value that indicates local-machine-wide storage.
// The credential object is NOT yet persisted to the Windows credential vault.
func NewGenericCredential(targetName string) (result *GenericCredential) {
	result = new(GenericCredential)
	result.TargetName = targetName
	result.Persist = PersistLocalMachine
	return
}

// Write persists the generic credential object to Windows credential manager.
func (t *GenericCredential) Write() (err error) {
	err = sysCredWrite(&t.Credential, sysCRED_TYPE_GENERIC)
	return
}

// Delete removes the credential object from Windows credential manager.
func (t *GenericCredential) Delete() (err error) {
	err = sysCredDelete(&t.Credential, sysCRED_TYPE_GENERIC)
	return
}

// GetDomainPassword fetches the domain-password credential with the given target host name from Windows credential manager.
// It returns nil and an error if the credential could not be found or an error occurred.
func GetDomainPassword(targetName string) (*DomainPassword, error) {
	cred, err := sysCredRead(targetName, sysCRED_TYPE_DOMAIN_PASSWORD)
	if cred != nil {
		return &DomainPassword{Credential: *cred}, err
	}
	return nil, err
}

// NewDomainPassword creates a new domain-password credential used for login to the given target host name.
// The  persist mode of the newly created object is set to a default value that indicates local-machine-wide storage.
// The credential object is NOT yet persisted to the Windows credential vault.
func NewDomainPassword(targetName string) (result *DomainPassword) {
	result = new(DomainPassword)
	result.TargetName = targetName
	result.Persist = PersistLocalMachine
	return
}

// Write persists the domain-password credential to Windows credential manager.
func (t *DomainPassword) Write() (err error) {
	err = sysCredWrite(&t.Credential, sysCRED_TYPE_DOMAIN_PASSWORD)
	return
}

// Delete removes the domain-password credential from Windows credential manager.
func (t *DomainPassword) Delete() (err error) {
	err = sysCredDelete(&t.Credential, sysCRED_TYPE_DOMAIN_PASSWORD)
	return
}

// SetPassword sets the CredentialBlob field of a domain password credential to the given string.
func (t *DomainPassword) SetPassword(pw string) {
	t.CredentialBlob = utf16ToByte(utf16FromString(pw))
}

// List retrieves all credentials of the Credentials store.
func List() ([]*Credential, error) {
	creds, err := sysCredEnumerate("", true)
	if err != nil && errors.Is(err, ErrElementNotFound) {
		// Ignore ERROR_NOT_FOUND and return an empty list instead
		creds = []*Credential{}
		err = nil
	}
	return creds, err
}

// FilteredList retrieves the list of credentials from the Credentials store that match the given filter.
// The filter string defines the prefix followed by an asterisk for the `TargetName` attribute of the credentials.
func FilteredList(filter string) ([]*Credential, error) {
	creds, err := sysCredEnumerate(filter, false)
	if err != nil && errors.Is(err, ErrElementNotFound) {
		// Ignore ERROR_NOT_FOUND and return an empty list instead
		creds = []*Credential{}
		err = nil
	}
	return creds, err
}
