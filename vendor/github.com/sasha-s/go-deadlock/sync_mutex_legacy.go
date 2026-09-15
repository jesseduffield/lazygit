//go:build !goexperiment.synctest && !deadlock_synctest && !deadlock_disable && !go1.18
// +build !goexperiment.synctest,!deadlock_synctest,!deadlock_disable,!go1.18

package deadlock

import "sync"

// StandardMutex wraps sync.Mutex
type StandardMutex struct {
	mu sync.Mutex
}

func (m *StandardMutex) Lock() {
	m.mu.Lock()
}

func (m *StandardMutex) Unlock() {
	m.mu.Unlock()
}

func (m *StandardMutex) TryLock() bool {
	// TryLock is not available before Go 1.18
	panic("TryLock requires Go 1.18 or later")
}

// StandardRWMutex wraps sync.RWMutex
type StandardRWMutex struct {
	mu sync.RWMutex
}

func (m *StandardRWMutex) Lock() {
	m.mu.Lock()
}

func (m *StandardRWMutex) Unlock() {
	m.mu.Unlock()
}

func (m *StandardRWMutex) RLock() {
	m.mu.RLock()
}

func (m *StandardRWMutex) RUnlock() {
	m.mu.RUnlock()
}

func (m *StandardRWMutex) TryLock() bool {
	// TryLock is not available before Go 1.18
	panic("TryLock requires Go 1.18 or later")
}

func (m *StandardRWMutex) TryRLock() bool {
	// TryRLock is not available before Go 1.18
	panic("TryRLock requires Go 1.18 or later")
}

func (m *StandardRWMutex) RLocker() sync.Locker {
	return m.mu.RLocker()
}

// Default factory functions
func newStandardMutex() MutexImpl {
	return &StandardMutex{}
}

func newStandardRWMutex() RWMutexImpl {
	return &StandardRWMutex{}
}