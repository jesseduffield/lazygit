package helpers

import (
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type BranchesHelper struct {
	c *HelperCommon
}

func NewBranchesHelper(c *HelperCommon) *BranchesHelper {
	return &BranchesHelper{
		c: c,
	}
}

func (self *BranchesHelper) ConfirmDeleteRemote(remoteName string, branchName string) error {
	title := utils.ResolvePlaceholderString(
		self.c.Tr.DeleteBranchTitle,
		map[string]string{
			"selectedBranchName": branchName,
		},
	)
	prompt := utils.ResolvePlaceholderString(
		self.c.Tr.DeleteRemoteBranchPrompt,
		map[string]string{
			"selectedBranchName": branchName,
			"upstream":           remoteName,
		},
	)
	return self.c.Confirm(types.ConfirmOpts{
		Title:  title,
		Prompt: prompt,
		HandleConfirm: func() error {
			return self.c.WithWaitingStatus(self.c.Tr.DeletingStatus, func(task gocui.Task) error {
				self.c.LogAction(self.c.Tr.Actions.DeleteRemoteBranch)
				if err := self.c.Git().Remote.DeleteRemoteBranch(task, remoteName, branchName); err != nil {
					return err
				}
				return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.BRANCHES, types.REMOTES}})
			})
		},
	})
}

func ShortBranchName(fullBranchName string) string {
	return strings.TrimPrefix(strings.TrimPrefix(fullBranchName, "refs/heads/"), "refs/remotes/")
}
