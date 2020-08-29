package gui

import (
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
)

// getCredentialsHandler returns a handler for the git credentials
func (gui *Gui) getCredentialsHandler(title string) *commands.AuthInput {
	gui.credentialsInput = make(chan string)
	gui.credentialsInputOpen = true
	res := &commands.AuthInput{
		StdIn:   gui.credentialsInput,
		Updates: make(chan commands.AuthUpdate),
		Open:    true,
	}

	go func(auth *commands.AuthInput) {
		for {
			update, open := <-res.Updates
			if !open {
				break
			}

			gui.g.Update(func(g *gocui.Gui) error {
				v, _ := g.View("credentials")

				question := update.MightBeQuestion
				if question == nil {
					v.Title = title
					v.Editable = false
					v.HasLoader = true
				} else {
					if !v.Editable {
						v.Editable = true
						v.HasLoader = false
						gui.clearEditorView(v)
					}
					v.Title = *question
				}

				if update.MaskInput {
					v.Mask = '*'
				} else {
					v.Mask = 0
				}

				if err := gui.switchContext(gui.Contexts.Credentials.Context); err != nil {
					return err
				}

				gui.RenderCommitLength()
				return nil
			})
		}

		gui.credentialsInputOpen = false
		_, _ = gui.g.SetViewOnBottom("credentials")
		_ = gui.returnFromContext()
	}(res)
	return res
}

func (gui *Gui) handleSubmitCredential(g *gocui.Gui, v *gocui.View) error {
	message := gui.trimmedContent(v)
	if gui.credentialsInputOpen {
		gui.credentialsInput <- message
	}
	gui.clearEditorView(v)
	if err := gui.returnFromContext(); err != nil {
		return err
	}

	return gui.refreshSidePanels(refreshOptions{})
}

func (gui *Gui) handleCloseCredentialsView(g *gocui.Gui, v *gocui.View) error {
	if gui.credentialsInputOpen {
		gui.credentialsInput <- ""
	}

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
