package keyring

import "errors"

// provider set in the init function by the relevant os file e.g.:
// keyring_unix.go
var provider Keyring = fallbackServiceProvider{}

var (
	// ErrNotFound is the expected error if the secret isn't found in the
	// keyring.
	ErrNotFound = errors.New("secret not found in keyring")
	// ErrSetDataTooBig is returned if `Set` was called with too much data.
	// On MacOS: The combination of service, username & password should not exceed ~3000 bytes
	// On Windows: The service is limited to 32KiB while the password is limited to 2560 bytes
	// On Linux/Unix: There is no theoretical limit but performance suffers with big values (>100KiB)
	ErrSetDataTooBig = errors.New("data passed to Set was too big")
)

// Keyring provides a simple set/get interface for a keyring service.
type Keyring interface {
	// Set password in keyring for user.
	Set(service, user, password string) error
	// Get password from keyring given service and user name.
	Get(service, user string) (string, error)
	// Delete secret from keyring.
	Delete(service, user string) error
	// DeleteAll deletes all secrets for a given service
	DeleteAll(service string) error
}

// Set password in keyring for user.
func Set(service, user, password string) error {
	return provider.Set(service, user, password)
}

// Get password from keyring given service and user name.
func Get(service, user string) (string, error) {
	return provider.Get(service, user)
}

// Delete secret from keyring.
func Delete(service, user string) error {
	return provider.Delete(service, user)
}

// DeleteAll deletes all secrets for a given service
func DeleteAll(service string) error {
	return provider.DeleteAll(service)
}
