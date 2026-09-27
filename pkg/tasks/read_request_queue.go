package tasks

import "sync"

// readRequestQueue is an unbounded, order-preserving FIFO of the read requests a
// view's running command task serves (see LinesToRead), with a reader that comes
// and goes.
//
// It's unbounded, rather than a fixed-size channel, for the same reasons as the
// user-event queue in gocui. Requests are handed over from the UI thread, where a
// blocking send would deadlock against the task that is waiting to be let go, and
// a fixed channel that fills up leaves only bad choices: blocking, dropping,
// reordering, or panicking on overflow. Appending to a slice does none of those.
//
// The reader coming and going is the other half of what it's for. A request is
// how a caller asks for content to be read and hears, through the request's Then,
// that it has been; a request nobody answers leaves that caller waiting for good.
// So asking whether a task is there and handing it the request are one step, and
// so are taking the task away and handing back what it never answered. A request
// made in between finds no task and goes back to its caller to answer.
//
// enqueue appends under the mutex and rings the doorbell; the task selects on the
// doorbell to wake, then takes requests until there are none left. The doorbell is
// buffered(1) and rung with a non-blocking send, so it's a coalescing "work
// pending" flag rather than a per-request signal: a burst of appends leaves at
// most one token, and the task takes everything the token stands for on a single
// wake. A token left over after the queue empties causes one harmless empty wake.
type readRequestQueue struct {
	mutex    sync.Mutex
	requests []LinesToRead
	doorbell chan struct{}

	// Whether a task is there to serve the requests. False before the first task
	// starts, and between one task ending and the next starting.
	serving bool
}

func newReadRequestQueue() *readRequestQueue {
	return &readRequestQueue{doorbell: make(chan struct{}, 1)}
}

// beginServing says that a task is now there to serve the queue, and returns the
// doorbell that tells it when there is something to serve.
func (self *readRequestQueue) beginServing() <-chan struct{} {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	self.serving = true
	return self.doorbell
}

// stopServing takes the task away and hands back the requests it never answered,
// for the caller to answer in its place.
func (self *readRequestQueue) stopServing() []LinesToRead {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	self.serving = false
	unanswered := self.requests
	self.requests = nil
	return unanswered
}

// enqueue gives a request to the task serving the queue, and reports whether
// there was one to give it to. When there wasn't, the request is the caller's to
// answer.
func (self *readRequestQueue) enqueue(request LinesToRead) bool {
	self.mutex.Lock()
	if !self.serving {
		self.mutex.Unlock()
		return false
	}
	self.requests = append(self.requests, request)
	self.mutex.Unlock()

	select {
	case self.doorbell <- struct{}{}:
	default:
	}
	return true
}

// dequeue takes the oldest request, reporting false when there are none.
func (self *readRequestQueue) dequeue() (LinesToRead, bool) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if len(self.requests) == 0 {
		return LinesToRead{}, false
	}
	request := self.requests[0]
	if len(self.requests) == 1 {
		// Release the backing array whenever the queue drains, so a one-off burst
		// doesn't pin its peak size for the rest of the session.
		self.requests = nil
	} else {
		self.requests[0] = LinesToRead{}
		self.requests = self.requests[1:]
	}
	return request, true
}
