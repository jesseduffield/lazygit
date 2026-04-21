package lo

import (
	"context"
	"sync"
	"time"
)

type synchronize struct {
	locker sync.Locker
}

func (s *synchronize) Do(callback func()) {
	s.locker.Lock()
	Try0(callback)
	s.locker.Unlock()
}

// Synchronize wraps the underlying callback in a mutex. It receives an optional mutex.
// Play: https://go.dev/play/p/X3cqROSpQmu
func Synchronize(opt ...sync.Locker) *synchronize { //nolint:revive
	if len(opt) > 1 {
		panic("lo.Synchronize: unexpected arguments")
	} else if len(opt) == 0 {
		opt = append(opt, &sync.Mutex{})
	}

	return &synchronize{
		locker: opt[0],
	}
}

// Async executes a function in a goroutine and returns the result in a channel.
// Play: https://go.dev/play/p/uo35gosuTLw
func Async[A any](f func() A) <-chan A {
	ch := make(chan A, 1)
	go func() {
		ch <- f()
	}()
	return ch
}

// Async0 executes a function in a goroutine and returns a channel set once the function finishes.
// Play: https://go.dev/play/p/tNqf1cClG_o
func Async0(f func()) <-chan struct{} {
	ch := make(chan struct{}, 1)
	go func() {
		f()
		ch <- struct{}{}
	}()
	return ch
}

// Async1 is an alias to Async.
// Play: https://go.dev/play/p/RBQWtIn4PsF
func Async1[A any](f func() A) <-chan A {
	return Async(f)
}

// Async2 has the same behavior as Async, but returns the 2 results as a tuple inside the channel.
// Play: https://go.dev/play/p/5SzzDjssXOH
func Async2[A, B any](f func() (A, B)) <-chan Tuple2[A, B] {
	ch := make(chan Tuple2[A, B], 1)
	go func() {
		ch <- T2(f())
	}()
	return ch
}

// Async3 has the same behavior as Async, but returns the 3 results as a tuple inside the channel.
// Play: https://go.dev/play/p/cZpZsDXNmlx
func Async3[A, B, C any](f func() (A, B, C)) <-chan Tuple3[A, B, C] {
	ch := make(chan Tuple3[A, B, C], 1)
	go func() {
		ch <- T3(f())
	}()
	return ch
}

// Async4 has the same behavior as Async, but returns the 4 results as a tuple inside the channel.
// Play: https://go.dev/play/p/9X5O2VrLzkR
func Async4[A, B, C, D any](f func() (A, B, C, D)) <-chan Tuple4[A, B, C, D] {
	ch := make(chan Tuple4[A, B, C, D], 1)
	go func() {
		ch <- T4(f())
	}()
	return ch
}

// Async5 has the same behavior as Async, but returns the 5 results as a tuple inside the channel.
// Play: https://go.dev/play/p/MqnUJpkmopA
func Async5[A, B, C, D, E any](f func() (A, B, C, D, E)) <-chan Tuple5[A, B, C, D, E] {
	ch := make(chan Tuple5[A, B, C, D, E], 1)
	go func() {
		ch <- T5(f())
	}()
	return ch
}

// Async6 has the same behavior as Async, but returns the 6 results as a tuple inside the channel.
// Play: https://go.dev/play/p/kM1X67JPdSP
func Async6[A, B, C, D, E, F any](f func() (A, B, C, D, E, F)) <-chan Tuple6[A, B, C, D, E, F] {
	ch := make(chan Tuple6[A, B, C, D, E, F], 1)
	go func() {
		ch <- T6(f())
	}()
	return ch
}

// WaitFor runs periodically until a condition is validated.
// Play: https://go.dev/play/p/t_wTDmubbK3
func WaitFor(condition func(i int) bool, timeout, heartbeatDelay time.Duration) (totalIterations int, elapsed time.Duration, conditionFound bool) {
	conditionWithContext := func(_ context.Context, currentIteration int) bool {
		return condition(currentIteration)
	}
	return WaitForWithContext(context.Background(), conditionWithContext, timeout, heartbeatDelay)
}

// WaitForWithContext runs periodically until a condition is validated or context is canceled.
// Play: https://go.dev/play/p/t_wTDmubbK3
func WaitForWithContext(ctx context.Context, condition func(ctx context.Context, currentIteration int) bool, timeout, heartbeatDelay time.Duration) (totalIterations int, elapsed time.Duration, conditionFound bool) {
	start := time.Now()

	if ctx.Err() != nil {
		return totalIterations, time.Since(start), false
	}

	ctx, cleanCtx := context.WithTimeout(ctx, timeout)
	ticker := time.NewTicker(heartbeatDelay)

	defer func() {
		cleanCtx()
		ticker.Stop()
	}()

	for {
		select {
		case <-ctx.Done():
			return totalIterations, time.Since(start), false
		case <-ticker.C:
			totalIterations++
			if condition(ctx, totalIterations-1) {
				return totalIterations, time.Since(start), true
			}
		}
	}
}
