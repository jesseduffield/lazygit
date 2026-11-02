package config

import (
	"errors"
	"math"
	"os"

	"golang.org/x/sys/windows"
)

// lockFile blocks until it has locked file. The lock is released by unlockFile,
// or when the file is closed or the process exits.
func lockFile(file *os.File, kind lockKind) error {
	var flags uint32
	if kind == exclusiveLock {
		flags = windows.LOCKFILE_EXCLUSIVE_LOCK
	}
	return windows.LockFileEx(windows.Handle(file.Fd()), flags, 0, math.MaxUint32, math.MaxUint32, new(windows.Overlapped))
}

func unlockFile(file *os.File) error {
	return windows.UnlockFileEx(windows.Handle(file.Fd()), 0, math.MaxUint32, math.MaxUint32, new(windows.Overlapped))
}

func isLockingUnsupported(err error) bool {
	return errors.Is(err, windows.ERROR_NOT_SUPPORTED)
}
