package gui

import (
	"os"

	"github.com/jesseduffield/gocui"
	. "github.com/jesseduffield/lazygit/pkg/gui/types"
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

	return gui.OS.CreateFileWithContent(os.Getenv("LAZYGIT_NEW_DIR_FILE"), dirName)
}

func (gui *Gui) handleQuitWithoutChangingDirectory() error {
	gui.State.RetainOriginalDir = true
	return gui.quit()
}

func (gui *Gui) handleQuit() error {
	gui.State.RetainOriginalDir = false
	return gui.quit()
}

func (gui *Gui) handleTopLevelReturn() error {
	currentContext := gui.currentContext()

	parentContext, hasParent := currentContext.GetParentContext()
	if hasParent && currentContext != nil && parentContext != nil {
		// TODO: think about whether this should be marked as a return rather than adding to the stack
		return gui.pushContext(parentContext)
	}

	for _, mode := range gui.modeStatuses() {
		if mode.isActive() {
			return mode.reset()
		}
	}

	repoPathStack := gui.RepoPathStack
	if len(repoPathStack) > 0 {
		n := len(repoPathStack) - 1

		path := repoPathStack[n]

		gui.RepoPathStack = repoPathStack[:n]

		return gui.dispatchSwitchToRepo(path, true)
	}

	if gui.Config.GetUserConfig().QuitOnTopLevelReturn {
		return gui.handleQuit()
	}

	return nil
}

func (gui *Gui) quit() error {
	if gui.State.Updating {
		return gui.createUpdateQuitConfirmation()
	}

	if gui.Config.GetUserConfig().ConfirmOnQuit {
		return gui.Ask(AskOpts{
			Title:  "",
			Prompt: gui.Tr.ConfirmQuit,
			HandleConfirm: func() error {
				return gocui.ErrQuit
			},
		})
	}

	return gocui.ErrQuit
}
