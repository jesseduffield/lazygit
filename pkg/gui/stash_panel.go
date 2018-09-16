package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
)

// refreshStashEntries refreshes the stash entries.
// returns an error if something goes wrong.
func (gui *Gui) refreshStashEntries() error {
	gui.g.Update(func(g *gocui.Gui) error {

		v, err := gui.g.View("stash")
		if err != nil {
			gui.Log.Errorf("Failed to get stash view at refreshStashEntries: %s\n", err)
		}

		gui.State.StashEntries = gui.GitCommand.GetStashEntries()

		v.Clear()

		// TODO use renderString or something
		for _, stashEntry := range gui.State.StashEntries {
			fmt.Fprintln(v, stashEntry.DisplayString)
		}

		err = gui.resetOrigin(v)
		if err != nil {
			gui.Log.Errorf("Failed to resetOrigin at ")
			return err
		}

		return nil
	})

	return nil
}

// getSelectedStashEntry returns the selected stash entry.
// returns a stash entry.
func (gui *Gui) getSelectedStashEntry() *commands.StashEntry {
	if len(gui.State.StashEntries) == 0 {
		return nil
	}

	stashView, _ := gui.g.View("stash")
	lineNumber := gui.getItemPosition(stashView)

	return &gui.State.StashEntries[lineNumber]
}

// handleStashEntrySelect is called when the user selects an item in the stash
// view.
// returns an error if something went wrong.
func (gui *Gui) handleStashEntrySelect() error {
	err := gui.renderGlobalOptions()
	if err != nil {
		gui.Log.Errorf("Failed to renderGlobalOptions at handleStashEntrySelect: %s\n", err)
		return err
	}

	go func() {
		stashEntry := gui.getSelectedStashEntry()
		if stashEntry == nil {

			err = gui.renderString("main", gui.Tr.SLocalize("NoStashEntries"))
			if err != nil {
				gui.Log.Errorf("Failed to renderString at handleStashEntrySelect: %s\n", err)
			}

			return
		}

		diff, _ := gui.GitCommand.GetStashEntryDiff(stashEntry.Index)

		err = gui.renderString("main", diff)
		if err != nil {
			gui.Log.Errorf("Failed to renderString at handleStashEntrySelect: %s\n", err)
		}
	}()

	return nil
}

// handleStashApply is a wrapper for the stashDo function sothat the
// gocui library can call it.
// g and v are passed by the gocui library.
// returns an error if something went wrong.
func (gui *Gui) handleStashApply(g *gocui.Gui, v *gocui.View) error {
	return gui.stashDo("apply")
}

// handleStashPop is a wrapper for the stashDo function sothat the
// gocui library can call it.
// g and v are passed by the gocui library.
// returns an error if something went wrong.
func (gui *Gui) handleStashPop(g *gocui.Gui, v *gocui.View) error {
	return gui.stashDo("pop")
}

// handleStashDrop is called when the user wants to drop the stash.
// g and v are passed by the gocui library.
// returns an error if something goes wrong.
func (gui *Gui) handleStashDrop(g *gocui.Gui, v *gocui.View) error {
	title := gui.Tr.SLocalize("StashDrop")
	message := gui.Tr.SLocalize("SureDropStashEntry")

	err := gui.createConfirmationPanel(v, title, message,
		func(g *gocui.Gui, v *gocui.View) error {
			return gui.stashDo("drop")
		}, nil)
	if err != nil {
		gui.Log.Errorf("Failed to createConfirmationPanel at handleStashDrop: %s\n", err)
		return err
	}

	return nil
}

// stashDo is the generic "actor" in the stash view.
// It actually performs the tasks that are requested.
// method: what to do.
// return an error if something goes wrong.
func (gui *Gui) stashDo(method string) error {
	stashEntry := gui.getSelectedStashEntry()

	if stashEntry == nil {
		errorMessage := gui.Tr.TemplateLocalize(
			"NoStashTo",
			Teml{
				"method": method,
			},
		)

		err := gui.createErrorPanel(errorMessage)
		if err != nil {
			gui.Log.Errorf("Failed to createErrorPanel at stashDo: %s\n", err)
			return err
		}

		return nil
	}

	err := gui.GitCommand.StashDo(stashEntry.Index, method)
	if err != nil {
		err = gui.createErrorPanel(err.Error())
		if err != nil {
			gui.Log.Errorf("Failed to createErrorPanel at statshDo: %s\n", err)
			return err
		}
	}

	err = gui.refreshStashEntries()
	if err != nil {
		gui.Log.Errorf("Failed to refreshStashEntries at stashDo: %s\n", err)
		return err
	}

	err = gui.refreshFiles()
	if err != nil {
		gui.Log.Errorf("Failed to refreshFiles at stashDo: %s\n", err)
		return err
	}

	return nil
}

// handleStashSave gets called when the user saves the stash.
// g and v are passed by the gocui library.
// returns an error if something goes wrong.
func (gui *Gui) handleStashSave(g *gocui.Gui, v *gocui.View) error {
	if len(gui.trackedFiles()) == 0 && len(gui.stagedFiles()) == 0 {
		err := gui.createErrorPanel(gui.Tr.SLocalize("NoTrackedStagedFilesStash"))
		if err != nil {
			gui.Log.Errorf("Failed to createErrorPanel at handleStashSave: %s\n", err)
			return err
		}
	}

	err := gui.createPromptPanel(v, gui.Tr.SLocalize("StashChanges"),
		func(g *gocui.Gui, v *gocui.View) error {

			err := gui.GitCommand.StashSave(gui.trimmedContent(v))
			if err != nil {
				err = gui.createErrorPanel(err.Error())
				if err != nil {
					gui.Log.Errorf("Failed to createErrorPanel at handleStashSave: %s\n", err)
					return err
				}
			}

			err = gui.refreshStashEntries()
			if err != nil {
				gui.Log.Errorf("Failed to refreshStashEntries at handleStashSave: %s\n", err)
				return err
			}

			err = gui.refreshFiles()
			if err != nil {
				gui.Log.Errorf("Failed to refreshFiles at handleStashSave: %s\n", err)
				return err
			}

			return nil
		})
	if err != nil {
		gui.Log.Errorf("Failed to createPromptPanel at handleStashSave: %s\n", err)
		return err
	}

	return nil
}
