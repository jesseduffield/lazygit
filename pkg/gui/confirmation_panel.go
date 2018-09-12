package gui

import (
	"strings"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
)

func (gui *Gui) createConfirmationPanel(currentView *gocui.View, title, prompt string, handleConfirm, handleClose func(*gocui.Gui, *gocui.View) error) error {

	gui.g.SetViewOnBottom("commitMessage")

	gui.g.Update(func(g *gocui.Gui) error {

		if view, _ := gui.g.View("confirmation"); view != nil {

			err := gui.closeConfirmationPrompt()
			if err != nil {

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
			gui.Log.Errorf("Failed to prepare confirmation panel at createConfirmationPanel: %s\n", err)
			return err
		}

		confirmationView.Editable = false

		err = gui.renderString(gui.g, "confirmation", prompt)
		if err != nil {
			gui.Log.Errorf("Failed to render string at createConfirmationPanel: %s\n", err)
			return err
		}

		err = gui.setKeyBindings(handleConfirm, handleClose)
		if err != nil {
			gui.Log.Errorf("Failed to set keybindings at createConfirmationPanel: %s\n", err)
			return err
		}

		return nil
	})

	return nil
}

func (gui *Gui) createPromptPanel(currentView *gocui.View, title string, handleConfirm func(*gocui.Gui, *gocui.View) error) error {

	gui.g.SetViewOnBottom("commitMessage")

	confirmationView, err := gui.prepareConfirmationPanel(currentView, title, "")
	if err != nil {
		return err
	}

	confirmationView.Editable = true
	return gui.setKeyBindings(handleConfirm, nil)
}

func (gui *Gui) createErrorPanel(message string) error {

	gui.Log.Error(message)

	currentView := gui.g.CurrentView()

	colorFunction := color.New(color.FgRed).SprintFunc()

	coloredMessage := colorFunction(strings.TrimSpace(message))

	err := gui.createConfirmationPanel(currentView, gui.Tr.SLocalize("Error"), coloredMessage, nil, nil)
	if err != nil {

	}

	return nil
}

func (gui *Gui) createMessagePanel(currentView *gocui.View, title, prompt string) error {
	return gui.createConfirmationPanel(currentView, title, prompt, nil, nil)
}

func (gui *Gui) wrappedConfirmationFunction(function func(*gocui.Gui, *gocui.View) error) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		if function != nil {
			if err := function(gui.g, v); err != nil {
				return err
			}
		}
		return gui.closeConfirmationPrompt()
	}
}

func (gui *Gui) closeConfirmationPrompt() error {

	view, err := gui.g.View("confirmation")
	if err != nil {
		panic(err)
	}

	if err := gui.returnFocus(gui.g, view); err != nil {
		panic(err)
	}

	gui.g.DeleteKeybindings("confirmation")

	return gui.g.DeleteView("confirmation")
}

func (gui *Gui) getMessageHeight(message string, width int) int {
	lines := strings.Split(message, "\n")
	lineCount := 0
	for _, line := range lines {
		lineCount += len(line)/width + 1
	}
	return lineCount
}

func (gui *Gui) getConfirmationPanelDimensions(prompt string) (int, int, int, int) {
	width, height := gui.g.Size()
	panelWidth := width / 2
	panelHeight := gui.getMessageHeight(prompt, panelWidth)
	return width/2 - panelWidth/2,
		height/2 - panelHeight/2 - panelHeight%2 - 1,
		width/2 + panelWidth/2,
		height/2 + panelHeight/2
}

func (gui *Gui) prepareConfirmationPanel(currentView *gocui.View, title, prompt string) (*gocui.View, error) {

	x0, y0, x1, y1 := gui.getConfirmationPanelDimensions(prompt)

	confirmationView, err := gui.g.SetView("confirmation", x0, y0, x1, y1, 0)
	if err != nil {

		if err != gocui.ErrUnknownView {
			return nil, err
		}

		confirmationView.Title = title
		confirmationView.FgColor = gocui.ColorWhite
	}

	confirmationView.Clear()

	err = gui.switchFocus(gui.g, currentView, confirmationView)
	if err != nil {
		return nil, err
	}

	return confirmationView, nil
}

func (gui *Gui) setKeyBindings(handleConfirm, handleClose func(*gocui.Gui, *gocui.View) error) error {

	actions := gui.Tr.TemplateLocalize(
		"CloseConfirm",
		Teml{
			"keyBindClose":   "esc",
			"keyBindConfirm": "enter",
		},
	)

	err := gui.renderString(gui.g, "options", actions)
	if err != nil {
		return err
	}

	err = gui.g.SetKeybinding("confirmation", gocui.KeyEnter, gocui.ModNone, gui.wrappedConfirmationFunction(handleConfirm))
	if err != nil {
		return err
	}

	err = gui.g.SetKeybinding("confirmation", gocui.KeyEsc, gocui.ModNone, gui.wrappedConfirmationFunction(handleClose))
	if err != nil {
		return err
	}

	return nil
}
