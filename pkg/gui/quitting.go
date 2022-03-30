package gui

import (
	"os"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// when a user runs lazygit with the LAZYGIT_NEW_DIR_FILE env variable defined
// we will write the current directory to that file on exit so that their
// shell can then change to that directory. That means you don't get kicked
// back to the directory that you started with.
func (gui *Gui) recordCurrentDirectory() error {
	// determine current directory, set it in LAZYGIT_NEW_DIR_FILE
	dirName, err := os.Getwd()
	if err != nil {
		return err
	}
	return gui.recordDirectory(dirName)
}

func (gui *Gui) recordDirectory(dirName string) error {
	newDirFilePath := os.Getenv("LAZYGIT_NEW_DIR_FILE")
	if newDirFilePath == "" {
		return nil
	}
	return gui.os.CreateFileWithContent(newDirFilePath, dirName)
}

func (gui *Gui) handleQuitWithoutChangingDirectory() error {
	gui.RetainOriginalDir = true
	return gui.quit()
}

func (gui *Gui) handleQuit() error {
	gui.RetainOriginalDir = false
	return gui.quit()
}

func (gui *Gui) handleTopLevelReturn() error {
	currentContext := gui.currentContext()

	parentContext, hasParent := currentContext.GetParentContext()
	if hasParent && currentContext != nil && parentContext != nil {
		// TODO: think about whether this should be marked as a return rather than adding to the stack
		return gui.c.PushContext(parentContext)
	}

	for _, mode := range gui.modeStatuses() {
		if mode.isActive() {
			return mode.reset()
		}
	}

	repoPathStack := gui.RepoPathStack
	if !repoPathStack.IsEmpty() {
		path := repoPathStack.Pop()

		return gui.dispatchSwitchToRepo(path, true)
	}

	if gui.c.UserConfig.QuitOnTopLevelReturn {
		return gui.handleQuit()
	}

	return nil
}

func (gui *Gui) quit() error {
	if gui.State.Updating {
		return gui.createUpdateQuitConfirmation()
	}

	if gui.c.UserConfig.ConfirmOnQuit {
		return gui.c.Confirm(types.ConfirmOpts{
			Title:  "",
			Prompt: gui.c.Tr.ConfirmQuit,
			HandleConfirm: func() error {
				return gocui.ErrQuit
			},
		})
	}

	return gocui.ErrQuit
}
