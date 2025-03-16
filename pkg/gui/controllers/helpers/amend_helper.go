package helpers

import "github.com/jesseduffield/lazygit/pkg/commands/git_commands"

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
	cmdObj := self.c.Git().Commit.AmendHeadCmdObj()
	self.c.LogAction(self.c.Tr.Actions.AmendCommit)
	return self.gpg.WithGpgHandling(cmdObj, git_commands.CommitGpgSign, self.c.Tr.AmendingStatus, nil, nil)
}
