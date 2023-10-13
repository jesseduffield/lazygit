package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type ScreenModeActions struct {
	c *ControllerCommon
}

func (self *ScreenModeActions) Next() error {
	self.c.State().GetRepoState().SetScreenMode(
		nextIntInCycle(
			[]types.WindowMaximisation{types.SCREEN_NORMAL, types.SCREEN_HALF, types.SCREEN_FULL},
			self.c.State().GetRepoState().GetScreenMode(),
		),
	)

	return nil
}

func (self *ScreenModeActions) Prev() error {
	self.c.State().GetRepoState().SetScreenMode(
		prevIntInCycle(
			[]types.WindowMaximisation{types.SCREEN_NORMAL, types.SCREEN_HALF, types.SCREEN_FULL},
			self.c.State().GetRepoState().GetScreenMode(),
		),
	)

	return nil
}

func nextIntInCycle(sl []types.WindowMaximisation, current types.WindowMaximisation) types.WindowMaximisation {
	for i, val := range sl {
		if val == current {
			if i == len(sl)-1 {
				return sl[0]
			}
			return sl[i+1]
		}
	}
	return sl[0]
}

func prevIntInCycle(sl []types.WindowMaximisation, current types.WindowMaximisation) types.WindowMaximisation {
	for i, val := range sl {
		if val == current {
			if i > 0 {
				return sl[i-1]
			}
			return sl[len(sl)-1]
		}
	}
	return sl[len(sl)-1]
}
