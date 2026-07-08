package gocui

import "sync"

// Tracks whether the program is busy (i.e. either something is happening on
// the main goroutine or a worker goroutine). Used by integration tests
// to wait until the program is idle before progressing.
type TaskManager struct {
	// each of these listeners will be notified when the program goes from busy to idle
	idleListeners []chan struct{}
	tasks         map[int]Task
	// auto-incrementing id for new tasks
	nextId int

	mutex sync.Mutex
}

func newTaskManager() *TaskManager {
	return &TaskManager{
		tasks:         make(map[int]Task),
		idleListeners: []chan struct{}{},
	}
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

func (self *TaskManager) addIdleListener(c chan struct{}) {
	self.idleListeners = append(self.idleListeners, c)
}

func (self *TaskManager) withMutex(f func()) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	f()

	// Check if all tasks are done
	for _, task := range self.tasks {
		if task.isBusy() {
			return
		}
	}

	// If we get here, all tasks are done, so
	// notify listeners that the program is idle
	for _, listener := range self.idleListeners {
		listener <- struct{}{}
	}
}

func (self *TaskManager) delete(taskId int) {
	self.withMutex(func() {
		delete(self.tasks, taskId)
	})
}
