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
// needs to be set on the gui object, it does so, and then returns the error
// the bool returned tells us whether the calling code should continue
func (gui *Gui) runSyncOrAsyncCommand(sub *exec.Cmd, err error) (bool, error) {
	if err != nil {
		if err != gui.Errors.ErrSubProcess {
			return false, gui.surfaceError(err)
		}
	}
	if sub != nil {
		gui.SubProcess = sub
		return false, gui.Errors.ErrSubProcess
	}
	return true, nil
}

func (gui *Gui) handleCommitConfirm(g *gocui.Gui, v *gocui.View) error {
	message := gui.trimmedContent(v)
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

	gui.clearEditorView(v)
	_ = gui.returnFromContext()
	return gui.refreshSidePanels(refreshOptions{mode: ASYNC})
}

func (gui *Gui) handleCommitClose(g *gocui.Gui, v *gocui.View) error {
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
