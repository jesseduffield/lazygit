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
// A caller acts on the content the view shows once it is answered, so a request
// is held for the task asked for last, from the moment it is asked for. When a
// task is replaced, the requests it hasn't answered stay for the task replacing it.
// Only once the task asked for last has given the view all of its content does the
// queue stop holding requests, and hand back the ones it still holds to answer. A
// request made after that finds no task and goes back to its caller to answer.
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

	// Whether requests are held for a task to serve: the one reading the content, or
	// one that has been asked for and is going to replace it. False before the first
	// task is asked for, and once the task asked for last has read all of its
	// content.
	holding bool
}

func newReadRequestQueue() *readRequestQueue {
	return &readRequestQueue{doorbell: make(chan struct{}, 1)}
}

// hold says that a task has been asked for, so that requests are held for it until
// it serves them.
func (self *readRequestQueue) hold() {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	self.holding = true
}

// beginServing says that a task is now there to serve the queue, and returns the
// doorbell that tells it when there is something to serve. The task's own first
// request is served first, ahead of those it takes over from the task it replaces,
// so that it shows its content as soon as it would have without them.
func (self *readRequestQueue) beginServing(firstRequest LinesToRead) <-chan struct{} {
	self.mutex.Lock()
	self.holding = true
	self.requests = append([]LinesToRead{firstRequest}, self.requests...)
	self.mutex.Unlock()

	self.ringDoorbell()
	return self.doorbell
}

// stopServing stops holding requests and hands back the ones it still holds, for
// the caller to answer in place of a task.
func (self *readRequestQueue) stopServing() []LinesToRead {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	self.holding = false
	unanswered := self.requests
	self.requests = nil
	return unanswered
}

// enqueue gives a request to the task that is going to serve it, and reports
// whether there is one. When there isn't, the request is the caller's to answer.
func (self *readRequestQueue) enqueue(request LinesToRead) bool {
	self.mutex.Lock()
	if !self.holding {
		self.mutex.Unlock()
		return false
	}
	self.requests = append(self.requests, request)
	self.mutex.Unlock()

	self.ringDoorbell()
	return true
}

// requeue puts back a request that a task was serving when it was stopped, ahead of
// the ones made after it, for the task replacing it to serve.
func (self *readRequestQueue) requeue(request LinesToRead) {
	self.mutex.Lock()
	self.requests = append([]LinesToRead{request}, self.requests...)
	self.mutex.Unlock()

	self.ringDoorbell()
}

func (self *readRequestQueue) ringDoorbell() {
	select {
	case self.doorbell <- struct{}{}:
	default:
	}
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
