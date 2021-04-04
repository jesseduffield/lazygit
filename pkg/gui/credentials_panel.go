package gui

import (
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type credentials chan string

// promptUserForCredential wait for a username, password or passphrase input from the credentials popup
func (gui *Gui) promptUserForCredential(passOrUname string) string {
	gui.credentials = make(chan string)
	gui.g.Update(func(g *gocui.Gui) error {
		credentialsView := gui.Views.Credentials
		switch passOrUname {
		case "username":
			credentialsView.Title = gui.Tr.CredentialsUsername
			credentialsView.Mask = 0
		case "password":
			credentialsView.Title = gui.Tr.CredentialsPassword
			credentialsView.Mask = '*'
		default:
			credentialsView.Title = gui.Tr.CredentialsPassphrase
			credentialsView.Mask = '*'
		}

		if err := gui.pushContext(gui.State.Contexts.Credentials); err != nil {
			return err
		}

		gui.RenderCommitLength()
		return nil
	})

	// wait for username/passwords/passphrase input
	userInput := <-gui.credentials
	return userInput + "\n"
}

func (gui *Gui) handleSubmitCredential() error {
	credentialsView := gui.Views.Credentials
	message := gui.trimmedContent(credentialsView)
	gui.credentials <- message
	gui.clearEditorView(credentialsView)
	if err := gui.returnFromContext(); err != nil {
		return err
	}

	return gui.refreshSidePanels(refreshOptions{})
}

func (gui *Gui) handleCloseCredentialsView() error {
	gui.credentials <- ""
	return gui.returnFromContext()
}

func (gui *Gui) handleCredentialsViewFocused() error {
	keybindingConfig := gui.Config.GetUserConfig().Keybinding

	message := utils.ResolvePlaceholderString(
		gui.Tr.CloseConfirm,
		map[string]string{
			"keyBindClose":   gui.getKeyDisplay(keybindingConfig.Universal.Return),
			"keyBindConfirm": gui.getKeyDisplay(keybindingConfig.Universal.Confirm),
		},
	)

	gui.renderString(gui.Views.Options, message)
	return nil
}

// handleCredentialsPopup handles the views after executing a command that might ask for credentials
func (gui *Gui) handleCredentialsPopup(cmdErr error) {
	if cmdErr != nil {
		errMessage := cmdErr.Error()
		if strings.Contains(errMessage, "Invalid username, password or passphrase") {
			errMessage = gui.Tr.PassUnameWrong
		}
		// we are not logging this error because it may contain a password or a passphrase
		_ = gui.createErrorPanel(errMessage)
	} else {
		_ = gui.closeConfirmationPrompt(false)
	}
}
