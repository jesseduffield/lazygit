package gocui

import "sync"

// Tracks whether the program is busy (i.e. either something is happening on
// the main goroutine or a worker goroutine). Used by integration tests
// to wait until the program is idle before progressing.
type TaskManager struct {
	tasks map[int]Task
	// auto-incrementing id for new tasks
	nextId int

	mutex sync.Mutex
	// signalled whenever the program transitions from busy to idle; used by
	// WaitUntilIdle
	idleCond *sync.Cond
}

func newTaskManager() *TaskManager {
	self := &TaskManager{
		tasks: make(map[int]Task),
	}
	self.idleCond = sync.NewCond(&self.mutex)

	return self
}

func (self *TaskManager) NewTask(background bool) *TaskImpl {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	self.nextId++
	taskId := self.nextId

	onDone := func() { self.delete(taskId) }
	task := &TaskImpl{id: taskId, busy: true, background: background, onDone: onDone, withMutex: self.withMutex}
	self.tasks[taskId] = task

	return task
}

// hasBusyForegroundTaskExcept reports whether any task other than `ignore` is
// currently busy and not a background task. It's used to decide whether a repo
// switch is safe: a foreground operation (or the refresh it triggers, or that
// refresh's follow-up callbacks) still in flight means the switch must wait, so
// it doesn't run against a repo that's about to be swapped out.
//
// `ignore` is the event currently being processed on the UI thread — the switch
// attempt itself — which is always busy and so must not count as a reason to
// refuse itself.
func (self *TaskManager) hasBusyForegroundTaskExcept(ignore Task) bool {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	for _, task := range self.tasks {
		if task != ignore && task.isBusy() && !task.isBackground() {
			return true
		}
	}

	return false
}

// WaitUntilIdle blocks until no task is busy. Integration tests use it to wait
// for the program to finish processing before taking the next step.
func (self *TaskManager) WaitUntilIdle() {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	for self.hasBusyTask() {
		self.idleCond.Wait()
	}
}

// caller must hold self.mutex
func (self *TaskManager) hasBusyTask() bool {
	for _, task := range self.tasks {
		if task.isBusy() {
			return true
		}
	}

	return false
}

func (self *TaskManager) withMutex(f func()) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	f()

	// Wake up any goroutine blocked in WaitUntilIdle. This must not block on
	// the waiter (we hold the mutex, and the waiter may itself be trying to
	// acquire it, e.g. by creating a task, before it next waits) — which is
	// exactly what Broadcast guarantees.
	if !self.hasBusyTask() {
		self.idleCond.Broadcast()
	}
}

func (self *TaskManager) delete(taskId int) {
	self.withMutex(func() {
		delete(self.tasks, taskId)
	})
}
