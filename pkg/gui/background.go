package gui

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type BackgroundRoutineMgr struct {
	gui *Gui

	// When this is greater than zero, the background routines (e.g. file refresh)
	// skip their work. We pause them while the gui is suspended (e.g. for a
	// subprocess) and while lazygit is itself driving a git operation that would
	// otherwise be caught mid-flight (see the waiting-status helpers). It's a
	// count rather than a bool because these pause scopes can overlap.
	pauseRefreshesCount atomic.Int32

	// a channel to trigger an immediate background fetch; we use this when switching repos
	triggerFetch chan struct{}
}

func (self *BackgroundRoutineMgr) PauseBackgroundRefreshes(pause bool) {
	if pause {
		self.pauseRefreshesCount.Add(1)
	} else {
		self.pauseRefreshesCount.Add(-1)
	}
}

func (self *BackgroundRoutineMgr) backgroundRefreshesPaused() bool {
	return self.pauseRefreshesCount.Load() > 0
}

func (self *BackgroundRoutineMgr) startBackgroundRoutines() {
	userConfig := self.gui.UserConfig()

	if userConfig.Git.AutoFetch {
		fetchInterval := userConfig.Refresher.FetchInterval
		if fetchInterval > 0 {
			// The channel must be created here, on the UI thread and before
			// the fetch goroutine spawns, so that triggerImmediateFetch (also
			// running on the UI thread) can read the field without racing the
			// write. See triggerImmediateFetch for why it is buffered.
			self.triggerFetch = make(chan struct{}, 1)
			go utils.Safe(self.startBackgroundFetch)
		} else {
			self.gui.c.Log.Errorf(
				"Value of config option 'refresher.fetchInterval' (%d) is invalid, disabling auto-fetch",
				fetchInterval)
		}
	}

	if userConfig.Git.AutoRefresh {
		refreshInterval := userConfig.Refresher.RefreshInterval
		if refreshInterval > 0 {
			go utils.Safe(self.startBackgroundFilesRefresh)
		} else {
			self.gui.c.Log.Errorf(
				"Value of config option 'refresher.refreshInterval' (%d) is invalid, disabling auto-refresh",
				refreshInterval)
		}
	}

	if userConfig.Git.AutoDetectExternalChanges {
		interval := userConfig.Refresher.ExternalChangeCheckInterval
		if interval > 0 {
			go utils.Safe(self.startBackgroundExternalChangeDetection)
		} else {
			self.gui.c.Log.Errorf(
				"Value of config option 'refresher.externalChangeCheckInterval' (%d) is invalid, disabling external change detection",
				interval)
		}
	}

	if self.gui.Config.GetDebug() {
		self.goEvery(time.Second*time.Duration(10), self.gui.stopChan, nil, func(_ bool) error {
			formatBytes := func(b uint64) string {
				const unit = 1000
				if b < unit {
					return fmt.Sprintf("%d B", b)
				}
				div, exp := uint64(unit), 0
				for n := b / unit; n >= unit; n /= unit {
					div *= unit
					exp++
				}
				return fmt.Sprintf("%.1f %cB",
					float64(b)/float64(div), "kMGTPE"[exp])
			}

			m := runtime.MemStats{}
			runtime.ReadMemStats(&m)
			self.gui.c.Log.Infof("Heap memory in use: %s", formatBytes(m.HeapAlloc))
			return nil
		})
	}
}

func (self *BackgroundRoutineMgr) startBackgroundFetch() {
	self.gui.waitForIntro.Wait()

	fetch := func(firstTimeOrRetriggered bool) error {
		// Capture what the fetch needs from the gui's per-repo state in a
		// single UI-thread hop: gui.git, gui.helpers and gui.State are all
		// replaced on a repo switch (which runs on the UI thread), so reading
		// them from this background goroutine would race the reassignment.
		// Capturing them together also ties the fetch, the post-fetch
		// refresh's generation baseline, and the recorded fetch time to the
		// same repo.
		var git *commands.GitCommand
		var appStatusHelper *helpers.AppStatusHelper
		var branchesHelper *helpers.BranchesHelper
		var fetchGeneration int
		if err := self.gui.g.OnUIThreadAndWaitBackground(func() error {
			git = self.gui.git
			appStatusHelper = self.gui.helpers.AppStatus
			branchesHelper = self.gui.helpers.BranchesHelper
			fetchGeneration = self.gui.c.State().GetRepoGeneration()
			self.gui.State.LastBackgroundFetchTime = time.Now()
			return nil
		}); err != nil {
			return err
		}

		if self.gui.UserConfig().Gui.ShowBottomLine || firstTimeOrRetriggered {
			return appStatusHelper.WithWaitingStatusImpl(self.gui.Tr.FetchingStatus, func(gocui.Task) error {
				return self.backgroundFetch(git, branchesHelper, fetchGeneration)
			}, nil)
		}

		return self.backgroundFetch(git, branchesHelper, fetchGeneration)
	}

	// We want an immediate fetch at startup, and since goEvery starts by
	// waiting for the interval, we need to trigger one manually first
	_ = fetch(true)

	userConfig := self.gui.UserConfig()
	self.goEvery(userConfig.Refresher.FetchIntervalDuration(), self.gui.stopChan, self.triggerFetch, fetch)
}

