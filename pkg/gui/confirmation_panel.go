package gui

import (
	"strings"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
)

type handler func(gui *gocui.Gui, view *gocui.View) error

// createConfirmationPanel creates a new confirmation panel.
// currentView: is used to return the focus.
// title: sets the title of the panel.
// prompt: what to ask the user.
// handleConfirm: what to call on confirmation.
// handleClose: what to call on close.
// returns an error if something goes wrong.
func (gui *Gui) createConfirmationPanel(currentView *gocui.View, title, prompt string, handleConfirm, handleClose handler) error {
	_, err := gui.g.SetViewOnBottom("commitMessage")
	if err != nil {
		gui.Log.Errorf("Failed to set view on bottom at createConfirmationPanel: %s\n", err)
		return err
	}

	gui.g.Update(func(g *gocui.Gui) error {

		view, _ := gui.g.View("confirmation")
		if view != nil {

			err := gui.closeConfirmationPrompt()
			if err != nil {

				errMessage := gui.Tr.TemplateLocalize(
					"CantCloseConfirmationPrompt",
					Teml{
						"error": err.Error(),
					},
				)

				gui.Log.Errorf("%s\n", errMessage)
				return err
			}
		}

		confirmationView, err := gui.prepareConfirmationPanel(currentView, title, prompt)
		if err != nil {
			gui.Log.Errorf("Failed to prepare confirmation panel at createConfirmationPanel: %s\n", err)
			return err
		}

		confirmationView.Editable = false

		err = gui.renderString("confirmation", prompt)
		if err != nil {
			gui.Log.Errorf("Failed to render string at createConfirmationPanel: %s\n", err)
			return err
		}

		err = gui.setConfirmationHandlers(handleConfirm, handleClose)
		if err != nil {
			gui.Log.Errorf("Failed to set keybindings at createConfirmationPanel: %s\n", err)
			return err
		}

		return nil
	})

	return nil
}

// createPromptPanel creates a prompt panel which is just a confirmation panel
// without any handlers.
// currentView: is used to return focus.
// title: what to call the view.
// handleConfirm: what to do on confirmation.
// returns an error if something goes wrong.
func (gui *Gui) createPromptPanel(currentView *gocui.View, title string, handleConfirm handler) error {
	_, err := gui.g.SetViewOnBottom("commitMessage")
	if err != nil {
		gui.Log.Errorf("Failed to set view on bottom at createPromptPanel: %s\n", err)
		return err
	}
	confirmationView, err := gui.prepareConfirmationPanel(currentView, title, "")
	if err != nil {
		gui.Log.Errorf("Failed to prepare confirmation panel at createPromptPanel: %s\n", err)
		return err
	}

	confirmationView.Editable = true

	err = gui.setConfirmationHandlers(handleConfirm, nil)
	if err != nil {
		gui.Log.Errorf("Failed to setConfirmationHandlers at createPromptPanel: %s\n", err)
		return err
	}

	return nil
}

// createErrorPanel creates a error panel which is just a confirmation panel.
// message: the error to display.
// returns an error if something goes wrong.
func (gui *Gui) createErrorPanel(message string) error {
	gui.Log.Error(message)

	currentView := gui.g.CurrentView()

	colorFunction := color.New(color.FgRed).SprintFunc()
	coloredMessage := colorFunction(strings.TrimSpace(message))

	err := gui.createConfirmationPanel(currentView, gui.Tr.SLocalize("Error"), coloredMessage, nil, nil)
	if err != nil {
		gui.Log.Errorf("Failed to create confirmation panel at createErrorPanel: %s\n", err)
		return err
	}

	return nil
}

// createMessagePanel creates a message panel which is just a confirmation panel.
// without any handlers.
// currentView: is used to return focus.
// title: what to call the view.
// prompt: what to display.
// returns an error if something goes wrong.
func (gui *Gui) createMessagePanel(currentView *gocui.View, title, prompt string) error {
	return gui.createConfirmationPanel(currentView, title, prompt, nil, nil)
}

// wrappedConfirmationFunction creates a function that is acceptable for some other functions.
// TO BE REMOVED.
func (gui *Gui) wrappedConfirmationFunction(function handler) handler {
	// TODO change this
	return func(g *gocui.Gui, v *gocui.View) error {
		if function != nil {
			if err := function(gui.g, v); err != nil {
				return err
			}
		}
		return gui.closeConfirmationPrompt()
	}

}

