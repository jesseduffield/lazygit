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
	if sub != nil {
		return false, gui.runSubprocessWithSuspense(sub)
	}
	return true, nil
}

func (gui *Gui) handleCommitConfirm() error {
	commitMessageView := gui.getCommitMessageView()
	message := gui.trimmedContent(commitMessageView)
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
	if !ok {
		return nil
	}

	gui.clearEditorView(commitMessageView)
	_ = gui.returnFromContext()
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

	gui.renderString("options", message)
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
	v := gui.getCommitMessageView()
	v.Subtitle = gui.getBufferLength(v)
}
