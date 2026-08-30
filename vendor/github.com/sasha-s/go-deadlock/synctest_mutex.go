//go:build (goexperiment.synctest || deadlock_synctest) && !deadlock_disable

package deadlock

import (
	"sync"
	"sync/atomic"
)

// ChannelMutex implements MutexImpl using channels for synctest compatibility
type ChannelMutex struct {
	ch     chan struct{}
	locked int32 // atomic
	once   sync.Once
}

func (m *ChannelMutex) init() {
	m.once.Do(func() {
		m.ch = make(chan struct{}, 1)
	})
}

func (m *ChannelMutex) Lock() {
	m.init()
	m.ch <- struct{}{}
	atomic.StoreInt32(&m.locked, 1)
}

func (m *ChannelMutex) Unlock() {
	if atomic.LoadInt32(&m.locked) == 0 {
		panic("unlock of unlocked mutex")
	}
	atomic.StoreInt32(&m.locked, 0)
	<-m.ch
}

func (m *ChannelMutex) TryLock() bool {
	m.init()
	select {
	case m.ch <- struct{}{}:
		atomic.StoreInt32(&m.locked, 1)
		return true
	default:
		return false
	}
}

// ChannelRWMutex implements RWMutexImpl with writer priority using channels.
// Implements writer-priority semantics using the "Third Readers-Writers Problem" solution.
type ChannelRWMutex struct {
	resource    chan struct{} // The actual resource being protected
	readTry     chan struct{} // Gate that closes when writers waiting
	rmutex      chan struct{} // Protects readCount modifications
	wmutex      chan struct{} // Protects writeCount modifications
	readCount   int32         // Number of active readers
	writeCount  int32         // Number of waiting/active writers
	once        sync.Once
}

func (m *ChannelRWMutex) init() {
	m.once.Do(func() {
		m.resource = make(chan struct{}, 1)
		m.readTry = make(chan struct{}, 1)
		m.rmutex = make(chan struct{}, 1)
		m.wmutex = make(chan struct{}, 1)
		// Initially, all semaphores are "released" (have a token)
		m.resource <- struct{}{}
		m.readTry <- struct{}{}
		m.rmutex <- struct{}{}
		m.wmutex <- struct{}{}
	})
}

func (m *ChannelRWMutex) Lock() {
	m.init()
	// Protect writeCount modification
	<-m.wmutex
	count := atomic.AddInt32(&m.writeCount, 1)
	if count == 1 {
		// First writer: close the gate to block new readers
		<-m.readTry
	}
	m.wmutex <- struct{}{} // Release wmutex

	// Acquire the resource (wait for existing readers to finish)
	<-m.resource
}

func (m *ChannelRWMutex) Unlock() {
	// Release the resource
	m.resource <- struct{}{}

	// Protect writeCount modification
	<-m.wmutex
	count := atomic.AddInt32(&m.writeCount, -1)
	if count == 0 {
		// Last writer: reopen the gate for readers
		m.readTry <- struct{}{}
	}
	m.wmutex <- struct{}{}
}

func (m *ChannelRWMutex) RLock() {
	m.init()
	// Wait at the gate (blocks if writers are waiting)
	<-m.readTry

	// Protect readCount modification
	<-m.rmutex
	count := atomic.AddInt32(&m.readCount, 1)
	if count == 1 {
		// First reader: acquire the resource to block writers
		<-m.resource
	}
	m.rmutex <- struct{}{} // Release rmutex

	// Release the gate so other readers can pass
	m.readTry <- struct{}{}
}

func (m *ChannelRWMutex) RUnlock() {
	<-m.rmutex
	count := atomic.AddInt32(&m.readCount, -1)
	if count < 0 {
		m.rmutex <- struct{}{}
		panic("RUnlock of unlocked RWMutex")
	}
	if count == 0 {
		// Last reader: release the resource for writers
		m.resource <- struct{}{}
	}
	m.rmutex <- struct{}{}
}

func (m *ChannelRWMutex) TryLock() bool {
	m.init()
	// Try to acquire wmutex
	select {
	case <-m.wmutex:
	default:
		return false
	}

	count := atomic.AddInt32(&m.writeCount, 1)
	if count == 1 {
		// First writer: try to close gate
		select {
		case <-m.readTry:
		default:
			// Failed, rollback
			atomic.AddInt32(&m.writeCount, -1)
			m.wmutex <- struct{}{}
			return false
		}
	}
	m.wmutex <- struct{}{}

	// Try to acquire resource
	select {
	case <-m.resource:
		return true
	default:
		// Failed, rollback writer count
		<-m.wmutex
		count = atomic.AddInt32(&m.writeCount, -1)
		if count == 0 {
			m.readTry <- struct{}{}
		}
		m.wmutex <- struct{}{}
		return false
	}
}

func (m *ChannelRWMutex) TryRLock() bool {
	m.init()
	// Try to pass through the gate
	select {
	case <-m.readTry:
	default:
		return false
	}

	// Try to acquire rmutex
	select {
	case <-m.rmutex:
	default:
		// Failed, release gate
		m.readTry <- struct{}{}
		return false
	}

	count := atomic.AddInt32(&m.readCount, 1)
	if count == 1 {
		// First reader: try to acquire resource
		select {
		case <-m.resource:
		default:
			// Failed, rollback
			atomic.AddInt32(&m.readCount, -1)
			m.rmutex <- struct{}{}
			m.readTry <- struct{}{}
			return false
		}
	}
	m.rmutex <- struct{}{}
	m.readTry <- struct{}{}
	return true
}

func (m *ChannelRWMutex) RLocker() sync.Locker {
	return (*channelRLocker)(m)
}

type channelRLocker ChannelRWMutex

func (r *channelRLocker) Lock()   { (*ChannelRWMutex)(r).RLock() }
func (r *channelRLocker) Unlock() { (*ChannelRWMutex)(r).RUnlock() }

// Factory functions for synctest
func newChannelMutex() MutexImpl {
	return &ChannelMutex{
		ch: make(chan struct{}, 1),
	}
}

func newChannelRWMutex() RWMutexImpl {
	m := &ChannelRWMutex{}
	m.init()
	return m
}

// Type aliases to override the standard mutex types for synctest
type StandardMutex = ChannelMutex
type StandardRWMutex = ChannelRWMutex
