//go:build deadlock_disable && !go1.18

package deadlock

import "sync"

// StandardMutex wraps sync.Mutex with no deadlock detection
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
	panic("TryLock requires Go 1.18 or later")
}

// StandardRWMutex wraps sync.RWMutex with no deadlock detection
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
	panic("TryLock requires Go 1.18 or later")
}

func (m *StandardRWMutex) TryRLock() bool {
	panic("TryRLock requires Go 1.18 or later")
}

func (m *StandardRWMutex) RLocker() sync.Locker {
	return m.mu.RLocker()
}
