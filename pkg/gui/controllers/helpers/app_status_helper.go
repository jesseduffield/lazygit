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

// WithWaitingStatus wraps a function and shows a waiting status while the function is still executing
func (self *AppStatusHelper) WithWaitingStatus(message string, f func(gocui.Task) error) {
	self.c.OnWorker(func(task gocui.Task) error {
		return self.WithWaitingStatusImpl(message, f, task)
	})
}

// WithWaitingStatusImpl is WithWaitingStatus for callers that already run on a
// goroutine of their own (e.g. the auto-fetch poller) rather than wanting the
// work dispatched to a worker. task is used to hide the status while the task
// is paused; it may be nil for callers whose f ignores its task.
func (self *AppStatusHelper) WithWaitingStatusImpl(message string, f func(gocui.Task) error, task gocui.Task) error {
	// A waiting status means lazygit is driving a git operation itself (often
	// one that internally runs a rebase and continues it). Pause the background
	// routines for its duration so they don't refresh from an intermediate
	// state and reveal, say, the half-finished history of a reword.
	self.c.PauseBackgroundRefreshes(true)
	defer self.c.PauseBackgroundRefreshes(false)

	return self.statusMgr().WithWaitingStatus(message, self.renderAppStatus, func(waitingStatusHandle *status.WaitingStatusHandle) error {
		return f(appStatusHelperTask{task, waitingStatusHandle})
	})
}

// WithWaitingStatusBlockingInput is like WithWaitingStatus, but it also blocks
// keyboard input for the whole duration of the operation: keys the user presses
// while it runs are buffered and replayed against the post-operation state (see
// gocui.BeginBlockingEvents). Use it for operations that manipulate an
// in-progress rebase or otherwise rewrite commits, where a racing keypress
// would target the wrong commit or todo.
//
// Must be called on the UI thread: the block is begun synchronously here, before
// the operation is dispatched to a worker, so no keypress can slip through in
// between.
func (self *AppStatusHelper) WithWaitingStatusBlockingInput(message string, f func(gocui.Task) error) {
	self.c.GocuiGui().BeginBlockingEvents()
	// Hide the rebasing-mode indicator (and its reset button) while we drive the
	// rebase ourselves; it reflects the transient on-disk state and would
	// otherwise flash on for the duration of the operation.
	self.modeHelper.SetSuppressRebasingMode(true)
	self.c.OnWorker(func(task gocui.Task) error {
		// End the block and restore the mode indicator once the operation and its
		// refresh have applied their UI updates: OnUIThread queues this after the
		// refresh's model bounces and Then (which RefreshFromWorker has already
		// enqueued by the time f returns), so the replayed keys act on the
		// refreshed state and any resulting rebase state shows correctly.
		defer self.c.OnUIThread(func() error {
			self.modeHelper.SetSuppressRebasingMode(false)
			return self.c.GocuiGui().EndBlockingEvents()
		})
		return self.WithWaitingStatusImpl(message, f, task)
	})
}

func (self *AppStatusHelper) HasStatus() bool {
	return self.statusMgr().HasStatus()
}

func (self *AppStatusHelper) GetStatusString() string {
	appStatus, _ := self.statusMgr().GetStatusString(self.c.UserConfig())
	return appStatus
}

// renderAppStatus ensures the render loop that keeps the app-status view up to
// date is running. There is one loop for the whole status stack, no matter how
// many statuses are showing: it draws whatever the top status currently is,
// and exits after drawing a final empty frame once the last status is removed.
//
// The loop always runs as a background task, regardless of what kind of
// operation owns a status: rendering runs no git commands, so it must never
// count towards lazygit being busy — otherwise it would block repo switching
// for as long as anything is showing (e.g. for the whole duration of a hung
// background fetch, or of a toast fading). A foreground operation's busy-ness
// is carried by its own worker task, not by the renderer.
func (self *AppStatusHelper) renderAppStatus() {
	if !self.statusMgr().ClaimRenderLoop() {
		return
	}

	self.c.OnWorkerBackground(func(_ gocui.Task) error {
		ticker := time.NewTicker(time.Millisecond * time.Duration(self.c.UserConfig().Gui.Spinner.Rate))
		defer ticker.Stop()
		prevAppStatus := ""
		for range ticker.C {
			appStatus, color := self.statusMgr().GetStatusString(self.c.UserConfig())

			update := self.c.OnUIThreadContentOnlyBackground
			if utils.StringWidth(appStatus) != utils.StringWidth(prevAppStatus) {
				// Need a full layout whenever the width of the status string changes. This can't
				// happen during normal spinning because we validate that all spinner frames have
				// the same width, so typically this will only be triggered at the beginning and end
				// of a status, or if the status string changes midway for some reason.
				update = self.c.OnUIThreadBackground
			}
			update(func() error {
				self.c.Views().AppStatus.FgColor = color
				self.c.SetViewContent(self.c.Views().AppStatus, appStatus)
				return nil
			})
			prevAppStatus = appStatus

			// Checked after rendering, so that the frame which clears the view
			// has already been drawn when we exit.
			if self.statusMgr().ReleaseRenderLoopIfEmpty() {
				break
			}
		}
		return nil
	})
}
