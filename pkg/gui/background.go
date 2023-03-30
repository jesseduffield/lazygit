package gui

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) startBackgroundRoutines() {
	userConfig := gui.UserConfig

	if userConfig.Git.AutoFetch {
		fetchInterval := userConfig.Refresher.FetchInterval
		if fetchInterval > 0 {
			go utils.Safe(gui.startBackgroundFetch)
		} else {
			gui.c.Log.Errorf(
				"Value of config option 'refresher.fetchInterval' (%d) is invalid, disabling auto-fetch",
				fetchInterval)
		}
	}

	if userConfig.Git.AutoRefresh {
		refreshInterval := userConfig.Refresher.RefreshInterval
		if refreshInterval > 0 {
			gui.goEvery(time.Second*time.Duration(refreshInterval), gui.stopChan, gui.refreshFilesAndSubmodules)
		} else {
			gui.c.Log.Errorf(
				"Value of config option 'refresher.refreshInterval' (%d) is invalid, disabling auto-refresh",
				refreshInterval)
		}
	}

	if gui.Config.GetDebug() {
		gui.goEvery(time.Second*time.Duration(10), gui.stopChan, func() error {
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
			gui.c.Log.Infof("Heap memory in use: %s", formatBytes(m.HeapAlloc))
			return nil
		})
	}
}

func (gui *Gui) startBackgroundFetch() {
	gui.waitForIntro.Wait()
	isNew := gui.IsNewRepo
	userConfig := gui.UserConfig
	if !isNew {
		time.After(time.Duration(userConfig.Refresher.FetchInterval) * time.Second)
	}
	err := gui.backgroundFetch()
	if err != nil && strings.Contains(err.Error(), "exit status 128") && isNew {
		_ = gui.c.Alert(gui.c.Tr.NoAutomaticGitFetchTitle, gui.c.Tr.NoAutomaticGitFetchBody)
	} else {
		gui.goEvery(time.Second*time.Duration(userConfig.Refresher.FetchInterval), gui.stopChan, func() error {
			err := gui.backgroundFetch()
			gui.render()
			return err
		})
	}
}

func (gui *Gui) goEvery(interval time.Duration, stop chan struct{}, function func() error) {
	go utils.Safe(func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if gui.PauseBackgroundThreads {
					continue
				}
				_ = function()
			case <-stop:
				return
			}
		}
	})
}
