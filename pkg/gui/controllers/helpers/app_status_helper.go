package helpers

import (
	"time"

	"github.com/jesseduffield/lazygit/pkg/gui/status"
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
	self.statusMgr().AddToastStatus(message)

	self.renderAppStatus()
}

// withWaitingStatus wraps a function and shows a waiting status while the function is still executing
func (self *AppStatusHelper) WithWaitingStatus(message string, f func() error) {
	self.c.OnWorker(func() {
		self.statusMgr().WithWaitingStatus(message, func() {
			self.renderAppStatus()

			if err := f(); err != nil {
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
	self.c.OnWorker(func() {
		ticker := time.NewTicker(time.Millisecond * 50)
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
