// +build go1.18

package deadlock

// TryLock tries to lock the mutex.
// Returns false if the lock is already in use, true otherwise.
func (m *Mutex) TryLock() bool {
	return trylock(m.mu.TryLock, m)
}

// TryLock tries to lock rw for writing.
// Returns false if the lock is already locked for reading or writing, true otherwise.
func (m *RWMutex) TryLock() bool {
	return trylock(m.mu.TryLock, m)
}

// TryRLock tries to lock rw for reading.
// Returns false if the lock is already locked for writing, true otherwise.
func (m *RWMutex) TryRLock() bool {
	return trylock(m.mu.TryRLock, m)
}

// trylock can not deadlock, so there is no deadlock detection.
// lock ordering is still supported by calling into preLock/postLock,
// and in failed attempt into postUnlock to unroll the state added by preLock.
func trylock(lockFn func() bool, ptr interface{}) bool {
	if Opts.Disable {
		return lockFn()
	}
	stack := callers(1)
	preLock(stack, ptr)
	ret := lockFn()
	if ret {
		postLock(stack, ptr)
	} else {
		postUnlock(ptr)
	}		
	return ret
}
