// +build !windows

package osfs

import (
	"syscall"
)

func (f *file) Lock() error {
	f.m.Lock()
	defer f.m.Unlock()

	return syscall.Flock(int(f.File.Fd()), syscall.LOCK_EX)
}

func (f *file) Unlock() error {
	f.m.Lock()
	defer f.m.Unlock()

	return syscall.Flock(int(f.File.Fd()), syscall.LOCK_UN)
}
