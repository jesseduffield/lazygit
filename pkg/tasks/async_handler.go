package tasks

import (
	"github.com/jesseduffield/gocui"
	"github.com/sasha-s/go-deadlock"
)

// the purpose of an AsyncHandler is to ensure that if we have multiple long-running
// requests, we only handle the result of the latest one. For example, if I am
// searching for 'abc' and I have to type 'a' then 'b' then 'c' and each keypress
// dispatches a request to search for things with the string so-far, we'll be searching
// for 'a', 'ab', and 'abc', and it may be that 'abc' comes back first, then 'ab',
// then 'a' and we don't want to display the result for 'a' just because it came
// back last. AsyncHandler keeps track of the order in which things were dispatched
// so that we can ignore anything that comes back late.
type AsyncHandler struct {
	currentId int
	lastId    int
	mutex     deadlock.Mutex
	onReject  func()
	onWorker  func(func(gocui.Task))
}

func NewAsyncHandler(onWorker func(func(gocui.Task))) *AsyncHandler {
	return &AsyncHandler{
		mutex:    deadlock.Mutex{},
		onWorker: onWorker,
	}
}

func (self *AsyncHandler) Do(f func() func()) {
	self.mutex.Lock()
	self.currentId++
	id := self.currentId
	self.mutex.Unlock()

	self.onWorker(func(gocui.Task) {
		after := f()
		self.handle(after, id)
	})
}

// f here is expected to be a function that doesn't take long to run
func (self *AsyncHandler) handle(f func(), id int) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if id < self.lastId {
		if self.onReject != nil {
			self.onReject()
		}
		return
	}

	self.lastId = id
	f()
}
