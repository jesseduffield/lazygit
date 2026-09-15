package deadlock

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/petermattis/goid"
)

// TimerPoolMode controls timer pooling behavior
type TimerPoolMode int

const (
	// TimerPoolDefault automatically chooses based on build environment
	TimerPoolDefault TimerPoolMode = iota
	// TimerPoolEnabled always uses timer pooling for performance
	TimerPoolEnabled
	// TimerPoolDisabled disables timer pooling (required for testing/synctest)
	TimerPoolDisabled
)

// Opts control how deadlock detection behaves.
// Options are supposed to be set once at a startup (say, when parsing flags).
var Opts = struct {
	// Mutex/RWMutex would work exactly as their sync counterparts
	// -- almost no runtime penalty, no deadlock detection if Disable == true.
	Disable bool
	// Would disable lock order based deadlock detection if DisableLockOrderDetection == true.
	DisableLockOrderDetection bool
	// Waiting for a lock for longer than DeadlockTimeout is considered a deadlock.
	// Ignored if DeadlockTimeout <= 0.
	DeadlockTimeout time.Duration
	// OnPotentialDeadlock is called each time a potential deadlock is detected -- either based on
	// lock order or on lock wait time.
	OnPotentialDeadlock func()
	// Will keep MaxMapSize lock pairs (happens before // happens after) in the map.
	// The map resets once the threshold is reached.
	MaxMapSize int
	// Will dump stacktraces of all goroutines when inconsistent locking is detected.
	PrintAllCurrentGoroutines bool
	// Controls timer pooling behavior.
	// TimerPoolDefault: Automatically choose based on build environment
	// TimerPoolEnabled: Always use timer pooling
	// TimerPoolDisabled: Never use timer pooling
	TimerPool TimerPoolMode
	mu        *sync.Mutex // Protects the LogBuf.
	// Will print deadlock info to log buffer.
	LogBuf io.Writer
}{
	DeadlockTimeout: time.Second * 30,
	OnPotentialDeadlock: func() {
		os.Exit(2)
	},
	MaxMapSize: 1024 * 64,
	mu:         &sync.Mutex{},
	LogBuf:     os.Stderr,
}

// Cond is sync.Cond wrapper
type Cond struct {
	sync.Cond
}

// Locker is sync.Locker wrapper
type Locker struct {
	sync.Locker
}

// Once is sync.Once wrapper
type Once struct {
	sync.Once
}

// Pool is sync.Poll wrapper
type Pool struct {
	sync.Pool
}

// WaitGroup is sync.WaitGroup wrapper
type WaitGroup struct {
	sync.WaitGroup
}

// NewCond is a sync.NewCond wrapper
var NewCond = sync.NewCond

// A Mutex is a drop-in replacement for sync.Mutex.
// Performs deadlock detection unless disabled in Opts.
type Mutex struct {
	mu StandardMutex
}

// Lock locks the mutex.
// If the lock is already in use, the calling goroutine
// blocks until the mutex is available.
//
// Unless deadlock detection is disabled, logs potential deadlocks to Opts.LogBuf,
// calling Opts.OnPotentialDeadlock on each occasion.
func (m *Mutex) Lock() {
	lock(m.mu.Lock, m)
}

// Unlock unlocks the mutex.
// It is a run-time error if m is not locked on entry to Unlock.
//
// A locked Mutex is not associated with a particular goroutine.
// It is allowed for one goroutine to lock a Mutex and then
// arrange for another goroutine to unlock it.
func (m *Mutex) Unlock() {
	m.mu.Unlock()
	if !Opts.Disable {
		postUnlock(m)
	}
}

// An RWMutex is a drop-in replacement for sync.RWMutex.
// Performs deadlock detection unless disabled in Opts.
type RWMutex struct {
	mu StandardRWMutex
}

// Lock locks rw for writing.
// If the lock is already locked for reading or writing,
// Lock blocks until the lock is available.
// To ensure that the lock eventually becomes available,
// a blocked Lock call excludes new readers from acquiring
// the lock.
//
// Unless deadlock detection is disabled, logs potential deadlocks to Opts.LogBuf,
// calling Opts.OnPotentialDeadlock on each occasion.
func (m *RWMutex) Lock() {
	lock(m.mu.Lock, m)
}

// Unlock unlocks the mutex for writing.  It is a run-time error if rw is
// not locked for writing on entry to Unlock.
//
// As with Mutexes, a locked RWMutex is not associated with a particular
// goroutine.  One goroutine may RLock (Lock) an RWMutex and then
// arrange for another goroutine to RUnlock (Unlock) it.
func (m *RWMutex) Unlock() {
	m.mu.Unlock()
	if !Opts.Disable {
		postUnlock(m)
	}
}

