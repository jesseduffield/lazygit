//go:build wasm || js
// +build wasm js

package tty

import "errors"

// NewDevTty obtains a default tty from the console or TTY (e.g. /dev/tty) for the process.
func NewDevTty() (Tty, error) {
	return nil, errors.New("No tty device on wasm")
}

// NewDevTtyFromDev obtains a tty from the given device path. Not supported on Windows.
func NewDevTtyFromDev(dev string) (Tty, error) {
	return nil, errors.New("No tty device on wasm")
}

// NewStdIoTty obtains a tty from stdin and stdout.
func NewStdIoTty() (Tty, error) {
	return nil, errors.New("No tty device on wasm")
}
