// lots of this has been directly ported from one of the example files, will brush up later

// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"strings"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type createPopupPanelOpts struct {
	hasLoader           bool
	editable            bool
	title               string
	prompt              string
	handleConfirm       func() error
	handleConfirmPrompt func(string) error
	handleClose         func() error

	// when handlersManageFocus is true, do not return from the confirmation context automatically. It's expected that the handlers will manage focus, whether that means switching to another context, or manually returning the context.
	handlersManageFocus bool

	findSuggestionsFunc func(string) []*types.Suggestion
}

type askOpts struct {
	title               string
	prompt              string
	handleConfirm       func() error
	handleClose         func() error
	handlersManageFocus bool
	findSuggestionsFunc func(string) []*types.Suggestion
}

type promptOpts struct {
	title               string
	initialContent      string
	handleConfirm       func(string) error
	findSuggestionsFunc func(string) []*types.Suggestion
}

func (gui *Gui) ask(opts askOpts) error {
	return gui.createPopupPanel(createPopupPanelOpts{
		title:               opts.title,
		prompt:              opts.prompt,
		handleConfirm:       opts.handleConfirm,
		handleClose:         opts.handleClose,
		handlersManageFocus: opts.handlersManageFocus,
		findSuggestionsFunc: opts.findSuggestionsFunc,
	})
}

func (gui *Gui) prompt(opts promptOpts) error {
	return gui.createPopupPanel(createPopupPanelOpts{
		title:               opts.title,
		prompt:              opts.initialContent,
		editable:            true,
		handleConfirmPrompt: opts.handleConfirm,
		findSuggestionsFunc: opts.findSuggestionsFunc,
	})
}

func (gui *Gui) createLoaderPanel(prompt string) error {
	return gui.createPopupPanel(createPopupPanelOpts{
		prompt:    prompt,
		hasLoader: true,
	})
}

func (gui *Gui) wrappedConfirmationFunction(handlersManageFocus bool, function func() error) func() error {
	return func() error {
		if function != nil {
			if err := function(); err != nil {
				return err
			}
		}

		if err := gui.closeConfirmationPrompt(handlersManageFocus); err != nil {
			return err
		}

		return nil
	}
}

func (gui *Gui) wrappedPromptConfirmationFunction(handlersManageFocus bool, function func(string) error, getResponse func() string) func() error {
	return func() error {
		if function != nil {
			if err := function(getResponse()); err != nil {
				return gui.surfaceError(err)
			}
		}

		if err := gui.closeConfirmationPrompt(handlersManageFocus); err != nil {
			return err
		}

		return nil
	}
}

func (gui *Gui) clearConfirmationViewKeyBindings() {
	keybindingConfig := gui.Config.GetUserConfig().Keybinding
	_ = gui.g.DeleteKeybinding("confirmation", gui.getKey(keybindingConfig.Universal.Confirm), gocui.ModNone)
	_ = gui.g.DeleteKeybinding("confirmation", gui.getKey(keybindingConfig.Universal.ConfirmAlt1), gocui.ModNone)
	_ = gui.g.DeleteKeybinding("confirmation", gui.getKey(keybindingConfig.Universal.Return), gocui.ModNone)
}

