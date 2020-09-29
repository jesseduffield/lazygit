package utils

import "github.com/go-errors/errors"

// WrapError wraps an error for the sake of showing a stack trace at the top level
// the go-errors package, for some reason, does not return nil when you try to wrap
// a non-error, so we're just doing it here
func WrapError(err error) error {
	if err == nil {
		return err
	}

	return errors.Wrap(err, 0)
}
