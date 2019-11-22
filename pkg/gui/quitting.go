package gui

import (
	"os"

	"github.com/jesseduffield/gocui"
)

// when a user runs lazygit with the LAZYGIT_NEW_DIR_FILE env variable defined
// we will write the current directory to that file on exit so that their
// shell can then change to that directory. That means you don't get kicked
// back to the directory that you started with.
func (gui *Gui) recordCurrentDirectory() error {
	if os.Getenv("LAZYGIT_NEW_DIR_FILE") == "" {
		return nil
	}

	// determine current directory, set it in LAZYGIT_NEW_DIR_FILE
	dirName, err := os.Getwd()
	if err != nil {
		return err
	}

	return gui.OSCommand.CreateFileWithContent(os.Getenv("LAZYGIT_NEW_DIR_FILE"), dirName)
}

func (gui *Gui) handleQuitWithoutChangingDirectory(g *gocui.Gui, v *gocui.View) error {
	gui.State.RetainOriginalDir = true
	return gui.quit(v)
}

func (gui *Gui) handleQuit(g *gocui.Gui, v *gocui.View) error {
	gui.State.RetainOriginalDir = false
	return gui.quit(v)
}

func (gui *Gui) quit(v *gocui.View) error {
	if gui.State.Updating {
		return gui.createUpdateQuitConfirmation(gui.g, v)
	}
	if gui.Config.GetUserConfig().GetBool("confirmOnQuit") {
		return gui.createConfirmationPanel(gui.g, v, true, "", gui.Tr.SLocalize("ConfirmQuit"), func(g *gocui.Gui, v *gocui.View) error {
			return gocui.ErrQuit
		}, nil)
	}

	return gocui.ErrQuit
}
