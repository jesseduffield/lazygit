package helpers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type BisectHelper struct {
	c *HelperCommon
}

func NewBisectHelper(c *HelperCommon) *BisectHelper {
	return &BisectHelper{c: c}
}

func (self *BisectHelper) Reset() error {
	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.Bisect.ResetTitle,
		Prompt: self.c.Tr.Bisect.ResetPrompt,
		HandleConfirm: func() error {
			self.c.LogAction(self.c.Tr.Actions.ResetBisect)
			if err := self.c.Git().Bisect.Reset(); err != nil {
				return err
			}

			self.PostBisectCommandRefresh()
			return nil
		},
	})

	return nil
}

func (self *BisectHelper) PostBisectCommandRefresh() {
	self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{}})
}
