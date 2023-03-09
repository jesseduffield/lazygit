package helpers

import (
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type AmendHelper struct {
	c   *types.HelperCommon
	git *commands.GitCommand
	gpg *GpgHelper
}

func NewAmendHelper(
	c *types.HelperCommon,
	git *commands.GitCommand,
	gpg *GpgHelper,
) *AmendHelper {
	return &AmendHelper{
		c:   c,
		git: git,
		gpg: gpg,
	}
}

func (self *AmendHelper) AmendHead() error {
	return self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.AmendLastCommitTitle,
		Prompt: self.c.Tr.SureToAmend,
		HandleConfirm: func() error {
			cmdObj := self.git.Commit.AmendHeadCmdObj()
			self.c.LogAction(self.c.Tr.Actions.AmendCommit)
			return self.gpg.WithGpgHandling(cmdObj, self.c.Tr.AmendingStatus, nil)
		},
	})
}
