package deadlock

import "sync"

// MutexImpl defines the interface for mutex implementations
type MutexImpl interface {
	Lock()
	Unlock()
	TryLock() bool
}

// RWMutexImpl defines the interface for rwmutex implementations
type RWMutexImpl interface {
	Lock()
	Unlock()
	RLock()
	RUnlock()
	TryLock() bool
	TryRLock() bool
	RLocker() sync.Locker
}