func (self *BackgroundRoutineMgr) startBackgroundFilesRefresh() {
	self.gui.waitForIntro.Wait()

	userConfig := self.gui.UserConfig()
	self.goEvery(userConfig.Refresher.RefreshIntervalDuration(), self.gui.stopChan, nil, func(_ bool) error {
		self.gui.c.RefreshFromWorker(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}, Background: true})
		return nil
	})
}

func (self *BackgroundRoutineMgr) startBackgroundExternalChangeDetection() {
	self.gui.waitForIntro.Wait()

	// We don't seed the snapshot here. The startup refresh captures one on
	// entry (like every refs-touching refresh), and until one has been
	// captured RefsSnapshotChangedSince treats the empty baseline as
	// "unchanged", so we never fire a spurious refresh before a baseline
	// exists — no need to depend on the timing of that startup refresh.

	userConfig := self.gui.UserConfig()
	self.goEvery(
		userConfig.Refresher.ExternalChangeCheckIntervalDuration(),
		self.gui.stopChan,
		nil,
		func(_ bool) error {
			self.checkForExternalChanges()
			return nil
		},
	)
}

func (self *BackgroundRoutineMgr) checkForExternalChanges() {
	// Capture the per-repo objects in a UI-thread hop, like the background
	// fetch does: gui.git and gui.helpers are replaced on a repo switch, so
	// reading them from this background goroutine would race the reassignment.
	var git *commands.GitCommand
	var refreshHelper *helpers.RefreshHelper
	if err := self.gui.g.OnUIThreadAndWaitBackground(func() error {
		git = self.gui.git
		refreshHelper = self.gui.helpers.Refresh
		return nil
	}); err != nil {
		return
	}

	current, err := git.Status.RefsSnapshot()
	if err != nil {
		// Transient error (e.g. git process couldn't start). Don't update the
		// stored snapshot; we'll retry next tick.
		self.gui.c.Log.Warnf("RefsSnapshot failed: %v", err)
		return
	}

	if !refreshHelper.RefsSnapshotChangedSince(current) {
		return
	}

	// goEvery checks the pause count before starting us, but a git operation
	// may have begun (and paused refreshes) after that check, while we were
	// reading the snapshot above. In that case the change we detected is the
	// operation's own intermediate state, so back off: the operation will
	// refresh and re-snapshot when it finishes, and if the change was really
	// external we'll catch it on the next tick after the pause lifts. We don't
	// update the stored snapshot, so nothing is swallowed.
	if self.backgroundRefreshesPaused() {
		return
	}

	// No need to update the stored snapshot here; Refresh does that.
	self.gui.c.Log.Info("External ref change detected — refreshing")
	self.gui.c.RefreshFromWorker(types.RefreshOptions{Background: true})
}

// Runs function every interval until stop is closed. A send on retrigger (if
// non-nil) runs the callback immediately and restarts the interval.
func (self *BackgroundRoutineMgr) goEvery(interval time.Duration, stop, retrigger chan struct{}, function func(bool) error) {
	done := make(chan struct{})
	go utils.Safe(func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		doit := func(retriggered bool) {
			if self.backgroundRefreshesPaused() {
				return
			}
			// OnWorkerBackground, not OnWorker: these routines and the refreshes
			// they trigger must not count towards lazygit being busy, or they'd
			// spuriously block a repo switch every time one happens to be running.
			self.gui.c.OnWorkerBackground(func(gocui.Task) error {
				_ = function(retriggered)
				done <- struct{}{}
				return nil
			})
			// waiting so that we don't bunch up refreshes if the refresh takes longer than the
			// interval, or if a retrigger comes in while we're still processing a timer-based one
			// (or vice versa)
			<-done
		}
		for {
			select {
			case <-ticker.C:
				doit(false)
			case <-retrigger:
				ticker.Reset(interval)
				doit(true)
			case <-stop:
				return
			}
		}
	})
}

// The parameters are captured by the caller before the fetch starts, not read
// here after it: the fetch is a network call during which the user may switch
// repos, and the post-fetch refresh needs to be able to tell (see
// PostFetchRefresh).
func (self *BackgroundRoutineMgr) backgroundFetch(git *commands.GitCommand, branchesHelper *helpers.BranchesHelper, fetchGeneration int) error {
	err := git.Sync.FetchBackground()

	return branchesHelper.PostFetchRefresh(err, true, fetchGeneration)
}

func (self *BackgroundRoutineMgr) triggerImmediateFetch() {
	if self.triggerFetch != nil {
		// This runs on the UI thread, which must never block waiting for a
		// background routine; in particular, the goEvery loop only receives
		// between callbacks, and an in-flight fetch can itself be waiting for
		// the UI thread to perform its post-fetch refresh, so a blocking send
		// here would deadlock. The channel has a buffer of one, so the trigger
		// is latched even when the loop isn't currently receiving; if one is
		// already pending, the two coalesce.
		select {
		case self.triggerFetch <- struct{}{}:
		default:
		}
	}
}
