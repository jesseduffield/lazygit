package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type BisectHelper struct {
	c   *types.ControllerCommon
	git *commands.GitCommand
}

func NewBisectHelper(
	c *types.ControllerCommon,
	git *commands.GitCommand,
) *BisectHelper {
	return &BisectHelper{
		c:   c,
		git: git,
	}
}

func (self *BisectHelper) Reset() error {
	return self.c.Ask(types.AskOpts{
		Title:  self.c.Tr.Bisect.ResetTitle,
		Prompt: self.c.Tr.Bisect.ResetPrompt,
		HandleConfirm: func() error {
			self.c.LogAction(self.c.Tr.Actions.ResetBisect)
			if err := self.git.Bisect.Reset(); err != nil {
				return self.c.Error(err)
			}

			return self.PostBisectCommandRefresh()
		},
	})
}

func (self *BisectHelper) PostBisectCommandRefresh() error {
	return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{}})
}
