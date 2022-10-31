package helpers

import (
	"sync"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type CredentialsHelper struct {
	c *types.HelperCommon
}

func NewCredentialsHelper(
	c *types.HelperCommon,
) *CredentialsHelper {
	return &CredentialsHelper{
		c: c,
	}
}

// promptUserForCredential wait for a username, password or passphrase input from the credentials popup
func (self *CredentialsHelper) PromptUserForCredential(passOrUname oscommands.CredentialType) string {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)

	userInput := ""

	self.c.OnUIThread(func() error {
		title, mask := self.getTitleAndMask(passOrUname)

		return self.c.Prompt(types.PromptOpts{
			Title: title,
			Mask:  mask,
			HandleConfirm: func(input string) error {
				userInput = input

				waitGroup.Done()

				return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
			},
			HandleClose: func() error {
				waitGroup.Done()

				return nil
			},
		})
	})

	// wait for username/passwords/passphrase input
	waitGroup.Wait()

	return userInput + "\n"
}

func (self *CredentialsHelper) getTitleAndMask(passOrUname oscommands.CredentialType) (string, bool) {
	switch passOrUname {
	case oscommands.Username:
		return self.c.Tr.CredentialsUsername, false
	case oscommands.Password:
		return self.c.Tr.CredentialsPassword, true
	case oscommands.Passphrase:
		return self.c.Tr.CredentialsPassphrase, true
	case oscommands.PIN:
		return self.c.Tr.CredentialsPIN, true
	}

	// should never land here
	panic("unexpected credential request")
}
