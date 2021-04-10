package gui

// Currently there is a bug where if we switch to a subprocess from within
// WithWaitingStatus we get stuck there and can't return to lazygit. We could
// fix this bug, or just stop running subprocesses from within there, given that
// we don't need to see a loading status if we're in a subprocess.
func (gui *Gui) withGpgHandling(cmdStr string, waitingStatus string, onSuccess func() error) error {
	useSubprocess := gui.GitCommand.UsingGpg()
	if useSubprocess {
		// Need to remember why we use the shell for the subprocess but not in the other case
		// Maybe there's no good reason
		success, err := gui.runSubprocessWithSuspense(gui.OSCommand.ShellCommandFromString(cmdStr))
		if success && onSuccess != nil {
			if err := onSuccess(); err != nil {
				return err
			}
		}
		if err := gui.refreshSidePanels(refreshOptions{mode: ASYNC}); err != nil {
			return err
		}

		if err != nil {
			return err
		}
	} else {
		return gui.WithWaitingStatus(waitingStatus, func() error {
			err := gui.OSCommand.RunCommand(cmdStr)
			if err != nil {
				return err
			} else if onSuccess != nil {
				if err := onSuccess(); err != nil {
					return err
				}
			}

			return gui.refreshSidePanels(refreshOptions{mode: ASYNC})
		})
	}

	return nil
}
