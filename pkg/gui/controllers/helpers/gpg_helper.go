package helpers

import (
	"fmt"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type GpgHelper struct {
	c *HelperCommon
}

func NewGpgHelper(c *HelperCommon) *GpgHelper {
	return &GpgHelper{
		c: c,
	}
}

// Currently there is a bug where if we switch to a subprocess from within
// WithWaitingStatus we get stuck there and can't return to lazygit. We could
// fix this bug, or just stop running subprocesses from within there, given that
// we don't need to see a loading status if we're in a subprocess.
func (self *GpgHelper) WithGpgHandling(cmdObj oscommands.ICmdObj, waitingStatus string, onSuccess func() error) error {
	useSubprocess := self.c.Git().Config.UsingGpg()
	if useSubprocess {
		success, err := self.c.RunSubprocess(cmdObj)
		if success && onSuccess != nil {
			if err := onSuccess(); err != nil {
				return err
			}
		}
		if err := self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC}); err != nil {
			return err
		}

		return err
	} else {
		return self.runAndStream(cmdObj, waitingStatus, onSuccess)
	}
}

func (self *GpgHelper) runAndStream(cmdObj oscommands.ICmdObj, waitingStatus string, onSuccess func() error) error {
	return self.c.WithWaitingStatus(waitingStatus, func(gocui.Task) error {
		if err := cmdObj.StreamOutput().Run(); err != nil {
			_ = self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
			return fmt.Errorf(
				self.c.Tr.GitCommandFailed, self.c.UserConfig.Keybinding.Universal.ExtrasMenu,
			)
		}

		if onSuccess != nil {
			if err := onSuccess(); err != nil {
				return err
			}
		}

		return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
	})
}
