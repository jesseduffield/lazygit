package helpers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type AmendHelper struct {
	c   *HelperCommon
	gpg *GpgHelper
}

func NewAmendHelper(
	c *HelperCommon,
	gpg *GpgHelper,
) *AmendHelper {
	return &AmendHelper{
		c:   c,
		gpg: gpg,
	}
}

func (self *AmendHelper) AmendHead() error {
	return self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.AmendLastCommitTitle,
		Prompt: self.c.Tr.SureToAmend,
		HandleConfirm: func() error {
			cmdObj := self.c.Git().Commit.AmendHeadCmdObj()
			self.c.LogAction(self.c.Tr.Actions.AmendCommit)
			return self.gpg.WithGpgHandling(cmdObj, self.c.Tr.AmendingStatus, nil)
		},
	})
}
