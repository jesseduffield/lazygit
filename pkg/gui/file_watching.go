package gui

import (
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/jesseduffield/lazygit/pkg/commands"
)

// NOTE: given that we often edit files ourselves, this may make us end up refreshing files too often
// TODO: consider watching the whole directory recursively (could be more expensive)
func (gui *Gui) watchFilesForChanges() {
	var err error
	gui.fileWatcher, err = fsnotify.NewWatcher()
	if err != nil {
		gui.Log.Error(err)
		return
	}
	go func() {
		for {
			select {
			// watch for events
			case event := <-gui.fileWatcher.Events:
				if event.Op == fsnotify.Chmod {
					// for some reason we pick up chmod events when they don't actually happen
					continue
				}
				// only refresh if we're not already
				if !gui.State.IsRefreshingFiles {
					if err := gui.refreshFiles(); err != nil {
						err = gui.createErrorPanel(gui.g, err.Error())
						if err != nil {
							gui.Log.Error(err)
						}
					}
				}

			// watch for errors
			case err := <-gui.fileWatcher.Errors:
				if err != nil {
					gui.Log.Warn(err)
				}
			}
		}
	}()
}

func (gui *Gui) addFilesToFileWatcher(files []*commands.File) error {
	// watch the files for changes
	dirName, err := os.Getwd()
	if err != nil {
		return err
	}

	for _, file := range files {
		if err := gui.fileWatcher.Add(filepath.Join(dirName, file.Name)); err != nil {
			// swallowing errors here because it doesn't really matter if we can't watch a file
			gui.Log.Warn(err)
		}
	}

	return nil
}
