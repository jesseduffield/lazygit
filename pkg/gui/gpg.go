package gui

import (
	. "github.com/jesseduffield/lazygit/pkg/commands/types"
	. "github.com/jesseduffield/lazygit/pkg/gui/types"
)

// Currently there is a bug where if we switch to a subprocess from within
// WithWaitingStatus we get stuck there and can't return to lazygit. We could
// fix this bug, or just stop running subprocesses from within there, given that
// we don't need to see a loading status if we're in a subprocess.
func (gui *Gui) withGpgHandling(cmdObj ICmdObj, waitingStatus string, onSuccess func() error) error {
	useSubprocess := gui.GitCommand.UsingGpg()
	if useSubprocess {
		success, err := gui.runSubprocessWithSuspense(cmdObj)
		if success && onSuccess != nil {
			if err := onSuccess(); err != nil {
				return err
			}
		}
		if err := gui.RefreshSidePanels(RefreshOptions{Mode: ASYNC}); err != nil {
			return err
		}

		if err != nil {
			return err
		}
	} else {
		return gui.WithWaitingStatus(waitingStatus, func() error {
			err := gui.OSCommand.RunExecutable(cmdObj)
			if err != nil {
				return err
			} else if onSuccess != nil {
				if err := onSuccess(); err != nil {
					return err
				}
			}

			return gui.RefreshSidePanels(RefreshOptions{Mode: ASYNC})
		})
	}

	return nil
}
