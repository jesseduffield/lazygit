package helpers

import (
	"time"

	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/status"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
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

	self.renderAppStatus(false)
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

// WithWaitingStatus wraps a function and shows a waiting status while the function is still executing
func (self *AppStatusHelper) WithWaitingStatus(message string, f func(gocui.Task) error) {
	self.c.OnWorker(func(task gocui.Task) error {
		return self.WithWaitingStatusImpl(message, f, task, false)
	})
}

// background reports whether this waiting status belongs to a background routine
// (the auto-fetch poller); when it does, the spinner it drives must not count
// towards lazygit being busy, or it'd block repo switches while a fetch runs.
func (self *AppStatusHelper) WithWaitingStatusImpl(message string, f func(gocui.Task) error, task gocui.Task, background bool) error {
	// A waiting status means lazygit is driving a git operation itself (often
	// one that internally runs a rebase and continues it). Pause the background
	// routines for its duration so they don't refresh from an intermediate
	// state and reveal, say, the half-finished history of a reword.
	self.c.PauseBackgroundRefreshes(true)
	defer self.c.PauseBackgroundRefreshes(false)

	return self.statusMgr().WithWaitingStatus(message, func() { self.renderAppStatus(background) }, func(waitingStatusHandle *status.WaitingStatusHandle) error {
		return f(appStatusHelperTask{task, waitingStatusHandle})
	})
}

func (self *AppStatusHelper) WithWaitingStatusSync(message string, f func() error) error {
	self.c.PauseBackgroundRefreshes(true)
	defer self.c.PauseBackgroundRefreshes(false)

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
	appStatus, _ := self.statusMgr().GetStatusString(self.c.UserConfig())
	return appStatus
}

func (self *AppStatusHelper) renderAppStatus(background bool) {
	// A background waiting status (auto-fetch) must not count towards lazygit
	// being busy, so its spinner worker and per-frame UI updates go through the
	// background variants.
	onWorker := self.c.OnWorker
	onUIThread := self.c.OnUIThread
	onUIThreadContentOnly := self.c.OnUIThreadContentOnly
	if background {
		onWorker = self.c.OnWorkerBackground
		onUIThread = self.c.OnUIThreadBackground
		onUIThreadContentOnly = self.c.OnUIThreadContentOnlyBackground
	}

	onWorker(func(_ gocui.Task) error {
		ticker := time.NewTicker(time.Millisecond * time.Duration(self.c.UserConfig().Gui.Spinner.Rate))
		defer ticker.Stop()
		prevAppStatus := ""
		for range ticker.C {
			appStatus, color := self.statusMgr().GetStatusString(self.c.UserConfig())

			update := onUIThreadContentOnly
			if utils.StringWidth(appStatus) != utils.StringWidth(prevAppStatus) {
				// Need a full layout whenever the width of the status string changes. This can't
				// happen during normal spinning because we validate that all spinner frames have
				// the same width, so typically this will only be triggered at the beginning and end
				// of a status, or if the status string changes midway for some reason.
				update = onUIThread
			}
			update(func() error {
				self.c.Views().AppStatus.FgColor = color
				self.c.SetViewContent(self.c.Views().AppStatus, appStatus)
				return nil
			})
			prevAppStatus = appStatus

			if appStatus == "" {
				break
			}
		}
		return nil
	})
}

func (self *AppStatusHelper) renderAppStatusSync(stop chan struct{}) {
	go func() {
		ticker := time.NewTicker(time.Millisecond * time.Duration(self.c.UserConfig().Gui.Spinner.Rate))
		defer ticker.Stop()

		// Write the status into the view before the first layout below, so that
		// layout (which sizes the bottom line based on the actual content of the
		// AppStatus view) leaves room for it and it shows right away. The ticker
		// only updates the spinner frame using ForceFlushViewsContentOnly, so this
		// doesn't re-layout.
		self.setAppStatusContent()

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
				self.setAppStatusContent()
				// Redraw all views of the bottom line:
				bottomLineViews := []*gocui.View{
					self.c.Views().AppStatus, self.c.Views().Options, self.c.Views().Information,
					self.c.Views().StatusSpacer1, self.c.Views().StatusSpacer2,
				}
				_ = self.c.GocuiGui().ForceFlushViewsContentOnly(bottomLineViews)
			case <-stop:
				// Clear the status from the view and re-layout, otherwise the
				// stale content would keep layout reserving room for it forever.
				// The UI thread is free again at this point, so we go through
				// OnUIThread like the async renderAppStatus does.
				self.c.OnUIThread(func() error {
					self.c.SetViewContent(self.c.Views().AppStatus, "")
					return nil
				})
				break outer
			}
		}
	}()
}

func (self *AppStatusHelper) setAppStatusContent() {
	appStatus, color := self.statusMgr().GetStatusString(self.c.UserConfig())
	self.c.Views().AppStatus.FgColor = color
	self.c.SetViewContent(self.c.Views().AppStatus, appStatus)
}
