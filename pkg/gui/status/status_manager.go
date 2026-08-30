package status

import (
	"time"

	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
	"github.com/sasha-s/go-deadlock"
)

// StatusManager's job is to handle queuing of loading states and toast notifications
// that you see at the bottom left of the screen.
type StatusManager struct {
	statuses []appStatus
	nextId   int
	mutex    deadlock.Mutex

	// Whether a render loop is currently drawing the statuses. Guarded by
	// mutex, so that claiming and releasing the loop stay atomic with the
	// changes to statuses; see ClaimRenderLoop and ReleaseRenderLoopIfEmpty.
	renderLoopRunning bool
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

func (self *StatusManager) WithWaitingStatus(message string, renderFunc func(), f func(*WaitingStatusHandle) error) error {
	handle := &WaitingStatusHandle{statusManager: self, message: message, renderFunc: renderFunc, id: -1}
	handle.Show()
	defer handle.Hide()

	return f(handle)
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

func (self *StatusManager) GetStatusString(userConfig *config.UserConfig) (string, gocui.Attribute) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if len(self.statuses) == 0 {
		return "", gocui.ColorDefault
	}
	topStatus := self.statuses[0]
	if topStatus.statusType == "waiting" {
		return topStatus.message + " " + presentation.Loader(time.Now(), userConfig.Gui.Spinner), topStatus.color
	}
	return topStatus.message, topStatus.color
}

func (self *StatusManager) HasStatus() bool {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	return len(self.statuses) > 0
}

// ClaimRenderLoop is called by whoever just added a status; it reports whether
// they must start the render loop. When it returns false, a loop is already
// running and will pick the new status up on its next tick.
func (self *StatusManager) ClaimRenderLoop() bool {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if self.renderLoopRunning {
		return false
	}

	self.renderLoopRunning = true
	return true
}

// ReleaseRenderLoopIfEmpty is called by the render loop after each frame it
// draws; a true result releases the loop's claim and tells it to exit, because
// there are no statuses left to draw. The emptiness check and the release are
// atomic with respect to ClaimRenderLoop, so a status added around this moment
// either sees the still-running loop or starts a fresh one — it can't end up
// unrendered.
func (self *StatusManager) ReleaseRenderLoopIfEmpty() bool {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if len(self.statuses) > 0 {
		return false
	}

	self.renderLoopRunning = false
	return true
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
