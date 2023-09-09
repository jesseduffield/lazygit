//go:build go1.13
// +build go1.13

package errors

import (
	baseErrors "errors"
)

// As finds the first error in err's tree that matches target, and if one is found, sets
// target to that error value and returns true. Otherwise, it returns false.
//
// For more information see stdlib errors.As.
func As(err error, target interface{}) bool {
	return baseErrors.As(err, target)
}

// Is detects whether the error is equal to a given error. Errors
// are considered equal by this function if they are matched by errors.Is
// or if their contained errors are matched through errors.Is.
func Is(e error, original error) bool {
	if baseErrors.Is(e, original) {
		return true
	}

	if e, ok := e.(*Error); ok {
		return Is(e.Err, original)
	}

	if original, ok := original.(*Error); ok {
		return Is(e, original.Err)
	}

	return false
}

// Join returns an error that wraps the given errors.
// Any nil error values are discarded.
// Join returns nil if every value in errs is nil.
// The error formats as the concatenation of the strings obtained
// by calling the Error method of each element of errs, with a newline
// between each string.
//
// A non-nil error returned by Join implements the Unwrap() []error method.
//
// For more information see stdlib errors.Join.
func Join(errs ...error) error {
	return baseErrors.Join(errs...)
}

// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
//
// Unwrap only calls a method of the form "Unwrap() error".
// In particular Unwrap does not unwrap errors returned by [Join].
//
// For more information see stdlib errors.Unwrap.
func Unwrap(err error) error {
	return baseErrors.Unwrap(err)
}
