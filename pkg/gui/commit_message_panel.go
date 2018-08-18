package gui

import (
	"github.com/jesseduffield/gocui"
)

func (gui *Gui) handleCommitConfirm(g *gocui.Gui, v *gocui.View) error {
	message := gui.trimmedContent(v)
	if message == "" {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("CommitWithoutMessageErr"))
	}
	sub, err := gui.GitCommand.Commit(g, message)
	if err != nil {
		// TODO need to find a way to send through this error
		if err != gui.Errors.ErrSubProcess {
			return gui.createErrorPanel(g, err.Error())
		}
	}
	if sub != nil {
		gui.SubProcess = sub
		return gui.Errors.ErrSubProcess
	}
	gui.refreshFiles(g)
	v.Clear()
	v.SetCursor(0, 0)
	g.SetViewOnBottom("commitMessage")
	gui.switchFocus(g, v, gui.getFilesView(g))
	return gui.refreshCommits(g)
}

func (gui *Gui) handleCommitClose(g *gocui.Gui, v *gocui.View) error {
	g.SetViewOnBottom("commitMessage")
	return gui.switchFocus(g, v, gui.getFilesView(g))
}

func (gui *Gui) handleNewlineCommitMessage(g *gocui.Gui, v *gocui.View) error {
	// resising ahead of time so that the top line doesn't get hidden to make
	// room for the cursor on the second line
	x0, y0, x1, y1 := gui.getConfirmationPanelDimensions(g, v.Buffer())
	if _, err := g.SetView("commitMessage", x0, y0, x1, y1+1, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}

	v.EditNewLine()
	return nil
}

func (gui *Gui) handleCommitFocused(g *gocui.Gui, v *gocui.View) error {
	message := gui.Tr.TemplateLocalize(
		"CloseConfirm",
		Teml{
			"keyBindClose":   "esc",
			"keyBindConfirm": "enter",
		},
	)
	return gui.renderString(g, "options", message)
}
