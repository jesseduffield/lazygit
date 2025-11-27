package gui

import (
	"fmt"
	"runtime"
	"time"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type BackgroundRoutineMgr struct {
	gui *Gui

	// if we've suspended the gui (e.g. because we've switched to a subprocess)
	// we typically want to pause some things that are running like background
	// file refreshes
	pauseBackgroundRefreshes bool

	// a channel to trigger an immediate background fetch; we use this when switching repos
	triggerFetch chan struct{}
}

func (self *BackgroundRoutineMgr) PauseBackgroundRefreshes(pause bool) {
	self.pauseBackgroundRefreshes = pause
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

	if self.gui.Config.GetDebug() {
		self.goEvery(time.Second*time.Duration(10), self.gui.stopChan, func() error {
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

	fetch := func() error {
		// Do this on the UI thread so that we don't have to deal with synchronization around the
		// access of the repo state.
		self.gui.onUIThread(func() error {
			// There's a race here, where we might be recording the time stamp for a different repo
			// than where the fetch actually ran. It's not very likely though, and not harmful if it
			// does happen; guarding against it would be more effort than it's worth.
			self.gui.State.LastBackgroundFetchTime = time.Now()
			return nil
		})

		return self.gui.helpers.AppStatus.WithWaitingStatusImpl(self.gui.Tr.FetchingStatus, func(gocui.Task) error {
			return self.backgroundFetch()
		}, nil)
	}

	// We want an immediate fetch at startup, and since goEvery starts by
	// waiting for the interval, we need to trigger one manually first
	_ = fetch()

	userConfig := self.gui.UserConfig()
	self.triggerFetch = self.goEvery(userConfig.Refresher.FetchIntervalDuration(), self.gui.stopChan, fetch)
}

func (self *BackgroundRoutineMgr) startBackgroundFilesRefresh() {
	self.gui.waitForIntro.Wait()

	userConfig := self.gui.UserConfig()
	self.goEvery(userConfig.Refresher.RefreshIntervalDuration(), self.gui.stopChan, func() error {
		self.gui.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}})
		return nil
	})
}

// returns a channel that can be used to trigger the callback immediately
func (self *BackgroundRoutineMgr) goEvery(interval time.Duration, stop chan struct{}, function func() error) chan struct{} {
	done := make(chan struct{})
	retrigger := make(chan struct{})
	go utils.Safe(func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		doit := func() {
			if self.pauseBackgroundRefreshes {
				return
			}
			self.gui.c.OnWorker(func(gocui.Task) error {
				_ = function()
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
				doit()
			case <-retrigger:
				ticker.Reset(interval)
				doit()
			case <-stop:
				return
			}
		}
	})
	return retrigger
}

func (self *BackgroundRoutineMgr) backgroundFetch() (err error) {
	err = self.gui.git.Sync.FetchBackground()

	self.gui.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.BRANCHES, types.COMMITS, types.REMOTES, types.TAGS}, Mode: types.SYNC})

	if err == nil {
		err = self.gui.helpers.BranchesHelper.AutoForwardBranches()
	}

	return err
}

func (self *BackgroundRoutineMgr) triggerImmediateFetch() {
	if self.triggerFetch != nil {
		self.triggerFetch <- struct{}{}
	}
}
