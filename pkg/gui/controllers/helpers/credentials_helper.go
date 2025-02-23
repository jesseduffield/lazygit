package helpers

import (
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type CredentialsHelper struct {
	c *HelperCommon
}

func NewCredentialsHelper(
	c *HelperCommon,
) *CredentialsHelper {
	return &CredentialsHelper{
		c: c,
	}
}

// promptUserForCredential wait for a username, password or passphrase input from the credentials popup
// We return a channel rather than returning the string directly so that the calling function knows
// when the prompt has been created (before the user has entered anything) so that it can
// note that we're now waiting on user input and lazygit isn't processing anything.
func (self *CredentialsHelper) PromptUserForInput(inputType oscommands.InputType) <-chan string {
	ch := make(chan string)

	self.c.OnUIThread(func() error {
		if inputType == oscommands.Ack {
			return self.c.Menu(types.CreateMenuOptions{
				Title: self.c.Tr.CommandLog + ": " + self.c.Tr.Actions.AckToContinue,
				Items: []*types.MenuItem{
					{
						Label: self.c.Tr.Yes,
						OnPress: func() error {
							ch <- "yes" + "\n"
							return nil
						},
					},
					{
						Label: self.c.Tr.No,
						OnPress: func() error {
							ch <- "no" + "\n"
							return nil
						},
					},
				},
			})
		}
		title, mask := self.getTitleAndMask(inputType)

		self.c.Prompt(types.PromptOpts{
			Title: title,
			Mask:  mask,
			HandleConfirm: func(input string) error {
				ch <- input + "\n"

				return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
			},
			HandleClose: func() error {
				ch <- "\n"

				return nil
			},
		})

		return nil
	})

	return ch
}

func (self *CredentialsHelper) getTitleAndMask(passOrUname oscommands.InputType) (string, bool) {
	switch passOrUname {
	case oscommands.Username:
		return self.c.Tr.CredentialsUsername, false
	case oscommands.Password:
		return self.c.Tr.CredentialsPassword, true
	case oscommands.Passphrase:
		return self.c.Tr.CredentialsPassphrase, true
	case oscommands.PIN:
		return self.c.Tr.CredentialsPIN, true
	case oscommands.Token:
		return self.c.Tr.CredentialsToken, true
	case oscommands.Ack:
		return self.c.Tr.Actions.AckToContinue, false
	}

	// should never land here
	panic("unexpected credential request")
}
