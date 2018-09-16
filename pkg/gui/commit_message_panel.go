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

		err := gui.createErrorPanel(gui.Tr.SLocalize("CommitWithoutMessageErr"))
		if err != nil {
			gui.Log.Errorf("Failed to create error panel at handleCommitConfirm: %s\n", err)
			return err
		}

		return nil
	}
  
	sub, err := gui.GitCommand.Commit(message)
	if err != nil {
		if err != gui.Errors.ErrSubProcess {
      err = gui.createErrorPanel(err.Error())
			if err != nil {
				gui.Log.Errorf("Failed to create error panel at handleCommitConfirm: %s\n", err)
				return err
			}
			return nil
		}
	}

	if sub != nil {
		gui.SubProcess = sub
		return gui.Errors.ErrSubProcess
	}

	err = gui.refreshFiles()
	if err != nil {
		gui.Log.Errorf("Failed to refresh files at handleCommitConfirm: %s\n", err)
		return err
	}

	v.Clear()

	err = v.SetCursor(0, 0)
	if err != nil {
		gui.Log.Errorf("Failed to setcuror at handleCommitConfirm: %s\n", err)
		return err
	}

	_, err = gui.g.SetViewOnBottom("commitMessage")
	if err != nil {
		gui.Log.Errorf("Failed to set view to bottom at handleCommitConfirm: %s\n", err)
		return err
	}

	vv, err := gui.g.View("files")
	if vv == nil || err != nil {
		gui.Log.Errorf("Failed to get the files view at handleCommitConfirm: %s\n", err)
		return err
	}

	err = gui.switchFocus(v, vv)
	if err != nil {
		gui.Log.Errorf("Failed to switch focus at handleCommitConfirm: %s\n", err)
		return err
	}

	err = gui.refreshCommits()
	if err != nil {
		gui.Log.Errorf("Failed to refresh commits at handleCommitConfirm: %s\n", err)
		return err
	}

	return nil
}

// handleCommitClose gets called when a user presses exit on a commit message view.
// g and v are passes by the gocui library, but only v is used.
// If anything goes wrong it returns an error.
func (gui *Gui) handleCommitClose(g *gocui.Gui, v *gocui.View) error {
	_, err := gui.g.SetViewOnBottom("commitMessage")
	if err != nil {
		gui.Log.Errorf("Failed to set view on bottom at handleCommitClose: %s\n", err)
		return err
	}

	vv, err := gui.g.View("files")
	if vv == nil || err != nil {
		gui.Log.Errorf("Failed to get the files view at handleCommitClose: %s\n", err)
		return err
	}

	err = gui.switchFocus(v, vv)
	if err != nil {
		gui.Log.Errorf("Failed to switch focus at handleCommitClose: %s\n", err)
		return err
	}

	return nil
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

	err := gui.renderString("options", message)
	if err != nil {
		gui.Log.Errorf("Failed render string at handleCommitFocused: %s\n", err)
		return err
	}

	return nil
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
		gui.Log.Errorf("Failed render get commitMessage view at RenderCommitLength: %s\n", err)
		return
	}

	v.Subtitle = gui.getBufferLength(v)

}
