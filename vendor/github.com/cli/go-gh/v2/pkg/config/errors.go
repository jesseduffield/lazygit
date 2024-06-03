package config

import (
	"fmt"
)

// InvalidConfigFileError represents an error when trying to read a config file.
type InvalidConfigFileError struct {
	Path string
	Err  error
}

// Allow InvalidConfigFileError to satisfy error interface.
func (e *InvalidConfigFileError) Error() string {
	return fmt.Sprintf("invalid config file %s: %s", e.Path, e.Err)
}

// Allow InvalidConfigFileError to be unwrapped.
func (e *InvalidConfigFileError) Unwrap() error {
	return e.Err
}

// KeyNotFoundError represents an error when trying to find a config key
// that does not exist.
type KeyNotFoundError struct {
	Key string
}

// Allow KeyNotFoundError to satisfy error interface.
func (e *KeyNotFoundError) Error() string {
	return fmt.Sprintf("could not find key %q", e.Key)
}
