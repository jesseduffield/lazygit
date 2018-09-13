package gui

import (
	"errors"
	"fmt"
)

var refreshables = make(map[string]func() error)

// refresh refreshes the gui
func (gui *Gui) refresh() {

	for n, s := range refreshables {
		err := s()
		if err != nil {
			gui.Log.Errorf("Failed to refresh %s at refresh: %s\n", n, err)
		}
	}

	return
}

// registerRefresher adds a new refresh job to be run when refresh is
// called.
// name: what to call the refresher.
// f: what function to register
func (gui *Gui) registerRefresher(name string, f func() error) error {

	if refreshables[name] != nil {
		err := errors.New(fmt.Sprintf("refresher %s already registered", name))
		gui.Log.Errorf("Failed to register refresher: %s\n", err)
		return err
	}

	refreshables[name] = f

	return nil
}

// removeRefresher removes a refresh job.
// name: what to remove.
// returns an error if something goes wrong.
func (gui *Gui) removeRefresher(name string) error {

	if refreshables[name] == nil {
		err := errors.New(fmt.Sprintf("refresher %s not registered", name))
		gui.Log.Errorf("Failed to remove refresher: %s\n", err)
		return err
	}

	return nil
}
