package gui

import (
	"strings"

	"github.com/jesseduffield/gocui"
	. "github.com/jesseduffield/lazygit/pkg/commands/types"
	. "github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type credentials chan string

// PromptUserForCredential wait for a username, password or passphrase input from the credentials popup
func (gui *Gui) PromptUserForCredential(credentialKind CredentialKind) string {
	gui.credentials = make(chan string)
	gui.g.Update(func(g *gocui.Gui) error {
		credentialsView := gui.Views.Credentials
		switch credentialKind {
		case USERNAME:
			credentialsView.Title = gui.Tr.CredentialsUsername
			credentialsView.Mask = 0
		case PASSWORD:
			credentialsView.Title = gui.Tr.CredentialsPassword
			credentialsView.Mask = '*'
		case PASSPHRASE:
			credentialsView.Title = gui.Tr.CredentialsPassphrase
			credentialsView.Mask = '*'
		default:
			panic("unknown credential requested")
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

	return gui.RefreshSidePanels(RefreshOptions{Mode: ASYNC})
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

// InformOnCredentialsOutcome handles the views after executing a command that might ask for credentials
func (gui *Gui) InformOnCredentialsOutcome(cmdErr error) {
	if cmdErr != nil {
		errMessage := cmdErr.Error()
		if strings.Contains(errMessage, "Invalid username, password or passphrase") {
			errMessage = gui.Tr.PassUnameWrong
		}
		_ = gui.returnFromContext()
		// we are not logging this error because it may contain a password or a passphrase
		_ = gui.CreateErrorPanel(errMessage)
	} else {
		_ = gui.closeConfirmationPrompt(false)
	}
}
