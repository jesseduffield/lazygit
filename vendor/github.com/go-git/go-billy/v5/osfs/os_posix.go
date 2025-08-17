//go:build !plan9 && !windows && !wasm
// +build !plan9,!windows,!wasm

package osfs

import (
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

func (f *file) Lock() error {
	f.m.Lock()
	defer f.m.Unlock()

	return unix.Flock(int(f.File.Fd()), unix.LOCK_EX)
}

func (f *file) Unlock() error {
	f.m.Lock()
	defer f.m.Unlock()

	return unix.Flock(int(f.File.Fd()), unix.LOCK_UN)
}

func rename(from, to string) error {
	return os.Rename(from, to)
}

// umask sets umask to a new value, and returns a func which allows the
// caller to reset it back to what it was originally.
func umask(new int) func() {
	old := syscall.Umask(new)
	return func() {
		syscall.Umask(old)
	}
}
