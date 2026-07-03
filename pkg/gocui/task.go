package gocui

// A task represents the fact that the program is busy doing something, which
// is useful for integration tests which only want to proceed when the program
// is idle.

type Task interface {
	Done()
	Pause()
	Continue()
	// not exporting these because we don't need to
	isBusy() bool
	isBackground() bool
}

type TaskImpl struct {
	id        int
	busy      bool
	onDone    func()
	withMutex func(func())
	// Background tasks don't count towards the program being "busy" for the
	// purpose of deciding whether a repo switch is safe (see
	// TaskManager.hasBusyForegroundTaskExcept). They're the ongoing background
	// routines (auto-fetch, files refresh, external-change detection) and the
	// refreshes they trigger, whose model writes are already guarded against a
	// concurrent repo switch by the repo generation.
	background bool
}

func (self *TaskImpl) Done() {
	self.onDone()
}

func (self *TaskImpl) Pause() {
	self.withMutex(func() {
		self.busy = false
	})
}

func (self *TaskImpl) Continue() {
	self.withMutex(func() {
		self.busy = true
	})
}

func (self *TaskImpl) isBusy() bool {
	return self.busy
}

func (self *TaskImpl) isBackground() bool {
	return self.background
}

type TaskStatus int

const (
	TaskStatusBusy TaskStatus = iota
	TaskStatusPaused
	TaskStatusDone
)

type FakeTask struct {
	status TaskStatus
}

func NewFakeTask() *FakeTask {
	return &FakeTask{
		status: TaskStatusBusy,
	}
}

func (self *FakeTask) Done() {
	self.status = TaskStatusDone
}

func (self *FakeTask) Pause() {
	self.status = TaskStatusPaused
}

func (self *FakeTask) Continue() {
	self.status = TaskStatusBusy
}

func (self *FakeTask) isBusy() bool {
	return self.status == TaskStatusBusy
}

func (self *FakeTask) isBackground() bool {
	return false
}

func (self *FakeTask) Status() TaskStatus {
	return self.status
}

func (self *FakeTask) FormatStatus() string {
	return formatTaskStatus(self.status)
}

func formatTaskStatus(status TaskStatus) string {
	switch status {
	case TaskStatusBusy:
		return "busy"
	case TaskStatusPaused:
		return "paused"
	case TaskStatusDone:
		return "done"
	}
	return "unknown"
}
