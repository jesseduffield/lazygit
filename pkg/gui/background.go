package gui

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/jesseduffield/lazygit/pkg/gocui"
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
		self.goEvery(time.Second*time.Duration(10), self.gui.stopChan, func(_ bool) error {
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
		// Do this on the UI thread so that we don't have to deal with synchronization around the
		// access of the repo state.
		self.gui.onUIThread(func() error {
			// There's a race here, where we might be recording the time stamp for a different repo
			// than where the fetch actually ran. It's not very likely though, and not harmful if it
			// does happen; guarding against it would be more effort than it's worth.
			self.gui.State.LastBackgroundFetchTime = time.Now()
			return nil
		})

		if self.gui.UserConfig().Gui.ShowBottomLine || firstTimeOrRetriggered {
			return self.gui.helpers.AppStatus.WithWaitingStatusImpl(self.gui.Tr.FetchingStatus, func(gocui.Task) error {
				return self.backgroundFetch()
			}, nil, true)
		}

		return self.backgroundFetch()
	}

	// We want an immediate fetch at startup, and since goEvery starts by
	// waiting for the interval, we need to trigger one manually first
	_ = fetch(true)

	userConfig := self.gui.UserConfig()
	self.triggerFetch = self.goEvery(userConfig.Refresher.FetchIntervalDuration(), self.gui.stopChan, fetch)
}

func (self *BackgroundRoutineMgr) startBackgroundFilesRefresh() {
	self.gui.waitForIntro.Wait()

	userConfig := self.gui.UserConfig()
	self.goEvery(userConfig.Refresher.RefreshIntervalDuration(), self.gui.stopChan, func(_ bool) error {
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
		func(_ bool) error {
			self.checkForExternalChanges()
			return nil
		},
	)
}

func (self *BackgroundRoutineMgr) checkForExternalChanges() {
	current, err := self.gui.git.Status.RefsSnapshot()
	if err != nil {
		// Transient error (e.g. git process couldn't start). Don't update the
		// stored snapshot; we'll retry next tick.
		self.gui.c.Log.Warnf("RefsSnapshot failed: %v", err)
		return
	}

	if !self.gui.helpers.Refresh.RefsSnapshotChangedSince(current) {
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

// returns a channel that can be used to trigger the callback immediately
func (self *BackgroundRoutineMgr) goEvery(interval time.Duration, stop chan struct{}, function func(bool) error) chan struct{} {
	done := make(chan struct{})
	// Buffered so that a retrigger arriving while the callback is running is
	// latched rather than lost: the loop below doesn't receive again until the
	// callback has finished, and the callback (a fetch) may be for the wrong
	// repo if the retrigger came from a repo switch.
	retrigger := make(chan struct{}, 1)
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
	return retrigger
}

func (self *BackgroundRoutineMgr) backgroundFetch() (err error) {
	err = self.gui.git.Sync.FetchBackground()

	return self.gui.helpers.BranchesHelper.PostFetchRefresh(err, true)
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
