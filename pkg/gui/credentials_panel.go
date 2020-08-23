package gui

import (
	"strings"

	"github.com/jesseduffield/gocui"
)

type credentials chan string

// promptUserForCredential wait for a username or password input from the credentials popup
func (gui *Gui) promptUserForCredential(passOrUname string) string {
	gui.credentials = make(chan string)
	gui.g.Update(func(g *gocui.Gui) error {
		credentialsView, _ := g.View("credentials")
		if passOrUname == "username" {
			credentialsView.Title = gui.Tr.SLocalize("CredentialsUsername")
			credentialsView.Mask = 0
		} else {
			credentialsView.Title = gui.Tr.SLocalize("CredentialsPassword")
			credentialsView.Mask = '*'
		}

		if err := gui.switchContext(gui.Contexts.Credentials.Context); err != nil {
			return err
		}

		gui.RenderCommitLength()
		return nil
	})

	// wait for username/passwords input
	userInput := <-gui.credentials
	return userInput + "\n"
}

func (gui *Gui) handleSubmitCredential(g *gocui.Gui, v *gocui.View) error {
	message := gui.trimmedContent(v)
	gui.credentials <- message
	gui.clearEditorView(v)
	if err := gui.returnFromContext(); err != nil {
		return err
	}

	return gui.refreshSidePanels(refreshOptions{})
}

func (gui *Gui) handleCloseCredentialsView(g *gocui.Gui, v *gocui.View) error {
	gui.credentials <- ""
	return gui.returnFromContext()
}

func (gui *Gui) handleCredentialsViewFocused() error {
	message := gui.Tr.TemplateLocalize(
		"CloseConfirm",
		Teml{
			"keyBindClose":   gui.getKeyDisplay("universal.return"),
			"keyBindConfirm": gui.getKeyDisplay("universal.confirm"),
		},
	)
	gui.renderString("options", message)
	return nil
}

// handleCredentialsPopup handles the views after executing a command that might ask for credentials
func (gui *Gui) handleCredentialsPopup(cmdErr error) {
	if cmdErr != nil {
		errMessage := cmdErr.Error()
		if strings.Contains(errMessage, "Invalid username or password") {
			errMessage = gui.Tr.SLocalize("PassUnameWrong")
		}
		// we are not logging this error because it may contain a password
		gui.createErrorPanel(errMessage)
	} else {
		_ = gui.closeConfirmationPrompt(false)
	}
}
