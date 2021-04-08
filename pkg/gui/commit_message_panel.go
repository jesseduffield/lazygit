package gui

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// runSyncOrAsyncCommand takes the output of a command that may have returned
// either no error, an error, or a subprocess to execute, and if a subprocess
// needs to be run, it runs it
func (gui *Gui) runSyncOrAsyncCommand(sub *exec.Cmd, err error) (bool, error) {
	if err != nil {
		return false, gui.surfaceError(err)
	}
	if sub == nil {
		return true, nil
	}

	err = gui.runSubprocessWithSuspense(sub)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (gui *Gui) handleCommitConfirm() error {
	message := gui.trimmedContent(gui.Views.CommitMessage)
	if message == "" {
		return gui.createErrorPanel(gui.Tr.CommitWithoutMessageErr)
	}
	flags := ""
	skipHookPrefix := gui.Config.GetUserConfig().Git.SkipHookPrefix
	if skipHookPrefix != "" && strings.HasPrefix(message, skipHookPrefix) {
		flags = "--no-verify"
	}
	ok, err := gui.runSyncOrAsyncCommand(gui.GitCommand.Commit(message, flags))
	if err != nil {
		return err
	}

	_ = gui.returnFromContext()

	if !ok {
		return nil
	}

	gui.clearEditorView(gui.Views.CommitMessage)
	return gui.refreshSidePanels(refreshOptions{mode: ASYNC})
}

func (gui *Gui) handleCommitClose() error {
	return gui.returnFromContext()
}

func (gui *Gui) handleCommitMessageFocused() error {
	message := utils.ResolvePlaceholderString(
		gui.Tr.CommitMessageConfirm,
		map[string]string{
			"keyBindClose":   "esc",
			"keyBindConfirm": "enter",
			"keyBindNewLine": "tab",
		},
	)

	gui.renderString(gui.Views.Options, message)
	return nil
}

func (gui *Gui) getBufferLength(view *gocui.View) string {
	return " " + strconv.Itoa(strings.Count(view.Buffer(), "")-1) + " "
}

// RenderCommitLength is a function.
func (gui *Gui) RenderCommitLength() {
	if !gui.Config.GetUserConfig().Gui.CommitLength.Show {
		return
	}

	gui.Views.CommitMessage.Subtitle = gui.getBufferLength(gui.Views.CommitMessage)
}
