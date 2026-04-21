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

func (self *TaskManager) NewTask() *TaskImpl {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	self.nextId++
	taskId := self.nextId

	onDone := func() { self.delete(taskId) }
	task := &TaskImpl{id: taskId, busy: true, onDone: onDone, withMutex: self.withMutex}
	self.tasks[taskId] = task

	return task
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