func (gui *Gui) closeConfirmationPrompt(handlersManageFocus bool) error {
	// we've already closed it so we can just return
	if !gui.Views.Confirmation.Visible {
		return nil
	}

	if !handlersManageFocus {
		if err := gui.returnFromContext(); err != nil {
			return err
		}
	}

	gui.clearConfirmationViewKeyBindings()
	gui.Views.Confirmation.Visible = false
	gui.Views.Suggestions.Visible = false

	return nil
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

func (gui *Gui) getConfirmationPanelDimensions(wrap bool, prompt string) (int, int, int, int) {
	width, height := gui.g.Size()
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

func (gui *Gui) prepareConfirmationPanel(title, prompt string, hasLoader bool, findSuggestionsFunc func(string) []*types.Suggestion) error {
	x0, y0, x1, y1 := gui.getConfirmationPanelDimensions(true, prompt)
	// calling SetView on an existing view returns the same view, so I'm not bothering
	// to reassign to gui.Views.Confirmation
	_, err := gui.g.SetView("confirmation", x0, y0, x1, y1, 0)
	if err != nil {
		return err
	}
	gui.Views.Confirmation.HasLoader = hasLoader
	if hasLoader {
		gui.g.StartTicking()
	}
	gui.Views.Confirmation.Title = title
	gui.Views.Confirmation.Wrap = true
	gui.Views.Confirmation.FgColor = theme.GocuiDefaultTextColor

	gui.findSuggestions = findSuggestionsFunc
	if findSuggestionsFunc != nil {
		suggestionsViewHeight := 11
		suggestionsView, err := gui.g.SetView("suggestions", x0, y1, x1, y1+suggestionsViewHeight, 0)
		if err != nil {
			return err
		}
		suggestionsView.Wrap = true
		suggestionsView.FgColor = theme.GocuiDefaultTextColor
		gui.setSuggestions([]*types.Suggestion{})
		suggestionsView.Visible = true
	}

	gui.g.Update(func(g *gocui.Gui) error {
		return gui.pushContext(gui.State.Contexts.Confirmation)
	})
	return nil
}

func (gui *Gui) createPopupPanel(opts createPopupPanelOpts) error {
	gui.g.Update(func(g *gocui.Gui) error {
		// remove any previous keybindings
		gui.clearConfirmationViewKeyBindings()

		err := gui.prepareConfirmationPanel(opts.title, opts.prompt, opts.hasLoader, opts.findSuggestionsFunc)
		if err != nil {
			return err
		}
		gui.Views.Confirmation.Editable = opts.editable
		gui.Views.Confirmation.Editor = gocui.EditorFunc(gui.defaultEditor)

		if opts.editable {
			if err := gui.Views.Confirmation.SetEditorContent(opts.prompt); err != nil {
				return err
			}
		} else {
			if err := gui.renderStringSync(gui.Views.Confirmation, opts.prompt); err != nil {
				return err
			}
		}

		return gui.setKeyBindings(opts)
	})
	return nil
}

func (gui *Gui) setKeyBindings(opts createPopupPanelOpts) error {
	actions := utils.ResolvePlaceholderString(
		gui.Tr.CloseConfirm,
		map[string]string{
			"keyBindClose":   "esc",
			"keyBindConfirm": "enter",
		},
	)

	gui.renderString(gui.Views.Options, actions)
	var onConfirm func() error
	if opts.handleConfirmPrompt != nil {
		onConfirm = gui.wrappedPromptConfirmationFunction(opts.handlersManageFocus, opts.handleConfirmPrompt, func() string { return gui.Views.Confirmation.Buffer() })
	} else {
		onConfirm = gui.wrappedConfirmationFunction(opts.handlersManageFocus, opts.handleConfirm)
	}

	type confirmationKeybinding struct {
		viewName string
		key      interface{}
		handler  func() error
	}

	keybindingConfig := gui.Config.GetUserConfig().Keybinding
	onSuggestionConfirm := gui.wrappedPromptConfirmationFunction(opts.handlersManageFocus, opts.handleConfirmPrompt, func() string { return gui.getSelectedSuggestionValue() })

	confirmationKeybindings := []confirmationKeybinding{
		{
			viewName: "confirmation",
			key:      gui.getKey(keybindingConfig.Universal.Confirm),
			handler:  onConfirm,
		},
		{
			viewName: "confirmation",
			key:      gui.getKey(keybindingConfig.Universal.ConfirmAlt1),
			handler:  onConfirm,
		},
		{
			viewName: "confirmation",
			key:      gui.getKey(keybindingConfig.Universal.Return),
			handler:  gui.wrappedConfirmationFunction(opts.handlersManageFocus, opts.handleClose),
		},
		{
			viewName: "confirmation",
			key:      gui.getKey(keybindingConfig.Universal.TogglePanel),
			handler:  func() error { return gui.replaceContext(gui.State.Contexts.Suggestions) },
		},
		{
			viewName: "suggestions",
			key:      gui.getKey(keybindingConfig.Universal.Confirm),
			handler:  onSuggestionConfirm,
		},
		{
			viewName: "suggestions",
			key:      gui.getKey(keybindingConfig.Universal.ConfirmAlt1),
			handler:  onSuggestionConfirm,
		},
		{
			viewName: "suggestions",
			key:      gui.getKey(keybindingConfig.Universal.Return),
			handler:  gui.wrappedConfirmationFunction(opts.handlersManageFocus, opts.handleClose),
		},
		{
			viewName: "suggestions",
			key:      gui.getKey(keybindingConfig.Universal.TogglePanel),
			handler:  func() error { return gui.replaceContext(gui.State.Contexts.Confirmation) },
		},
	}

	for _, binding := range confirmationKeybindings {
		if err := gui.g.SetKeybinding(binding.viewName, nil, binding.key, gocui.ModNone, gui.wrappedHandler(binding.handler)); err != nil {
			return err
		}
	}

	return nil
}

func (gui *Gui) wrappedHandler(f func() error) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		return f()
	}
}

func (gui *Gui) createErrorPanel(message string) error {
	colorFunction := color.New(color.FgRed).SprintFunc()
	coloredMessage := colorFunction(strings.TrimSpace(message))
	if err := gui.refreshSidePanels(refreshOptions{mode: ASYNC}); err != nil {
		return err
	}

	return gui.ask(askOpts{
		title:  gui.Tr.Error,
		prompt: coloredMessage,
	})
}

func (gui *Gui) surfaceError(err error) error {
	if err == nil {
		return nil
	}

	return gui.createErrorPanel(err.Error())
}
