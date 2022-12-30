package gui

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) handleQuitWithoutChangingDirectory() error {
	gui.RetainOriginalDir = true
	return gui.quit()
}

func (gui *Gui) handleQuit() error {
	gui.RetainOriginalDir = false
	return gui.quit()
}

func (gui *Gui) handleTopLevelReturn() error {
	currentContext := gui.c.CurrentContext()

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

	repoPathStack := gui.c.State().GetRepoPathStack()
	if !repoPathStack.IsEmpty() {
		return gui.helpers.Repos.DispatchSwitchToRepo(repoPathStack.Pop(), true)
	}

	if gui.c.UserConfig.QuitOnTopLevelReturn {
		return gui.handleQuit()
	}

	return nil
}

func (gui *Gui) quit() error {
	if gui.c.State().GetUpdating() {
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

func (gui *Gui) createUpdateQuitConfirmation() error {
	return gui.c.Confirm(types.ConfirmOpts{
		Title:  gui.Tr.ConfirmQuitDuringUpdateTitle,
		Prompt: gui.Tr.ConfirmQuitDuringUpdate,
		HandleConfirm: func() error {
			return gocui.ErrQuit
		},
	})
}
