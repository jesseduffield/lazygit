package gui

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// refreshStatus is called to refresh the status view.
// returns an error when something goes wrong.
func (gui *Gui) refreshStatus() error {
	v, err := gui.g.View("status")
	if err != nil {
		gui.Log.Errorf("Failed to get the status view at refreshStatus: %s\n", err)
		return err
	}

	// for some reason if this isn't wrapped in an update the clear seems to
	// be applied after the other things or something like that; the panel's
	// contents end up cleared
	gui.g.Update(func(*gocui.Gui) error {

		v.Clear()

		pushables, pullables := gui.GitCommand.UpstreamDifferenceCount()
		fmt.Fprint(v, "↑"+pushables+"↓"+pullables)
		branches := gui.State.Branches

		err := gui.updateHasMergeConflictStatus()
		if err != nil {
			gui.Log.Errorf("Failed to updateHasMergeConflictStatus at refreshStatus: %s\n", err)
			return err
		}

		if gui.State.HasMergeConflicts {
			fmt.Fprint(v, utils.ColoredString(" (merging)", color.FgYellow))
		}

		if len(branches) == 0 {
			return nil
		}

		branch := branches[0]
		name := utils.ColoredString(branch.Name, branch.GetColor())
		repo := utils.GetCurrentRepoName()
		fmt.Fprint(v, " "+repo+" → "+name)

		return nil
	})

	return nil
}

// handleCheckForUpdate is called when the user wants to check for updates.
// g and v are passed by the gocui library.
// returns an error when something goes wrong.
func (gui *Gui) handleCheckForUpdate(g *gocui.Gui, v *gocui.View) error {
	gui.Updater.CheckForNewUpdate(gui.onUserUpdateCheckFinish, true)

	err := gui.createMessagePanel(v, "", gui.Tr.SLocalize("CheckingForUpdates"))
	if err != nil {
		gui.Log.Errorf("Failed to createMessagePanel at handleCheckForUpdate: %s\n", err)
		return err
	}

	return nil
}

// handleStatusSelect is called when the status view is selected.
// returns an error when something goes wrong.
func (gui *Gui) handleStatusSelect() error {
	err := gui.renderString("main", dashboardString)
	if err != nil {
		gui.Log.Errorf("Failed to renderString at ")
		return err
	}

	err = gui.renderGlobalOptions()
	if err != nil {
		gui.Log.Errorf("Failed to renderGlobalOptions at handleStatusSelect: %s\n", err)
		return err
	}

	return nil
}

// handleOpenConfig is called when the user wants to open the config.
// g and v are passed by the gocui library.
// returns an error when something goes wrong.
func (gui *Gui) handleOpenConfig(g *gocui.Gui, v *gocui.View) error {
	return gui.openFile(gui.Config.GetUserConfig().ConfigFileUsed())
}

// handleEditCOnfig is called when the user wants to edit the config.
// g and v are passed by the gocui library.
// returns an error when something goes wrong.
func (gui *Gui) handleEditConfig(g *gocui.Gui, v *gocui.View) error {
	filename := gui.Config.GetUserConfig().ConfigFileUsed()
	return gui.editFile(filename)
}
