package gui

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/jesseduffield/gocui"
)

// runSyncOrAsyncCommand takes the output of a command that may have returned
// either no error, an error, or a subprocess to execute, and if a subprocess
// needs to be set on the gui object, it does so, and then returns the error
// the bool returned tells us whether the calling code should continue
func (gui *Gui) runSyncOrAsyncCommand(sub *exec.Cmd, err error) (bool, error) {
	if err != nil {
		if err != gui.Errors.ErrSubProcess {
			return false, gui.createErrorPanel(gui.g, err.Error())
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
		return gui.createErrorPanel(g, gui.Tr.SLocalize("CommitWithoutMessageErr"))
	}
	flags := ""
	skipHookPrefix := gui.Config.GetUserConfig().GetString("git.skipHookPrefix")
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

	v.Clear()
	_ = v.SetCursor(0, 0)
	_ = v.SetOrigin(0, 0)
	_, _ = g.SetViewOnBottom("commitMessage")
	_ = gui.switchFocus(g, v, gui.getFilesView())
	return gui.refreshSidePanels(g)
}

func (gui *Gui) handleCommitClose(g *gocui.Gui, v *gocui.View) error {
	g.SetViewOnBottom("commitMessage")
	return gui.switchFocus(g, v, gui.getFilesView())
}

func (gui *Gui) handleCommitFocused(g *gocui.Gui, v *gocui.View) error {
	if _, err := g.SetViewOnTop("commitMessage"); err != nil {
		return err
	}

	message := gui.Tr.TemplateLocalize(
		"CloseConfirm",
		Teml{
			"keyBindClose":   "esc",
			"keyBindConfirm": "enter",
		},
	)
	return gui.renderString(g, "options", message)
}

func (gui *Gui) getBufferLength(view *gocui.View) string {
	return " " + strconv.Itoa(strings.Count(view.Buffer(), "")-1) + " "
}

// RenderCommitLength is a function.
func (gui *Gui) RenderCommitLength() {
	if !gui.Config.GetUserConfig().GetBool("gui.commitLength.show") {
		return
	}
	v := gui.getCommitMessageView()
	v.Subtitle = gui.getBufferLength(v)
}