// RLock locks the mutex for reading.
//
// Unless deadlock detection is disabled, logs potential deadlocks to Opts.LogBuf,
// calling Opts.OnPotentialDeadlock on each occasion.
func (m *RWMutex) RLock() {
	lock(m.mu.RLock, m)
}

// RUnlock undoes a single RLock call;
// it does not affect other simultaneous readers.
// It is a run-time error if rw is not locked for reading
// on entry to RUnlock.
func (m *RWMutex) RUnlock() {
	m.mu.RUnlock()
	if !Opts.Disable {
		postUnlock(m)
	}
}

// RLocker returns a Locker interface that implements
// the Lock and Unlock methods by calling RLock and RUnlock.
func (m *RWMutex) RLocker() sync.Locker {
	return m.mu.RLocker()
}

func preLock(stack []uintptr, p interface{}) {
	lo.preLock(stack, p)
}

func postLock(stack []uintptr, buf *[stackBufSize]uintptr, p interface{}) {
	lo.postLock(stack, buf, p)
}

func postUnlock(p interface{}) {
	lo.postUnlock(p)
}

func lock(lockFn func(), ptr interface{}) {
	if Opts.Disable {
		lockFn()
		return
	}
	stack, buf := callers(1)
	// Cache timeout before preLock so all Opts reads complete before preLock
	// may call OnPotentialDeadlock. If preLock detects a problem (recursive
	// lock, order violation) the goroutine may block forever in lockFn below,
	// and reading Opts after preLock would race with any later Opts write.
	timeout := Opts.DeadlockTimeout
	preLock(stack, ptr)
	if timeout <= 0 {
		lockFn()
	} else {
		currentID := goid.Get()
		e := dw.register(stack, ptr, currentID, timeout)
		lockFn()
		dw.deregister(e)
		postLock(stack, buf, ptr)
		return
	}
	postLock(stack, buf, ptr)
}

// pendingEntry tracks a goroutine that is waiting to acquire a lock. Entries are
// pooled to avoid per-lock heap allocations (goroutine stacks, channels, closures).
//
// Timer safety invariants:
//   - checkFn is allocated once per entry and reused across pool cycles, so recycling
//     an entry does not allocate a new closure.
//   - The done flag synchronizes the callback with deregister: deregister sets done=1
//     before calling Stop(), and the callback checks done before acting. Because both
//     use atomic operations, the callback is guaranteed to observe done=1 if deregister
//     has already run, even if the runtime already scheduled the callback.
//   - An entry is only returned to the pool when timer.Stop() returns true, meaning
//     the timer was successfully cancelled and the callback will never run. This prevents
//     a recycled entry from being mutated by an in-flight callback.
//   - When Stop() returns false (callback already firing or queued), the entry is
//     intentionally leaked to GC. This only happens in the rare deadlock-timeout path.
type pendingEntry struct {
	stack   []uintptr
	ptr     interface{}
	gid     int64
	done    int32 // atomic: 0=pending, 1=acquired
	timer   *time.Timer
	checkFn func()
}

func newPendingEntry() *pendingEntry {
	e := &pendingEntry{}
	// Capture e by pointer so the closure is stable across pool reuse, no new
	// closure allocation when the entry is recycled.
	e.checkFn = func() {
		// If the lock was acquired (done=1), the entry may already be back in the
		// pool or being reused. Bail out unconditionally.
		if atomic.LoadInt32(&e.done) != 0 {
			return
		}
		onDeadlockTimeout(e)
	}
	return e
}

var pendingPool = sync.Pool{
	New: func() interface{} {
		return newPendingEntry()
	},
}

type deadlockWatcher struct{}

var dw deadlockWatcher

func (w *deadlockWatcher) register(stack []uintptr, ptr interface{}, gid int64, timeout time.Duration) *pendingEntry {
	var e *pendingEntry
	if shouldDisableTimerPool() {
		e = newPendingEntry()
	} else {
		e = pendingPool.Get().(*pendingEntry)
	}
	e.stack = stack
	e.ptr = ptr
	e.gid = gid
	atomic.StoreInt32(&e.done, 0)
	if e.timer == nil {
		// First use (freshly allocated entry): create the AfterFunc timer.
		// AfterFunc avoids the channel-drain problems of channel-based timers,
		// which are especially problematic under testing/synctest.
		e.timer = time.AfterFunc(timeout, e.checkFn)
	} else {
		// Reused from pool: the timer was previously Stop()'d successfully
		// (guaranteed by deregister), so Reset is safe here.
		e.timer.Reset(timeout)
	}
	return e
}

