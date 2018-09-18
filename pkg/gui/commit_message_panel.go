package gui

import (
	"strconv"
	"strings"

	"github.com/jesseduffield/gocui"
)

// handleCommitConfirm gets called when a user presses enter on
// a commit confirmation panel.
// g and v are passed by the gocui.
// If anything goes wrong it returns an error.
func (gui *Gui) handleCommitConfirm(g *gocui.Gui, v *gocui.View) error {
	message := gui.trimmedContent(v)
	if message == "" {
		return gui.createErrorPanel(gui.Tr.SLocalize("CommitWithoutMessageErr"))
	}

	sub, err := gui.GitCommand.Commit(message)
	if err != nil {
		if err != gui.Errors.ErrSubProcess {
			return gui.createErrorPanel(err.Error())
		}
	}

	if sub != nil {
		gui.SubProcess = sub
		return gui.Errors.ErrSubProcess
	}

	if err := gui.refreshFiles(); err != nil {
		return err
	}

	v.Clear()

	if err = v.SetCursor(0, 0); err != nil {
		return err
	}

	if _, err = gui.g.SetViewOnBottom("commitMessage"); err != nil {
		return err
	}

	filesView, err := gui.g.View("files")
	if filesView == nil || err != nil {
		return err
	}

	if err := gui.switchFocus(v, filesView); err != nil {
		return err
	}

	return gui.refreshCommits()
}

// handleCommitClose gets called when a user presses exit on a commit message view.
// g and v are passes by the gocui library, but only v is used.
// If anything goes wrong it returns an error.
func (gui *Gui) handleCommitClose(g *gocui.Gui, v *gocui.View) error {
	if _, err := gui.g.SetViewOnBottom("commitMessage"); err != nil {
		return err
	}

	filesView, err := gui.g.View("files")
	if filesView == nil || err != nil {
		return err
	}

	return gui.switchFocus(v, filesView)
}

// handleCommitFocused gets called when the commitMessageView is called.
// If anything goes wrong it returns an error.
func (gui *Gui) handleCommitFocused() error {
	message := gui.Tr.TemplateLocalize(
		"CloseConfirm",
		Teml{
			"keyBindClose":   "esc",
			"keyBindConfirm": "enter",
		},
	)

	return gui.renderString("options", message)
}

// simpleEditor is a simple implementation to provide custom key handling for
// the gocui library.
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

// getBufferLength calculates the size of the buffer.
// Takes the view to check.
// returns the count as a string.
func (gui *Gui) getBufferLength(view *gocui.View) string {
	return " " + strconv.Itoa(strings.Count(view.Buffer(), "")-1) + " "
}

// RenderCommitLength renders the commit length
func (gui *Gui) RenderCommitLength() {
	if !gui.Config.GetUserConfig().GetBool("gui.commitLength.show") {
		return
	}

	v, err := gui.g.View("commitMessage")
	if err != nil {
		return
	}

	v.Subtitle = gui.getBufferLength(v)
}
