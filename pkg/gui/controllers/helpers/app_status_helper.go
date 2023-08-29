package helpers

import (
	"time"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/status"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type AppStatusHelper struct {
	c *HelperCommon

	statusMgr func() *status.StatusManager
}

func NewAppStatusHelper(c *HelperCommon, statusMgr func() *status.StatusManager) *AppStatusHelper {
	return &AppStatusHelper{
		c:         c,
		statusMgr: statusMgr,
	}
}

func (self *AppStatusHelper) Toast(message string) {
	if self.c.RunningIntegrationTest() {
		// Don't bother showing toasts in integration tests. You can't check for
		// them anyway, and they would only slow down the test unnecessarily by
		// two seconds.
		return
	}

	self.statusMgr().AddToastStatus(message)

	self.renderAppStatus()
}

// A custom task for WithWaitingStatus calls; it wraps the original one and
// hides the status whenever the task is paused, and shows it again when
// continued.
type appStatusHelperTask struct {
	gocui.Task
	waitingStatusHandle *status.WaitingStatusHandle
}

// poor man's version of explicitly saying that struct X implements interface Y
var _ gocui.Task = appStatusHelperTask{}

func (self appStatusHelperTask) Pause() {
	self.waitingStatusHandle.Hide()
	self.Task.Pause()
}

func (self appStatusHelperTask) Continue() {
	self.Task.Continue()
	self.waitingStatusHandle.Show()
}

// withWaitingStatus wraps a function and shows a waiting status while the function is still executing
func (self *AppStatusHelper) WithWaitingStatus(message string, f func(gocui.Task) error) {
	self.c.OnWorker(func(task gocui.Task) {
		self.statusMgr().WithWaitingStatus(message, self.renderAppStatus, func(waitingStatusHandle *status.WaitingStatusHandle) {
			if err := f(appStatusHelperTask{task, waitingStatusHandle}); err != nil {
				self.c.OnUIThread(func() error {
					return self.c.Error(err)
				})
			}
		})
	})
}

func (self *AppStatusHelper) HasStatus() bool {
	return self.statusMgr().HasStatus()
}

func (self *AppStatusHelper) GetStatusString() string {
	return self.statusMgr().GetStatusString()
}

func (self *AppStatusHelper) renderAppStatus() {
	self.c.OnWorker(func(_ gocui.Task) {
		ticker := time.NewTicker(time.Millisecond * utils.LoaderAnimationInterval)
		defer ticker.Stop()
		for range ticker.C {
			appStatus := self.statusMgr().GetStatusString()
			self.c.OnUIThread(func() error {
				self.c.SetViewContent(self.c.Views().AppStatus, appStatus)
				return nil
			})

			if appStatus == "" {
				return
			}
		}
	})
}
