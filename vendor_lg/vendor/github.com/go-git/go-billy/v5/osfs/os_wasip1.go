//go:build wasip1
// +build wasip1

package osfs

import (
	"os"
	"syscall"
)

func (f *file) Lock() error {
	f.m.Lock()
	defer f.m.Unlock()
	return nil
}

func (f *file) Unlock() error {
	f.m.Lock()
	defer f.m.Unlock()
	return nil
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
