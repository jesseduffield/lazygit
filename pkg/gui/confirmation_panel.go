// lots of this has been directly ported from one of the example files, will brush up later

// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"strings"

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
		panic(err)
	}
	if err := gui.returnFocus(g, view); err != nil {
		panic(err)
	}
	g.DeleteKeybindings("confirmation")
	return g.DeleteView("confirmation")
}

func (gui *Gui) getMessageHeight(message string, width int) int {
	lines := strings.Split(message, "\n")
	lineCount := 0
	for _, line := range lines {
		lineCount += len(line)/width + 1
	}
	return lineCount
}

func (gui *Gui) getConfirmationPanelDimensions(g *gocui.Gui, prompt string) (int, int, int, int) {
	width, height := g.Size()
	panelWidth := width / 2
	panelHeight := gui.getMessageHeight(prompt, panelWidth)
	return width/2 - panelWidth/2,
		height/2 - panelHeight/2 - panelHeight%2 - 1,
		width/2 + panelWidth/2,
		height/2 + panelHeight/2
}

func (gui *Gui) createPromptPanel(g *gocui.Gui, currentView *gocui.View, title string, handleConfirm func(*gocui.Gui, *gocui.View) error) error {
	gui.onNewPopupPanel()
	confirmationView, err := gui.prepareConfirmationPanel(currentView, title, "")
	if err != nil {
		return err
	}
	confirmationView.Editable = true
	return gui.setKeyBindings(g, handleConfirm, nil)
}

func (gui *Gui) prepareConfirmationPanel(currentView *gocui.View, title, prompt string) (*gocui.View, error) {
	x0, y0, x1, y1 := gui.getConfirmationPanelDimensions(gui.g, prompt)
	confirmationView, err := gui.g.SetView("confirmation", x0, y0, x1, y1, 0)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return nil, err
		}
		confirmationView.Title = title
		confirmationView.FgColor = gocui.ColorWhite
	}
	confirmationView.Clear()

	if err := gui.switchFocus(gui.g, currentView, confirmationView); err != nil {
		return nil, err
	}
	return confirmationView, nil
}

func (gui *Gui) onNewPopupPanel() {
	gui.g.SetViewOnBottom("commitMessage")
}

func (gui *Gui) createConfirmationPanel(g *gocui.Gui, currentView *gocui.View, title, prompt string, handleConfirm, handleClose func(*gocui.Gui, *gocui.View) error) error {
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
		confirmationView, err := gui.prepareConfirmationPanel(currentView, title, prompt)
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
	return gui.createConfirmationPanel(g, currentView, title, prompt, nil, nil)
}

func (gui *Gui) createErrorPanel(g *gocui.Gui, message string) error {
	gui.Log.Error(message)
	currentView := g.CurrentView()
	colorFunction := color.New(color.FgRed).SprintFunc()
	coloredMessage := colorFunction(strings.TrimSpace(message))
	return gui.createConfirmationPanel(g, currentView, gui.Tr.SLocalize("Error"), coloredMessage, nil, nil)
}
