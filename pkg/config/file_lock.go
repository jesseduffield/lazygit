//go:build !windows

package config

import (
	"errors"
	"os"

	"golang.org/x/sys/unix"
)

// lockFile blocks until it has locked file. The lock is released by unlockFile,
// or when the file is closed or the process exits.
func lockFile(file *os.File, kind lockKind) error {
	how := unix.LOCK_SH
	if kind == exclusiveLock {
		how = unix.LOCK_EX
	}
	for {
		err := unix.Flock(int(file.Fd()), how)
		if !errors.Is(err, unix.EINTR) {
			return err
		}
	}
}

func unlockFile(file *os.File) error {
	return unix.Flock(int(file.Fd()), unix.LOCK_UN)
}

func isLockingUnsupported(err error) bool {
	return errors.Is(err, unix.ENOTSUP) || errors.Is(err, unix.EOPNOTSUPP) || errors.Is(err, unix.ENOLCK)
}
