package status

import (
	"time"

	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
	"github.com/sasha-s/go-deadlock"
)

// StatusManager's job is to handle queuing of loading states and toast notifications
// that you see at the bottom left of the screen.
type StatusManager struct {
	statuses []appStatus
	nextId   int
	mutex    deadlock.Mutex
}

// Can be used to manipulate a waiting status while it is running (e.g. pause
// and resume it)
type WaitingStatusHandle struct {
	statusManager *StatusManager
	message       string
	renderFunc    func()
	id            int
}

func (self *WaitingStatusHandle) Show() {
	self.id = self.statusManager.addStatus(self.message, "waiting")
	self.renderFunc()
}

func (self *WaitingStatusHandle) Hide() {
	self.statusManager.removeStatus(self.id)
}

type appStatus struct {
	message    string
	statusType string
	id         int
}

func NewStatusManager() *StatusManager {
	return &StatusManager{}
}

func (self *StatusManager) WithWaitingStatus(message string, renderFunc func(), f func()) {
	handle := &WaitingStatusHandle{statusManager: self, message: message, renderFunc: renderFunc, id: -1}
	handle.Show()

	f()

	handle.Hide()
}

func (self *StatusManager) AddToastStatus(message string) int {
	id := self.addStatus(message, "toast")

	go func() {
		time.Sleep(time.Second * 2)

		self.removeStatus(id)
	}()

	return id
}

func (self *StatusManager) GetStatusString() string {
	if len(self.statuses) == 0 {
		return ""
	}
	topStatus := self.statuses[0]
	if topStatus.statusType == "waiting" {
		return topStatus.message + " " + utils.Loader()
	}
	return topStatus.message
}

func (self *StatusManager) HasStatus() bool {
	return len(self.statuses) > 0
}

func (self *StatusManager) addStatus(message string, statusType string) int {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	self.nextId++
	id := self.nextId

	newStatus := appStatus{
		message:    message,
		statusType: statusType,
		id:         id,
	}
	self.statuses = append([]appStatus{newStatus}, self.statuses...)

	return id
}

func (self *StatusManager) removeStatus(id int) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	self.statuses = lo.Filter(self.statuses, func(status appStatus, _ int) bool {
		return status.id != id
	})
}