// deregister marks the lock as acquired and cancels the deadlock timer.
// Must be called exactly once per register call. The entry pointer is
// stack-local in lock(), so concurrent or duplicate calls cannot occur.
func (w *deadlockWatcher) deregister(e *pendingEntry) {
	// Mark done BEFORE stopping the timer. The callback checks done with an
	// atomic load, so even if the timer fires concurrently, the callback will
	// see done=1 and return without acting.
	atomic.StoreInt32(&e.done, 1)
	stopped := e.timer.Stop()
	// Only recycle the entry if Stop() confirmed the callback won't run.
	// If Stop() returned false the callback is already executing or queued;
	// recycling would race with the callback reading entry fields.
	if stopped && !shouldDisableTimerPool() {
		e.stack = nil
		e.ptr = nil
		e.gid = 0
		pendingPool.Put(e)
	}
}

func onDeadlockTimeout(e *pendingEntry) {
	lo.mu.Lock()
	holders, ok := lo.cur[e.ptr]
	if !ok || len(holders) == 0 {
		// Lock appears unheld (transient state, holder may have just released).
		// Reschedule if the waiter is still pending. Note: this creates a new timer
		// (e.timer is not updated), so if deregister runs later it will Stop() the
		// original (already-fired) timer, get false, and skip pooling. The new timer's
		// callback will then observe done=1 and no-op. This is safe but means the
		// entry won't be recycled, acceptable since this is the rare timeout path.
		lo.mu.Unlock()
		if atomic.LoadInt32(&e.done) == 0 {
			time.AfterFunc(Opts.DeadlockTimeout, e.checkFn)
		}
		return
	}
	Opts.mu.Lock()
	fmt.Fprintln(Opts.LogBuf, header)
	for _, prev := range holders {
		fmt.Fprintln(Opts.LogBuf, "Previous place where the lock was grabbed")
		fmt.Fprintf(Opts.LogBuf, "goroutine %v lock %p\n", prev.gid, e.ptr)
		printStack(Opts.LogBuf, prev.stack)
	}
	fmt.Fprintln(Opts.LogBuf, "Have been trying to lock it again for more than", Opts.DeadlockTimeout)
	fmt.Fprintf(Opts.LogBuf, "goroutine %v lock %p\n", e.gid, e.ptr)
	printStack(Opts.LogBuf, e.stack)
	stacks := stacks()
	grs := bytes.Split(stacks, []byte("\n\n"))
	for _, prev := range holders {
		for _, g := range grs {
			if goid.ExtractGID(g) == prev.gid {
				fmt.Fprintln(Opts.LogBuf, "Here is what goroutine", prev.gid, "doing now")
				Opts.LogBuf.Write(g)
				fmt.Fprintln(Opts.LogBuf)
			}
		}
	}
	lo.other(e.ptr)
	if Opts.PrintAllCurrentGoroutines {
		fmt.Fprintln(Opts.LogBuf, "All current goroutines:")
		Opts.LogBuf.Write(stacks)
	}
	fmt.Fprintln(Opts.LogBuf)
	if buf, ok := Opts.LogBuf.(*bufio.Writer); ok {
		buf.Flush()
	}
	Opts.mu.Unlock()
	lo.mu.Unlock()
	Opts.OnPotentialDeadlock()
}

type lockOrder struct {
	mu    sync.Mutex
	cur   map[interface{}][]stackGID // stacktraces + gids for the locks currently taken.
	order map[beforeAfter]ss         // expected order of locks.
}

type stackGID struct {
	stack []uintptr
	gid   int64
	buf   *[stackBufSize]uintptr // pooled backing array; returned via releaseStackBuf in postUnlock
}

type ss struct {
	before []uintptr
	after  []uintptr
}

var lo = newLockOrder()

func newLockOrder() *lockOrder {
	return &lockOrder{
		cur:   map[interface{}][]stackGID{},
		order: map[beforeAfter]ss{},
	}
}

// holdersPool recycles []stackGID slices used by lockOrder.cur to track which
// goroutines currently hold each lock. Slices are returned to the pool in
// postUnlock when a lock's holder count drops to zero, and reused in postLock
// for the next lock acquisition, avoiding a new slice allocation per mutex.
var holdersPool sync.Pool

// postLock records the current goroutine as a holder of lock p. It tries to
// reuse a pooled []stackGID slice before allocating, and stores the pooled
// stack buffer in the entry so postUnlock can release it later.
func (l *lockOrder) postLock(stack []uintptr, buf *[stackBufSize]uintptr, p interface{}) {
	gid := goid.Get()
	entry := stackGID{stack, gid, buf}
	l.mu.Lock()
	holders := l.cur[p]
	if holders == nil {
		if s, ok := holdersPool.Get().([]stackGID); ok {
			holders = s[:0]
		}
	}
	l.cur[p] = append(holders, entry)
	l.mu.Unlock()
}

