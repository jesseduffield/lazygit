package gui

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
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

	return gui.refreshSidePanels(types.RefreshOptions{Mode: types.ASYNC})
}

func (gui *Gui) handleCloseCredentialsView() error {
	gui.Views.Credentials.ClearTextArea()
	gui.credentials <- ""
	return gui.returnFromContext()
}

func (gui *Gui) handleAskFocused() error {
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
