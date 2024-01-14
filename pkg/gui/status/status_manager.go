package status

import (
	"time"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
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
	self.id = self.statusManager.addStatus(self.message, "waiting", types.ToastKindStatus)
	self.renderFunc()
}

func (self *WaitingStatusHandle) Hide() {
	self.statusManager.removeStatus(self.id)
}

type appStatus struct {
	message    string
	statusType string
	color      gocui.Attribute
	id         int
}

func NewStatusManager() *StatusManager {
	return &StatusManager{}
}

func (self *StatusManager) WithWaitingStatus(message string, renderFunc func(), f func(*WaitingStatusHandle)) {
	handle := &WaitingStatusHandle{statusManager: self, message: message, renderFunc: renderFunc, id: -1}
	handle.Show()

	f(handle)

	handle.Hide()
}

func (self *StatusManager) AddToastStatus(message string, kind types.ToastKind) int {
	id := self.addStatus(message, "toast", kind)

	go func() {
		delay := lo.Ternary(kind == types.ToastKindError, time.Second*4, time.Second*2)
		time.Sleep(delay)

		self.removeStatus(id)
	}()

	return id
}

func (self *StatusManager) GetStatusString() (string, gocui.Attribute) {
	if len(self.statuses) == 0 {
		return "", gocui.ColorDefault
	}
	topStatus := self.statuses[0]
	if topStatus.statusType == "waiting" {
		return topStatus.message + " " + utils.Loader(time.Now()), topStatus.color
	}
	return topStatus.message, topStatus.color
}

func (self *StatusManager) HasStatus() bool {
	return len(self.statuses) > 0
}

func (self *StatusManager) addStatus(message string, statusType string, kind types.ToastKind) int {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	self.nextId++
	id := self.nextId

	color := gocui.ColorCyan
	if kind == types.ToastKindError {
		color = gocui.ColorRed
	}

	newStatus := appStatus{
		message:    message,
		statusType: statusType,
		color:      color,
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