func (l *lockOrder) preLock(stack []uintptr, p interface{}) {
	if Opts.DisableLockOrderDetection {
		return
	}
	gid := goid.Get()
	l.mu.Lock()
	for b, holders := range l.cur {
		if b == p {
			for _, bs := range holders {
				if bs.gid == gid {
					Opts.mu.Lock()
					fmt.Fprintln(Opts.LogBuf, header, "Recursive locking:")
					fmt.Fprintf(Opts.LogBuf, "current goroutine %d lock %p\n", gid, b)
					printStack(Opts.LogBuf, stack)
					fmt.Fprintln(Opts.LogBuf, "Previous place where the lock was grabbed (same goroutine)")
					printStack(Opts.LogBuf, bs.stack)
					l.other(p)
					if buf, ok := Opts.LogBuf.(*bufio.Writer); ok {
						buf.Flush()
					}
					Opts.mu.Unlock()
					Opts.OnPotentialDeadlock()
					break
				}
			}
			continue
		}
		for _, bs := range holders {
			if bs.gid != gid { // We want locks taken in the same goroutine only.
				continue
			}
			if s, ok := l.order[newBeforeAfter(p, b)]; ok {
				Opts.mu.Lock()
				fmt.Fprintln(Opts.LogBuf, header, "Inconsistent locking. saw this ordering in one goroutine:")
				fmt.Fprintln(Opts.LogBuf, "happened before")
				printStack(Opts.LogBuf, s.before)
				fmt.Fprintln(Opts.LogBuf, "happened after")
				printStack(Opts.LogBuf, s.after)
				fmt.Fprintln(Opts.LogBuf, "in another goroutine: happened before")
				printStack(Opts.LogBuf, bs.stack)
				fmt.Fprintln(Opts.LogBuf, "happened after")
				printStack(Opts.LogBuf, stack)
				l.other(p)
				fmt.Fprintln(Opts.LogBuf)
				if buf, ok := Opts.LogBuf.(*bufio.Writer); ok {
					buf.Flush()
				}
				Opts.mu.Unlock()
				Opts.OnPotentialDeadlock()
			}
			// Copy both stacks: they're backed by pooled buffers that will be
			// recycled in postUnlock, but l.order entries persist until MaxMapSize.
			l.order[newBeforeAfter(b, p)] = ss{copyStack(bs.stack), copyStack(stack)}
			if len(l.order) == Opts.MaxMapSize { // Reset the map to keep memory footprint bounded.
				l.order = map[beforeAfter]ss{}
			}
		}
	}
	l.mu.Unlock()
}

func (l *lockOrder) postUnlock(p interface{}) {
	gid := goid.Get()
	l.mu.Lock()
	holders := l.cur[p]
	idx := -1
	for i, h := range holders {
		if h.gid == gid {
			idx = i
			break
		}
	}
	if idx >= 0 {
		removedBuf := holders[idx].buf
		holders[idx] = holders[len(holders)-1]
		holders[len(holders)-1] = stackGID{}
		holders = holders[:len(holders)-1]
		releaseStackBuf(removedBuf)
	} else if len(holders) > 0 {
		// Cross-goroutine unlock: Go permits one goroutine to Lock and a different
		// goroutine to Unlock, so the unlocking gid may not match any holder entry.
		// This is a rare edge case in practice, the vast majority of code unlocks
		// from the same goroutine that locked. We remove an arbitrary entry to keep
		// the holder count consistent with the real lock state (the lock *was*
		// released, so one entry must go). The trade-off: for RWMutex with multiple
		// concurrent readers we may discard the wrong reader's stack trace, making a
		// future deadlock report show a slightly misleading "previous lock site".
		// Detection correctness is unaffected.
		removedBuf := holders[len(holders)-1].buf
		holders[len(holders)-1] = stackGID{}
		holders = holders[:len(holders)-1]
		releaseStackBuf(removedBuf)
	}
	if len(holders) == 0 {
		// Delete the map key so the mutex pointer is not retained, allowing GC of
		// the struct it's embedded in. Recycle the backing slice via pool so the
		// next postLock on any mutex can reuse it instead of allocating.
		if cap(holders) > 0 {
			holdersPool.Put(holders[:0])
		}
		delete(l.cur, p)
	} else {
		l.cur[p] = holders
	}
	l.mu.Unlock()
}

// Under lo.mu Locked.
func (l *lockOrder) other(ptr interface{}) {
	empty := true
	for k, holders := range l.cur {
		if k == ptr {
			continue
		}
		if len(holders) > 0 {
			empty = false
			break
		}
	}
	if empty {
		return
	}
	fmt.Fprintln(Opts.LogBuf, "Other goroutines holding locks:")
	for k, holders := range l.cur {
		if k == ptr {
			continue
		}
		for _, pp := range holders {
			fmt.Fprintf(Opts.LogBuf, "goroutine %v lock %p\n", pp.gid, k)
			printStack(Opts.LogBuf, pp.stack)
		}
	}
	fmt.Fprintln(Opts.LogBuf)
}

const header = "POTENTIAL DEADLOCK:"
