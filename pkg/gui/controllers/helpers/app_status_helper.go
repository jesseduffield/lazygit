package helpers

import (
	"time"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/status"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type AppStatusHelper struct {
	c *HelperCommon

	statusMgr  func() *status.StatusManager
	modeHelper *ModeHelper
}

func NewAppStatusHelper(c *HelperCommon, statusMgr func() *status.StatusManager, modeHelper *ModeHelper) *AppStatusHelper {
	return &AppStatusHelper{
		c:          c,
		statusMgr:  statusMgr,
		modeHelper: modeHelper,
	}
}

func (self *AppStatusHelper) Toast(message string, kind types.ToastKind) {
	if self.c.RunningIntegrationTest() {
		// Don't bother showing toasts in integration tests. You can't check for
		// them anyway, and they would only slow down the test unnecessarily by
		// two seconds.
		return
	}

	self.statusMgr().AddToastStatus(message, kind)

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
	self.c.OnWorker(func(task gocui.Task) error {
		return self.statusMgr().WithWaitingStatus(message, self.renderAppStatus, func(waitingStatusHandle *status.WaitingStatusHandle) error {
			return f(appStatusHelperTask{task, waitingStatusHandle})
		})
	})
}

func (self *AppStatusHelper) WithWaitingStatusSync(message string, f func() error) error {
	return self.statusMgr().WithWaitingStatus(message, func() {}, func(*status.WaitingStatusHandle) error {
		stop := make(chan struct{})
		defer func() { close(stop) }()
		self.renderAppStatusSync(stop)

		return f()
	})
}

func (self *AppStatusHelper) HasStatus() bool {
	return self.statusMgr().HasStatus()
}

func (self *AppStatusHelper) GetStatusString() string {
	appStatus, _ := self.statusMgr().GetStatusString(self.c.UserConfig)
	return appStatus
}

func (self *AppStatusHelper) renderAppStatus() {
	self.c.OnWorker(func(_ gocui.Task) error {
		ticker := time.NewTicker(time.Millisecond * time.Duration(self.c.UserConfig.Gui.Spinner.Rate))
		defer ticker.Stop()
		for range ticker.C {
			appStatus, color := self.statusMgr().GetStatusString(self.c.UserConfig)
			self.c.Views().AppStatus.FgColor = color
			self.c.OnUIThread(func() error {
				self.c.SetViewContent(self.c.Views().AppStatus, appStatus)
				return nil
			})

			if appStatus == "" {
				break
			}
		}
		return nil
	})
}

func (self *AppStatusHelper) renderAppStatusSync(stop chan struct{}) {
	go func() {
		ticker := time.NewTicker(time.Millisecond * 50)
		defer ticker.Stop()

		// Forcing a re-layout and redraw after we added the waiting status;
		// this is needed in case the gui.showBottomLine config is set to false,
		// to make sure the bottom line appears. It's also useful for redrawing
		// once after each of several consecutive keypresses, e.g. pressing
		// ctrl-j to move a commit down several steps.
		_ = self.c.GocuiGui().ForceLayoutAndRedraw()

		self.modeHelper.SetSuppressRebasingMode(true)
		defer func() { self.modeHelper.SetSuppressRebasingMode(false) }()

	outer:
		for {
			select {
			case <-ticker.C:
				appStatus, color := self.statusMgr().GetStatusString(self.c.UserConfig)
				self.c.Views().AppStatus.FgColor = color
				self.c.SetViewContent(self.c.Views().AppStatus, appStatus)
				// Redraw all views of the bottom line:
				bottomLineViews := []*gocui.View{
					self.c.Views().AppStatus, self.c.Views().Options, self.c.Views().Information,
					self.c.Views().StatusSpacer1, self.c.Views().StatusSpacer2,
				}
				_ = self.c.GocuiGui().ForceRedrawViews(bottomLineViews...)
			case <-stop:
				break outer
			}
		}
	}()
}
