package lo

import (
	"sync"
	"time"

	"github.com/samber/lo/internal/xtime"
)

type debounce struct {
	after     time.Duration
	mu        *sync.Mutex
	timer     *time.Timer
	done      bool
	callbacks []func()
}

func (d *debounce) reset() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.done {
		return
	}

	if d.timer != nil {
		d.timer.Stop()
	}

	d.timer = time.AfterFunc(d.after, func() {
		// We need to lock the mutex here to avoid race conditions with 2 concurrent calls to reset()
		d.mu.Lock()
		callbacks := append([]func(){}, d.callbacks...)
		d.mu.Unlock()

		for i := range callbacks {
			callbacks[i]()
		}
	})
}

func (d *debounce) cancel() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
		d.timer = nil
	}

	d.done = true
}

// NewDebounce creates a debounced instance that delays invoking functions given until after wait milliseconds have elapsed.
// Play: https://go.dev/play/p/mz32VMK2nqe
func NewDebounce(duration time.Duration, f ...func()) (func(), func()) {
	d := &debounce{
		after:     duration,
		mu:        new(sync.Mutex),
		timer:     nil,
		done:      false,
		callbacks: f,
	}

	return func() {
		d.reset()
	}, d.cancel
}

type debounceByItem struct {
	mu    *sync.Mutex
	timer *time.Timer
	count int
}

type debounceBy[T comparable] struct {
	after     time.Duration
	mu        *sync.Mutex
	items     map[T]*debounceByItem
	callbacks []func(key T, count int)
}

func (d *debounceBy[T]) reset(key T) {
	d.mu.Lock()
	if _, ok := d.items[key]; !ok {
		d.items[key] = &debounceByItem{
			mu:    new(sync.Mutex),
			timer: nil,
		}
	}

	item := d.items[key]

	d.mu.Unlock()

	item.mu.Lock()
	defer item.mu.Unlock()

	item.count++

	if item.timer != nil {
		item.timer.Stop()
	}

	item.timer = time.AfterFunc(d.after, func() {
		// We need to lock the mutex here to avoid race conditions with 2 concurrent calls to reset()
		item.mu.Lock()
		count := item.count
		item.count = 0
		callbacks := append([]func(key T, count int){}, d.callbacks...)
		item.mu.Unlock()

		for i := range callbacks {
			callbacks[i](key, count)
		}
	})
}

func (d *debounceBy[T]) cancel(key T) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if item, ok := d.items[key]; ok {
		item.mu.Lock()

		if item.timer != nil {
			item.timer.Stop()
			item.timer = nil
		}

		item.mu.Unlock()

		delete(d.items, key)
	}
}

// NewDebounceBy creates a debounced instance for each distinct key, that delays invoking functions given until after wait milliseconds have elapsed.
// Play: https://go.dev/play/p/d3Vpt6pxhY8
func NewDebounceBy[T comparable](duration time.Duration, f ...func(key T, count int)) (func(key T), func(key T)) {
	d := &debounceBy[T]{
		after:     duration,
		mu:        new(sync.Mutex),
		items:     map[T]*debounceByItem{},
		callbacks: f,
	}

	return func(key T) {
		d.reset(key)
	}, d.cancel
}

// Attempt invokes a function N times until it returns valid output. Returns either the caught error or nil.
// When the first argument is less than `1`, the function runs until a successful response is returned.
// Play: https://go.dev/play/p/3ggJZ2ZKcMj
func Attempt(maxIteration int, f func(index int) error) (int, error) {
	var err error

	for i := 0; maxIteration <= 0 || i < maxIteration; i++ {
		// for retries >= 0 {
		err = f(i)
		if err == nil {
			return i + 1, nil
		}
	}

	return maxIteration, err
}

// AttemptWithDelay invokes a function N times until it returns valid output,
// with a pause between each call. Returns either the caught error or nil.
// When the first argument is less than `1`, the function runs until a successful
// response is returned.
// Play: https://go.dev/play/p/tVs6CygC7m1
func AttemptWithDelay(maxIteration int, delay time.Duration, f func(index int, duration time.Duration) error) (int, time.Duration, error) {
	var err error

	start := xtime.Now()

	for i := 0; maxIteration <= 0 || i < maxIteration; i++ {
		err = f(i, xtime.Since(start))
		if err == nil {
			return i + 1, xtime.Since(start), nil
		}

		if maxIteration <= 0 || i+1 < maxIteration {
			xtime.Sleep(delay)
		}
	}

	return maxIteration, xtime.Since(start), err
}

// AttemptWhile invokes a function N times until it returns valid output.
// Returns either the caught error or nil, along with a bool value to determine
// whether the function should be invoked again. It will terminate the invoke
// immediately if the second return value is false. When the first
// argument is less than `1`, the function runs until a successful response is
// returned.
// Play: https://go.dev/play/p/1VS7HxlYMOG
func AttemptWhile(maxIteration int, f func(int) (error, bool)) (int, error) {
	var err error
	var shouldContinueInvoke bool

	for i := 0; maxIteration <= 0 || i < maxIteration; i++ {
		// for retries >= 0 {
		err, shouldContinueInvoke = f(i)
		if !shouldContinueInvoke { // if shouldContinueInvoke is false, then return immediately
			return i + 1, err
		}
		if err == nil {
			return i + 1, nil
		}
	}

	return maxIteration, err
}

