package gui

import (
	"errors"
	"fmt"
)

var refreshables map[string]func() error

// refresh refreshes the gui
func (gui *Gui) refresh() error {

	for _, s := range refreshables {
		err := s()
		if err != nil {
			gui.Log.Error(err.Error())
		}
	}

	gui.refreshBranches(gui.g)
	gui.refreshFiles()
	gui.refreshCommits(gui.g)
	return nil
}

// registerRefresher adds a new refresh job to be run when refresh is
// called
func (gui *Gui) registerRefresher(name string, f func() error) error {

	if refreshables[name] != nil {
		err := errors.New(fmt.Sprintf("refresher %s already registered", name))
		gui.Log.Error(err.Error())
		return err
	}

	refreshables[name] = f
	return nil
}

// removeRefresher removes a refresh job
func (gui *Gui) removeRefresher(name string) error {

	if refreshables[name] == nil {
		err := errors.New(fmt.Sprintf("refresher %s not registered", name))
		gui.Log.Error(err.Error())
		return err
	}

	return nil
}
