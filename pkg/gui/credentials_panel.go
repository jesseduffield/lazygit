package gui

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type credentials chan string

// promptUserForCredential wait for a username, password or passphrase input from the credentials popup
func (gui *Gui) promptUserForCredential(passOrUname oscommands.CredentialType) string {
	gui.credentials = make(chan string)
	gui.OnUIThread(func() error {
		credentialsView := gui.Views.Credentials
		switch passOrUname {
		case oscommands.Username:
			credentialsView.Title = gui.Tr.CredentialsUsername
			credentialsView.Mask = 0
		case oscommands.Password:
			credentialsView.Title = gui.Tr.CredentialsPassword
			credentialsView.Mask = '*'
		case oscommands.Passphrase:
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
	message := strings.TrimSpace(credentialsView.TextArea.GetContent())
	gui.credentials <- message
	credentialsView.ClearTextArea()
	if err := gui.returnFromContext(); err != nil {
		return err
	}

	return gui.refreshSidePanels(refreshOptions{mode: ASYNC})
}

func (gui *Gui) handleCloseCredentialsView() error {
	gui.credentials <- ""
	return gui.returnFromContext()
}

func (gui *Gui) handleCredentialsViewFocused() error {
	keybindingConfig := gui.UserConfig.Keybinding

	message := utils.ResolvePlaceholderString(
		gui.Tr.CloseConfirm,
		map[string]string{
			"keyBindClose":   gui.getKeyDisplay(keybindingConfig.Universal.Return),
			"keyBindConfirm": gui.getKeyDisplay(keybindingConfig.Universal.Confirm),
		},
	)

	return gui.renderString(gui.Views.Options, message)
}

// handleCredentialsPopup handles the views after executing a command that might ask for credentials
func (gui *Gui) handleCredentialsPopup(cmdErr error) {
	if cmdErr != nil {
		errMessage := cmdErr.Error()
		if strings.Contains(errMessage, "Invalid username, password or passphrase") {
			errMessage = gui.Tr.PassUnameWrong
		}
		_ = gui.returnFromContext()
		// we are not logging this error because it may contain a password or a passphrase
		_ = gui.PopupHandler.ErrorMsg(errMessage)
	} else {
		_ = gui.closeConfirmationPrompt(false)
	}
}