// AttemptWhileWithDelay invokes a function N times until it returns valid output,
// with a pause between each call. Returns either the caught error or nil, along
// with a bool value to determine whether the function should be invoked again.
// It will terminate the invoke immediately if the second return value is false.
// When the first argument is less than `1`, the function runs until a successful
// response is returned.
// Play: https://go.dev/play/p/mhufUjJfLEF
func AttemptWhileWithDelay(maxIteration int, delay time.Duration, f func(int, time.Duration) (error, bool)) (int, time.Duration, error) {
	var err error
	var shouldContinueInvoke bool

	start := xtime.Now()

	for i := 0; maxIteration <= 0 || i < maxIteration; i++ {
		err, shouldContinueInvoke = f(i, xtime.Since(start))
		if !shouldContinueInvoke { // if shouldContinueInvoke is false, then return immediately
			return i + 1, xtime.Since(start), err
		}
		if err == nil {
			return i + 1, xtime.Since(start), nil
		}

		if maxIteration <= 0 || i+1 < maxIteration {
			xtime.Sleep(delay)
		}
	}

	return maxIteration, xtime.Since(start), err
}

type transactionStep[T any] struct {
	exec       func(T) (T, error)
	onRollback func(T) T
}

// NewTransaction instantiate a new transaction.
// Play: https://go.dev/play/p/Qxrd7MGQGh1
func NewTransaction[T any]() *Transaction[T] {
	return &Transaction[T]{
		steps: []transactionStep[T]{},
	}
}

// Transaction implements a Saga pattern.
type Transaction[T any] struct {
	steps []transactionStep[T]
}

// Then adds a step to the chain of callbacks. Returns the same Transaction.
// Play: https://go.dev/play/p/Qxrd7MGQGh1 https://go.dev/play/p/xrHb2_kMvTY
func (t *Transaction[T]) Then(exec func(T) (T, error), onRollback func(T) T) *Transaction[T] {
	t.steps = append(t.steps, transactionStep[T]{
		exec:       exec,
		onRollback: onRollback,
	})

	return t
}

// Process runs the Transaction steps and rollbacks in case of errors.
// Play: https://go.dev/play/p/Qxrd7MGQGh1 https://go.dev/play/p/xrHb2_kMvTY
func (t *Transaction[T]) Process(state T) (T, error) {
	var i int
	var err error

	for i < len(t.steps) {
		state, err = t.steps[i].exec(state)
		if err != nil {
			break
		}

		i++
	}

	if err == nil {
		return state, nil
	}

	for i > 0 {
		i--
		state = t.steps[i].onRollback(state)
	}

	return state, err
}

// @TODO: single mutex per key?
type throttleBy[T comparable] struct {
	mu         *sync.Mutex
	timer      *time.Timer
	interval   time.Duration
	callbacks  []func(key T)
	countLimit int
	count      map[T]int
}

func (th *throttleBy[T]) throttledFunc(key T) {
	th.mu.Lock()
	defer th.mu.Unlock()

	if th.count[key] < th.countLimit {
		th.count[key]++

		for _, f := range th.callbacks {
			f(key)
		}
	}
	if th.timer == nil {
		th.timer = time.AfterFunc(th.interval, func() {
			th.reset()
		})
	}
}

func (th *throttleBy[T]) reset() {
	th.mu.Lock()
	defer th.mu.Unlock()

	if th.timer != nil {
		th.timer.Stop()
	}

	th.count = map[T]int{}
	th.timer = nil
}

// NewThrottle creates a throttled instance that invokes given functions only once in every interval.
// This returns 2 functions, First one is throttled function and Second one is a function to reset interval.
// Play: https://go.dev/play/p/qQn3fm8Z7jS
func NewThrottle(interval time.Duration, f ...func()) (throttle, reset func()) {
	return NewThrottleWithCount(interval, 1, f...)
}

// NewThrottleWithCount is NewThrottle with count limit, throttled function will be invoked count times in every interval.
// Play: https://go.dev/play/p/w5nc0MgWtjC
func NewThrottleWithCount(interval time.Duration, count int, f ...func()) (throttle, reset func()) {
	callbacks := Map(f, func(item func(), _ int) func(struct{}) {
		return func(struct{}) {
			item()
		}
	})

	throttleFn, reset := NewThrottleByWithCount(interval, count, callbacks...)
	return func() {
		throttleFn(struct{}{})
	}, reset
}

// NewThrottleBy creates a throttled instance that invokes given functions only once in every interval.
// This returns 2 functions, First one is throttled function and Second one is a function to reset interval.
// Play: https://go.dev/play/p/0Wv6oX7dHdC
func NewThrottleBy[T comparable](interval time.Duration, f ...func(key T)) (throttle func(key T), reset func()) {
	return NewThrottleByWithCount(interval, 1, f...)
}

// NewThrottleByWithCount is NewThrottleBy with count limit, throttled function will be invoked count times in every interval.
// Play: https://go.dev/play/p/vQk3ECH7_EW
func NewThrottleByWithCount[T comparable](interval time.Duration, count int, f ...func(key T)) (throttle func(key T), reset func()) {
	if count <= 0 {
		count = 1
	}

	th := &throttleBy[T]{
		mu:         new(sync.Mutex),
		interval:   interval,
		callbacks:  f,
		countLimit: count,
		count:      map[T]int{},
	}
	return th.throttledFunc, th.reset
}
