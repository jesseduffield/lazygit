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
	sub, err := gui.GitCommand.Commit(message, false)
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

type credentials chan string

// waitForPassUname wait for a username or password input from the pushPassUname popup
func (gui *Gui) waitForPassUname(g *gocui.Gui, currentView *gocui.View, passOrUname string) string {
	gui.credentials = make(chan string)
	pushPassUnameView, _ := g.View("pushPassUname")
	if passOrUname == "username" {
		pushPassUnameView.Title = gui.Tr.SLocalize("PushUsername")
		pushPassUnameView.Mask = 0
	} else {
		pushPassUnameView.Title = gui.Tr.SLocalize("PushPassword")
		pushPassUnameView.Mask = '*'
	}
	g.Update(func(g *gocui.Gui) error {
		_, err := g.SetViewOnTop("pushPassUname")
		if err != nil {
			return err
		}
		err = gui.switchFocus(g, currentView, pushPassUnameView)
		if err != nil {
			return err
		}
		gui.RenderCommitLength()
		return nil
	})

	// wait for username/passwords input
	userInput := <-gui.credentials
	return userInput
}

func (gui *Gui) handlePushConfirm(g *gocui.Gui, v *gocui.View) error {
	message := gui.trimmedContent(v)
	if message == "" {
		// make sure to input something
		// if not dune the push progress will run forever
		message = "-"
	}
	gui.credentials <- message
	err := gui.refreshFiles(g)
	if err != nil {
		return err
	}
	v.Clear()
	err = v.SetCursor(0, 0)
	if err != nil {
		return err
	}
	_, err = g.SetViewOnBottom("pushPassUname")
	if err != nil {
		return err
	}
	err = gui.switchFocus(g, v, gui.getFilesView(g))
	if err != nil {
		return err
	}
	return gui.refreshCommits(g)
}

func (gui *Gui) handlePushClose(g *gocui.Gui, v *gocui.View) error {
	_, err := g.SetViewOnBottom("pushPassUname")
	if err != nil {
		return err
	}
	gui.credentials <- "-"
	return gui.switchFocus(g, v, gui.getFilesView(g))
}

func (gui *Gui) handlePushFocused(g *gocui.Gui, v *gocui.View) error {
	if _, err := g.SetViewOnTop("pushPassUname"); err != nil {
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

// RenderCommitLength is a function.
func (gui *Gui) RenderCommitLength() {
	if !gui.Config.GetUserConfig().GetBool("gui.commitLength.show") {
		return
	}
	v := gui.getCommitMessageView(gui.g)
	v.Subtitle = gui.getBufferLength(v)
}
