// lots of this has been directly ported from one of the example files, will brush up later

// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
)

func (gui *Gui) wrappedConfirmationFunction(function func(*gocui.Gui, *gocui.View) error) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		if function != nil {
			if err := function(g, v); err != nil {
				return err
			}
		}
		return gui.closeConfirmationPrompt(g)
	}
}

func (gui *Gui) closeConfirmationPrompt(g *gocui.Gui) error {
	view, err := g.View("confirmation")
	if err != nil {
		return nil // if it's already been closed we can just return
	}
	if err := gui.returnFocus(g, view); err != nil {
		panic(err)
	}
	g.DeleteKeybindings("confirmation")
	return g.DeleteView("confirmation")
}

func (gui *Gui) getMessageHeight(wrap bool, message string, width int) int {
	lines := strings.Split(message, "\n")
	lineCount := 0
	// if we need to wrap, calculate height to fit content within view's width
	if wrap {
		for _, line := range lines {
			lineCount += len(line)/width + 1
		}
	} else {
		lineCount = len(lines)
	}
	return lineCount
}

func (gui *Gui) getConfirmationPanelDimensions(g *gocui.Gui, wrap bool, prompt string) (int, int, int, int) {
	width, height := g.Size()
	panelWidth := width / 2
	panelHeight := gui.getMessageHeight(wrap, prompt, panelWidth)
	return width/2 - panelWidth/2,
		height/2 - panelHeight/2 - panelHeight%2 - 1,
		width/2 + panelWidth/2,
		height/2 + panelHeight/2
}

func (gui *Gui) createPromptPanel(g *gocui.Gui, currentView *gocui.View, title string, handleConfirm func(*gocui.Gui, *gocui.View) error) error {
	gui.onNewPopupPanel()
	confirmationView, err := gui.prepareConfirmationPanel(currentView, title, "", false)
	if err != nil {
		return err
	}
	confirmationView.Editable = true
	return gui.setKeyBindings(g, handleConfirm, nil)
}

func (gui *Gui) prepareConfirmationPanel(currentView *gocui.View, title, prompt string, hasLoader bool) (*gocui.View, error) {
	x0, y0, x1, y1 := gui.getConfirmationPanelDimensions(gui.g, true, prompt)
	confirmationView, err := gui.g.SetView("confirmation", x0, y0, x1, y1, 0)
	if err != nil {
		if err.Error() != "unknown view" {
			return nil, err
		}
		confirmationView.HasLoader = hasLoader
		confirmationView.Title = title
		confirmationView.Wrap = true
		confirmationView.FgColor = gocui.ColorWhite
	}
	gui.g.Update(func(g *gocui.Gui) error {
		return gui.switchFocus(gui.g, currentView, confirmationView)
	})
	return confirmationView, nil
}

func (gui *Gui) onNewPopupPanel() {
	viewNames := []string{"commitMessage",
		"credentials",
		"menu"}
	for _, viewName := range viewNames {
		_, _ = gui.g.SetViewOnBottom(viewName)
	}
}

func (gui *Gui) createLoaderPanel(g *gocui.Gui, currentView *gocui.View, prompt string) error {
	return gui.createPopupPanel(g, currentView, "", prompt, true, nil, nil)
}

// it is very important that within this function we never include the original prompt in any error messages, because it may contain e.g. a user password
func (gui *Gui) createConfirmationPanel(g *gocui.Gui, currentView *gocui.View, title, prompt string, handleConfirm, handleClose func(*gocui.Gui, *gocui.View) error) error {
	return gui.createPopupPanel(g, currentView, title, prompt, false, handleConfirm, handleClose)
}

func (gui *Gui) createPopupPanel(g *gocui.Gui, currentView *gocui.View, title, prompt string, hasLoader bool, handleConfirm, handleClose func(*gocui.Gui, *gocui.View) error) error {
	gui.onNewPopupPanel()
	g.Update(func(g *gocui.Gui) error {
		// delete the existing confirmation panel if it exists
		if view, _ := g.View("confirmation"); view != nil {
			if err := gui.closeConfirmationPrompt(g); err != nil {
				errMessage := gui.Tr.TemplateLocalize(
					"CantCloseConfirmationPrompt",
					Teml{
						"error": err.Error(),
					},
				)
				gui.Log.Error(errMessage)
			}
		}
		confirmationView, err := gui.prepareConfirmationPanel(currentView, title, prompt, hasLoader)
		if err != nil {
			return err
		}
		confirmationView.Editable = false
		if err := gui.renderString(g, "confirmation", prompt); err != nil {
			return err
		}
		return gui.setKeyBindings(g, handleConfirm, handleClose)
	})
	return nil
}

func (gui *Gui) setKeyBindings(g *gocui.Gui, handleConfirm, handleClose func(*gocui.Gui, *gocui.View) error) error {
	actions := gui.Tr.TemplateLocalize(
		"CloseConfirm",
		Teml{
			"keyBindClose":   "esc",
			"keyBindConfirm": "enter",
		},
	)
	if err := gui.renderString(g, "options", actions); err != nil {
		return err
	}
	if err := g.SetKeybinding("confirmation", gocui.KeyEnter, gocui.ModNone, gui.wrappedConfirmationFunction(handleConfirm)); err != nil {
		return err
	}
	return g.SetKeybinding("confirmation", gocui.KeyEsc, gocui.ModNone, gui.wrappedConfirmationFunction(handleClose))
}

func (gui *Gui) createMessagePanel(g *gocui.Gui, currentView *gocui.View, title, prompt string) error {
	return gui.createPopupPanel(g, currentView, title, prompt, false, nil, nil)
}

// createSpecificErrorPanel allows you to create an error popup, specifying the
//  view to be focused when the user closes the popup, and a boolean specifying
// whether we will log the error. If the message may include a user password,
// this function is to be used over the more generic createErrorPanel, with
// willLog set to false
func (gui *Gui) createSpecificErrorPanel(message string, nextView *gocui.View, willLog bool) error {
	if willLog {
		go func() {
			// when reporting is switched on this log call sometimes introduces
			// a delay on the error panel popping up. Here I'm adding a second wait
			// so that the error is logged while the user is reading the error message
			time.Sleep(time.Second)
			gui.Log.Error(message)
		}()
	}

	colorFunction := color.New(color.FgRed).SprintFunc()
	coloredMessage := colorFunction(strings.TrimSpace(message))
	return gui.createConfirmationPanel(gui.g, nextView, gui.Tr.SLocalize("Error"), coloredMessage, nil, nil)
}

func (gui *Gui) createErrorPanel(g *gocui.Gui, message string) error {
	return gui.createSpecificErrorPanel(message, g.CurrentView(), true)
}
