package gui

func (gui *Gui) inFilterMode() bool {
	return gui.State.Modes.Filtering.Path != ""
}

func (gui *Gui) validateNotInFilterMode() (bool, error) {
	if gui.inFilterMode() {
		err := gui.ask(askOpts{
			returnToView:       gui.g.CurrentView(),
			returnFocusOnClose: true,
			title:              gui.Tr.SLocalize("MustExitFilterModeTitle"),
			prompt:             gui.Tr.SLocalize("MustExitFilterModePrompt"),
			handleConfirm: func() error {
				return gui.exitFilterMode()
			},
		})

		return false, err
	}
	return true, nil
}

func (gui *Gui) exitFilterMode() error {
	gui.State.Modes.Filtering.Path = ""
	return gui.Errors.ErrRestart
}
