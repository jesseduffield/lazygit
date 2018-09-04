package gui

import (
	"strconv"
	"strings"

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

	err = gui.refreshFiles(g)
	if err != nil {

	}

	v.Clear()

	err = v.SetCursor(0, 0)
	if err != nil {
		return err
	}

	_, err = g.SetViewOnBottom("commitMessage")
	if err != nil {
		return err
	}

	err = gui.switchFocus(g, v, gui.getFilesView(g))
	if err != nil {
		return err
	}

	return gui.refreshCommits(g)

}

func (gui *Gui) handleCommitClose(g *gocui.Gui, v *gocui.View) error {

	_, err := g.SetViewOnBottom("commitMessage")
	if err != nil {
		return err
	}

	return gui.switchFocus(g, v, gui.getFilesView(g))

}

func (gui *Gui) handleNewlineCommitMessage(g *gocui.Gui, v *gocui.View) error {

	// resising ahead of time so that the top line doesn't get hidden to make
	// room for the cursor on the second line
	x0, y0, x1, y1 := gui.getConfirmationPanelDimensions(g, v.Buffer())

	_, err := g.SetView("commitMessage", x0, y0, x1, y1+1, 0)
	if err != nil {

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

func (gui *Gui) simpleEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	case key == gocui.KeyDelete:
		v.EditDelete(false)
	case key == gocui.KeyArrowDown:
		v.MoveCursor(0, 1, false)
	case key == gocui.KeyArrowUp:
		v.MoveCursor(0, -1, false)
	case key == gocui.KeyArrowLeft:
		v.MoveCursor(-1, 0, false)
	case key == gocui.KeyArrowRight:
		v.MoveCursor(1, 0, false)
	case key == gocui.KeyTab:
		v.EditNewLine()
	case key == gocui.KeySpace:
		v.EditWrite(' ')
	case key == gocui.KeyInsert:
		v.Overwrite = !v.Overwrite
	default:
		v.EditWrite(ch)
	}

	gui.RenderCommitLength()
}

func (gui *Gui) getBufferLength(view *gocui.View) string {
	return " " + strconv.Itoa(strings.Count(view.Buffer(), "")-1) + " "
}

func (gui *Gui) RenderCommitLength() {
	if !gui.Config.GetUserConfig().GetBool("gui.commitLength.show") {
		return
	}
	v := gui.getCommitMessageView(gui.g)
	v.Subtitle = gui.getBufferLength(v)
}