// closeConfirmationPrompt closes the confirmation prompt.
// returns an error if something goes wrong.
func (gui *Gui) closeConfirmationPrompt() error {
	view, err := gui.g.View("confirmation")
	if err != nil {
		gui.Log.Errorf("Failed to get confirmation view at closeConfirmationPrompt: %s\n", err)
		return err
	}

	err = gui.returnFocus(view)
	if err != nil {
		gui.Log.Errorf("Failed to return focus at closeConfirmationPrompt: %s\n", err)
		return err
	}

	gui.g.DeleteKeybindings("confirmation")

	err = gui.g.DeleteView("confirmation")
	if err != nil {
		gui.Log.Errorf("Failed to delete confirmation view at closeConfirmationPrompt: %s\n", err)
		return err
	}

	return nil
}

// getMessageHeight returns the height of the message
// message: what to check.
// width: the width.
// returns the size.
func (gui *Gui) getMessageHeight(message string, width int) int {
	lines := strings.Split(message, "\n")
	lineCount := 0
	for _, line := range lines {
		lineCount += len(line)/width + 1
	}

	return lineCount
}

// getConfirmationPanelDimensions calculates the dimensions of the panel.
// prompt: what to prompt the user.
// returns: begin x, begin y, offset width, offset height.
func (gui *Gui) getConfirmationPanelDimensions(prompt string) (int, int, int, int) {
	width, height := gui.g.Size()
	panelWidth := width / 2
	panelHeight := gui.getMessageHeight(prompt, panelWidth)

	return width/2 - panelWidth/2,
		height/2 - panelHeight/2 - panelHeight%2 - 1,
		width/2 + panelWidth/2,
		height/2 + panelHeight/2
}

// prepareConfirmationPanel prepares a confirmation panel... who would have thought...
// currentView: used to return focus.
// title: the title of the panel.
// prompt: what to prompt the user.
// returns the view and if any occurred an error
func (gui *Gui) prepareConfirmationPanel(currentView *gocui.View, title, prompt string) (*gocui.View, error) {
	x0, y0, x1, y1 := gui.getConfirmationPanelDimensions(prompt)
	confirmationView, err := gui.g.SetView("confirmation", x0, y0, x1, y1, 0)
	if err != nil {
		if err != gocui.ErrUnknownView {
			gui.Log.Errorf("Failed to set view at prepareConfirmationPanel: %s\n", err)
			return nil, err
		}

		confirmationView.Title = title
		confirmationView.FgColor = gocui.ColorWhite
	}

	confirmationView.Clear()

	err = gui.switchFocus(currentView, confirmationView)
	if err != nil {
		gui.Log.Errorf("Failed to switch focus at prepareConfirmationPanel: %s\n", err)
		return nil, err
	}

	return confirmationView, nil
}

// setConfirmationHandlers makes sure that the keybindings for the panel is set.
// handleConfirm: what to do when the user presses confirm.
// handleClose: what to do when the user presses close.
func (gui *Gui) setConfirmationHandlers(handleConfirm, handleClose handler) error {
	actions := gui.Tr.TemplateLocalize(
		"CloseConfirm",
		Teml{
			"keyBindClose":   "esc",
			"keyBindConfirm": "enter",
		},
	)

	err := gui.renderString("options", actions)
	if err != nil {
		gui.Log.Errorf("Failed to renderString at setConfirmationHandlers: %s\n", err)
		return err
	}

	err = gui.g.SetKeybinding("confirmation", gocui.KeyEnter, gocui.ModNone, gui.wrappedConfirmationFunction(handleConfirm))
	if err != nil {
		gui.Log.Errorf("Failed to setConfirmationHandlers at closeConfirmationPrompt: %s\n", err)
		return err
	}

	err = gui.g.SetKeybinding("confirmation", gocui.KeyEsc, gocui.ModNone, gui.wrappedConfirmationFunction(handleClose))
	if err != nil {
		gui.Log.Errorf("Failed to setConfirmationHandlers at closeConfirmationPrompt: %s\n", err)
		return err
	}

	return nil
}
