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
	"github.com/jesseduffield/lazygit/pkg/theme"
)

type createPopupPanelOpts struct {
	returnToView        *gocui.View
	hasLoader           bool
	returnFocusOnClose  bool
	editable            bool
	title               string
	prompt              string
	handleConfirm       func() error
	handleConfirmPrompt func(string) error
	handleClose         func() error
}

type askOpts struct {
	returnToView       *gocui.View
	returnFocusOnClose bool
	title              string
	prompt             string
	handleConfirm      func() error
	handleClose        func() error
}

func (gui *Gui) createLoaderPanel(currentView *gocui.View, prompt string) error {
	return gui.createPopupPanel(createPopupPanelOpts{
		returnToView:       currentView,
		prompt:             prompt,
		hasLoader:          true,
		returnFocusOnClose: true,
	})
}

func (gui *Gui) ask(opts askOpts) error {
	return gui.createPopupPanel(createPopupPanelOpts{
		returnToView:       opts.returnToView,
		title:              opts.title,
		prompt:             opts.prompt,
		returnFocusOnClose: opts.returnFocusOnClose,
		handleConfirm:      opts.handleConfirm,
		handleClose:        opts.handleClose,
	})
}

func (gui *Gui) prompt(currentView *gocui.View, title string, initialContent string, handleConfirm func(string) error) error {
	return gui.createPopupPanel(createPopupPanelOpts{
		returnToView:        currentView,
		title:               title,
		prompt:              initialContent,
		returnFocusOnClose:  true,
		editable:            true,
		handleConfirmPrompt: handleConfirm,
	})
}

func (gui *Gui) wrappedConfirmationFunction(function func() error, returnFocusOnClose bool) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {

		if function != nil {
			if err := function(); err != nil {
				return err
			}
		}

		return gui.closeConfirmationPrompt(g, returnFocusOnClose)
	}
}

func (gui *Gui) wrappedPromptConfirmationFunction(function func(string) error, returnFocusOnClose bool) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {

		if function != nil {
			if err := function(v.Buffer()); err != nil {
				return err
			}
		}

		return gui.closeConfirmationPrompt(g, returnFocusOnClose)
	}
}

func (gui *Gui) closeConfirmationPrompt(g *gocui.Gui, returnFocusOnClose bool) error {
	view, err := g.View("confirmation")
	if err != nil {
		return nil // if it's already been closed we can just return
	}
	view.Editable = false
	if returnFocusOnClose {
		if err := gui.returnFocus(g, view); err != nil {
			panic(err)
		}
	}
	g.DeleteKeybinding("confirmation", gocui.KeyEnter, gocui.ModNone)
	g.DeleteKeybinding("confirmation", gocui.KeyEsc, gocui.ModNone)
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
	// we want a minimum width up to a point, then we do it based on ratio.
	panelWidth := 4 * width / 7
	minWidth := 80
	if panelWidth < minWidth {
		if width-2 < minWidth {
			panelWidth = width - 2
		} else {
			panelWidth = minWidth
		}
	}
	panelHeight := gui.getMessageHeight(wrap, prompt, panelWidth)
	if panelHeight > height*3/4 {
		panelHeight = height * 3 / 4
	}
	return width/2 - panelWidth/2,
		height/2 - panelHeight/2 - panelHeight%2 - 1,
		width/2 + panelWidth/2,
		height/2 + panelHeight/2
}

func (gui *Gui) prepareConfirmationPanel(currentView *gocui.View, title, prompt string, hasLoader bool) (*gocui.View, error) {
	x0, y0, x1, y1 := gui.getConfirmationPanelDimensions(gui.g, true, prompt)
	confirmationView, err := gui.g.SetView("confirmation", x0, y0, x1, y1, 0)
	if err != nil {
		if err.Error() != "unknown view" {
			return nil, err
		}
		confirmationView.HasLoader = hasLoader
		if hasLoader {
			gui.g.StartTicking()
		}
		confirmationView.Title = title
		confirmationView.Wrap = true
		confirmationView.FgColor = theme.GocuiDefaultTextColor
	}
	gui.g.Update(func(g *gocui.Gui) error {
		return gui.switchFocus(currentView, confirmationView)
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

func (gui *Gui) createPopupPanel(opts createPopupPanelOpts) error {
	gui.onNewPopupPanel()
	gui.g.Update(func(g *gocui.Gui) error {
		// delete the existing confirmation panel if it exists
		if view, _ := g.View("confirmation"); view != nil {
			if err := gui.closeConfirmationPrompt(g, true); err != nil {
				gui.Log.Error(err)
			}
		}
		confirmationView, err := gui.prepareConfirmationPanel(opts.returnToView, opts.title, opts.prompt, opts.hasLoader)
		if err != nil {
			return err
		}
		confirmationView.Editable = opts.editable
		if opts.editable {
			go func() {
				// TODO: remove this wait (right now if you remove it the EditGotoToEndOfLine method doesn't seem to work)
				time.Sleep(time.Millisecond)
				gui.g.Update(func(g *gocui.Gui) error {
					confirmationView.EditGotoToEndOfLine()
					return nil
				})
			}()
		}

		gui.renderString("confirmation", opts.prompt)
		return gui.setKeyBindings(opts)
	})
	return nil
}

func (gui *Gui) setKeyBindings(opts createPopupPanelOpts) error {
	actions := gui.Tr.TemplateLocalize(
		"CloseConfirm",
		Teml{
			"keyBindClose":   "esc",
			"keyBindConfirm": "enter",
		},
	)

	gui.renderString("options", actions)
	if opts.handleConfirmPrompt != nil {
		if err := gui.g.SetKeybinding("confirmation", nil, gocui.KeyEnter, gocui.ModNone, gui.wrappedPromptConfirmationFunction(opts.handleConfirmPrompt, opts.returnFocusOnClose)); err != nil {
			return err
		}
	} else {
		if err := gui.g.SetKeybinding("confirmation", nil, gocui.KeyEnter, gocui.ModNone, gui.wrappedConfirmationFunction(opts.handleConfirm, opts.returnFocusOnClose)); err != nil {
			return err
		}
	}

	return gui.g.SetKeybinding("confirmation", nil, gocui.KeyEsc, gocui.ModNone, gui.wrappedConfirmationFunction(opts.handleClose, opts.returnFocusOnClose))
}

func (gui *Gui) createErrorPanel(message string) error {
	colorFunction := color.New(color.FgRed).SprintFunc()
	coloredMessage := colorFunction(strings.TrimSpace(message))
	if err := gui.refreshSidePanels(refreshOptions{mode: ASYNC}); err != nil {
		return err
	}

	return gui.ask(askOpts{
		returnToView:       gui.g.CurrentView(),
		title:              gui.Tr.SLocalize("Error"),
		prompt:             coloredMessage,
		returnFocusOnClose: true,
	})
}

func (gui *Gui) surfaceError(err error) error {
	return gui.createErrorPanel(err.Error())
}
