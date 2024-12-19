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
			go utils.Safe(func() { self.startBackgroundFilesRefresh(refreshInterval) })
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
		err := self.backgroundFetch()
		self.gui.c.Render()
		return err
	}

	// We want an immediate fetch at startup, and since goEvery starts by
	// waiting for the interval, we need to trigger one manually first
	_ = fetch()

	userConfig := self.gui.UserConfig()
	self.goEvery(time.Second*time.Duration(userConfig.Refresher.FetchInterval), self.gui.stopChan, fetch)
}

func (self *BackgroundRoutineMgr) startBackgroundFilesRefresh(refreshInterval int) {
	self.gui.waitForIntro.Wait()

	self.goEvery(time.Second*time.Duration(refreshInterval), self.gui.stopChan, func() error {
		return self.gui.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}})
	})
}

func (self *BackgroundRoutineMgr) goEvery(interval time.Duration, stop chan struct{}, function func() error) {
	done := make(chan struct{})
	go utils.Safe(func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if self.pauseBackgroundRefreshes {
					continue
				}
				self.gui.c.OnWorker(func(gocui.Task) error {
					_ = function()
					done <- struct{}{}
					return nil
				})
				// waiting so that we don't bunch up refreshes if the refresh takes longer than the interval
				<-done
			case <-stop:
				return
			}
		}
	})
}

func (self *BackgroundRoutineMgr) backgroundFetch() (err error) {
	err = self.gui.git.Sync.FetchBackground()

	_ = self.gui.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.BRANCHES, types.COMMITS, types.REMOTES, types.TAGS}, Mode: types.ASYNC})

	return err
}
