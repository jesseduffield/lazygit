// +build !windows

package osfs

import (
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
